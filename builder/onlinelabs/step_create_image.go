package onlinelabs

import (
	"fmt"
	"log"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type stepCreateImage struct{}

func (s *stepCreateImage) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(ClientInterface)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(Config)
	snapshotID := state.Get("snapshot_id").(string)
	serverArch := state.Get("server_arch").(string)

	ui.Say(fmt.Sprintf("Creating image: %v", c.ImageArtifactName))
	image, err := client.CreateImage(c.OrganizationID, c.ImageArtifactName, serverArch, snapshotID)
	if err != nil {
		err := fmt.Errorf("error creating image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Printf("Image ID: %s", image.ID)

	state.Put("image_id", image.ID)
	state.Put("image_name", image.Name)

	return multistep.ActionContinue
}

func (s *stepCreateImage) Cleanup(state multistep.StateBag) {
}
