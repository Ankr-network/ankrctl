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
	ankrctl "github.com/Ankr-network/dccn-cli"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
	"time"
)

const (
	MockResultSuccess = "Success"
	MockAppid        = "100"
	MockAppName      = "app"
	MockAppImage     = "nginx:1.12"
	MockReplica       = "2"
	MockUpdateImage   = "nginx:1.13"
	MockUpdateReplica = "3"
	MockAppType      = "Deploy"
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

	lc := ankrctl.NewLiveCommand("./ankrctl_linux_amd64")

	MockUserEmail := MockUserName + "@mailinator.com"

	fmt.Println("user register test..")
	registerRes, err := lc.Run( "user", "register", MockUserName,
		"--email", MockUserEmail, "--password", MockPassword)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(registerRes))
	assert.True(t, strings.Contains(string(registerRes), MockResultSuccess))

}
