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
	"io"

	pb "github.com/Ankr-network/dccn-common/protos/common"
)

type Chart struct {
	Charts []*pb.Chart
}

var _ Displayable = &Chart{}

func (c *Chart) JSON(out io.Writer) error {
	return writeJSON(c.Charts, out)
}

func (c *Chart) Cols() []string {
	cols := []string{
		"Repo", "Name", "LatestVersion", "LatestAppVersion", "Description",
	}
	return cols
}

func (c *Chart) ColMap() map[string]string {
	return map[string]string{
		"Repo": "Repo", "Name": "Name", "LatestVersion": "Latest Version",
		"LatestAppVersion": "Latest App Version", "Description": "Description",
	}
}

func (c *Chart) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, c := range c.Charts {
		m := map[string]interface{}{
			"Repo": c.ChartRepo, "Name": c.ChartName, "LatestVersion": c.ChartLatestVersion,
			"LatestAppVersion": c.ChartLatestAppVersion, "Description": c.ChartDescription,
		}
		out = append(out, m)
	}

	return out
}
