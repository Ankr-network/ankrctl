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
	//"fmt"
	"io"
	//"strings"
	pb "github.com/Ankr-network/dccn-rpc/protocol_new/cli"
)

type Task struct {
	//Tasks pb.Tasks
	Tasks []pb.TaskInfo
}

var _ Displayable = &Task{}

func (d *Task) JSON(out io.Writer) error {
	return writeJSON(d.Tasks, out)
}

func (d *Task) Cols() []string {
	cols := []string{
		"Taskid", "Taskname", "Uptime", "Creationdate", "Replica", "Datacenter", "Status",
	}
	return cols
}

func (d *Task) ColMap() map[string]string {
	return map[string]string{
		"Taskid": "Taskid", "Taskname": "Taskname", "Uptime": "Uptime",
		"Creationdate": "Creationdate", "Replica": "Replica",
		"Datacenter": "Datacenter", "Status": "Status",
	}
}

func (d *Task) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, d := range d.Tasks {
		m := map[string]interface{}{
			"Taskid": d.Taskid, "Taskname": d.Taskname, "Uptime": d.Uptime,
			"Creationdate": d.Creationdate, "Replica": d.Replica,
			"Datacenter": d.Datacenter, "Status": d.Status,
		}
		out = append(out, m)
	}

	return out
}
