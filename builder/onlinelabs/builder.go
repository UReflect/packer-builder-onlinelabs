package onlinelabs

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/common/uuid"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/mitchellh/multistep"
)

const BuilderId = "meatballhat.onlinelabs"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	AccountURL string `mapstructure:"account_url"`
	APIURL     string `mapstructure:"api_url"`
	APIToken   string `mapstructure:"api_token"`

	ImageID        string `mapstructure:"image_id"`
	OrganizationID string `mapstructure:"organization_id"`

	ServerName        string    `mapstructure:"server_name"`
	ServerTags        []string  `mapstructure:"server_tags"`
	ServerVolumes     []*Volume `mapstructure:"volumes"`
	DynamicPublicIP   bool      `mapstructure:"dynamic_public_ip"`
	SnapshotName      string    `mapstructure:"snapshot_name"`
	ImageArtifactName string    `mapstructure:"image_artifact_name"`
	InstanceType      string    `mapstructure:"instance_type"`
	Region            string    `mapstructure:"region"`
	SourceImage       string    `mapstructure:"source_image"`

	RawStateTimeout string `mapstructure:"state_timeout"`
	StateTimeout    time.Duration

	ctx interpolate.Context
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func NewBuilder() *Builder {
	return &Builder{
		config: Config{},
		runner: nil,
	}
}

func getenvDefault(key, dflt string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return dflt
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
	}, raws...)
	if err != nil {
		return nil, err
	}

	if b.config.AccountURL == "" {
		b.config.AccountURL = getenvDefault("ONLINELABS_ACCOUNT_URL", AccountURL.String())
	}

	if b.config.APIURL == "" {
		b.config.APIURL = getenvDefault("ONLINELABS_API_URL", APIURL.String())
	}

	if b.config.APIToken == "" {
		b.config.APIToken = os.Getenv("ONLINELABS_API_TOKEN")
	}

	if b.config.ImageID == "" {
		b.config.ImageID = os.Getenv("ONLINELABS_IMAGE_ID")
	}

	if b.config.OrganizationID == "" {
		b.config.OrganizationID = os.Getenv("ONLINELABS_ORGANIZATION_ID")
	}

	if b.config.ServerName == "" {
		b.config.ServerName = getenvDefault("ONLINELABS_SERVER_NAME", fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID()))
	}

	if b.config.SnapshotName == "" {
		b.config.SnapshotName = getenvDefault("ONLINELABS_SNAPSHOT_NAME", "packer-snapshot-{{ timestamp }}")
	}

	if b.config.ImageArtifactName == "" {
		b.config.ImageArtifactName = getenvDefault("ONLINELABS_IMAGE_ARTIFACT_NAME", "packer-image-{{ timestamp }}")
	}

	if b.config.Comm.SSHUsername == "" {
		b.config.Comm.SSHUsername = getenvDefault("ONLINELABS_SSH_USERNAME", "root")
	}

	if b.config.RawStateTimeout == "" {
		b.config.RawStateTimeout = getenvDefault("ONLINELABS_RAW_STATE_TIMEOUT", "6m")
	}

	var errs *packer.MultiError
	errs = packer.MultiErrorAppend(errs, b.config.Comm.Prepare(&b.config.ctx)...)

	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}

	newTags := []string{}
	log.Println(b.config.ServerTags)
	for _, v := range b.config.ServerTags {
		// Check errors here
		newTags = append(newTags, v)
	}

	b.config.ServerTags = newTags

	if b.config.APIToken == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("an api_token must be specified"))
	}

	stateTimeout, err := time.ParseDuration(b.config.RawStateTimeout)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing state_timeout: %s", err))
	}
	b.config.StateTimeout = stateTimeout

	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}

	log.Println(common.ScrubConfig(b.config, b.config.APIToken))

	return nil, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	accountURL := AccountURL
	apiURL := APIURL
	if u, err := url.Parse(b.config.AccountURL); err == nil {
		accountURL = u
	}
	if u, err := url.Parse(b.config.APIURL); err == nil {
		apiURL = u
	}
	client := NewClient(b.config.APIToken, b.config.OrganizationID, accountURL, apiURL)

	state := &multistep.BasicStateBag{}
	state.Put("config", b.config)
	state.Put("client", client)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		// TODO: &stepCreateSSHKey{},
		&stepCreateServer{},
		&stepStartServer{},
		&stepServerInfo{},
		&communicator.StepConnect{
			Host:      sshAddress,
			SSHConfig: sshConfig,
			Config:    &b.config.Comm,
		},
		&common.StepProvision{},
		&stepShutdown{},
		&stepPowerOff{},
		&stepCreateSnapshot{},
		&stepCreateImage{},
	}

	// Run the steps
	if b.config.PackerDebug {
		b.runner = &multistep.DebugRunner{
			Steps:   steps,
			PauseFn: common.MultistepDebugFn(ui),
		}
	} else {
		b.runner = &multistep.BasicRunner{Steps: steps}
	}

	b.runner.Run(state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	if _, ok := state.GetOk("image_name"); !ok {
		log.Println("Failed to find image_name in state. Bug?")
		return nil, nil
	}

	artifact := &Artifact{
		id:     state.Get("image_id").(string),
		name:   state.Get("image_name").(string),
		client: client,
	}

	return artifact, nil
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}
