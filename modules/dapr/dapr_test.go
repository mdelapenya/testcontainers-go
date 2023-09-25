package dapr

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/testcontainers/testcontainers-go"
)

func TestDapr(t *testing.T) {
	ctx := context.Background()

	container, err := RunContainer(ctx, testcontainers.WithImage("daprio/daprd:1.11.3"))
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	// verify that the dapr network is created and the container is attached to it
	cli, err := testcontainers.NewDockerClientWithOpts(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	response, err := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.Arg("name", container.Settings.NetworkName)),
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(response) != 1 {
		t.Fatalf("expected 1 network, got %d", len(response))
	}

	daprNetwork := response[0]

	if daprNetwork.Name != container.Settings.NetworkName {
		t.Fatalf("expected network name %s, got %s", container.Settings.NetworkName, daprNetwork.Name)
	}

	// verify that the container is attached to the dapr network
	nws, err := container.Networks(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(nws) != 1 {
		t.Fatalf("expected 1 network, got %d", len(nws))
	}

	if nws[0] != container.Settings.NetworkName {
		t.Fatalf("expected network name %s, got %s", container.Settings.NetworkName, nws[0])
	}
}
