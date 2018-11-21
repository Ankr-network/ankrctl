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

package main

import (
	"flag"
	"log"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/install"
)

var (
	ver    = flag.String("ver", dccncli.AnkrVersion.String(), "ankr version")
	path   = flag.String("path", "", "upload path")
	user   = flag.String("user", "", "bintray user")
	apikey = flag.String("apikey", "", "bintray apikey")
)

func main() {
	flag.Parse()

	if *path == "" {
		log.Fatal("path is required")
	}

	ui := install.UserInfo{
		User:   *user,
		Apikey: *apikey,
	}

	err := install.Upload(ui, *ver, *path)
	if err != nil {
		log.Fatalf("upload failed: %v", err)
	}
}
