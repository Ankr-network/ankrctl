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
	"testing"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/pkg/runner"
	"github.com/Ankr-network/dccn-cli/pkg/runner/mocks"
	"github.com/Ankr-network/dccn-cli/pkg/ssh"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestSSHComand(t *testing.T) {
	parent := &Command{
		Command: &cobra.Command{
			Use:   "compute",
			Short: "compute commands",
			Long:  "compute commands are for controlling and managing infrastructure",
		},
	}
	cmd := SSH(parent)
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd)
}

func TestSSH_ID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Get", testTask.ID).Return(&testTask, nil)

		config.Args = append(config.Args, strconv.Itoa(testTask.ID))

		err := RunSSH(config)
		assert.NoError(t, err)
	})
}

func TestSSH_InvalidID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunSSH(config)
		assert.Error(t, err)
	})
}

func TestSSH_UnknownTask(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("List").Return(testTaskList, nil)

		config.Args = append(config.Args, "missing")

		err := RunSSH(config)
		assert.EqualError(t, err, "could not find task")
	})
}

func TestSSH_TaskWithNoPublic(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("List").Return(testPrivateTaskList, nil)

		config.Args = append(config.Args, testPrivateTask.Name)

		err := RunSSH(config)
		assert.EqualError(t, err, "could not find task address")
	})
}

func TestSSH_CustomPort(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		rm := &mocks.Runner{}
		rm.On("Run").Return(nil)

		tc := config.Ankr.(*TestConfig)
		tc.SSHFn = func(user, host, keyPath string, port int, opts ssh.Options) runner.Runner {
			assert.Equal(t, 2222, port)
			return rm
		}

		tm.tasks.On("List").Return(testTaskList, nil)

		config.Ankr.Set(config.NS, dccncli.ArgsSSHPort, "2222")
		config.Args = append(config.Args, testTask.Name)

		err := RunSSH(config)
		assert.NoError(t, err)
	})
}

func TestSSH_CustomUser(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		rm := &mocks.Runner{}
		rm.On("Run").Return(nil)

		tc := config.Ankr.(*TestConfig)
		tc.SSHFn = func(user, host, keyPath string, port int, opts ssh.Options) runner.Runner {
			assert.Equal(t, "foobar", user)
			return rm
		}

		tm.tasks.On("List").Return(testTaskList, nil)

		config.Ankr.Set(config.NS, dccncli.ArgSSHUser, "foobar")
		config.Args = append(config.Args, testTask.Name)

		err := RunSSH(config)
		assert.NoError(t, err)
	})
}

func TestSSH_AgentForwarding(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		rm := &mocks.Runner{}
		rm.On("Run").Return(nil)

		tc := config.Ankr.(*TestConfig)
		tc.SSHFn = func(user, host, keyPath string, port int, opts ssh.Options) runner.Runner {
			assert.Equal(t, true, opts[dccncli.ArgsSSHAgentForwarding])
			return rm
		}

		tm.tasks.On("List").Return(testTaskList, nil)

		config.Ankr.Set(config.NS, dccncli.ArgsSSHAgentForwarding, true)
		config.Args = append(config.Args, testTask.Name)

		err := RunSSH(config)
		assert.NoError(t, err)
	})
}

func TestSSH_CommandExecuting(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		rm := &mocks.Runner{}
		rm.On("Run").Return(nil)

		tc := config.Ankr.(*TestConfig)
		tc.SSHFn = func(user, host, keyPath string, port int, opts ssh.Options) runner.Runner {
			assert.Equal(t, "uptime", opts[dccncli.ArgSSHCommand])
			return rm
		}

		tm.tasks.On("List").Return(testTaskList, nil)
		config.Ankr.Set(config.NS, dccncli.ArgSSHCommand, "uptime")
		config.Args = append(config.Args, testTask.Name)

		err := RunSSH(config)
		assert.NoError(t, err)
	})
}

func Test_extractHostInfo(t *testing.T) {
	cases := []struct {
		s string
		e sshHostInfo
	}{
		{s: "host", e: sshHostInfo{host: "host"}},
		{s: "root@host", e: sshHostInfo{user: "root", host: "host"}},
		{s: "root@host:22", e: sshHostInfo{user: "root", host: "host", port: "22"}},
		{s: "host:22", e: sshHostInfo{host: "host", port: "22"}},
		{s: "dokku@simple-task-02efb9c544", e: sshHostInfo{host: "simple-task-02efb9c544", user: "dokku"}},
	}

	for _, c := range cases {
		i := extractHostInfo(c.s)
		assert.Equal(t, c.e, i)
	}
}
