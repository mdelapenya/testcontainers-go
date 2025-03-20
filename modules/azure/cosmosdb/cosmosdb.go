package cosmosdb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultImage      = "mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:latest"
	defaultPort       = "8081/tcp"
	defaultPartitions = 10

	// EmulatorCredentials is the only-well known key for the emulator
	EmulatorCredentials = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

	connectionStringFormat = "AccountEndpoint=%s/;AccountKey=%s;"
)

type Container struct {
	testcontainers.Container
	pemBytes []byte
}

// Run creates an instance of the CosmosDB container type
func Run(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        img,
		ExposedPorts: []string{defaultPort},
		Env:          map[string]string{},
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	// 1. Gather all config options
	for _, opt := range opts {
		if err := opt.Customize(&genericContainerReq); err != nil {
			return nil, fmt.Errorf("customize: %w", err)
		}
	}

	waitStragies := []wait.Strategy{
		wait.ForListeningPort(defaultPort).SkipInternalCheck(),
	}

	// 1. Set the wait strategy for the emulator
	// buffer for the certificate bytes
	var certificateBuffer bytes.Buffer

	if strings.HasSuffix(img, "vnext-preview") {
		genericContainerReq.Env["LOG_LEVEL"] = "debug"

		certPath := "/home/cosmosdev/cosmosdev.crt"
		if genericContainerReq.Env["CERT_PATH"] != "" {
			certPath = genericContainerReq.Env["CERT_PATH"]
		}

		waitStragies = append(waitStragies,
			wait.ForFile(certPath).WithMatcher(func(r io.Reader) error {
				if _, err := io.Copy(&certificateBuffer, r); err != nil {
					return err
				}
				return nil
			}),
		)
	} else {
		// The certificate must be available as an HTTPS endpoint, and we extract it the certificate
		// bytes from there.
		waitStragies = append(waitStragies,
			wait.ForHTTP("/_explorer/emulator.pem").
				WithPort(defaultPort).
				WithStartupTimeout(2*time.Minute).
				WithTLS(true).
				WithAllowInsecure(true).
				WithStatusCodeMatcher(func(status int) bool {
					return status == http.StatusOK
				}).
				WithResponseMatcher(func(r io.Reader) bool {
					if _, err := io.Copy(&certificateBuffer, r); err != nil {
						return false
					}
					return true
				}),
		)
	}

	genericContainerReq.WaitingFor = wait.ForAll(waitStragies...)

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	var c *Container
	if container != nil {
		c = &Container{Container: container, pemBytes: certificateBuffer.Bytes()}
	}

	if err != nil {
		return c, fmt.Errorf("generic container: %w", err)
	}

	return c, nil
}

// ConnectionString returns the connection string of the emulator
func (c *Container) ConnectionString(ctx context.Context) (string, error) {
	emulatorURL, err := c.EmulatorURL(ctx)
	if err != nil {
		return "", fmt.Errorf("port endpoint: %w", err)
	}

	return fmt.Sprintf(connectionStringFormat, emulatorURL, EmulatorCredentials), nil
}

// EmulatorURL returns the URL of the emulator
func (c *Container) EmulatorURL(ctx context.Context) (string, error) {
	return c.PortEndpoint(ctx, defaultPort, "https")
}

// Certificate returns the PEM bytes of the emulator certificate
func (c *Container) Certificate() []byte {
	return c.pemBytes
}
