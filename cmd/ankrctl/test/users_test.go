package test

import (
	"fmt"
	ankrctl "github.com/Ankr-network/dccn-cli"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
	"time"
)


func TestMockCommand(t *testing.T) {

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

	lc := ankrctl.NewLiveCommand("../../../build/ankrctl_linux_amd64")
	MockUserEmail := "test12345@mailinator.com"
	MockPassword := "test12345"
	fmt.Println("user login test..")
	loginRes, err := lc.Run( "user", "login", "--email", MockUserEmail, "--password", MockPassword)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(loginRes))
	assert.True(t, strings.Contains(string(loginRes), "success"))

}