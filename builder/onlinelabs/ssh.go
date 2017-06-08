package onlinelabs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/multistep"
	"golang.org/x/crypto/ssh"
)

func sshAddress(state multistep.StateBag) (string, error) {
	ipAddress := state.Get("server_ip").(string)
	return fmt.Sprintf("%s", ipAddress), nil
}

func sshConfig(state multistep.StateBag) (*ssh.ClientConfig, error) {
	config := state.Get("config").(Config)
	// privateKey := state.Get("privateKey").(string)

	// signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	log.Println(config.Comm.SSHPrivateKey)
	pkFile, err := os.Open(config.Comm.SSHPrivateKey)
	if err != nil {
		return nil, err
	}
	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(pkBytes)
	if err != nil {
		return nil, fmt.Errorf("Error setting up SSH config: %s", err)
	}

	return &ssh.ClientConfig{
		User: config.Comm.SSHUsername,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}
