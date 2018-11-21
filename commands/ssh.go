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
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/dccn-cli/pkg/ssh"
)

var (
	sshHostRE = regexp.MustCompile(`^((?P<m1>\w+)@)?(?P<m2>.*?)(:(?P<m3>\d+))?$`)
)

// SSH creates the ssh commands heirarchy
func SSH(parent *Command) *Command {
	usr, err := user.Current()
	checkErr(err)

	path := filepath.Join(usr.HomeDir, ".ssh", "id_rsa")

	cmdSSH := CmdBuilder(parent, RunSSH, "ssh <task-id|host>", "ssh to task", Writer,
		docCategories("task"))
	AddStringFlag(cmdSSH, dccncli.ArgSSHUser, "", "root", "ssh user")
	AddStringFlag(cmdSSH, dccncli.ArgsSSHKeyPath, "", path, "path to private ssh key")
	AddIntFlag(cmdSSH, dccncli.ArgsSSHPort, "", 22, "port sshd is running on")
	AddBoolFlag(cmdSSH, dccncli.ArgsSSHAgentForwarding, "", false, "enable ssh agent forwarding")
	AddBoolFlag(cmdSSH, dccncli.ArgsSSHPrivateIP, "", false, "ssh to private ip instead of public ip")
	AddStringFlag(cmdSSH, dccncli.ArgSSHCommand, "", "", "command to execute")

	return cmdSSH
}

// RunSSH finds a task to ssh to given input parameters (name or id).
func RunSSH(c *CmdConfig) error {
	if len(c.Args) == 0 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	taskID := c.Args[0]

	if taskID == "" {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	user, err := c.Ankr.GetString(c.NS, dccncli.ArgSSHUser)
	if err != nil {
		return err
	}

	keyPath, err := c.Ankr.GetString(c.NS, dccncli.ArgsSSHKeyPath)
	if err != nil {
		return err
	}

	port, err := c.Ankr.GetInt(c.NS, dccncli.ArgsSSHPort)
	if err != nil {
		return err
	}

	var opts = make(ssh.Options)
	opts[dccncli.ArgsSSHAgentForwarding], err = c.Ankr.GetBool(c.NS, dccncli.ArgsSSHAgentForwarding)
	if err != nil {
		return err
	}

	opts[dccncli.ArgSSHCommand], err = c.Ankr.GetString(c.NS, dccncli.ArgSSHCommand)
	if err != nil {
		return nil
	}

	privateIPChoice, err := c.Ankr.GetBool(c.NS, dccncli.ArgsSSHPrivateIP)
	if err != nil {
		return err
	}

	var task *do.Task

	ds := c.Tasks()
	if id, err := strconv.Atoi(taskID); err == nil {
		// taskID is an integer

		doTask, err := ds.Get(id)
		if err != nil {
			return err
		}

		task = doTask
	} else {
		// taskID is a string
		tasks, err := ds.List()
		if err != nil {
			return err
		}

		shi := extractHostInfo(taskID)

		if shi.user != "" {
			user = shi.user
		}

		if i, err := strconv.Atoi(shi.port); shi.port != "" && err != nil {
			port = i
		}

		for _, d := range tasks {
			if d.Name == shi.host {
				task = &d
				break
			}
			if strconv.Itoa(d.ID) == shi.host {
				task = &d
				break
			}
		}

		if task == nil {
			return errors.New("could not find task")
		}

	}

	if user == "" {
		user = defaultSSHUser(task)
	}

	ip, err := privateIPElsePub(task, privateIPChoice)
	if err != nil {
		return err
	}

	if ip == "" {
		return errors.New("could not find task address")
	}

	runner := c.Ankr.SSH(user, ip, keyPath, port, opts)
	return runner.Run()
}

func defaultSSHUser(task *do.Task) string {
	slug := strings.ToLower(task.Image.Slug)
	if strings.Contains(slug, "coreos") {
		return "core"
	}

	return "root"
}

type sshHostInfo struct {
	user string
	host string
	port string
}

func extractHostInfo(in string) sshHostInfo {
	m := sshHostRE.FindStringSubmatch(in)
	r := map[string]string{}
	for i, n := range sshHostRE.SubexpNames() {
		r[n] = m[i]
	}

	return sshHostInfo{
		user: r["m1"],
		host: r["m2"],
		port: r["m3"],
	}
}

func privateIPElsePub(task *do.Task, choice bool) (string, error) {
	if choice {
		return task.PrivateIPv4()
	}
	return task.PublicIPv4()
}
