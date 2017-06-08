package onlinelabs

import (
	"fmt"
	"log"

	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
)

type stepCreateSnapshot struct{}

func (s *stepCreateSnapshot) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(ClientInterface)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	serverID := state.Get("server_id").(string)
	server, err := client.GetServer(serverID)
	if err != nil {
		err := fmt.Errorf("error fetching server metadata: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Creating snapshot: %v", c.SnapshotName))
	log.Println(server)
	snapshot, err := client.CreateSnapshot(c.SnapshotName, c.OrganizationID, server.Volumes["0"].ID)
	if err != nil {
		err := fmt.Errorf("error creating snapshot: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Printf("Snapshot ID: %s", snapshot.ID)

	state.Put("snapshot_id", snapshot.ID)
	state.Put("snapshot_name", c.SnapshotName)

	return multistep.ActionContinue
}

func (s *stepCreateSnapshot) Cleanup(state multistep.StateBag) {
	ok := false
	snapshotID := ""

	if snapshotID, ok = state.Get("snapshot_id").(string); !ok || snapshotID == "" {
		log.Printf("no snapshot id found; skipping cleanup")
		return
	}

	client := state.Get("client").(ClientInterface)
	err := client.DestroySnapshot(snapshotID)
	if err != nil {
		log.Printf("error destroying snapshot %v: %v", snapshotID, err)
	}
}
