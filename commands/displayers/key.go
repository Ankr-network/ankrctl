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
)

type Key struct {
	Keystores []*KeyStore
}

type KeyStore struct {
	Name      string `json:"name,omitempty"`
	Address   string `json:"address"`
	PublicKey string `json:"publickey"`
}

var _ Displayable = &Key{}

func (c *Key) JSON(out io.Writer) error {
	return writeJSON(c.Keystores, out)
}

func (c *Key) Cols() []string {
	cols := []string{
		"Name", "Address", "PublicKey",
	}
	return cols
}

func (c *Key) ColMap() map[string]string {
	return map[string]string{
		"Name": "Name", "Address": "Address",
		"PublicKey": "Public Key",
	}
}

func (c *Key) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, c := range c.Keystores {
		m := map[string]interface{}{
			"Name": c.Name, "Address": c.Address,
			"PublicKey": c.PublicKey,
		}
		out = append(out, m)
	}

	return out
}
