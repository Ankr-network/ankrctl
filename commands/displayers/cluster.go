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
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	common "github.com/Ankr-network/dccn-common/protos/common"
)

type Cluster struct {
	Clusters []common.DataCenterStatus
}

type Metrics struct {
	TotalCPU      int64
	UsedCPU       int64
	TotalMemory   int64
	UsedMemory    int64
	TotalStorage  int64
	UsedStorage   int64
	ImageCount    int64
	EndPointCount int64
	NetworkIO     int64
}

var _ Displayable = &Cluster{}

func (c *Cluster) JSON(out io.Writer) error {
	return writeJSON(c.Clusters, out)
}

func (c *Cluster) Cols() []string {
	cols := []string{
		"ID", "Name", "CPU", "MEM", "Storage", "Lat", "Lng", "Status", "WalletAddress",
	}
	return cols
}

func (c *Cluster) ColMap() map[string]string {
	return map[string]string{
		"ID": "ID", "Name": "Name", "CPU": "CPU", "MEM": "Memory", "Storage": "Storage",
		"Lat": "Latitude", "Lng": "Longitude", "Status": "Status", "WalletAddress": "WalletAddress",
	}
}

func (c *Cluster) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, c := range c.Clusters {
		metrics := Metrics{}
		_ = json.Unmarshal([]byte(c.DcHeartbeatReport.Metrics), &metrics)
		m := map[string]interface{}{
			"ID": c.DcId, "Name": c.DcName, "CPU": strconv.Itoa(int(metrics.TotalCPU)) + "CPU(s)",
			"MEM":     fmt.Sprintf("%.2f", float64(metrics.TotalMemory)/1073741824) + "GB",
			"Storage": fmt.Sprintf("%.2f", float64(metrics.TotalStorage)/1073741824) + "GB",
			"Lat":     c.GeoLocation.Lat, "Lng": c.GeoLocation.Lng, "Status": strings.ToLower(c.DcStatus.String()),
			"WalletAddress": c.DcAttributes.WalletAddress,
		}
		out = append(out, m)
	}

	return out
}
