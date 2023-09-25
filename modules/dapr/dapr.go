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
}

// RunContainer creates an instance of the Dapr container type
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*DaprContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "daprio/daprd:1.11.3",
		ExposedPorts: []string{defaultDaprPort},
		Cmd:          []string{"./daprd", "-app-id", defaultDaprAppName, "--dapr-listen-addresses=0.0.0.0", "-components-path", "/components"},
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		opt.Customize(&genericContainerReq)
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	return &DaprContainer{Container: container}, nil
}
