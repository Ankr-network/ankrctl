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
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"

	ankrctl "github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"

	"context"
	"strconv"

	ankr_const "github.com/Ankr-network/dccn-common"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	gwtaskmgr "github.com/Ankr-network/dccn-common/protos/gateway/taskmgr/v1"
	gwusermgr "github.com/Ankr-network/dccn-common/protos/gateway/usermgr/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Namespace creates the namespace command.
func namespaceCmd() *Command {
	//DCCN-CLI cluster
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "namespace",
			Aliases: []string{"n"},
			Short:   "namespace commands",
			Long:    "namespace is used to access namespace commands",
		},
		DocCategories: []string{"namespace"},
		IsIndex:       true,
	}

	//DCCN-CLI namespace create
	cmdRunNamespaceCreate := CmdBuilder(cmd, RunNamespaceCreate, "create <ns-name> [ns-name ...]", "create namespace", Writer,
		aliasOpt("cr"), docCategories("namespace"))
	AddStringFlag(cmdRunNamespaceCreate, ankrctl.ArgNsCpuLimitSlug, "", "", "Namespace CPU Limit (in vCPUs)", requiredOpt())
	AddStringFlag(cmdRunNamespaceCreate, ankrctl.ArgNsMemLimitSlug, "", "", "Namespace MEM Limit (in GiB)", requiredOpt())
	AddStringFlag(cmdRunNamespaceCreate, ankrctl.ArgNsStorageLimitSlug, "", "", "Namespace Storage Limit (in GiB)", requiredOpt())
	AddStringFlag(cmdRunNamespaceCreate, ankrctl.ArgNsClusterIDSlug, "", "", "Namespace Cluster Id")

	//DCCN-CLI namespace list
	cmdRunNamespaceList := CmdBuilder(cmd, RunNamespaceList, "list [GLOB]", "list namespace", Writer,
		aliasOpt("ls"), displayerType(&displayers.Cluster{}), docCategories("namespace"))
	_ = cmdRunNamespaceList

	//DCCN-CLI namespace update
	cmdRunNamespaceUpdate := CmdBuilder(cmd, RunNamespaceUpdate, "update <namespace-id> [namespace-id ...]", "update namespace", Writer,
		aliasOpt("ud"), docCategories("namespace"))
	AddStringFlag(cmdRunNamespaceUpdate, ankrctl.ArgNsCpuLimitSlug, "", "", "Namespace CPU Limit (in vCPUs)", requiredOpt())
	AddStringFlag(cmdRunNamespaceUpdate, ankrctl.ArgNsMemLimitSlug, "", "", "Namespace MEM Limit (in GiB)", requiredOpt())
	AddStringFlag(cmdRunNamespaceUpdate, ankrctl.ArgNsStorageLimitSlug, "", "", "Namespace Storage Limit (in GiB)", requiredOpt())

	//DCCN-CLI namespace delete
	cmdRunNamespaceDelete := CmdBuilder(cmd, RunNamespaceDelete, "delete <namespace-id> [namespace-id ...]", "delete namespace",
		Writer, aliasOpt("dl"), docCategories("namespace"))
	AddBoolFlag(cmdRunNamespaceDelete, ankrctl.ArgForce, ankrctl.ArgShortForce, false, "Force namespace delete")

	return cmd
}

// RunNamespaceCreate create a namespace.
func RunNamespaceCreate(c *CmdConfig) error {

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
		log.Fatalf("did not connect: %v", err)
		return nil
	}
	defer conn.Close()

	appClient := gwtaskmgr.NewAppMgrClient(conn)
	if err != nil {
		return err
	}

	cpuLimit, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsCpuLimitSlug)
	if err != nil {
		return err
	}

	nsCpuLimit, err := strconv.ParseFloat(cpuLimit, 32)
	if err != nil {
		return err
	}

	memLimit, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsMemLimitSlug)
	if err != nil {
		return err
	}
	nsMemLimit, err := strconv.ParseFloat(memLimit, 32)
	if err != nil {
		return err
	}

	storageLimit, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsStorageLimitSlug)
	if err != nil {
		return err
	}

	nsStorageLimit, err := strconv.ParseFloat(storageLimit, 32)
	if err != nil {
		return err
	}

	nsClusterId, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsClusterIDSlug)
	if err != nil {
		return err
	}

	createNamespaceRequest := &gwtaskmgr.CreateNamespaceRequest{
		NsCpuLimit:     uint32(nsCpuLimit * 1000),
		NsMemLimit:     uint32(nsMemLimit * 1024),
		NsStorageLimit: uint32(nsStorageLimit * 1024),
		ClusterId:      nsClusterId,
	}

	var wg sync.WaitGroup
	errs := make(chan error, len(c.Args))
	for _, name := range c.Args {
		createNamespaceRequest.NsName = name

		wg.Add(1)
		go func() {
			defer wg.Done()
			rsp, err := appClient.CreateNamespace(tokenctx, createNamespaceRequest)
			if err != nil {
				errs <- err
			} else {
				if rsp != nil {
					fmt.Printf("Namespace %s create success. \n", rsp.NsId)
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

// RunNamespaceList returns a list of namespace.
func RunNamespaceList(c *CmdConfig) error {

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

	var matchedList []common_proto.NamespaceReport

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil
	}
	defer conn.Close()
	appClient := gwtaskmgr.NewAppMgrClient(conn)

	r, err := appClient.NamespaceList(tokenctx, &common_proto.Empty{})
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}

	for _, nsReport := range r.NsReports {
		var skip = true
		if len(matches) == 0 {
			skip = false
		} else {
			for _, m := range matches {
				if m.Match(nsReport.Namespace.NsName) {
					skip = false
				}
			}
		}

		if !skip {
			matchedList = append(matchedList, *nsReport)
		}
	}
	item := &displayers.Namespace{Namespaces: matchedList}
	return c.Display(item)
}

// RunNamespaceUpdate update the namespace setting.
func RunNamespaceUpdate(c *CmdConfig) error {

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
		log.Fatalf("did not connect: %v", err)
		return nil
	}
	defer conn.Close()

	appClient := gwtaskmgr.NewAppMgrClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}

	updateNamespaceRequest := &gwtaskmgr.UpdateNamespaceRequest{}
	updateNamespaceRequest.NsId = c.Args[0]

	cpuLimit, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsCpuLimitSlug)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}
	nsCpuLimit, err := strconv.ParseUint(cpuLimit, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}

	updateNamespaceRequest.NsCpuLimit = uint32(nsCpuLimit)

	memLimit, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsMemLimitSlug)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}
	nsMemLimit, err := strconv.ParseUint(memLimit, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}
	updateNamespaceRequest.NsMemLimit = uint32(nsMemLimit)

	storageLimit, err := c.Ankr.GetString(c.NS, ankrctl.ArgNsStorageLimitSlug)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}

	nsStorageLimit, err := strconv.ParseUint(storageLimit, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}
	updateNamespaceRequest.NsStorageLimit = uint32(nsStorageLimit)

	fn := func(ids []string) error {
		for _, id := range ids {
			_, err := appClient.UpdateNamespace(tokenctx, updateNamespaceRequest)
			if err != nil {
				return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
			}
			fmt.Printf("Namespace %s update success.\n", id)
		}
		return nil
	}

	return fn(c.Args)

}

// RunNamespaceDelete delete a namespace.
func RunNamespaceDelete(c *CmdConfig) error {

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

	force, err := c.Ankr.GetBool(c.NS, ankrctl.ArgForce)
	if err != nil {
		return err
	}

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenctx, cancel := context.WithTimeout(ctx, ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	if force || AskForConfirm(fmt.Sprintf("Are you sure you want to Cancel %d namespace(s) (y/N) ? ", len(c.Args))) == nil {
		url := viper.GetString("hub-url")

		conn, err := grpc.Dial(url+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
			return nil
		}

		defer conn.Close()
		appClient := gwtaskmgr.NewAppMgrClient(conn)

		fn := func(ids []string) error {
			for _, id := range ids {
				_, err := appClient.DeleteNamespace(tokenctx, &gwtaskmgr.DeleteNamespaceRequest{NsId: id})
				if err != nil {
					return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
				}
				fmt.Printf("Namespace %s delete success.\n", id)
			}
			return nil
		}

		return fn(c.Args)
	}

	return nil
}
