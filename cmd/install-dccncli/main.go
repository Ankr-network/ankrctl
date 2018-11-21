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
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Ankr-network/dccn-cli/install"
	"github.com/fatih/color"
)

var (
	ver = "0.6.0"
)

func main() {

	var err error
	defer func() {
		if err != nil {
			log.Fatalf("error encountered: %v", err)
		}
	}()

	bold := color.New(color.Bold, color.FgWhite).SprintfFunc()

	// get install directory
	home, err := homeDir()
	if err != nil {
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("ankr installation directory (this will create a ankr subdirectory) (%s): ", bold(home))
	installDir, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	installDir = strings.TrimSpace(installDir)

	if installDir == "" {
		installDir = home
	}

	// create install directory
	fmt.Printf("creating %s/ankr directory...\n\n", installDir)
	err = os.MkdirAll(filepath.Join(installDir, "bin"), 0755)
	if err != nil {
		return
	}

	// create temp directory
	tmpDir, err := ioutil.TempDir("", "ankr-install-")
	if err != nil {
		return
	}
	defer func() {
		err := os.Remove(tmpDir)
		if err != nil {
			fmt.Printf("could not remove temp directory (%s): %v", tmpDir, err)
		}
	}()

	// retrieve ankr binary
	filename := archiveName(ver)

	fmt.Println("retrieving ankr...")
	ankrPath := filepath.Join(tmpDir, filename)
	file, err := install.Download(ankrPath, install.URL(filename))
	if err != nil {
		return
	}
	file.Close()
	fmt.Println()

	fmt.Println("retrieving ankr checksum...")
	checksumPath := filepath.Join(tmpDir, filename+".sha256")
	checksumFile, err := install.Download(checksumPath, install.URL(filename+".sha256"))
	if err != nil {
		log.Fatalf("could not download ankr checksum file: %v", err)
	}
	checksumFile.Close()
	fmt.Println(" ")

	// validate binary
	fmt.Println("validating ankr checksum...")
	f, err := os.Open(ankrPath)
	if err != nil {
		return
	}
	defer f.Close()

	cs, err := os.Open(checksumPath)
	if err != nil {
		return
	}
	defer func() {
		cs.Close()
		os.Remove(checksumPath)
	}()

	err = install.Validate(f, cs)
	if err != nil {
		return
	}

	fmt.Println("checksum was valid")

	// place binary in install directory
	ankrInstallPath := filepath.Join(installDir, "bin", "ankr")
	fmt.Println("placing ankr in install path...")
	err = os.Rename(ankrPath, ankrInstallPath)
	if err != nil {
		return
	}
	os.Chmod(ankrInstallPath, 0755)

	fmt.Println("install complete!")
}

func homeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

func archiveName(ver string) string {
	var suffix string

	if runtime.GOOS == "darwin" {
		suffix = "darwin-10.6-amd64"
	} else {
		suffix = fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	}

	return fmt.Sprintf("ankr-%s-%s", ver, suffix)
}
