#!/bin/bash
#

curl -sf -H "X-Vault-Token: $VAULT_TOKEN" -d '{"public_key":"'"$(cat ~/.ssh/talk_rsa.pub)"'"}' $VAULT_ADDR/v1/ssh/sign/user | jq -r '.data.signed_key' > /tmp/talk_rsa-cert

server=$1

ssh -i ~/.ssh/talk_rsa -o CertificateFile=/tmp/talk_rsa-cert -L 2375:/var/run/docker.sock -N user@$server &
ssh_pid=$!

sleep 3

export DOCKER_HOST=127.0.0.1

docker version
docker run --name run-aweful hello-world

kill $ssh_pid

