#!/bin/bash

VAULT_TOKEN=ROOT VAULT_ADDR=http://127.0.0.1:8200 vault secrets enable pki
VAULT_TOKEN=ROOT VAULT_ADDR=http://127.0.0.1:8200 vault secret tune -max-lease-ttl=8760h pki/
VAULT_TOKEN=ROOT VAULT_ADDR=http://127.0.0.1:8200 vault write pki/root/generate/internal common_name=example.com ttl=8760h
VAULT_TOKEN=ROOT VAULT_ADDR=http://127.0.0.1:8200 vault write pki/roles/example-dot-com allowed_domains=example.com allow_subdomains=true max_ttl=72h