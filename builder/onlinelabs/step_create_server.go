package onlinelabs

import (
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type stepCreateServer struct {
	serverID string
}

func (s *stepCreateServer) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(ClientInterface)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)

	ui.Say("Creating server...")
	log.Printf("server creation params: name=%q org=%q image=%q volumes=%v tags=%v",
		c.ServerName, c.OrganizationID, c.ImageID, c.ServerVolumes, c.ServerTags)

	server, err := client.CreateServer(c.ServerName, c.OrganizationID, c.ImageID, c.ServerVolumes, c.ServerTags)

	if err != nil {
		err := fmt.Errorf("Error creating server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if server == nil {
		err := fmt.Errorf("No server created")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	s.serverID = server.ID

	state.Put("server_id", server.ID)
	state.Put("server_arch", server.Image.Arch)

	return multistep.ActionContinue
}

func (s *stepCreateServer) Cleanup(state multistep.StateBag) {
	if s.serverID == "" {
		return
	}

	client := state.Get("client").(ClientInterface)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Destroying server...")

	err := client.PowerOffServer(s.serverID)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error powering off server. Please destroy it manually: %v", s.serverID))
		return
	}

	err = waitForServerState("stopped", s.serverID, client, 30*time.Second)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error waiting for server to stop. Please destroy it manually: %v", s.serverID))
		return
	}

	err = client.DestroyServer(s.serverID)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error destroying server. Please destroy it manually: %v", s.serverID))
	}
}
