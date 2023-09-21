package dapr_test

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/dapr"
)

func ExampleRunContainer() {
	// runDaprContainer {
	ctx := context.Background()

	daprContainer, err := dapr.RunContainer(ctx,
		testcontainers.WithImage("daprio/daprd:1.11.3"),
		dapr.WithAppName("dapr-app"),
		dapr.WithComponents(
			dapr.NewComponent("pubsub", "pubsub.in-memory", map[string]string{"foo": "bar", "bar": "baz"}),
			dapr.NewComponent("statestore", "statestore.in-memory", map[string]string{"baz": "qux", "quux": "quuz"}),
		),
		dapr.WithComponentsPath("/components"),
	)
	if err != nil {
		panic(err)
	}

	// Clean up the container
	defer func() {
		if err := daprContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}()
	// }

	state, err := daprContainer.State(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(state.Running)

	// Output:
	// true
}
