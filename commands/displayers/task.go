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

type Task struct {
	//Tasks pb.Tasks
	Tasks []*pb.Task
}

var _ Displayable = &Task{}

func (d *Task) JSON(out io.Writer) error {
	return writeJSON(d.Tasks, out)
}

func (d *Task) Cols() []string {
	cols := []string{
		"TaskId", "TaskName", "Type", "Image", "Uptime", "CreationDate", "Replica", "DataCenter", "Status",
	}
	return cols
}

func (d *Task) ColMap() map[string]string {
	return map[string]string{
		"TaskId": "TaskId", "TaskName": "TaskName", "Type": "Type", "Image": "Image", "Uptime": "Uptime",
		"CreationDate": "CreationDate", "Replica": "Replica", "DataCenter": "DataCenter", "Status": "Status",
	}
}

func (d *Task) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, d := range d.Tasks {
		m := map[string]interface{}{
			"TaskId": d.Id, "TaskName": d.Name, "Type": d.Type, "Image": d.Image, "Uptime": d.Uptime,
			"CreationDate": d.CreationDate, "Replica": d.Replica, "DataCenter": d.DataCenter, "Status": d.Status,
		}
		out = append(out, m)
	}

	return out
}
