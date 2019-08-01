package test

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var (
	MockAppName = "app_create_cli_test"
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
	appCreateRes, err := lc.Run("app", "create", MockAppName, "--chart-name", "wordpress", "--chart-repo", "stable", "--chart-version", "5.6.0",  "--ns-id", test_ns_id)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(string(appCreateRes))
	assert.True(t, len(string(appCreateRes)) > 0)
	assert.True(t, strings.Contains(string(appCreateRes), "success"))
	app_id := strings.Fields(string(appCreateRes))[1]
	assert.True(t, len(app_id) > 0)

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
	appCreateRes_1, err_1 := lc.Run("app", "create", MockAppName, "--chart-name", "wordpress", "--chart-repo", "stable", "--chart-version", "5.6.0",  "--ns-name", "app_create_cli_test", "--cpu-limit", "1000", "--mem-limit", "2048", "--storage-limit","8")

	if err_1 != nil {
		t.Error(err_1.Error())
	}
	t.Log(string(appCreateRes_1))
	assert.True(t, len(string(appCreateRes_1)) > 0)
	assert.True(t, strings.Contains(string(appCreateRes_1), "success"))
	app_id_1 := strings.Fields(string(appCreateRes_1))[1]
	assert.True(t, len(app_id_1) > 0)

	// wait for statues changed
	time.Sleep(10 * time.Second)

	// purge the app created
	lc.Run("app", "purge", app_id_1, "-f")

	// wait for statues changed
	time.Sleep(2 * time.Second)

}


