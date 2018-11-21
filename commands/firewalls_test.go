package commands

import (
	"strconv"
	"testing"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/godo"

	"github.com/stretchr/testify/assert"
)

var (
	testFirewall = do.Firewall{
		Firewall: &godo.Firewall{
			Name: "my firewall",
		},
	}

	testFirewallList = do.Firewalls{
		testFirewall,
	}
)

func TestFirewallCommand(t *testing.T) {
	cmd := Firewall()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "get", "create", "update", "list", "list-by-task", "delete", "add-tasks", "remove-tasks", "add-tags", "remove-tags", "add-rules", "remove-rules")
}

func TestFirewallGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tm.firewalls.On("Get", fID).Return(&testFirewall, nil)

		config.Args = append(config.Args, fID)

		err := RunFirewallGet(config)
		assert.NoError(t, err)
	})
}

func TestFirewallCreate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		firewallCreateRequest := &godo.FirewallRequest{
			Name: "firewall",
			InboundRules: []godo.InboundRule{
				{
					Protocol:  "icmp",
					PortRange: "",
					Sources:   &godo.Sources{},
				},
				{
					Protocol:  "tcp",
					PortRange: "8000-9000",
					Sources: &godo.Sources{
						Addresses: []string{"127.0.0.0", "0::/0", "::/1"},
					},
				},
			},
			Tags:       []string{"backend"},
			TaskIDs: []int{1, 2},
		}
		tm.firewalls.On("Create", firewallCreateRequest).Return(&testFirewall, nil)

		config.Ankr.Set(config.NS, dccncli.ArgFirewallName, "firewall")
		config.Ankr.Set(config.NS, dccncli.ArgTagNames, []string{"backend"})
		config.Ankr.Set(config.NS, dccncli.ArgTaskIDs, []string{"1", "2"})
		config.Ankr.Set(config.NS, dccncli.ArgInboundRules, "protocol:icmp protocol:tcp,ports:8000-9000,address:127.0.0.0,address:0::/0,address:::/1")

		err := RunFirewallCreate(config)
		assert.NoError(t, err)
	})
}

func TestFirewallUpdate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		firewallUpdateRequest := &godo.FirewallRequest{
			Name: "firewall",
			InboundRules: []godo.InboundRule{
				{
					Protocol:  "tcp",
					PortRange: "8000-9000",
					Sources: &godo.Sources{
						Addresses: []string{"127.0.0.0"},
					},
				},
			},
			OutboundRules: []godo.OutboundRule{
				{
					Protocol:  "tcp",
					PortRange: "8080",
					Destinations: &godo.Destinations{
						LoadBalancerUIDs: []string{"lb-uuid"},
						Tags:             []string{"new-tasks"},
					},
				},
				{
					Protocol:  "tcp",
					PortRange: "80",
					Destinations: &godo.Destinations{
						Addresses: []string{"192.168.0.0"},
					},
				},
			},
			TaskIDs: []int{1},
		}
		tm.firewalls.On("Update", fID, firewallUpdateRequest).Return(&testFirewall, nil)

		config.Args = append(config.Args, fID)
		config.Ankr.Set(config.NS, dccncli.ArgFirewallName, "firewall")
		config.Ankr.Set(config.NS, dccncli.ArgTaskIDs, []string{"1"})
		config.Ankr.Set(config.NS, dccncli.ArgInboundRules, "protocol:tcp,ports:8000-9000,address:127.0.0.0")
		config.Ankr.Set(config.NS, dccncli.ArgOutboundRules, "protocol:tcp,ports:8080,load_balancer_uid:lb-uuid,tag:new-tasks protocol:tcp,ports:80,address:192.168.0.0")

		err := RunFirewallUpdate(config)
		assert.NoError(t, err)
	})
}

func TestFirewallList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.firewalls.On("List").Return(testFirewallList, nil)

		err := RunFirewallList(config)
		assert.NoError(t, err)
	})
}

func TestFirewallListByTask(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		dID := 124
		tm.firewalls.On("ListByTask", dID).Return(testFirewallList, nil)
		config.Args = append(config.Args, strconv.Itoa(dID))

		err := RunFirewallListByTask(config)
		assert.NoError(t, err)
	})
}

func TestFirewallDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tm.firewalls.On("Delete", fID).Return(nil)

		config.Args = append(config.Args, fID)
		config.Ankr.Set(config.NS, dccncli.ArgForce, true)

		err := RunFirewallDelete(config)
		assert.NoError(t, err)
	})
}

func TestFirewallAddTasks(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		taskIDs := []int{1, 2}
		tm.firewalls.On("AddTasks", fID, taskIDs[0], taskIDs[1]).Return(nil)

		config.Args = append(config.Args, fID)
		config.Ankr.Set(config.NS, dccncli.ArgTaskIDs, []string{"1", "2"})

		err := RunFirewallAddTasks(config)
		assert.NoError(t, err)
	})
}

func TestFirewallRemoveTasks(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "cde2c0d6-41e3-479e-ba60-ad971227232c"
		taskIDs := []int{1}
		tm.firewalls.On("RemoveTasks", fID, taskIDs[0]).Return(nil)

		config.Args = append(config.Args, fID)
		config.Ankr.Set(config.NS, dccncli.ArgTaskIDs, []string{"1"})

		err := RunFirewallRemoveTasks(config)
		assert.NoError(t, err)
	})
}

func TestFirewallAddTags(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tags := []string{"frontend", "backend"}
		tm.firewalls.On("AddTags", fID, tags[0], tags[1]).Return(nil)

		config.Args = append(config.Args, fID)
		config.Ankr.Set(config.NS, dccncli.ArgTagNames, []string{"frontend", "backend"})

		err := RunFirewallAddTags(config)
		assert.NoError(t, err)
	})
}

func TestFirewallRemoveTags(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tags := []string{"backend"}
		tm.firewalls.On("RemoveTags", fID, tags[0]).Return(nil)

		config.Args = append(config.Args, fID)
		config.Ankr.Set(config.NS, dccncli.ArgTagNames, []string{"backend"})

		err := RunFirewallRemoveTags(config)
		assert.NoError(t, err)
	})
}

func TestFirewallAddRules(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		inboundRules := []godo.InboundRule{
			{
				Protocol:  "tcp",
				PortRange: "80",
				Sources: &godo.Sources{
					Addresses: []string{"127.0.0.0", "0.0.0.0/0", "2604:A880:0002:00D0:0000:0000:32F1:E001"},
				},
			},
			{
				Protocol:  "tcp",
				PortRange: "8080",
				Sources: &godo.Sources{
					Tags:       []string{"backend"},
					TaskIDs: []int{1, 2, 3},
				},
			},
		}
		outboundRules := []godo.OutboundRule{
			{
				Protocol:  "tcp",
				PortRange: "22",
				Destinations: &godo.Destinations{
					LoadBalancerUIDs: []string{"lb-uuid"},
				},
			},
		}
		firewallRulesRequest := &godo.FirewallRulesRequest{
			InboundRules:  inboundRules,
			OutboundRules: outboundRules,
		}

		tm.firewalls.On("AddRules", fID, firewallRulesRequest).Return(nil)

		config.Args = append(config.Args, fID)
		config.Ankr.Set(config.NS, dccncli.ArgInboundRules, "protocol:tcp,ports:80,address:127.0.0.0,address:0.0.0.0/0,address:2604:A880:0002:00D0:0000:0000:32F1:E001 protocol:tcp,ports:8080,tag:backend,task_id:1,task_id:2,task_id:3")
		config.Ankr.Set(config.NS, dccncli.ArgOutboundRules, "protocol:tcp,ports:22,load_balancer_uid:lb-uuid")

		err := RunFirewallAddRules(config)
		assert.NoError(t, err)
	})
}

func TestFirewallRemoveRules(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		inboundRules := []godo.InboundRule{
			{
				Protocol:  "tcp",
				PortRange: "80",
				Sources: &godo.Sources{
					Addresses: []string{"0.0.0.0/0"},
				},
			},
		}
		outboundRules := []godo.OutboundRule{
			{
				Protocol:  "tcp",
				PortRange: "22",
				Destinations: &godo.Destinations{
					Tags:      []string{"back:end"},
					Addresses: []string{"::/0"},
				},
			},
		}
		firewallRulesRequest := &godo.FirewallRulesRequest{
			InboundRules:  inboundRules,
			OutboundRules: outboundRules,
		}

		tm.firewalls.On("RemoveRules", fID, firewallRulesRequest).Return(nil)

		config.Args = append(config.Args, fID)
		config.Ankr.Set(config.NS, dccncli.ArgInboundRules, "protocol:tcp,ports:80,address:0.0.0.0/0")
		config.Ankr.Set(config.NS, dccncli.ArgOutboundRules, "protocol:tcp,ports:22,tag:back:end,address:::/0")

		err := RunFirewallRemoveRules(config)
		assert.NoError(t, err)
	})
}
