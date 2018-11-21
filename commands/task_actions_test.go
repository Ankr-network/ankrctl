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
	"testing"

	"github.com/Ankr-network/dccn-cli"
	"github.com/stretchr/testify/assert"
)

func TestTaskActionCommand(t *testing.T) {
	cmd := TaskAction()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "change-kernel", "enable-backups", "disable-backups", "enable-ipv6", "enable-private-networking", "get", "power-cycle", "power-off", "power-on", "password-reset", "reboot", "rebuild", "rename", "resize", "restore", "shutdown", "snapshot")
}

func TestTaskActionsChangeKernel(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("ChangeKernel", 1, 2).Return(&testAction, nil)

		config.Ankr.Set(config.NS, dccncli.ArgKernelID, 2)
		config.Args = append(config.Args, "1")

		err := RunTaskActionChangeKernel(config)
		assert.NoError(t, err)
	})
}
func TestTaskActionsEnableBackups(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("EnableBackups", 1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActionEnableBackups(config)
		assert.NoError(t, err)
	})

}
func TestTaskActionsDisableBackups(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("DisableBackups", 1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActionDisableBackups(config)
		assert.NoError(t, err)
	})

}
func TestTaskActionsEnableIPv6(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("EnableIPv6", 1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActionEnableIPv6(config)
		assert.NoError(t, err)
	})
}

func TestTaskActionsEnablePrivateNetworking(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("EnablePrivateNetworking", 1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActionEnablePrivateNetworking(config)
		assert.NoError(t, err)
	})
}
func TestTaskActionsGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("Get", 1, 2).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Ankr.Set(config.NS, dccncli.ArgActionID, 2)

		err := RunTaskActionGet(config)
		assert.NoError(t, err)
	})
}

func TestTaskActionsPasswordReset(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("PasswordReset", 1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActionPasswordReset(config)
		assert.NoError(t, err)
	})
}

func TestTaskActionsPowerCycle(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("PowerCycle", 1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActionPowerCycle(config)
		assert.NoError(t, err)
	})

}
func TestTaskActionsPowerOff(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("PowerOff", 1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActionPowerOff(config)
		assert.NoError(t, err)
	})
}
func TestTaskActionsPowerOn(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("PowerOn", 1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActionPowerOn(config)
		assert.NoError(t, err)
	})

}
func TestTaskActionsReboot(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("Reboot", 1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActionReboot(config)
		assert.NoError(t, err)
	})
}

func TestTaskActionsRebuildByImageID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("RebuildByImageID", 1, 2).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Ankr.Set(config.NS, dccncli.ArgImage, "2")

		err := RunTaskActionRebuild(config)
		assert.NoError(t, err)

		assert.True(t, tm.taskActions.AssertExpectations(t))
	})
}

func TestTaskActionsRebuildByImageSlug(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("RebuildByImageSlug", 1, "slug").Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Ankr.Set(config.NS, dccncli.ArgImage, "slug")

		err := RunTaskActionRebuild(config)
		assert.NoError(t, err)

		assert.True(t, tm.taskActions.AssertExpectations(t))
	})

}
func TestTaskActionsRename(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("Rename", 1, "name").Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Ankr.Set(config.NS, dccncli.ArgTaskName, "name")

		err := RunTaskActionRename(config)
		assert.NoError(t, err)
	})
}

func TestTaskActionsResize(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("Resize", 1, "1gb", true).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Ankr.Set(config.NS, dccncli.ArgSizeSlug, "1gb")
		config.Ankr.Set(config.NS, dccncli.ArgResizeDisk, true)

		err := RunTaskActionResize(config)
		assert.NoError(t, err)
	})
}

func TestTaskActionsRestore(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("Restore", 1, 2).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Ankr.Set(config.NS, dccncli.ArgImageID, 2)

		err := RunTaskActionRestore(config)
		assert.NoError(t, err)
	})
}

func TestTaskActionsShutdown(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("Shutdown", 1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActionShutdown(config)
		assert.NoError(t, err)
	})
}

func TestTaskActionsSnapshot(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.taskActions.On("Snapshot", 1, "name").Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Ankr.Set(config.NS, dccncli.ArgSnapshotName, "name")

		err := RunTaskActionSnapshot(config)
		assert.NoError(t, err)
	})
}
