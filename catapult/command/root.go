package command

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"
	"strings"

	"github.com/marcboudreau/go-devops-talk/catapult/tunnel"

	"github.com/marcboudreau/go-devops-talk/catapult"
	"github.com/marcboudreau/go-devops-talk/catapult/vault"
	"github.com/spf13/cobra"
)

var privateKeyFilename string

var publicKeyFilename string

var localAddressStr string

var remoteAddressStr string

var keySigningService catapult.KeySigningService

var rootCmd = &cobra.Command{
	Use:   "catapult username@server",
	Short: "Catapult signs SSH keys and then uses them to establish a tunnel (forward a local port) to the specified server.",
	Long: `Catapult uses a key signing service to sign a given public SSH key.  
It then uses that signed public key (certificate) to connect with the specified server.
Once connected, it establishes a tunnel by opening a local port and forwarding all data
it receives to the specified remote port, and vice-versa.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Error: missing argument")
			fmt.Fprintln(os.Stderr, cmd.UsageString())
			return
		}

		if len(args) > 1 {
			fmt.Fprintln(os.Stderr, "Warning: extra arguments will be ignored")
		}

		username, serverAddress, err := parseArg(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to parse command argument %s.  Error: %s\n", args[0], err)
			return
		}

		publicKey, err := os.Open(publicKeyFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to open public key file %s.  Error: %s\n", publicKeyFilename, err)
			return
		}
		defer publicKey.Close()

		keySigningService, err = vault.New("user")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create Vault client for key signing.  Error: %s\n", err)
			return
		}

		certificate, err := keySigningService.SignKey(publicKey, username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to sign public key.  Error: %s\n", err)
			return
		}

		privateKey, err := os.Open(privateKeyFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to open private key file %s.  Error: %s\n", privateKeyFilename, err)
			return
		}

		certificateSigner, err := tunnel.CreateSigner(privateKey, certificate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create public key signer.  Error: %s\n", err)
			return
		}

		if !strings.Contains(serverAddress, ":") {
			serverAddress = serverAddress + ":22"
		}

		server, err := net.ResolveTCPAddr("tcp", serverAddress)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to parse server address %s.  Error: %s\n", serverAddress, err)
			return
		}

		local, err := parseAddress(localAddressStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to parse local address %s.  Error: %s\n", localAddressStr, err)
			return
		}

		remote, err := parseAddress(remoteAddressStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to parse remote address %s.  Error: %s\n", remoteAddressStr, err)
			return
		}

		tunnel.Create(username, certificateSigner, server, local, remote)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&privateKeyFilename, "privateKey", "k", "", "File containing the private SSH key used to connect to the server.")
	rootCmd.Flags().StringVarP(&publicKeyFilename, "publicKey", "p", "", "File containing the public SSH key to sign.")
	rootCmd.Flags().StringVarP(&localAddressStr, "localAddress", "l", "", "Network address of local port of the tunnel to establish.")
	rootCmd.Flags().StringVarP(&remoteAddressStr, "remoteAddress", "r", "", "Network address of remote port of the tunnel to establish.")
}

// Execute executes the rootCmd Command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encountered during execution: %s\n", err)
		os.Exit(1)
	}
}

func parseArg(arg string) (username, server string, err error) {
	pos := strings.Index(arg, "@")
	if pos == -1 {
		user, err := user.Current()
		if err != nil {
			return "", "", err
		}

		return user.Username, arg, nil
	}

	if pos == 0 {
		// arg value starts with @, so no username provided
		return "", "", errors.New("a username must be specified before the @ sign in the server address")
	}

	if pos == len(arg)-1 {
		// arg value ends with @, so no server address provided
		return "", "", errors.New("a server address must be specified after the @ sign")
	}

	return arg[:pos], arg[pos+1:], nil
}

func parseAddress(address string) (net.Addr, error) {
	var network string

	pos := strings.Index(address, ":")
	if pos == -1 {
		network = "tcp"
	} else {
		network = address[:pos]
	}
	rest := address[pos+1:]

	switch network[:2] {
	case "tc":
		return net.ResolveTCPAddr(network, rest)
	case "ud":
		return net.ResolveUDPAddr(network, rest)
	case "ip":
		return net.ResolveIPAddr(network, rest)
	case "un":
		return net.ResolveUnixAddr(network, rest)
	}

	return nil, fmt.Errorf("unknown network %s", network)
}
