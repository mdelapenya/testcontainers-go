package eventstoredb

import (
	"context"
	"fmt"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	httpPort = "2113/tcp"
	tcpPort  = "1113/tcp"
)

// Container represents the EventStoreDB container type used in the module.
type Container struct {
	testcontainers.Container
	insecure bool
}

// WithInsecure sets EVENTSTORE_INSECURE=true on the container. EventStoreDB runs
// in insecure mode by default; this option is provided for explicitness.
func WithInsecure() testcontainers.CustomizeRequestOption {
	return testcontainers.WithEnv(map[string]string{
		"EVENTSTORE_INSECURE": "true",
	})
}

// Run creates an instance of the EventStoreDB container type.
func Run(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*Container, error) {
	moduleOpts := make([]testcontainers.ContainerCustomizer, 0, 3+len(opts))
	moduleOpts = append(moduleOpts,
		testcontainers.WithExposedPorts(httpPort, tcpPort),
		testcontainers.WithEnv(map[string]string{
			"EVENTSTORE_INSECURE":                  "true",
			"EVENTSTORE_ENABLE_ATOM_PUB_OVER_HTTP": "true",
		}),
		testcontainers.WithWaitStrategy(
			wait.ForHTTP("/health/live").
				WithPort(httpPort).
				WithStatusCodeMatcher(func(status int) bool {
					return status == 204
				}),
		),
	)
	moduleOpts = append(moduleOpts, opts...)

	ctr, err := testcontainers.Run(ctx, img, moduleOpts...)
	var c *Container
	if ctr != nil {
		c = &Container{Container: ctr}
	}

	if err != nil {
		return c, fmt.Errorf("run eventstoredb: %w", err)
	}

	// Inspect the container environment to determine insecure mode.
	inspect, err := ctr.Inspect(ctx)
	if err != nil {
		return c, fmt.Errorf("inspect eventstoredb: %w", err)
	}

	for _, env := range inspect.Config.Env {
		if v, ok := strings.CutPrefix(env, "EVENTSTORE_INSECURE="); ok {
			c.insecure = strings.EqualFold(v, "true")
			break
		}
	}

	return c, nil
}

// ConnectionString returns the gRPC connection string for the EventStoreDB container.
// When running in insecure mode it returns "esdb://host:port?tls=false";
// otherwise it returns "esdb+discover://host:port".
func (c *Container) ConnectionString(ctx context.Context) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("get host: %w", err)
	}

	port, err := c.MappedPort(ctx, httpPort)
	if err != nil {
		return "", fmt.Errorf("get mapped port: %w", err)
	}

	if c.insecure {
		return fmt.Sprintf("esdb://%s:%s?tls=false", host, port.Port()), nil
	}

	return fmt.Sprintf("esdb+discover://%s:%s", host, port.Port()), nil
}
