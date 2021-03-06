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
	"github.com/Ankr-network/ankrctl/types"
	"strings"
)

func handleColumns(ns string, config types.Config) ([]string, error) {
	colStr, err := config.GetString(ns, types.ArgFormat)
	if err != nil {
		return nil, err
	}

	var cols []string
	for _, c := range strings.Split(strings.Join(strings.Fields(colStr), ""), ",") {
		if c != "" {
			cols = append(cols, c)
		}
	}

	hh, err := config.GetBool(ns, types.ArgNoHeader)
	if err != nil {
		return nil, err
	}

	hc.HideHeader(hh)

	return cols, nil
}
