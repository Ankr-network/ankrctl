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

	"github.com/spf13/cobra"

	"context"

	"github.com/Ankr-network/ankrctl/types"
	ankr_const "github.com/Ankr-network/dccn-common"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	gwusermgr "github.com/Ankr-network/dccn-common/protos/gateway/usermgr/v1"
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
	AddStringFlag(cmdUserRegister, types.ArgEmailSlug, "", "", "User email", requiredOpt())
	AddStringFlag(cmdUserRegister, types.ArgPasswordSlug, "", "", "User password", requiredOpt())

	//DCCN-CLI user comfirm registration
	cmdUserConfirmRegistration := CmdBuilder(cmd, RunUserConfirmRegistration,
		"confirm-registration <user-email>", "user registration confirmation", Writer,
		aliasOpt("rc"), docCategories("user"))
	AddStringFlag(cmdUserConfirmRegistration, types.ArgRegisterCodeSlug,
		"", "", "User registration confirmation code", requiredOpt())

	//DCCN-CLI user forgot password
	cmdUserForgotPassword := CmdBuilder(cmd, RunUserForgotPassword, "forgot-password <user-email>",
		"user password forgot", Writer, aliasOpt("fp"), docCategories("user"))
	_ = cmdUserForgotPassword

	//DCCN-CLI user comfirm password
	cmdUserConfirmPassword := CmdBuilder(cmd, RunUserConfirmPassword, "confirm-password <user-email>",
		"user password change confirmation", Writer, aliasOpt("pc"), docCategories("user"))
	AddStringFlag(cmdUserConfirmPassword, types.ArgPasswordCodeSlug,
		"", "", "User password change confirmation code", requiredOpt())
	AddStringFlag(cmdUserConfirmPassword, types.ArgConfirmPasswordSlug,
		"", "", "User confirm new password", requiredOpt())

	//DCCN-CLI user change password
	cmdUserChangePassword := CmdBuilder(cmd, RunUserChangePassword, "change-password <user-email>",
		"user password change", Writer, aliasOpt("cp"), docCategories("user"))
	AddStringFlag(cmdUserChangePassword, types.ArgOldPasswordSlug,
		"", "", "User old password", requiredOpt())
	AddStringFlag(cmdUserChangePassword, types.ArgNewPasswordSlug,
		"", "", "User new password", requiredOpt())

	//DCCN-CLI user change email
	cmdUserChangeEmail := CmdBuilder(cmd, RunUserChangeEmail, "email-change <new-email>",
		"user email change", Writer, aliasOpt("ec"), docCategories("user"))
	_ = cmdUserChangeEmail

	//DCCN-CLI user confirm email
	cmdUserConfirmEmail := CmdBuilder(cmd, RunUserConfirmEmail, "email-confirm <new-email>",
		"user confirm email change", Writer, aliasOpt("ce"), docCategories("user"))
	AddStringFlag(cmdUserConfirmEmail, types.ArgEmailCodeSlug,
		"", "", "User email change confirmation code", requiredOpt())

	//DCCN-CLI user update attribute
	cmdUserUpdate := CmdBuilder(cmd, RunUserUpdate, "update", "user update attribute",
		Writer, aliasOpt("ua"), docCategories("user"))
	AddStringFlag(cmdUserUpdate, types.ArgUpdateKeySlug, "", "", "User attribute key", requiredOpt())
	AddStringFlag(cmdUserUpdate, types.ArgUpdateValueSlug, "", "", "User attribute value", requiredOpt())

	//DCCN-CLI user login
	cmdUserLogin := CmdBuilder(cmd, RunUserLogin, "login", "user login", Writer,
		aliasOpt("li"), docCategories("user"))
	_ = cmdUserLogin

	//DCCN-CLI user logout
	cmdUserLogout := CmdBuilder(cmd, RunUserLogout, "logout", "user logout", Writer,
		aliasOpt("lo"), docCategories("user"))
	_ = cmdUserLogout

	//DCCN-CLI get user detail with wallet address
	cmdUserDetail := CmdBuilder(cmd, RunUserDetail, "detail",
		"get user detail with wallet address", Writer, aliasOpt("ud"), docCategories("user"))
	_ = cmdUserDetail

	return cmd

}

// RunUserRegister register a user.
func RunUserRegister(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return types.NewMissingArgsErr(c.NS)
	}

	email, err := c.Ankr.GetString(c.NS, types.ArgEmailSlug)
	if err != nil {
		fmt.Println(err)
		return err
	}

	password, err := c.Ankr.GetString(c.NS, types.ArgPasswordSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	userClient := gwusermgr.NewUserMgrClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	urr := &gwusermgr.RegisterRequest{
		Password: password,
		Email:    email,
		Name:     c.Args[0],
	}

	if _, err := userClient.Register(ctx, urr); err != nil {
		return err
	}

	fmt.Printf("User %s Register Requested, Please Check Your Email Box.\n", email)

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
	userClient := gwusermgr.NewUserMgrClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(),
		ankr_const.ClientTimeOut*time.Second)
	defer cancel()

	rsp, err := userClient.Login(ctx,
		&gwusermgr.LoginRequest{Email: email, Password: string(password)})
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

	authResult := gwusermgr.AuthenticationResult{}
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
	userClient := gwusermgr.NewUserMgrClient(conn)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()
	if _, err := userClient.Logout(tokenContext,
		&gwusermgr.RefreshToken{RefreshToken: authResult.RefreshToken}); err != nil {
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
		return types.NewMissingArgsErr(c.NS)
	}

	confirmationCode, err := c.Ankr.GetString(c.NS, types.ArgRegisterCodeSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	userClient := gwusermgr.NewUserMgrClient(conn)
	if _, err := userClient.ConfirmRegistration(context.Background(),
		&gwusermgr.ConfirmRegistrationRequest{Email: c.Args[0],
			ConfirmationCode: confirmationCode}); err != nil {
		return err
	}
	fmt.Println("Confirm Registration Success.")

	return nil
}

// RunUserForgotPassword send request to request new password.
func RunUserForgotPassword(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return types.NewMissingArgsErr(c.NS)
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	userClient := gwusermgr.NewUserMgrClient(conn)
	if _, err := userClient.ForgotPassword(context.Background(),
		&gwusermgr.ForgotPasswordRequest{Email: c.Args[0]}); err != nil {
		return err
	}

	fmt.Println("Forgot Password Requested, Please Check Your Email Box.")

	return nil
}

// RunUserConfirmPassword confirm password after reset.
func RunUserConfirmPassword(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return types.NewMissingArgsErr(c.NS)
	}

	confirmationCode, err := c.Ankr.GetString(c.NS, types.ArgPasswordCodeSlug)
	if err != nil {
		return err
	}

	confirmPassword, err := c.Ankr.GetString(c.NS, types.ArgConfirmPasswordSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	userClient := gwusermgr.NewUserMgrClient(conn)
	if _, err := userClient.ConfirmPassword(context.Background(),
		&gwusermgr.ConfirmPasswordRequest{Email: c.Args[0], ConfirmationCode: confirmationCode,
			NewPassword: confirmPassword}); err != nil {
		return err
	}

	fmt.Println("Confirm Password Success.")

	return nil
}

// RunUserChangePassword change password with new password.
func RunUserChangePassword(c *CmdConfig) error {

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

	oldPassword, err := c.Ankr.GetString(c.NS, types.ArgOldPasswordSlug)
	if err != nil {
		return err
	}

	newPassword, err := c.Ankr.GetString(c.NS, types.ArgNewPasswordSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")
	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	userClient := gwusermgr.NewUserMgrClient(conn)
	if _, err := userClient.ChangePassword(tokenctx,
		&gwusermgr.ChangePasswordRequest{NewPassword: newPassword, OldPassword: oldPassword}); err != nil {
		return err
	}

	fmt.Println("Change Password Success.")

	return nil
}

// RunUserTokenRefresh refresh token with new one.
func RunUserTokenRefresh(c *CmdConfig) error {

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
	userClient := gwusermgr.NewUserMgrClient(conn)
	rsp, err := userClient.RefreshSession(tokenctx,
		&gwusermgr.RefreshToken{RefreshToken: authResult.RefreshToken})
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

	authResult := gwusermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)
	user := gwusermgr.User{}
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
		return types.NewMissingArgsErr(c.NS)
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	userClient := gwusermgr.NewUserMgrClient(conn)
	if _, err := userClient.ChangeEmail(tokenctx,
		&gwusermgr.ChangeEmailRequest{NewEmail: c.Args[0]}); err != nil {
		return err
	}
	user.Email = c.Args[0]
	viper.Set("User", user)
	if err := writeConfig(); err != nil {
		return err
	}

	fmt.Println("Change Email Requested, Please Check Your Email Box.")

	return nil
}

// RunUserConfirmEmail confirm user registration.
func RunUserConfirmEmail(c *CmdConfig) error {

	authResult := gwusermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)
	user := gwusermgr.User{}
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
		return types.NewMissingArgsErr(c.NS)
	}

	confirmationCode, err := c.Ankr.GetString(c.NS, types.ArgEmailCodeSlug)
	if err != nil {
		return err
	}

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	userClient := gwusermgr.NewUserMgrClient(conn)
	if _, err := userClient.ConfirmEmail(tokenctx,
		&gwusermgr.ConfirmEmailRequest{NewEmail: c.Args[0],
			ConfirmationCode: confirmationCode}); err != nil {
		return err
	}
	fmt.Println("Email Change Confirm Success.")

	return nil
}

// RunUserUpdate update user attribute.
func RunUserUpdate(c *CmdConfig) error {

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

	updateKey, err := c.Ankr.GetString(c.NS, types.ArgUpdateKeySlug)
	if err != nil {
		return err
	}

	updateValue, err := c.Ankr.GetString(c.NS, types.ArgUpdateValueSlug)
	if err != nil {
		return err
	}

	attributeArray := []*gwusermgr.UserAttribute{}
	attribute := &gwusermgr.UserAttribute{}

	keys := map[string]bool{
		"AvatarBackgroundColor": true,
		"MainnetToErcAddr":      true,
		"ErcToMainnetAddr":      true,
		"MainnetToBepAddr":      true,
		"ErcToBepAddr":          true,
		"BepToErcAddr":          true,
		"BepToMainnetAddr":      true,
		"BepPubKey":             true,
		"ErcPubKey":             true,
		"BepToMainnetMemo":      true,
		"Avatar":                true,
		"Name":                  true,
		"PubKey":                true,
	}

	if _, ok := keys[updateKey]; !ok {
		return fmt.Errorf("not correct user attribute for update")
	}
	attribute.Key = updateKey
	attribute.Value = updateValue

	attributeArray = append(attributeArray, attribute)

	url := viper.GetString("hub-url")

	conn, err := grpc.Dial(url+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	userClient := gwusermgr.NewUserMgrClient(conn)
	rsp, err := userClient.UpdateAttributes(tokenctx,
		&gwusermgr.UpdateAttributesRequest{UserAttributes: attributeArray})
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

// RunUserDetail get user tail with wallet address.
func RunUserDetail(c *CmdConfig) error {

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
	userClient := gwusermgr.NewUserMgrClient(conn)
	rsp, err := userClient.UserDetail(tokenctx, &common_proto.Empty{})
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s \n", rsp.Attributes.Name)
	fmt.Printf("Email: %s \n", rsp.Email)
	fmt.Printf("Status: %s \n", rsp.Status.String())
	fmt.Printf("Creation Date: %s \n",
		time.Unix(int64(rsp.Attributes.CreationDate), 0).Format("Mon Jan 2 15:04:05 MST 2006"))
	fmt.Printf("Last Modified Date: %s \n",
		time.Unix(int64(rsp.Attributes.LastModifiedDate), 0).Format("Mon Jan 2 15:04:05 MST 2006"))
	fmt.Printf("Pubkey: %s \n", rsp.Attributes.PubKey)

	for _, a := range rsp.Attributes.ExtraFields {
		fmt.Printf("%s: %s \n", a.Key, a.Value)
	}

	return nil

}
