package onlinelabs

import (
	"fmt"
	"log"

	"github.com/mitchellh/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepPowerOff struct{}

func (s *stepPowerOff) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(ClientInterface)
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	serverID := state.Get("server_id").(string)

	server, err := client.GetServer(serverID)
	if err != nil {
		err := fmt.Errorf("Error checking server state: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if server.State == "stopped" {
		return multistep.ActionContinue
	}

	ui.Say("Forcefully shutting down server...")
	err = client.PowerOffServer(serverID)
	if err != nil {
		err := fmt.Errorf("Error powering off server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Println("Waiting for poweroff event to complete...")
	err = waitForServerState("stopped", serverID, client, c.StateTimeout)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepPowerOff) Cleanup(state multistep.StateBag) {
	// no cleanup
}
