Feature: Tidy Command

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
    And I have enabled vault audits log
    And I have the PKI backend enabled at "pki" with a test cert
    When I run `bin/vault-cert-info-int-test tidy` 
    Then the output should contain:
      """
      Tidy command complete
      """
    And the exit status should be 0
    And the Docker log output for "dummyvault" should contain "Tidy operation successfully started"