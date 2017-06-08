package onlinelabs

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
)

type stepShutdown struct{}

func (s *stepShutdown) Run(state multistep.StateBag) multistep.StepAction {
	comm := state.Get("communicator").(packer.Communicator)
	ui := state.Get("ui").(packer.Ui)
	serverID := state.Get("server_id").(string)
	client := state.Get("client").(ClientInterface)

	ui.Say("Gracefully shutting down server...")

	cmd := &packer.RemoteCmd{Command: "halt"}

	if err := comm.Start(cmd); err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	cmd.Wait()
	if cmd.ExitStatus != 0 {
		err := fmt.Errorf("shutdown exited with non-zero exit status: %d", cmd.ExitStatus)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	err := waitForServerState("stopped", serverID, client, 30*time.Second)
	if err != nil {
		log.Printf("Error waiting for graceful off: %s", err)
	}

	return multistep.ActionContinue
}

func (s *stepShutdown) Cleanup(state multistep.StateBag) {
	// no cleanup
}
