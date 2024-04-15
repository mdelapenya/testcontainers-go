package buildkit_test

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/buildkit"
)

func ExampleBuildKitOptionsModifier() {
	ctx := context.Background()

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
		Started: false, // no need to start the container
	})
	// }
	if err != nil {
		log.Fatalf("failed to start container: %v", err)
	}
	defer func() {
		if err := c.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %v", err) // nolint:gocritic
		}
	}()

	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		log.Fatalf("Could not access the docker client: %s", err) // nolint:gocritic
	}
	defer cli.Close()

	inspect, _, err := cli.ImageInspectWithRaw(ctx, expectedTag)
	if err != nil {
		log.Fatalf("Image %s should exist", expectedTag) // nolint:gocritic
	}

	fmt.Printf("%v\n", inspect.RepoTags[0])
	fmt.Printf("%v\n", inspect.Comment)

	// Output:
	// test-repo:test-tag
	// buildkit.dockerfile.v0
}
