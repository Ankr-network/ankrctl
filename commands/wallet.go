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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"strings"
	"syscall"
	"time"

	ankrctl "github.com/Ankr-network/dccn-cli"
	ankr_const "github.com/Ankr-network/dccn-common"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	gwusermgr "github.com/Ankr-network/dccn-common/protos/gateway/usermgr/v1"
	wallet "github.com/Ankr-network/dccn-common/wallet"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var tendermintURL string

// walletCmd creates the wallet command.
func walletCmd() *Command {
	//DCCN-CLI wallet
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "wallet",
			Aliases: []string{"w"},
			Short:   "wallet commands",
			Long:    "wallet is used to access wallet commands",
		},
		DocCategories: []string{"wallet"},
		IsIndex:       true,
	}

	//DCCN-CLI wallet genkey
	cmdWalletGenkey := CmdBuilder(cmd, RunWalletGenkey, "genkey", "generate key pair for Mainnet",
		Writer, aliasOpt("gk"), docCategories("wallet"))
	_ = cmdWalletGenkey

	//DCCN-CLI wallet importkey
	cmdWalletImportkey := CmdBuilder(cmd, RunWalletImportkey, "importkey <keystore>",
		"import key from keystore", Writer, aliasOpt("ik"), docCategories("wallet"))
	_ = cmdWalletImportkey

	/*DCCN-CLI wallet exportkey
	cmdWalletExportkey := CmdBuilder(cmd, RunWalletExportkey, "exportkey <filename>",
		"export key to keystore", Writer, aliasOpt("ek"), docCategories("wallet"))
	_ = cmdWalletExportkey
	*/
	//DCCN-CLI wallet send coins
	cmdWalletSendCoins := CmdBuilder(cmd, RunWalletSendCoins, "sendcoins <coins-amount>",
		"send token to address", Writer, aliasOpt("st"), docCategories("wallet"))
	AddStringFlag(cmdWalletSendCoins, ankrctl.ArgTargetAddressSlug, "", "", "send token to wallet address",
		requiredOpt())
	AddStringFlag(cmdWalletSendCoins, ankrctl.ArgKeyStoreSlug, "", "", "wallet keystore", requiredOpt())

	//DCCN-CLI wallet get balance
	cmdWalletGetbalance := CmdBuilder(cmd, RunWalletGetbalance, "getbalance <address>",
		"get balance of wallet by address", Writer, aliasOpt("gb"), docCategories("wallet"))
	_ = cmdWalletGetbalance

	//DCCN-CLI wallet generate erc address
	cmdWalletGenAddress := CmdBuilder(cmd, RunWalletGenAddress, "genaddr",
		"generate wallet address for deposit and withdraw", Writer, aliasOpt("ga"), docCategories("wallet"))
	AddStringFlag(cmdWalletGenAddress, ankrctl.ArgAddressTypeSlug, "", "", "wallet address type (MAINNET/ERC20/BEP2)", requiredOpt())
	AddStringFlag(cmdWalletGenAddress, ankrctl.ArgAddressPurposeSlug, "", "", "wallet address purpose (MAINNET/ERC20/BEP2)", requiredOpt())

	//DCCN-CLI wallet search deposit in a period
	cmdWalletSearchDeposit := CmdBuilder(cmd, RunWalletSearchDeposit, "search",
		"wallet search deposit in a period", Writer, aliasOpt("sd"), docCategories("wallet"))
	AddStringFlag(cmdWalletSearchDeposit, ankrctl.ArgSearchDepositStartSlug, "", "", "wallet search deposit start date (format: `mm/dd/yyyy`)", requiredOpt())
	AddStringFlag(cmdWalletSearchDeposit, ankrctl.ArgSearchDepositEndSlug, "", "", "wallet address deposit end date (format: `mm/dd/yyyy`)", requiredOpt())

	//DCCN-CLI wallet get deposit history
	cmdWalletDepositHistory := CmdBuilder(cmd, RunWalletDepositHistory, "history",
		"retrieve wallet deposit history", Writer, aliasOpt("dh"), docCategories("wallet"))
	_ = cmdWalletDepositHistory

	return cmd

}

// RunWalletGenkey generate wallet key.
func RunWalletGenkey(c *CmdConfig) error {

	if AskForConfirm(fmt.Sprintf(`please record and backup keystore once it is generated, we don’t store your private key! 
	 type 'yes' to confirm that you understood the result of this action: `)) == nil {

		fmt.Println("\ngenerating keys...\n")

		privateKey, publicKey, address := wallet.GenerateKeys()

		if privateKey == "" || publicKey == "" || address == "" {
			return fmt.Errorf("generated keys error: empty secrets")
		}

		fmt.Println("private key: ", privateKey, "\npublic key: ", publicKey, "\naddress: ", address)

		fmt.Print("\nabout to export to keystore..\nplease input the keystore encryption password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil
		}

		fmt.Print("\nplease input password again: ")
		confirmPassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		if string(password) != string(confirmPassword) {
			return errors.New("\npassword and confirm password not match")
		}

		cryptoStruct, err := EncryptDataV3([]byte(privateKey), []byte(password), StandardScryptN, StandardScryptP)
		if err != nil {
			return err
		}

		encryptedKeyJSONV3 := encryptedKeyJSONV3{
			address,
			publicKey,
			cryptoStruct,
			keyJSONVersion,
		}

		jsonKey, err := json.Marshal(encryptedKeyJSONV3)
		if err != nil {
			return err
		}

		fmt.Println("\n\nexporting to keystore...")

		ts := time.Now().UTC()

		kfw, err := KeyFileWriter(fmt.Sprintf("UTC--%s--%s", toISO8601(ts), address))
		if err != nil {
			return err
		}

		defer kfw.Close()

		_, err = kfw.Write(jsonKey)
		if err != nil {
			return errors.New("unable to write keystore")
		}

		fmt.Printf("\nkeystore: UTC--%s--%s\n\n", toISO8601(ts), address)

	}

	return nil
}

// RunWalletImportkey import wallet key.
func RunWalletImportkey(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	authResult := gwusermgr.AuthenticationResult{}
	viper.UnmarshalKey("AuthResult", &authResult)

	if authResult.AccessToken == "" {
		return fmt.Errorf("no ankr network access token found, cannot update key to ankr network")
	}

	if AskForConfirm(fmt.Sprintf(`if you have previously generated or imported wallet address, the former ones will be replaced!
	 type 'yes' to confirm that you understood the result of this action: `)) == nil {

		kf, err := ioutil.ReadFile(c.Args[0])
		if err != nil {
			return err
		}

		var key encryptedKeyJSONV3

		err = json.Unmarshal(kf, &key)
		if err != nil {
			return err
		}
		fmt.Println("\nimporting...\nwallet address: ", key.Address)

		fmt.Println("\nupdating wallet...")

		viper.Set(ankrctl.ArgPublicKeySlug, key.PublicKey)
		viper.Set(ankrctl.ArgAddressSlug, key.Address)
		if err := writeConfig(); err != nil {
			return fmt.Errorf(err.Error())
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
		userAttributes := []*gwusermgr.UserAttribute{}
		attribute := &gwusermgr.UserAttribute{Key: "PubKey", Value: key.Address}
		userAttributes = append(userAttributes, attribute)
		rsp, err := userClient.UpdateAttributes(tokenctx,
			&gwusermgr.UpdateAttributesRequest{UserAttributes: userAttributes})
		if err != nil {
			return err
		}
		fmt.Printf("updated user %s with wallet address: %s\n", rsp.Email, rsp.Attributes.PubKey)

	}

	return nil
}

/* RunWalletExportkey export wallet key.
func RunWalletExportkey(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	if AskForConfirm(fmt.Sprintf(`About to export privateKey/publicKey/address to keystore.
	Type 'yes' to confirm that you would save this key file: `)) == nil {

		key := encryptedKeyJSONV3{}
		key.PublicKey = viper.GetString(ankrctl.ArgPublicKeySlug)
		key.Address = viper.GetString(ankrctl.ArgAddressSlug)

		if key.PrivateKey == "" || key.PublicKey == "" || key.Address == "" {
			return errors.New("No existing key to export")
		}

		fmt.Print("Please input the key file encryption secret: ")
		secret, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil
		}

		fmt.Print("\nPlease input passcode again to confirm: ")
		confirmSecret, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		if string(secret) != string(confirmSecret) {
			return errors.New("\nSecret and confirm secret not match")
		}

		fmt.Println("\n\nexporting keystore...")

		kfw, err := KeyFileWriter(c.Args[0])
		if err != nil {
			return err
		}

		defer kfw.Close()

		kf, err := json.Marshal(key)
		if err != nil {
			return err
		}

		ekf, err := AesEncrypt(kf, secret)
		if err != nil {
			return errors.New("unable to encrypt key file")
		}

		_, err = kfw.Write(ekf)
		if err != nil {
			return errors.New("unable to write key file")
		}

		fmt.Println("\nDone.")

	}

	return nil
}
*/

// RunWalletSendtoken send token to other wallet address.
func RunWalletSendCoins(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	target, err := c.Ankr.GetString(c.NS, ankrctl.ArgTargetAddressSlug)
	if err != nil {
		return err
	}

	amount := c.Args[0]
	amountInt := strings.Split(amount, ".")[0]
	lenPow := 18
	if len(amountInt) > 10 {
		return errors.New("Amount range should be within 10^10")
	}

	if strings.Contains(amount, ".") {
		if len(strings.Split(amount, ".")) > 2 {
			return errors.New("Invalid amount format")
		}
		amountDecimal := strings.Split(amount, ".")[1]
		if len(amountDecimal) > 18 {
			return errors.New("Amount should retain less than 18 digits after decimal point")
		}
		amountInt = amountInt + amountDecimal
		lenPow = lenPow - len(amountDecimal)
	}
	tokenAmount, ok := new(big.Int).SetString(amountInt, 10)
	if !ok {
		return fmt.Errorf("Parsing amount format error: %s", amountInt)
	}
	tokenAmount = tokenAmount.Mul(tokenAmount, big.NewInt(int64(math.Pow10(lenPow))))

	keystore, err := c.Ankr.GetString(c.NS, ankrctl.ArgKeyStoreSlug)
	if err != nil {
		return err
	}
	ks, err := ioutil.ReadFile(keystore)
	if err != nil {
		return err
	}

	var key encryptedKeyJSONV3

	err = json.Unmarshal(ks, &key)
	if err != nil {
		return err
	}
	address := key.Address
	publicKey := key.PublicKey

	fmt.Print("Please input the keystore password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil
	}
	privateKeyDecrypt, err := DecryptDataV3(key.Crypto, string(password))
	if err != nil {
		return err
	}
	privateKey := string(privateKeyDecrypt)

	if address == "" || publicKey == "" || privateKey == "" {

		fmt.Println("\nNo key found, please approve this transaction with your wallet address, public key and private key!")

		fmt.Print("\nAddress: ")
		address, err = retrieveUserInput()
		if err != nil {
			return err
		}

		fmt.Print("\nPublic key: ")
		publicKey, err = retrieveUserInput()
		if err != nil {
			return err
		}

		fmt.Print("\nPrivate key: ")
		privateKey, err = retrieveUserInput()
		if err != nil {
			return err
		}
	}

	if address == "" || publicKey == "" || privateKey == "" {
		return errors.New("Wrong wallet address, public key and private key")
	}
	if tendermintURL == "" {
		tendermintURL = "chain-dev.dccn.ankr.com"
	}
	fmt.Println("")
	if AskForConfirm(fmt.Sprintf("About to send %s from wallet address %s to wallet address %s, Type 'yes' to confirm this action: ", c.Args[0], address, target)) == nil {
		txhash, err := wallet.SendCoins(tendermintURL, "26657", privateKey, address, target, tokenAmount.String())
		if err != nil {
			return err
		}
		fmt.Printf("\ndone. tx hash: %s\n", txhash)
	} else {
		fmt.Println("\nabort.")
	}

	return nil
}

// RunWalletGetbalance get balance from chain.
func RunWalletGetbalance(c *CmdConfig) error {

	address := viper.GetString(ankrctl.ArgAddressSlug)
	if len(c.Args) > 0 {
		address = c.Args[0]
	}

	if address == "" && len(c.Args) < 1 {
		return ankrctl.NewMissingArgsErr(c.NS)
	}

	fmt.Printf("Query balance by address %s\n", address)
	if tendermintURL == "" {
		tendermintURL = "chain-dev.dccn.ankr.com"
	}
	balance, err := wallet.GetBalance(tendermintURL, "26657", address)
	if err != nil {
		return err
	}
	balanceDecimal := balance[:len(balance)-18] + "." + balance[len(balance)-18:]
	fmt.Printf("The balance is: %s\n", balanceDecimal)

	return nil
}

// RunWalletGenAddress generate wallet key for deposit/withdraw.
func RunWalletGenAddress(c *CmdConfig) error {

	s := map[string]bool{"MAINNET": true, "ERC20": true, "BEP2": true}

	addressType, err := c.Ankr.GetString(c.NS, ankrctl.ArgAddressTypeSlug)
	if err != nil {
		return err
	}

	addressPurpose, err := c.Ankr.GetString(c.NS, ankrctl.ArgAddressPurposeSlug)
	if err != nil {
		return err
	}

	_, typeOk := s[addressType]
	_, purposeOk := s[addressPurpose]

	if !typeOk || !purposeOk {
		return fmt.Errorf("type or purpose not one of MAINNET/ERC20/BEP2..")
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
	userClient := gwusermgr.NewUserMgrClient(conn)

	rsp, err := userClient.CreateAddress(tokenctx,
		&gwusermgr.GenerateAddressRequest{
			Type:    addressType,
			Purpose: addressPurpose,
		})
	if err != nil {
		return err
	}

	fmt.Printf("Generated Address Type %s %s for Purpose %s \n",
		addressType, rsp.Typeaddress, addressPurpose)

	return nil
}

// RunWalletSearchDeposit search deposit for certain period.
func RunWalletSearchDeposit(c *CmdConfig) error {

	start, err := c.Ankr.GetString(c.NS, ankrctl.ArgSearchDepositStartSlug)
	if err != nil {
		return err
	}
	startTime, err := time.Parse("01/02/2006", start)
	if err != nil {
		return err
	}
	end, err := c.Ankr.GetString(c.NS, ankrctl.ArgSearchDepositEndSlug)
	if err != nil {
		return err
	}
	endTime, err := time.Parse("01/02/2006", end)
	if err != nil {
		return err
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
	userClient := gwusermgr.NewUserMgrClient(conn)

	rsp, err := userClient.SearchDeposit(tokenctx,
		&gwusermgr.SearchDepositRequest{
			Start: &timestamp.Timestamp{
				Seconds: startTime.Unix(),
			},
			End: &timestamp.Timestamp{
				Seconds: endTime.Unix(),
			},
		})
	if err != nil {
		return err
	}

	for _, v := range rsp.Deposits {
		amount, ok := new(big.Float).SetString(v.Amount)
		if !ok {
			return fmt.Errorf("Parsing amount format error: %s", v.Amount)
		}
		fmt.Printf("Time: %s\nHash: %s\nState: %s\nConfirmed Block Height: %s\nFrom Account Address Type: %s\nFrom Account Address: %s\nTo Account Address Type: %s\nTo Account Address: %s\nAmount: %v\n\n",
			ptypes.TimestampString(v.Time), v.TxHash, v.TxState, v.ConfirmedBlockHeight, v.FromAccountAddressType, v.FromAccountAddress, v.ToAccountAddressType, v.ToAccountAddress, new(big.Float).Quo(amount, big.NewFloat(float64(1000000000000000000.0))).String())
	}

	return nil
}

// RunWalletDepositHistory return deposit history for certain period.
func RunWalletDepositHistory(c *CmdConfig) error {

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

	rsp, err := userClient.DepositHistory(tokenctx, &common_proto.Empty{})
	if err != nil {
		return err
	}

	for _, v := range rsp.Deposits {
		amount, ok := new(big.Float).SetString(v.Amount)
		if !ok {
			return fmt.Errorf("Parsing amount format error: %s", v.Amount)
		}
		fmt.Printf("Time: %s\nHash: %s\nState: %s\nConfirmed Block Height: %s\nFrom Account Address Type: %s\nFrom Account Address: %s\nTo Account Address Type: %s\nTo Account Address: %s\nAmount: %v\n\n",
			ptypes.TimestampString(v.Time), v.TxHash, v.TxState, v.ConfirmedBlockHeight, v.FromAccountAddressType, v.FromAccountAddress, v.ToAccountAddressType, v.ToAccountAddress, new(big.Float).Quo(amount, big.NewFloat(float64(1000000000000000000.0))).String())
	}

	return nil
}
