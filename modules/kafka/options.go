package kafka

import (
	"github.com/testcontainers/testcontainers-go"
)

type options struct {
	// KafkaAuthenticationMethod is either "none" for plaintext or "sasl"
	// for SASL (scram) authentication.
	KafkaAuthenticationMethod string

	// Listeners is a list of custom listeners that can be provided to access the
	// containers form within docker networks
	Listeners []Listener
}

func defaultOptions() options {
	return options{
		KafkaAuthenticationMethod: "none",
		Listeners:                 make([]Listener, 0),
	}
}

// Compiler check to ensure that Option implements the testcontainers.ContainerCustomizer interface.
var _ testcontainers.ContainerCustomizer = (Option)(nil)

// Option is an option for the Kafka container.
type Option func(*options)

// Customize is a NOOP. It's defined to satisfy the testcontainers.ContainerCustomizer interface.
func (o Option) Customize(*testcontainers.GenericContainerRequest) error {
	// NOOP to satisfy interface.
	return nil
}

// WithListener adds a custom listener to the Kafka container. Listener
// will be aliases to all networks, so they can be accessed from within docker
// networks. At least one network must be attached to the container, if not an
// error will be thrown when starting the container.
func WithListener(l Listener) Option {
	return func(o *options) {
		o.Listeners = append(o.Listeners, l)
	}
}
