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
	test_ns_id := strings.Split(string(nsCreateRes), " ")[1]
	if err != nil {
		t.Error(err)
	}else{
	t.Log(string(nsCreateRes))
	assert.True(t, strings.Contains(string(nsCreateRes), "success"))
	t.Log("create namespace successfully")
	}

	// wait for status changed
	time.Sleep(10 * time.Second)

	// delete the namespace created
	lc.Run("namespace", "delete", test_ns_id, "-f")

	// wait for statues changed
	time.Sleep(2 * time.Second)
}

func TestRunNamespaceUpdate(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// namespace create for update test
	MockNamespaceName = "ns_update_cli_test"
	nsCreateRes, _ := lc.Run( "namespace", "create", MockNamespaceName, "--cpu-limit", MockNamespaceCpu, "--mem-limit", MockNamespaceMem, "--storage-limit", MockNamespaceStorage)
	test_ns_id := strings.Split(string(nsCreateRes), " ")[1]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// namespace update test
	t.Log("namespace update test ...")
	nsUpdateRes, err := lc.Run("namespace", "update", test_ns_id, "--cpu-limit", "1024", "--mem-limit", "2048", "--storage-limit", "16")
	if err != nil {
		t.Error(err)
	}else{
	t.Log(string(nsUpdateRes))
	assert.True(t, strings.Contains(string(nsUpdateRes), "success"))
	}

	// wait for statues changed
	time.Sleep(5 * time.Second)

	// delete the namespace created
	lc.Run("namespace", "delete", test_ns_id, "-f")

	// wait for statues changed
	time.Sleep(2 * time.Second)
}

func TestRunNamespaceDelete(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// namespace create for delete test
	MockNamespaceName = "ns_delete_cli_test"
	nsCreateRes, _ := lc.Run( "namespace", "create", MockNamespaceName, "--cpu-limit", MockNamespaceCpu, "--mem-limit", MockNamespaceMem, "--storage-limit", MockNamespaceStorage)
	test_ns_id := strings.Split(string(nsCreateRes), " ")[1]

	// wait for status changed
	time.Sleep(10 * time.Second)

	// namespace delete test
	t.Log("namespace delete test ...")
	nsDeleteRes, err := lc.Run("namespace", "delete", test_ns_id, "-f")
	if err != nil {
		t.Error(err)
	}else{
	t.Log(string(nsDeleteRes))
	assert.True(t, strings.Contains(string(nsDeleteRes), "success"))
	t.Log("delete namespace successfully")
	}

	// wait for statues changed
	time.Sleep(2 * time.Second)
}

func TestRunNamespaceList(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// namespace list test
	t.Log("namespace list test ...")
	nsListRes, err := lc.Run("namespace", "list")
	if err != nil {
		t.Error(err)
	}else{
	t.Log(string(nsListRes))
	assert.True(t, strings.Contains(string(nsListRes), "Name"))
	assert.True(t, strings.Contains(string(nsListRes), "ID"))
	}

	// wait for statues changed
	time.Sleep(2 * time.Second)
}

