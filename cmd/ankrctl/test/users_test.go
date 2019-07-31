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

package test

import (
	ankrctl "github.com/Ankr-network/dccn-cli"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

const MockResultSuccess = "Success"

var(
	CorrectUserEmail = "test12345@mailinator.com"
	CorrectPassword = "test12345"
	CorrectUserName = "test12345"
	lc = ankrctl.NewLiveCommand("../../../build/ankrctl_linux_amd64")
)

func TestRunUserLogin(t *testing.T) {

	t.Log("user login test ...")

	// case 1: correct input
	loginRes, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(loginRes))
	assert.True(t, strings.Contains(string(loginRes), MockResultSuccess))

	// case 2: invalid inputs
	_, err_invalid := lc.Run( "user", "login", "--email", "", "--password", "")
	if err_invalid == nil {
		t.Error(err_invalid)
	}
	t.Log("Cannot login successfully for invalid email or password")

}

func TestRunUserLogout(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// wait for status changed
	time.Sleep(2 * time.Second)

	// test user logout
	t.Log("user logout test ...")
	logoutRes, err := lc.Run( "user", "logout")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(logoutRes))
	assert.True(t, strings.Contains(string(logoutRes), MockResultSuccess))

}

func TestRunUserDetail(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// wait for status changed
	time.Sleep(2 * time.Second)

	// user detail test
	t.Log("user detail test ...")
	detailRes, err := lc.Run( "user", "detail")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(detailRes))
	assert.True(t, strings.Contains(string(detailRes), "Name"))
	assert.True(t, strings.Contains(string(detailRes), "Email"))
	assert.True(t, strings.Contains(string(detailRes), "Status"))
}

func TestRunUserUpdate(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// wait for status changed
	time.Sleep(2 * time.Second)

	// user update test
	t.Log("user update test ...")
	updateRes, err := lc.Run( "user", "update", CorrectUserEmail, "--update-key", "Name", "--update-value", "user_name_update_test")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(updateRes))
	assert.True(t, strings.Contains(string(updateRes), MockResultSuccess))

	// wait for status changed
	time.Sleep(2 * time.Second)

	// check the update result
	_, err_update := lc.Run( "user", "detail")
	if err_update != nil {
		t.Error(err_update)
	}

	// recovery
	lc.Run( "user", "update", CorrectUserEmail, "--update-key", "Name", "--update-value", CorrectUserName)

}

func TestRunUserChangePassword(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// wait for status changed
	time.Sleep(2 * time.Second)

	// user change password test
	t.Log("user change password test ...")
	changePasswordRes, err := lc.Run( "user", "change-password", CorrectUserEmail, "--old-password", CorrectPassword, "--new-password", "ankr12345")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(changePasswordRes))
	assert.True(t, strings.Contains(string(changePasswordRes), MockResultSuccess))

	// wait for status changed
	time.Sleep(2 * time.Second)

	// use logincli api to test
	_, err_change_password := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", "ankr12345")
	if err_change_password != nil {
		t.Error(err_change_password)
	}

	// recovery
	lc.Run( "user", "change-password", CorrectUserEmail, "--old-password", "ankr12345", "--new-password", CorrectPassword)

}

func TestRunUserChangeEmail(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// wait for status changed
	time.Sleep(2 * time.Second)

	// user change email test
	t.Log("user change email test ...")
	changeEmailRes, err := lc.Run( "user", "email-change", "ankrtestfake@mailinator.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(changeEmailRes))
	assert.True(t, strings.Contains(string(changeEmailRes), "Requested"))
}

func TestRunUserForgotPassword(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// wait for status changed
	time.Sleep(2 * time.Second)

	// user change email test
	t.Log("user forgot password test ...")
	ForgotPasswordRes, err := lc.Run( "user", "forgot-password", CorrectUserEmail)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(ForgotPasswordRes))
	assert.True(t, strings.Contains(string(ForgotPasswordRes), "Requested"))
}