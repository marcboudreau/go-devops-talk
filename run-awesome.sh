#!/bin/bash
set -euo pipefail

server=$1
private_key=${SSH_PRIVATE_KEY_FILE:-${HOME:-.}/.ssh/id_rsa}
public_key=${SSH_PUBLIC_KEY_FILE:-${private_key}.pub}


go run catapult/cmd/catapult/main.go -k $private_key 