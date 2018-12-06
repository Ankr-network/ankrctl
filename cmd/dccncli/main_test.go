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
	"strings"
	"testing"

	"github.com/Ankr-network/dccn-cli"

	"github.com/stretchr/testify/assert"

	"fmt"
	"log"
	"net"

	pb "github.com/Ankr-network/dccn-rpc/protocol"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	lc := dccncli.NewLiveCommand("go")
	taskCreate, err := lc.Run("run", "main.go", "compute", "task", "create", "nginx:1.12", "-u", "localhost")
	assert.NoError(t, err)
	assert.True(t, len(string(taskCreate)) > 0)
	assert.True(t, strings.Contains(string(taskCreate), "created successfully"))
	id := strings.Fields(string(taskCreate))[2]
	assert.True(t, len(id) > 0)

	taskList, err := lc.Run("run", "main.go", "compute", "task", "list", "-u", "localhost")
	assert.NoError(t, err)
	assert.True(t, len(string(taskList)) > 0)
	taskInfo := strings.Split(string(taskList), "\n")
	taskFound := false
	for _, task := range taskInfo {
		if task != "" && id == strings.Fields(string(task))[0] {
			assert.Equal(t, strings.Fields(string(task))[1], "nginx:1.12")
			assert.Equal(t, strings.Fields(string(task))[4], "running")
			taskFound = true
		}
	}
	assert.True(t, taskFound)

	taskDelete, err := lc.Run("run", "main.go", "compute", "task", "delete", "-f", id, "-u", "localhost")
	assert.NoError(t, err)
	assert.True(t, len(string(taskDelete)) > 0)
	assert.Equal(t, id, string(bytes.Split(taskDelete, []byte(" "))[3]))
	assert.True(t, strings.Contains(string(taskDelete), "...Success!"))
}

type server struct {
	Taskid   int64
	TaskName string
	Status   string
}

func (s *server) AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	fmt.Printf("received add task request, creating task with id 100\n")
	s.Taskid = 100
	s.TaskName = in.Name
	s.Status = "running"
	//util.Task{Name: in.Name, Region: in.Region, Zone: in.Zone, ID: 100, Status: "running"}
	return &pb.AddTaskResponse{Status: "Success", Taskid: s.Taskid}, nil
}

func (s *server) TaskList(ctx context.Context, in *pb.TaskListRequest) (*pb.TaskListResponse, error) {
	fmt.Printf("task list reqeust, returning with task id 100\n")
	var taskList []*pb.TaskInfo
	taskInfo := &pb.TaskInfo{}
	taskInfo.Taskid = s.Taskid
	taskInfo.Taskname = s.TaskName
	taskInfo.Status = s.Status
	taskList = append(taskList, taskInfo)
	return &pb.TaskListResponse{Tasksinfo: taskList}, nil
}

func (s *server) CancelTask(ctx context.Context, in *pb.CancelTaskRequest) (*pb.CancelTaskResponse, error) {
	fmt.Printf("received cancel task request, delete task id %d\n", s.Taskid)
	if in.Taskid != s.Taskid {
		fmt.Printf("can not find task\n")
		return &pb.CancelTaskResponse{Status: "Failure"}, nil
	}
	return &pb.CancelTaskResponse{Status: "Success"}, nil
}

func (s *server) K8ReportStatus(ctx context.Context, in *pb.ReportRequest) (*pb.ReportResponse, error) {
	return &pb.ReportResponse{Status: "Success"}, nil
}

func (s *server) K8Task(stream pb.Dccncli_K8TaskServer) error {
	return nil
}
