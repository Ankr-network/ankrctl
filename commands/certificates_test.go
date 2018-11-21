package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Ankr-network/dccn-cli"
	"github.com/Ankr-network/dccn-cli/do"
	"github.com/Ankr-network/godo"

	"github.com/stretchr/testify/assert"
)

var (
	testCertificate = do.Certificate{
		Certificate: &godo.Certificate{
			ID:              "892071a0-bb95-49bc-8021-3afd67a210bf",
			Name:            "web-cert-01",
			NotAfter:        "2017-02-22T00:23:00Z",
			SHA1Fingerprint: "dfcc9f57d86bf58e321c2c6c31c7a971be244ac7",
			Created:         "2017-02-08T16:02:37Z",
		},
	}

	testCertificateList = do.Certificates{testCertificate}
)

func TestCertificateCommand(t *testing.T) {
	cmd := Certificate()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "get", "create", "list", "delete")
}

func TestCertificateGetNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunCertificateGet(config)
		assert.Error(t, err)
	})
}

func TestCertificateGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		cID := "892071a0-bb95-49bc-8021-3afd67a210bf"
		tm.certificates.On("Get", cID).Return(&testCertificate, nil)

		config.Args = append(config.Args, cID)

		err := RunCertificateGet(config)
		assert.NoError(t, err)
	})
}

func TestCertificatesCreate(t *testing.T) {
	tests := []struct {
		desc           string
		certName       string
		DNSNames       []string
		privateKeyPath string
		certLeafPath   string
		certChainPath  string
		certType       string
		certificate    godo.CertificateRequest
	}{
		{
			desc:           "creates custom certificate",
			certName:       "custom-cert",
			privateKeyPath: filepath.Join(os.TempDir(), "cer-key.crt"),
			certLeafPath:   filepath.Join(os.TempDir(), "leaf-cer.crt"),
			certChainPath:  filepath.Join(os.TempDir(), "cer-chain.crt"),
			certificate: godo.CertificateRequest{
				Name:             "custom-cert",
				PrivateKey:       "-----BEGIN PRIVATE KEY-----",
				LeafCertificate:  "-----BEGIN CERTIFICATE-----",
				CertificateChain: "-----BEGIN CERTIFICATE-----",
			},
		},

		{
			desc:           "creates custom cerftificate without specifying chain",
			certName:       "cert-without-chain",
			privateKeyPath: filepath.Join(os.TempDir(), "cer-key.crt"),
			certLeafPath:   filepath.Join(os.TempDir(), "leaf-cer.crt"),
			certificate: godo.CertificateRequest{
				Name:             "cert-without-chain",
				PrivateKey:       "-----BEGIN PRIVATE KEY-----",
				LeafCertificate:  "-----BEGIN CERTIFICATE-----",
				CertificateChain: "",
			},
		},
		{
			desc:     "creates lets_encrypt cerftificate",
			certName: "lets-encrypt-cert",
			DNSNames: []string{"sampledomain.org", "api.sampledomain.org"},
			certType: "lets_encrypt",
			certificate: godo.CertificateRequest{
				Name:     "lets-encrypt-cert",
				DNSNames: []string{"sampledomain.org", "api.sampledomain.org"},
				Type:     "lets_encrypt",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
				if test.privateKeyPath != "" {
					pkErr := ioutil.WriteFile(test.privateKeyPath, []byte("-----BEGIN PRIVATE KEY-----"), 0600)
					assert.NoError(t, pkErr)

					defer os.Remove(test.privateKeyPath)
				}

				if test.certLeafPath != "" {
					certErr := ioutil.WriteFile(test.certLeafPath, []byte("-----BEGIN CERTIFICATE-----"), 0600)
					assert.NoError(t, certErr)

					defer os.Remove(test.certLeafPath)
				}

				if test.certChainPath != "" {
					certErr := ioutil.WriteFile(test.certChainPath, []byte("-----BEGIN CERTIFICATE-----"), 0600)
					assert.NoError(t, certErr)

					defer os.Remove(test.certChainPath)
				}

				tm.certificates.On("Create", &test.certificate).Return(&testCertificate, nil)

				config.Ankr.Set(config.NS, dccncli.ArgCertificateName, test.certName)

				if test.privateKeyPath != "" {
					config.Ankr.Set(config.NS, dccncli.ArgPrivateKeyPath, test.privateKeyPath)
				}

				if test.certLeafPath != "" {
					config.Ankr.Set(config.NS, dccncli.ArgLeafCertificatePath, test.certLeafPath)
				}

				if test.certChainPath != "" {
					config.Ankr.Set(config.NS, dccncli.ArgCertificateChainPath, test.certChainPath)
				}

				if test.DNSNames != nil {
					config.Ankr.Set(config.NS, dccncli.ArgCertificateDNSNames, test.DNSNames)
				}

				if test.certType != "" {
					config.Ankr.Set(config.NS, dccncli.ArgCertificateType, test.certType)
				}

				err := RunCertificateCreate(config)
				assert.NoError(t, err)
			})
		})
	}
}

func TestCertificateList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.certificates.On("List").Return(testCertificateList, nil)

		err := RunCertificateList(config)
		assert.NoError(t, err)
	})
}

func TestCertificateDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		cID := "892071a0-bb95-49bc-8021-3afd67a210bf"
		tm.certificates.On("Delete", cID).Return(nil)

		config.Args = append(config.Args, cID)
		config.Ankr.Set(config.NS, dccncli.ArgForce, true)

		err := RunCertificateDelete(config)
		assert.NoError(t, err)
	})
}

func TestCertificateDeleteNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunCertificateDelete(config)
		assert.Error(t, err)
	})
}
