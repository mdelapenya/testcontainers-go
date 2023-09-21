# Dapr

Not available until the next release of testcontainers-go <a href="https://github.com/testcontainers/testcontainers-go"><span class="tc-version">:material-tag: main</span></a>

## Introduction

The Testcontainers module for Dapr.

## Adding this module to your project dependencies

Please run the following command to add the Dapr module to your Go dependencies:

```
go get github.com/testcontainers/testcontainers-go/modules/dapr
```

## Usage example

<!--codeinclude-->
[Creating a Dapr container](../../modules/dapr/examples_test.go) inside_block:runDaprContainer
<!--/codeinclude-->

## Module reference

The Dapr module exposes one entrypoint function to create the Dapr container, and this function receives two parameters:

```golang
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*DaprContainer, error)
```

- `context.Context`, the Go context.
- `testcontainers.ContainerCustomizer`, a variadic argument for passing options.

### Container Options

When starting the Dapr container, you can pass options in a variadic way to configure it.

#### Image

If you need to set a different Dapr Docker image, you can use `testcontainers.WithImage` with a valid Docker image
for Dapr. E.g. `testcontainers.WithImage("daprio/daprd:1.11.3")`.

{% include "../features/common_functional_options.md" %}

#### Application Name

It's possible to define the application name used by Dapr with the `WithAppName(name string)` functional option. If not passed, the default value is `dapr-app`.

#### Components

You can add components to the Dapr container with the `WithComponents(components ...Component)` functional option. If not passed, the default value is an empty map.

The `Component` struct has the following fields:

<!--codeinclude-->
[Dapr Component](../../modules/dapr/options.go) inside_block:componentStruct
<!--/codeinclude-->

- The key used to internally identify a Component is a string formed by the component name and the component type, separated by a colon. E.g. `my-pubsub:pubsub`.
- Metadata it's a map of strings, where the key is the metadata name and the value is the metadata value. It will be used to render a YAML file with the component configuration.

Each component will result in a configuration file that will be uploaded to the Dapr container, under the `/components` directory. It's possible to change this file path with the `WithComponentsPath(path string)` functional option. If not passed, the default value is `/components`.

The file will be named as the component name, and the content will be a YAML file with the following structure:

```yaml
apiVersion: dapr.io/v1alpha1
  kind: Component
  metadata:
    name: statestore
  spec:
    type: state.in-memory
    version: v1
  metadata:
    - name: foo1
      value: bar1
    - name: foo2
      value: bar2
```

### Container Methods

The Dapr container exposes the following methods:

#### GRPCPort

This method returns the integer representation of the exposed port for the Dapr gRPC API, which internally is `50001`, and an error if something went wrong while retrieving the port.
