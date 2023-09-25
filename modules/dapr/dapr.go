package dapr

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
)

const (
	// defaultComponentsPath is the path where the components are mounted in the Dapr container
	defaultComponentsPath string = "/components"
	defaultDaprPort       string = "50001/tcp"
	defaultDaprAppName    string = "dapr-app"
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
	Settings options
}

// GRPCPort returns the port used by the Dapr container
func (c *DaprContainer) GRPCPort(ctx context.Context) (int, error) {
	port, err := c.MappedPort(ctx, nat.Port(defaultDaprPort))
	if err != nil {
		return 0, err
	}

	return port.Int(), nil
}

// RunContainer creates an instance of the Dapr container type
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

	opts = append(opts, WithComponents(NewComponent("statestore", "state.in-memory", map[string]string{})))

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

	genericContainerReq.Cmd = []string{"./daprd", "-app-id", settings.AppName, "--dapr-listen-addresses=0.0.0.0", "-components-path", settings.ComponentsPath}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	return &DaprContainer{Container: container, Settings: settings}, nil
}

// renderComponents renders the configuration file for each component, creating a temporary file for each one under a default
// temporary directory. The entire directory is then uploaded to the container.
func renderComponents(settings options, req *testcontainers.GenericContainerRequest) error {
	for _, component := range settings.Components {
		content, err := component.Render()

		tmpComponentFile := filepath.Join(componentsTmpDir, component.FileName())
		err = os.WriteFile(tmpComponentFile, content, 0o600)
		if err != nil {
			return err
		}

	}

	req.Files = append(req.Files, testcontainers.ContainerFile{
		HostFilePath:      componentsTmpDir,
		ContainerFilePath: settings.ComponentsPath,
		FileMode:          0o600,
	})

	return nil
}
