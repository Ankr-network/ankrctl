package commands

import (
	"fmt"
	ankrctl "github.com/Ankr-network/dccn-cli"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
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

var MockUserName string
var MockPassword string
var MockUserEmail string
var url = os.Getenv("URL_BRANCH")
var MockRegisterCode string

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
	MockUserName = "test" + c.String()
	MockPassword = b.String()
	MockUserEmail = MockUserName + "@mailinator.com"
	fmt.Println("url: " + url + "\n")
}


func TestRunUserRegister(t *testing.T) {

	lc := ankrctl.NewLiveCommand("go")
	fmt.Println("user register test..")

	registerRes, err := lc.Run("run", "main.go", "user", "register", MockUserName,
		"--email", MockUserEmail, "--password", MockPassword, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(registerRes))
	assert.True(t, strings.Contains(string(registerRes), MockResultSuccess))
}


func TestRunUserLogin(t *testing.T) {
	// use a valid account
	MockUserName = "testabcd1234"
	MockUserEmail = "testabcd1234@mailinator.com"
	MockPassword = "abcd1234"

	lc := ankrctl.NewLiveCommand("go")
	fmt.Println("user login test..")

	loginRes, err := lc.Run("run", "main.go", "user", "login",
		MockUserEmail, "--password", MockPassword, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(loginRes))
	assert.True(t, strings.Contains(string(loginRes), MockResultSuccess))
}

func TestRunUserDetail(t *testing.T) {
	lc := ankrctl.NewLiveCommand("go")
	fmt.Println("user detail test..")

	// login at first
	_, err := lc.Run("run", "main.go", "user", "login", MockUserEmail, "--password", MockPassword, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}

	// test user detail api
	detailRes, err := lc.Run("run", "main.go", "user", "detail")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(detailRes))
	assert.True(t, strings.Contains(string(detailRes), MockResultSuccess))

}

func TestRunUserConfirmRegistration(t *testing.T) {

	lc := ankrctl.NewLiveCommand("go")
	MockRegisterCode = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTMyNDI0MTgsImp0aSI6IjEzNGZjZjg4LTg4NTUtNGM3Zi05NmZmLWU1OTU2ZDg4OTZkMCIsImlzcyI6ImFua3IubmV0d29yayJ9.n6InDb5RhwOduTc9-Vo-1rS6CVhqYAF2AnwEnn2B0LE"
	fmt.Println("user confirm registration test..")

	confirmRes, err := lc.Run("run", "main.go", "user", "confirm-registration",
		MockUserEmail, "--register-code", MockRegisterCode, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(confirmRes))
	assert.True(t, strings.Contains(string(confirmRes), MockResultSuccess))
}

func TestRunUserUpdate(t *testing.T) {
	lc := ankrctl.NewLiveCommand("go")
	fmt.Println("user update test..")

	updateRes, err := lc.Run("run", "main.go", "user", "update",
		MockUserEmail, "--update-key", MockRegisterCode, "-u", url)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(updateRes))
	assert.True(t, strings.Contains(string(updateRes), MockResultSuccess))

}

func TestRunUserLogout(t *testing.T) {

}