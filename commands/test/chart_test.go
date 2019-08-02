package test

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var (
	chartUploadVersion = "8.8.8"
	chartUploadName = "chart_update_test"
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


func TestRunChartDownload(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// chart download test
	t.Log("chart download test ...")
	chartDownloadRes, err := lc.Run( "chart", "download", "wordpress", "--download-repo", "stable", "--download-version", "5.6.2")
	if err != nil {
		t.Error(err)
	}else{
		t.Log(string(chartDownloadRes))
		assert.True(t, strings.Contains(string(chartDownloadRes), "Successfully"))
		t.Log("download chart successfully")
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

	// download a chart for upload
	chartDownloadRes, _ := lc.Run( "chart", "download", "wordpress", "--download-repo", "stable", "--download-version", "5.6.2")
	t.Log(string(chartDownloadRes))

	// chart upload test
	t.Log("chart upload test ...")
	chartUploadRes, err := lc.Run( "chart", "upload", chartUploadName, "--upload-file", "/go/src/github.com/Ankr-network/dccn-cli/commands/test/wordpress-5.6.2.tgz", "--upload-version", chartUploadVersion)
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

func TestRunChartDelete(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// download a chart for upload
	lc.Run( "chart", "download", "wordpress", "--download-repo", "stable", "--download-version", "5.6.2")

	// upload a chart for delete test
	lc.Run( "chart", "upload", chartUploadName, "--upload-file", "/go/src/github.com/Ankr-network/dccn-cli/commands/test/wordpress-5.6.2.tgz", "--upload-version", chartUploadVersion)

	// wait for status changed
	time.Sleep(5 * time.Second)

	// chart delete test
	t.Log("chart delete test ...")
	chartDeleteRes, err := lc.Run("chart", "delete", chartUploadName, "--delete-version", chartUploadVersion, "-f")
	if err != nil {
		t.Error(err)
	}else{
		t.Log(string(chartDeleteRes))
		assert.True(t, strings.Contains(string(chartDeleteRes), "success"))
	}

	// wait for status changed
	time.Sleep(2 * time.Second)
}

func TestRunChartDetail(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// chart detail test
	t.Log("chart detail test ...")
	chartDetailRes, err := lc.Run("chart", "detail", "wordpress", "--detail-repo", "stable", "--show-version", "5.6.0")
	if err != nil {
		t.Error(err)
	}else{
		t.Log(string(chartDetailRes))
		assert.True(t, strings.Contains(string(chartDetailRes), "Repo"))
	}

	// wait for status changed
	time.Sleep(2 * time.Second)
}

/*func TestRunChartSaveas(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// chart saveas test
	t.Log("chart saveas test ...")
	chartSaveasRes, err := lc.Run("chart", "saveas", "wordpress-5.6.0", "--delete-version", chartUploadVersion, "-f")
	if err != nil {
		t.Error(err)
	}else{
		t.Log(string(chartSaveasRes))
		assert.True(t, strings.Contains(string(chartSaveasRes), "success"))
	}

	// wait for status changed
	time.Sleep(2 * time.Second)
}*/
