package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/certutil"

	// helpers
	"github.com/cloudflare/cfssl/certinfo"
	"github.com/hokaccha/go-prettyjson"
	"github.com/olekukonko/tablewriter"

	"github.com/urfave/cli/v2"

	vltcrthlpr "github.com/petems/vault-cert-helpers"
)

// Version is what is returned by the `-v` flag
const Version = "0.1.0"

// gitCommit is the gitcommit its built from
var gitCommit = "development"

func main() {
	app := &cli.App{
		Name:    "vault-cert-info",
		Usage:   "A simple cli app to return certificates from a Vault PKI mount",
		Version: Version + "-" + gitCommit,
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all certificates",
				Action: func(c *cli.Context) error {
					err := cmdVaultListCerts(c)
					return err
				},
			},
			{
				Name:  "cert",
				Usage: "Get information from one certificate",
				Action: func(c *cli.Context) error {
					err := cmdVaultCert(c)
					return err
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "pki",
				Value: "pki",
				Usage: fmt.Sprintf("The path to your pki engine"),
			},
			&cli.BoolFlag{
				Name:  "serial",
				Value: true,
				Usage: fmt.Sprintf("Convert the cert serial from bigint to a HEX formatted string\n\t\teg. 611168959691622484330452100723265332608845077531 -> 6b:0d:c3:94:c9:e1:20:d1:9a:eb:76:66:db:3d:8a:37:23:75:dc:1b"),
			},
			&cli.StringFlag{
				Name:  "format",
				Value: "pretty_json",
				Usage: fmt.Sprintf("The format you want them returned in, valid values are: table, json, pretty_json"),
			},
			&cli.BoolFlag{
				Name:  "silent",
				Value: true,
				Usage: fmt.Sprintf("Do not output anything other than errors or returned data"),
			},
		},
	}

	cli.AppHelpTemplate = `NAME:
	{{.Name}} - {{.Usage}}
USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
	{{if len .Authors}}
AUTHOR:
	{{range .Authors}}{{ . }}{{end}}
	{{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}
VAULT CONFIGURATION:
	Vault configuration is set by the common Vault environmental variables: 
		VAULT_ADDR: The address for the Vault server (Required)
		VAULT_TOKEN: The token for the Vault server (Required) 
		VAULT_NAMESPACE: The namespace where PKI is mounted (Optional) 
	See https://www.vaultproject.io/docs/commands/#environment-variables 

GLOBAL OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
	{{.Copyright}}
	{{end}}{{if .Version}}
VERSION:
	{{.Version}}
	{{end}}
 `

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func getENV(value string) (string, error) {
	envValue := os.Getenv(value)
	if len(envValue) == 0 {
		return "", fmt.Errorf("No ENV value for %s", value)
	}
	return envValue, nil
}

// NewVaultClient creates a new Vault API client for Vault
func NewVaultClient(vaultAddr, vaultToken string) (*api.Client, error) {
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	client, err := api.NewClient(&api.Config{Address: vaultAddr, HttpClient: httpClient})
	if err != nil {
		return nil, err
	}

	client.SetToken(vaultToken)

	if namespace := os.Getenv("VAULT_NAMESPACE"); namespace != "" {
		client.SetNamespace(namespace)
	}

	return client, nil
}

func serialConvert(cert *certinfo.Certificate) {
	bignum, _ := new(big.Int).SetString(cert.SerialNumber, 0)
	convertedSerial := certutil.GetHexFormatted(bignum.Bytes(), ":")
	reflect.ValueOf(cert).Elem().FieldByName("SerialNumber").SetString(convertedSerial)
}

func tablePrint(arrayOfCerts []*certinfo.Certificate) {

	data := [][]string{}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeader([]string{"Common Name", "Organization", "Expires", "Serial"})

	utc := time.FixedZone("UTC+0", 9*60*60)

	for _, cert := range arrayOfCerts {

		data = append(data, []string{cert.Subject.CommonName, cert.Subject.Organization, cert.NotAfter.In(utc).Format(time.RFC3339), cert.SerialNumber})

		for _, v := range data {
			table.Append(v)
		}
	}
	table.Render()
}

func cmdVaultListCerts(ctx *cli.Context) (err error) {

	silent := ctx.Bool("silent")

	vaultAddr, err := getENV("VAULT_ADDR")

	if err != nil {
		return err
	}

	vaultToken, err := getENV("VAULT_TOKEN")

	if err != nil {
		return err
	}

	client, err := NewVaultClient(vaultAddr, vaultToken)
	if err != nil {
		return err
	}

	pkiPath := ctx.String("pki")
	convertSerial := ctx.Bool("serial")

	if silent {

	} else {
		fmt.Printf("Reading certs from %s/%s\n", vaultAddr, pkiPath)
	}

	secret, err := vltcrthlpr.GetListOfCerts(client, pkiPath)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	arrayOfCerts, err := vltcrthlpr.GetArrayOfCertsFromVault(client, secret, pkiPath)

	if err != nil {
		return err
	}

	var arrayOfCertInfo = []*certinfo.Certificate{}

	for _, cert := range arrayOfCerts {

		certinfoCert := certinfo.ParseCertificate(cert)

		if convertSerial {
			serialConvert(certinfoCert)
		}

		arrayOfCertInfo = append(arrayOfCertInfo, certinfoCert)
	}

	switch ctx.String("format") {
	case "json":
		certsAsMarshall, err := json.Marshal(arrayOfCertInfo)
		if err != nil {
			return err
		}
		fmt.Println(string(certsAsMarshall))
	case "pretty_json":
		s, err := prettyjson.Marshal(arrayOfCertInfo)
		if err != nil {
			return err
		}
		fmt.Println(string(s))
	case "table":
		tablePrint(arrayOfCertInfo)
	}

	return nil

}

func cmdVaultCert(ctx *cli.Context) (err error) {

	certSerial := ctx.Args().Get(0)

	if certSerial == "" {
		return fmt.Errorf("cert argument requires a serial to lookup")
	}

	silent := ctx.Bool("silent")

	vaultAddr, err := getENV("VAULT_ADDR")

	if err != nil {
		return err
	}

	vaultToken, err := getENV("VAULT_TOKEN")

	if err != nil {
		return err
	}

	client, err := NewVaultClient(vaultAddr, vaultToken)
	if err != nil {
		return err
	}

	pkiPath := ctx.String("pki")
	convertSerial := ctx.Bool("serial")

	if silent {

	} else {
		fmt.Printf("Reading cert from %s/%s\n", vaultAddr, pkiPath)
	}

	secret, err := client.Logical().Read(fmt.Sprintf("%s/cert/%s", pkiPath, certSerial))

	if err != nil {
		return err
	}

	if secret == nil {
		return fmt.Errorf("No value found for cert at '%s'", fmt.Sprintf("%s/cert/%s", pkiPath, certSerial))
	}

	certParse, err := vltcrthlpr.ParseCertFromVaultSecret(secret)

	if err != nil {
		return err
	}

	certinfoCert := certinfo.ParseCertificate(certParse)

	if convertSerial {
		serialConvert(certinfoCert)
	}

	switch ctx.String("format") {
	case "json":
		certAsMarshall, err := json.Marshal(certinfoCert)
		if err != nil {
			return err
		}
		fmt.Println(string(certAsMarshall))
	case "pretty_json":
		s, err := prettyjson.Marshal(certinfoCert)
		if err != nil {
			return err
		}
		fmt.Println(string(s))
	case "table":
		tablePrint([]*certinfo.Certificate{certinfoCert})
	}

	return nil

}
