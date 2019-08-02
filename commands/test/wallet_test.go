package test

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestRunWalletKeylist(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// wallet list test
	walletListRes, err := lc.Run( "wallet", "listkey")
	if err != nil {
		t.Error(err)
	}else{
		t.Log(string(walletListRes))
		assert.True(t, strings.Contains(string(walletListRes), "Name"))
		assert.True(t, strings.Contains(string(walletListRes), "Address"))
		assert.True(t, strings.Contains(string(walletListRes), "Public Key"))
		t.Log("list wallet successfully")
	}

	// wait for status changed
	time.Sleep(2 * time.Second)
}

func TestRunWalletGenkey(t *testing.T) {

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// wallet genkey test
	walletGenKeyRes, err := lc.Run( "wallet", "genkey", "my_new_key", "-f")
	if err != nil {
		t.Error(err)
	}else{
		t.Log(string(walletGenKeyRes))
		assert.True(t, strings.Contains(string(walletGenKeyRes), "private key"))
		assert.True(t, strings.Contains(string(walletGenKeyRes), "public key"))
		assert.True(t, strings.Contains(string(walletGenKeyRes), "address"))
		t.Log("wallet genkey successfully")
	}

	// wait for status changed
	time.Sleep(2 * time.Second)
}

func TestRunWalletImportkey(t *testing.T) {

	// need a break, sleep 60s
	time.Sleep(60*time.Second)

	// user login at first
	_, err := lc.Run( "user", "login", "--email", CorrectUserEmail, "--password", CorrectPassword)
	if err != nil {
		t.Error(err)
	}

	// wallet importkey test
	walletImportKeyRes, err := lc.Run( "wallet", "importkey", "my_import_key", "--keystore", "UTC--2019-07-06T22-29-42.344574000Z--219B0A5F896B7A1949128B8F5136362BF939994D769D58", "-f")
	if err != nil {
		t.Error(err)
	}else{
		t.Log(string(walletImportKeyRes))
		assert.True(t, strings.Contains(string(walletImportKeyRes), "private key"))
		assert.True(t, strings.Contains(string(walletImportKeyRes), "public key"))
		assert.True(t, strings.Contains(string(walletImportKeyRes), "address"))
		t.Log("wallet genkey successfully")
	}

	// wait for status changed
	time.Sleep(2 * time.Second)
}