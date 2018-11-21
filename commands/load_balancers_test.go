package commands

import (
	"testing"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/godo"

	"github.com/stretchr/testify/assert"
)

var (
	testLoadBalancer = do.LoadBalancer{
		LoadBalancer: &godo.LoadBalancer{
			Algorithm: "round_robin",
			Region: &godo.Region{
				Slug: "nyc1",
			},
			StickySessions: &godo.StickySessions{},
			HealthCheck:    &godo.HealthCheck{},
		}}

	testLoadBalancerList = do.LoadBalancers{
		testLoadBalancer,
	}
)

func TestLoadBalancerCommand(t *testing.T) {
	cmd := LoadBalancer()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "get", "list", "create", "update", "delete", "add-tasks", "remove-tasks", "add-forwarding-rules", "remove-forwarding-rules")
}

func TestLoadBalancerGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		lbID := "cde2c0d6-41e3-479e-ba60-ad971227232c"
		tm.loadBalancers.On("Get", lbID).Return(&testLoadBalancer, nil)

		config.Args = append(config.Args, lbID)

		err := RunLoadBalancerGet(config)
		assert.NoError(t, err)
	})
}

func TestLoadBalancerGetNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunLoadBalancerGet(config)
		assert.Error(t, err)
	})
}

func TestLoadBalancerList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.loadBalancers.On("List").Return(testLoadBalancerList, nil)

		err := RunLoadBalancerList(config)
		assert.NoError(t, err)
	})
}

func TestLoadBalancerCreateWithInvalidTaskIDsArgs(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		config.Ankr.Set(config.NS, dccncli.ArgTaskIDs, []string{"bogus"})

		err := RunLoadBalancerCreate(config)
		assert.Error(t, err)
	})
}

func TestLoadBalancerCreateWithMalformedForwardingRulesArgs(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		config.Ankr.Set(config.NS, dccncli.ArgForwardingRules, "something,something")

		err := RunLoadBalancerCreate(config)
		assert.Error(t, err)
	})
}

func TestLoadBalancerCreate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		r := godo.LoadBalancerRequest{
			Name:       "lb-name",
			Region:     "nyc1",
			TaskIDs: []int{1, 2},
			StickySessions: &godo.StickySessions{
				Type: "none",
			},
			HealthCheck: &godo.HealthCheck{
				Protocol:               "http",
				Port:                   80,
				CheckIntervalSeconds:   4,
				ResponseTimeoutSeconds: 23,
				HealthyThreshold:       5,
				UnhealthyThreshold:     10,
			},
			ForwardingRules: []godo.ForwardingRule{
				{
					EntryProtocol:  "tcp",
					EntryPort:      3306,
					TargetProtocol: "tcp",
					TargetPort:     3306,
					TlsPassthrough: true,
				},
			},
		}
		tm.loadBalancers.On("Create", &r).Return(&testLoadBalancer, nil)

		config.Ankr.Set(config.NS, dccncli.ArgRegionSlug, "nyc1")
		config.Ankr.Set(config.NS, dccncli.ArgLoadBalancerName, "lb-name")
		config.Ankr.Set(config.NS, dccncli.ArgTaskIDs, []string{"1", "2"})
		config.Ankr.Set(config.NS, dccncli.ArgStickySessions, "type:none")
		config.Ankr.Set(config.NS, dccncli.ArgHealthCheck, "protocol:http,port:80,check_interval_seconds:4,response_timeout_seconds:23,healthy_threshold:5,unhealthy_threshold:10")
		config.Ankr.Set(config.NS, dccncli.ArgForwardingRules, "entry_protocol:tcp,entry_port:3306,target_protocol:tcp,target_port:3306,tls_passthrough:true")

		err := RunLoadBalancerCreate(config)
		assert.NoError(t, err)
	})
}

func TestLoadBalancerUpdate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		lbID := "cde2c0d6-41e3-479e-ba60-ad971227232c"
		r := godo.LoadBalancerRequest{
			Name:       "lb-name",
			Region:     "nyc1",
			TaskIDs: []int{1, 2},
			StickySessions: &godo.StickySessions{
				Type:             "cookies",
				CookieName:       "DO-LB",
				CookieTtlSeconds: 5,
			},
			HealthCheck: &godo.HealthCheck{
				Protocol:               "http",
				Port:                   80,
				CheckIntervalSeconds:   4,
				ResponseTimeoutSeconds: 23,
				HealthyThreshold:       5,
				UnhealthyThreshold:     10,
			},
			ForwardingRules: []godo.ForwardingRule{
				{
					EntryProtocol:  "http",
					EntryPort:      80,
					TargetProtocol: "http",
					TargetPort:     80,
				},
			},
		}

		tm.loadBalancers.On("Update", lbID, &r).Return(&testLoadBalancer, nil)

		config.Args = append(config.Args, lbID)
		config.Ankr.Set(config.NS, dccncli.ArgRegionSlug, "nyc1")
		config.Ankr.Set(config.NS, dccncli.ArgLoadBalancerName, "lb-name")
		config.Ankr.Set(config.NS, dccncli.ArgTaskIDs, []string{"1", "2"})
		config.Ankr.Set(config.NS, dccncli.ArgStickySessions, "type:cookies,cookie_name:DO-LB,cookie_ttl_seconds:5")
		config.Ankr.Set(config.NS, dccncli.ArgHealthCheck, "protocol:http,port:80,check_interval_seconds:4,response_timeout_seconds:23,healthy_threshold:5,unhealthy_threshold:10")
		config.Ankr.Set(config.NS, dccncli.ArgForwardingRules, "entry_protocol:http,entry_port:80,target_protocol:http,target_port:80")

		err := RunLoadBalancerUpdate(config)
		assert.NoError(t, err)
	})
}

func TestLoadBalancerUpdateNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunLoadBalancerUpdate(config)
		assert.Error(t, err)
	})
}

func TestLoadBalancerDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		lbID := "cde2c0d6-41e3-479e-ba60-ad971227232c"
		tm.loadBalancers.On("Delete", lbID).Return(nil)

		config.Args = append(config.Args, lbID)
		config.Ankr.Set(config.NS, dccncli.ArgForce, true)

		err := RunLoadBalancerDelete(config)
		assert.NoError(t, err)
	})
}

func TestLoadBalancerDeleteNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunLoadBalancerDelete(config)
		assert.Error(t, err)
	})
}

func TestLoadBalancerAddTasks(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		lbID := "cde2c0d6-41e3-479e-ba60-ad971227232c"
		tm.loadBalancers.On("AddTasks", lbID, 1, 23).Return(nil)

		config.Args = append(config.Args, lbID)
		config.Ankr.Set(config.NS, dccncli.ArgTaskIDs, []string{"1", "23"})

		err := RunLoadBalancerAddTasks(config)
		assert.NoError(t, err)
	})
}

func TestLoadBalancerAddTasksNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunLoadBalancerAddTasks(config)
		assert.Error(t, err)
	})
}

func TestLoadBalancerRemoveTasks(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		lbID := "cde2c0d6-41e3-479e-ba60-ad971227232c"
		tm.loadBalancers.On("RemoveTasks", lbID, 321).Return(nil)

		config.Args = append(config.Args, lbID)
		config.Ankr.Set(config.NS, dccncli.ArgTaskIDs, []string{"321"})

		err := RunLoadBalancerRemoveTasks(config)
		assert.NoError(t, err)
	})
}

func TestLoadBalancerRemoveTasksNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunLoadBalancerRemoveTasks(config)
		assert.Error(t, err)
	})
}

func TestLoadBalancerAddForwardingRules(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		lbID := "cde2c0d6-41e3-479e-ba60-ad971227232c"
		forwardingRule := godo.ForwardingRule{
			EntryProtocol:  "http",
			EntryPort:      80,
			TargetProtocol: "http",
			TargetPort:     80,
		}
		tm.loadBalancers.On("AddForwardingRules", lbID, forwardingRule).Return(nil)

		config.Args = append(config.Args, lbID)
		config.Ankr.Set(config.NS, dccncli.ArgForwardingRules, "entry_protocol:http,entry_port:80,target_protocol:http,target_port:80")

		err := RunLoadBalancerAddForwardingRules(config)
		assert.NoError(t, err)
	})
}

func TestLoadBalancerAddForwardingRulesNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunLoadBalancerAddForwardingRules(config)
		assert.Error(t, err)
	})
}

func TestLoadBalancerRemoveForwardingRules(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		lbID := "cde2c0d6-41e3-479e-ba60-ad971227232c"
		forwardingRules := []godo.ForwardingRule{
			{
				EntryProtocol:  "http",
				EntryPort:      80,
				TargetProtocol: "http",
				TargetPort:     80,
			},
			{
				EntryProtocol:  "tcp",
				EntryPort:      3306,
				TargetProtocol: "tcp",
				TargetPort:     3306,
				TlsPassthrough: true,
			},
		}
		tm.loadBalancers.On("RemoveForwardingRules", lbID, forwardingRules[0], forwardingRules[1] ).Return(nil)

		config.Args = append(config.Args, lbID)
		config.Ankr.Set(config.NS, dccncli.ArgForwardingRules, "entry_protocol:http,entry_port:80,target_protocol:http,target_port:80 entry_protocol:tcp,entry_port:3306,target_protocol:tcp,target_port:3306,tls_passthrough:true")

		err := RunLoadBalancerRemoveForwardingRules(config)
		assert.NoError(t, err)
	})
}

func TestLoadBalancerRemoveForwardingRulesNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunLoadBalancerRemoveForwardingRules(config)
		assert.Error(t, err)
	})
}
