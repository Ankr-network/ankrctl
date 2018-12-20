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

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"

	"context"

	pb "github.com/Ankr-network/dccn-rpc/protocol"
	"google.golang.org/grpc"
)

const (
	port = "50051"
)

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
	AddStringFlag(cmdTaskCreate, akrctl.ArgDcSlug, "", "", "Task dc")
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
	return cmd
}

// RunTaskCreate creates a task.
//DCCN-CLI comput task create
func RunTaskCreate(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	taskdc, err := c.Ankr.GetString(c.NS, akrctl.ArgDcSlug)
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
	if url == "" {
		url += clientURL
	}

	conn, err := grpc.Dial(url+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	dc := pb.NewDccncliClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	errs := make(chan error, len(c.Args))
	for _, name := range c.Args {
		tcrq := &pb.AddTaskRequest{
			Name:       name,
			Type:       tasktype,
			Datacenter: taskdc,
			Usertoken:  "ed1605e17374bde6c68864d072c9f5c9",
		}
		if replica != "" {
			replicaCount, err := strconv.Atoi(replica)
			if err != nil {
				return fmt.Errorf("replica count %s is not an int", replica)
			}
			tcrq.Replica = int64(replicaCount)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			tcrp, err := dc.AddTask(ctx, tcrq)
			if err != nil {
				errs <- err
				return
			}
			if tcrp.Status == "Success" {
				fmt.Printf("Task id %d created successfully. \n", tcrp.Taskid)
			} else {
				fmt.Printf("Fail to create task. \n")
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

// RunTaskDelete destroy a task by id.
func RunTaskDelete(c *CmdConfig) error {

	force, err := c.Ankr.GetBool(c.NS, akrctl.ArgForce)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	if force || AskForConfirm(fmt.Sprintf("delete %d task(s)", len(c.Args))) == nil {
		url := viper.GetString("hub-url")
		if url == "" {
			url += clientURL
		}
		conn, err := grpc.Dial(url+":"+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		dc := pb.NewDccncliClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fn := func(ids []int) error {
			for _, id := range ids {
				if ctr, err := dc.CancelTask(ctx, &pb.CancelTaskRequest{Taskid: int64(id), Usertoken: "ed1605e17374bde6c68864d072c9f5c9"}); err != nil {
					return fmt.Errorf("unable to delete task %d: %v", id, err)
				} else {
					fmt.Printf("Delete task id %d ...%s! \n", id, ctr.Status)
				}
			}
			return nil
		}
		if extractedIDs, err := allInt(c.Args); err == nil {
			return fn(extractedIDs)
		}
		return err

	}
	return fmt.Errorf("operation aborted")

}

// RunTaskList returns a list of tasks.
func RunTaskList(c *CmdConfig) error {

	matches := []glob.Glob{}
	for _, globStr := range c.Args {
		g, err := glob.Compile(globStr)
		if err != nil {
			return fmt.Errorf("unknown glob %q", globStr)
		}

		matches = append(matches, g)
	}

	var matchedList []pb.TaskInfo

	url := viper.GetString("hub-url")
	if url == "" {
		url += clientURL
	}
	conn, err := grpc.Dial(url+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	dc := pb.NewDccncliClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	r, err := dc.TaskList(ctx, &pb.TaskListRequest{Usertoken: "ed1605e17374bde6c68864d072c9f5c9"})
	if err != nil {
		log.Fatalf("Client: could not send: %v", err)
	}
	Taskinfos := r.Tasksinfo

	for _, taskinfo := range Taskinfos {
		var skip = true
		if len(matches) == 0 {
			skip = false
		} else {
			for _, m := range matches {
				if m.Match(taskinfo.Taskname) {
					skip = false
				}
			}
		}

		var task pb.TaskInfo
		task.Taskid = taskinfo.Taskid
		task.Taskname = taskinfo.Taskname
		task.Uptime = taskinfo.Uptime
		task.Creationdate = taskinfo.Creationdate
		task.Status = taskinfo.Status
		task.Datacenter = taskinfo.Datacenter
		task.Replica = taskinfo.Replica

		if !skip {
			matchedList = append(matchedList, task)
		}
	}

	item := &displayers.Task{Tasks: matchedList}
	return c.Display(item)
}

// RunTaskUpdate updates a task.
//DCCN-CLI comput task update
func RunTaskUpdate(c *CmdConfig) error {

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
	if url == "" {
		url += clientURL
	}

	conn, err := grpc.Dial(url+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	dc := pb.NewDccncliClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fn := func(ids []int) error {
		for _, id := range ids {
			utrq := &pb.UpdateTaskRequest{
				Taskid:    int64(id),
				Usertoken: "ed1605e17374bde6c68864d072c9f5c9",
			}
			if replica != "" {
				replicaCount, err := strconv.Atoi(replica)
				if err != nil {
					return fmt.Errorf("replica count %s is not an int", replica)
				}
				utrq.Replica = int64(replicaCount)
			}
			if name != "" {
				utrq.Name = name
			}
			if utrp, err := dc.UpdateTask(ctx, utrq); err != nil {
				return fmt.Errorf("unable to update task %d: %v", id, err)
			} else {
				fmt.Printf("Update task id %d ...%s! \n", id, utrp.Status)
			}
		}
		return nil
	}
	if extractedIDs, err := allInt(c.Args); err == nil {
		return fn(extractedIDs)
	}
	return err
}

// RunTaskDetail show a task detail by id.
func RunTaskDetail(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	url := viper.GetString("hub-url")
	if url == "" {
		url += clientURL
	}
	conn, err := grpc.Dial(url+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	dc := pb.NewDccncliClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fn := func(ids []int) error {
		for _, id := range ids {
			if ctr, err := dc.TaskDetail(ctx, &pb.TaskDetailRequest{Taskid: int64(id), Usertoken: "ed1605e17374bde6c68864d072c9f5c9"}); err != nil {
				return fmt.Errorf("unable to get task %d detail: %v", id, err)
			} else {
				fmt.Printf("task id %d detail:\n%s\n", id, ctr.Body)
			}
		}
		return nil
	}
	if extractedIDs, err := allInt(c.Args); err == nil {
		return fn(extractedIDs)
	}
	return err

}
