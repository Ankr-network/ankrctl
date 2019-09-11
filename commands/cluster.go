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

	"github.com/Ankr-network/ankrctl/commands/displayers"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"

	"context"

	ankr_const "github.com/Ankr-network/dccn-common"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	gwdcmgr "github.com/Ankr-network/dccn-common/protos/gateway/dcmgr/v1"
	gwusermgr "github.com/Ankr-network/dccn-common/protos/gateway/usermgr/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Cluster creates the cluster command.
func clusterCmd() *Command {
	//DCCN-CLI cluster
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "cluster",
			Aliases: []string{"c"},
			Short:   "cluster commands",
			Long:    "cluster is used to access datacenter commands",
		},
		DocCategories: []string{"cluster"},
		IsIndex:       true,
	}

	//DCCN-CLI cluster list
	cmdRunClusterList := CmdBuilder(cmd, RunClusterList, "list [GLOB]", "list cluster", Writer,
		aliasOpt("ls"), displayerType(&displayers.Cluster{}), docCategories("cluster"))
	_ = cmdRunClusterList

	//DCCN-CLI cluster network info
	cmdRunNetworkInfo := CmdBuilder(cmd, RunNetworkInfo, "network", "list network info", Writer,
		aliasOpt("ni"), docCategories("cluster"))
	_ = cmdRunNetworkInfo

	return cmd
}

// RunClusterList returns a list of cluster.
func RunClusterList(c *CmdConfig) error {

	authResult := gwusermgr.AuthenticationResult{}
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

	var matchedList []common_proto.DataCenterStatus

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	dcMgr := gwdcmgr.NewDCAPIClient(conn)

	r, err := dcMgr.DataCenterList(tokenctx, &common_proto.Empty{})
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}

	for _, cluster := range r.DcList {
		var skip = true
		if len(matches) == 0 {
			skip = false
		} else {
			for _, m := range matches {
				if m.Match(cluster.DcName) {
					skip = false
				}
			}
		}

		if !skip {
			if cluster.GeoLocation == nil {
				cluster.GeoLocation = &common_proto.GeoLocation{}
			}
			matchedList = append(matchedList, *cluster)
		}
	}
	item := &displayers.Cluster{Clusters: matchedList}
	return c.Display(item)
}

// RunNetworkInfo returns a overview of apps.
func RunNetworkInfo(c *CmdConfig) error {

	authResult := gwusermgr.AuthenticationResult{}
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

	dcMgr := gwdcmgr.NewDCAPIClient(conn)
	if err != nil {
		return err
	}
	resp, err := dcMgr.NetworkInfo(tokenctx, &common_proto.Empty{})
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}
	fmt.Printf("User Count:\t\t%v\nHost Count:\t\t%v\nNamespace Count:\t%v\nContainer Count:\t%v\nTraffic:\t%v\n",
		resp.UserCount, resp.HostCount, resp.NsCount, resp.ContainerCount, resp.Traffic)

	return nil
}
