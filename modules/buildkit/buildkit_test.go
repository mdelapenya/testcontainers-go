package buildkit_test

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/buildkit"
)

func TestBuildImageFromDockerfileBuildkit(t *testing.T) {
	ctx := context.Background()

	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		t.Fatal(err)
	}

	testArg := "testFile"

	expectedTag := "test-repo:test-tag"

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "testdata",
				Dockerfile: "buildx.Dockerfile",
				Repo:       "test-repo",
				Tag:        "test-tag",
				BuildArgs: map[string]*string{
					"FILENAME": &testArg,
				},
				PrintBuildLog:        true,
				BuildOptionsModifier: buildkit.BuildKitOptionsModifier,
			},
		},
		Started: false, // do not start the container
	})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = cli.ImageInspectWithRaw(ctx, expectedTag)
	if err != nil {
		t.Fatalf("Image %s should exist", expectedTag)
	}

	t.Cleanup(func() {
		if err := c.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})
}
