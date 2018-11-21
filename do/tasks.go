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

package do

import (
	"context"

	"github.com/Ankr-network/godo"
	//"github.com/Ankr-network/godo/util"
)

// TaskIPTable is a table of interface IPS.
type TaskIPTable map[InterfaceType]string

// InterfaceType is a an interface type.
type InterfaceType string

const (
	// InterfacePublic is a public interface.
	InterfacePublic InterfaceType = "public"
	// InterfacePrivate is a private interface.
	InterfacePrivate InterfaceType = "private"
)

// Task is a wrapper for godo.Task
type Task struct {
	*godo.Task
}

// Tasks is a slice of Task.
type Tasks []Task

// Kernel is a wrapper for godo.Kernel
type Kernel struct {
	*godo.Kernel
}

// Kernels is a slice of Kernel.
type Kernels []Kernel

// TasksService is an interface for interacting with AnkrNetwork's task api.
type TasksService interface {
	//DCCN-CLI task list
	List() (Tasks, error)
	ListByTag(string) (Tasks, error)
	Get(int) (*Task, error)
	Create(*godo.TaskCreateRequest, bool) (*Task, error)
	CreateMultiple(*godo.TaskMultiCreateRequest) (Tasks, error)
	Delete(int) (string, error)
	DeleteByTag(string) error
	Kernels(int) (Kernels, error)
	Snapshots(int) (Images, error)
	Backups(int) (Images, error)
	Actions(int) (Actions, error)
	Neighbors(int) (Tasks, error)
}

type tasksService struct {
	client *godo.Client
}

var _ TasksService = &tasksService{}

// NewTasksService builds a TasksService instance.
func NewTasksService(client *godo.Client) TasksService {
	return &tasksService{
		client: client,
	}
}
//DCCN-CLI task list
func (ds *tasksService) List() (Tasks, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Tasks.List(context.TODO(), opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}
	
	list := make(Tasks, len(si))
	for i := range si {
		a := si[i].(godo.Task)
		list[i] = Task{Task: &a}
	}

	return list, nil
}

func (ds *tasksService) ListByTag(tagName string) (Tasks, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Tasks.ListByTag(context.TODO(), tagName, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Tasks, len(si))
	for i := range si {
		a := si[i].(godo.Task)
		list[i] = Task{Task: &a}
	}

	return list, nil
}

func (ds *tasksService) Get(id int) (*Task, error) {
	d, _, err := ds.client.Tasks.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	return &Task{Task: d}, nil
}
//DCCN-CLI comput task create
func (ds *tasksService) Create(dcr *godo.TaskCreateRequest, wait bool) (*Task, error) {
	d, resp, err := ds.client.Tasks.Create(context.TODO(), dcr)
	if err != nil {
		return nil, err
	}
	_ = resp
	/*
	if wait {
		var action *godo.LinkAction
		for _, a := range resp.Links.Actions {
			if a.Rel == "create" {
				action = &a
				break
			}
		}

		if action != nil {
			_ = util.WaitForActive(context.TODO(), ds.client, action.HREF)
			doTask, err := ds.Get(d.ID)
			if err != nil {
				return nil, err
			}
			d = doTask.Task
		}
	}
	*/
	return &Task{Task: d}, nil
}

func (ds *tasksService) CreateMultiple(dmcr *godo.TaskMultiCreateRequest) (Tasks, error) {
	godoTasks, _, err := ds.client.Tasks.CreateMultiple(context.TODO(), dmcr)
	if err != nil {
		return nil, err
	}

	var tasks Tasks
	for _, d := range godoTasks {
		tasks = append(tasks, Task{Task: &d})
	}

	return tasks, nil
}
//DCCN-CLI comput task delete
func (ds *tasksService) Delete(id int) (string, error) {
	status, err := ds.client.Tasks.Delete(context.TODO(), id)
	return status, err
}

func (ds *tasksService) DeleteByTag(tag string) error {
	_, err := ds.client.Tasks.DeleteByTag(context.TODO(), tag)
	return err
}

func (ds *tasksService) Kernels(id int) (Kernels, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Tasks.Kernels(context.TODO(), id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Kernels, len(si))
	for i := range si {
		a := si[i].(godo.Kernel)
		list[i] = Kernel{Kernel: &a}
	}

	return list, nil
}

func (ds *tasksService) Snapshots(id int) (Images, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Tasks.Snapshots(context.TODO(), id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Images, len(si))
	for i := range si {
		a := si[i].(godo.Image)
		list[i] = Image{Image: &a}
	}

	return list, nil
}

func (ds *tasksService) Backups(id int) (Images, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Tasks.Backups(context.TODO(), id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Images, len(si))
	for i := range si {
		a := si[i].(godo.Image)
		list[i] = Image{Image: &a}
	}

	return list, nil
}

func (ds *tasksService) Actions(id int) (Actions, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ds.client.Tasks.Actions(context.TODO(), id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Actions, len(si))
	for i := range si {
		a := si[i].(godo.Action)
		list[i] = Action{Action: &a}
	}

	return list, nil
}

func (ds *tasksService) Neighbors(id int) (Tasks, error) {
	list, _, err := ds.client.Tasks.Neighbors(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	var tasks Tasks
	for _, d := range list {
		tasks = append(tasks, Task{Task: &d})
	}

	return tasks, nil
}
