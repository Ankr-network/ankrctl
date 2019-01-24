/*
Copyright 2018 The Dccncli Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"

	akrctl "github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"

	"context"

	ankr_const "github.com/Ankr-network/dccn-common"
	pb "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	"google.golang.org/grpc"
)

var port = ":" + strconv.Itoa(ankr_const.DefaultPort)

var clientURL string

// Task creates the task command.
func Task() *Command {
	//DCCN-CLI task
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "task",
			Aliases: []string{"t"},
			Short:   "task commands",
			Long:    "task is used to access task commands",
		},
		DocCategories: []string{"task"},
		IsIndex:       true,
	}

	//DCCN-CLI comput task create
	cmdTaskCreate := CmdBuilder(cmd, RunTaskCreate, "create <task-name> [task-name ...]", "create task", Writer,
		aliasOpt("cr"), displayerType(&displayers.Task{}), docCategories("task"))
	AddStringFlag(cmdTaskCreate, akrctl.ArgDcidSlug, "", "", "Task dc-id")
	AddStringFlag(cmdTaskCreate, akrctl.ArgTypeSlug, "", "", "Task type")
	AddStringFlag(cmdTaskCreate, akrctl.ArgReplicaSlug, "", "", "Task replica")

	//DCCN-CLI comput task delete
	cmdRunTaskDelete := CmdBuilder(cmd, RunTaskDelete, "delete <task-id> [task-id ...]", "Delete task by id", Writer,
		aliasOpt("dl", "del", "rm"), docCategories("Task"))
	AddBoolFlag(cmdRunTaskDelete, akrctl.ArgForce, akrctl.ArgShortForce, false, "Force task delete")

	//DCCN-CLI comput task update
	cmdRunTaskUpdate := CmdBuilder(cmd, RunTaskUpdate, "update <task-id> [task-id ...]", "Update task by id", Writer,
		aliasOpt("ud", "udt", "ch"), docCategories("Task"))
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgNameSlug, "", "", "Task name")
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgReplicaSlug, "", "", "Task replica")

	//DCCN-CLI task list
	cmdRunTaskList := CmdBuilder(cmd, RunTaskList, "list [GLOB]", "list tasks", Writer,
		aliasOpt("ls"), displayerType(&displayers.Task{}), docCategories("task"))
	_ = cmdRunTaskList

	//DCCN-CLI task detail
	cmdRunTaskDetail := CmdBuilder(cmd, RunTaskDetail, "detail <task-id>", "list tasks detail", Writer,
		aliasOpt("dt"), displayerType(&displayers.Task{}), docCategories("task"))
	_ = cmdRunTaskDetail

	//DCCN-CLI comput task purge
	cmdRunTaskPurge := CmdBuilder(cmd, RunTaskPurge, "purge <task-id> [task-id ...]", "Purge task by id", Writer,
		aliasOpt("pg"), docCategories("Task"))
	AddBoolFlag(cmdRunTaskPurge, akrctl.ArgForce, akrctl.ArgShortForce, false, "Force task purge")

	return cmd

}

// RunTaskCreate creates a task.
//DCCN-CLI comput task create
func RunTaskCreate(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	userid, err := c.Ankr.GetInt(c.NS, akrctl.ArgUserID)
	if err != nil {
		return err
	}

	taskdcid, err := c.Ankr.GetString(c.NS, akrctl.ArgDcidSlug)
	if err != nil {
		return err
	}

	tasktype, err := c.Ankr.GetString(c.NS, akrctl.ArgTypeSlug)
	if err != nil {
		return err
	}

	replica, err := c.Ankr.GetString(c.NS, akrctl.ArgReplicaSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	dc := pb.NewTaskMgrClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	errs := make(chan error, len(c.Args))
	for _, name := range c.Args {
		tcrq := &pb.AddTaskRequest{
			UserId: int64(userid),
			Task: &common_proto.Task{
				Name: name,
				Type: tasktype,
			},
		}
		if taskdcid != "" {
			dcid, err := strconv.Atoi(taskdcid)
			if err != nil {
				return fmt.Errorf("dc id %s is not an int", taskdcid)
			}
			tcrq.Task.DataCenterId = int64(dcid)
		}
		if replica != "" {
			replicaCount, err := strconv.Atoi(replica)
			if err != nil {
				return fmt.Errorf("replica count %s is not an int", replica)
			}
			tcrq.Task.Replica = int32(replicaCount)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			tcrp, err := dc.CreateTask(ctx, tcrq)
			if err != nil {
				errs <- err
				return
			}
			if tcrp.Error != nil {
				fmt.Printf("Fail to initialize task, %s.\n", tcrp.Error.Details)
			} else {
				fmt.Printf("Task id %d initialize successfully. \n", tcrp.TaskId)
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func allInt(in []string) ([]int, error) {
	out := []int{}
	seen := map[string]bool{}

	for _, i := range in {
		if seen[i] {
			continue
		}

		seen[i] = true

		id, err := strconv.Atoi(i)
		if err != nil {
			return nil, fmt.Errorf("%s is not an int", i)
		}
		out = append(out, id)
	}
	return out, nil
}

// RunTaskPurge purge a task from hub.
func RunTaskPurge(c *CmdConfig) error {

	userid, err := c.Ankr.GetInt(c.NS, akrctl.ArgUserID)
	if err != nil {
		return err
	}

	force, err := c.Ankr.GetBool(c.NS, akrctl.ArgForce)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	if force || AskForConfirm(fmt.Sprintf("purge %d task(s)", len(c.Args))) == nil {
		url := viper.GetString("hub-url")

		conn, err := grpc.Dial(url+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		dc := pb.NewTaskMgrClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), ankr_const.ClientTimeOut*time.Second)
		defer cancel()

		fn := func(ids []string) error {
			for _, id := range ids {
				if ctr, _ := dc.PurgeTask(ctx, &pb.Request{UserId: int64(userid), TaskId: id}); err != nil {
					return fmt.Errorf("unable to purge task %d: %v", id, ctr.Details)
				}
			}
			return nil
		}
		return fn(c.Args)

	}
	return fmt.Errorf("operation aborted")

}

// RunTaskDelete destroy a task by id.
func RunTaskDelete(c *CmdConfig) error {

	userid, err := c.Ankr.GetInt(c.NS, akrctl.ArgUserID)
	if err != nil {
		return err
	}

	force, err := c.Ankr.GetBool(c.NS, akrctl.ArgForce)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	if force || AskForConfirm(fmt.Sprintf("delete %d task(s)", len(c.Args))) == nil {
		url := viper.GetString("hub-url")

		conn, err := grpc.Dial(url+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		defer conn.Close()
		dc := pb.NewTaskMgrClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), ankr_const.ClientTimeOut*time.Second)
		defer cancel()

		fn := func(ids []string) error {
			for _, id := range ids {
				if ctr, _ := dc.CancelTask(ctx, &pb.Request{UserId: int64(userid), TaskId: id}); ctr != nil {
					return fmt.Errorf("unable to delete task %d: %v", id, ctr.Details)
				}
			}
			return nil
		}

		return fn(c.Args)
	}
	return fmt.Errorf("operation aborted")

}

// RunTaskList returns a list of tasks.
func RunTaskList(c *CmdConfig) error {

	userid, err := c.Ankr.GetInt(c.NS, akrctl.ArgUserID)
	if err != nil {
		return err
	}

	matches := []glob.Glob{}
	for _, globStr := range c.Args {
		g, err := glob.Compile(globStr)
		if err != nil {
			return fmt.Errorf("unknown glob %q", globStr)
		}

		matches = append(matches, g)
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	dc := pb.NewTaskMgrClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ankr_const.ClientTimeOut*time.Second)
	defer cancel()
	r, err := dc.TaskList(ctx, &pb.ID{UserId: int64(userid)})
	if err != nil {
		log.Fatalf("Client: could not send: %v", err)
	}

	item := &displayers.Task{Tasks: r.Tasks}
	return c.Display(item)
}

// RunTaskUpdate updates a task.
//DCCN-CLI comput task update
func RunTaskUpdate(c *CmdConfig) error {

	userid, err := c.Ankr.GetInt(c.NS, akrctl.ArgUserID)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	name, err := c.Ankr.GetString(c.NS, akrctl.ArgNameSlug)
	if err != nil {
		return err
	}

	replica, err := c.Ankr.GetString(c.NS, akrctl.ArgReplicaSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	dc := pb.NewTaskMgrClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	fn := func(ids []string) error {
		for _, id := range ids {
			utrq := &pb.UpdateTaskRequest{
				UserId: int64(userid),
				Task: &common_proto.Task{
					Id: id,
				},
			}
			if replica != "" {
				replicaCount, err := strconv.Atoi(replica)
				if err != nil {
					return fmt.Errorf("replica count %s is not an int", replica)
				}
				utrq.Task.Replica = int32(replicaCount)
			}
			if name != "" {
				utrq.Task.Name = name
			}
			if utrp, _ := dc.UpdateTask(ctx, utrq); utrp != nil {
				return fmt.Errorf("unable to update task %d: %v", id, utrp.Details)
			}
		}
		return nil
	}
	return fn(c.Args)
}

// RunTaskDetail show a task detail by id.
func RunTaskDetail(c *CmdConfig) error {

	userid, err := c.Ankr.GetInt(c.NS, akrctl.ArgUserID)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	dc := pb.NewTaskMgrClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ankr_const.ClientTimeOut*time.Second)
	defer cancel()
	fn := func(ids []string) error {
		for _, id := range ids {
			if ctr, _ := dc.TaskDetail(ctx, &pb.Request{UserId: int64(userid)}); ctr.Error != nil {
				return fmt.Errorf("unable to get task %d detail: %v", id, ctr.Error.Details)
			}
		}
		return nil
	}
	return fn(c.Args)
}
