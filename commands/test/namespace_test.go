package test

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var (
	MockNamespaceName string
	MockNamespaceCpu = "1000"
	MockNamespaceMem = "512"
	MockNamespaceStorage = "8"
)

func TestRunNamespaceCreate(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// namespace create test
	t.Log("namespace create test ...")
	MockNamespaceName = "ns_create_cli_test"
	nsCreateRes, err := lc.Run( "namespace", "create", MockNamespaceName, "--cpu-limit", MockNamespaceCpu, "--mem-limit", MockNamespaceMem, "--storage-limit", MockNamespaceStorage)

	if err != nil {
		t.Error(err)
	}
	t.Log(string(nsCreateRes))
	assert.True(t, strings.Contains(string(nsCreateRes), "success"))
	t.Log("list charts successfully")

	test_ns_id := strings.Split(string(nsCreateRes), " ")[1]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// delete the namespace created
	lc.Run("namespace", "delete", test_ns_id)

	// wait for statues changed
	time.Sleep(2 * time.Second)
}