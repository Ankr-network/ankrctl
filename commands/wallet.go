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
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	akrctl "github.com/Ankr-network/dccn-cli"
	wallet "github.com/Ankr-network/dccn-common/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tendermintURL string

type Key struct {
	PrivateKey string `json:"privatekey"`
	PublicKey  string `json:"publickey"`
	Address    string `json:"address"`
}

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
	cmdWalletGenkey := CmdBuilder(cmd, RunWalletGenkey, "genkey", "generate key pair on chain",
		Writer, aliasOpt("gk"), docCategories("wallet"))
	_ = cmdWalletGenkey

	//DCCN-CLI wallet importkey
	cmdWalletImportkey := CmdBuilder(cmd, RunWalletImportkey, "importkey <filename>",
		"import key from key file", Writer, aliasOpt("ik"), docCategories("wallet"))
	_ = cmdWalletImportkey

	//DCCN-CLI wallet exportkey
	cmdWalletExportkey := CmdBuilder(cmd, RunWalletExportkey, "exportkey <filename>",
		"export key to key file", Writer, aliasOpt("ek"), docCategories("wallet"))
	_ = cmdWalletExportkey

	//DCCN-CLI wallet send token
	cmdWalletSendtoken := CmdBuilder(cmd, RunWalletSendtoken, "sendtoken <token-amount>",
		"send token to address", Writer, aliasOpt("st"), docCategories("wallet"))
	AddStringFlag(cmdWalletSendtoken, akrctl.ArgTargetSlug, "", "", "send token to wallet address",
		requiredOpt())
	AddStringFlag(cmdWalletSendtoken, akrctl.ArgPublicKeySlug, "", "", "wallet public key")
	AddStringFlag(cmdWalletSendtoken, akrctl.ArgPrivateKeySlug, "", "", "wallet private key")
	AddStringFlag(cmdWalletSendtoken, akrctl.ArgAddressSlug, "", "", "wallet address")

	//DCCN-CLI wallet get balance
	cmdWalletGetbalance := CmdBuilder(cmd, RunWalletGetbalance, "getbal <address>",
		"get balance of wallet by address", Writer, aliasOpt("gb"), docCategories("wallet"))
	_ = cmdWalletGetbalance

	return cmd

}

// RunWalletGenkey generate wallet key.
func RunWalletGenkey(c *CmdConfig) error {

	if AskForConfirm(fmt.Sprintf(`About to generate wallet address, public key and private key..
	Please record and backup wallet address and keys once generated!! 
	Note: If these keys lost, you will lost access to your tokens in the wallet!!
	Note: If you have previously generated these keys, the former ones will be replaced!!  
	
	Type 'yes' to confirm that you understood the result of this action: `)) == nil {

		fmt.Println("\n\nGenerating keys...\n")

		privateKey, publicKey, address := wallet.GenerateKeys()

		if privateKey == "" || publicKey == "" || address == "" {
			return fmt.Errorf("Generated keys error, empty secrets..")
		}

		fmt.Println("Updating wallet...\n")

		viper.Set(akrctl.ArgPrivateKeySlug, privateKey)
		viper.Set(akrctl.ArgPublicKeySlug, publicKey)
		viper.Set(akrctl.ArgAddressSlug, address)
		if err := writeConfig(); err != nil {
			return fmt.Errorf(err.Error())
		}
		// Todo: update user with publicKey/address

		fmt.Println("Private Key: ", privateKey, "\nPublic Key: ", publicKey, "\nAddress: ", address)
	}

	return nil
}

// RunWalletImportkey import wallet key.
func RunWalletImportkey(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	if AskForConfirm(fmt.Sprintf(`About to import address, public key and private key from key file.
	Note: If you have previously generated or imported these keys, the former ones will be replaced!
	Type 'yes' to confirm that you understood the result of this action: `)) == nil {

		kf, err := ioutil.ReadFile(c.Args[0])
		if err != nil {
			return err
		}

		fmt.Print("\nPlease input the keyfile secret: ")
		secret, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil
		}

		dkf, err := AesDecrypt(kf, []byte(secret))
		if err != nil {
			return err
		}

		var key Key

		err = json.Unmarshal(dkf, &key)
		if err != nil {
			return err
		}
		fmt.Println("\nImporting...\nPrivate Key: ", key.PrivateKey,
			"\nPublic Key: ", key.PublicKey, "\nAddress: ", key.Address)

		fmt.Println("\nUpdating wallet...")

		viper.Set(akrctl.ArgPrivateKeySlug, key.PrivateKey)
		viper.Set(akrctl.ArgPublicKeySlug, key.PublicKey)
		viper.Set(akrctl.ArgAddressSlug, key.Address)
		if err := writeConfig(); err != nil {
			return fmt.Errorf(err.Error())
		}

		// Todo: update user with publicKey/address

		fmt.Println("\nDone.")
	}

	return nil
}

// RunWalletExportkey export wallet key.
func RunWalletExportkey(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	if AskForConfirm(fmt.Sprintf(`About to export privateKey/publicKey/address to key file.
	Type 'yes' to confirm that you would save this key file: `)) == nil {

		key := Key{}
		key.PrivateKey = viper.GetString(akrctl.ArgPrivateKeySlug)
		key.PublicKey = viper.GetString(akrctl.ArgPublicKeySlug)
		key.Address = viper.GetString(akrctl.ArgAddressSlug)

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

		fmt.Println("\n\nExporting keys...")

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

// RunWalletSendtoken send token to other wallet address.
func RunWalletSendtoken(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	target, err := c.Ankr.GetString(c.NS, akrctl.ArgTargetSlug)
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
	tokenAmount, _ := new(big.Int).SetString(amountInt, 10)
	tokenAmount = tokenAmount.Mul(tokenAmount, big.NewInt(int64(math.Pow10(lenPow))))
	address, err := c.Ankr.GetString(c.NS, akrctl.ArgAddressSlug)
	if err != nil {
		return err
	}

	publicKey, err := c.Ankr.GetString(c.NS, akrctl.ArgPublicKeySlug)
	if err != nil {
		return err
	}

	privateKey, err := c.Ankr.GetString(c.NS, akrctl.ArgPrivateKeySlug)
	if err != nil {
		return err
	}

	if address == "" || publicKey == "" || privateKey == "" {

		address = viper.GetString(akrctl.ArgAddressSlug)
		publicKey = viper.GetString(akrctl.ArgPublicKeySlug)
		privateKey = viper.GetString(akrctl.ArgPrivateKeySlug)

		if address == "" || publicKey == "" || privateKey == "" {

			fmt.Println("\nPlease approve this transaction with your wallet address, public key and private key!")

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
	}

	if address == "" || publicKey == "" || privateKey == "" {
		return errors.New("Wrong wallet address, public key and private key")
	}
	if tendermintURL == "" {
		tendermintURL = "chain-dev.dccn.ankr.network"
	}
	if AskForConfirm(fmt.Sprintf("About to send %s from wallet address %s to wallet address %s, Type 'yes' to confirm this action: ", tokenAmount.String(), address, target)) == nil {
		if err := wallet.SendCoins(tendermintURL, "26657", privateKey, address, target, tokenAmount.String(), publicKey); err != nil {
			return err
		}
		fmt.Println("\nDone.")
	} else {
		fmt.Println("\nAbort.")
	}

	return nil
}

// RunWalletGetbalance get balance from chain.
func RunWalletGetbalance(c *CmdConfig) error {

	address := viper.GetString(akrctl.ArgAddressSlug)
	if len(c.Args) > 0 {
		address = c.Args[0]
	}

	if address == "" && len(c.Args) < 1 {
		return akrctl.NewMissingArgsErr(c.NS)
	}

	fmt.Printf("Query balance by address %s\n", address)
	if tendermintURL == "" {
		tendermintURL = "chain-dev.dccn.ankr.network"
	}
	balance, err := wallet.GetBalance(tendermintURL, "26657", address)
	if err != nil {
		return err
	}

	fmt.Printf("The balance is: %s\n", balance)

	return nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	keyPatch := []byte(fmt.Sprintf("%16s", string(key)))
	block, err := aes.NewCipher(keyPatch)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, keyPatch[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	keyPatch := []byte(fmt.Sprintf("%16s", string(key)))
	block, err := aes.NewCipher(keyPatch)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyPatch[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func KeyFileWriter(keyFile string) (io.WriteCloser, error) {
	f, err := os.Create(keyFile)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(keyFile, 0600); err != nil {
		return nil, err
	}

	return f, nil
}
