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
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/godo"
	"github.com/gobwas/glob"
	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
)

// Task creates the task command.
func Task() *Command {
	//DCCN-CLI task
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "task",
			Aliases: []string{"t"},
			Short:   "task commands",
			Long:    "task is used to access task commands",
		},
		DocCategories: []string{"task"},
		IsIndex:       true,
	}

	//CmdBuilder(cmd, RunTaskActions, "actions <task-id>", "task actions", Writer,
	//	aliasOpt("a"), displayerType(&displayers.Action{}), docCategories("task"))

	//CmdBuilder(cmd, RunTaskBackups, "backups <task-id>", "task backups", Writer,
	//	aliasOpt("b"), displayerType(&displayers.Image{}), docCategories("task"))
	//DCCN-CLI comput task create
	cmdTaskCreate := CmdBuilder(cmd, RunTaskCreate, "create <task-name> [task-name ...]", "create task", Writer,
		aliasOpt("cr"), displayerType(&displayers.Task{}), docCategories("task"))
	//AddStringSliceFlag(cmdTaskCreate, dccncli.ArgSSHKeys, "", []string{}, "SSH Keys or fingerprints")
	//AddStringFlag(cmdTaskCreate, dccncli.ArgUserData, "", "", "User data")
	//AddStringFlag(cmdTaskCreate, dccncli.ArgUserDataFile, "", "", "User data file")
	AddBoolFlag(cmdTaskCreate, dccncli.ArgCommandWait, "", false, "Wait for task to be created")
	AddStringFlag(cmdTaskCreate, dccncli.ArgRegionSlug, "", "", "Task region",
		requiredOpt())
	AddStringFlag(cmdTaskCreate, dccncli.ArgZoneSlug, "", "", "Task zone",
		requiredOpt())
	//AddStringFlag(cmdTaskCreate, dccncli.ArgSizeSlug, "", "", "Task size")
	//	requiredOpt())
	//AddBoolFlag(cmdTaskCreate, dccncli.ArgBackups, "", false, "Backup task")
	//AddBoolFlag(cmdTaskCreate, dccncli.ArgIPv6, "", false, "IPv6 support")
	//AddBoolFlag(cmdTaskCreate, dccncli.ArgPrivateNetworking, "", false, "Private networking")
	//AddBoolFlag(cmdTaskCreate, dccncli.ArgMonitoring, "", false, "Monitoring")
	//AddStringFlag(cmdTaskCreate, dccncli.ArgImage, "", "", "Task image")
	//	requiredOpt())
	//AddStringFlag(cmdTaskCreate, dccncli.ArgTagName, "", "", "Tag name")
	//AddStringSliceFlag(cmdTaskCreate, dccncli.ArgTagNames, "", []string{}, "Tag names")

	//AddStringSliceFlag(cmdTaskCreate, dccncli.ArgVolumeList, "", []string{}, "Volumes to attach")
	//DCCN-CLI comput task delete
	cmdRunTaskDelete := CmdBuilder(cmd, RunTaskDelete, "delete <task-id|task-name> [task-id|task-name ...]", "Delete task by id or name", Writer,
		aliasOpt("d", "del", "rm"), docCategories("Task"))
	AddBoolFlag(cmdRunTaskDelete, dccncli.ArgForce, dccncli.ArgShortForce, false, "Force task delete")

	//cmdRunTaskGet := CmdBuilder(cmd, RunTaskGet, "get <task-id>", "get task", Writer,
	//	aliasOpt("g"), displayerType(&displayers.Task{}), docCategories("task"))
	//AddStringFlag(cmdRunTaskGet, dccncli.ArgTemplate, "", "", "Go template format. Few sample values:{{.ID}} {{.Name}} {{.Memory}} {{.Region.Name}} {{.Image}} {{.Tags}}")

	//CmdBuilder(cmd, RunTaskKernels, "kernels <task-id>", "task kernels", Writer,
	//	aliasOpt("k"), displayerType(&displayers.Kernel{}), docCategories("task"))
	//DCCN-CLI task list
	cmdRunTaskList := CmdBuilder(cmd, RunTaskList, "list [GLOB]", "list tasks", Writer,
		aliasOpt("ls"), displayerType(&displayers.Task{}), docCategories("task"))
	_ = cmdRunTaskList
	//AddStringFlag(cmdRunTaskList, dccncli.ArgRegionSlug, "", "", "Task region")
	//AddStringFlag(cmdRunTaskList, dccncli.ArgTagName, "", "", "Tag name")

	//CmdBuilder(cmd, RunTaskNeighbors, "neighbors <task-id>", "task neighbors", Writer,
	//	aliasOpt("n"), displayerType(&displayers.Task{}), docCategories("task"))

	//CmdBuilder(cmd, RunTaskSnapshots, "snapshots <task-id>", "snapshots", Writer,
	//	aliasOpt("s"), displayerType(&displayers.Image{}), docCategories("task"))

	//cmdRunTaskTag := CmdBuilder(cmd, RunTaskTag, "tag <task-id|task-name>", "tag", Writer,
	//	docCategories("task"))
	//AddStringFlag(cmdRunTaskTag, dccncli.ArgTagName, "", "", "Tag name",
	//	requiredOpt())

	//cmdRunTaskUntag := CmdBuilder(cmd, RunTaskUntag, "untag <task-id|task-name>", "untag", Writer,
	//	docCategories("task"))
	//AddStringSliceFlag(cmdRunTaskUntag, dccncli.ArgTagName, "", []string{}, "tag names")

	return cmd
}

// RunTaskActions returns a list of actions for a task.
func RunTaskActions(c *CmdConfig) error {

	ds := c.Tasks()

	id, err := getTaskIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	list, err := ds.Actions(id)
	if err != nil {
		return err
	}
	item := &displayers.Action{Actions: list}
	return c.Display(item)
}

// RunTaskBackups returns a list of backup images for a task.
func RunTaskBackups(c *CmdConfig) error {

	ds := c.Tasks()

	id, err := getTaskIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	list, err := ds.Backups(id)
	if err != nil {
		return err
	}

	item := &displayers.Image{Images: list}
	return c.Display(item)
}

// RunTaskCreate creates a task.
//DCCN-CLI comput task create
func RunTaskCreate(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	region, err := c.Ankr.GetString(c.NS, dccncli.ArgRegionSlug)
	if err != nil {
		return err
	}
	
	zone, err := c.Ankr.GetString(c.NS, dccncli.ArgZoneSlug)
	if err != nil {
		return err
	}

	size, err := c.Ankr.GetString(c.NS, dccncli.ArgSizeSlug)
	if err != nil {
		return err
	}
	size = "s-1vcpu-3gb"
	backups, err := c.Ankr.GetBool(c.NS, dccncli.ArgBackups)
	if err != nil {
		return err
	}

	ipv6, err := c.Ankr.GetBool(c.NS, dccncli.ArgIPv6)
	if err != nil {
		return err
	}

	privateNetworking, err := c.Ankr.GetBool(c.NS, dccncli.ArgPrivateNetworking)
	if err != nil {
		return err
	}

	monitoring, err := c.Ankr.GetBool(c.NS, dccncli.ArgMonitoring)
	if err != nil {
		return err
	}

	keys, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgSSHKeys)
	if err != nil {
		return err
	}

	tagName, err := c.Ankr.GetString(c.NS, dccncli.ArgTagName)
	if err != nil {
		return err
	}

	tagNames, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgTagNames)
	if err != nil {
		return err
	}

	sshKeys := extractSSHKeys(keys)

	userData, err := c.Ankr.GetString(c.NS, dccncli.ArgUserData)
	if err != nil {
		return err
	}

	volumeList, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgVolumeList)
	if err != nil {
		return err
	}
	volumes := extractVolumes(volumeList)

	filename, err := c.Ankr.GetString(c.NS, dccncli.ArgUserDataFile)
	if err != nil {
		return err
	}

	userData, err = extractUserData(userData, filename)
	if err != nil {
		return err
	}

	imageStr, err := c.Ankr.GetString(c.NS, dccncli.ArgImage)
	if err != nil {
		return err
	}
	imageStr = "ubuntu-16-04-x64"
	createImage := godo.TaskCreateImage{Slug: imageStr}

	i, err := strconv.Atoi(imageStr)
	if err == nil {
		createImage = godo.TaskCreateImage{ID: i}
	}

	wait, err := c.Ankr.GetBool(c.NS, dccncli.ArgCommandWait)
	if err != nil {
		return err
	}

	ds := c.Tasks()
	ts := c.Tags()

	var wg sync.WaitGroup
	var createdList do.Tasks
	errs := make(chan error, len(c.Args))
	for _, name := range c.Args {
		dcr := &godo.TaskCreateRequest{
			Name:              name,
			Region:            region,
			Zone:			   zone,
			Size:              size,
			Image:             createImage,
			Volumes:           volumes,
			Backups:           backups,
			IPv6:              ipv6,
			PrivateNetworking: privateNetworking,
			Monitoring:        monitoring,
			SSHKeys:           sshKeys,
			UserData:          userData,
			Tags:              tagNames,
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if tagName != "" {
				tag, err := ts.Get(tagName)
				if err != nil {
					errs <- err
					return
				}
				if tag == nil {
					errs <- fmt.Errorf("Specified Tag must exist")
					return
				}
			}
			d, err := ds.Create(dcr, wait)
			if err != nil {
				errs <- err
				return
			}
			if (d.Status == "Success"){
				fmt.Printf("Task id %d created successfully. \n", d.ID)
			}else{
				fmt.Printf("Fail to create task. \n")
			}
			
			if tagName != "" {
				trr := &godo.TagResourcesRequest{
					Resources: []godo.Resource{
						{ID: strconv.Itoa(d.ID), Type: godo.TaskResourceType},
					},
				}

				err := ts.TagResources(tagName, trr)
				if err != nil {
					errs <- err
				}

			}

			createdList = append(createdList, *d)
		}()
	}

	wg.Wait()
	close(errs)

	//item := &displayers.Task{Tasks: createdList}

	for err := range errs {
		if err != nil {
			return err
		}
	}
	
	//c.Display(item)

	return nil
}

// RunTaskTag adds a tag to a task.
func RunTaskTag(c *CmdConfig) error {
	ds := c.Tasks()
	ts := c.Tags()

	if len(c.Args) < 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	tag, err := c.Ankr.GetString(c.NS, dccncli.ArgTagName)
	if err != nil {
		return err
	}

	fn := func(ids []int) error {
		trr := &godo.TagResourcesRequest{}
		for _, id := range ids {
			r := godo.Resource{
				ID:   strconv.Itoa(id),
				Type: godo.TaskResourceType,
			}
			trr.Resources = append(trr.Resources, r)
		}

		return ts.TagResources(tag, trr)
	}

	return matchTasks(c.Args, ds, fn)
}

// RunTaskUntag untags a task.
func RunTaskUntag(c *CmdConfig) error {
	ds := c.Tasks()
	ts := c.Tags()

	if len(c.Args) < 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	taskIDStrs := c.Args

	tagNames, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgTagName)
	if err != nil {
		return err
	}

	fn := func(ids []int) error {
		urr := &godo.UntagResourcesRequest{}

		for _, id := range ids {
			for _, tagName := range tagNames {
				r := godo.Resource{
					ID:   strconv.Itoa(id),
					Type: godo.TaskResourceType,
				}

				urr.Resources = append(urr.Resources, r)

				err := ts.UntagResources(tagName, urr)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	return matchTasks(taskIDStrs, ds, fn)
}

func extractSSHKeys(keys []string) []godo.TaskCreateSSHKey {
	sshKeys := []godo.TaskCreateSSHKey{}

	for _, k := range keys {
		if i, err := strconv.Atoi(k); err == nil {
			if i > 0 {
				sshKeys = append(sshKeys, godo.TaskCreateSSHKey{ID: i})
			}
			continue
		}

		if k != "" {
			sshKeys = append(sshKeys, godo.TaskCreateSSHKey{Fingerprint: k})
		}
	}

	return sshKeys
}

func extractUserData(userData, filename string) (string, error) {
	if userData == "" && filename != "" {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", err
		}
		userData = string(data)
	}

	return userData, nil
}

func extractVolumes(volumeList []string) []godo.TaskCreateVolume {
	var volumes []godo.TaskCreateVolume

	for _, v := range volumeList {
		var req godo.TaskCreateVolume
		if uuid.Parse(v) != nil {
			req.ID = v
		} else {
			req.Name = v
		}
		volumes = append(volumes, req)
	}

	return volumes
}

func allInt(in []string) ([]int, error) {
	out := []int{}
	seen := map[string]bool{}

	for _, i := range in {
		if seen[i] {
			continue
		}

		seen[i] = true

		id, err := strconv.Atoi(i)
		if err != nil {
			return nil, fmt.Errorf("%s is not an int", i)
		}
		out = append(out, id)
	}
	return out, nil
}

// RunTaskDelete destroy a task by id.
func RunTaskDelete(c *CmdConfig) error {
	ds := c.Tasks()

	force, err := c.Ankr.GetBool(c.NS, dccncli.ArgForce)
	if err != nil {
		return err
	}

	tagName, err := c.Ankr.GetString(c.NS, dccncli.ArgTagName)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 && tagName == "" {
		return dccncli.NewMissingArgsErr(c.NS)
	} else if len(c.Args) > 0 && tagName != "" {
		return fmt.Errorf("please specify task identifiers or a tag name")
	} else if tagName != "" {
		if force || AskForConfirm("delete task by \""+tagName+"\" tag") == nil {
			return ds.DeleteByTag(tagName)
		}
		return nil
	}

	if force || AskForConfirm(fmt.Sprintf("delete %d task(s)", len(c.Args))) == nil {

		fn := func(ids []int) error {
			for _, id := range ids {
				if status, err := ds.Delete(id); err != nil {
					return fmt.Errorf("unable to delete task %d: %v", id, err)
				}else{
					fmt.Printf("Delete task id %d...%s! \n", id, status)
				}
			}
			return nil
		}
		if extractedIDs, err := allInt(c.Args); err == nil {
			return fn(extractedIDs)
		}
		return err
		//return matchTasks(c.Args, ds, fn)
	}
	return fmt.Errorf("operation aborted")

	return nil

}

type matchTasksFn func(ids []int) error

func matchTasks(ids []string, ds do.TasksService, fn matchTasksFn) error {
	if extractedIDs, err := allInt(ids); err == nil {
		return fn(extractedIDs)
	}

	sum, err := buildTaskSummary(ds)
	if err != nil {
		return err
	}

	matchedMap := map[int]bool{}
	for _, idStr := range ids {
		count := sum.count[idStr]
		if count > 1 {
			return fmt.Errorf("there are %d Tasks with the name %q, please delete by id. [%s]",
				count, idStr, strings.Join(sum.byName[idStr], ", "))
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			id, ok := sum.byID[idStr]
			if !ok {
				return fmt.Errorf("task with name %q could not be found", idStr)
			}

			matchedMap[id] = true
			continue
		}

		matchedMap[id] = true
	}

	var extractedIDs []int
	for id := range matchedMap {
		extractedIDs = append(extractedIDs, id)
	}

	sort.Ints(extractedIDs)
	return fn(extractedIDs)
}

// RunTaskGet returns a task.
func RunTaskGet(c *CmdConfig) error {
	id, err := getTaskIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	getTemplate, err := c.Ankr.GetString(c.NS, dccncli.ArgTemplate)
	if err != nil {
		return err
	}

	ds := c.Tasks()

	d, err := ds.Get(id)
	if err != nil {
		return err
	}

	item := &displayers.Task{Tasks: do.Tasks{*d}}
	if getTemplate != "" {
		t := template.New("get template")
		t, err = t.Parse(getTemplate)
		if err != nil {
			return err
		}
		return t.Execute(c.Out, d)
	}
	return c.Display(item)
}

// RunTaskKernels returns a list of available kernels for a task.
func RunTaskKernels(c *CmdConfig) error {

	ds := c.Tasks()
	id, err := getTaskIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	list, err := ds.Kernels(id)
	if err != nil {
		return err
	}

	item := &displayers.Kernel{Kernels: list}
	return c.Display(item)
}

// RunTaskList returns a list of tasks.
func RunTaskList(c *CmdConfig) error {

	ds := c.Tasks()

	region, err := c.Ankr.GetString(c.NS, dccncli.ArgRegionSlug)
	if err != nil {
		return err
	}

	tagName, err := c.Ankr.GetString(c.NS, dccncli.ArgTagName)
	if err != nil {
		return err
	}

	matches := []glob.Glob{}
	for _, globStr := range c.Args {
		g, err := glob.Compile(globStr)
		if err != nil {
			return fmt.Errorf("unknown glob %q", globStr)
		}

		matches = append(matches, g)
	}

	var matchedList do.Tasks

	var list do.Tasks
	if tagName == "" {
		list, err = ds.List()
		if err != nil {
			return err
		}
	} else {
		list, err = ds.ListByTag(tagName)
	}

	for _, task := range list {
		var skip = true
		if len(matches) == 0 {
			skip = false
		} else {
			for _, m := range matches {
				if m.Match(task.Taskname) {
					skip = false
				}
			}
		}

		if !skip && region != "" {
			if region != task.Region.Slug {
				skip = true
			}
		}

		if !skip {
			matchedList = append(matchedList, task)
		}
	}

	item := &displayers.Task{Tasks: matchedList}
	return c.Display(item)
}

// RunTaskNeighbors returns a list of task neighbors.
func RunTaskNeighbors(c *CmdConfig) error {

	ds := c.Tasks()

	id, err := getTaskIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	list, err := ds.Neighbors(id)
	if err != nil {
		return err
	}

	item := &displayers.Task{Tasks: list}
	return c.Display(item)
}

// RunTaskSnapshots returns a list of available kernels for a task.
func RunTaskSnapshots(c *CmdConfig) error {

	ds := c.Tasks()
	id, err := getTaskIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	list, err := ds.Snapshots(id)
	if err != nil {
		return err
	}

	item := &displayers.Image{Images: list}
	return c.Display(item)
}

func getTaskIDArg(ns string, args []string) (int, error) {
	if len(args) != 1 {
		return 0, dccncli.NewMissingArgsErr(ns)
	}

	return strconv.Atoi(args[0])
}

type taskSummary struct {
	count  map[string]int
	byID   map[string]int
	byName map[string][]string
}

func buildTaskSummary(ds do.TasksService) (*taskSummary, error) {
	list, err := ds.List()
	if err != nil {
		return nil, err
	}

	var sum taskSummary

	sum.count = map[string]int{}
	sum.byID = map[string]int{}
	sum.byName = map[string][]string{}
	for _, d := range list {
		sum.count[d.Name]++
		sum.byID[d.Name] = d.ID
		sum.byName[d.Name] = append(sum.byName[d.Name], strconv.Itoa(d.ID))
	}

	return &sum, nil
}
