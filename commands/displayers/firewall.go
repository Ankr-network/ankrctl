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

package displayers

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Ankr-network/dccn-cli/do"
)

type Firewall struct {
	Firewalls do.Firewalls
}

var _ Displayable = &Firewall{}

func (f *Firewall) JSON(out io.Writer) error {
	return writeJSON(f.Firewalls, out)
}

func (f *Firewall) Cols() []string {
	return []string{
		"ID",
		"Name",
		"Status",
		"Created",
		"InboundRules",
		"OutboundRules",
		"TaskIDs",
		"Tags",
		"PendingChanges",
	}
}

func (f *Firewall) ColMap() map[string]string {
	return map[string]string{
		"ID":             "ID",
		"Name":           "Name",
		"Status":         "Status",
		"Created":        "Created At",
		"InboundRules":   "Inbound Rules",
		"OutboundRules":  "Outbound Rules",
		"TaskIDs":     "Task IDs",
		"Tags":           "Tags",
		"PendingChanges": "Pending Changes",
	}
}

func (f *Firewall) KV() []map[string]interface{} {
	out := []map[string]interface{}{}

	for _, fw := range f.Firewalls {
		irs, ors := firewallRulesPrintHelper(fw)
		o := map[string]interface{}{
			"ID":             fw.ID,
			"Name":           fw.Name,
			"Status":         fw.Status,
			"Created":        fw.Created,
			"InboundRules":   irs,
			"OutboundRules":  ors,
			"TaskIDs":     taskListHelper(fw.TaskIDs),
			"Tags":           strings.Join(fw.Tags, ","),
			"PendingChanges": firewallPendingChangesPrintHelper(fw),
		}
		out = append(out, o)
	}

	return out
}

func firewallRulesPrintHelper(fw do.Firewall) (string, string) {
	var irs, ors []string

	for _, ir := range fw.InboundRules {
		ss := firewallInAndOutboundRulesPrintHelper(ir.Sources.Addresses, ir.Sources.Tags, ir.Sources.TaskIDs, ir.Sources.LoadBalancerUIDs)
		if ir.Protocol == "icmp" {
			irs = append(irs, fmt.Sprintf("%v:%v,%v", "protocol", ir.Protocol, ss))
		} else {
			irs = append(irs, fmt.Sprintf("%v:%v,%v:%v,%v", "protocol", ir.Protocol, "ports", ir.PortRange, ss))
		}
	}

	for _, or := range fw.OutboundRules {
		ds := firewallInAndOutboundRulesPrintHelper(or.Destinations.Addresses, or.Destinations.Tags, or.Destinations.TaskIDs, or.Destinations.LoadBalancerUIDs)
		if or.Protocol == "icmp" {
			ors = append(ors, fmt.Sprintf("%v:%v,%v", "protocol", or.Protocol, ds))
		} else {
			ors = append(ors, fmt.Sprintf("%v:%v,%v:%v,%v", "protocol", or.Protocol, "ports", or.PortRange, ds))
		}
	}

	return strings.Join(irs, " "), strings.Join(ors, " ")
}

func firewallInAndOutboundRulesPrintHelper(addresses []string, tags []string, taskIDs []int, loadBalancerUIDs []string) string {
	output := []string{}
	resources := map[string][]string{
		"address":           addresses,
		"tag":               tags,
		"load_balancer_uid": loadBalancerUIDs,
	}

	for k, vs := range resources {
		for _, r := range vs {
			output = append(output, fmt.Sprintf("%v:%v", k, r))
		}
	}

	for _, dID := range taskIDs {
		output = append(output, fmt.Sprintf("%v:%v", "task_id", dID))
	}

	return strings.Join(output, ",")
}

func firewallPendingChangesPrintHelper(fw do.Firewall) string {
	output := []string{}

	for _, pc := range fw.PendingChanges {
		output = append(output, fmt.Sprintf("%v:%v,%v:%v,%v:%v", "task_id", pc.TaskID, "removing", pc.Removing, "status", pc.Status))
	}

	return strings.Join(output, " ")
}

func taskListHelper(IDs []int) string {
	output := []string{}

	for _, id := range IDs {
		output = append(output, strconv.Itoa(id))
	}

	return strings.Join(output, ",")
}
