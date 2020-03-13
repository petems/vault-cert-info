Feature: Version Command

  Background:
    Given I have "go" command installed
    When I run `go build -o ../../bin/vault-cert-info-int-test ../../main.go`
    Then the exit status should be 0

  Scenario:
    Given a build of vault-cert-info
    When I run `bin/vault-cert-info-int-test`
    Then the output should contain:
      """"
      NAME:
        vault-cert-info - A simple cli app to return certificates from a Vault PKI mount
      """"
