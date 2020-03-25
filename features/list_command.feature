Feature: List Command

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
    And I have the PKI backend enabled with a test cert
    When I run `bin/vault-cert-info-int-test list` 
    Then the output should contain:
      """
      "common_name": "arubatest.com",
      """
    And the exit status should be 0
  
  Scenario: Listing certs when they are present with --format table
    Given no old "dummyvault" containers exist
    And I set the environment variables to:
      | variable           | value               |
      | VAULT_TOKEN        | ROOT                |
      | VAULT_ADDR         | http://0.0.0.0:8008 |
    And I have a dummy vault server running called "dummyvault" running on port "8008" with root token "ROOT"
    And I have the PKI backend enabled with a test cert
    When I run `bin/vault-cert-info-int-test --format=table list` 
    Then the output should contain "COMMON NAME"
    And the output should contain "arubatest.com"
    And the exit status should be 0

  Scenario: Error when wrong endpoint is given
    Given no old "dummyvault" containers exist
    And I set the environment variables to:
      | variable           | value               |
      | VAULT_TOKEN        | ROOT                |
      | VAULT_ADDR         | http://0.0.0.0:8008 |
    And I have a dummy vault server running called "dummyvault" running on port "8008" with root token "ROOT"
    And I have the PKI backend enabled with a test cert
    When I run `bin/vault-cert-info-int-test --pki=not_exist list` 
    Then the output should contain:
      """
      No certs found at not_exist/certs/
      """
    And the exit status should be 1