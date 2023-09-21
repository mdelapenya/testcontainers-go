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

	daprContainer, err := dapr.RunContainer(ctx, testcontainers.WithImage("daprio/daprd:1.11.3"))
	if err != nil {
		panic(err)
	}

	// Clean up the container after
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
