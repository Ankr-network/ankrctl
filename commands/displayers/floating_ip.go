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

	"github.com/Ankr-network/dccn-cli/do"
)

type FloatingIP struct {
	FloatingIPs do.FloatingIPs
}

var _ Displayable = &FloatingIP{}

func (fi *FloatingIP) JSON(out io.Writer) error {
	return writeJSON(fi.FloatingIPs, out)
}

func (fi *FloatingIP) Cols() []string {
	return []string{
		"IP", "Region", "TaskID", "TaskName",
	}
}

func (fi *FloatingIP) ColMap() map[string]string {
	return map[string]string{
		"IP": "IP", "Region": "Region", "TaskID": "Task ID", "TaskName": "Task Name",
	}
}

func (fi *FloatingIP) KV() []map[string]interface{} {
	out := []map[string]interface{}{}

	for _, f := range fi.FloatingIPs {
		var taskID, taskName string
		if f.Task != nil {
			taskID = fmt.Sprintf("%d", f.Task.ID)
			taskName = f.Task.Name
		}

		o := map[string]interface{}{
			"IP": f.IP, "Region": f.Region.Slug,
			"TaskID": taskID, "TaskName": taskName,
		}

		out = append(out, o)
	}

	return out
}
