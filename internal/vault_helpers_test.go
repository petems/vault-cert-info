package vaulthelpers

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/vault/api"
)

var recorderMode recorder.Mode

func TestGetArrayOfCertsFromVault(t *testing.T) {

	_, ok := os.LookupEnv("RECORDING")
	if !ok {
		recorderMode = recorder.ModeReplaying
	} else {
		recorderMode = recorder.ModeRecording
	}

	// Start our recorder
	r, err := recorder.NewAsMode("fixtures/get_array_of_certs_from_vault", recorderMode, nil)

	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	// Create an HTTP client and inject our transport
	client := &http.Client{
		Transport: r, // Inject as transport!
	}

	token, tokenOK := os.LookupEnv("VAULT_TOKEN")
	vaultAddr, vaultAddrOK := os.LookupEnv("VAULT_ADDR")
	if !tokenOK || !vaultAddrOK {
		t.Fatal("environment variables VAULT_TOKEN and VAULT_ADDR not set")
	}

	vaultClient, err := api.NewClient(&api.Config{Address: vaultAddr, HttpClient: client})
	if err != nil {
		t.Fatalf("Failed to get new Vault client: %s", err)
	}

	vaultClient.SetToken(token)

	listOfCerts, err := vaultClient.Logical().List(fmt.Sprintf("%s/certs/", "pki"))

	assert.NoError(t, err, "vaultClient.Logical().List(\"pki/certs/\") returned an error")

	assert.NotNil(t, listOfCerts, "vaultClient.Logical().List(\"pki/certs/\") returned nil")

	arrayOfCerts, err := getArrayOfCertsFromVault(vaultClient, listOfCerts, vaultAddr, "pki", true)

	assert.NoError(t, err, "vaultClient.Logical().List(\"pki/certs/\") returned an error")

	assert.NotNil(t, arrayOfCerts, "vaultClient.Logical().List(\"pki/certs/\") returned nil")
}
