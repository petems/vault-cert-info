Feature: Version Command

  Background:
    Given I have "go" command installed
    When I run `go build -o ../../bin/vault-cert-info-int-test ../../main.go`
    Then the exit status should be 0

  Scenario: Version with no flags
    Given a build of vault-cert-info
    When I run `bin/vault-cert-info-int-test --version`
    Then the output should contain:
      """""

      vault-cert-info version 0.1.0-development
      """""