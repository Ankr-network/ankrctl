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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/godo"

	"github.com/spf13/cobra"
)

// Firewall creates the firewall command.
func Firewall() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "firewall",
			Short: "firewall commands",
			Long:  "firewall is used to access firewall commands",
		},
	}

	CmdBuilder(cmd, RunFirewallGet, "get <id>", "get firewall", Writer, aliasOpt("g"), displayerType(&displayers.Firewall{}))

	cmdFirewallCreate := CmdBuilder(cmd, RunFirewallCreate, "create", "create firewall", Writer, aliasOpt("c"), displayerType(&displayers.Firewall{}))
	AddStringFlag(cmdFirewallCreate, dccncli.ArgFirewallName, "", "", "firewall name", requiredOpt())
	AddStringFlag(cmdFirewallCreate, dccncli.ArgInboundRules, "", "", "comma-separated key:value list, example value: protocol:tcp,ports:22,task_id:1,task_id:2,tag:frontend, use quoted string of space-separated values for multiple rules")
	AddStringFlag(cmdFirewallCreate, dccncli.ArgOutboundRules, "", "", "comma-separated key:value list, example value: protocol:tcp,ports:22,address:0.0.0.0/0, use quoted string of space-separated values for multiple rules")
	AddStringSliceFlag(cmdFirewallCreate, dccncli.ArgTaskIDs, "", []string{}, "comma-separated list of task IDs, example value: 123,456")
	AddStringSliceFlag(cmdFirewallCreate, dccncli.ArgTagNames, "", []string{}, "comma-separated list of tag names, example value: frontend,backend")

	cmdFirewallUpdate := CmdBuilder(cmd, RunFirewallUpdate, "update <id>", "update firewall", Writer, aliasOpt("u"), displayerType(&displayers.Firewall{}))
	AddStringFlag(cmdFirewallUpdate, dccncli.ArgFirewallName, "", "", "firewall name", requiredOpt())
	AddStringFlag(cmdFirewallUpdate, dccncli.ArgInboundRules, "", "", "comma-separated key:value list, example value: protocol:tcp,ports:22,task_id:123, use quoted string of space-separated values for multiple rules")
	AddStringFlag(cmdFirewallUpdate, dccncli.ArgOutboundRules, "", "", "comma-separated key:value list, example value: protocol:tcp,ports:22,address:0.0.0.0/0, use quoted string of space-separated values for multiple rules")
	AddStringSliceFlag(cmdFirewallUpdate, dccncli.ArgTaskIDs, "", []string{}, "comma-separated list of task IDs, example value: 123,456")
	AddStringSliceFlag(cmdFirewallUpdate, dccncli.ArgTagNames, "", []string{}, "comma-separated list of tag names, example value: frontend,backend")

	CmdBuilder(cmd, RunFirewallList, "list", "list firewalls", Writer, aliasOpt("ls"), displayerType(&displayers.Firewall{}))

	CmdBuilder(cmd, RunFirewallListByTask, "list-by-task <task_id>", "list firewalls by task ID", Writer, displayerType(&displayers.Firewall{}))

	cmdRunRecordDelete := CmdBuilder(cmd, RunFirewallDelete, "delete <id> [id ...]", "delete firewall", Writer, aliasOpt("d", "rm"))
	AddBoolFlag(cmdRunRecordDelete, dccncli.ArgForce, dccncli.ArgShortForce, false, "Force firewall delete")

	cmdAddTasks := CmdBuilder(cmd, RunFirewallAddTasks, "add-tasks <id>", "add tasks to the firewall", Writer)
	AddStringSliceFlag(cmdAddTasks, dccncli.ArgTaskIDs, "", []string{}, "comma-separated list of task IDs, example valus: 123,456")

	cmdRemoveTasks := CmdBuilder(cmd, RunFirewallRemoveTasks, "remove-tasks <id>", "remove tasks from the firewall", Writer)
	AddStringSliceFlag(cmdRemoveTasks, dccncli.ArgTaskIDs, "", []string{}, "comma-separated list of task IDs, example value: 123,456")

	cmdAddTags := CmdBuilder(cmd, RunFirewallAddTags, "add-tags <id>", "add tags to the firewall", Writer)
	AddStringSliceFlag(cmdAddTags, dccncli.ArgTagNames, "", []string{}, "comma-separated list of tag names, example valus: frontend,backend")

	cmdRemoveTags := CmdBuilder(cmd, RunFirewallRemoveTags, "remove-tags <id>", "remove tags from the firewall", Writer)
	AddStringSliceFlag(cmdRemoveTags, dccncli.ArgTagNames, "", []string{}, "comma-separated list of tag names, example value: frontend,backend")

	cmdAddRules := CmdBuilder(cmd, RunFirewallAddRules, "add-rules <id>", "add inbound/outbound rules to the firewall", Writer)
	AddStringFlag(cmdAddRules, dccncli.ArgInboundRules, "", "", "comma-separated key:value list, example value: protocol:tcp,ports:22,address:0.0.0.0/0, use quoted string of space-separated values for multiple rules")
	AddStringFlag(cmdAddRules, dccncli.ArgOutboundRules, "", "", "comma-separated key:value list, example value: protocol:tcp,ports:22,address:0.0.0.0/0, use quoted string of space-separated values for multiple rules")

	cmdRemoveRules := CmdBuilder(cmd, RunFirewallRemoveRules, "remove-rules <id>", "remove inbound/outbound rules from the firewall", Writer)
	AddStringFlag(cmdRemoveRules, dccncli.ArgInboundRules, "", "", "comma-separated key:value list, example value: protocol:tcp,ports:22,load_balancer_uid:lb-uuid, use quoted string of space-separated values for multiple rules")
	AddStringFlag(cmdRemoveRules, dccncli.ArgOutboundRules, "", "", "comma-separated key:value list, example value: protocol:tcp,ports:22,address:0.0.0.0/0, use quoted string of space-separated values for multiple rules")

	return cmd
}

// RunFirewallGet retrieves an existing Firewall by its identifier.
func RunFirewallGet(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	id := c.Args[0]

	fs := c.Firewalls()
	f, err := fs.Get(id)
	if err != nil {
		return err
	}

	item := &displayers.Firewall{Firewalls: do.Firewalls{*f}}
	return c.Display(item)
}

// RunFirewallCreate creates a new Firewall with a given configuration.
func RunFirewallCreate(c *CmdConfig) error {
	r := new(godo.FirewallRequest)
	if err := buildFirewallRequestFromArgs(c, r); err != nil {
		return err
	}

	fs := c.Firewalls()
	f, err := fs.Create(r)
	if err != nil {
		return err
	}

	item := &displayers.Firewall{Firewalls: do.Firewalls{*f}}
	return c.Display(item)
}

// RunFirewallUpdate updates an existing Firewall with new configuration.
func RunFirewallUpdate(c *CmdConfig) error {
	if len(c.Args) == 0 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	fID := c.Args[0]

	r := new(godo.FirewallRequest)
	if err := buildFirewallRequestFromArgs(c, r); err != nil {
		return err
	}

	fs := c.Firewalls()
	f, err := fs.Update(fID, r)
	if err != nil {
		return err
	}

	item := &displayers.Firewall{Firewalls: do.Firewalls{*f}}
	return c.Display(item)
}

// RunFirewallList lists Firewalls.
func RunFirewallList(c *CmdConfig) error {
	fs := c.Firewalls()
	list, err := fs.List()
	if err != nil {
		return err
	}

	items := &displayers.Firewall{Firewalls: list}
	return c.Display(items)
}

// RunFirewallListByTask lists Firewalls for a given Task.
func RunFirewallListByTask(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	dID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return fmt.Errorf("invalid task id: [%v]", c.Args[0])
	}

	fs := c.Firewalls()
	list, err := fs.ListByTask(dID)
	if err != nil {
		return err
	}

	items := &displayers.Firewall{Firewalls: list}
	return c.Display(items)
}

// RunFirewallDelete deletes a Firewall by its identifier.
func RunFirewallDelete(c *CmdConfig) error {
	if len(c.Args) < 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	force, err := c.Ankr.GetBool(c.NS, dccncli.ArgForce)
	if err != nil {
		return err
	}

	fs := c.Firewalls()
	if force || AskForConfirm(fmt.Sprintf("delete %d firewall(s)", len(c.Args))) == nil {
		for _, id := range c.Args {
			if err := fs.Delete(id); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("operation aborted")
	}

	return nil
}

// RunFirewallAddTasks adds tasks to a Firewall.
func RunFirewallAddTasks(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	fID := c.Args[0]

	taskIDsList, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgTaskIDs)
	if err != nil {
		return err
	}

	taskIDs, err := extractTaskIDs(taskIDsList)
	if err != nil {
		return err
	}

	return c.Firewalls().AddTasks(fID, taskIDs...)
}

// RunFirewallRemoveTasks removes tasks from a Firewall.
func RunFirewallRemoveTasks(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	fID := c.Args[0]

	taskIDsList, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgTaskIDs)
	if err != nil {
		return err
	}

	taskIDs, err := extractTaskIDs(taskIDsList)
	if err != nil {
		return err
	}

	return c.Firewalls().RemoveTasks(fID, taskIDs...)
}

// RunFirewallAddTags adds tags to a Firewall.
func RunFirewallAddTags(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	fID := c.Args[0]

	tagList, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgTagNames)
	if err != nil {
		return err
	}

	return c.Firewalls().AddTags(fID, tagList...)
}

// RunFirewallRemoveTags removes tags from a Firewall.
func RunFirewallRemoveTags(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	fID := c.Args[0]

	tagList, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgTagNames)
	if err != nil {
		return err
	}

	return c.Firewalls().RemoveTags(fID, tagList...)
}

// RunFirewallAddRules adds rules to a Firewall.
func RunFirewallAddRules(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	fID := c.Args[0]

	rr := new(godo.FirewallRulesRequest)
	if err := buildFirewallRulesRequestFromArgs(c, rr); err != nil {
		return err
	}

	return c.Firewalls().AddRules(fID, rr)
}

// RunFirewallRemoveRules removes rules from a Firewall.
func RunFirewallRemoveRules(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	fID := c.Args[0]

	rr := new(godo.FirewallRulesRequest)
	if err := buildFirewallRulesRequestFromArgs(c, rr); err != nil {
		return err
	}

	return c.Firewalls().RemoveRules(fID, rr)
}

func buildFirewallRequestFromArgs(c *CmdConfig, r *godo.FirewallRequest) error {
	name, err := c.Ankr.GetString(c.NS, dccncli.ArgFirewallName)
	if err != nil {
		return err
	}
	r.Name = name

	ira, err := c.Ankr.GetString(c.NS, dccncli.ArgInboundRules)
	if err != nil {
		return err
	}

	inboundRules, err := extractInboundRules(ira)
	if err != nil {
		return err
	}
	r.InboundRules = inboundRules

	ora, err := c.Ankr.GetString(c.NS, dccncli.ArgOutboundRules)
	if err != nil {
		return err
	}

	outboundRules, err := extractOutboundRules(ora)
	if err != nil {
		return err
	}
	r.OutboundRules = outboundRules

	taskIDsList, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgTaskIDs)
	if err != nil {
		return err
	}

	taskIDs, err := extractTaskIDs(taskIDsList)
	if err != nil {
		return err
	}
	r.TaskIDs = taskIDs

	tagsList, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgTagNames)
	if err != nil {
		return err
	}
	r.Tags = tagsList

	return nil
}

func buildFirewallRulesRequestFromArgs(c *CmdConfig, rr *godo.FirewallRulesRequest) error {
	ira, err := c.Ankr.GetString(c.NS, dccncli.ArgInboundRules)
	if err != nil {
		return err
	}

	inboundRules, err := extractInboundRules(ira)
	if err != nil {
		return err
	}
	rr.InboundRules = inboundRules

	ora, err := c.Ankr.GetString(c.NS, dccncli.ArgOutboundRules)
	if err != nil {
		return err
	}

	outboundRules, err := extractOutboundRules(ora)
	if err != nil {
		return err
	}
	rr.OutboundRules = outboundRules

	return nil
}

func extractInboundRules(s string) (rules []godo.InboundRule, err error) {
	if len(s) == 0 {
		return nil, nil
	}

	list := strings.Split(s, " ")
	for _, v := range list {
		rule, err := extractRule(v, "sources")
		if err != nil {
			return nil, err
		}
		mr, _ := json.Marshal(rule)
		ir := &godo.InboundRule{}
		json.Unmarshal(mr, ir)
		rules = append(rules, *ir)
	}

	return rules, nil
}

func extractOutboundRules(s string) (rules []godo.OutboundRule, err error) {
	if len(s) == 0 {
		return nil, nil
	}

	list := strings.Split(s, " ")
	for _, v := range list {
		rule, err := extractRule(v, "destinations")
		if err != nil {
			return nil, err
		}
		mr, _ := json.Marshal(rule)
		or := &godo.OutboundRule{}
		json.Unmarshal(mr, or)
		rules = append(rules, *or)
	}

	return rules, nil
}

func extractRule(ruleStr string, sd string) (map[string]interface{}, error) {
	rule := map[string]interface{}{}
	var taskIDs []int
	var addresses, lbUIDs, tags []string

	kvs := strings.Split(ruleStr, ",")
	for _, v := range kvs {
		pair := strings.SplitN(v, ":", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("Unexpected input value [%v], must be a key:value pair", pair)
		}

		switch pair[0] {
		case "address":
			addresses = append(addresses, pair[1])
		case "task_id":
			i, err := strconv.Atoi(pair[1])
			if err != nil {
				return nil, fmt.Errorf("Provided value [%v] for task id is not of type int", pair[0])
			}
			taskIDs = append(taskIDs, i)
		case "load_balancer_uid":
			lbUIDs = append(lbUIDs, pair[1])
		case "tag":
			tags = append(tags, pair[1])
		default:
			rule[pair[0]] = pair[1]
		}
	}

	rule[sd] = map[string]interface{}{
		"addresses":          addresses,
		"task_ids":        taskIDs,
		"load_balancer_uids": lbUIDs,
		"tags":               tags,
	}

	return rule, nil
}
