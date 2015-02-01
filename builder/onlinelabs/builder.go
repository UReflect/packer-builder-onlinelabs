package onlinelabs

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/common/uuid"
	"github.com/mitchellh/packer/packer"
)

const BuilderId = "meatballhat.onlinelabs"

type config struct {
	common.PackerConfig `mapstructure:",squash"`

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
	SSHUsername       string    `mapstructure:"ssh_username"`
	SSHPassword       string    `mapstructure:"ssh_password"`
	SSHPrivateKeyFile string    `mapstructure:"ssh_private_key_file"`
	SSHPort           uint      `mapstructure:"ssh_port"`

	RawSSHTimeout   string `mapstructure:"ssh_timeout"`
	RawStateTimeout string `mapstructure:"state_timeout"`

	sshTimeout   time.Duration
	stateTimeout time.Duration

	tpl *packer.ConfigTemplate
}

type Builder struct {
	config *config
	runner multistep.Runner
}

func NewBuilder() *Builder {
	return &Builder{
		config: &config{},
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
	md, err := common.DecodeConfig(&b.config, raws...)
	if err != nil {
		return nil, err
	}

	b.config.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return nil, err
	}
	b.config.tpl.UserVars = b.config.PackerUserVars

	errs := common.CheckUnusedConfig(md)

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
		b.config.SnapshotName = getenvDefault("ONLINELABS_SNAPSHOT_NAME", "packer-{{timestamp}}")
	}

	if b.config.SSHUsername == "" {
		b.config.SSHUsername = getenvDefault("ONLINELABS_SSH_USERNAME", "root")
	}

	if b.config.SSHPassword == "" {
		b.config.SSHPassword = os.Getenv("ONLINELABS_SSH_PASSWORD")
	}

	if b.config.SSHPrivateKeyFile == "" {
		b.config.SSHPrivateKeyFile = os.Getenv("ONLINELABS_SSH_PRIVATE_KEY_FILE")
	}

	if b.config.SSHPort == 0 {
		v, err := strconv.ParseUint(getenvDefault("ONLINELABS_SSH_PORT", "22"), 0, 64)
		if err == nil {
			b.config.SSHPort = uint(v)
		} else {
			b.config.SSHPort = 22
		}
	}

	if b.config.RawSSHTimeout == "" {
		b.config.RawSSHTimeout = getenvDefault("ONLINELABS_RAW_SSH_TIMEOUT", "1m")
	}

	if b.config.RawStateTimeout == "" {
		b.config.RawStateTimeout = getenvDefault("ONLINELABS_RAW_STATE_TIMEOUT", "6m")
	}

	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}

	templates := map[string]*string{
		"image_id":      &b.config.ImageID,
		"account_url":   &b.config.AccountURL,
		"api_url":       &b.config.APIURL,
		"api_token":     &b.config.APIToken,
		"snapshot_name": &b.config.SnapshotName,
		"server_name":   &b.config.ServerName,
		"ssh_username":  &b.config.SSHUsername,
		"ssh_timeout":   &b.config.RawSSHTimeout,
		"state_timeout": &b.config.RawStateTimeout,
	}

	for n, ptr := range templates {
		var err error
		*ptr, err = b.config.tpl.Process(*ptr, nil)
		if err != nil {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("Error processing %s: %s", n, err))
		}
	}

	if b.config.APIToken == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("an api_token must be specified"))
	}

	sshTimeout, err := time.ParseDuration(b.config.RawSSHTimeout)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing ssh_timeout: %s", err))
	}
	b.config.sshTimeout = sshTimeout

	stateTimeout, err := time.ParseDuration(b.config.RawStateTimeout)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing state_timeout: %s", err))
	}
	b.config.stateTimeout = stateTimeout

	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}

	common.ScrubConfig(b.config, b.config.APIToken)
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
		&common.StepConnectSSH{
			SSHAddress:     sshAddress,
			SSHConfig:      sshConfig,
			SSHWaitTimeout: 5 * time.Minute,
		},
		&common.StepProvision{},
		&stepShutdown{},
		&stepPowerOff{},
		&stepSnapshot{},
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

	if _, ok := state.GetOk("snapshot_name"); !ok {
		log.Println("Failed to find snapshot_name in state. Bug?")
		return nil, nil
	}

	artifact := &Artifact{
		id:     state.Get("snapshot_image_id").(string),
		name:   state.Get("snapshot_name").(string),
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
