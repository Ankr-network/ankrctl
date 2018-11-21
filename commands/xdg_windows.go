// Copyright 2018 The Dccncli Authors All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build windows

package commands

import (
	"os"
	"path/filepath"
)

func configHome() string {
	configHome := os.Getenv("LOCALAPPDATA")
	if configHome == "" {
		// Resort to APPDATA for Windows XP users.
		configHome = os.Getenv("APPDATA")
		if configHome == "" {
			// If still empty, use the default path
			userName := os.Getenv("USERNAME")
			configHome = filepath.Join("C:/", "Users", userName, "AppData", "Local")
		}
	}

	return filepath.Join(configHome, "dccncli", "config")
}

// legacyConfigCheck is a no-op on windows since go doesn't have a chmod
// on this platform.
func legacyConfigCheck() {
}
