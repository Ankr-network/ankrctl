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
	"syscall"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"

	ankrctl "github.com/Ankr-network/dccn-cli"

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
	cmdUserRegister := CmdBuilder(cmd, RunUserRegister, "register <user-name>", "user register",
		Writer, aliasOpt("rg"), docCategories("user"))
	AddStringFlag(cmdUserRegister, ankrctl.ArgEmailSlug, "", "", "User email", requiredOpt())
	AddStringFlag(cmdUserRegister, ankrctl.ArgPasswordSlug, "", "", "User password", requiredOpt())

	//DCCN-CLI user comfirm registration
	cmdUserConfirmRegistration := CmdBuilder(cmd, RunUserConfirmRegistration,
		"confirm-registration <user-email>", "user registration confirmation", Writer,
		aliasOpt("rc"), docCategories("user"))
	AddStringFlag(cmdUserConfirmRegistration, ankrctl.ArgRegisterCodeSlug,
		"", "", "User registration confirmation code", requiredOpt())

	//DCCN-CLI user forgot password
	cmdUserForgotPassword := CmdBuilder(cmd, RunUserForgotPassword, "forgot-password <user-email>",
		"user password forgot", Writer, aliasOpt("fp"), docCategories("user"))
	_ = cmdUserForgotPassword

	//DCCN-CLI user comfirm password
	cmdUserConfirmPassword := CmdBuilder(cmd, RunUserConfirmPassword, "confirm-password <user-email>",
		"user password change confirmation", Writer, aliasOpt("pc"), docCategories("user"))
	AddStringFlag(cmdUserConfirmPassword, ankrctl.ArgPasswordCodeSlug,
		"", "", "User password change confirmation code", requiredOpt())
	AddStringFlag(cmdUserConfirmPassword, ankrctl.ArgConfirmPasswordSlug,
		"", "", "User confirm new password", requiredOpt())

	//DCCN-CLI user change password
	cmdUserChangePassword := CmdBuilder(cmd, RunUserChangePassword, "change-password <user-email>",
		"user password change", Writer, aliasOpt("cp"), docCategories("user"))
	AddStringFlag(cmdUserChangePassword, ankrctl.ArgOldPasswordSlug,
		"", "", "User old password", requiredOpt())
	AddStringFlag(cmdUserChangePassword, ankrctl.ArgNewPasswordSlug,
		"", "", "User new password", requiredOpt())

	//DCCN-CLI user change email
	cmdUserChangeEmail := CmdBuilder(cmd, RunUserChangeEmail, "email-change <new-email>",
		"user email change", Writer, aliasOpt("ec"), docCategories("user"))
	_ = cmdUserChangeEmail

	//DCCN-CLI user cnfirm email
	cmdUserConfirmEmail := CmdBuilder(cmd, RunUserConfirmEmail, "email-confirm <new-email>",
		"user confirm email change", Writer, aliasOpt("ce"), docCategories("user"))
	AddStringFlag(cmdUserConfirmEmail, ankrctl.ArgEmailCodeSlug,
		"", "", "User email change confirmation code", requiredOpt())

	//DCCN-CLI user update attribute
	cmdUserUpdate := CmdBuilder(cmd, RunUserUpdate, "update <user-email>", "user update attribute",
		Writer, aliasOpt("ua"), docCategories("user"))
	AddStringFlag(cmdUserUpdate, ankrctl.ArgUpdateKeySlug, "", "", "User attribute key", requiredOpt())
	AddStringFlag(cmdUserUpdate, ankrctl.ArgUpdateValueSlug, "", "", "User attribute value", requiredOpt())

	//DCCN-CLI user login
	cmdUserLogin := CmdBuilder(cmd, RunUserLogin, "login", "user login", Writer,
		aliasOpt("li"), docCategories("user"))
	_ = cmdUserLogin

	//DCCN-CLI user logout
	cmdUserLogout := CmdBuilder(cmd, RunUserLogout, "logout", "user logout", Writer,
		aliasOpt("lo"), docCategories("user"))
	_ = cmdUserLogout

	return cmd

}

// RunUserRegister register a user.
func RunUserRegister(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	email, err := c.Ankr.GetString(c.NS, ankrctl.ArgEmailSlug)
	if err != nil {
		fmt.Println(err)
		return err
	}

	password, err := c.Ankr.GetString(c.NS, ankrctl.ArgPasswordSlug)
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

	urr := &usermgr.RegisterRequest{
		Password: password,
		User: &usermgr.User{
			Email: email,
			Attributes: &usermgr.UserAttributes{
				Name: c.Args[0],
			},
		},
	}

	if _, err := userClient.Register(ctx, urr); err != nil {
		return err
	}

	fmt.Printf("User %s Register Success.\n", email)

	return nil
}

// RunUserLogin login user by email and password.
func RunUserLogin(c *CmdConfig) error {

	fmt.Print("\nEmail: ")
	email, err := retrieveUserInput()
	if err != nil {
		return err
	}

	fmt.Print("\nPassword: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
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
	ctx, cancel := context.WithTimeout(context.Background(),
		ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	rsp, err := userClient.Login(ctx,
		&usermgr.LoginRequest{Email: email, Password: string(password)})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("\n\nLogin Successful!\n\n")
	viper.Set("UserDetail", rsp.User)
	viper.Set("AuthResult", rsp.AuthenticationResult)
	if err := writeConfig(); err != nil {
		return err
	}
	return nil

}

// RunUserLogout logout user.
func RunUserLogout(c *CmdConfig) error {

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found")
	}

	md := metadata.New(map[string]string{
		"token": authResult.AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	userClient := usermgr.NewUserMgrClient(conn)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()
	if _, err := userClient.Logout(tokenContext,
		&usermgr.RefreshToken{RefreshToken: authResult.RefreshToken}); err != nil {
		return err
	}
	viper.Set("UserDetail", "")
	viper.Set("AuthResult", "")
	if err := writeConfig(); err != nil {
		return err
	}
	fmt.Println("Logout Success.")

	return nil
}

// RunUserConfirmRegistration confirm user registration.
func RunUserConfirmRegistration(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	confirmationCode, err := c.Ankr.GetString(c.NS, ankrctl.ArgRegisterCodeSlug)
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
	if _, err := userClient.ConfirmRegistration(context.Background(),
		&usermgr.ConfirmRegistrationRequest{Email: c.Args[0],
			ConfirmationCode: confirmationCode}); err != nil {
		return err
	}
	fmt.Println("Confirm Registration Success.")

	return nil
}

// RunUserForgotPassword send request to request new password.
func RunUserForgotPassword(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	userClient := usermgr.NewUserMgrClient(conn)
	if _, err := userClient.ForgotPassword(context.Background(),
		&usermgr.ForgotPasswordRequest{Email: c.Args[0]}); err != nil {
		return err
	}

	fmt.Println("Forgot Password Request Success.")

	return nil
}

// RunUserConfirmPassword confirm password after reset.
func RunUserConfirmPassword(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	confirmationCode, err := c.Ankr.GetString(c.NS, ankrctl.ArgPasswordCodeSlug)
	if err != nil {
		return err
	}

	confirmPassword, err := c.Ankr.GetString(c.NS, ankrctl.ArgConfirmPasswordSlug)
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
	if _, err := userClient.ConfirmPassword(context.Background(),
		&usermgr.ConfirmPasswordRequest{Email: c.Args[0], ConfirmationCode: confirmationCode,
			NewPassword: confirmPassword}); err != nil {
		return err
	}

	fmt.Println("Confirm Password Success.")

	return nil
}

// RunUserChangePassword change password with new password.
func RunUserChangePassword(c *CmdConfig) error {

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

	oldPassword, err := c.Ankr.GetString(c.NS, ankrctl.ArgOldPasswordSlug)
	if err != nil {
		return err
	}

	newPassword, err := c.Ankr.GetString(c.NS, ankrctl.ArgNewPasswordSlug)
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
	if _, err := userClient.ChangePassword(tokenctx,
		&usermgr.ChangePasswordRequest{NewPassword: newPassword, OldPassword: oldPassword}); err != nil {
		return err
	}

	fmt.Println("Change Password Success.")

	return nil
}

// RunUserTokenRefresh refresh token with new one.
func RunUserTokenRefresh(c *CmdConfig) error {

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
	userClient := usermgr.NewUserMgrClient(conn)
	rsp, err := userClient.RefreshSession(tokenctx,
		&usermgr.RefreshToken{RefreshToken: authResult.RefreshToken})
	if err != nil {
		return err
	}
	viper.Set("AuthResult", rsp)
	if err := writeConfig(); err != nil {
		return err
	}
	fmt.Println("Refresh Session Success.")

	return nil
}

// RunUserChangeEmail change password with new password.
func RunUserChangeEmail(c *CmdConfig) error {

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)
	user := usermgr.User{}
	viper.UnmarshalKey("User", &user)
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

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	userClient := usermgr.NewUserMgrClient(conn)
	if _, err := userClient.ChangeEmail(tokenctx,
		&usermgr.ChangeEmailRequest{NewEmail: c.Args[0]}); err != nil {
		return err
	}
	user.Email = c.Args[0]
	viper.Set("User", user)
	if err := writeConfig(); err != nil {
		return err
	}

	fmt.Println("Change Email Success.")

	return nil
}

// RunUserConfirmEmail confirm user registration.
func RunUserConfirmEmail(c *CmdConfig) error {

	authResult := usermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)
	user := usermgr.User{}
	viper.UnmarshalKey("User", &user)
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

	confirmationCode, err := c.Ankr.GetString(c.NS, ankrctl.ArgEmailCodeSlug)
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
	if _, err := userClient.ConfirmEmail(tokenctx,
		&usermgr.ConfirmEmailRequest{NewEmail: c.Args[0],
			ConfirmationCode: confirmationCode}); err != nil {
		return err
	}
	fmt.Println("Email Change Confirm Success.")

	return nil
}

// RunUserUpdate update user attribute.
func RunUserUpdate(c *CmdConfig) error {

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

	updateKey, err := c.Ankr.GetString(c.NS, ankrctl.ArgUpdateKeySlug)
	if err != nil {
		return err
	}

	updateValue, err := c.Ankr.GetString(c.NS, ankrctl.ArgUpdateValueSlug)
	if err != nil {
		return err
	}

	taskarray := []*usermgr.UserAttribute{}
	task := &usermgr.UserAttribute{}

	switch updateKey {
	case "name":
		task.Key = "Name"
		task.Value = &usermgr.UserAttribute_StringValue{StringValue: updateValue}
	case "pubkey":
		task.Key = "PubKey"
		task.Value = &usermgr.UserAttribute_StringValue{StringValue: updateValue}
	default:
		return fmt.Errorf("not correct user attribute for update")
	}

	taskarray = append(taskarray, task)

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	userClient := usermgr.NewUserMgrClient(conn)
	rsp, err := userClient.UpdateAttributes(tokenctx,
		&usermgr.UpdateAttributesRequest{UserAttributes: taskarray})
	if err != nil {
		return err
	}

	viper.Set("User", rsp)
	if err := writeConfig(); err != nil {
		return err
	}

	fmt.Println("User Update Attribute Success.")

	return nil
}
