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
	"io/ioutil"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/commands/displayers"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/godo"

	"github.com/spf13/cobra"
)

// Certificate creates the certificate command.
func Certificate() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "certificate",
			Short: "certificate commands",
			Long:  "certificate is used to access certificate commands",
		},
	}

	CmdBuilder(cmd, RunCertificateGet, "get <id>", "get certificate", Writer, aliasOpt("g"))

	cmdCertificateCreate := CmdBuilder(cmd, RunCertificateCreate, "create", "create new certificate", Writer, aliasOpt("c"))
	AddStringFlag(cmdCertificateCreate, dccncli.ArgCertificateName, "", "", "certificate name", requiredOpt())
	AddStringSliceFlag(cmdCertificateCreate, dccncli.ArgCertificateDNSNames, "", []string{}, "comma-separated list of domain names, required for lets_encrypt certificate")
	AddStringFlag(cmdCertificateCreate, dccncli.ArgPrivateKeyPath, "", "", "path to a private key for the certificate, required for custom certificate")
	AddStringFlag(cmdCertificateCreate, dccncli.ArgLeafCertificatePath, "", "", "path to a certificate leaf, required for custom certificate")
	AddStringFlag(cmdCertificateCreate, dccncli.ArgCertificateChainPath, "", "", "path to a certificate chain")
	AddStringFlag(cmdCertificateCreate, dccncli.ArgCertificateType, "", "", "certificate type, possible values: custom or lets_encrypt")

	CmdBuilder(cmd, RunCertificateList, "list", "list certificates", Writer, aliasOpt("ls"))

	cmdCertificateDelete := CmdBuilder(cmd, RunCertificateDelete, "delete <id>", "delete certificate", Writer, aliasOpt("d", "rm"))
	AddBoolFlag(cmdCertificateDelete, dccncli.ArgForce, dccncli.ArgShortForce, false, "Force certificate delete")

	return cmd
}

// RunCertificateGet retrieves an existing certificate by its identifier.
func RunCertificateGet(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	cID := c.Args[0]

	cs := c.Certificates()
	cer, err := cs.Get(cID)
	if err != nil {
		return err
	}

	item := &displayers.Certificate{Certificates: do.Certificates{*cer}}
	return c.Display(item)
}

// RunCertificateCreate creates a certificate.
func RunCertificateCreate(c *CmdConfig) error {
	name, err := c.Ankr.GetString(c.NS, dccncli.ArgCertificateName)
	if err != nil {
		return err
	}

	domainList, err := c.Ankr.GetStringSlice(c.NS, dccncli.ArgCertificateDNSNames)
	if err != nil {
		return err
	}

	cType, err := c.Ankr.GetString(c.NS, dccncli.ArgCertificateType)
	if err != nil {
		return err
	}

	r := &godo.CertificateRequest{
		Name:     name,
		DNSNames: domainList,
		Type:     cType,
	}

	pkPath, err := c.Ankr.GetString(c.NS, dccncli.ArgPrivateKeyPath)
	if err != nil {
		return err
	}

	if len(pkPath) > 0 {
		pc, err := readInputFromFile(pkPath)
		if err != nil {
			return err
		}

		r.PrivateKey = pc
	}

	lcPath, err := c.Ankr.GetString(c.NS, dccncli.ArgLeafCertificatePath)
	if err != nil {
		return err
	}

	if len(lcPath) > 0 {
		lc, err := readInputFromFile(lcPath)
		if err != nil {
			return err
		}

		r.LeafCertificate = lc
	}

	ccPath, err := c.Ankr.GetString(c.NS, dccncli.ArgCertificateChainPath)
	if err != nil {
		return err
	}

	if len(ccPath) > 0 {
		cc, err := readInputFromFile(ccPath)
		if err != nil {
			return err
		}

		r.CertificateChain = cc
	}

	cs := c.Certificates()
	cer, err := cs.Create(r)
	if err != nil {
		return err
	}

	item := &displayers.Certificate{Certificates: do.Certificates{*cer}}
	return c.Display(item)
}

// RunCertificateList lists certificates.
func RunCertificateList(c *CmdConfig) error {
	cs := c.Certificates()
	list, err := cs.List()
	if err != nil {
		return err
	}

	item := &displayers.Certificate{Certificates: list}
	return c.Display(item)
}

// RunCertificateDelete deletes a certificate by its identifier.
func RunCertificateDelete(c *CmdConfig) error {
	if len(c.Args) != 1 {
		return dccncli.NewMissingArgsErr(c.NS)
	}
	cID := c.Args[0]

	force, err := c.Ankr.GetBool(c.NS, dccncli.ArgForce)
	if err != nil {
		return err
	}

	if force || AskForConfirm("delete this certificate") == nil {
		cs := c.Certificates()
		if err := cs.Delete(cID); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("operation aborted")
	}

	return nil
}

func readInputFromFile(path string) (string, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(fileBytes), nil
}
