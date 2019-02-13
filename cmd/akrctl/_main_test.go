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
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	akrctl "github.com/Ankr-network/dccn-cli"

	"github.com/stretchr/testify/assert"
)

const (
	MockUserName      = "testuser"
	MockUserEmail     = "@mailinator.com"
	MockPassword      = "11111"
	MockResultSuccess = "Success"
	MockTaskid        = "100"
	MockTaskName      = "task"
	MockTaskImage     = "nginx:1.12"
	MockReplica       = "2"
	MockUpdateImage   = "nginx:1.13"
	MockUpdateReplica = "3"
	MockTaskType      = "Deploy"
)

func TestMockCommand_Run(t *testing.T) {

	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	randstr := b.String() // E.g. "ExcbsVQs"

	var url = os.Getenv("URL_BRANCH")
	t.Log("url: " + url + "\n")

	lc := akrctl.NewLiveCommand("go")

	t.Log("user register test..")
	registerRes, err := lc.Run("run", "main.go", "user", "register", MockUserName, "--email", MockUserName+"_"+randstr+MockUserEmail, "--password", MockPassword, "-u", url)
	t.Log(string(registerRes))
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(registerRes), MockResultSuccess))

	t.Log("user login test..")
	loginRes, err := lc.Run("run", "main.go", "user", "login", MockUserName+"_"+randstr+MockUserEmail, "--password", MockPassword, "-u", url)
	t.Log(string(loginRes))
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(loginRes), MockResultSuccess))

	t.Log("compute dc list test..")
	dcList, err := lc.Run("run", "main.go", "compute", "dc", "list", "-u", url)
	t.Log(string(dcList))
	assert.NoError(t, err)
	assert.True(t, len(string(dcList)) > 0)
	dcInfo := strings.Split(string(dcList), "\n")
	dcid := ""
	if len(dcInfo) < 2 {
		t.Error("no dc available..")
	} else {
		dcid = strings.Fields(dcInfo[1])[0]
	}

	t.Log("compute task create test..")
	taskCreate, err := lc.Run("run", "main.go", "compute", "task", "create",
		MockTaskName, "--image", MockTaskImage, "--dc-id", dcid, "--type", MockTaskType, "--replica", MockReplica,
		"-u", url)
	t.Log(string(taskCreate))
	assert.NoError(t, err)
	assert.True(t, len(string(taskCreate)) > 0)
	assert.True(t, strings.Contains(string(taskCreate), MockResultSuccess))
	id := strings.Fields(string(taskCreate))[1]
	assert.True(t, len(id) > 0)

	time.Sleep(10 * time.Second)

	t.Log("compute task list test..")
	taskList, err := lc.Run("run", "main.go", "compute", "task", "list", "-u", url)
	t.Log(string(taskList))
	assert.NoError(t, err)
	assert.True(t, len(string(taskList)) > 1)
	taskInfo := strings.Split(string(taskList), "\n")
	taskFound := false
	for _, task := range taskInfo {
		if task != "" && id == strings.Fields(task)[0] {
			assert.Equal(t, strings.Fields(task)[1], MockTaskName)
			assert.Equal(t, strings.Fields(task)[3], MockTaskImage)
			assert.Equal(t, strings.Fields(task)[6], MockReplica)
			assert.Equal(t, strings.Fields(task)[7], dcid)
			taskFound = true
		}
	}
	assert.True(t, taskFound)

	t.Log("compute task update test..")
	taskUpdate, err := lc.Run("run", "main.go", "compute", "task", "update", id,
		"--image", MockUpdateImage, "--replica", MockUpdateReplica, "-u", url)
	t.Log(string(taskUpdate))
	assert.NoError(t, err)
	assert.True(t, len(string(taskUpdate)) > 0)
	assert.True(t, strings.Contains(string(taskUpdate), MockResultSuccess))
	/*
		t.Log("compute task list after update test..")
		taskUpdateList, err := lc.Run("run", "main.go", "compute", "task", "list", "-u", url)
		t.Log(string(taskUpdateList))
		assert.NoError(t, err)
		assert.True(t, len(string(taskList)) > 1)
		taskUpdateInfo := strings.Split(string(taskUpdateList), "\n")
		taskUpdateFound := false
		for _, task := range taskUpdateInfo {
			if task != "" && id == strings.Fields(task)[0] {
				assert.Equal(t, strings.Fields(task)[3], MockUpdateImage)
				assert.Equal(t, strings.Fields(task)[6], MockUpdateReplica)
				taskUpdateFound = true
			}
		}
		assert.True(t, taskUpdateFound)
	*/
	t.Log("compute task detail test..")
	taskDetail, err := lc.Run("run", "main.go", "compute", "task", "detail", id, "-u", url)
	t.Log(string(taskDetail))
	assert.NoError(t, err)
	assert.True(t, len(string(taskDetail)) > 0)
	assert.True(t, strings.Contains(string(taskDetail), MockResultSuccess))

	t.Log("compute task cancel test..")
	taskCancel, err := lc.Run("run", "main.go", "compute", "task", "cancel", "-f", id, "-u", url)
	t.Log(string(taskCancel))
	assert.NoError(t, err)
	assert.True(t, len(string(taskCancel)) > 0)
	assert.True(t, strings.Contains(string(taskCancel), MockResultSuccess))

	t.Log("compute task list after cancel test..")
	taskCancelList, err := lc.Run("run", "main.go", "compute", "task", "list", "-u", url)
	t.Log(string(taskCancelList))
	assert.NoError(t, err)
	assert.True(t, len(string(taskList)) > 1)
	taskCancelInfo := strings.Split(string(taskCancelList), "\n")
	taskCancelFound := false
	for _, task := range taskCancelInfo {
		if task != "" && id == strings.Fields(task)[0] {
			//assert.Equal(t, strings.Fields(task)[9], common_proto.TaskStatus_CANCELLED)
			taskCancelFound = true
		}
	}
	assert.True(t, taskCancelFound)

	t.Log("compute task purge test..")
	taskPurge, err := lc.Run("run", "main.go", "compute", "task", "purge", "-f", id, "-u", url)
	t.Log(string(taskPurge))
	assert.NoError(t, err)
	assert.True(t, len(string(taskPurge)) > 0)
	assert.True(t, strings.Contains(string(taskPurge), MockResultSuccess))

	/*
		t.Log("compute task list after purge test..")
		taskPurgeList, err := lc.Run("run", "main.go", "compute", "task", "list", "-u", url)
		t.Log(string(taskPurgeList))
		assert.NoError(t, err)
		assert.True(t, len(string(taskList)) > 1)
		taskPurgeInfo := strings.Split(string(taskPurgeList), "\n")
		taskPurgeFound := false
		for _, task := range taskPurgeInfo {
			if task != "" && id == strings.Fields(task)[0] {
				assert.Equal(t, strings.Fields(task)[9], common_proto.TaskStatus_CANCELLED)
				taskPurgeFound = true
			}
		}
		assert.False(t, taskPurgeFound)
	*/
}
