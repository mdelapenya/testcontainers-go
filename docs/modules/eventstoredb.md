# EventStoreDB

Not available until the next release <a href="https://github.com/testcontainers/testcontainers-go"><span class="tc-version">:material-tag: main</span></a>

## Introduction

The Testcontainers module for [EventStoreDB](https://www.eventstore.com/), an open-source database built for event sourcing. It stores data as streams of immutable events, making it well-suited for event-driven architectures, CQRS, and audit logging.

## Adding this module to your project dependencies

Please run the following command to add the EventStoreDB module to your Go dependencies:

```
go get github.com/testcontainers/testcontainers-go/modules/eventstoredb
```

## Usage example

<!--codeinclude-->
[Creating a EventStoreDB container](../../modules/eventstoredb/examples_test.go) inside_block:runEventStoreDBContainer
<!--/codeinclude-->

## Module Reference

### Run function

- Not available until the next release <a href="https://github.com/testcontainers/testcontainers-go"><span class="tc-version">:material-tag: main</span></a>

The EventStoreDB module exposes one entrypoint function to create the EventStoreDB container, and this function receives three parameters:

```golang
func Run(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*Container, error)
```

- `context.Context`, the Go context.
- `string`, the Docker image to use.
- `testcontainers.ContainerCustomizer`, a variadic argument for passing options.

#### Image

Use the second argument in the `Run` function to set a valid Docker image.
In example: `Run(context.Background(), "eventstore/eventstore:24.2")`.

### Container Options

When starting the EventStoreDB container, you can pass options in a variadic way to configure it.

{% include "../features/common_functional_options_list.md" %}

#### WithInsecure

- Not available until the next release <a href="https://github.com/testcontainers/testcontainers-go"><span class="tc-version">:material-tag: main</span></a>

Sets `EVENTSTORE_INSECURE=true` on the container. EventStoreDB defaults to insecure mode in this module; this option is provided for explicitness when you want to make the intent clear in your test code.

```golang
ctr, err := eventstoredb.Run(ctx, "eventstore/eventstore:24.2",
    eventstoredb.WithInsecure(),
)
```

### Container Methods

The EventStoreDB container exposes the following methods:

#### ConnectionString

- Not available until the next release <a href="https://github.com/testcontainers/testcontainers-go"><span class="tc-version">:material-tag: main</span></a>

The `ConnectionString(ctx)` method returns the gRPC connection string for the EventStoreDB container on port `2113`.

- Insecure mode: `esdb://host:port?tls=false`
- Secure mode: `esdb+discover://host:port`

<!--codeinclude-->
[Get connection string](../../modules/eventstoredb/examples_test.go) inside_block:connectionString
<!--/codeinclude-->
