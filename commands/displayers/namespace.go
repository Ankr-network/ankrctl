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

	common "github.com/Ankr-network/dccn-common/protos/common"
)

type Namespace struct {
	Namespaces []common.NamespaceReport
}

var _ Displayable = &Namespace{}

func (n *Namespace) JSON(out io.Writer) error {
	return writeJSON(n.Namespaces, out)
}

func (n *Namespace) Cols() []string {
	cols := []string{
		"ID", "Name", "CpuLimit", "MemLimit", "StorageLimit", "ClusterID", "ClusterName", "Status", "Event",
	}
	return cols
}

func (n *Namespace) ColMap() map[string]string {
	return map[string]string{
		"ID": "ID", "Name": "Name", "CpuLimit": "CPU Limit", "MemLimit": "Memory Limit",
		"StorageLimit": "Storage Limit", "ClusterID": "Cluster ID", "ClusterName": "Cluster Name",
		"Status": "Status", "Event": "Event",
	}
}

func (n *Namespace) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, n := range n.Namespaces {
		m := map[string]interface{}{
			"ID": n.Namespace.NsId, "Name": n.Namespace.NsName,
			"CpuLimit":     fmt.Sprintf("%v vCPU(s)", float64(n.Namespace.NsCpuLimit)/1000),
			"MemLimit":     fmt.Sprintf("%v GB", float64(n.Namespace.NsMemLimit)/1000),
			"StorageLimit": fmt.Sprintf("%v GB", float64(n.Namespace.NsStorageLimit)),
			"ClusterID":    n.Namespace.ClusterId, "ClusterName": n.Namespace.ClusterName,
			"Status": n.NsStatus.String(), "Event": n.NsEvent,
		}
		out = append(out, m)
	}

	return out
}
