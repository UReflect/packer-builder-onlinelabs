package onlinelabs

import (
	"fmt"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type stepServerInfo struct{}

func (s *stepServerInfo) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(ClientInterface)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	serverID := state.Get("server_id").(string)

	ui.Say("Waiting for server to become active...")

	err := waitForServerState("running", serverID, client, c.StateTimeout)
	if err != nil {
		err := fmt.Errorf("Error waiting for server to become active: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	server, err := client.GetServer(serverID)
	if err != nil {
		err := fmt.Errorf("Error retrieving server %s: %s", serverID, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if server.PublicIP == nil {
		err := fmt.Errorf("Error getting server public IP%s: %s", serverID, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("server_ip", server.PublicIP.Address)

	return multistep.ActionContinue
}

func (s *stepServerInfo) Cleanup(state multistep.StateBag) {
}
