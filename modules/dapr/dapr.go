package dapr

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"

	"github.com/testcontainers/testcontainers-go"
)

const (
	// defaultComponentsPath is the path where the components are mounted in the Dapr container
	defaultComponentsPath string = "/components"
	defaultDaprPort       string = "50001/tcp"
	defaultDaprAppName    string = "dapr-app"
	// defaultDaprNetworkName is the name of the network created by the Dapr container, in which the app container is connected
	// and all the components will be attached to.
	defaultDaprNetworkName string = "dapr-network"
)

var (
	//go:embed mounts/component.yaml.tpl
	componentYamlTpl string

	// componentsTmpDir is the directory where the components are created before being mounted in the container
	componentsTmpDir string
)

// DaprContainer represents the Dapr container type used in the module
type DaprContainer struct {
	testcontainers.Container
	Network             testcontainers.Network
	ComponentContainers map[string]testcontainers.Container
	Settings            options
}

// GRPCPort returns the port used by the Dapr container
func (c *DaprContainer) GRPCPort(ctx context.Context) (int, error) {
	port, err := c.MappedPort(ctx, nat.Port(defaultDaprPort))
	if err != nil {
		return 0, err
	}

	return port.Int(), nil
}

// Terminate terminates the Dapr container and removes the component containers and the Dapr network,
// in that particular order.
func (c *DaprContainer) Terminate(ctx context.Context) error {
	if err := c.Container.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate Dapr container %w", err)
	}

	for key, componentContainer := range c.ComponentContainers {
		// do not terminate the component container if it has no image defined
		if c.Settings.Components[key].Image == "" {
			continue
		}

		if err := componentContainer.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate component container %w", err)
		}
	}

	if err := c.Network.Remove(ctx); err != nil {
		return fmt.Errorf("failed to terminate Dapr network %w", err)
	}

	return nil
}

// RunContainer creates an instance of the Dapr container type, creating the following elements:
// - a Dapr network
// - a Dapr container
// - a component container for each component defined in the options. The component must have an image defined.
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*DaprContainer, error) {
	componentsTmpDir = filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixMilli()), "components")
	err := os.MkdirAll(componentsTmpDir, 0o700)
	if err != nil {
		return nil, err
	}

	// make sure the temporary components directory is removed after the container is run.
	defer func() {
		_ = os.Remove(componentsTmpDir)
	}()

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

	// Transfer the components to the container in the form of a YAML file for each component
	if err := renderComponents(settings, &genericContainerReq); err != nil {
		return nil, err
	}

	genericContainerReq.Cmd = []string{"./daprd", "-app-id", settings.AppName, "--dapr-listen-addresses=0.0.0.0", "-components-path", defaultComponentsPath}

	nw, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name: settings.NetworkName,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Dapr network: %w", err)
	}

	// attach Dapr container to the Dapr network
	genericContainerReq.Networks = []string{settings.NetworkName}
	// setting the network alias to the application name will make it easier to connect to the Dapr container
	genericContainerReq.NetworkAliases = map[string][]string{
		settings.NetworkName: {settings.AppName},
	}

	daprContainer, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	// we must start the component containers in container mode, so that they can connect to the Dapr container
	networkMode := fmt.Sprintf("container:%v", daprContainer.GetContainerID())

	componentContainers := map[string]testcontainers.Container{}
	for _, component := range settings.Components {
		if component.Image == "" {
			continue
		}

		componentContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:    component.Image,
				Networks: []string{settings.NetworkName},
				NetworkAliases: map[string][]string{
					settings.NetworkName: {component.Name},
				},
				HostConfigModifier: func(hc *container.HostConfig) {
					hc.NetworkMode = container.NetworkMode(networkMode)
				},
			},
			Started: true,
		})
		if err != nil {
			return nil, err
		}

		componentContainers[component.Key()] = componentContainer
	}

	return &DaprContainer{
		Container:           daprContainer,
		Settings:            settings,
		ComponentContainers: componentContainers,
		Network:             nw,
	}, nil
}

// renderComponents renders the configuration file for each component, creating a temporary file for each one under a default
// temporary directory. The entire directory is then uploaded to the container, including the
// right permissions (0o777) for Dapr to access the files.
func renderComponents(settings options, req *testcontainers.GenericContainerRequest) error {
	execPermissions := os.FileMode(0o777)

	for _, component := range settings.Components {
		content, err := component.Render()
		if err != nil {
			return err
		}

		tmpComponentFile := filepath.Join(componentsTmpDir, component.FileName())
		err = os.WriteFile(tmpComponentFile, content, execPermissions)
		if err != nil {
			return err
		}

	}

	req.Files = append(req.Files, testcontainers.ContainerFile{
		HostFilePath:      componentsTmpDir,
		ContainerFilePath: defaultComponentsPath,
		FileMode:          int64(execPermissions),
	})

	return nil
}
