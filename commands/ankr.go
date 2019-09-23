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

package commands

import (
	"errors"
	"fmt"
	"github.com/Ankr-network/ankrctl/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ankr-network/ankrctl/commands/displayers"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	// defaultConfigName is the name of the config file when no alternative is supplied.
	defaultConfigName = "config.yaml"
)

// AnkrCmd is the base command.
var AnkrCmd = &Command{
	Command: &cobra.Command{
		Use:   "ankrctl",
		Short: "ankrctl is a command line interface for the Ankr DCCN Hub.",
	},
}

// Context holds the current auth context
var Context string

// HubURL holds the HUB URL to use.
var HubURL string

// Token holds the global authorization token.
var Token string

// Output holds the global output format.
var Output string

// Verbose toggles verbose output.
var Verbose bool

var requiredColor = color.New(color.Bold).SprintfFunc()

// Writer is where output should be written to.
var Writer = os.Stdout

// Trace toggles http tracing output.
var Trace bool

// cfgFile is the location of the config file
var cfgFile string

// cfgFileWriter is the config file writer
var cfgFileWriter = defaultConfigFileWriter

// ErrNoAccessToken is an error for when there is no access token.
var ErrNoAccessToken = errors.New("no access token has been configured")

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetEnvPrefix("ANKR")
	viper.BindEnv("hub-url", "ANKR_HUB_URL")
	viper.SetDefault("hub-url", clientURL)
	addCommands()
}

func initConfig() {
	var err error
	cfgFile, err = findConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	legacyConfigCheck()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	if _, err := os.Stat(cfgFile); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("reading initialization failed:", err)
		}
	}

	viper.SetDefault("output", "text")
	viper.SetDefault("context", "default")
}

func findConfig() (string, error) {
	if cfgFile != "" {
		return cfgFile, nil
	}

	if homeDir() != "" {
		legacyConfigPath := filepath.Join(homeDir(), ".ankrctlcfg")
		if _, err := os.Stat(legacyConfigPath); err == nil {
			msg := fmt.Sprintf("Configuration detected at %q. Please move .ankrctlcfg to %s",
				legacyConfigPath, configPath())
			warn(msg)
		}
	}

	if os.Getenv("XDG_CONFIG_HOME") != "" {
		legacyXDGPath := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "config.yaml")
		if _, err := os.Stat(legacyXDGPath); err == nil {
			msg := fmt.Sprintf("Configuration detected at %q. Please move config.yaml to %s",
				legacyXDGPath, configPath())
			warn(msg)
		}
	}

	ch := configHome()
	if err := os.MkdirAll(ch, 0755); err != nil {
		return "", err
	}

	return filepath.Join(ch, defaultConfigName), nil
}

func configPath() string {
	return fmt.Sprintf("%s/%s", configHome(), defaultConfigName)
}

// Execute executes the current command using AnkrCmd.
func Execute() {
	if err := AnkrCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// AddCommands adds sub commands to the base command.
func addCommands() {
	AnkrCmd.AddCommand(appCmd())
	AnkrCmd.AddCommand(namespaceCmd())
	AnkrCmd.AddCommand(clusterCmd())
	AnkrCmd.AddCommand(chartCmd())
	AnkrCmd.AddCommand(userCmd())
	AnkrCmd.AddCommand(walletCmd())
}

type flagOpt func(c *Command, name, key string)

func requiredOpt() flagOpt {
	return func(c *Command, name, key string) {
		c.MarkFlagRequired(key)

		key = requiredKey(key)
		viper.Set(key, true)

		u := c.Flag(name).Usage
		c.Flag(name).Usage = fmt.Sprintf("%s %s", u, requiredColor("(required)"))
	}
}

func requiredKey(key string) string {
	return fmt.Sprintf("required.%s", key)
}

func betaOpt() flagOpt {
	return func(c *Command, name, key string) {
		c.Flag(name).Hidden = !isBeta()
	}
}

func isBeta() bool {
	return viper.GetBool("enable-beta")
}

// AddStringFlag adds a string flag to a command.
func AddStringFlag(cmd *Command, name, shorthand, dflt, desc string, opts ...flagOpt) {
	fn := flagName(cmd, name)
	cmd.Flags().StringP(name, shorthand, dflt, desc)

	for _, o := range opts {
		o(cmd, name, fn)
	}

	viper.BindPFlag(fn, cmd.Flags().Lookup(name))
}

// AddIntFlag adds an integr flag to a command.
func AddIntFlag(cmd *Command, name, shorthand string, def int, desc string, opts ...flagOpt) {
	fn := flagName(cmd, name)
	cmd.Flags().IntP(name, shorthand, def, desc)
	viper.BindPFlag(fn, cmd.Flags().Lookup(name))

	for _, o := range opts {
		o(cmd, name, fn)
	}
}

// AddBoolFlag adds a boolean flag to a command.
func AddBoolFlag(cmd *Command, name, shorthand string, def bool, desc string, opts ...flagOpt) {
	fn := flagName(cmd, name)
	cmd.Flags().BoolP(name, shorthand, def, desc)
	viper.BindPFlag(fn, cmd.Flags().Lookup(name))

	for _, o := range opts {
		o(cmd, name, fn)
	}
}

// AddStringSliceFlag adds a string slice flag to a command.
func AddStringSliceFlag(cmd *Command, name, shorthand string, def []string, desc string, opts ...flagOpt) {
	fn := flagName(cmd, name)
	cmd.Flags().StringSliceP(name, shorthand, def, desc)
	viper.BindPFlag(fn, cmd.Flags().Lookup(name))

	for _, o := range opts {
		o(cmd, name, fn)
	}
}

func flagName(cmd *Command, name string) string {
	parentName := types.NSRoot
	if cmd.Parent() != nil {
		parentName = cmd.Parent().Name()
	}

	return fmt.Sprintf("%s.%s.%s", parentName, cmd.Name(), name)
}

func cmdNS(cmd *cobra.Command) string {
	parentName := types.NSRoot
	if cmd.Parent() != nil {
		parentName = cmd.Parent().Name()
	}

	return fmt.Sprintf("%s.%s", parentName, cmd.Name())
}

// CmdRunner runs a command and passes in a cmdConfig.
type CmdRunner func(*CmdConfig) error

// CmdConfig is a command configuration.
type CmdConfig struct {
	NS   string
	Ankr types.Config
	Out  io.Writer
	Args []string

	getContextAccessToken func() (string, string)
	setContextAccessToken func(string, string)
}

// NewCmdConfig creates an instance of a CmdConfig.
func NewCmdConfig(ns string, dc types.Config, out io.Writer, args []string) (*CmdConfig, error) {

	cmdConfig := &CmdConfig{
		NS:   ns,
		Ankr: dc,
		Out:  out,
		Args: args,

		getContextAccessToken: func() (string, string) {
			context := Context
			if context == "" {
				context = viper.GetString("context")
			}
			token := ""
			userid := ""

			switch context {
			case "default":
				token = viper.GetString(types.ArgAccessToken)
				userid = viper.GetString(types.ArgUserID)
			default:
				contexts := viper.GetStringMapString("auth-contexts")
				userid = contexts[types.ArgUserID]
				token = contexts[types.ArgAccessToken]
			}
			return token, userid
		},

		setContextAccessToken: func(token string, userid string) {
			context := Context
			if context == "" {
				context = viper.GetString("context")
			}

			switch context {
			case "default":
				viper.Set(types.ArgAccessToken, token)
				viper.Set(types.ArgUserID, userid)
			default:
				contexts := viper.GetStringMapString("auth-contexts")
				contexts[types.ArgAccessToken] = token
				contexts[types.ArgUserID] = userid

				viper.Set("auth-contexts", contexts)
			}
		},
	}

	return cmdConfig, nil
}

// Display displayes the output from a command.
func (c *CmdConfig) Display(d displayers.Displayable) error {
	dc := &displayers.Displayer{
		NS:     c.NS,
		Config: c.Ankr,
		Item:   d,
		Out:    c.Out,
	}

	return dc.Display()
}

// CmdBuilder builds a new command.
func CmdBuilder(parent *Command, cr CmdRunner, cliText, desc string, out io.Writer, options ...cmdOption) *Command {
	return cmdBuilderWithInit(parent, cr, cliText, desc, out, options...)
}

func cmdBuilderWithInit(parent *Command, cr CmdRunner, cliText, desc string, out io.Writer, options ...cmdOption) *Command {
	cc := &cobra.Command{
		Use:   cliText,
		Short: desc,
		Long:  desc,
		Run: func(cmd *cobra.Command, args []string) {
			c, err := NewCmdConfig(
				cmdNS(cmd),
				types.AnkrConfig,
				out,
				args,
			)
			checkErr(err, cmd)

			err = cr(c)
			checkErr(err, cmd)
		},
	}

	c := &Command{Command: cc}

	if parent != nil {
		parent.AddCommand(c)
	}

	for _, co := range options {
		co(c)
	}

	if cols := c.fmtCols; cols != nil {
		formatHelp := fmt.Sprintf("Columns for output in a comma separated list. Possible values: %s",
			strings.Join(cols, ","))
		AddStringFlag(c, types.ArgFormat, "", "", formatHelp)
		AddBoolFlag(c, types.ArgNoHeader, "", false, "hide headers")
	}

	return c

}

func writeConfig() error {
	f, err := cfgFileWriter()
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		return errors.New("unable to encode configuration to YAML format")
	}

	_, err = f.Write(b)
	if err != nil {
		return errors.New("unable to write configuration")
	}

	return nil
}

func defaultConfigFileWriter() (io.WriteCloser, error) {
	f, err := os.Create(cfgFile)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(cfgFile, 0600); err != nil {
		return nil, err
	}

	return f, nil
}
