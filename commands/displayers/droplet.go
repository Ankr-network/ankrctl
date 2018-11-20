/*
Copyright 2018 The Doctl Authors All rights reserved.
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

	"github.com/Ankr-network/dccn-cli/do"
)

type Droplet struct {
	Droplets do.Droplets
}

var _ Displayable = &Droplet{}

func (d *Droplet) JSON(out io.Writer) error {
	return writeJSON(d.Droplets, out)
}

func (d *Droplet) Cols() []string {
	cols := []string{
		"ID", "Taskname", "Uptime", "Creationdate", "Status",
	}
	return cols
}

func (d *Droplet) ColMap() map[string]string {
	return map[string]string{
		"ID":"ID", "Taskname": "Taskname", "Uptime": "Uptime", 
		"Creationdate": "Creationdate", "Status": "Status",
	}
}

func (d *Droplet) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, d := range d.Droplets {
		m := map[string]interface{}{
			"ID": d.ID,"Taskname": d.Taskname, "Uptime": d.Uptime, 
			"Creationdate": d.Creationdate, "Status": d.Status,
		}
		out = append(out, m)
	}

	return out
}
