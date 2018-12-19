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

package akrctl

const (
	// ArgAccessToken is the access token to be used for the operations
	ArgAccessToken = "access-token"
	// ArgContext is the name of the auth context to use
	ArgContext = "context"
	// ArgTaskID is a task id argument.
	ArgTaskID = "task-id"
	// ArgTaskIDs is a list of task IDs.
	ArgTaskIDs = "task-ids"
	// ArgTaskName is a task name argument.
	ArgTaskName = "task-name"
	// ArgRegionSlug is a region slug argument.
	ArgRegionSlug = "region"
	// ArgZoneSlug is a zone slug argument.
	ArgZoneSlug = "zone"
	// ArgFormat is columns to include in output argment.
	ArgFormat = "format"
	// ArgNoHeader hides the output header.
	ArgNoHeader = "no-header"
	// ArgPollTime is how long before the next poll argument.
	ArgPollTime = "poll-timeout"
	// ArgOutput is an output type argument.
	ArgOutput = "output"
	// ArgForce forces confirmation on actions
	ArgForce = "force"
)
