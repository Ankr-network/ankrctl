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
	"strconv"

	"github.com/Ankr-network/dccn-cli"
	"github.com/spf13/cobra"
)

// Version creates a version command.
func Version() *Command {
	return &Command{
		Command: &cobra.Command{
			Use:   "version",
			Short: "show the current version",
			Run: func(cmd *cobra.Command, args []string) {
				if dccncli.Build != "" {
					dccncli.AnkrVersion.Build = dccncli.Build
				}
				if dccncli.Major != "" {
					i, _ := strconv.Atoi(dccncli.Major)
					dccncli.AnkrVersion.Major = i
				}
				if dccncli.Minor != "" {
					i, _ := strconv.Atoi(dccncli.Minor)
					dccncli.AnkrVersion.Minor = i
				}
				if dccncli.Patch != "" {
					i, _ := strconv.Atoi(dccncli.Patch)
					dccncli.AnkrVersion.Patch = i
				}
				if dccncli.Label != "" {
					dccncli.AnkrVersion.Label = dccncli.Label
				}

				fmt.Println(dccncli.AnkrVersion.Complete(&dccncli.GithubLatestVersioner{}))
			},
		},
	}
}
