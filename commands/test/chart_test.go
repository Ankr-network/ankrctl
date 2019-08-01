package test

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var (
	chartUploadFile = "../dccn-appmgr/examples/test/wordpress-5.7.1.tgz"
	chartUploadVersion = "5.7.1"
	chartUploadName = "wordpress"
)

func TestRunChartList(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// chart list test

	chartListRes, err := lc.Run( "chart", "list", "--list-repo", "stable")
	if err != nil {
		t.Error(err)
	}else{
	t.Log(string(chartListRes))
	t.Log("list charts successfully")
	assert.True(t, strings.Contains(string(chartListRes), "Name"))
	assert.True(t, strings.Contains(string(chartListRes), "Repo"))
	}

	// wait for status changed
	time.Sleep(2 * time.Second)
}


func TestRunChartUpload(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// chart upload test
	t.Log("chart upload test ...")
	chartUploadRes, err := lc.Run( "chart", "upload", chartUploadName, "--upload-file", chartUploadFile, "--upload-version", chartUploadVersion)
	if err != nil {
		t.Error(err)
	}else{
	t.Log(string(chartUploadRes))
	assert.True(t, strings.Contains(string(chartUploadRes), "success"))
	t.Log("upload chart successfully")
	}

	// wait for status changed
	time.Sleep(5 * time.Second)

	// delete the chart uploaded
	lc.Run("chart", "delete", chartUploadName, "--delete-version", chartUploadVersion, "-f")

	// wait for status changed
	time.Sleep(2 * time.Second)

}

