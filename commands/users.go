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
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"

	akrctl "github.com/Ankr-network/dccn-cli"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"

	"github.com/spf13/cobra"

	"context"

	ankr_const "github.com/Ankr-network/dccn-common"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// userCmd creates the user command.
func userCmd() *Command {
	//DCCN-CLI user
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "user",
			Aliases: []string{"u"},
			Short:   "user commands",
			Long:    "user is used to access user commands",
		},
		DocCategories: []string{"user"},
		IsIndex:       true,
	}

	//DCCN-CLI user register
	cmdUserRegister := CmdBuilder(cmd, RunUserRegister, "register <user-name>", "user register", Writer,
		aliasOpt("rg"), docCategories("user"))
	AddStringFlag(cmdUserRegister, akrctl.ArgNicknameSlug, "", "", "User nickname")
	AddStringFlag(cmdUserRegister, akrctl.ArgPasswordSlug, "", "", "User password")
	AddStringFlag(cmdUserRegister, akrctl.ArgEmailSlug, "", "", "User email")
	AddStringFlag(cmdUserRegister, akrctl.ArgBalanceSlug, "", "", "User balance")

	//DCCN-CLI user Login
	cmdUserLogin := CmdBuilder(cmd, RunUserLogin, "login <user-email>", "user login", Writer,
		aliasOpt("li"), docCategories("user"))
	AddStringFlag(cmdUserLogin, akrctl.ArgPasswordSlug, "", "", "User password")

	//DCCN-CLI user logout
	cmdUserLogout := CmdBuilder(cmd, RunUserLogout, "logout", "user logout", Writer,
		aliasOpt("lo"), docCategories("user"))
	_ = cmdUserLogout

	return cmd

}

// RunUserRegister register a user.
func RunUserRegister(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	nickname, err := c.Ankr.GetString(c.NS, akrctl.ArgNicknameSlug)
	if err != nil {
		return err
	}

	password, err := c.Ankr.GetString(c.NS, akrctl.ArgPasswordSlug)
	if err != nil {
		return err
	}

	email, err := c.Ankr.GetString(c.NS, akrctl.ArgEmailSlug)
	if err != nil {
		return err
	}

	balance, err := c.Ankr.GetString(c.NS, akrctl.ArgBalanceSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	userClient := usermgr.NewUserMgrClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	errs := make(chan *common_proto.Error, len(c.Args))
	for _, name := range c.Args {
		user := &usermgr.User{
			Name:     name,
			Nickname: nickname,
			Email:    email,
			Password: password,
		}
		if balance != "" {
			b, err := strconv.Atoi(balance)
			if err != nil {
				return fmt.Errorf("balance %s is not an int", balance)
			}
			user.Balance = int32(b)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := userClient.Register(ctx, user); err != nil {
				log.Fatal(err.Error())
			} else {
				fmt.Printf("User %s Register Success.\n", email)
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			return errors.New(err.Details)
		}
	}

	return nil
}

// RunUserLogin login user by email and password.
func RunUserLogin(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	password, err := c.Ankr.GetString(c.NS, akrctl.ArgPasswordSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	userClient := usermgr.NewUserMgrClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ankr_const.ClientTimeOut*time.Second)
	defer cancel()
	fn := func(emails []string) error {
		for _, email := range emails {
			if rsp, err := userClient.Login(ctx, &usermgr.LoginRequest{Email: email, Password: password}); err != nil {
				log.Fatal(err.Error())
			} else {
				fmt.Printf("User %s Login Success, Token: %s\n", email, rsp.Token)
				c.setContextAccessToken(rsp.Token, rsp.UserId)
			}
		}
		if err := writeConfig(); err != nil {
			return fmt.Errorf(err.Error())
		}
		return nil
	}
	return fn(c.Args)
}

// RunUserLogout logout user.
func RunUserLogout(c *CmdConfig) error {

	url := viper.GetString("hub-url")

	token, _ := c.getContextAccessToken()

	if token == "" {
		return fmt.Errorf("unable to read AnkrNetwork access token")
	}

	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	userClient := usermgr.NewUserMgrClient(conn)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()
	if _, err := userClient.Logout(tokenContext, &usermgr.LogoutRequest{}); err != nil {
		return fmt.Errorf(err.Error())
	} else {
		fmt.Println("Logout Success.")
	}
	return nil
}
