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

	"github.com/Ankr-network/ankrctl/types"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	"github.com/stretchr/testify/assert"
)

const (
	MockResultSuccess = "Success"
	MockAppid         = "100"
	MockAppName       = "app"
	MockAppImage      = "nginx:1.12"
	MockReplica       = "2"
	MockUpdateImage   = "nginx:1.13"
	MockUpdateReplica = "3"
	MockAppType       = "Deploy"
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

	lc := types.NewLiveCommand("go")

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

	fmt.Println("compute app create test..")
	appCreate, err := lc.Run("run", "main.go", "compute", "app", "create",
		MockAppName, "--image", MockAppImage, "--dc-name", dcid, "--type", MockAppType,
		"--replica", MockReplica, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(appCreate))
	assert.True(t, len(string(appCreate)) > 0)
	assert.True(t, strings.Contains(string(appCreate), MockResultSuccess))
	id := strings.Fields(string(appCreate))[1]
	assert.True(t, len(id) > 0)

	fmt.Println("compute app list test..")
	appList, err := lc.Run("run", "main.go", "compute", "app", "list", "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(appList))
	assert.True(t, len(string(appList)) > 1)
	appInfo := strings.Split(string(appList), "\n")
	appFound := false
	for _, app := range appInfo {
		if app != "" && id == strings.Fields(app)[0] {
			assert.Equal(t, strings.Fields(app)[1], MockAppName)
			assert.Equal(t, strings.Fields(app)[3], MockAppImage)
			assert.Equal(t, strings.Fields(app)[6], MockReplica)
			appFound = true
		}
	}
	assert.True(t, appFound)

	time.Sleep(5 * time.Second)

	fmt.Println("compute app update test..")
	appUpdate, err := lc.Run("run", "main.go", "compute", "app", "update", id,
		"--image", MockUpdateImage, "--replica", MockUpdateReplica, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(appUpdate))
	assert.True(t, strings.Contains(string(appUpdate), MockResultSuccess))

	time.Sleep(5 * time.Second)

	fmt.Println("compute app list after update test..")
	appUpdateList, err := lc.Run("run", "main.go", "compute", "app", "list", "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(appUpdateList))
	assert.True(t, len(string(appList)) > 1)
	appUpdateInfo := strings.Split(string(appUpdateList), "\n")
	appUpdateFound := false
	for _, app := range appUpdateInfo {
		if app != "" && id == strings.Fields(app)[0] {
			assert.Equal(t, strings.Fields(app)[3], MockUpdateImage)
			assert.Equal(t, strings.Fields(app)[6], MockUpdateReplica)
			appUpdateFound = true
		}
	}
	assert.True(t, appUpdateFound)

	fmt.Println("compute app cancel test..")
	appCancel, err := lc.Run("run", "main.go", "compute", "app", "cancel", "-f", id, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(appCancel))
	assert.True(t, len(string(appCancel)) > 0)
	assert.True(t, strings.Contains(string(appCancel), MockResultSuccess))

	time.Sleep(5 * time.Second)

	fmt.Println("compute app list after cancel test..")
	appCancelList, err := lc.Run("run", "main.go", "compute", "app", "list", "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(appCancelList))
	assert.True(t, len(string(appList)) > 1)
	appCancelInfo := strings.Split(string(appCancelList), "\n")
	appCancelFound := false
	for _, app := range appCancelInfo {
		if app != "" && id == strings.Fields(app)[0] {
			assert.True(t, strings.Contains(app, common_proto.AppStatus_CANCELLED.String()))
			appCancelFound = true
		}
	}
	assert.True(t, appCancelFound)

	fmt.Println("compute app purge test..")
	appPurge, err := lc.Run("run", "main.go", "compute", "app", "purge", "-f", id, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(appPurge))
	assert.True(t, len(string(appPurge)) > 0)
	assert.True(t, strings.Contains(string(appPurge), MockResultSuccess))

	fmt.Println("compute app list after purge test..")
	appPurgeList, err := lc.Run("run", "main.go", "compute", "app", "list", "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(appPurgeList))
	assert.True(t, len(string(appList)) > 1)
	appPurgeInfo := strings.Split(string(appPurgeList), "\n")
	appPurgeFound := false
	for _, app := range appPurgeInfo {
		if app != "" && id == strings.Fields(app)[0] {
			appPurgeFound = true
		}
	}
	assert.False(t, appPurgeFound)

}
