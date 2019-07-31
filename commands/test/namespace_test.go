package test

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRunNamespaceCreate(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// namespace create test
	t.Log("namespace create test ...")
	chartListRes, err := lc.Run( "chart", "list", "--list-repo", "stable")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(chartListRes))
	t.Log("list charts successfully")
	assert.True(t, strings.Contains(string(chartListRes), "Name"))

}