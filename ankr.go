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

package dccncli

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

const (
	// NSRoot is a configuration key that signifies this value is at the root.
	NSRoot = "akrctl"
)

var (
	// AnkrConfig holds the app's current configuration.
	AnkrConfig Config = &LiveConfig{}

	// AnkrVersion is ankr's version.
	AnkrVersion = Version{
		Major: 1,
		Minor: 11,
		Patch: 0,
		Label: "dev",
	}

	// Build is ankr's build tag.
	Build string

	// Major is dccncli's major version.
	Major string

	// Minor is dccncli's minor version.
	Minor string

	// Patch is dccncli's patch version.
	Patch string

	// Label is dccncli's label.
	Label string
)

func init() {
	jww.SetStdoutThreshold(jww.LevelError)
}

// Version is the version info for ankr.
type Version struct {
	Major, Minor, Patch int
	Name, Build, Label  string
}

func (v Version) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))
	if v.Label != "" {
		buffer.WriteString("-" + v.Label)
	}

	return buffer.String()
}

// Config is an interface that represent ankr's config.
type Config interface {
	Set(ns, key string, val interface{})
	IsSet(key string) bool
	GetString(ns, key string) (string, error)
	GetBool(ns, key string) (bool, error)
	GetInt(ns, key string) (int, error)
	GetStringSlice(ns, key string) ([]string, error)
}

// LiveConfig is an implementation of Config for live values.
type LiveConfig struct {
	cliArgs map[string]bool
}

var _ Config = &LiveConfig{}

// Set sets a config key.
func (c *LiveConfig) Set(ns, key string, val interface{}) {
	nskey := fmt.Sprintf("%s-%s", ns, key)
	viper.Set(nskey, val)
}

func (c *LiveConfig) IsSet(key string) bool {
	r := regexp.MustCompile("\b*--([a-z-_]+)")
	matches := r.FindAllStringSubmatch(strings.Join(os.Args, " "), -1)
	if len(matches) == 0 {
		return false
	}

	if len(c.cliArgs) == 0 {
		args := make(map[string]bool)
		for _, match := range matches {
			args[match[1]] = true
		}
		c.cliArgs = args
	}

	return c.cliArgs[key]
}

// GetString returns a config value as a string.
func (c *LiveConfig) GetString(ns, key string) (string, error) {
	if ns == NSRoot {
		return viper.GetString(key), nil
	}

	nskey := fmt.Sprintf("%s.%s", ns, key)

	isRequired := viper.GetBool(fmt.Sprintf("required.%s", nskey))
	str := viper.GetString(nskey)

	if isRequired && strings.TrimSpace(str) == "" {
		return "", NewMissingArgsErr(nskey)
	}

	return str, nil
}

// GetBool returns a config value as a bool.
func (c *LiveConfig) GetBool(ns, key string) (bool, error) {
	if ns == NSRoot {
		return viper.GetBool(key), nil
	}

	nskey := fmt.Sprintf("%s.%s", ns, key)

	return viper.GetBool(nskey), nil
}

// GetInt returns a config value as an int.
func (c *LiveConfig) GetInt(ns, key string) (int, error) {
	if ns == NSRoot {
		return viper.GetInt(key), nil
	}

	nskey := fmt.Sprintf("%s.%s", ns, key)

	isRequired := viper.GetBool(fmt.Sprintf("required.%s", nskey))
	val := viper.GetInt(nskey)

	if isRequired && val == 0 {
		return 0, NewMissingArgsErr(nskey)
	}

	return val, nil
}

// GetStringSlice returns a config value as a string slice.
func (c *LiveConfig) GetStringSlice(ns, key string) ([]string, error) {
	if ns == NSRoot {
		return viper.GetStringSlice(key), nil
	}

	nskey := fmt.Sprintf("%s.%s", ns, key)

	isRequired := viper.GetBool(fmt.Sprintf("required.%s", nskey))
	val := viper.GetStringSlice(nskey)
	if isRequired && emptyStringSlice(val) {
		return nil, NewMissingArgsErr(nskey)
	}

	out := []string{}
	for _, item := range viper.GetStringSlice(nskey) {
		item = strings.TrimPrefix(item, "[")
		item = strings.TrimSuffix(item, "]")

		list := strings.Split(item, ",")
		for _, str := range list {
			if str == "" {
				continue
			}

			out = append(out, str)
		}
	}

	return out, nil
}

// This is needed because an empty StringSlice flag returns `["[]"]`
func emptyStringSlice(s []string) bool {
	return len(s) == 1 && s[0] == "[]"
}
