package onlinelabs

import (
	"fmt"
	"log"
	"time"
)

var (
	MaxWaitErrors = 5

	errInvalidServerState = fmt.Errorf("only 'up' and 'down' server states are supported")
)

func waitForServerState(desiredState, serverID string, client ClientInterface, timeout time.Duration) error {
	done := make(chan struct{})
	defer close(done)

	result := make(chan error, 1)
	go func() {
		errors := 0
		attempts := 0
		for {
			attempts += 1

			log.Printf("Checking server status... (attempt: %d)", attempts)
			server, err := client.GetServer(serverID)
			if err != nil {
				if errors < MaxWaitErrors {
					errors++
					log.Printf("Error caught while waiting for server status (%d out of %d): %v", errors, MaxWaitErrors, err)
					continue
				}
				result <- err
				return
			}

			if server.State == desiredState {
				result <- nil
				return
			}

			time.Sleep(3 * time.Second)

			select {
			case <-done:
				return
			default:
			}
		}
	}()

	log.Printf("Waiting for up to %d seconds for server to become %s", timeout/time.Second, desiredState)
	select {
	case err := <-result:
		return err
	case <-time.After(timeout):
		err := fmt.Errorf("Timeout while waiting to for server to become '%s'", desiredState)
		return err
	}
}
