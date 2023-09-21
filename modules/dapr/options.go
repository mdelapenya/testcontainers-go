package dapr

import (
	"github.com/testcontainers/testcontainers-go"
)

type options struct {
	AppName string
}

func defaultOptions() options {
	return options{
		AppName: defaultDaprAppName,
	}
}

// Compiler check to ensure that Option implements the testcontainers.ContainerCustomizer interface.
var _ testcontainers.ContainerCustomizer = (*Option)(nil)

// Option is an option for the Redpanda container.
type Option func(*options)

// Customize is a NOOP. It's defined to satisfy the testcontainers.ContainerCustomizer interface.
func (o Option) Customize(*testcontainers.GenericContainerRequest) {
	// NOOP to satisfy interface.
}

// WithAppName defines the app name added to the dapr config.
func WithAppName(name string) Option {
	return func(o *options) {
		o.AppName = name
	}
}
