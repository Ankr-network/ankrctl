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
	"strconv"
	"testing"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/godo"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	testImage = do.Image{Image: &godo.Image{
		ID:      1,
		Slug:    "slug",
		Regions: []string{"test0"},
	}}
	testImageSecondary = do.Image{Image: &godo.Image{
		ID:      2,
		Slug:    "slug-secondary",
		Regions: []string{"test0"},
	}}
	testImageList = do.Images{testImage, testImageSecondary}
)

func TestTaskCommand(t *testing.T) {
	cmd := Task()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "actions", "backups", "create", "delete", "get", "kernels", "list", "neighbors", "snapshots", "tag", "untag")
}

func TestTaskActionList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Actions", 1).Return(testActionList, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskActions(config)
		assert.NoError(t, err)
	})
}

func TestTaskBackupList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Backups", 1).Return(testImageList, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskBackups(config)
		assert.NoError(t, err)
	})
}

func TestTaskCreate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		volumeUUID := uuid.New()
		dcr := &godo.TaskCreateRequest{
			Name:    "task",
			Region:  "dev0",
			Size:    "1gb",
			Image:   godo.TaskCreateImage{ID: 0, Slug: "image"},
			SSHKeys: []godo.TaskCreateSSHKey{},
			Volumes: []godo.TaskCreateVolume{
				{Name: "test-volume"},
				{ID: volumeUUID},
			},
			Backups:           false,
			IPv6:              false,
			PrivateNetworking: false,
			Monitoring:        false,
			UserData:          "#cloud-config",
			Tags:              []string{"one", "two"},
		}
		tm.tasks.On("Create", dcr, false).Return(&testTask, nil)

		config.Args = append(config.Args, "task")

		config.Ankr.Set(config.NS, dccncli.ArgRegionSlug, "dev0")
		config.Ankr.Set(config.NS, dccncli.ArgSizeSlug, "1gb")
		config.Ankr.Set(config.NS, dccncli.ArgImage, "image")
		config.Ankr.Set(config.NS, dccncli.ArgUserData, "#cloud-config")
		config.Ankr.Set(config.NS, dccncli.ArgVolumeList, []string{"test-volume", volumeUUID})
		config.Ankr.Set(config.NS, dccncli.ArgTagNames, []string{"one", "two"})

		err := RunTaskCreate(config)
		assert.NoError(t, err)
	})
}

func TestTaskCreateWithTag(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		dcr := &godo.TaskCreateRequest{Name: "task", Region: "dev0", Size: "1gb", Image: godo.TaskCreateImage{ID: 0, Slug: "image"}, SSHKeys: []godo.TaskCreateSSHKey{}, Backups: false, IPv6: false, PrivateNetworking: false, UserData: "#cloud-config"}
		tm.tasks.On("Create", dcr, false).Return(&testTask, nil)
		tm.tags.On("Get", "my-tag").Return(&testTag, nil)

		trr := &godo.TagResourcesRequest{
			Resources: []godo.Resource{
				{ID: "1", Type: godo.TaskResourceType},
			},
		}
		tm.tags.On("TagResources", "my-tag", trr).Return(nil)

		config.Args = append(config.Args, "task")

		config.Ankr.Set(config.NS, dccncli.ArgRegionSlug, "dev0")
		config.Ankr.Set(config.NS, dccncli.ArgSizeSlug, "1gb")
		config.Ankr.Set(config.NS, dccncli.ArgImage, "image")
		config.Ankr.Set(config.NS, dccncli.ArgUserData, "#cloud-config")
		config.Ankr.Set(config.NS, dccncli.ArgTagName, "my-tag")

		err := RunTaskCreate(config)
		assert.NoError(t, err)
	})
}

func TestTaskCreateUserDataFile(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		dcr := &godo.TaskCreateRequest{Name: "task", Region: "dev0", Size: "1gb", Image: godo.TaskCreateImage{ID: 0, Slug: "image"}, SSHKeys: []godo.TaskCreateSSHKey{}, Backups: false, IPv6: false, PrivateNetworking: false, UserData: "#cloud-config\n\ncoreos:\n  etcd2:\n    # generate a new token for each unique cluster from https://discovery.etcd.io/new?size=5\n    # specify the initial size of your cluster with ?size=X\n    discovery: https://discovery.etcd.io/<token>\n    # multi-region and multi-cloud deployments need to use $public_ipv4\n    advertise-client-urls: http://$private_ipv4:2379,http://$private_ipv4:4001\n    initial-advertise-peer-urls: http://$private_ipv4:2380\n    # listen on both the official ports and the legacy ports\n    # legacy ports can be omitted if your application doesn't depend on them\n    listen-client-urls: http://0.0.0.0:2379,http://0.0.0.0:4001\n    listen-peer-urls: http://$private_ipv4:2380\n  units:\n    - name: etcd2.service\n      command: start\n    - name: fleet.service\n      command: start\n"}
		tm.tasks.On("Create", dcr, false).Return(&testTask, nil)

		config.Args = append(config.Args, "task")

		config.Ankr.Set(config.NS, dccncli.ArgRegionSlug, "dev0")
		config.Ankr.Set(config.NS, dccncli.ArgSizeSlug, "1gb")
		config.Ankr.Set(config.NS, dccncli.ArgImage, "image")
		config.Ankr.Set(config.NS, dccncli.ArgUserDataFile, "../testdata/cloud-config.yml")

		err := RunTaskCreate(config)
		assert.NoError(t, err)
	})
}

func TestTaskDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Delete", 1).Return(nil)

		config.Args = append(config.Args, strconv.Itoa(testTask.ID))
		config.Ankr.Set(config.NS, dccncli.ArgForce, true)

		err := RunTaskDelete(config)
		assert.NoError(t, err)

	})
}

func TestTaskDeleteByTag(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("DeleteByTag", "my-tag").Return(nil)

		config.Ankr.Set(config.NS, dccncli.ArgTagName, "my-tag")
		config.Ankr.Set(config.NS, dccncli.ArgForce, true)

		err := RunTaskDelete(config)
		assert.NoError(t, err)
	})

}

func TestTaskDeleteRepeatedID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Delete", 1).Return(nil).Once()

		id := strconv.Itoa(testTask.ID)
		config.Args = append(config.Args, id, id)
		config.Ankr.Set(config.NS, dccncli.ArgForce, true)

		err := RunTaskDelete(config)
		assert.NoError(t, err)
	})
}

func TestTaskDeleteByName(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("List").Return(testTaskList, nil)
		tm.tasks.On("Delete", 1).Return(nil)

		config.Args = append(config.Args, testTask.Name)
		config.Ankr.Set(config.NS, dccncli.ArgForce, true)

		err := RunTaskDelete(config)
		assert.NoError(t, err)
	})
}

func TestTaskDeleteByName_Ambiguous(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		list := do.Tasks{testTask, testTask}
		tm.tasks.On("List").Return(list, nil)

		config.Args = append(config.Args, testTask.Name)
		config.Ankr.Set(config.NS, dccncli.ArgForce, true)

		err := RunTaskDelete(config)
		t.Log(err)
		assert.Error(t, err)
	})
}

func TestTaskDelete_MixedNameAndType(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("List").Return(testTaskList, nil)
		tm.tasks.On("Delete", 1).Return(nil).Once()

		id := strconv.Itoa(testTask.ID)
		config.Args = append(config.Args, id, testTask.Name)
		config.Ankr.Set(config.NS, dccncli.ArgForce, true)

		err := RunTaskDelete(config)
		assert.NoError(t, err)
	})

}

func TestTaskGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Get", testTask.ID).Return(&testTask, nil)

		config.Args = append(config.Args, strconv.Itoa(testTask.ID))

		err := RunTaskGet(config)
		assert.NoError(t, err)
	})
}

func TestTaskGet_Template(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Get", testTask.ID).Return(&testTask, nil)

		config.Args = append(config.Args, strconv.Itoa(testTask.ID))
		config.Ankr.Set(config.NS, dccncli.ArgTemplate, "{{.Name}}")

		err := RunTaskGet(config)
		assert.NoError(t, err)
	})
}

func TestTaskKernelList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Kernels", testTask.ID).Return(testKernelList, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskKernels(config)
		assert.NoError(t, err)
	})
}

func TestTaskNeighbors(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Neighbors", testTask.ID).Return(testTaskList, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskNeighbors(config)
		assert.NoError(t, err)
	})
}

func TestTaskSnapshotList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Snapshots", testTask.ID).Return(testImageList, nil)

		config.Args = append(config.Args, "1")

		err := RunTaskSnapshots(config)
		assert.NoError(t, err)
	})
}

func TestTasksList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("List").Return(testTaskList, nil)

		err := RunTaskList(config)
		assert.NoError(t, err)
	})
}

func TestTasksListByTag(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("ListByTag", "my-tag").Return(testTaskList, nil)

		config.Ankr.Set(config.NS, dccncli.ArgTagName, "my-tag")

		err := RunTaskList(config)
		assert.NoError(t, err)
	})
}

func TestTasksTag(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		trr := &godo.TagResourcesRequest{
			Resources: []godo.Resource{
				{ID: "1", Type: godo.TaskResourceType},
			},
		}
		tm.tags.On("TagResources", "my-tag", trr).Return(nil)

		config.Args = append(config.Args, "1")
		config.Ankr.Set(config.NS, dccncli.ArgTagName, "my-tag")

		err := RunTaskTag(config)
		assert.NoError(t, err)
	})
}

func TestTasksTagMultiple(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		trr := &godo.TagResourcesRequest{
			Resources: []godo.Resource{
				{ID: "1", Type: godo.TaskResourceType},
				{ID: "2", Type: godo.TaskResourceType},
			},
		}
		tm.tags.On("TagResources", "my-tag", trr).Return(nil)

		config.Args = append(config.Args, "1")
		config.Args = append(config.Args, "2")
		config.Ankr.Set(config.NS, dccncli.ArgTagName, "my-tag")

		err := RunTaskTag(config)
		assert.NoError(t, err)
	})
}

func TestTasksTagByName(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		trr := &godo.TagResourcesRequest{
			Resources: []godo.Resource{
				{ID: "1", Type: godo.TaskResourceType},
			},
		}
		tm.tags.On("TagResources", "my-tag", trr).Return(nil)
		tm.tasks.On("List").Return(testTaskList, nil)

		config.Args = append(config.Args, testTask.Name)
		config.Ankr.Set(config.NS, dccncli.ArgTagName, "my-tag")

		err := RunTaskTag(config)
		assert.NoError(t, err)
	})
}

func TestTasksTagMultipleNameAndID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		trr := &godo.TagResourcesRequest{
			Resources: []godo.Resource{
				{ID: "1", Type: godo.TaskResourceType},
				{ID: "3", Type: godo.TaskResourceType},
			},
		}
		tm.tags.On("TagResources", "my-tag", trr).Return(nil)
		tm.tasks.On("List").Return(testTaskList, nil)

		config.Args = append(config.Args, testTask.Name)
		config.Args = append(config.Args, strconv.Itoa(anotherTestTask.ID))
		config.Ankr.Set(config.NS, dccncli.ArgTagName, "my-tag")

		err := RunTaskTag(config)
		assert.NoError(t, err)
	})
}

func TestTasksUntag(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		urr := &godo.UntagResourcesRequest{
			Resources: []godo.Resource{
				{ID: "1", Type: godo.TaskResourceType},
			},
		}

		tm.tags.On("UntagResources", "my-tag", urr).Return(nil)
		tm.tasks.On("List").Return(testTaskList, nil)

		config.Args = []string{testTask.Name}
		config.Ankr.Set(config.NS, dccncli.ArgTagName, "my-tag")

		err := RunTaskUntag(config)
		assert.NoError(t, err)
	})
}

func Test_extractSSHKey(t *testing.T) {
	cases := []struct {
		in       []string
		expected []godo.TaskCreateSSHKey
	}{
		{
			in:       []string{"1"},
			expected: []godo.TaskCreateSSHKey{{ID: 1}},
		},
		{
			in:       []string{"fingerprint"},
			expected: []godo.TaskCreateSSHKey{{Fingerprint: "fingerprint"}},
		},
		{
			in:       []string{"1", "2"},
			expected: []godo.TaskCreateSSHKey{{ID: 1}, {ID: 2}},
		},
		{
			in:       []string{"1", "fingerprint"},
			expected: []godo.TaskCreateSSHKey{{ID: 1}, {Fingerprint: "fingerprint"}},
		},
	}

	for _, c := range cases {
		got := extractSSHKeys(c.in)
		assert.Equal(t, c.expected, got)
	}
}
