package onlinelabs

import (
	"fmt"
	"log"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type stepSnapshot struct{}

func (s *stepSnapshot) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(ClientInterface)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(*config)
	serverID := state.Get("server_id").(string)
	server, err := client.GetServer(serverID)
	if err != nil {
		err := fmt.Errorf("Error fetching server metadata: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Creating snapshot: %v", c.SnapshotName))
	snapshot, err := client.CreateSnapshot(c.SnapshotName, c.OrganizationID, server.Volumes["0"].ID)
	if err != nil {
		err := fmt.Errorf("Error creating snapshot: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Printf("Snapshot image ID: %s", snapshot.ID)

	state.Put("snapshot_image_id", snapshot.ID)
	state.Put("snapshot_name", c.SnapshotName)

	return multistep.ActionContinue
}

func (s *stepSnapshot) Cleanup(state multistep.StateBag) {
	// no cleanup
}
