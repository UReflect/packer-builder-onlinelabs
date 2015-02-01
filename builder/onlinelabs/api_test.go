package onlinelabs

import (
	"encoding/json"
	"fmt"
	"testing"
)

const (
	serverJSON = `{
	"bootscript": null,
	"dynamic_public_ip": false,
	"id": "741db378-6b87-46d4-a8c5-4e46a09ab1f8",
	"image": {
	 "id": "85917034-46b0-4cc5-8b48-f0a2245e357e",
	 "name": "ubuntu working"
	},
	"name": "my_server",
	"organization": "000a115d-2852-4b0a-9ce8-47f1134ba95a",
	"private_ip": null,
	"public_ip": null,
	"state": "running",
	"tags": [
	 "test",
	 "www"
	],
	"volumes": {
	 "0": {
		 "export_uri": null,
		 "id": "c1eb8f3a-4f0b-4b95-a71c-93223e457f5a",
		 "name": "vol simple snapshot",
		 "organization": "000a115d-2852-4b0a-9ce8-47f1134ba95a",
		 "server": {
			 "id": "741db378-6b87-46d4-a8c5-4e46a09ab1f8",
			 "name": "my_server"
		 },
		 "size": 10000000000,
		 "volume_type": "l_hdd"
	 }
	}
}`
)

func TestServerSerialization(t *testing.T) {
	srv := &Server{}
	err := json.Unmarshal([]byte(serverJSON), srv)
	if err != nil {
		t.Error(err)
	}

	//	if srv.Bootscript != nil {
	//		t.Fatal("Bootscript != nil")
	//	}

	if srv.DynamicPublicIP != false {
		t.Fatal("DynamicPublicIP != false")
	}

	if srv.ID != "741db378-6b87-46d4-a8c5-4e46a09ab1f8" {
		t.Fatal("ID != expected")
	}

	if srv.Image == nil {
		t.Fatal("Image was not deserialized")
	}

	if srv.Image.ID != "85917034-46b0-4cc5-8b48-f0a2245e357e" {
		t.Fatal("Image.ID != expected")
	}

	if srv.Image.Name != "ubuntu working" {
		t.Fatal("Image.Name != expected")
	}

	if srv.Name != "my_server" {
		t.Fatal("Name != my_server")
	}

	if srv.Organization != "000a115d-2852-4b0a-9ce8-47f1134ba95a" {
		t.Fatal("Organization != expected")
	}

	if srv.PrivateIP != nil {
		t.Fatal("PrivateIP != nil")
	}

	if srv.PublicIP != nil {
		t.Fatal("PublicIP != nil")
	}

	if srv.State != "running" {
		t.Fatal("State != running")
	}

	if len(srv.Tags) != 2 {
		t.Fatal("len(Tags) != 2")
	}

	if srv.Tags[0] != "test" {
		t.Fatal("Tags[0] != test")
	}

	if srv.Tags[1] != "www" {
		t.Fatal("Tags[1] != www")
	}

	if srv.Volumes == nil {
		t.Fatal("Volumes == nil")
	}

	if len(srv.Volumes) != 1 {
		t.Fatal("len(Volumes) != 1")
	}

	vol0, ok := srv.Volumes["0"]
	if !ok {
		t.Fatal("Volumes[\"0\"] not present")
	}

	if vol0.ExportURI != nil {
		t.Fatal("vol0.ExportURI != nil")
	}

	if vol0.ID != "c1eb8f3a-4f0b-4b95-a71c-93223e457f5a" {
		t.Fatal("vol0.ID != expected")
	}

	if vol0.Name != "vol simple snapshot" {
		t.Fatal("vol0.Name != expected")
	}

	if vol0.Organization != "000a115d-2852-4b0a-9ce8-47f1134ba95a" {
		t.Fatal("vol0.Organization != expected")
	}

	if vol0.Server == nil {
		t.Fatal("vol0.Server == nil")
	}

	if vol0.Server.ID != "741db378-6b87-46d4-a8c5-4e46a09ab1f8" {
		t.Fatal("vol0.Server.ID != expected")
	}

	if vol0.Server.Name != "my_server" {
		t.Fatal("vol0.Server.Name != my_server")
	}

	if vol0.Size != uint64(10000000000) {
		t.Fatal("vol0.Size != expected")
	}

	if vol0.VolumeType != "l_hdd" {
		t.Fatal("vol0.VolumeType != l_hdd")
	}

	b, err := json.MarshalIndent(srv, "", "  ")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf(string(b))
}
