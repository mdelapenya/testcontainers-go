package testcontainers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/internal/config"
)

const (
	dockerSock         = "unix:///var/run/docker.sock"
	tcpDockerHost1234  = "tcp://127.0.0.1:1234"
	tcpDockerHost33293 = "tcp://127.0.0.1:33293"
	tcpDockerHost4711  = "tcp://127.0.0.1:4711"
)

// unset environment variables to avoid side effects
// execute this function before each test
func resetTestEnv(t *testing.T) {
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "")
	t.Setenv("TESTCONTAINERS_RYUK_CONTAINER_PRIVILEGED", "")
}

func TestReadConfig(t *testing.T) {
	resetTestEnv(t)

	t.Run("Config is read just once", func(t *testing.T) {
		t.Setenv("HOME", "")
		t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

		cfg := ReadConfig()

		expected := TestcontainersConfig{
			Config: config.Config{
				RyukDisabled: true,
				Host:         dockerSock,
			},
		}

		assert.Equal(t, expected, cfg)

		t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "false")
		cfg = ReadConfig()
		assert.Equal(t, expected, cfg)
	})
}
