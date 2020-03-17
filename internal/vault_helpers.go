package vaulthelpers

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/pkg/errors"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/certutil"

	// helpers
	"github.com/cloudflare/cfssl/certinfo"
)

// NewVaultClient creates a new Vault API client for Vault
func NewVaultClient(vaultAddr, vaultToken string, client *http.Client) (*api.Client, error) {

	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	vaultClient, err := api.NewClient(&api.Config{Address: vaultAddr, HttpClient: client})
	if err != nil {
		return nil, err
	}

	vaultClient.SetToken(vaultToken)

	if namespace := os.Getenv("VAULT_NAMESPACE"); namespace != "" {
		vaultClient.SetNamespace(namespace)
	}

	return vaultClient, nil
}

func getArrayOfCertsFromVault(client *api.Client, secret *api.Secret, vaultAddr string, pkiPath string, serial bool) (arrayOfCerts []*certinfo.Certificate, err error) {

	var certArray = []*certinfo.Certificate{}

	for _, key := range secret.Data["keys"].([]interface{}) {
		secret, err := client.Logical().Read(fmt.Sprintf("%s/cert/%s", pkiPath, key))
		if err != nil {
			return nil, err
		}

		certParse, err := parseCertFromVaultSecret(secret)

		if err != nil {
			return nil, err
		}

		if serial {
			bignum, _ := new(big.Int).SetString(certParse.SerialNumber, 0)
			convertedSerial := certutil.GetHexFormatted(bignum.Bytes(), ":")
			reflect.ValueOf(certParse).Elem().FieldByName("SerialNumber").SetString(convertedSerial)
		}

		certArray = append(certArray, certParse)

	}

	return certArray, err
}

func parseCertFromVaultSecret(secret *api.Secret) (*certinfo.Certificate, error) {
	rawCert := secret.Data["certificate"].(string)
	block, _ := pem.Decode([]byte(rawCert))
	if block == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse certificate")
	}

	certParse := certinfo.ParseCertificate(cert)

	return certParse, nil
}
