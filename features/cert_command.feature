Feature: Cert Command

  Background:
    Given I have "go" command installed
    And I have "docker" command installed
    And nothings running on port "8008"
    When I run `go build -o ../../bin/vault-cert-info-int-test ../../main.go`
    Then the exit status should be 0
  
  Scenario: Listing certs when they are present
    Given no old "dummyvault" containers exist
    And I set the environment variables to:
      | variable           | value               |
      | VAULT_TOKEN        | ROOT                |
      | VAULT_ADDR         | http://0.0.0.0:8008 |
    And I have a dummy vault server running called "dummyvault" running on port "8008" with root token "ROOT"
    And I have the PKI backend enabled at "pki" with a test cert
    When I run "bin/vault-cert-info-int-test cert" with the first serial
    Then the output should contain "arubatest.com"
    And the exit status should be 0

  Scenario: Looking up a non-existance cert
    Given no old "dummyvault" containers exist
    And I set the environment variables to:
      | variable           | value               |
      | VAULT_TOKEN        | ROOT                |
      | VAULT_ADDR         | http://0.0.0.0:8008 |
    And I have a dummy vault server running called "dummyvault" running on port "8008" with root token "ROOT"
    And I have the PKI backend enabled at "pki" with a test cert
    When I run `bin/vault-cert-info-int-test cert "NON-EXISTANT"`
    Then the output should contain "No value found for cert at 'pki/cert/NON-EXISTANT'"
    And the exit status should be 1