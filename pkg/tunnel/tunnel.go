package tunnel

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/colonyos/colonies/pkg/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type Tunnel struct {
	client       *client.ColoniesClient
	processID    string
	JumpHost     string
	JumpHostPort int
	User         string
	SSHKey       string
	LocalPort    int
	RemotePort   int
	prvKey       string
}

func NewTunnel(client *client.ColoniesClient, processID string, jumpHost string, jumpHostPort int, user string, sshKey string, localPort int, remotePort int, prvKey string) *Tunnel {
	return &Tunnel{
		client:       client,
		processID:    processID,
		JumpHost:     jumpHost,
		JumpHostPort: jumpHostPort,
		User:         user,
		SSHKey:       sshKey,
		LocalPort:    localPort,
		RemotePort:   remotePort,
		prvKey:       prvKey,
	}
}

func (t *Tunnel) Start() {
	log.Info("Starting tunnel")
	nodes := ""
outerBound:
	for {
		process, err := t.client.GetProcess(t.processID, t.prvKey)
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Failed to get process")
			continue
		}
		if process != nil {
			for _, attr := range process.Attributes {
				if attr.Key == "NODES" {
					nodes = attr.Value
					break outerBound
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Process running at nodes", nodes)

	signer, err := loadKey(t)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Error("Failed to load key")
	}
	config := &ssh.ClientConfig{
		User: t.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", t.JumpHost+":"+strconv.Itoa(t.JumpHostPort), config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}

	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(t.LocalPort))
	if err != nil {
		log.Fatalf("Failed to listen on local port: %s", err)
	}

	fmt.Println("Remote process can be access at http://localhost:" + strconv.Itoa(t.LocalPort))

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept local connection: %s", err)
		}

		go handleConnection(t, nodes, client, localConn)
	}
}

func handleConnection(t *Tunnel, node string, client *ssh.Client, localConn net.Conn) {
	defer localConn.Close()

	remoteConn, err := client.Dial("tcp", node+":"+strconv.Itoa(t.RemotePort))
	if err != nil {
		log.Printf("Failed to dial remote server: %s", err)
		return
	}
	defer remoteConn.Close()

	go io.Copy(remoteConn, localConn)
	io.Copy(localConn, remoteConn)
}

func loadKey(t *Tunnel) (ssh.Signer, error) {
	file, err := os.Open(t.SSHKey)
	if err != nil {
		return nil, fmt.Errorf("unable to open SSH key: %v", err)
	}
	defer file.Close()

	key, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read SSH key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse SSH key: %v", err)
	}

	return signer, nil
}
