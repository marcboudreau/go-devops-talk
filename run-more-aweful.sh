#!/bin/bash
set -euo pipefail

${DEBUG:+set -x}

if [[ -z ${VAULT_ADDR:+x} ]]; then
    echo "ERROR: The VAULT_ADDR environment variable is not set."
    exit 1
fi

if ! curl -fsL -m 10 $VAULT_ADDR/v1/sys/health >/dev/null ; then
    echo "ERROR: The Vault server at $VAULT_ADDR is either not running or unhealthy."
    exit 1
fi

if [[ -z ${VAULT_TOKEN:+x} ]]; then
    echo "ERROR: The VAULT_TOKEN environment variable is not set."
    exit 1
fi

if ! curl -fsL -m 10 -H "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/auth/token/lookup-self >/dev/null ; then
    echo "ERROR: The provided VAULT_TOKEN value is either invalid or expired."
    exit 1
fi

# Allow customizing the SSH backend's mountpoint and the role used to sign keys.
vault_path=${VAULT_SSH_MOUNTPOINT:-ssh}/sign/${VAULT_SSH_ROLE:-user}

if ! curl -fsL -m 10 -H "X-Vault-Token: $VAULT_TOKEN" -d '{"paths":["'$vault_path'"]}' $VAULT_ADDR/v1/sys/capabilities-self >/dev/null ; then
    echo "ERROR: The provided VAULT_TOKEN does not have access to make a request to the $vault_path endpoint."
    exit 1
fi

private_key_file=${SSH_PRIVATE_KEY_FILE:-$HOME/.ssh/talk_rsa}
public_key_file=${SSH_PUBLIC_KEY_FILE:-${private_key_file}.pub}

certificate_file=$(mktemp)
trap "rm $certificate_file" EXIT

if ! curl -fsL -m 10 -H "X-Vault-Token: $VAULT_TOKEN" -d '{"public_key":"'"$(cat $public_key_file)"'"}' $VAULT_ADDR/v1/$vault_path | jq -r '.data.signed_key' > $certificate_file ; then
    echo "ERROR: Failed to obtained signed SSH keys from Vault server."
    exit 1
fi

local_port=${SSH_TUNNEL_LOCAL_PORT:-2375}
remote_addr=${SSH_TUNNEL_REMOTE_ADDR:-/var/run/docker.sock}
username=${SSH_USER:-user}
server=$1

ssh -i $private_key_file -o CertificateFile=$certificate_file -L $local_port:$remote_addr -N $username@$server &
ssh_pid=$!

sleep 3

if ! ps $ssh_pid >/dev/null ; then
    echo "ERROR: Failed to setup SSH tunnel."
    exit 1
fi

export DOCKER_HOST=tcp://127.0.0.1:$local_port

docker version
docker run --name run-more-aweful hello-world

kill $ssh_pid
