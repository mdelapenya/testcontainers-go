package kafka_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/network"
)

func ExampleRunContainer() {
	// runKafkaContainer {
	ctx := context.Background()

	kafkaContainer, err := kafka.RunContainer(ctx,
		kafka.WithClusterID("test-cluster"),
		testcontainers.WithImage("confluentinc/confluent-local:7.5.0"),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container after
	defer func() {
		if err := kafkaContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
	// }

	state, err := kafkaContainer.State(ctx)
	if err != nil {
		log.Fatalf("failed to get container state: %s", err) // nolint:gocritic
	}

	fmt.Println(kafkaContainer.ClusterID)
	fmt.Println(state.Running)

	// Output:
	// test-cluster
	// true
}

func ExampleRunContainer_withListeners() {
	ctx := context.Background()

	// 1. Create network
	kafkaNw, err := network.New(ctx)
	if err != nil {
		log.Fatalf("failed to create network: %s", err)
	}

	const (
		kafkaTopic        string = "msgs"
		kcatMessage       string = "Message produced by kcat"
		kcatMsgsFile      string = "/tmp/msgs.txt"
		kafkaAliasName    string = "kafka0"
		kafkaListenerPort int    = 29092
	)

	// 2. Start kafka container with listeners
	kafkaContainer, err := kafka.RunContainer(ctx,
		kafka.WithClusterID("test-cluster"),
		testcontainers.WithImage("confluentinc/confluent-local:7.5.0"),
		network.WithNetwork([]string{kafkaAliasName}, kafkaNw),
		kafka.WithTopic(kafkaTopic),
		kafka.WithListeners(kafka.NewListener("LISTENER_ALICE", kafkaAliasName, kafkaListenerPort)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// 3. Start KCat container in the same network as Kafka
	kcat, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "confluentinc/cp-kcat:7.4.1",
			Networks: []string{
				kafkaNw.Name,
			},
			Entrypoint: []string{
				"sh",
			},
			Cmd: []string{
				"-c",
				"tail -f /dev/null",
			},
			Files: []testcontainers.ContainerFile{
				{
					Reader:            bytes.NewReader([]byte(kcatMessage)),
					ContainerFilePath: kcatMsgsFile,
					FileMode:          0o777,
				},
			},
		},
		Started: true,
	})
	// }
	if err != nil {
		log.Fatalf("failed to start kcat container: %s", err)
	}

	defer func() {
		if err := kcat.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate kcat container: %s", err)
		}
		if err := kafkaContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate redpanda container: %s", err)
		}

		if err := kafkaNw.Remove(context.Background()); err != nil {
			log.Fatalf("failed to remove network: %s", err)
		}
	}()

	bootstrapServer := fmt.Sprintf("%s:%d", kafkaAliasName, kafkaListenerPort)

	// 4. Produce message to Kafka
	_, _, err = kcat.Exec(ctx, []string{"kcat", "-b", bootstrapServer, "-t", kafkaTopic, "-P", "-l", kcatMsgsFile})
	if err != nil {
		log.Fatalf("failed to produce message to Kafka: %s", err)
	}

	// 5. Consume message from Kafka
	_, stdout, err := kcat.Exec(ctx, []string{"kcat", "-b", bootstrapServer, "-C", "-t", kafkaTopic, "-c", "1"}, tcexec.Multiplexed())
	if err != nil {
		log.Fatalf("failed to consume message from Kafka: %s", err)
	}

	// 6. Read Message from stdout
	out, err := io.ReadAll(stdout)
	if err != nil {
		log.Fatalf("failed to read message from stdout: %s", err)
	}

	// 7. Assert message
	fmt.Println(string(out))

	// Output:
	// Message produced by kcat
}
