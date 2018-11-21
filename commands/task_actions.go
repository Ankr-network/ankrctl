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
	"strconv"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/spf13/cobra"
)

type actionFn func(das do.TaskActionsService) (*do.Action, error)

func performAction(c *CmdConfig, fn actionFn) error {
	das := c.TaskActions()

	a, err := fn(das)
	if err != nil {
		return err
	}

	wait, err := c.Ankr.GetBool(c.NS, dccncli.ArgCommandWait)
	if err != nil {
		return err
	}

	if wait {
		a, err = actionWait(c, a.ID, 5)
		if err != nil {
			return err
		}

	}

	item := &displayers.Action{Actions: do.Actions{*a}}
	return c.Display(item)
}

// TaskAction creates the task-action command.
func TaskAction() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "task-action",
			Aliases: []string{"da"},
			Short:   "task action commands",
			Long:    "task-action is used to access task action commands",
		},
	}

	cmdTaskActionGet := CmdBuilder(cmd, RunTaskActionGet, "get <task-id>", "get task action", Writer,
		aliasOpt("g"), displayerType(&displayers.Action{}), docCategories("task"))
	AddIntFlag(cmdTaskActionGet, dccncli.ArgActionID, "", 0, "Action ID", requiredOpt())

	cmdTaskActionEnableBackups := CmdBuilder(cmd, RunTaskActionEnableBackups,
		"enable-backups <task-id>", "enable backups", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionEnableBackups, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionDisableBackups := CmdBuilder(cmd, RunTaskActionDisableBackups,
		"disable-backups <task-id>", "disable backups", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionDisableBackups, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionReboot := CmdBuilder(cmd, RunTaskActionReboot,
		"reboot <task-id>", "reboot task", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionReboot, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionPowerCycle := CmdBuilder(cmd, RunTaskActionPowerCycle,
		"power-cycle <task-id>", "power cycle task", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionPowerCycle, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionShutdown := CmdBuilder(cmd, RunTaskActionShutdown,
		"shutdown <task-id>", "shutdown task", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionShutdown, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionPowerOff := CmdBuilder(cmd, RunTaskActionPowerOff,
		"power-off <task-id>", "power off task", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionPowerOff, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionPowerOn := CmdBuilder(cmd, RunTaskActionPowerOn,
		"power-on <task-id>", "power on task", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionPowerOn, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionPasswordReset := CmdBuilder(cmd, RunTaskActionPasswordReset,
		"password-reset <task-id>", "password reset task", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionPasswordReset, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionEnableIPv6 := CmdBuilder(cmd, RunTaskActionEnableIPv6,
		"enable-ipv6 <task-id>", "enable ipv6", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionEnableIPv6, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionEnablePrivateNetworking := CmdBuilder(cmd, RunTaskActionEnablePrivateNetworking,
		"enable-private-networking <task-id>", "enable private networking", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionEnablePrivateNetworking, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionRestore := CmdBuilder(cmd, RunTaskActionRestore,
		"restore <task-id>", "restore backup", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddIntFlag(cmdTaskActionRestore, dccncli.ArgImageID, "", 0, "Image ID", requiredOpt())
	AddBoolFlag(cmdTaskActionRestore, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionResize := CmdBuilder(cmd, RunTaskActionResize,
		"resize <task-id>", "resize task", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddBoolFlag(cmdTaskActionResize, dccncli.ArgResizeDisk, "", false, "Resize disk")
	AddStringFlag(cmdTaskActionResize, dccncli.ArgSizeSlug, "", "", "New size")
	AddBoolFlag(cmdTaskActionResize, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionRebuild := CmdBuilder(cmd, RunTaskActionRebuild,
		"rebuild <task-id>", "rebuild task", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddStringFlag(cmdTaskActionRebuild, dccncli.ArgImage, "", "", "Image ID or Slug", requiredOpt())
	AddBoolFlag(cmdTaskActionRebuild, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionRename := CmdBuilder(cmd, RunTaskActionRename,
		"rename <task-id>", "rename task", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddStringFlag(cmdTaskActionRename, dccncli.ArgTaskName, "", "", "Task name", requiredOpt())
	AddBoolFlag(cmdTaskActionRename, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionChangeKernel := CmdBuilder(cmd, RunTaskActionChangeKernel,
		"change-kernel <task-id>", "change kernel", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddIntFlag(cmdTaskActionChangeKernel, dccncli.ArgKernelID, "", 0, "Kernel ID", requiredOpt())
	AddBoolFlag(cmdTaskActionChangeKernel, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdTaskActionSnapshot := CmdBuilder(cmd, RunTaskActionSnapshot,
		"snapshot <task-id>", "snapshot task", Writer,
		displayerType(&displayers.Action{}), docCategories("task"))
	AddStringFlag(cmdTaskActionSnapshot, dccncli.ArgSnapshotName, "", "", "Snapshot name", requiredOpt())
	AddBoolFlag(cmdTaskActionSnapshot, dccncli.ArgCommandWait, "", false, "Wait for action to complete")

	return cmd
}

// RunTaskActionGet returns a task action by id.
func RunTaskActionGet(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		taskID, err := strconv.Atoi(c.Args[0])
		if err != nil {
			return nil, err
		}

		actionID, err := c.Ankr.GetInt(c.NS, dccncli.ArgActionID)
		if err != nil {
			return nil, err
		}

		a, err := das.Get(taskID, actionID)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionEnableBackups disables backups for a task.
func RunTaskActionEnableBackups(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])
		if err != nil {
			return nil, err
		}

		a, err := das.EnableBackups(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionDisableBackups disables backups for a task.
func RunTaskActionDisableBackups(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])
		if err != nil {
			return nil, err
		}

		a, err := das.DisableBackups(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionReboot reboots a task.
func RunTaskActionReboot(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])
		if err != nil {
			return nil, err
		}

		a, err := das.Reboot(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionPowerCycle power cycles a task.
func RunTaskActionPowerCycle(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.PowerCycle(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionShutdown shuts a task down.
func RunTaskActionShutdown(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		a, err := das.Shutdown(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionPowerOff turns task power off.
func RunTaskActionPowerOff(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.PowerOff(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionPowerOn turns task power on.
func RunTaskActionPowerOn(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.PowerOn(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionPasswordReset resets the task root password.
func RunTaskActionPasswordReset(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.PasswordReset(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionEnableIPv6 enables IPv6 for a task.
func RunTaskActionEnableIPv6(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.EnableIPv6(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionEnablePrivateNetworking enables private networking for a task.
func RunTaskActionEnablePrivateNetworking(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.EnablePrivateNetworking(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionRestore restores a task using an image id.
func RunTaskActionRestore(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		image, err := c.Ankr.GetInt(c.NS, dccncli.ArgImageID)
		if err != nil {
			return nil, err
		}

		a, err := das.Restore(id, image)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionResize resizesx a task giving a size slug and
// optionally expands the disk.
func RunTaskActionResize(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		size, err := c.Ankr.GetString(c.NS, dccncli.ArgSizeSlug)
		if err != nil {
			return nil, err
		}

		disk, err := c.Ankr.GetBool(c.NS, dccncli.ArgResizeDisk)
		if err != nil {
			return nil, err
		}

		a, err := das.Resize(id, size, disk)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionRebuild rebuilds a task using an image id or slug.
func RunTaskActionRebuild(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		image, err := c.Ankr.GetString(c.NS, dccncli.ArgImage)
		if err != nil {
			return nil, err
		}

		var a *do.Action
		if i, aerr := strconv.Atoi(image); aerr == nil {
			a, err = das.RebuildByImageID(id, i)
		} else {
			a, err = das.RebuildByImageSlug(id, image)
		}
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionRename renames a task.
func RunTaskActionRename(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		name, err := c.Ankr.GetString(c.NS, dccncli.ArgTaskName)
		if err != nil {
			return nil, err
		}

		a, err := das.Rename(id, name)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionChangeKernel changes the kernel for a task.
func RunTaskActionChangeKernel(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		kernel, err := c.Ankr.GetInt(c.NS, dccncli.ArgKernelID)
		if err != nil {
			return nil, err
		}

		a, err := das.ChangeKernel(id, kernel)
		return a, err
	}

	return performAction(c, fn)
}

// RunTaskActionSnapshot creates a snapshot for a task.
func RunTaskActionSnapshot(c *CmdConfig) error {
	fn := func(das do.TaskActionsService) (*do.Action, error) {
		if len(c.Args) != 1 {
			return nil, dccncli.NewMissingArgsErr(c.NS)
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		name, err := c.Ankr.GetString(c.NS, dccncli.ArgSnapshotName)
		if err != nil {
			return nil, err
		}

		a, err := das.Snapshot(id, name)
		return a, err
	}

	return performAction(c, fn)
}
