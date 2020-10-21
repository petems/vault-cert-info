package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/certutil"

	// helpers
	"github.com/cloudflare/cfssl/certinfo"
	"github.com/hokaccha/go-prettyjson"
	"github.com/olekukonko/tablewriter"
	vltcrthlpr "github.com/petems/vault-cert-helpers"

	"github.com/urfave/cli/v2"
)

// Version is what is returned by the `-v` flag
const Version = "0.1.0"

// gitCommit is the gitcommit its built from
var gitCommit = "development"

func main() {
	//nolint:dupl // CLI config is repetitive and flags as duplicates
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
				Name:  "expiry",
				Usage: "List certificates expiring within a certain amount of days",
				Action: func(c *cli.Context) error {
					err := cmdVaultExpiringCerts(c)
					return err
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "expiry_days",
						Usage:    "Number of days until cert expiry to show",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "remain_days",
						Value: false,
						Usage: "Output table with remaining days instead of the expiry date",
					},
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
				Usage: "The path to your pki engine",
			},
			&cli.BoolFlag{
				Name:  "sort",
				Value: true,
				Usage: "Sort certs A-Z by cert.Subject.CommonName",
			},
			&cli.BoolFlag{
				Name:  "serial",
				Value: true,
				Usage: "Convert the cert serial from bigint to a HEX formatted string\n\t\teg. 611168959691622484330452100723265332608845077531 -> 6b:0d:c3:94:c9:e1:20:d1:9a:eb:76:66:db:3d:8a:37:23:75:dc:1b",
			},
			&cli.StringFlag{
				Name:  "format",
				Value: "pretty_json",
				Usage: "The format you want them returned in, valid values are: table, json, pretty_json",
			},
			&cli.BoolFlag{
				Name:  "silent",
				Value: true,
				Usage: "Do not output anything other than errors or returned data",
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

func daysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}

	days := -a.YearDay()
	for year := a.Year(); year < b.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += b.YearDay()

	return days
}

func tablePrint(arrayOfCerts []*certinfo.Certificate) {

	data := [][]string{}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeader([]string{"Common Name", "Organization", "Expires", "Serial"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	utc := time.FixedZone("UTC+0", 9*60*60)

	for _, cert := range arrayOfCerts {
		data = append(data, []string{cert.Subject.CommonName, cert.Subject.Organization, cert.NotAfter.In(utc).Format(time.RFC3339), cert.SerialNumber})
	}

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}

func tablePrintDaysToExpiry(arrayOfCerts []*certinfo.Certificate) {

	data := [][]string{}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeader([]string{"Common Name", "Organization", "Expires", "Days Until Expiry", "Serial"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	utc := time.FixedZone("UTC+0", 9*60*60)

	for _, cert := range arrayOfCerts {

		currentTime := time.Now()
		daysUntilExpiry := daysBetween(cert.NotAfter.In(utc), currentTime)

		data = append(data, []string{cert.Subject.CommonName, cert.Subject.Organization, cert.NotAfter.In(utc).Format(time.RFC3339), fmt.Sprintf("%v", daysUntilExpiry), cert.SerialNumber})
	}

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}

func printResults(format string, certArray []*certinfo.Certificate) error {

	//nolint:dupl // JSON case gets flagged here
	switch format {
	case "json":
		certAsMarshall, err := json.Marshal(certArray)
		if err != nil {
			return err
		}
		fmt.Println(string(certAsMarshall))
	case "pretty_json":
		s, err := prettyjson.Marshal(certArray)
		if err != nil {
			return err
		}
		fmt.Println(string(s))
	case "table":
		tablePrint(certArray)
	default:
		return fmt.Errorf("invalid format given. valid formats: json, pretty_json, table, got: %s", format)
	}

	return nil

}

func wrongPkiPath(pkiPath string, command string) {
	if pkiPath == "pki" {
		fmt.Printf("No certs found under 'pki', have you mounted with a different path? Use the parameter 'vault-cert-info --pki=alt_pki_path %s'\n", command)
	}
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

	if !silent {
		fmt.Printf("Reading certs from %s/%s\n", vaultAddr, pkiPath)
	}

	listOfCertsSecret, err := vltcrthlpr.GetListOfCerts(client, pkiPath)

	if err != nil {
		wrongPkiPath(pkiPath, "list")
		return err
	}

	arrayOfCerts, err := vltcrthlpr.GetArrayOfCertsFromVault(client, listOfCertsSecret, pkiPath)

	if err != nil {
		return err
	}

	sortOption := ctx.Bool("sort")

	if sortOption {
		sort.Slice(arrayOfCerts, func(i, j int) bool { return arrayOfCerts[i].Subject.CommonName < arrayOfCerts[j].Subject.CommonName })
	}

	var arrayOfCertInfo = []*certinfo.Certificate{}

	for _, cert := range arrayOfCerts {

		certinfoCert := certinfo.ParseCertificate(cert)

		if convertSerial {
			serialConvert(certinfoCert)
		}

		arrayOfCertInfo = append(arrayOfCertInfo, certinfoCert)
	}

	err = printResults(ctx.String("format"), arrayOfCertInfo)

	if err != nil {
		return err
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

	err = printResults(ctx.String("format"), []*certinfo.Certificate{certinfoCert})

	if err != nil {
		return err
	}

	return nil

}

func cmdVaultExpiringCerts(ctx *cli.Context) (err error) {

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

	if !silent {
		fmt.Printf("Reading certs from %s/%s\n", vaultAddr, pkiPath)
	}

	listOfCertsSecret, err := vltcrthlpr.GetListOfCerts(client, pkiPath)

	if err != nil {
		wrongPkiPath(pkiPath, "expiry --max_days=31")
		return err
	}

	arrayOfCerts, err := vltcrthlpr.GetArrayOfCertsFromVault(client, listOfCertsSecret, pkiPath)

	if err != nil {
		return err
	}

	var arrayOfCertExpiringInfo = []*certinfo.Certificate{}

	expiryDays := ctx.String("expiry_days")
	expiryDaysInt, err := strconv.Atoi(expiryDays)

	if err != nil {
		return err
	}

	for _, cert := range arrayOfCerts {

		certinfoCert := certinfo.ParseCertificate(cert)

		currentTime := time.Now()

		if daysBetween(certinfoCert.NotAfter, currentTime) <= expiryDaysInt {
			if convertSerial {
				serialConvert(certinfoCert)
			}

			arrayOfCertExpiringInfo = append(arrayOfCertExpiringInfo, certinfoCert)
		}

	}

	remainDays := ctx.Bool("remain_days")

	if remainDays {
		tablePrintDaysToExpiry(arrayOfCertExpiringInfo)
	} else {
		err = printResults(ctx.String("format"), arrayOfCertExpiringInfo)
	}

	if err != nil {
		return err
	}

	return nil

}
