package firebase_test

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go/modules/firebase"
)

func ExampleRunContainer() {
	// runFirebaseContainer {
	ctx := context.Background()

	firebaseContainer, err := firebase.RunContainer(
		ctx,
		firebase.WithRoot(filepath.Join(BasePath(), "firebase")),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container
	defer func() {
		if err := firebaseContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err) // nolint:gocritic
		}
	}()
	// }

	state, err := firebaseContainer.State(ctx)
	if err != nil {
		log.Fatalf("failed to get container state: %s", err) // nolint:gocritic
	}

	fmt.Println(state.Running)

	// Output:
	// true
}
