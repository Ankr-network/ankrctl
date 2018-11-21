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
)

// TaskActionsService is an interface for interacting with AnkrNetwork's task action api.
type TaskActionsService interface {
	Shutdown(int) (*Action, error)
	ShutdownByTag(string) (Actions, error)
	PowerOff(int) (*Action, error)
	PowerOffByTag(string) (Actions, error)
	PowerOn(int) (*Action, error)
	PowerOnByTag(string) (Actions, error)
	PowerCycle(int) (*Action, error)
	PowerCycleByTag(string) (Actions, error)
	Reboot(int) (*Action, error)
	Restore(int, int) (*Action, error)
	Resize(int, string, bool) (*Action, error)
	Rename(int, string) (*Action, error)
	Snapshot(int, string) (*Action, error)
	SnapshotByTag(string, string) (Actions, error)
	EnableBackups(int) (*Action, error)
	EnableBackupsByTag(string) (Actions, error)
	DisableBackups(int) (*Action, error)
	DisableBackupsByTag(string) (Actions, error)
	PasswordReset(int) (*Action, error)
	RebuildByImageID(int, int) (*Action, error)
	RebuildByImageSlug(int, string) (*Action, error)
	ChangeKernel(int, int) (*Action, error)
	EnableIPv6(int) (*Action, error)
	EnableIPv6ByTag(string) (Actions, error)
	EnablePrivateNetworking(int) (*Action, error)
	EnablePrivateNetworkingByTag(string) (Actions, error)
	Get(int, int) (*Action, error)
	GetByURI(string) (*Action, error)
}

type taskActionsService struct {
	client *godo.Client
}

var _ TaskActionsService = &taskActionsService{}

// NewTaskActionsService builds an instance of TaskActionsService.
func NewTaskActionsService(godoClient *godo.Client) TaskActionsService {
	return &taskActionsService{
		client: godoClient,
	}
}

func (das *taskActionsService) handleActionResponse(a *godo.Action, err error) (*Action, error) {
	if err != nil {
		return nil, err
	}

	return &Action{Action: a}, nil
}

func (das *taskActionsService) handleTagActionResponse(a []godo.Action, err error) (Actions, error) {
	if err != nil {
		return nil, err
	}

	actions := make([]Action, 0, len(a))

	for _, action := range a {
		actions = append(actions, Action{Action: &action})
	}

	return actions, nil
}

func (das *taskActionsService) Shutdown(id int) (*Action, error) {
	a, _, err := das.client.TaskActions.Shutdown(context.TODO(), id)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) ShutdownByTag(tag string) (Actions, error) {
	a, _, err := das.client.TaskActions.ShutdownByTag(context.TODO(), tag)
	return das.handleTagActionResponse(a, err)
}

func (das *taskActionsService) PowerOff(id int) (*Action, error) {
	a, _, err := das.client.TaskActions.PowerOff(context.TODO(), id)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) PowerOffByTag(tag string) (Actions, error) {
	a, _, err := das.client.TaskActions.PowerOffByTag(context.TODO(), tag)
	return das.handleTagActionResponse(a, err)
}

func (das *taskActionsService) PowerOn(id int) (*Action, error) {
	a, _, err := das.client.TaskActions.PowerOn(context.TODO(), id)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) PowerOnByTag(tag string) (Actions, error) {
	a, _, err := das.client.TaskActions.PowerOnByTag(context.TODO(), tag)
	return das.handleTagActionResponse(a, err)
}

func (das *taskActionsService) PowerCycle(id int) (*Action, error) {
	a, _, err := das.client.TaskActions.PowerCycle(context.TODO(), id)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) PowerCycleByTag(tag string) (Actions, error) {
	a, _, err := das.client.TaskActions.PowerCycleByTag(context.TODO(), tag)
	return das.handleTagActionResponse(a, err)
}

func (das *taskActionsService) Reboot(id int) (*Action, error) {
	a, _, err := das.client.TaskActions.Reboot(context.TODO(), id)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) Restore(id, imageID int) (*Action, error) {
	a, _, err := das.client.TaskActions.Restore(context.TODO(), id, imageID)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) Resize(id int, sizeSlug string, resizeDisk bool) (*Action, error) {
	a, _, err := das.client.TaskActions.Resize(context.TODO(), id, sizeSlug, resizeDisk)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) Rename(id int, name string) (*Action, error) {
	a, _, err := das.client.TaskActions.Rename(context.TODO(), id, name)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) Snapshot(id int, name string) (*Action, error) {
	a, _, err := das.client.TaskActions.Snapshot(context.TODO(), id, name)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) SnapshotByTag(tag string, name string) (Actions, error) {
	a, _, err := das.client.TaskActions.SnapshotByTag(context.TODO(), tag, name)
	return das.handleTagActionResponse(a, err)
}

func (das *taskActionsService) EnableBackups(id int) (*Action, error) {
	a, _, err := das.client.TaskActions.EnableBackups(context.TODO(), id)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) EnableBackupsByTag(tag string) (Actions, error) {
	a, _, err := das.client.TaskActions.EnableBackupsByTag(context.TODO(), tag)
	return das.handleTagActionResponse(a, err)
}

func (das *taskActionsService) DisableBackups(id int) (*Action, error) {
	a, _, err := das.client.TaskActions.DisableBackups(context.TODO(), id)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) DisableBackupsByTag(tag string) (Actions, error) {
	a, _, err := das.client.TaskActions.DisableBackupsByTag(context.TODO(), tag)
	return das.handleTagActionResponse(a, err)
}

func (das *taskActionsService) PasswordReset(id int) (*Action, error) {
	a, _, err := das.client.TaskActions.PasswordReset(context.TODO(), id)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) RebuildByImageID(id, imageID int) (*Action, error) {
	a, _, err := das.client.TaskActions.RebuildByImageID(context.TODO(), id, imageID)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) RebuildByImageSlug(id int, slug string) (*Action, error) {
	a, _, err := das.client.TaskActions.RebuildByImageSlug(context.TODO(), id, slug)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) ChangeKernel(id, kernelID int) (*Action, error) {
	a, _, err := das.client.TaskActions.ChangeKernel(context.TODO(), id, kernelID)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) EnableIPv6(id int) (*Action, error) {
	a, _, err := das.client.TaskActions.EnableIPv6(context.TODO(), id)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) EnableIPv6ByTag(tag string) (Actions, error) {
	a, _, err := das.client.TaskActions.EnableIPv6ByTag(context.TODO(), tag)
	return das.handleTagActionResponse(a, err)
}

func (das *taskActionsService) EnablePrivateNetworking(id int) (*Action, error) {
	a, _, err := das.client.TaskActions.EnablePrivateNetworking(context.TODO(), id)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) EnablePrivateNetworkingByTag(tag string) (Actions, error) {
	a, _, err := das.client.TaskActions.EnablePrivateNetworkingByTag(context.TODO(), tag)
	return das.handleTagActionResponse(a, err)
}

func (das *taskActionsService) Get(id int, actionID int) (*Action, error) {
	a, _, err := das.client.TaskActions.Get(context.TODO(), id, actionID)
	return das.handleActionResponse(a, err)
}

func (das *taskActionsService) GetByURI(uri string) (*Action, error) {
	a, _, err := das.client.TaskActions.GetByURI(context.TODO(), uri)
	return das.handleActionResponse(a, err)
}
