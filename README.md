# go-devops-talk

Source material for my DevOps talk presented at the [Ottawa Go Meetup event](https://www.meetup.com/Ottawa-Go-Meetup/events/kbcvhqyzcbmc/).

This material has 3 component directories:
1. catapult
2. packer
3. terraform

As well as some bash scripts in the root directory that are executed as part of the demo.

## catapult

Catapult is a tool written in Go that handles the functionality contained in the **run-aweful.sh** and **run-more-aweful.sh** scripts.
In a nutshell, it submits a request to a Vault server to sign a public SSH key, so that it can be used to authenticate to a server
that has never been given the public SSH key, but that does trust the Certificate Authority that Vault used for signing.

By codifying this functionality in an application, the amount of error handling and testability are greatly improved.

## packer

The packer directory contains a Packer specification to build a Google Compute Engine Image.  The image is built such that
sshd will trust public SSH keys that have been signed by the Vault server's Certificate Authority.

## terraform

The terraform directory contains a Terraform project that manages the infrastructure for this demo: a single Google Compute Engine
VM Instance.

## Bash Scripts

There are 3 bash scripts in the root directory:
1. run-aweful.sh
2. run-more-aweful.sh
3. run-awesome.sh

The idea is that the **run-aweful.sh** script is the initial take at scripting the task of establishing an SSH
tunnel and running some docker commands.

The **run-more-aweful.sh** script shows an attempt to improve the robustness of the script by detecting failure modes,
but the complexity goes way up and there are still lots of gaps.

The **run-awesome.sh** script is much simpler, and leaves the complexity of detecting various errors to the **catapult**
application, which is better suited to do so, and can be properly unit tested.

## License

This material is made available under the MIT license.  See the [LICENSE](./LICENSE) file for details.