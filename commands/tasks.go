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
	"errors"
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
	"google.golang.org/grpc/metadata"
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
	AddStringFlag(cmdTaskCreate, akrctl.ArgDcidSlug, "", "", "Task data center id")
	AddStringFlag(cmdTaskCreate, akrctl.ArgImageSlug, "", "", "Task image")
	AddStringFlag(cmdTaskCreate, akrctl.ArgTypeSlug, "", "", "Task type")
	AddStringFlag(cmdTaskCreate, akrctl.ArgReplicaSlug, "", "", "Task replica")

	//DCCN-CLI comput task cancel
	cmdRunTaskCancel := CmdBuilder(cmd, RunTaskCancel, "cancel <task-id> [task-id ...]", "Cancel task by id", Writer,
		aliasOpt("dl", "del", "rm"), docCategories("Task"))
	AddBoolFlag(cmdRunTaskCancel, akrctl.ArgForce, akrctl.ArgShortForce, false, "Force task cancel")

	//DCCN-CLI comput task update
	cmdRunTaskUpdate := CmdBuilder(cmd, RunTaskUpdate, "update <task-id> [task-id ...]", "Update task by id", Writer,
		aliasOpt("ud", "udt", "ch"), docCategories("Task"))
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgImageSlug, "", "", "Task image")
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgReplicaSlug, "", "", "Task replica")
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgTypeSlug, "", "", "Task type")
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgDcidSlug, "", "", "Task data center id")

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

	image, err := c.Ankr.GetString(c.NS, akrctl.ArgImageSlug)
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
	if tasktype == "" {
		tasktype = "Default"
	}

	replica, err := c.Ankr.GetString(c.NS, akrctl.ArgReplicaSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")

	token, userid := c.getContextAccessToken()

	if token == "" {
		return fmt.Errorf("Unable to read AnkrNetwork access token.")
	}
	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	dc := pb.NewTaskMgrClient(conn)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	errs := make(chan *common_proto.Error, len(c.Args))
	for _, name := range c.Args {
		tcrq := &pb.CreateTaskRequest{
			UserId: userid,
			Task: &common_proto.Task{
				Name:  name,
				Image: image,
				Type:  common_proto.TaskType(common_proto.TaskType_value[tasktype+"TaskType"]),
			},
		}
		if taskdcid != "" {
			tcrq.Task.DataCenterId = taskdcid
		}
		if replica != "" {
			r, err := strconv.Atoi(replica)
			if err != nil {
				return fmt.Errorf("Replica count %s is not an int\n", replica)
			}
			tcrq.Task.Replica = int32(r)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			tcrp, err := dc.CreateTask(tokenctx, tcrq)
			if err != nil {
				log.Fatal(err.Error())
			} else {
				fmt.Printf("Task %s Create Success. \n", tcrp.TaskId)
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			return errors.New(err.Details)
		}
	}

	return nil
}

// RunTaskPurge purge a task from hub.
func RunTaskPurge(c *CmdConfig) error {

	userid, err := c.Ankr.GetString(c.NS, akrctl.ArgUserID)
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

	token, userid := c.getContextAccessToken()

	if token == "" {
		return fmt.Errorf("Unable to read Ankr Network access token")
	}

	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if force || AskForConfirm(fmt.Sprintf("Purge %d task(s)", len(c.Args))) == nil {
		url := viper.GetString("hub-url")

		conn, err := grpc.Dial(url+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()
		dc := pb.NewTaskMgrClient(conn)
		tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
		defer cancel()

		fn := func(ids []string) error {
			for _, id := range ids {
				if ctr, _ := dc.PurgeTask(tokenctx, &pb.Request{UserId: userid, TaskId: id}); ctr != nil && ctr.Status == common_proto.Status_FAILURE {
					return fmt.Errorf("Unable to purge task %s: %v", id, ctr.Details)
				} else {
					fmt.Printf("Task %s Purge Success.\n", id)
				}
			}
			return nil
		}
		return fn(c.Args)

	}
	return fmt.Errorf("Operation aborted")

}

// RunTaskCancel destroy a task by id.
func RunTaskCancel(c *CmdConfig) error {

	userid, err := c.Ankr.GetString(c.NS, akrctl.ArgUserID)
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

	token, userid := c.getContextAccessToken()

	if token == "" {
		return fmt.Errorf("Unable to read AnkrNetwork access token")
	}

	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if force || AskForConfirm(fmt.Sprintf("Cancel %d task(s)", len(c.Args))) == nil {
		url := viper.GetString("hub-url")

		conn, err := grpc.Dial(url+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}

		defer conn.Close()
		dc := pb.NewTaskMgrClient(conn)
		tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
		defer cancel()

		fn := func(ids []string) error {
			for _, id := range ids {
				if ctr, _ := dc.CancelTask(tokenctx, &pb.Request{UserId: userid, TaskId: id}); ctr != nil && ctr.Status == common_proto.Status_FAILURE {
					return fmt.Errorf("Unable to cancel task %s: %v", id, ctr.Details)
				} else {
					fmt.Printf("Task %s Cancel Success.\n", id)
				}
			}
			return nil
		}

		return fn(c.Args)
	}
	return fmt.Errorf("Operation aborted")

}

// RunTaskList returns a list of tasks.
func RunTaskList(c *CmdConfig) error {

	userid, err := c.Ankr.GetString(c.NS, akrctl.ArgUserID)
	if err != nil {
		return err
	}

	matches := []glob.Glob{}
	for _, globStr := range c.Args {
		g, err := glob.Compile(globStr)
		if err != nil {
			return fmt.Errorf("Unknown glob %q", globStr)
		}

		matches = append(matches, g)
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	token, userid := c.getContextAccessToken()

	if token == "" {
		return fmt.Errorf("Unable to read Ankr Network access token.")
	}

	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	defer conn.Close()
	dc := pb.NewTaskMgrClient(conn)
	defer cancel()
	r, _ := dc.TaskList(tokenctx, &pb.ID{UserId: userid})
	if r.Error != nil && r.Error.Status == common_proto.Status_FAILURE {
}

// RunTaskUpdate updates a task.
//DCCN-CLI comput task update
func RunTaskUpdate(c *CmdConfig) error {

	userid, err := c.Ankr.GetString(c.NS, akrctl.ArgUserID)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	image, err := c.Ankr.GetString(c.NS, akrctl.ArgNameSlug)
	if err != nil {
		return err
	}

	replica, err := c.Ankr.GetString(c.NS, akrctl.ArgReplicaSlug)
	if err != nil {
		return err
	}

	dcid, err := c.Ankr.GetString(c.NS, akrctl.ArgDcidSlug)
	if err != nil {
		return err
	}

	tasktype, err := c.Ankr.GetString(c.NS, akrctl.ArgTypeSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	token, userid := c.getContextAccessToken()

	if token == "" {
		return fmt.Errorf("Unable to read Ankr Network access token")
	}

	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	defer conn.Close()
	dc := pb.NewTaskMgrClient(conn)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	fn := func(ids []string) error {
		for _, id := range ids {
			utrq := &pb.UpdateTaskRequest{
				UserId: userid,
				Task: &common_proto.Task{
					Id: id,
				},
			}
			if replica != "" {
				replicaCount, err := strconv.Atoi(replica)
				if err != nil {
					return fmt.Errorf("Replica count %s is not an int", replica)
				}
				utrq.Task.Replica = int32(replicaCount)
			}
			if image != "" {
				utrq.Task.Name = image
			}
			if tasktype != "" {
				utrq.Task.Type = common_proto.TaskType(common_proto.TaskType_value[tasktype+"TaskType"])
			}
			if dcid != "" {
				utrq.Task.DataCenterId = dcid
			}
			if utrp, _ := dc.UpdateTask(tokenctx, utrq); utrp != nil && utrp.Status == common_proto.Status_FAILURE {
				return fmt.Errorf("Unable to update task %s: %v", id, utrp.Details)
			} else {
				fmt.Printf("Task %s Update Success.\n", id)
			}
		}
		return nil
	}
	return fn(c.Args)
}

// RunTaskDetail show a task detail by id.
func RunTaskDetail(c *CmdConfig) error {

	userid, err := c.Ankr.GetString(c.NS, akrctl.ArgUserID)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	token, userid := c.getContextAccessToken()

	if token == "" {
		return fmt.Errorf("Unable to read Ankr Network access token")
	}

	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	defer conn.Close()
	dc := pb.NewTaskMgrClient(conn)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()
	fn := func(ids []string) error {
		for _, id := range ids {
			if ctr, _ := dc.TaskDetail(tokenctx, &pb.Request{UserId: userid, TaskId: id}); ctr.Error != nil && ctr.Error.Status == common_proto.Status_FAILURE {
				return fmt.Errorf("Unable to get task %s detail: %v", id, ctr.Error.Details)
			} else {
				fmt.Printf("Task %s Detail Success.\n", ctr.Task.Id)
			}
		}
		return nil
	}
	return fn(c.Args)
}
