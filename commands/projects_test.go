package commands

import (
	"testing"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/godo"
	"github.com/stretchr/testify/assert"
)

var (
	testProject = do.Project{
		Project: &godo.Project{
			Name:        "my project",
			Description: "my project description",
			Purpose:     "my project purpose",
			Environment: "Development",
			IsDefault:   false,
		},
	}

	testProjectList = do.Projects{testProject}

	testProjectResourcesList = do.ProjectResources{
		{
			ProjectResource: &godo.ProjectResource{URN: "do:task:1234"},
		},
		{
			ProjectResource: &godo.ProjectResource{URN: "do:floatingip:1.2.3.4"},
		},
	}
	testProjectResourcesListSingle = do.ProjectResources{
		{
			ProjectResource: &godo.ProjectResource{URN: "do:task:1234"},
		},
	}
)

func TestProjectsCommand(t *testing.T) {
	cmd := Projects()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "list", "get", "create", "update", "delete", "resources")
}

func TestProjectsList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.projects.On("List").Return(testProjectList, nil)

		err := RunProjectsList(config)
		assert.NoError(t, err)
	})
}

func TestProjectsGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		projectUUID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tm.projects.On("Get", projectUUID).Return(&testProject, nil)

		config.Args = append(config.Args, projectUUID)

		err := RunProjectsGet(config)
		assert.NoError(t, err)
	})
}

func TestProjectsCreate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		projectCreateRequest := &godo.CreateProjectRequest{
			Name:        "project name",
			Description: "project description",
			Purpose:     "personal use",
			Environment: "Staging",
		}
		tm.projects.On("Create", projectCreateRequest).Return(&testProject, nil)

		config.Ankr.Set(config.NS, dccncli.ArgProjectName, "project name")
		config.Ankr.Set(config.NS, dccncli.ArgProjectDescription, "project description")
		config.Ankr.Set(config.NS, dccncli.ArgProjectPurpose, "personal use")
		config.Ankr.Set(config.NS, dccncli.ArgProjectEnvironment, "Staging")

		err := RunProjectsCreate(config)
		assert.NoError(t, err)
	})
}

func TestProjectsUpdateAllAttributes(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		projectUUID := "ab06e011-6dd1-4034-9293-201f71aba299"
		updateReq := &godo.UpdateProjectRequest{
			Name:        "project name",
			Description: "project description",
			Purpose:     "project purpose",
			Environment: "Production",
			IsDefault:   false,
		}
		tm.projects.On("Update", projectUUID, updateReq).Return(&testProject, nil)

		config.Ankr.(*TestConfig).IsSetMap = map[string]bool{
			dccncli.ArgProjectName:        true,
			dccncli.ArgProjectDescription: true,
			dccncli.ArgProjectPurpose:     true,
			dccncli.ArgProjectEnvironment: true,
			dccncli.ArgProjectIsDefault:   true,
		}

		config.Args = append(config.Args, projectUUID)
		config.Ankr.Set(config.NS, dccncli.ArgProjectName, "project name")
		config.Ankr.Set(config.NS, dccncli.ArgProjectDescription, "project description")
		config.Ankr.Set(config.NS, dccncli.ArgProjectPurpose, "project purpose")
		config.Ankr.Set(config.NS, dccncli.ArgProjectEnvironment, "Production")
		config.Ankr.Set(config.NS, dccncli.ArgProjectIsDefault, false)

		err := RunProjectsUpdate(config)
		assert.NoError(t, err)
	})
}

func TestProjectsUpdateSomeAttributes(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		projectUUID := "ab06e011-6dd1-4034-9293-201f71aba299"
		updateReq := &godo.UpdateProjectRequest{
			Name:        "project name",
			Description: "project description",
			Purpose:     nil,
			Environment: nil,
			IsDefault:   nil,
		}
		tm.projects.On("Update", projectUUID, updateReq).Return(&testProject, nil)

		config.Ankr.(*TestConfig).IsSetMap = map[string]bool{
			dccncli.ArgProjectName:        true,
			dccncli.ArgProjectDescription: true,
			dccncli.ArgProjectPurpose:     false,
			dccncli.ArgProjectEnvironment: false,
			dccncli.ArgProjectIsDefault:   false,
		}

		config.Args = append(config.Args, projectUUID)
		config.Ankr.Set(config.NS, dccncli.ArgProjectName, "project name")
		config.Ankr.Set(config.NS, dccncli.ArgProjectDescription, "project description")

		err := RunProjectsUpdate(config)
		assert.NoError(t, err)
	})
}

func TestProjectsUpdateOneAttribute(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		projectUUID := "ab06e011-6dd1-4034-9293-201f71aba299"
		updateReq := &godo.UpdateProjectRequest{
			Name:        "project name",
			Description: nil,
			Purpose:     nil,
			Environment: nil,
			IsDefault:   nil,
		}
		tm.projects.On("Update", projectUUID, updateReq).Return(&testProject, nil)

		config.Ankr.(*TestConfig).IsSetMap = map[string]bool{
			dccncli.ArgProjectName:        true,
			dccncli.ArgProjectDescription: false,
			dccncli.ArgProjectPurpose:     false,
			dccncli.ArgProjectEnvironment: false,
			dccncli.ArgProjectIsDefault:   false,
		}

		config.Args = append(config.Args, projectUUID)
		config.Ankr.Set(config.NS, dccncli.ArgProjectName, "project name")
		config.Ankr.Set(config.NS, dccncli.ArgProjectDescription, "project description")

		err := RunProjectsUpdate(config)
		assert.NoError(t, err)
	})
}

func TestProjectsDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		projectUUID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tm.projects.On("Delete", projectUUID).Return(nil)

		config.Args = append(config.Args, projectUUID)
		config.Ankr.Set(config.NS, dccncli.ArgForce, true)

		err := RunProjectsDelete(config)
		assert.NoError(t, err)
	})
}

func TestProjectResourcesCommand(t *testing.T) {
	cmd := ProjectResourcesCmd()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "list", "get", "assign")
}

func TestProjectResourcesList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		projectUUID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tm.projects.On("ListResources", projectUUID).Return(testProjectResourcesList, nil)

		config.Args = append(config.Args, projectUUID)
		err := RunProjectResourcesList(config)
		assert.NoError(t, err)
	})
}

func TestProjectResourcesGetWithValidURN(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tasks.On("Get", 1234).Return(&testTask, nil)

		config.Args = append(config.Args, "do:task:1234")
		err := RunProjectResourcesGet(config)
		assert.NoError(t, err)
	})
}

func TestProjectResourcesGetWithInvalidURN(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		config.Args = append(config.Args, "fakeurn")
		err := RunProjectResourcesGet(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "URN must be in the format")
	})
}

func TestProjectResourcesAssignOneResource(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		projectUUID := "ab06e011-6dd1-4034-9293-201f71aba299"
		urn := "do:task:1234"
		tm.projects.On("AssignResources", projectUUID, []string{urn}).Return(testProjectResourcesListSingle, nil)

		config.Args = append(config.Args, projectUUID)
		config.Ankr.Set(config.NS, dccncli.ArgProjectResource, []string{urn})

		err := RunProjectResourcesAssign(config)
		assert.NoError(t, err)
	})
}

func TestProjectResourcesAssignMultipleResources(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		projectUUID := "ab06e011-6dd1-4034-9293-201f71aba299"
		urn := "do:task:1234"
		otherURN := "do:floatingip:1.2.3.4"
		tm.projects.On("AssignResources", projectUUID, []string{urn, otherURN}).Return(testProjectResourcesList, nil)

		config.Args = append(config.Args, projectUUID)
		config.Ankr.Set(config.NS, dccncli.ArgProjectResource, []string{urn, otherURN})

		err := RunProjectResourcesAssign(config)
		assert.NoError(t, err)
	})
}
