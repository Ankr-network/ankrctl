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
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"

	"context"

	ankrctl "github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	gwtaskmgr "github.com/Ankr-network/dccn-common/protos/gateway/taskmgr/v1"
	gwusermgr "github.com/Ankr-network/dccn-common/protos/gateway/usermgr/v1"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"

	ankr_const "github.com/Ankr-network/dccn-common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var port = ":" + strconv.Itoa(ankr_const.DefaultPort)

var clientURL string

// App creates the app command.
func appCmd() *Command {
	//DCCN-CLI app
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "app",
			Aliases: []string{"a"},
			Short:   "app commands",
			Long:    "app is used to manage application",
		},
		DocCategories: []string{"app"},
		IsIndex:       true,
	}

	//DCCN-CLI comput app create
	cmdRunAppCreate := CmdBuilder(cmd, RunAppCreate, "create <app-name> [app-name ...]",
		"create app", Writer, aliasOpt("cr"), docCategories("app"))
	AddStringFlag(cmdRunAppCreate, ankrctl.ArgChartNameSlug, "", "", "Chart name", requiredOpt())
	AddStringFlag(cmdRunAppCreate, ankrctl.ArgChartRepoSlug, "", "", "Chart repo", requiredOpt())
	AddStringFlag(cmdRunAppCreate, ankrctl.ArgChartVersionSlug, "", "", "Chart version", requiredOpt())
	AddStringFlag(cmdRunAppCreate, ankrctl.ArgNsIDSlug, "", "", "Namespace ID")
	AddStringFlag(cmdRunAppCreate, ankrctl.ArgNsNameSlug, "", "", "Namespace Name")
	AddStringFlag(cmdRunAppCreate, ankrctl.ArgNsCpuLimitSlug, "", "", "Namespace CPU Limit (mCPUs)")
	AddStringFlag(cmdRunAppCreate, ankrctl.ArgNsMemLimitSlug, "", "", "Namespace MEM Limit (MBs)")
	AddStringFlag(cmdRunAppCreate, ankrctl.ArgNsStorageLimitSlug, "", "", "Namespace Storage Limit (GBs)")

	//DCCN-CLI comput app cancel
	cmdRunAppCancel := CmdBuilder(cmd, RunAppCancel, "cancel <app-id> [app-id ...]",
		"Cancel app by id", Writer, aliasOpt("dl"), docCategories("app"))
	AddBoolFlag(cmdRunAppCancel, ankrctl.ArgForce, ankrctl.ArgShortForce, false, "Force app cancel")

	//DCCN-CLI comput app purge
	cmdRunAppPurge := CmdBuilder(cmd, RunAppPurge, "purge <app-id> [app-id ...]", "Purge app by id",
		Writer, aliasOpt("rm"), docCategories("app"))
	AddBoolFlag(cmdRunAppPurge, ankrctl.ArgForce, ankrctl.ArgShortForce, false, "Force app purge")

	//DCCN-CLI comput app update
	cmdRunAppUpdate := CmdBuilder(cmd, RunAppUpdate, "update <app-id> [app-id ...]",
		"Update app by id", Writer, aliasOpt("ud"), docCategories("app"))
	AddStringFlag(cmdRunAppUpdate, ankrctl.ArgAppNameSlug, "", "", "App name")
	AddStringFlag(cmdRunAppUpdate, ankrctl.ArgUpdateVersionSlug, "", "", "Update version")

	//DCCN-CLI app list
	cmdRunAppList := CmdBuilder(cmd, RunAppList, "list [GLOB]", "list apps", Writer,
		aliasOpt("ls"), displayerType(&displayers.AppReport{}), docCategories("app"))
	_ = cmdRunAppList

	//DCCN-CLI app detail
	cmdRunAppDetail := CmdBuilder(cmd, RunAppDetail, "detail <app-id>", "list app detail", Writer,
		aliasOpt("dt"), docCategories("app"))
	_ = cmdRunAppDetail

	//DCCN-CLI app overview
	cmdRunAppOverview := CmdBuilder(cmd, RunAppOverview, "overview", "show apps overview", Writer,
		aliasOpt("ov"), docCategories("app"))
	_ = cmdRunAppOverview

	return cmd

}

// RunAppCreate creates a app.
//DCCN-CLI comput app create
func RunAppCreate(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	createAppRequest := &gwtaskmgr.CreateAppRequest{}

	chartname, err := c.Ankr.GetString(c.NS, ankrctl.ArgChartNameSlug)
	if err != nil {
		return err
	}

	chartrepo, err := c.Ankr.GetString(c.NS, ankrctl.ArgChartRepoSlug)
	if err != nil {
		return err
	}

	chartver, err := c.Ankr.GetString(c.NS, ankrctl.ArgChartVersionSlug)
	if err != nil {
		return err
	}

	createAppRequest.Chart = &gwtaskmgr.Chart{
		ChartName: chartname,
		ChartRepo: chartrepo,
		ChartVer:  chartver,
	}

	nsID, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsIDSlug)
	if err != nil {
		return err
	}

	if nsID != "" {
		createAppRequest.NsId = nsID
	} else {
		nsname, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsNameSlug)
		if err != nil {
			return err
		}

		if len(nsname) == 0 {
			return ankrctl.NewMissingArgsErr(c.NS)
		}

		cpuLimit, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsCpuLimitSlug)
		if err != nil {
			return err
		}

		nsCpuLimit, err := strconv.ParseUint(cpuLimit, 10, 32)
		if nsCpuLimit == 0 || err != nil {
			return fmt.Errorf("Cpu Limit %s is not a valid number", cpuLimit)
		}

		memLimit, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsMemLimitSlug)
		if err != nil {
			return err
		}

		nsMemLimit, err := strconv.ParseUint(memLimit, 10, 32)
		if nsMemLimit == 0 || err != nil {
			return fmt.Errorf("Mem Limit %s is not a valid number", memLimit)
		}

		storageLimit, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsStorageLimitSlug)
		if err != nil {
			return err
		}

		nsStorageLimit, err := strconv.ParseUint(storageLimit, 10, 32)
		if nsStorageLimit == 0 || err != nil {
			return fmt.Errorf("Storage Limit %s is not a valid number", storageLimit)
		}

		createAppRequest.Namespace = &gwtaskmgr.Namespace{
			NsName:         nsname,
			NsCpuLimit:     uint32(nsCpuLimit),
			NsMemLimit:     uint32(nsMemLimit),
			NsStorageLimit: uint32(nsStorageLimit),
		}
	}

	url := viper.GetString("hub-url")

	authResult := gwusermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found")
	}

	md := metadata.New(map[string]string{
		"token": authResult.AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	appClient := gwtaskmgr.NewAppMgrClient(conn)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	errs := make(chan error, len(c.Args))
	for _, name := range c.Args {
		createAppRequest.AppName = name

		wg.Add(1)
		go func() {
			defer wg.Done()
			rsp, err := appClient.CreateApp(tokenctx, createAppRequest)
			if err != nil {
				errs <- err
			} else {
				if rsp != nil {
					fmt.Printf("App %s create success. \n", rsp.AppId)
				}
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
		}
	}

	return nil
}

// RunAppPurge purge a app from hub.
func RunAppPurge(c *CmdConfig) error {

	force, err := c.Ankr.GetBool(c.NS, ankrctl.ArgForce)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

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

	if force || AskForConfirm(fmt.Sprintf("Are you sure you want to Purge %d app(s) (y/N) ? ", len(c.Args))) == nil {
		url := viper.GetString("hub-url")

		conn, err := grpc.Dial(url+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()
		appClient := gwtaskmgr.NewAppMgrClient(conn)

		fn := func(ids []string) error {
			for _, id := range ids {
				_, err := appClient.PurgeApp(tokenctx, &gwtaskmgr.AppID{AppId: id})
				if err != nil {
					return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
				}
				fmt.Printf("App %s purge success.\n", id)
			}
			return nil
		}
		return fn(c.Args)

	}
	return fmt.Errorf("Operation aborted")

}

// RunAppCancel destroy a app by id.
func RunAppCancel(c *CmdConfig) error {

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

	force, err := c.Ankr.GetBool(c.NS, ankrctl.ArgForce)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	if force || AskForConfirm(fmt.Sprintf("Are you sure you want to Cancel %d app(s) (y/N) ? ", len(c.Args))) == nil {
		url := viper.GetString("hub-url")

		conn, err := grpc.Dial(url+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}

		defer conn.Close()
		appClient := gwtaskmgr.NewAppMgrClient(conn)

		fn := func(ids []string) error {
			for _, id := range ids {
				_, err := appClient.CancelApp(tokenctx, &gwtaskmgr.AppID{AppId: id})
				if err != nil {
					return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
				}
				fmt.Printf("App %s cancel success.\n", id)
			}
			return nil
		}

		return fn(c.Args)
	}
	return fmt.Errorf("Operation aborted")

}

// RunAppList returns a list of apps.
func RunAppList(c *CmdConfig) error {

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
			return fmt.Errorf("Unknown glob %q", globStr)
		}

		matches = append(matches, g)
	}

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	appClient := gwtaskmgr.NewAppMgrClient(conn)

	r, err := appClient.AppList(tokenctx, &common_proto.Empty{})
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}

	item := &displayers.AppReport{Apps: r.AppReports}
	return c.Display(item)
}

// RunAppDetail returns a list of apps.
func RunAppDetail(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

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
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	appClient := gwtaskmgr.NewAppMgrClient(conn)

	r, err := appClient.AppDetail(tokenctx, &gwtaskmgr.AppID{AppId: c.Args[0]})
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}
	fmt.Printf("Application %s resource detail:\n %s \n", c.Args[0], r.AppReport.Detail)

	return nil
}

// RunAppOverview returns a overview of apps.
func RunAppOverview(c *CmdConfig) error {

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
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	appClient := gwtaskmgr.NewAppMgrClient(conn)
	if err != nil {
		return err
	}
	tor, err := appClient.AppOverview(tokenctx, &common_proto.Empty{})
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}
	fmt.Printf("Cluster Count:\t\t%v\nNamespace Count:\t%v\nNetwork Count:\t\t%v\nTotal App Count:\t%v\nCluster Count:\t%v\n",
		tor.ClusterCount, tor.NamespaceCount, tor.NetworkCount, tor.TotalAppCount, tor.ClusterCount)
	fmt.Printf("Cpu Total:\t\t%v\nCpu Usage:\t%v\nMem Total:\t\t%v\nMem Usage:\t%v\nStorage Total:\t%v\nStorage Usage:\t%v\n",
		tor.CpuTotal, tor.CpuUsage, tor.MemTotal, tor.MemUsage, tor.StorageTotal, tor.StorageUsage)

	return nil
}

// RunAppUpdate updates a app.
func RunAppUpdate(c *CmdConfig) error {

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

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	updateAppRequest := &gwtaskmgr.UpdateAppRequest{}

	appname, err := c.Ankr.GetString(c.NS, ankrctl.ArgAppNameSlug)
	if err != nil {
		return err
	}
	if len(appname) > 0 {
		updateAppRequest.AppName = appname
	}

	chartver, err := c.Ankr.GetString(c.NS, ankrctl.ArgUpdateVersionSlug)
	if err != nil {
		return err
	}

	if len(chartver) > 0 {
		updateAppRequest.ChartVer = chartver
	}

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	appClient := gwtaskmgr.NewAppMgrClient(conn)

	fn := func(ids []string) error {
		for _, id := range ids {
			updateAppRequest.AppId = id

			_, err := appClient.UpdateApp(tokenctx, updateAppRequest)
			if err != nil {
				return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
			}
			fmt.Printf("App %s update success.\n", id)
		}
		return nil
	}
	return fn(c.Args)
}
