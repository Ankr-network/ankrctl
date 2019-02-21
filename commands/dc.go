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
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"

	"github.com/Ankr-network/dccn-cli/commands/displayers"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"

	"context"

	ankr_const "github.com/Ankr-network/dccn-common"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Dc creates the dc command.
func Dc() *Command {
	//DCCN-CLI dc
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "dc",
			Aliases: []string{"d"},
			Short:   "dc commands",
			Long:    "dc is used to access datacenter commands",
		},
		DocCategories: []string{"dc"},
		IsIndex:       true,
	}

	//DCCN-CLI dc list
	cmdRunDcList := CmdBuilder(cmd, RunDcList, "list [GLOB]", "list dc", Writer,
		aliasOpt("ls"), displayerType(&displayers.Dc{}), docCategories("dc"))
	_ = cmdRunDcList

	//DCCN-CLI dc network info
	cmdRunNetworkInfo := CmdBuilder(cmd, RunNetworkInfo, "network", "list network info", Writer,
		aliasOpt("ni"), docCategories("dc"))
	_ = cmdRunNetworkInfo

	//DCCN-CLI dc leader board
	cmdRunDcLeader := CmdBuilder(cmd, RunDcLeader, "leader", "list dc leader board", Writer,
		aliasOpt("lb"), docCategories("dc"))
	_ = cmdRunDcLeader

	return cmd
}

// RunDcList returns a list of dc.
func RunDcList(c *CmdConfig) error {

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found")
	}

	md := metadata.New(map[string]string{
		"token": authResult.AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	matches := []glob.Glob{}
	for _, globStr := range c.Args {
		g, err := glob.Compile(globStr)
		if err != nil {
			return fmt.Errorf("unknown glob %q", globStr)
		}

		matches = append(matches, g)
	}

	var matchedList []common_proto.DataCenter

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	dcClient := dcmgr.NewDCAPIClient(conn)

	r, err := dcClient.DataCenterList(tokenctx, &common_proto.Empty{})
	if err != nil {
		log.Fatalf("Client: could not send: %v", err)
	}

	for _, dc := range r.DcList {
		var skip = true
		if len(matches) == 0 {
			skip = false
		} else {
			for _, m := range matches {
				if m.Match(dc.Name) {
					skip = false
				}
			}
		}

		if !skip {
			if dc.GeoLocation == nil {
				dc.GeoLocation = &common_proto.GeoLocation{}
			}
			matchedList = append(matchedList, *dc)
		}
	}
	fmt.Println(matchedList)
	item := &displayers.Dc{Dcs: matchedList}
	return c.Display(item)
}

// RunNetworkInfo returns a overview of tasks.
func RunNetworkInfo(c *CmdConfig) error {

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found")
	}
	md := metadata.New(map[string]string{
		"token": authResult.AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	dcMgr := dcmgr.NewDCAPIClient(conn)
	if err != nil {
		return err
	}
	resp, err := dcMgr.NetworkInfo(tokenctx, &common_proto.Empty{})
	if err != nil {
		return err
	}
	fmt.Printf("User Count:\t\t%v\nHost Count:\t\t%v\nEnvironment Count:\t%v\nContainer Count:\t%v\n",
		resp.UserCount, resp.HostCount, resp.EnvironmentCount, resp.ContainerCount)

	return nil
}

// RunDcLeader returns a overview of tasks.
func RunDcLeader(c *CmdConfig) error {

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found")
	}
	md := metadata.New(map[string]string{
		"token": authResult.AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	dcMgr := dcmgr.NewDCAPIClient(conn)
	if err != nil {
		return err
	}
	resp, err := dcMgr.DataCenterLeaderBoard(tokenctx, &common_proto.Empty{})
	if err != nil {
		return err
	}
	fmt.Println("Name\t\tNumber")
	for _, d := range resp.List {
		fmt.Printf("%s\t%v\n", d.Name, d.Number)
	}

	return nil
}
