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
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	ankrctl "github.com/Ankr-network/dccn-cli"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	"github.com/stretchr/testify/assert"
)

const (
	MockResultSuccess = "Success"
	MockTaskid        = "100"
	MockTaskName      = "task"
	MockTaskImage     = "nginx:1.12"
	MockReplica       = "2"
	MockUpdateImage   = "nginx:1.13"
	MockUpdateReplica = "3"
	MockTaskType      = "Deploy"
)

type mail struct {
	From    string `json:"f"`
	Subject string `json:"s"`
	HTML    string `json:"html"`
	Text    string `json:"text"`
}

type msg struct {
	UID string `json:"uid"`
}

type inbox struct {
	Msgs []msg `json:"msgs"`
}

func TestMockCommand_Run(t *testing.T) {

	rand.Seed(time.Now().UnixNano())
	charsA := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	charsa := []rune("abcdefghijklmnopqrstuvwxyz")
	nums := []rune("0123456789")
	var b strings.Builder
	var c strings.Builder
	for i := 0; i < 3; i++ {
		b.WriteRune(charsA[rand.Intn(len(charsA))])
		b.WriteRune(charsa[rand.Intn(len(charsa))])
		b.WriteRune(nums[rand.Intn(len(nums))])
		c.WriteRune(charsa[rand.Intn(len(charsa))])
		c.WriteRune(charsa[rand.Intn(len(charsa))])
		c.WriteRune(charsa[rand.Intn(len(charsa))])
	}

	MockPassword := b.String()
	MockUserName := "test" + c.String()

	var url = os.Getenv("URL_BRANCH")
	fmt.Println("url: " + url + "\n")

	lc := ankrctl.NewLiveCommand("go")

	MockUserEmail := MockUserName + "@mailinator.com"

	fmt.Println("user register test..")
	registerRes, err := lc.Run("run", "main.go", "user", "register", MockUserName,
		"--email", MockUserEmail, "--password", MockPassword,
		"-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(registerRes))
	assert.True(t, strings.Contains(string(registerRes), MockResultSuccess))

	MockUserName = "testabcd1234"
	MockUserEmail = "testabcd1234@mailinator.com"
	MockPassword = "abcd1234"
	fmt.Println("user login test..")
	loginRes, err := lc.Run("run", "main.go", "user", "login",
		MockUserEmail, "--password", MockPassword, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(loginRes))
	assert.True(t, strings.Contains(string(loginRes), MockResultSuccess))

	fmt.Println("compute dc list test..")
	dcList, err := lc.Run("run", "main.go", "compute", "dc", "list", "-u", url)
	fmt.Println(string(dcList))
	if err != nil {
		t.Error(err.Error())
	}
	assert.True(t, len(string(dcList)) > 0)
	dcInfo := strings.Split(string(dcList), "\n")
	dcid := ""
	if len(dcInfo) < 2 {
		t.Error("no dc available..")
	} else {
		dcid = strings.Fields(dcInfo[1])[0]
	}

	fmt.Println("compute task create test..")
	taskCreate, err := lc.Run("run", "main.go", "compute", "task", "create",
		MockTaskName, "--image", MockTaskImage, "--dc-name", dcid, "--type", MockTaskType,
		"--replica", MockReplica, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(taskCreate))
	assert.True(t, len(string(taskCreate)) > 0)
	assert.True(t, strings.Contains(string(taskCreate), MockResultSuccess))
	id := strings.Fields(string(taskCreate))[1]
	assert.True(t, len(id) > 0)

	fmt.Println("compute task list test..")
	taskList, err := lc.Run("run", "main.go", "compute", "task", "list", "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(taskList))
	assert.True(t, len(string(taskList)) > 1)
	taskInfo := strings.Split(string(taskList), "\n")
	taskFound := false
	for _, task := range taskInfo {
		if task != "" && id == strings.Fields(task)[0] {
			assert.Equal(t, strings.Fields(task)[1], MockTaskName)
			assert.Equal(t, strings.Fields(task)[3], MockTaskImage)
			assert.Equal(t, strings.Fields(task)[6], MockReplica)
			taskFound = true
		}
	}
	assert.True(t, taskFound)

	time.Sleep(5 * time.Second)

	fmt.Println("compute task update test..")
	taskUpdate, err := lc.Run("run", "main.go", "compute", "task", "update", id,
		"--image", MockUpdateImage, "--replica", MockUpdateReplica, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(taskUpdate))
	assert.True(t, strings.Contains(string(taskUpdate), MockResultSuccess))

	time.Sleep(5 * time.Second)

	fmt.Println("compute task list after update test..")
	taskUpdateList, err := lc.Run("run", "main.go", "compute", "task", "list", "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(taskUpdateList))
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

	fmt.Println("compute task cancel test..")
	taskCancel, err := lc.Run("run", "main.go", "compute", "task", "cancel", "-f", id, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(taskCancel))
	assert.True(t, len(string(taskCancel)) > 0)
	assert.True(t, strings.Contains(string(taskCancel), MockResultSuccess))

	time.Sleep(5 * time.Second)

	fmt.Println("compute task list after cancel test..")
	taskCancelList, err := lc.Run("run", "main.go", "compute", "task", "list", "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(taskCancelList))
	assert.True(t, len(string(taskList)) > 1)
	taskCancelInfo := strings.Split(string(taskCancelList), "\n")
	taskCancelFound := false
	for _, task := range taskCancelInfo {
		if task != "" && id == strings.Fields(task)[0] {
			assert.True(t, strings.Contains(task, common_proto.TaskStatus_CANCELLED.String()))
			taskCancelFound = true
		}
	}
	assert.True(t, taskCancelFound)

	fmt.Println("compute task purge test..")
	taskPurge, err := lc.Run("run", "main.go", "compute", "task", "purge", "-f", id, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(taskPurge))
	assert.True(t, len(string(taskPurge)) > 0)
	assert.True(t, strings.Contains(string(taskPurge), MockResultSuccess))

	fmt.Println("compute task list after purge test..")
	taskPurgeList, err := lc.Run("run", "main.go", "compute", "task", "list", "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(taskPurgeList))
	assert.True(t, len(string(taskList)) > 1)
	taskPurgeInfo := strings.Split(string(taskPurgeList), "\n")
	taskPurgeFound := false
	for _, task := range taskPurgeInfo {
		if task != "" && id == strings.Fields(task)[0] {
			taskPurgeFound = true
		}
	}
	assert.False(t, taskPurgeFound)

}
