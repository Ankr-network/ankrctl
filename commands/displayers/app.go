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

type AppReport struct {
	Apps []*pb.AppReport
}

var _ Displayable = &AppReport{}

func (d *AppReport) JSON(out io.Writer) error {
	return writeJSON(d.Apps, out)
}

func (d *AppReport) Cols() []string {
	cols := []string{
		"ID", "Name", "ChartRepo", "ChartName", "ChartVersion", "AppVersion", "Namespace",
		"Cluster", "LastModifyDate", "CreationDate", "Status", "Event", "Endpoint",
	}
	return cols
}

func (d *AppReport) ColMap() map[string]string {
	return map[string]string{
		"ID": "ID", "Name": "Name", "ChartRepo": "Chart Repo", "ChartName": "Chart Name",
		"ChartVersion": "Chart Version", "AppVersion": "App Version", "Namespace": "Namespace",
		"Cluster": "Cluster", "LastModifyDate": "Last Modified Date", "CreationDate": "Creation Date",
		"Status": "Status", "Event": "Event", "Endpoint": "Endpoint",
	}
}

func (d *AppReport) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, d := range d.Apps {
		m := map[string]interface{}{
			"ID": d.AppDeployment.AppId, "Name": d.AppDeployment.AppName,
			"ChartRepo": d.AppDeployment.ChartDetail.ChartRepo, "ChartName": d.AppDeployment.ChartDetail.ChartName,
			"ChartVersion": d.AppDeployment.ChartDetail.ChartVer, "AppVersion": d.AppDeployment.ChartDetail.ChartAppVer,
			"Namespace": d.AppDeployment.Namespace.NsName, "Cluster": d.AppDeployment.Namespace.ClusterName,
			"LastModifyDate": time.Unix(int64(d.AppDeployment.Attributes.LastModifiedDate.Seconds), 0).Format(time.RFC822),
			"CreationDate":   time.Unix(int64(d.AppDeployment.Attributes.CreationDate.Seconds), 0).Format(time.RFC822),
			"Status":         strings.ToLower(d.AppStatus.String()), "Event": strings.ToLower(d.AppEvent.String()),
			"Endpoint": d.Endpoint,
		}
		out = append(out, m)
	}

	return out
}
