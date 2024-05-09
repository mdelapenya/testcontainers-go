package kafka

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
)

type options struct {
	// AutoCreateTopics will automatically create topics when the container starts.
	AutoCreateTopics bool

	// TopicCreationHooks is a list of hooks that will be executed after the container is ready.
	TopicCreationHooks []testcontainers.ContainerHook

	// KafkaAuthenticationMethod is either "none" for plaintext or "sasl"
	// for SASL (scram) authentication.
	KafkaAuthenticationMethod string

	// Listeners is a list of custom listeners that can be provided to access the
	// containers form within docker networks
	Listeners []Listener
}

func defaultOptions() options {
	return options{
		KafkaAuthenticationMethod: "none",
		Listeners:                 make([]Listener, 0),
	}
}

// Compiler check to ensure that Option implements the testcontainers.ContainerCustomizer interface.
var _ testcontainers.ContainerCustomizer = (Option)(nil)

// Option is an option for the Kafka container.
type Option func(*options)

// Customize is a NOOP. It's defined to satisfy the testcontainers.ContainerCustomizer interface.
func (o Option) Customize(*testcontainers.GenericContainerRequest) error {
	// NOOP to satisfy interface.
	return nil
}

// WithListeners adds a custom listener to the Kafka container. Listener
// will be aliases to all networks, so they can be accessed from within docker
// networks. At least one network must be attached to the container, if not an
// error will be thrown when starting the container.
func WithListeners(l ...Listener) Option {
	return func(o *options) {
		o.Listeners = append(o.Listeners, l...)
	}
}

// WithTopic adds a topic to the list of topics that will be created when the
// container starts. It won't check the format of the topic, so it's up to the
// user to provide a valid topic name and configuration. E.g. "topic:1:1"
func WithTopic(topic string) Option {
	cmds := []string{
		"kafka-topics", "--create", "--topic", topic, "--bootstrap-server", "localhost:9092",
	}

	return func(o *options) {
		o.TopicCreationHooks = append(o.TopicCreationHooks, func(ctx context.Context, container testcontainers.Container) error {
			code, _, err := container.Exec(ctx, cmds, tcexec.Multiplexed())
			if err != nil {
				return err
			}

			if code != 0 {
				return fmt.Errorf("failed to create topic %s", topic)
			}

			return nil
		})
	}
}
