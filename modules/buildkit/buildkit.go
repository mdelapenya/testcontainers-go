package buildkit

import (
	"context"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions"
	"github.com/moby/buildkit/session"

	"github.com/testcontainers/testcontainers-go"
)

const minBuildKitApiVersion = "1.39"

// BuildKitOptionsModifier is a function that modifies the build options to use BuildKit.
// It checks if the docker client supports BuildKit and if it does, it creates a new build session.
// You can use this function as a BuildOptionsModifier in the ContainerRequest.FromDockerfile
// to build images using BuildKit.
func BuildKitOptionsModifier(buildOptions *types.ImageBuildOptions) {
	ctx := context.Background()

	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		testcontainers.Logger.Printf("🛠️ Could not access the docker client: %s", err)
		return
	}
	defer cli.Close()

	clientApiVersion := cli.ClientVersion()

	if !versions.GreaterThanOrEqualTo(clientApiVersion, minBuildKitApiVersion) {
		testcontainers.Logger.Printf("🛠️ BuildKit is not supported by the docker client")
		return
	}

	s, err := session.NewSession(ctx, "testcontainers", "")
	if err != nil {
		testcontainers.Logger.Printf("🛠️ Could not create a build session, building without buildkit: %s", err)
		return
	}

	testcontainers.Logger.Printf("🛠️ Building using BuildKit")

	dialSession := func(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error) {
		return cli.DialHijack(ctx, "/session", proto, meta)
	}

	go func() {
		if err := s.Run(ctx, dialSession); err != nil {
			testcontainers.Logger.Printf("🛠️ Failed to run the build session: %s", err)
		}
	}()
	defer s.Close()

	buildOptions.SessionID = s.ID()
	buildOptions.Version = types.BuilderBuildKit
}
