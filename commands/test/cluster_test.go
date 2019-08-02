package test

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestRunClusterList(t *testing.T) {
	// as the first test, sleep 60s
	time.Sleep(60 * time.Second)

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// cluster list test
	t.Log("cluster list test ...")
	clusterListRes, err := lc.Run( "cluster", "list")
	if err != nil {
		t.Error(err)
	}else{
		t.Log(string(clusterListRes))
		assert.True(t, strings.Contains(string(clusterListRes), "Name"))
		assert.True(t, strings.Contains(string(clusterListRes), "ID"))
		t.Log("list clusters successfully")
	}

	// wait for status changed
	time.Sleep(2 * time.Second)
}

func TestRunNetworkInfo(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// network info test
	t.Log("network info test ...")
	clusterNetworkRes, err := lc.Run( "cluster", "network")
	if err != nil {
		t.Error(err)
	}else{
		t.Log(string(clusterNetworkRes))
		assert.True(t, strings.Contains(string(clusterNetworkRes), "User"))
		assert.True(t, strings.Contains(string(clusterNetworkRes), "Host"))
		assert.True(t, strings.Contains(string(clusterNetworkRes), "Namespace"))
		assert.True(t, strings.Contains(string(clusterNetworkRes), "Container"))
		assert.True(t, strings.Contains(string(clusterNetworkRes), "Traffic"))
		t.Log("get network infos successfully")
	}

	// wait for status changed
	time.Sleep(2 * time.Second)
}

