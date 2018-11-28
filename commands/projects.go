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

/*
import (
	"fmt"
	"strings"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/godo"
	"github.com/spf13/cobra"
)

func Projects() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "projects",
			Short: "[beta] projects commands",
			Long:  "[beta] projects commands are for creating and managing projects",
		},
	}

	CmdBuilder(cmd, RunProjectsList, "list", "list projects", Writer, aliasOpt("ls"), displayerType(&displayers.Project{}), betaCmd())
	CmdBuilder(cmd, RunProjectsGet, "get <id>", "get a project; use \"default\" as ID to get default project", Writer, aliasOpt("g"), displayerType(&displayers.Project{}), betaCmd())

	cmdProjectsCreate := CmdBuilder(cmd, RunProjectsCreate, "create", "create project", Writer, aliasOpt("c"), displayerType(&displayers.Project{}), betaCmd())
	AddStringFlag(cmdProjectsCreate, dccncli.ArgProjectName, "", "", "project name", requiredOpt())
	AddStringFlag(cmdProjectsCreate, dccncli.ArgProjectPurpose, "", "", "project purpose", requiredOpt())
	AddStringFlag(cmdProjectsCreate, dccncli.ArgProjectDescription, "", "", "a description of your project")
	AddStringFlag(cmdProjectsCreate, dccncli.ArgProjectEnvironment, "", "", "the environment in which your project resides. Should be one of 'Development', 'Staging', 'Production'.")

	cmdProjectsUpdate := CmdBuilder(cmd, RunProjectsUpdate, "update <id>", "update project; use \"default\" as ID to update the default project", Writer, aliasOpt("u"), displayerType(&displayers.Project{}), betaCmd())
	AddStringFlag(cmdProjectsUpdate, dccncli.ArgProjectName, "", "", "project name")
	AddStringFlag(cmdProjectsUpdate, dccncli.ArgProjectPurpose, "", "", "project purpose")
	AddStringFlag(cmdProjectsUpdate, dccncli.ArgProjectDescription, "", "", "a description of your project")
	AddStringFlag(cmdProjectsUpdate, dccncli.ArgProjectEnvironment, "", "", "the environment in which your project resides. Should be one of 'Development', 'Staging', 'Production'.")
	AddBoolFlag(cmdProjectsUpdate, dccncli.ArgProjectIsDefault, "", false, "set the specified project as your default project")

	cmdProjectsDelete := CmdBuilder(cmd, RunProjectsDelete, "delete <id> [<id> ...]", "delete project", Writer, aliasOpt("d", "rm"), betaCmd())
	AddBoolFlag(cmdProjectsDelete, dccncli.ArgForce, dccncli.ArgShortForce, false, "Force project delete")

	cmd.AddCommand(ProjectResourcesCmd())

	return cmd
}

func ProjectResourcesCmd() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "resources",
			Short: "project resources commands",
			Long:  "project resources commands are for assigning and listing resources in projects",
		},
	}
	CmdBuilder(cmd, RunProjectResourcesList, "list <project-id>", "list project resources", Writer, aliasOpt("ls"), displayerType(&displayers.ProjectResource{}), betaCmd())
	CmdBuilder(cmd, RunProjectResourcesGet, "get <urn>", "get a project resource by its URN", Writer, aliasOpt("g"), betaCmd())

	cmdProjectResourcesAssign := CmdBuilder(cmd, RunProjectResourcesAssign, "assign <project-id> --resource=<urn> [--resource=<urn> ...]", "assign one or more resources to a project project", Writer, aliasOpt("a"), betaCmd())
	AddStringSliceFlag(cmdProjectResourcesAssign, dccncli.ArgProjectResource, "", []string{}, "resource URNs denoting resources to assign to the project")

	return cmd
}

// RunProjectsList lists Projects.
func RunProjectsList(c *CmdConfig) error {
	ps := c.Projects()
	list, err := ps.List()
	if err != nil {
		return err
	}

	return c.Display(&displayers.Project{Projects: list})
}

// RunProjectsGet retrieves an existing Project by its identifier. Use "default"
// as an identifier to retrieve your default project.
func RunProjectsGet(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	id := c.Args[0]

	ps := c.Projects()
	p, err := ps.Get(id)
	if err != nil {
		return err
	}

	return c.Display(&displayers.Project{Projects: do.Projects{*p}})
}

// RunProjectsCreate creates a new Project with a given configuration.
func RunProjectsCreate(c *CmdConfig) error {
	r := new(godo.CreateProjectRequest)
	if err := buildProjectsCreateRequestFromArgs(c, r); err != nil {
		return err
	}

	ps := c.Projects()
	p, err := ps.Create(r)
	if err != nil {
		return err
	}

	return c.Display(&displayers.Project{Projects: do.Projects{*p}})
}

func RunProjectsUpdate(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	id := c.Args[0]

	r := new(godo.UpdateProjectRequest)
	if err := buildProjectsUpdateRequestFromArgs(c, r); err != nil {
		return err
	}

	ps := c.Projects()
	p, err := ps.Update(id, r)
	if err != nil {
		return err
	}

	return c.Display(&displayers.Project{Projects: do.Projects{*p}})
}

func RunProjectsDelete(c *CmdConfig) error {
	if len(c.Args) < 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}

	force, err := c.Ankr.GetBool(c.NS, dccncli.ArgForce)
	if err != nil {
		return err
	}

	ps := c.Projects()
	var suffix string
	if len(c.Args) != 1 {
		suffix = "s"
	}
	if force || AskForConfirm(fmt.Sprintf("delete %d project%s", len(c.Args), suffix)) == nil {
		for _, id := range c.Args {
			if err := ps.Delete(id); err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("operation aborted")
}

func RunProjectResourcesList(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	id := c.Args[0]

	ps := c.Projects()
	list, err := ps.ListResources(id)
	if err != nil {
		return err
	}

	return c.Display(&displayers.ProjectResource{ProjectResources: list})
}

func RunProjectResourcesGet(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	urn := c.Args[0]

	parts, isValid := validateURN(urn)
	if !isValid {
		return dccncli.NewInvalidURNErr(urn)
	}

	c.Args = []string{parts[2]}
	switch parts[1] {
	case "task":
		return RunTaskGet(c)
	case "floatingip":
		return RunFloatingIPGet(c)
	case "loadbalancer":
		return RunLoadBalancerGet(c)
	case "domain":
		return RunDomainGet(c)
	case "volume":
		return RunVolumeGet(c)
	default:
		return fmt.Errorf("%q is an invalid resource type, consult the documentation", parts[1])
	}
}

func RunProjectResourcesAssign(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	projectUUID := c.Args[0]

	urns, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgProjectResource)
	if err != nil {
		return err
	}

	ps := c.Projects()
	list, err := ps.AssignResources(projectUUID, urns)
	if err != nil {
		return err
	}

	return c.Display(&displayers.ProjectResource{ProjectResources: list})
}

func validateURN(urn string) ([]string, bool) {
	parts := strings.Split(urn, ":")
	if len(parts) != 3 {
		return nil, false
	}

	if parts[0] != "do" {
		return nil, false
	}

	if strings.TrimSpace(parts[1]) == "" {
		return nil, false
	}

	if strings.TrimSpace(parts[2]) == "" {
		return nil, false
	}

	return parts, true
}

func buildProjectsCreateRequestFromArgs(c *CmdConfig, r *godo.CreateProjectRequest) error {
	name, err := c.Ankr.GetString(c.NS, dccncli.ArgProjectName)
	if err != nil {
		return err
	}
	r.Name = name

	purpose, err := c.Ankr.GetString(c.NS, dccncli.ArgProjectPurpose)
	if err != nil {
		return err
	}
	r.Purpose = purpose

	description, err := c.Ankr.GetString(c.NS, dccncli.ArgProjectDescription)
	if err != nil {
		return err
	}
	r.Description = description

	environment, err := c.Ankr.GetString(c.NS, dccncli.ArgProjectEnvironment)
	if err != nil {
		return err
	}
	r.Environment = environment

	return nil
}

func buildProjectsUpdateRequestFromArgs(c *CmdConfig, r *godo.UpdateProjectRequest) error {
	if c.Ankr.IsSet(dccncli.ArgProjectName) {
		name, err := c.Ankr.GetString(c.NS, dccncli.ArgProjectName)
		if err != nil {
			return err
		}
		r.Name = name
	}

	if c.Ankr.IsSet(dccncli.ArgProjectPurpose) {
		purpose, err := c.Ankr.GetString(c.NS, dccncli.ArgProjectPurpose)
		if err != nil {
			return err
		}
		r.Purpose = purpose
	}

	if c.Ankr.IsSet(dccncli.ArgProjectDescription) {
		description, err := c.Ankr.GetString(c.NS, dccncli.ArgProjectDescription)
		if err != nil {
			return err
		}
		r.Description = description
	}

	if c.Ankr.IsSet(dccncli.ArgProjectEnvironment) {
		environment, err := c.Ankr.GetString(c.NS, dccncli.ArgProjectEnvironment)
		if err != nil {
			return err
		}
		r.Environment = environment
	}

	if c.Ankr.IsSet(dccncli.ArgProjectIsDefault) {
		isDefault, err := c.Ankr.GetBool(c.NS, dccncli.ArgProjectIsDefault)
		if err != nil {
			return err
		}
		r.IsDefault = isDefault
	}

	return nil
}
*/
