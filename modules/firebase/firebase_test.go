package firebase_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/testcontainers/testcontainers-go/modules/firebase"
)

func TestFirebase(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Minute))

	thing := fmt.Sprintf("%s/test", os.Getenv("PWD"))
	fmt.Println(thing)

	container, err := firebase.RunContainer(
		ctx,
		firebase.WithRoot(filepath.Join(BasePath(), "firebase")),
		firebase.WithCache(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
		cancel()
	})

	// perform assertions
	firestoreUrl, err := container.FirestoreConnectionString(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, firestoreUrl)

	authUrl, err := container.AuthConnectionString(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, authUrl)
}
