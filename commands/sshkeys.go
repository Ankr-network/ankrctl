/*
Copyright 2018 The Dccncli Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/godo"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// SSHKeys creates the ssh key commands heirarchy.
func SSHKeys() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "ssh-key",
			Aliases: []string{"k"},
			Short:   "sshkey commands",
			Long:    "sshkey is used to access ssh key commands",
		},
		DocCategories: []string{"sshkeys"},
		IsIndex:       true,
	}

	CmdBuilder(cmd, RunKeyList, "list", "list ssh keys", Writer,
		aliasOpt("ls"), displayerType(&displayers.Key{}), docCategories("sshkeys"))

	CmdBuilder(cmd, RunKeyGet, "get <key-id|key-fingerprint>", "get ssh key", Writer,
		aliasOpt("g"), displayerType(&displayers.Key{}), docCategories("sshkeys"))

	cmdSSHKeysCreate := CmdBuilder(cmd, RunKeyCreate, "create <key-name>", "create ssh key", Writer,
		aliasOpt("c"), displayerType(&displayers.Key{}), docCategories("sshkeys"))
	AddStringFlag(cmdSSHKeysCreate, dccncli.ArgKeyPublicKey, "", "", "Key contents", requiredOpt())

	cmdSSHKeysImport := CmdBuilder(cmd, RunKeyImport, "import <key-name>", "import ssh key", Writer,
		aliasOpt("i"), displayerType(&displayers.Key{}), docCategories("sshkeys"))
	AddStringFlag(cmdSSHKeysImport, dccncli.ArgKeyPublicKeyFile, "", "", "Public key file", requiredOpt())

	cmdRunKeyDelete := CmdBuilder(cmd, RunKeyDelete, "delete <key-id|key-fingerprint>", "delete ssh key", Writer,
		aliasOpt("d"), docCategories("sshkeys"))
	AddBoolFlag(cmdRunKeyDelete, dccncli.ArgForce, dccncli.ArgShortForce, false, "Force ssh key delete")

	cmdSSHKeysUpdate := CmdBuilder(cmd, RunKeyUpdate, "update <key-id|key-fingerprint>", "update ssh key", Writer,
		aliasOpt("u"), displayerType(&displayers.Key{}), docCategories("sshkeys"))
	AddStringFlag(cmdSSHKeysUpdate, dccncli.ArgKeyName, "", "", "Key name", requiredOpt())

	return cmd
}

// RunKeyList lists keys.
func RunKeyList(c *CmdConfig) error {
	ks := c.Keys()

	list, err := ks.List()
	if err != nil {
		return err
	}

	item := &displayers.Key{Keys: list}
	return c.Display(item)
}

// RunKeyGet retrieves a key.
func RunKeyGet(c *CmdConfig) error {
	ks := c.Keys()

	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	rawKey := c.Args[0]
	k, err := ks.Get(rawKey)

	if err != nil {
		return err
	}

	item := &displayers.Key{Keys: do.SSHKeys{*k}}
	return c.Display(item)
}

// RunKeyCreate uploads a SSH key.
func RunKeyCreate(c *CmdConfig) error {
	ks := c.Keys()

	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	name := c.Args[0]

	publicKey, err := c.Ankr.GetString(c.NS, dccncli.ArgKeyPublicKey)
	if err != nil {
		return err
	}

	kcr := &godo.KeyCreateRequest{
		Name:      name,
		PublicKey: publicKey,
	}

	r, err := ks.Create(kcr)
	if err != nil {
		return err
	}

	item := &displayers.Key{Keys: do.SSHKeys{*r}}
	return c.Display(item)
}

// RunKeyImport imports a key from a file
func RunKeyImport(c *CmdConfig) error {
	ks := c.Keys()

	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	keyPath, err := c.Ankr.GetString(c.NS, dccncli.ArgKeyPublicKeyFile)
	if err != nil {
		return err
	}

	keyName := c.Args[0]

	keyFile, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return err
	}

	_, comment, _, _, err := ssh.ParseAuthorizedKey(keyFile)
	if err != nil {
		return err
	}

	if len(keyName) < 1 {
		keyName = comment
	}

	kcr := &godo.KeyCreateRequest{
		Name:      keyName,
		PublicKey: string(keyFile),
	}

	r, err := ks.Create(kcr)
	if err != nil {
		return err
	}

	item := &displayers.Key{Keys: do.SSHKeys{*r}}
	return c.Display(item)
}

// RunKeyDelete deletes a key.
func RunKeyDelete(c *CmdConfig) error {
	ks := c.Keys()

	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	force, err := c.Ankr.GetBool(c.NS, dccncli.ArgForce)
	if err != nil {
		return nil
	}

	if force || AskForConfirm("delete ssh key") == nil {
		rawKey := c.Args[0]
		return ks.Delete(rawKey)
	} else {
		return fmt.Errorf("operation aborted")
	}
	return nil
}

// RunKeyUpdate updates a key.
func RunKeyUpdate(c *CmdConfig) error {
	ks := c.Keys()

	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	rawKey := c.Args[0]

	name, err := c.Ankr.GetString(c.NS, dccncli.ArgKeyName)
	if err != nil {
		return err
	}

	req := &godo.KeyUpdateRequest{
		Name: name,
	}

	k, err := ks.Update(rawKey, req)
	if err != nil {
		return err
	}

	item := &displayers.Key{Keys: do.SSHKeys{*k}}
	return c.Display(item)
}
