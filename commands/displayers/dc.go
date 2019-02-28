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

type Dc struct {
	Dcs []common.DataCenter
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

var _ Displayable = &Dc{}

func (d *Dc) JSON(out io.Writer) error {
	return writeJSON(d.Dcs, out)
}

func (d *Dc) Cols() []string {
	cols := []string{
		"Id", "Name", "CPU", "RAM", "HDD", "Lat", "Lng", "Status", "WalletAddress",
	}
	return cols
}

func (d *Dc) ColMap() map[string]string {
	return map[string]string{
		"Id": "Id", "Name": "Name", "CPU": "CPU", "RAM": "RAM", "HDD": "HDD",
		"Lat": "Latitude", "Lng": "Longitude", "Status": "Status", "WalletAddress": "WalletAddress",
	}
}

func (d *Dc) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, d := range d.Dcs {
		metrics := Metrics{}
		_ = json.Unmarshal([]byte(d.DcHeartbeatReport.Metrics), &metrics)
		m := map[string]interface{}{
			"Id": d.Id, "Name": d.Name, "CPU": strconv.Itoa(int(metrics.TotalCPU)) + "CPU(s)",
			"RAM": fmt.Sprintf("%.2f", float64(metrics.TotalMemory)/1073741824) + "GB",
			"HDD": fmt.Sprintf("%.2f", float64(metrics.TotalStorage)/1073741824) + "GB",
			"Lat": d.GeoLocation.Lat, "Lng": d.GeoLocation.Lng, "Status": strings.ToLower(d.Status.String()),
			"WalletAddress": d.DcAttributes.WalletAddress,
		}
		out = append(out, m)
	}

	return out
}
