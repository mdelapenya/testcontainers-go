package dapr

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/testcontainers/testcontainers-go"
)

type options struct {
	AppName        string
	Components     map[string]Component
	ComponentsPath string
}

func defaultOptions() options {
	return options{
		AppName:        defaultDaprAppName,
		Components:     map[string]Component{},
		ComponentsPath: defaultComponentsPath,
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

// componentStruct {
type Component struct {
	Name     string
	Type     string
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

func NewComponent(name string, t string, metadata map[string]string) Component {
	return Component{
		Name:     name,
		Type:     t,
		Metadata: metadata,
	}
}

// WithComponents defines the components added to the dapr config, using a variadic list of Component.
func WithComponents(component ...Component) Option {
	return func(o *options) {
		for _, c := range component {
			o.Components[c.Key()] = c
		}
	}
}

// WithComponentsPath defines the container path where the components will be stored.
func WithComponentsPath(path string) Option {
	return func(o *options) {
		o.ComponentsPath = path
	}
}
