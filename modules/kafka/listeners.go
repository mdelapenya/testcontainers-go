package kafka

import (
	"fmt"
	"math"

	"github.com/testcontainers/testcontainers-go"
)

type Listener struct {
	Name                 string
	Address              string
	Port                 int
	AuthenticationMethod string
}

// NewListener creates a new listener with the provided name, address and port.
func NewListener(name, address string, port int) Listener {
	return Listener{
		Name:                 name,
		Address:              address,
		Port:                 port,
		AuthenticationMethod: "none",
	}
}

// String returns a string representation of the listener. E.g. LISTENER_BOB://localhost:9092
func (l Listener) String() string {
	return fmt.Sprintf("%s://%s:%d", l.Name, l.Address, l.Port)
}

// Parse validates the listener configuration:
// - port must be between 0 and 65535
func (l Listener) Parse() error {
	if l.Port < 0 || l.Port > math.MaxUint16 {
		return fmt.Errorf("invalid port on listener %s:%d (must be between 0 and 65535)", l.Address, l.Port)
	}
	return nil
}

// registerListeners validates that the provided listeners are valid and set network aliases for the provided addresses.
// It will set the KAFKA_ADVERTISED_LISTENERS environment variable with the provided listeners.
// This method will check that the container is attached to at least one network and that network aliases are defined.
func registerListeners(settings options, req *testcontainers.GenericContainerRequest) error {
	if len(settings.Listeners) == 0 {
		return nil
	}

	if len(req.Networks) == 0 {
		return fmt.Errorf("container must be attached to at least one network")
	}

	if len(req.NetworkAliases) == 0 {
		return fmt.Errorf("container must have network aliases defined")
	}

	advertisedListeners := req.Env["KAFKA_ADVERTISED_LISTENERS"]

	for _, listener := range settings.Listeners {
		if err := listener.Parse(); err != nil {
			return err
		}

		if advertisedListeners == "" {
			advertisedListeners = listener.String()
		} else {
			advertisedListeners += "," + listener.String()
		}

		for _, network := range req.Networks {
			req.NetworkAliases[network] = append(req.NetworkAliases[network], listener.Address)
		}
	}

	req.Env["KAFKA_ADVERTISED_LISTENERS"] = advertisedListeners

	return nil
}
