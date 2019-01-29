#!/bin/bash
set -euo pipefail

server=$1
private_key=${SSH_PRIVATE_KEY_FILE:-${HOME:-.}/.ssh/id_rsa}
public_key=${SSH_PUBLIC_KEY_FILE:-${private_key}.pub}

catapult/bin/catapult -k $private_key -p $public_key -l tcp:127.0.0.1:2375 -r unix:/var/run/docker.sock user@$server &
catapult_pid=$!

echo "Waiting for tunnel to be established..."

while ! netstat -ant | grep LISTEN | grep 2375 >/dev/null; do
  echo -n "."
done
echo ""

export DOCKER_HOST=tcp://127.0.0.1:2375

docker version
docker ps -a

kill $catapult_pid
