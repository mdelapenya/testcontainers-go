package dapr

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/testcontainers/testcontainers-go"
)

type options struct {
	AppName     string
	Components  map[string]Component
	NetworkName string
}

// defaultOptions returns the default options for the Dapr container, including an in-memory state store.
func defaultOptions() options {
	inMemoryStore := NewComponent("statestore", "state.in-memory", map[string]string{})

	return options{
		AppName: defaultDaprAppName,
		Components: map[string]Component{
			inMemoryStore.Key(): inMemoryStore,
		},
		NetworkName: defaultDaprNetworkName,
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

// WithNetworkName defines the network name in which the dapr container is attached.
func WithNetworkName(name string) Option {
	return func(o *options) {
		o.NetworkName = name
	}
}

// componentStruct {
type Component struct {
	Name     string
	Type     string
	Image    string
	Metadata map[string]string
}

// }

// Key returns the component name, which must be unique.
func (c Component) Key() string {
	return c.Name
}

// FileName returns the component file name, which must be unique.
func (c Component) FileName() string {
	return c.Name + ".yaml"
}

// Render returns the component configuration as a byte slice, obtained after the interpolation
// of the component template.
func (c Component) Render() ([]byte, error) {
	tpl, err := template.New(c.FileName()).Parse(componentYamlTpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse component file template: %w", err)
	}

	var componentConfig bytes.Buffer
	if err := tpl.Execute(&componentConfig, c); err != nil {
		return nil, fmt.Errorf("failed to render component template: %w", err)
	}

	return componentConfig.Bytes(), nil
}

// NewComponentWithImage returns a new Component without its Docker image.
// Those components without a Docker image won't be run as a separate container in the Dapr network.
func NewComponent(name string, t string, metadata map[string]string) Component {
	return Component{
		Name:     name,
		Type:     t,
		Metadata: metadata,
	}
}

// NewComponentWithImage returns a new Component including its Docker image.
// Those components with a Docker image will be run as a separate container in the Dapr network,
// and using Dapr's container ID as the container network mode.
func NewComponentWithImage(name string, t string, image string, metadata map[string]string) Component {
	c := NewComponent(name, t, metadata)

	c.Image = image

	return c
}

// WithComponents defines the components added to the dapr config, using a variadic list of Component.
func WithComponents(component ...Component) Option {
	return func(o *options) {
		for _, c := range component {
			o.Components[c.Key()] = c
		}
	}
}
