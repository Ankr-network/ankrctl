/*
Copyright 2018 The Doctl Authors All rights reserved.
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
	"testing"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/godo"
	"github.com/stretchr/testify/assert"
)

func TestFloatingIPCommands(t *testing.T) {
	cmd := FloatingIP()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "create", "delete", "get", "list")
}

func TestFloatingIPsList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.floatingIPs.On("List").Return(testFloatingIPList, nil)

		RunFloatingIPList(config)
	})
}

func TestFloatingIPsGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.floatingIPs.On("Get", "127.0.0.1").Return(&testFloatingIP, nil)

		config.Args = append(config.Args, "127.0.0.1")

		RunFloatingIPGet(config)
	})
}

func TestFloatingIPsCreate_Droplet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		ficr := &godo.FloatingIPCreateRequest{DropletID: 1}
		tm.floatingIPs.On("Create", ficr).Return(&testFloatingIP, nil)

		config.Doit.Set(config.NS, doctl.ArgDropletID, 1)

		err := RunFloatingIPCreate(config)
		assert.NoError(t, err)
	})
}

func TestFloatingIPsCreate_Region(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		ficr := &godo.FloatingIPCreateRequest{Region: "dev0"}
		tm.floatingIPs.On("Create", ficr).Return(&testFloatingIP, nil)

		config.Doit.Set(config.NS, doctl.ArgRegionSlug, "dev0")

		err := RunFloatingIPCreate(config)
		assert.NoError(t, err)
	})
}

func TestFloatingIPsCreate_fail_with_no_args(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunFloatingIPCreate(config)
		assert.Error(t, err)
	})
}

func TestFloatingIPsCreate_fail_with_both_args(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		config.Doit.Set(config.NS, doctl.ArgDropletID, 1)
		config.Doit.Set(config.NS, doctl.ArgRegionSlug, "dev0")

		err := RunFloatingIPCreate(config)
		assert.Error(t, err)
	})
}

func TestFloatingIPsDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.floatingIPs.On("Delete", "127.0.0.1").Return(nil)

		config.Args = append(config.Args, "127.0.0.1")

		config.Doit.Set(config.NS, doctl.ArgForce, true)

		RunFloatingIPDelete(config)
	})
}
