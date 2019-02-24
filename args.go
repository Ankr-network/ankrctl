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
	// ArgUserID is a user id argument.
	ArgUserID = "userid"
	// ArgDcNameSlug is a datacenter slug argument.
	ArgDcNameSlug = "dc-name"
	// ArgTypeSlug is a type slug argument.
	ArgTypeSlug = "type"
	// ArgFormat is columns to include in output argment.
	ArgFormat = "format"
	// ArgNoHeader hides the output header.
	ArgNoHeader = "no-header"
	// ArgOutput is an output type argument.
	ArgOutput = "output"
	// ArgForce forces confirmation on actions
	ArgForce = "force"
	// ArgTaskIdSlug is a task id slug argument.
	ArgTaskIdSlug = "task-id"
	// ArgImageSlug is a task image slug argument.
	ArgImageSlug = "image"
	// ArgReplicaSlug is a rtask eplica slug argument.
	ArgReplicaSlug = "replica"
	// ArgScheduleSlug is a task schedule slug argument.
	ArgScheduleSlug = "schedule"
	// ArgPasswordSlug is a user password slug argument.
	ArgPasswordSlug = "password"
	// ArgEmailSlug is a user email slug argument.
	ArgEmailSlug = "email"
	// ArgAddressSlug is a wallet address slug argument.
	ArgAddressSlug = "address"
	// ArgTargetSlug is a wallet send token target address slug argument.
	ArgTargetSlug = "target"
	// ArgPrivateKeySlug is a wallet private key slug argument.
	ArgPrivateKeySlug = "private-key"
	// ArgPublicKeySlug is a wallet public key slug argument.
	ArgPublicKeySlug = "public-key"
	// ArgRegisterCodeSlug is a user registration confirmation code slug argument.
	ArgRegisterCodeSlug = "register-code"
	// ArgPasswordCodeSlug is a password registration confirmation code slug argument.
	ArgPasswordCodeSlug = "password-code"
	// ArgConfirmPasswordSlug is a password reset confirmation new password slug argument.
	ArgConfirmPasswordSlug = "confirm-password"
	// ArgNewPasswordSlug is a change new password slug argument.
	ArgNewPasswordSlug = "new-password"
	// ArgOldPasswordSlug is a change old password slug argument.
	ArgOldPasswordSlug = "old-password"
	// ArgUpdateKeySlug is a update user attribute key slug argument.
	ArgUpdateKeySlug = "update-key"
	// ArgUpdateValueSlug is a update user attribute value slug argument.
	ArgUpdateValueSlug = "update-value"
)
