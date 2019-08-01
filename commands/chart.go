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
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"k8s.io/helm/pkg/chartutil"

	ankrctl "github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"

	"context"

	ankr_const "github.com/Ankr-network/dccn-common"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	gwtaskmgr "github.com/Ankr-network/dccn-common/protos/gateway/taskmgr/v1"
	gwusermgr "github.com/Ankr-network/dccn-common/protos/gateway/usermgr/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Chart creates the chart command.
func chartCmd() *Command {
	//DCCN-CLI cluster
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "chart",
			Aliases: []string{"n"},
			Short:   "chart commands",
			Long:    "chart is used to access chart commands",
		},
		DocCategories: []string{"chart"},
		IsIndex:       true,
	}

	//DCCN-CLI chart upload
	cmdRunChartUpload := CmdBuilder(cmd, RunChartUpload, "upload <upload-name>", "create chart", Writer,
		aliasOpt("cr"), docCategories("chart"))
	AddStringFlag(cmdRunChartUpload, ankrctl.ArgUploadVersionSlug, "", "", "Chart Version", requiredOpt())
	AddStringFlag(cmdRunChartUpload, ankrctl.ArgUploadFileSlug, "", "", "Chart File", requiredOpt())

	//DCCN-CLI chart list
	cmdRunChartList := CmdBuilder(cmd, RunChartList, "list [GLOB]", "list chart", Writer,
		aliasOpt("ls"), displayerType(&displayers.Chart{}), docCategories("chart"))
	AddStringFlag(cmdRunChartList, ankrctl.ArgListRepoSlug, "", "", "List Repo")

	//DCCN-CLI chart detail
	cmdRunChartDetail := CmdBuilder(cmd, RunChartDetail, "detail <detail-name>", "get chart details", Writer,
		aliasOpt("dt"), docCategories("chart"))
	AddStringFlag(cmdRunChartDetail, ankrctl.ArgDetailRepoSlug, "", "", "Detail Repo", requiredOpt())
	AddStringFlag(cmdRunChartDetail, ankrctl.ArgShowVersionSlug, "", "", "Show Version", requiredOpt())

	//DCCN-CLI chart update
	cmdRunChartSaveas := CmdBuilder(cmd, RunChartSaveas, "saveas <saveas-name>", "saveas chart", Writer,
		aliasOpt("ud"), docCategories("chart"))
	AddStringFlag(cmdRunChartSaveas, ankrctl.ArgSourceRepoSlug, "", "", "Source Repo", requiredOpt())
	AddStringFlag(cmdRunChartSaveas, ankrctl.ArgSourceVersionSlug, "", "", "Source Version", requiredOpt())
	AddStringFlag(cmdRunChartSaveas, ankrctl.ArgSourceNameSlug, "", "", "Source Name", requiredOpt())
	AddStringFlag(cmdRunChartSaveas, ankrctl.ArgSaveasVersionSlug, "", "", "SaveAs Version", requiredOpt())
	AddStringFlag(cmdRunChartSaveas, ankrctl.ArgValuesYamlSlug, "", "", "Values Yaml File", requiredOpt())

	//DCCN-CLI chart download
	cmdRunChartDownload := CmdBuilder(cmd, RunChartDownload, "download <download-name>",
		"download chart", Writer, aliasOpt("dl"), docCategories("chart"))
	AddStringFlag(cmdRunChartDownload, ankrctl.ArgDownloadRepoSlug, "", "", "Download Repo", requiredOpt())
	AddStringFlag(cmdRunChartDownload, ankrctl.ArgDownloadVersionSlug, "", "", "Download Version", requiredOpt())

	//DCCN-CLI chart delete
	cmdRunChartDelete := CmdBuilder(cmd, RunChartDelete, "delete <delete-name>", "delete chart",
		Writer, aliasOpt("dl"), docCategories("chart"))
	AddStringFlag(cmdRunChartDelete, ankrctl.ArgDeleteVersionSlug, "", "", "Chart Version", requiredOpt())
	AddBoolFlag(cmdRunChartDelete, ankrctl.ArgForce, ankrctl.ArgShortForce, false, "Force chart delete")

	return cmd
}

// RunChartUpload upload a new chart.
func RunChartUpload(c *CmdConfig) error {
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
	}
	defer conn.Close()

	appClient := gwtaskmgr.NewAppMgrClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}

	uploadChartRequest := &gwtaskmgr.UploadChartRequest{}
	uploadChartRequest.ChartName = c.Args[0]

	uploadChartRequest.ChartVer, err = c.Ankr.GetString(c.NS, ankrctl.ArgUploadVersionSlug)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}
	uploadChartRequest.ChartRepo = "user"

	file, err := c.Ankr.GetString(c.NS, ankrctl.ArgUploadFileSlug)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}
	uploadChartRequest.ChartFile, err = ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}

	_, err = appClient.UploadChart(tokenctx, uploadChartRequest)
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}

	fmt.Printf("Chart %s upload success. \n", c.Args[0])

	return nil

}

// RunChartList returns a list of chart.
func RunChartList(c *CmdConfig) error {

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

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	appClient := gwtaskmgr.NewAppMgrClient(conn)

	chartRepo, err := c.Ankr.GetString(c.NS, ankrctl.ArgListRepoSlug)
	if err != nil {
		return err
	}
	r, err := appClient.ChartList(tokenctx, &gwtaskmgr.ChartListRequest{ChartRepo: chartRepo})
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}
	var matchedList []*common_proto.Chart

	for _, chart := range r.Charts {
		var skip = true
		if len(matches) == 0 {
			skip = false
		} else {
			for _, m := range matches {
				if m.Match(chart.ChartName) {
					skip = false
				}
			}
		}

		if !skip {
			matchedList = append(matchedList, chart)
		}
	}
	item := &displayers.Chart{Charts: matchedList}
	return c.Display(item)
}

// RunChartDetail returns chart details.
func RunChartDetail(c *CmdConfig) error {
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
	}
	defer conn.Close()
	appClient := gwtaskmgr.NewAppMgrClient(conn)

	chartDetailRequest := &gwtaskmgr.ChartDetailRequest{ChartName: c.Args[0]}
	chartDetailRequest.ChartRepo, err = c.Ankr.GetString(c.NS, ankrctl.ArgDetailRepoSlug)
	if err != nil {
		return err
	}
	chartDetailRequest.ChartVer, err = c.Ankr.GetString(c.NS, ankrctl.ArgShowVersionSlug)
	if err != nil {
		return err
	}

	r, err := appClient.ChartDetail(tokenctx, chartDetailRequest)
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}
	if r.ChartVersionDetails != nil {
		fmt.Printf("Repo: %s\tChart: %s \n", r.ChartRepo, r.ChartName)
		fmt.Println("Version\t\tApp Version")
		for _, versiondetail := range r.ChartVersionDetails {
			fmt.Printf("%s\t\t%s\n", versiondetail.ChartVer, versiondetail.ChartAppVer)
		}
	}
	fmt.Printf("\n++++++++++ Chart versions %s readme.md ++++++++++\n%s\n", chartDetailRequest.ChartVer, r.ReadmeMd)
	fmt.Printf("\n++++++++++ Chart versions %s values.yaml ++++++++++\n%s\n", chartDetailRequest.ChartVer, r.ValuesYaml)
	return nil
}

// RunChartSaveas save as new version of chart with updated value.
func RunChartSaveas(c *CmdConfig) error {
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
	}
	defer conn.Close()

	appClient := gwtaskmgr.NewAppMgrClient(conn)
	if err != nil {
		return err
	}
	saveasChartRequest := &gwtaskmgr.SaveAsChartRequest{}

	saveasVer, err := c.Ankr.GetString(c.NS, ankrctl.ArgSaveasVersionSlug)
	if err != nil {
		return err
	}

	saveasChartRequest.Destination = &gwtaskmgr.Destination{
		SaveasName: c.Args[0],
		SaveasVer:  saveasVer,
	}

	sourceName, err := c.Ankr.GetString(c.NS, ankrctl.ArgSourceNameSlug)
	if err != nil {
		return err
	}

	sourceVer, err := c.Ankr.GetString(c.NS, ankrctl.ArgSourceVersionSlug)
	if err != nil {
		return err
	}

	sourceRepo, err := c.Ankr.GetString(c.NS, ankrctl.ArgSourceRepoSlug)
	if err != nil {
		return err
	}

	saveasChartRequest.Source = &gwtaskmgr.Source{
		ChartName: sourceName,
		ChartRepo: sourceRepo,
		ChartVer:  sourceVer,
	}
	file, err := c.Ankr.GetString(c.NS, ankrctl.ArgValuesYamlSlug)
	if err != nil {
		return err
	}

	saveasChartRequest.ValuesYaml, err = ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	_, err = appClient.SaveAsChart(tokenctx, saveasChartRequest)
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}
	fmt.Printf("Chart %s save success.\n", c.Args[0])

	return nil

}

// RunChartDownload download chart to local file system.
func RunChartDownload(c *CmdConfig) error {
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
	}
	defer conn.Close()

	appClient := gwtaskmgr.NewAppMgrClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}
	downloadChartRequest := &gwtaskmgr.DownloadChartRequest{}

	downloadChartRequest.ChartName = c.Args[0]
	downloadChartRequest.ChartVer, err = c.Ankr.GetString(c.NS, ankrctl.ArgDownloadVersionSlug)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}

	downloadChartRequest.ChartRepo, err = c.Ankr.GetString(c.NS, ankrctl.ArgDownloadRepoSlug)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}

	rsp, err := appClient.DownloadChart(tokenctx, downloadChartRequest)
	if err != nil {
		return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
	}

	chart, err := chartutil.LoadArchive(bytes.NewReader(rsp.ChartFile))
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}

	dest, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}

	name, err := chartutil.Save(chart, dest)
	if err == nil {
		fmt.Printf("Successfully download chart and saved it to: %s\n", name)
	} else {
		fmt.Fprintf(os.Stdout, "\nERROR: %s\n",err.Error())
		return nil
	}

	return nil

}

// RunChartDelete delete a chart.
func RunChartDelete(c *CmdConfig) error {
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
	chartVersion, err := c.Ankr.GetString(c.NS, ankrctl.ArgDeleteVersionSlug)
	if err != nil {
		return err
	}
	if force || AskForConfirm(fmt.Sprintf("Are you sure you want to Delete chart %s version %s (y/N) ? ", c.Args[0], chartVersion)) == nil {

		url := viper.GetString("hub-url")

		conn, err := grpc.Dial(url+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}

		defer conn.Close()
		appClient := gwtaskmgr.NewAppMgrClient(conn)

		_, err = appClient.DeleteChart(tokenctx, &gwtaskmgr.DeleteChartRequest{
			ChartName: c.Args[0],
			ChartRepo: "user",
			ChartVer:  chartVersion,
		})
		if err != nil {
			return fmt.Errorf("Status Code: %s  Message: %s", status.Code(err), err.Error())
		}
		fmt.Printf("Chart %s version %s delete success.\n", c.Args[0], chartVersion)

	}

	return nil
}
