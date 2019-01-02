// +build !windows

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

package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/Ankr-network/dccn-cli"

	"github.com/stretchr/testify/assert"

	"fmt"
	"log"
	"net"

	pb "github.com/Ankr-network/dccn-common/protocol/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	MockTaskid        = 100
	MockTaskname      = "nginx:1.12"
	MockReplica       = "2"
	MockUpdateName    = "nginx:1.13"
	MockUpdateReplica = "3"
	MockStatus        = "running"
	MockType          = "web"
	MockDcid          = "1"
	MockDcname        = "data-center1"
)

func TestMockCommand_Run(t *testing.T) {
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		ss := server{}
		pb.RegisterDccncliServer(s, &ss)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	lc := akrctl.NewLiveCommand("go")

	//compute task create test
	taskCreate, err := lc.Run("run", "main.go", "compute", "task", "create",
		MockTaskname, "--dc-id", MockDcid, "--type", MockType, "--replica", MockReplica,
		"-u", "localhost")
	assert.NoError(t, err)
	assert.True(t, len(string(taskCreate)) > 0)
	assert.True(t, strings.Contains(string(taskCreate), "initialize successfully"))
	id := strings.Fields(string(taskCreate))[2]
	assert.True(t, len(id) > 0)

	//compute task list test
	taskList, err := lc.Run("run", "main.go", "compute", "task", "list", "-u", "localhost")
	assert.NoError(t, err)
	assert.True(t, len(string(taskList)) > 0)
	taskInfo := strings.Split(string(taskList), "\n")
	taskFound := false
	for _, task := range taskInfo {
		if task != "" && id == strings.Fields(string(task))[0] {
			assert.Equal(t, strings.Fields(string(task))[1], MockTaskname)
			assert.Equal(t, strings.Fields(string(task))[4], MockReplica)
			assert.Equal(t, strings.Fields(string(task))[5], MockStatus)
			taskFound = true
		}
	}
	assert.True(t, taskFound)

	//compute task update test
	taskUpdate, err := lc.Run("run", "main.go", "compute", "task", "update", id,
		"--name", MockUpdateName, "--replica", MockUpdateReplica, "-u", "localhost")
	assert.NoError(t, err)
	assert.True(t, len(string(taskUpdate)) > 0)
	assert.Equal(t, id, string(bytes.Split(taskUpdate, []byte(" "))[3]))
	assert.True(t, strings.Contains(string(taskUpdate), "...Success!"))

	//compute task delete test
	taskDelete, err := lc.Run("run", "main.go", "compute", "task", "delete", "-f", id, "-u", "localhost")
	assert.NoError(t, err)
	assert.True(t, len(string(taskDelete)) > 0)
	assert.Equal(t, id, string(bytes.Split(taskDelete, []byte(" "))[3]))
	assert.True(t, strings.Contains(string(taskDelete), "...Success!"))

	//compute task detail test
	taskDetail, err := lc.Run("run", "main.go", "compute", "task", "detail", strconv.Itoa(MockTaskid), "-u", "localhost")
	assert.NoError(t, err)
	assert.True(t, len(string(taskDetail)) > 0)
	detailInfo := strings.Split(string(taskDetail), "\n")
	assert.Equal(t, id, strings.Split(detailInfo[0], " ")[2])
	assert.True(t, len(string(detailInfo[1])) > 0)

	//compute dc list test
	dcList, err := lc.Run("run", "main.go", "compute", "dc", "list", "-u", "localhost")
	assert.NoError(t, err)
	assert.True(t, len(string(dcList)) > 0)
	dcInfo := strings.Split(string(dcList), "\n")
	dcFound := false
	for _, dc := range dcInfo {
		if dc != "" {
			dcFound = true
		}
	}
	assert.True(t, dcFound)
}

type server struct {
	Taskid       int64
	Taskname     string
	Status       string
	Tasktype     string
	Taskdc       string
	Replica      int64
	Datacenterid int64
	Datacenter   string
}

func (s *server) AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	fmt.Printf("received add task request, creating task with id %d\n", MockTaskid)
	s.Taskid = MockTaskid
	s.Status = MockStatus
	s.Taskname = in.Name
	s.Replica = in.Replica
	s.Datacenterid = in.Datacenterid
	return &pb.AddTaskResponse{Status: "Success", Taskid: s.Taskid}, nil
}

func (s *server) TaskList(ctx context.Context, in *pb.TaskListRequest) (*pb.TaskListResponse, error) {
	fmt.Printf("task list reqeust, returning with task id %d\n", s.Taskid)
	var taskList []*pb.TaskInfo
	taskInfo := &pb.TaskInfo{}
	taskInfo.Taskid = s.Taskid
	taskInfo.Taskname = s.Taskname
	taskInfo.Status = s.Status
	taskInfo.Replica = s.Replica
	taskInfo.Datacenter = s.Datacenter
	taskList = append(taskList, taskInfo)
	return &pb.TaskListResponse{Tasksinfo: taskList}, nil
}

func (s *server) DataCenterList(ctx context.Context, in *pb.DataCenterListRequest) (*pb.DataCenterListResponse, error) {
	fmt.Printf("dc list reqeust, returning with dc list\n")
	var dcList []*pb.DataCenterInfo
	dataCenterInfo := &pb.DataCenterInfo{}
	dataCenterInfo.Id = 1
	dataCenterInfo.Name = MockDcname
	dcList = append(dcList, dataCenterInfo)
	return &pb.DataCenterListResponse{DcList: dcList}, nil
}

func (s *server) TaskDetail(ctx context.Context, in *pb.TaskDetailRequest) (*pb.TaskDetailResponse, error) {
	fmt.Printf("task detail list reqeust, returning with task detail\n")
	return &pb.TaskDetailResponse{Body: "task detail"}, nil
}

func (s *server) CancelTask(ctx context.Context, in *pb.CancelTaskRequest) (*pb.CancelTaskResponse, error) {
	fmt.Printf("received cancel task request, delete task id %d\n", s.Taskid)
	if in.Taskid != s.Taskid {
		return &pb.CancelTaskResponse{Status: "Failure"}, fmt.Errorf("Can not find task.\n")
	}
	return &pb.CancelTaskResponse{Status: "Success"}, nil
}

func (s *server) PurgeTask(ctx context.Context, in *pb.PurgeTaskRequest) (*pb.PurgeTaskResponse, error) {
	fmt.Printf("received purge task request, purge task id %d\n", s.Taskid)
	if in.Taskid != s.Taskid {
		return &pb.PurgeTaskResponse{Status: "Failure"}, fmt.Errorf("Can not find task.\n")
	}
	return &pb.PurgeTaskResponse{Status: "Success"}, nil
}

func (s *server) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	fmt.Printf("received update task request, update task id %d\n", s.Taskid)
	if in.Taskid != s.Taskid {
		return &pb.UpdateTaskResponse{Status: "Failure"}, fmt.Errorf("Can not find task.\n")
	}
	if in.Name == s.Taskname {
		return &pb.UpdateTaskResponse{Status: "Failure"}, fmt.Errorf("Update task name no change.\n")
	}

	if in.Replica <= 0 {
		return &pb.UpdateTaskResponse{Status: "Failure"}, fmt.Errorf("Update replica not valid.\n")
	}

	return &pb.UpdateTaskResponse{Status: "Success"}, nil
}
