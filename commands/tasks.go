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
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"

	akrctl "github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"

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
	cmdRunTaskCreate := CmdBuilder(cmd, RunTaskCreate, "create <task-name> [task-name ...]",
		"create task", Writer, aliasOpt("cr"), displayerType(&displayers.Task{}), docCategories("task"))
	AddStringFlag(cmdRunTaskCreate, akrctl.ArgImageSlug, "", "", "Task image", requiredOpt())
	AddStringFlag(cmdRunTaskCreate, akrctl.ArgTypeSlug, "", "", "Task type")
	AddStringFlag(cmdRunTaskCreate, akrctl.ArgDcNameSlug, "", "", "Task data center name")
	AddStringFlag(cmdRunTaskCreate, akrctl.ArgReplicaSlug, "", "", "Task replica")
	AddStringFlag(cmdRunTaskCreate, akrctl.ArgScheduleSlug, "", "", "Task schedule")

	//DCCN-CLI comput task cancel
	cmdRunTaskCancel := CmdBuilder(cmd, RunTaskCancel, "cancel <task-id> [task-id ...]",
		"Cancel task by id", Writer, aliasOpt("dl", "del", "rm"), docCategories("Task"))
	AddBoolFlag(cmdRunTaskCancel, akrctl.ArgForce, akrctl.ArgShortForce, false, "Force task cancel")

	//DCCN-CLI comput task update
	cmdRunTaskUpdate := CmdBuilder(cmd, RunTaskUpdate, "update <task-id> [task-id ...]",
		"Update task by id", Writer, aliasOpt("ud", "udt", "ch"), docCategories("Task"))
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgImageSlug, "", "", "Task image")
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgReplicaSlug, "", "", "Task replica")
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgTypeSlug, "", "", "Task type")
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgDcNameSlug, "", "", "Task data center name")
	AddStringFlag(cmdRunTaskUpdate, akrctl.ArgScheduleSlug, "", "", "Task schedule")

	//DCCN-CLI task list
	cmdRunTaskList := CmdBuilder(cmd, RunTaskList, "list [GLOB]", "list tasks", Writer,
		aliasOpt("ls"), displayerType(&displayers.Task{}), docCategories("task"))
	AddStringFlag(cmdRunTaskList, akrctl.ArgTaskIdSlug, "", "", "Task id")

	//DCCN-CLI comput task purge
	cmdRunTaskPurge := CmdBuilder(cmd, RunTaskPurge, "purge <task-id> [task-id ...]", "Purge task by id",
		Writer, aliasOpt("pg"), docCategories("Task"))
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

	dcName, err := c.Ankr.GetString(c.NS, akrctl.ArgDcNameSlug)
	if err != nil {
		return err
	}

	taskType, err := c.Ankr.GetString(c.NS, akrctl.ArgTypeSlug)
	if err != nil {
		return err
	}
	if taskType == "" {
		taskType = "Deployment"
	}
	taskType = strings.ToUpper(taskType)

	replica, err := c.Ankr.GetString(c.NS, akrctl.ArgReplicaSlug)
	if err != nil {
		return err
	}

	task := common_proto.Task{}
	if dcName != "" {
		task.DataCenterName = dcName
	}

	switch taskType {
	case "DEPLOYMENT":
		task.Type = common_proto.TaskType_DEPLOYMENT
		task.TypeData = &common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image: image}}
	case "JOB":
		task.Type = common_proto.TaskType_JOB
		task.TypeData = &common_proto.Task_TypeJob{TypeJob: &common_proto.TaskTypeJob{Image: image}}
	case "CRONJOB":
		task.Type = common_proto.TaskType_CRONJOB
		schedule, err := c.Ankr.GetString(c.NS, akrctl.ArgScheduleSlug)
		if err != nil {
			return err
		}
		task.TypeData = &common_proto.Task_TypeCronJob{TypeCronJob: &common_proto.TaskTypeCronJob{Image: image, Schedule: schedule}}
	default:
		task.Type = common_proto.TaskType_DEPLOYMENT
		task.TypeData = &common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image: image}}
	}

	if replica != "" {
		r, err := strconv.Atoi(replica)
		if err != nil {
			return fmt.Errorf("replica count %s is not an int32", replica)
		}
		task.Attributes = &common_proto.TaskAttributes{Replica: int32(r)}
	}

	url := viper.GetString("hub-url")

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found")
	}

	md := metadata.New(map[string]string{
		"token": authResult.AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	taskMgr := pb.NewTaskMgrClient(conn)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	errs := make(chan error, len(c.Args))
	for _, name := range c.Args {
		task.Name = name
		tcrq := &pb.CreateTaskRequest{Task: &task}

		wg.Add(1)
		go func() {
			defer wg.Done()
			tcrp, err := taskMgr.CreateTask(tokenctx, tcrq)
			if err != nil {
				errs <- err
			} else {
				if tcrp != nil {
					fmt.Printf("Task %s Create Success. \n", tcrp.TaskId)
				}
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

// RunTaskPurge purge a task from hub.
func RunTaskPurge(c *CmdConfig) error {

	force, err := c.Ankr.GetBool(c.NS, akrctl.ArgForce)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found")
	}

	md := metadata.New(map[string]string{
		"token": authResult.AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	if force || AskForConfirm(fmt.Sprintf("Are you sure you want to Purge %d task(s) (y/N) ? ", len(c.Args))) == nil {
		url := viper.GetString("hub-url")

		conn, err := grpc.Dial(url+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()
		dc := pb.NewTaskMgrClient(conn)

		fn := func(ids []string) error {
			for _, id := range ids {
				_, err := dc.PurgeTask(tokenctx, &pb.TaskID{TaskId: id})
				if err != nil {
					return err
				}
				fmt.Printf("Task %s Purge Success.\n", id)
			}
			return nil
		}
		return fn(c.Args)

	}
	return fmt.Errorf("Operation aborted")

}

// RunTaskCancel destroy a task by id.
func RunTaskCancel(c *CmdConfig) error {

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found")
	}

	md := metadata.New(map[string]string{
		"token": authResult.AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	force, err := c.Ankr.GetBool(c.NS, akrctl.ArgForce)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	if force || AskForConfirm(fmt.Sprintf("Are you sure you want to Cancel %d task(s) (y/N) ? ", len(c.Args))) == nil {
		url := viper.GetString("hub-url")

		conn, err := grpc.Dial(url+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}

		defer conn.Close()
		taskMgr := pb.NewTaskMgrClient(conn)

		fn := func(ids []string) error {
			for _, id := range ids {
				_, err := taskMgr.CancelTask(tokenctx, &pb.TaskID{TaskId: id})
				if err != nil {
					return err
				}
				fmt.Printf("Task %s Cancel Success.\n", id)
			}
			return nil
		}

		return fn(c.Args)
	}
	return fmt.Errorf("Operation aborted")

}

// RunTaskList returns a list of tasks.
func RunTaskList(c *CmdConfig) error {

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found")
	}
	md := metadata.New(map[string]string{
		"token": authResult.AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

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
	defer conn.Close()

	taskMgr := pb.NewTaskMgrClient(conn)
	tlr := pb.TaskListRequest{}
	taskId, err := c.Ankr.GetString(c.NS, akrctl.ArgTaskIdSlug)
	if err != nil {
		return err
	}
	if taskId != "" {
		tlr.TaskFilter = &pb.TaskFilter{TaskId: taskId}
	}
	r, err := taskMgr.TaskList(tokenctx, &tlr)
	if err != nil {
		return err
	}

	item := &displayers.Task{Tasks: r.Tasks}
	return c.Display(item)
}

// RunTaskUpdate updates a task.
func RunTaskUpdate(c *CmdConfig) error {

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found")
	}

	md := metadata.New(map[string]string{
		"token": authResult.AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	image, err := c.Ankr.GetString(c.NS, akrctl.ArgImageSlug)
	if err != nil {
		return err
	}

	replica, err := c.Ankr.GetString(c.NS, akrctl.ArgReplicaSlug)
	if err != nil {
		return err
	}

	dcName, err := c.Ankr.GetString(c.NS, akrctl.ArgDcNameSlug)
	if err != nil {
		return err
	}

	taskType, err := c.Ankr.GetString(c.NS, akrctl.ArgTypeSlug)
	if err != nil {
		return err
	}

	task := common_proto.Task{}

	if dcName != "" {
		task.DataCenterName = dcName
	}

	if taskType != "" {
		switch taskType {
		case "DEPLOYMENT":
			task.Type = common_proto.TaskType_DEPLOYMENT
			taskDeploy := common_proto.TaskTypeDeployment{}
			if image != "" {
				taskDeploy.Image = image
			}
			task.TypeData = &common_proto.Task_TypeDeployment{TypeDeployment: &taskDeploy}
		case "JOB":
			task.Type = common_proto.TaskType_JOB
			taskDeploy := common_proto.TaskTypeJob{}
			if image != "" {
				taskDeploy.Image = image
			}
			task.TypeData = &common_proto.Task_TypeJob{TypeJob: &taskDeploy}
		case "CRONJOB":
			task.Type = common_proto.TaskType_CRONJOB
			taskDeploy := common_proto.TaskTypeCronJob{}
			if image != "" {
				taskDeploy.Image = image
			}
			schedule, err := c.Ankr.GetString(c.NS, akrctl.ArgScheduleSlug)
			if err != nil {
				return err
			}
			if schedule != "" {
				taskDeploy.Schedule = schedule
			}
			task.TypeData = &common_proto.Task_TypeCronJob{TypeCronJob: &taskDeploy}
		default:
			task.Type = common_proto.TaskType_DEPLOYMENT
			taskDeploy := common_proto.TaskTypeDeployment{}
			if image != "" {
				taskDeploy.Image = image
			}
			task.TypeData = &common_proto.Task_TypeDeployment{TypeDeployment: &taskDeploy}
		}
	}

	if replica != "" {
		r, err := strconv.Atoi(replica)
		if err != nil {
			return fmt.Errorf("replica count %s is not an int32", replica)
		}
		task.Attributes = &common_proto.TaskAttributes{Replica: int32(r)}
	}

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	taskMgr := pb.NewTaskMgrClient(conn)

	fn := func(ids []string) error {
		for _, id := range ids {
			task.Id = id
			utrq := &pb.UpdateTaskRequest{Task: &task}
			_, err := taskMgr.UpdateTask(tokenctx, utrq)
			if err != nil {
				return err
			}
			fmt.Printf("Task %s Update Success.\n", id)
		}
		return nil
	}
	return fn(c.Args)
}
