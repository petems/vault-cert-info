#!/bin/bash
x=1
while [ $x -le 50 ]
do
  VAULT_TOKEN=ROOT VAULT_ADDR=http://127.0.0.1:8200 vault write pki/issue/example-dot-com common_name=vch$x.example.com
  x=$(( $x + 1 ))
done