package onlinelabs

import (
	"fmt"
	"log"
)

type Artifact struct {
	id   string
	name string

	client ClientInterface
}

func (*Artifact) BuilderId() string {
	return BuilderId
}

func (a *Artifact) Id() string {
	return a.id
}

func (*Artifact) Files() []string {
	return nil
}

func (a *Artifact) String() string {
	return fmt.Sprintf("An image was created: '%v' (%v)", a.id, a.name)
}

func (a *Artifact) State(name string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: %s (%s)", a.id, a.name)
	return a.client.DestroyImage(a.id)
}
