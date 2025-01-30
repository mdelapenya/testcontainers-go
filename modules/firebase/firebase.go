package firebase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/mount"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// FirebaseContainer represents the Firebase container type used in the module
type FirebaseContainer struct {
	testcontainers.Container
}

const defaultImageName = "ghcr.io/u-health/docker-firebase-emulator:13.29.2"

// WithRoot sets the directory which is copied to the destination container as firebase root
func WithRoot(rootPath string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		if !strings.HasSuffix(rootPath, "/firebase") {
			return fmt.Errorf("root path must end with '/firebase': %s", rootPath)
		}
		req.Files = append(req.Files, testcontainers.ContainerFile{
			HostFilePath:      rootPath,
			ContainerFilePath: "/srv/firebase",
		})

		return nil
	}
}

// WithData names the data directory in firebase root
func WithData(dataPath string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["DATA_DIRECTORY"] = dataPath
		return nil
	}
}

func cache(volumeName string, volumeOptions *mount.VolumeOptions) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		m := testcontainers.ContainerMount{
			Source: testcontainers.DockerVolumeMountSource{
				Name:          volumeName,
				VolumeOptions: volumeOptions,
			},
			Target: "/root/.cache/firebase",
		}
		req.Mounts = append(req.Mounts, m)
		return nil
	}
}

// WithCache enables firebase binary cache based on session (meaningful only when multiple tests are used)
func WithCache() testcontainers.CustomizeRequestOption {
	volumeName := fmt.Sprintf("firestore-cache-%s", testcontainers.SessionID())
	volumeOptions := &mount.VolumeOptions{
		Labels: testcontainers.GenericLabels(),
	}

	return cache(volumeName, volumeOptions)
}

// RunContainer creates an instance of the Firebase container type
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*FirebaseContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: defaultImageName,
		ExposedPorts: []string{
			UiPort,
			HubPort,
			LoggingPort,
			FunctionsPort,
			FirestorePort,
			PubsubPort,
			DatabasePort,
			AuthPort,
			StoragePort,
			HostingPort,
		},

		Env: map[string]string{},
		WaitingFor: wait.ForAll(
			wait.ForHTTP("/").WithPort(UiPort).WithStartupTimeout(3 * time.Minute),
		),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		if err := opt.Customize(&genericContainerReq); err != nil {
			return nil, fmt.Errorf("customize: %w", err)
		}
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	return &FirebaseContainer{Container: container}, nil
}
