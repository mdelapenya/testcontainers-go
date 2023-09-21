package dapr

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
)

const (
	defaultDaprPort    string = "50001/tcp"
	defaultDaprAppName string = "dapr-app"
)

// DaprContainer represents the Dapr container type used in the module
type DaprContainer struct {
	testcontainers.Container
	Settings options
}

// RunContainer creates an instance of the Dapr container type
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*DaprContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "daprio/daprd:1.11.3",
		ExposedPorts: []string{defaultDaprPort},
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	settings := defaultOptions()
	for _, opt := range opts {
		if apply, ok := opt.(Option); ok {
			apply(&settings)
		}
		opt.Customize(&genericContainerReq)
	}

	genericContainerReq.Cmd = []string{"./daprd", "-app-id", settings.AppName, "--dapr-listen-addresses=0.0.0.0", "-components-path", "/components"}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	return &DaprContainer{Container: container, Settings: settings}, nil
}
