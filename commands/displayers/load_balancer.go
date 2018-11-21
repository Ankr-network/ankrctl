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
	"strings"

	"github.com/Ankr-network/dccn-cli/do"
)

type LoadBalancer struct {
	LoadBalancers do.LoadBalancers
}

var _ Displayable = &LoadBalancer{}

func (lb *LoadBalancer) JSON(out io.Writer) error {
	return writeJSON(lb.LoadBalancers, out)
}

func (lb *LoadBalancer) Cols() []string {
	return []string{
		"ID",
		"IP",
		"Name",
		"Status",
		"Created",
		"Algorithm",
		"Region",
		"Tag",
		"TaskIDs",
		"RedirectHttpToHttps",
		"StickySessions",
		"HealthCheck",
		"ForwardingRules",
	}
}

func (lb *LoadBalancer) ColMap() map[string]string {
	return map[string]string{
		"ID":                  "ID",
		"IP":                  "IP",
		"Name":                "Name",
		"Status":              "Status",
		"Created":             "Created At",
		"Algorithm":           "Algorithm",
		"Region":              "Region",
		"Tag":                 "Tag",
		"TaskIDs":          "Task IDs",
		"RedirectHttpToHttps": "SSL",
		"StickySessions":      "Sticky Sessions",
		"HealthCheck":         "Health Check",
		"ForwardingRules":     "Forwarding Rules",
	}
}

func (lb *LoadBalancer) KV() []map[string]interface{} {
	out := []map[string]interface{}{}

	for _, l := range lb.LoadBalancers {
		forwardingRules := []string{}
		for _, r := range l.ForwardingRules {
			forwardingRules = append(forwardingRules, prettyPrintStruct(r))
		}

		o := map[string]interface{}{
			"ID":                  l.ID,
			"IP":                  l.IP,
			"Name":                l.Name,
			"Status":              l.Status,
			"Created":             l.Created,
			"Algorithm":           l.Algorithm,
			"Region":              l.Region.Slug,
			"Tag":                 l.Tag,
			"TaskIDs":          fmt.Sprintf(strings.Trim(strings.Replace(fmt.Sprint(l.TaskIDs), " ", ",", -1), "[]")),
			"RedirectHttpToHttps": l.RedirectHttpToHttps,
			"StickySessions":      prettyPrintStruct(l.StickySessions),
			"HealthCheck":         prettyPrintStruct(l.HealthCheck),
			"ForwardingRules":     fmt.Sprintf(strings.Join(forwardingRules, " ")),
		}
		out = append(out, o)
	}

	return out
}
