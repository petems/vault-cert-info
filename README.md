# vault-cert-info

[![Build Status](https://travis-ci.com/petems/vault-cert-info.svg?branch=master)](https://travis-ci.com/petems/vault-cert-info)

A simple cli app to return certificates from a Vault PKI mount

## Install 

```
# install it into ./bin/
curl -sSfL https://raw.githubusercontent.com/petems/vault-cert-info/master/install.sh | sh -s v0.1.0
```

## Example

### Help

```
$ vault-cert-info
NAME:
  vault-cert-info - A simple cli app to return certificates from a Vault PKI mount
USAGE:
  vault-cert-info [global options] command [command options] [arguments...]

COMMANDS:
   list     List all certificates
   cert     Get information from one certificate
   help, h  Shows a list of commands or help for one command

VAULT CONFIGURATION:
  Vault configuration is set by the common Vault environmental variables:
    VAULT_ADDR: The address for the Vault server (Required)
    VAULT_TOKEN: The token for the Vault server (Required)
    VAULT_NAMESPACE: The namespace where PKI is mounted (Optional)
  See https://www.vaultproject.io/docs/commands/#environment-variables

GLOBAL OPTIONS:
  --pki value     The path to your pki engine (default: "pki")
  --serial        Convert the cert serial from bigint to a HEX formatted string
                  eg. 611168959691622484330452100723265332608845077531 -> 6b:0d:c3:94:c9:e1:20:d1:9a:eb:76:66:db:3d:8a:37:23:75:dc:1b (default: true)
  --format value  The format you want them returned in, valid values are: table, json, pretty_json (default: "pretty_json")
  --silent        Do not output anything other than errors or returned data (default: true)
  --help, -h      show help (default: false)
  --version, -v   print the version (default: false)

VERSION:
  0.1.0-ca478c3
```

### List 

#### Pretty JSON 

```
$ vault-cert-info list
[
  {
    "authority_key_id": "8F:50:00:B3:62:DB:93:20:0E:55:36:1A:06:8A:D0:7C:E0:D3:CE:56",
    "issuer": {
      "common_name": "example.com",
      "names": [
        "example.com"
      ]
    },
    "not_after": "2025-03-12T14:38:01Z",
    "not_before": "2020-03-13T14:37:31Z",
    "pem": "-----BEGIN CERTIFICATE-----\nMIIDpjCCAo6gAwIBAgIUItD3L/bBJsosfPrXY6wrqX06iTAwDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMjAwMzEzMTQzNzMxWhcNMjUw\nMzEyMTQzODAxWjAtMSswKQYDVQQDEyJleGFtcGxlLmNvbSBJbnRlcm1lZGlhdGUg\nQXV0aG9yaXR5MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1veq6qgz\nX8X7efKNQLF7BzTKd5iFm7MypSZTpfd6kunUSKCrLoIPH+oNTUbxXLsGXPxsKvSt\nb4DNoZ/XJkCPjTjNY3m11AWDD3Yg/Ons/KBPlfIwPW/c0tQs3N1t+b83lSWbU98B\nFt/pmfQelsG2lP+N7YqGTYGkShhdgO1BApJizjlO0xOyrlnKqUZrm3ccIII+iHHo\n5CIHnwZoFXTrixuWDumE6nsCI7nQw4LJuuNCkOQfdVvVrcnWXK8fwRpHsZjcX4fL\nv6JpSkVkIfj3zpp47b2zhdwPi8MTthvlHcDwU7+iseHsClGDhJ0FfSOpvnwQ4Wis\nmHlPbCYMCzVXVQIDAQABo4HUMIHRMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8E\nBTADAQH/MB0GA1UdDgQWBBQTW6RW6565S3W0gqr8G+KFQADmVjAfBgNVHSMEGDAW\ngBSPUACzYtuTIA5VNhoGitB84NPOVjA7BggrBgEFBQcBAQQvMC0wKwYIKwYBBQUH\nMAKGH2h0dHA6Ly8xMjcuMC4wLjE6ODIwMC92MS9wa2kvY2EwMQYDVR0fBCowKDAm\noCSgIoYgaHR0cDovLzEyNy4wLjAuMTo4MjAwL3YxL3BraS9jcmwwDQYJKoZIhvcN\nAQELBQADggEBAEwrVmDoIkamedgRvLdiyUla+DP6L1FCLlg/G+MhyGqdaDdI9zZm\noEfF7b1BtgKG+G2GrCIyZdmafCkZbRnfn+qQLsPd8rHFrhqCmr8PKJckRMXFWniJ\np5Bd1N9pziVvnctsu9JatGTMzxYvvj14UJri9aMSfCcpDscxKz9sqh+l8QCxC9qJ\nbIjLj4hXgw7ggHGYVjhcqM8ifloGOsTZ1DAvNWEhoVRzw4t2083Ro0g9dS9i08VB\nnrgae+OMIdV+B6Xw14GXXqpIEe4al+vN+6l9hhGPal3W0qKNvAzxue8GRDil2D4b\neQj3+9rzqbUdkaIhZosSX9/iF32FEpCztt0=\n-----END CERTIFICATE-----\n",
    "serial_number": "198765774265194290508604545037765352926991780144",
    "sigalg": "SHA256WithRSA",
    "subject": {
      "common_name": "example.com Intermediate Authority",
      "names": [
        "example.com Intermediate Authority"
      ]
    },
    "subject_key_id": "13:5B:A4:56:EB:9E:B9:4B:75:B4:82:AA:FC:1B:E2:85:40:00:E6:56"
  },
  {
    "authority_key_id": "8F:50:00:B3:62:DB:93:20:0E:55:36:1A:06:8A:D0:7C:E0:D3:CE:56",
    "issuer": {
      "common_name": "example.com",
      "names": [
        "example.com"
      ]
    },
    "not_after": "2030-03-11T14:37:34Z",
    "not_before": "2020-03-13T14:37:04Z",
    "pem": "-----BEGIN CERTIFICATE-----\nMIIDNTCCAh2gAwIBAgIUaw3DlMnhINGa63Zm2z2KNyN13BswDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMjAwMzEzMTQzNzA0WhcNMzAw\nMzExMTQzNzM0WjAWMRQwEgYDVQQDEwtleGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAMkIUiPd9mJ+z9HA5sky0G9gWxME2q2A8Hk2Eaxw\nX/pqzh6WMj01WNgSkTQGIt+dHUhlAaDY4L2xATUzv+SU1PkDBL0JkdJck/f60ygi\ngiUjEyp9+VgCPgxKPUo9IBI3oJZHR0LTYvrwoqrBi4Ra0qOuZdr/5mqu07pSbWaE\nie935PH7jasCXoiHJKBSNsLjYNMV0XgGkh0dfxeSBVpb/degs23lQZIqsX8LiqMd\n5gZGtOddxf95o00sgjqBcAhf5DliknDB34GE8gVAAy8UXrm5bBt0nIgDml9srejt\nZfOc54adZnoWWVzSwQ2pcqOSYXD3oaMP96IQ7naXAh9Ws9ECAwEAAaN7MHkwDgYD\nVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFI9QALNi25Mg\nDlU2GgaK0Hzg085WMB8GA1UdIwQYMBaAFI9QALNi25MgDlU2GgaK0Hzg085WMBYG\nA1UdEQQPMA2CC2V4YW1wbGUuY29tMA0GCSqGSIb3DQEBCwUAA4IBAQAXnIGY3+ju\no3/eaHBJxkK+kxJnq29R/+iBFFh2XE6oT0kCIygnLovIZfNsEjiHUUYEcT8CmQte\nPobkWkdO8F0VDaOGRcva2W68iAR9JYb5pkT7hdt+0VehCyWuKJHWud2s54J7CCRL\np/pge+N5cEHx0u/DUNuvIxJDePTKN1BaSfVNwMWXrP1liA8sBsPXKhyneTI3Blve\nlnaNnWWNm7kRgy47jukBLGRgRWr8b0yPjrKlryRyu61OhXCmuNlSsz+cu77G+RN5\nWczVvypplc9QSNDWwAPgFNYLsVJWwwp166FVRF+AbzfkSUbjt95c5zDVe56HtApG\nFK3XQEVNAXWI\n-----END CERTIFICATE-----\n",
    "sans": [
      "example.com"
    ],
    "serial_number": "611168959691622484330452100723265332608845077531",
    "sigalg": "SHA256WithRSA",
    "subject": {
      "common_name": "example.com",
      "names": [
        "example.com"
      ]
    },
    "subject_key_id": "8F:50:00:B3:62:DB:93:20:0E:55:36:1A:06:8A:D0:7C:E0:D3:CE:56"
  }
]
```

#### Regular JSON

```
vault-cert-info --format=json list
[{"subject":{"common_name":"example.com Intermediate Authority","names":["example.com Intermediate Authority"]},"issuer":{"common_name":"example.com","names":["example.com"]},"serial_number":"198765774265194290508604545037765352926991780144","not_before":"2020-03-13T14:37:31Z","not_after":"2025-03-12T14:38:01Z","sigalg":"SHA256WithRSA","authority_key_id":"8F:50:00:B3:62:DB:93:20:0E:55:36:1A:06:8A:D0:7C:E0:D3:CE:56","subject_key_id":"13:5B:A4:56:EB:9E:B9:4B:75:B4:82:AA:FC:1B:E2:85:40:00:E6:56","pem":"-----BEGIN CERTIFICATE-----\nMIIDpjCCAo6gAwIBAgIUItD3L/bBJsosfPrXY6wrqX06iTAwDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMjAwMzEzMTQzNzMxWhcNMjUw\nMzEyMTQzODAxWjAtMSswKQYDVQQDEyJleGFtcGxlLmNvbSBJbnRlcm1lZGlhdGUg\nQXV0aG9yaXR5MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1veq6qgz\nX8X7efKNQLF7BzTKd5iFm7MypSZTpfd6kunUSKCrLoIPH+oNTUbxXLsGXPxsKvSt\nb4DNoZ/XJkCPjTjNY3m11AWDD3Yg/Ons/KBPlfIwPW/c0tQs3N1t+b83lSWbU98B\nFt/pmfQelsG2lP+N7YqGTYGkShhdgO1BApJizjlO0xOyrlnKqUZrm3ccIII+iHHo\n5CIHnwZoFXTrixuWDumE6nsCI7nQw4LJuuNCkOQfdVvVrcnWXK8fwRpHsZjcX4fL\nv6JpSkVkIfj3zpp47b2zhdwPi8MTthvlHcDwU7+iseHsClGDhJ0FfSOpvnwQ4Wis\nmHlPbCYMCzVXVQIDAQABo4HUMIHRMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8E\nBTADAQH/MB0GA1UdDgQWBBQTW6RW6565S3W0gqr8G+KFQADmVjAfBgNVHSMEGDAW\ngBSPUACzYtuTIA5VNhoGitB84NPOVjA7BggrBgEFBQcBAQQvMC0wKwYIKwYBBQUH\nMAKGH2h0dHA6Ly8xMjcuMC4wLjE6ODIwMC92MS9wa2kvY2EwMQYDVR0fBCowKDAm\noCSgIoYgaHR0cDovLzEyNy4wLjAuMTo4MjAwL3YxL3BraS9jcmwwDQYJKoZIhvcN\nAQELBQADggEBAEwrVmDoIkamedgRvLdiyUla+DP6L1FCLlg/G+MhyGqdaDdI9zZm\noEfF7b1BtgKG+G2GrCIyZdmafCkZbRnfn+qQLsPd8rHFrhqCmr8PKJckRMXFWniJ\np5Bd1N9pziVvnctsu9JatGTMzxYvvj14UJri9aMSfCcpDscxKz9sqh+l8QCxC9qJ\nbIjLj4hXgw7ggHGYVjhcqM8ifloGOsTZ1DAvNWEhoVRzw4t2083Ro0g9dS9i08VB\nnrgae+OMIdV+B6Xw14GXXqpIEe4al+vN+6l9hhGPal3W0qKNvAzxue8GRDil2D4b\neQj3+9rzqbUdkaIhZosSX9/iF32FEpCztt0=\n-----END CERTIFICATE-----\n"},{"subject":{"common_name":"example.com","names":["example.com"]},"issuer":{"common_name":"example.com","names":["example.com"]},"serial_number":"611168959691622484330452100723265332608845077531","sans":["example.com"],"not_before":"2020-03-13T14:37:04Z","not_after":"2030-03-11T14:37:34Z","sigalg":"SHA256WithRSA","authority_key_id":"8F:50:00:B3:62:DB:93:20:0E:55:36:1A:06:8A:D0:7C:E0:D3:CE:56","subject_key_id":"8F:50:00:B3:62:DB:93:20:0E:55:36:1A:06:8A:D0:7C:E0:D3:CE:56","pem":"-----BEGIN CERTIFICATE-----\nMIIDNTCCAh2gAwIBAgIUaw3DlMnhINGa63Zm2z2KNyN13BswDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMjAwMzEzMTQzNzA0WhcNMzAw\nMzExMTQzNzM0WjAWMRQwEgYDVQQDEwtleGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAMkIUiPd9mJ+z9HA5sky0G9gWxME2q2A8Hk2Eaxw\nX/pqzh6WMj01WNgSkTQGIt+dHUhlAaDY4L2xATUzv+SU1PkDBL0JkdJck/f60ygi\ngiUjEyp9+VgCPgxKPUo9IBI3oJZHR0LTYvrwoqrBi4Ra0qOuZdr/5mqu07pSbWaE\nie935PH7jasCXoiHJKBSNsLjYNMV0XgGkh0dfxeSBVpb/degs23lQZIqsX8LiqMd\n5gZGtOddxf95o00sgjqBcAhf5DliknDB34GE8gVAAy8UXrm5bBt0nIgDml9srejt\nZfOc54adZnoWWVzSwQ2pcqOSYXD3oaMP96IQ7naXAh9Ws9ECAwEAAaN7MHkwDgYD\nVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFI9QALNi25Mg\nDlU2GgaK0Hzg085WMB8GA1UdIwQYMBaAFI9QALNi25MgDlU2GgaK0Hzg085WMBYG\nA1UdEQQPMA2CC2V4YW1wbGUuY29tMA0GCSqGSIb3DQEBCwUAA4IBAQAXnIGY3+ju\no3/eaHBJxkK+kxJnq29R/+iBFFh2XE6oT0kCIygnLovIZfNsEjiHUUYEcT8CmQte\nPobkWkdO8F0VDaOGRcva2W68iAR9JYb5pkT7hdt+0VehCyWuKJHWud2s54J7CCRL\np/pge+N5cEHx0u/DUNuvIxJDePTKN1BaSfVNwMWXrP1liA8sBsPXKhyneTI3Blve\nlnaNnWWNm7kRgy47jukBLGRgRWr8b0yPjrKlryRyu61OhXCmuNlSsz+cu77G+RN5\nWczVvypplc9QSNDWwAPgFNYLsVJWwwp166FVRF+AbzfkSUbjt95c5zDVe56HtApG\nFK3XQEVNAXWI\n-----END CERTIFICATE-----\n"}]
```