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

### Container Methods

The Dapr container exposes the following methods:
