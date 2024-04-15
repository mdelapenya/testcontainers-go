package buildkit_test

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/buildkit"
)

func TestGenericContainerFromDockerfileBuildkit(t *testing.T) {
	ctx := context.Background()

	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		t.Fatal(err)
	}

	testArg := "testFile"

	expectedTag := "test-repo:test-tag"

	// buildFromDockerfileWithBuildKit {
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
	// }
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

func TestBuildImageFromDockerfileBuildkit(t *testing.T) {
	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer provider.Close()

	cli := provider.Client()

	ctx := context.Background()

	testArg := "testFile"

	tag, err := provider.BuildImage(ctx, &testcontainers.ContainerRequest{
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
	})
	if err != nil {
		t.Fatal(err)
	}

	if tag != "test-repo:test-tag" {
		t.Fatalf("Expected tag %s, got %s", "test-repo:test-tag", tag)
	}

	_, _, err = cli.ImageInspectWithRaw(ctx, tag)
	if err != nil {
		t.Fatalf("Image %s should exist", tag)
	}

	t.Cleanup(func() {
		_, err := cli.ImageRemove(ctx, tag, types.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}
