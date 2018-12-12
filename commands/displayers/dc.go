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
	pb "github.com/Ankr-network/dccn-rpc/protocol"
)

type Dc struct {
	//Tasks pb.Tasks
	Dcs []pb.DataCenterInfo
}

var _ Displayable = &Dc{}

func (d *Dc) JSON(out io.Writer) error {
	return writeJSON(d.Dcs, out)
}

func (d *Dc) Cols() []string {
	cols := []string{
		"Id", "Name",
	}
	return cols
}

func (d *Dc) ColMap() map[string]string {
	return map[string]string{
		"Id": "Id", "Name": "Name",
	}
}

func (d *Dc) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, d := range d.Dcs {
		m := map[string]interface{}{
			"Id": d.Id, "Name": d.Name,
		}
		out = append(out, m)
	}

	return out
}
