package onlinelabs

import "time"

type Server struct {
	StateDetail string `json:"state_detail"`
	Image       struct {
		DefaultBootscript interface{} `json:"default_bootscript"`
		CreationDate      time.Time   `json:"creation_date"`
		Name              string      `json:"name"`
		ModificationDate  time.Time   `json:"modification_date"`
		Organization      string      `json:"organization"`
		ExtraVolumes      string      `json:"extra_volumes"`
		Arch              string      `json:"arch"`
		ID                string      `json:"id"`
		RootVolume        struct {
			Size       int64  `json:"size"`
			ID         string `json:"id"`
			VolumeType string `json:"volume_type"`
			Name       string `json:"name"`
		} `json:"root_volume"`
		Public bool `json:"public"`
	} `json:"image"`
	CreationDate      time.Time   `json:"creation_date"`
	PublicIP          *IPAddress  `json:"public_ip"`
	PrivateIP         *NullString `json:"private_ip"`
	ID                string      `json:"id"`
	DynamicIPRequired bool        `json:"dynamic_ip_required"`
	ModificationDate  time.Time   `json:"modification_date"`
	EnableIpv6        bool        `json:"enable_ipv6"`
	Hostname          string      `json:"hostname"`
	State             string      `json:"state"`
	Bootscript        struct {
		Kernel       string `json:"kernel"`
		Initrd       string `json:"initrd"`
		Default      bool   `json:"default"`
		Bootcmdargs  string `json:"bootcmdargs"`
		Architecture string `json:"architecture"`
		Title        string `json:"title"`
		Dtb          string `json:"dtb"`
		Organization string `json:"organization"`
		ID           string `json:"id"`
		Public       bool   `json:"public"`
	} `json:"bootscript"`
	Location struct {
		PlatformID   string `json:"platform_id"`
		HypervisorID string `json:"hypervisor_id"`
		NodeID       string `json:"node_id"`
		ClusterID    string `json:"cluster_id"`
		ZoneID       string `json:"zone_id"`
	} `json:"location"`
	Ipv6           interface{}        `json:"ipv6"`
	CommercialType string             `json:"commercial_type"`
	Tags           []interface{}      `json:"tags"`
	Arch           string             `json:"arch"`
	ExtraNetworks  []interface{}      `json:"extra_networks"`
	Name           string             `json:"name"`
	Volumes        map[string]*Volume `json:"volume"`
	SecurityGroup  struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"security_group"`
	Organization string `json:"organization"`
}

type Image struct {
	DefaultBootscript *Bootscript `json:"default_bootscript,omitempty"`
	Arch              string      `json:"arch,omitempty"`
	CreationDate      string      `json:"creation_date,omitempty"`
	ExtraVolumes      string      `json:"extra_volumes,omitempty"` // is this a bug?? e.g.: "[]"
	FromImage         *NullString `json:"from_image,omitempty"`
	FromServer        *NullString `json:"from_server,omitempty"`
	ID                string      `json:"id"`
	MarketplaceKey    *NullString `json:"marketplace_key,omitempty"`
	ModificationDate  string      `json:"modification_date,omitempty"`
	Name              string      `json:"name"`
	Organization      string      `json:"organization,omitempty"`
	Public            bool        `json:"public"`
	RootVolume        *Volume     `json:"root_volume,omitempty"`
}

type Volume struct {
	Size             int64     `json:"size"`
	Name             string    `json:"name"`
	ModificationDate time.Time `json:"modification_date"`
	Organization     string    `json:"organization"`
	ExportURI        string    `json:"export_uri"`
	CreationDate     time.Time `json:"creation_date"`
	ID               string    `json:"id"`
	VolumeType       string    `json:"volume_type"`
	Server           struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"server"`
}

type AbbreviatedServer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Snapshot struct {
	BaseVolume   *Volume `json:"base_volume"`
	CreationDate string  `json:"creation_date"`
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Organization string  `json:"organization"`
	Size         uint64  `json:"size"`
	State        string  `json:"state"`
	VolumeType   string  `json:"volume_type"`
}

type Bootscript struct {
	Kernel       *Kernel      `json:"kernel"`
	Title        string       `json:"title"`
	Public       bool         `json:"public"`
	Initrd       *Initrd      `json:"initrd"`
	BootCmdArgs  *BootCmdArgs `json:"bootcmdargs"`
	Organization string       `json:"organization"`
	ID           string       `json:"id"`
}

type Kernel struct {
	Dtb   string `json:"dtb"`
	Path  string `json:"path"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Initrd struct {
	Path  string `json:"path"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

type BootCmdArgs struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type IPAddress struct {
	Dynamic bool   `json:"dynamic"`
	ID      string `json:"id"`
	Address string `json:"address"`
}

type createServerParams struct {
	Organization string             `json:"organization"`
	Name         string             `json:"name"`
	Image        string             `json:"image"`
	Tags         []string           `json:"tags,omitempty"`
	Volumes      map[string]*Volume `json:"volumes"`
}

type createSnapshotParams struct {
	Organization string `json:"organization"`
	Name         string `json:"name"`
	VolumeID     string `json:"volume_id"`
}

type createImageParams struct {
	Organization string `json:"organization"`
	Name         string `json:"name"`
	Arch         string `json:"arch"`
	RootVolume   string `json:"root_volume"`
}

type NullString struct {
	Value string
}

func (ns *NullString) String() string {
	return ns.Value
}

func (ns *NullString) UnmarshalJSON(j []byte) error {
	if string(j) == "null" {
		ns.Value = ""
	}
	ns.Value = string(j)
	return nil
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns.Value == "" {
		return []byte("null"), nil
	}

	return []byte(ns.Value), nil
}
