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
	"strings"
	"time"

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
		"TaskId", "TaskName", "Type", "Image", "LastModifyDate", "CreationDate", "Replica", "DataCenterName", "Status",
	}
	return cols
}

func (d *Task) ColMap() map[string]string {
	return map[string]string{
		"TaskId": "Task Id", "TaskName": "Task Name", "Type": "Type", "Image": "Image", "LastModifyDate": "Last Modify Date",
		"CreationDate": "Creation Date", "Replica": "Replica", "DataCenterName": "Data Center", "Status": "Status",
	}
}

func (d *Task) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, d := range d.Tasks {
		image := ""
		switch d.Type {
		case pb.TaskType_CRONJOB:
			image = d.GetTypeCronJob().Image
		case pb.TaskType_DEPLOYMENT:
			image = d.GetTypeDeployment().Image
		case pb.TaskType_JOB:
			image = d.GetTypeJob().Image
		}
		m := map[string]interface{}{
			"TaskId": d.Id, "TaskName": d.Name, "Type": strings.ToLower(d.Type.String()), "Image": image,
			"LastModifyDate": time.Unix(int64(d.Attributes.LastModifiedDate), 0).Format(time.RFC822),
			"CreationDate":   time.Unix(int64(d.Attributes.CreationDate), 0).Format(time.RFC822),
			"Replica":        d.Attributes.Replica, "DataCenterName": d.DataCenterName,
			"Status": strings.ToLower(d.Status.String()),
		}
		out = append(out, m)
	}

	return out
}
