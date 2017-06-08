package onlinelabs

import (
	"fmt"

	"github.com/mitchellh/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepStartServer struct{}

func (s *stepStartServer) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(ClientInterface)
	ui := state.Get("ui").(packer.Ui)
	serverID := state.Get("server_id").(string)

	err := client.PowerOnServer(serverID)
	if err != nil {
		err := fmt.Errorf("Error starting server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepStartServer) Cleanup(state multistep.StateBag) {
}
