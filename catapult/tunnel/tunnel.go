package tunnel

import (
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// Create connects to the specified server and establishes the tunnel.
func Create(username string, signer ssh.Signer, server, local, remote net.Addr) {

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial(server.Network(), server.String(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to connect to server %s:%s: %s\n", server.Network(), server.String(), err)
		return
	}
	defer client.Close()

	localListener, err := net.Listen(local.Network(), local.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to open listener socket %s:%s for local end of tunnel: %s\n", local.Network(), local.String(), err)
	}
	defer localListener.Close()

	for {
		localConn, err := localListener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to accept connection from listener socket %s:%s: %s", local.Network(), local.String(), err)
			continue
		}
		// localConn gets closed in the copyConnection(localConn, remoteConn) function below

		remoteConn, err := client.Dial(remote.Network(), remote.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to connect to remote end of tunnel %s:%s: %s", remote.Network(), remote.String(), err)
		}
		// remoteConn gets closed in the copyConnection(remoteConn, localConn) function below

		go copyConnection(remoteConn, localConn)
		go copyConnection(localConn, remoteConn)
	}

}

func copyConnection(writer, reader net.Conn) {
	if _, err := io.Copy(writer, reader); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to transfer data in tunnel: %s", err)
	}

	writer.Close()
}
