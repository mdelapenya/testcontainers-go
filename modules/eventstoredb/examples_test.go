package eventstoredb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/eventstoredb"
)

func ExampleRun() {
	// runEventStoreDBContainer {
	ctx := context.Background()

	eventstoredbContainer, err := eventstoredb.Run(ctx, "eventstore/eventstore:24.2")
	defer func() {
		if err := testcontainers.TerminateContainer(eventstoredbContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}
	// }

	state, err := eventstoredbContainer.State(ctx)
	if err != nil {
		log.Printf("failed to get container state: %s", err)
		return
	}

	fmt.Println(state.Running)

	// Output:
	// true
}

func ExampleRun_connectionString() {
	// connectionString {
	ctx := context.Background()

	eventstoredbContainer, err := eventstoredb.Run(ctx, "eventstore/eventstore:24.2",
		eventstoredb.WithInsecure(),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(eventstoredbContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	connStr, err := eventstoredbContainer.ConnectionString(ctx)
	if err != nil {
		log.Printf("failed to get connection string: %s", err)
		return
	}
	// }

	_ = connStr
	fmt.Println("connected")

	// Output:
	// connected
}
