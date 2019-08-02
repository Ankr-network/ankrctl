package test

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var (
	MockAppName = "app_create_cli_test"
	ChartName = "wordpress"
	ChartRepo = "stable"
	ChartVersion = "5.6.2"
)


func TestRunAppCreate(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// app create test
	// case 1: use a prepared namespace to create the app
	t.Log("app create test ... (case 1)")

	// create a namespace for app_create test
	nsCreateRes, _ := lc.Run( "namespace", "create", "app_create_cli_test", "--cpu-limit", MockNamespaceCpu, "--mem-limit", MockNamespaceMem, "--storage-limit", MockNamespaceStorage)
	test_ns_id := strings.Split(string(nsCreateRes), " ")[1]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// create app
	appCreateRes, err := lc.Run("app", "create", MockAppName, "--chart-name", ChartName, "--chart-repo", ChartRepo, "--chart-version", ChartVersion,  "--ns-id", test_ns_id)
	app_id_pre := strings.Split(string(appCreateRes), " ")[5]
	app_id := strings.Split(app_id_pre, ",")[0]
	if err != nil {
		t.Error(err.Error())
	}else{
		t.Log(string(appCreateRes))
		assert.True(t, len(string(appCreateRes)) > 0)
		assert.True(t, strings.Contains(string(appCreateRes), "success"))
		assert.True(t, len(app_id) > 0)
	}


	// wait for statues changed
	time.Sleep(10 * time.Second)

	// purge the app created
	lc.Run("app", "purge", app_id, "-f")

	// wait for statues changed
	time.Sleep(10 * time.Second)

	// cancel the namespace created
	lc.Run("namespace", "delete", test_ns_id, "-f")

	// wait for statues changed
	time.Sleep(2 * time.Second)


	// case 2: create app with namespace at the same time
	t.Log("app create test ... (case 2)")
	appCreateRes_1, err_1 := lc.Run("app", "create", MockAppName, "--chart-name", ChartName, "--chart-repo", ChartRepo, "--chart-version", ChartVersion,  "--ns-name", "app_create_cli_test", "--cpu-limit", "1000", "--mem-limit", "2048", "--storage-limit","8")
	app_id_pre_1 := strings.Split(string(appCreateRes), " ")[5]
	app_id_1 := strings.Split(app_id_pre_1, ",")[0]
	if err_1 != nil {
		t.Error(err_1.Error())
	}else{
		t.Log(string(appCreateRes_1))
		assert.True(t, len(string(appCreateRes_1)) > 0)
		assert.True(t, strings.Contains(string(appCreateRes_1), "success"))
		assert.True(t, len(app_id_1) > 0)
	}


	// wait for statues changed
	time.Sleep(10 * time.Second)

	// purge the app created
	lc.Run("app", "purge", app_id_1, "-f")

	// wait for statues changed
	time.Sleep(2 * time.Second)

}

func TestRunAppUpdate(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// app create for app_update test
	// create a namespace for app_create
	nsCreateRes, _ := lc.Run( "namespace", "create", "app_update_cli_test", "--cpu-limit", MockNamespaceCpu, "--mem-limit", MockNamespaceMem, "--storage-limit", MockNamespaceStorage)
	test_ns_id := strings.Split(string(nsCreateRes), " ")[1]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// download and upload a chart for app_create
	lc.Run( "chart", "download", "wordpress", "--download-repo", "stable", "--download-version", "5.6.2")
	lc.Run( "chart", "upload", "app_update_test", "--upload-file", "/go/src/github.com/Ankr-network/dccn-cli/commands/test/wordpress-5.6.2.tgz", "--upload-version", chartUploadVersion)

	// wait for status changed
	time.Sleep(10 * time.Second)

	// create app
	appCreateRes, _ := lc.Run("app", "create", "app_update_cli_test", "--chart-name", "app_update_test", "--chart-repo", "user", "--chart-version", chartUploadVersion,  "--ns-id", test_ns_id)
	app_id_pre := strings.Split(string(appCreateRes), " ")[5]
	app_id := strings.Split(app_id_pre, ",")[0]
	t.Log(app_id)

	// wait for status changed
	time.Sleep(15 * time.Second)

	// check
	lc.Run("app", "list")

	// update app test
	t.Log("app update test ... ")
	appUpdateRes, err := lc.Run("app", "update", app_id, "--app-name", "app_update_result", "--update-version", "6.6.6")
	if err != nil {
		t.Error(err.Error())
	}else{
		t.Log(string(appUpdateRes))
		assert.True(t, strings.Contains(string(appUpdateRes), "success"))
	}

	// wait for statues changed
	time.Sleep(10 * time.Second)

	// purge the app created
	lc.Run("app", "purge", app_id, "-f")

	// wait for statues changed
	time.Sleep(10 * time.Second)

	// cancel the namespace created
	lc.Run("namespace", "delete", test_ns_id, "-f")

	// delete the chart upload
	lc.Run("chart", "delete", chartUploadName, "--delete-version", "6.6.6", "-f")

	// wait for statues changed
	time.Sleep(2 * time.Second)

}

func TestRunAppCancel(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// app create for app_cancel test
	// create a namespace for app_create
	nsCreateRes, _ := lc.Run( "namespace", "create", "app_cancel_cli_test", "--cpu-limit", MockNamespaceCpu, "--mem-limit", MockNamespaceMem, "--storage-limit", MockNamespaceStorage)
	test_ns_id := strings.Split(string(nsCreateRes), " ")[1]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// create app
	appCreateRes, _ := lc.Run("app", "create", "app_cancel_cli_test", "--chart-name", ChartName, "--chart-repo", ChartRepo, "--chart-version", ChartVersion,  "--ns-id", test_ns_id)
	app_id_pre := strings.Split(string(appCreateRes), " ")[5]
	app_id := strings.Split(app_id_pre, ",")[0]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// cancel app test
	t.Log("app cancel test ... ")
	appCancelRes, err := lc.Run("app", "cancel", app_id, "-f")
	if err != nil {
		t.Error(err.Error())
	}else{
		t.Log(string(appCancelRes))
		assert.True(t, strings.Contains(string(appCancelRes), "success"))
	}

	// wait for statues changed
	time.Sleep(10 * time.Second)

	// cancel the namespace created
	lc.Run("namespace", "delete", test_ns_id, "-f")

	// wait for statues changed
	time.Sleep(2 * time.Second)

}

func TestRunAppPurge(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// app create for app_purge test
	// create a namespace for app_create
	nsCreateRes, _ := lc.Run( "namespace", "create", "app_cancel_cli_test", "--cpu-limit", MockNamespaceCpu, "--mem-limit", MockNamespaceMem, "--storage-limit", MockNamespaceStorage)
	test_ns_id := strings.Split(string(nsCreateRes), " ")[1]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// create app
	appCreateRes, _ := lc.Run("app", "create", "app_purge_cli_test", "--chart-name", ChartName, "--chart-repo", ChartRepo, "--chart-version", ChartVersion,  "--ns-id", test_ns_id)
	app_id_pre := strings.Split(string(appCreateRes), " ")[5]
	app_id := strings.Split(app_id_pre, ",")[0]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// purge app test
	t.Log("app purge test ... ")
	appPurgeRes, err := lc.Run("app", "purge", app_id, "-f")
	if err != nil {
		t.Error(err.Error())
	}else{
		t.Log(string(appPurgeRes))
		assert.True(t, strings.Contains(string(appPurgeRes), "success"))
	}

	// wait for statues changed
	time.Sleep(10 * time.Second)

	// cancel the namespace created
	lc.Run("namespace", "delete", test_ns_id, "-f")

	// wait for statues changed
	time.Sleep(2 * time.Second)
}

func TestRunAppList(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// app list test
	t.Log("app list test ... ")
	appListRes, err := lc.Run("app", "list")
	if err != nil {
		t.Error(err.Error())
	}else{
		t.Log(string(appListRes))
		assert.True(t, strings.Contains(string(appListRes), "ID"))
		assert.True(t, strings.Contains(string(appListRes), "Name"))
	}

	// wait for statues changed
	time.Sleep(2 * time.Second)
}

func TestRunAppDetail(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// app create for app_detail test
	// create a namespace for app_detail test
	nsCreateRes, _ := lc.Run( "namespace", "create", "app_detail_cli_test", "--cpu-limit", MockNamespaceCpu, "--mem-limit", MockNamespaceMem, "--storage-limit", MockNamespaceStorage)
	test_ns_id := strings.Split(string(nsCreateRes), " ")[1]

	t.Log(test_ns_id)

	// wait for status changed
	time.Sleep(10 * time.Second)

	// create app
	appCreateRes, _ := lc.Run("app", "create", MockAppName, "--chart-name", ChartName, "--chart-repo", ChartRepo, "--chart-version", ChartVersion,  "--ns-id", test_ns_id)
	app_id_pre := strings.Split(string(appCreateRes), " ")[5]
	app_id := strings.Split(app_id_pre, ",")[0]
	t.Log(app_id)

	// wait for status changed
	time.Sleep(15 * time.Second)

	// app detail test
	t.Log("app detail test ... ")
	appDetailRes, err := lc.Run("app", "detail", app_id)
	if err != nil {
		t.Error(err.Error())
	}else{
		t.Log(string(appDetailRes))
		assert.True(t, strings.Contains(string(appDetailRes), "detail"))
	}

	// wait for statues changed
	time.Sleep(10 * time.Second)

	// purge the app created
	lc.Run("app", "purge", app_id, "-f")

	// wait for statues changed
	time.Sleep(10 * time.Second)

	// cancel the namespace created
	lc.Run("namespace", "delete", test_ns_id, "-f")

	// wait for statues changed
	time.Sleep(2 * time.Second)
}

func TestRunAppOverview(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// app create for app_overview test
	// create a namespace for app_overview test
	nsCreateRes, _ := lc.Run( "namespace", "create", "app_overview_cli_test", "--cpu-limit", MockNamespaceCpu, "--mem-limit", MockNamespaceMem, "--storage-limit", MockNamespaceStorage)
	test_ns_id := strings.Split(string(nsCreateRes), " ")[1]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// create app
	appCreateRes, _ := lc.Run("app", "create", MockAppName, "--chart-name", ChartName, "--chart-repo", ChartRepo, "--chart-version", ChartVersion,  "--ns-id", test_ns_id)
	app_id_pre := strings.Split(string(appCreateRes), " ")[5]
	app_id := strings.Split(app_id_pre, ",")[0]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// app overview test
	t.Log("app overview test ... ")
	appOverviewRes, err := lc.Run("app", "overview", app_id)
	if err != nil {
		t.Error(err.Error())
	}else{
		t.Log(string(appOverviewRes))
		assert.True(t, strings.Contains(string(appOverviewRes), "Namespace"))
		assert.True(t, strings.Contains(string(appOverviewRes), "Cluster"))
	}

	// wait for statues changed
	time.Sleep(10 * time.Second)

	// purge the app created
	lc.Run("app", "purge", app_id, "-f")

	// wait for statues changed
	time.Sleep(10 * time.Second)

	// cancel the namespace created
	lc.Run("namespace", "delete", test_ns_id, "-f")

	// wait for statues changed
	time.Sleep(2 * time.Second)
}
