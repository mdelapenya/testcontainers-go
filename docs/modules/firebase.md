# Firebase

Not available until the next release of testcontainers-go <a href="https://github.com/testcontainers/testcontainers-go"><span class="tc-version">:material-tag: main</span></a>

## Introduction

The Testcontainers module for Firebase.

## Adding this module to your project dependencies

Please run the following command to add the Firebase module to your Go dependencies:

```
go get github.com/testcontainers/testcontainers-go/modules/firebase
```

## Usage example

<!--codeinclude-->
[Creating a Firebase container](../../modules/firebase/examples_test.go) inside_block:runFirebaseContainer
<!--/codeinclude-->

## Module reference

The Firebase module exposes one entrypoint function to create the Firebase container, and this function receives two parameters:

```golang
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*FirebaseContainer, error)
```

- `context.Context`, the Go context.
- `testcontainers.ContainerCustomizer`, a variadic argument for passing options.

### Container Options

When starting the Firebase container, you can pass options in a variadic way to configure it.

#### Image

If you need to set a different Firebase Docker image, you can use `testcontainers.WithImage` with a valid Docker image
for Firebase. E.g. `testcontainers.WithImage("ghcr.io/thoughtgears/docker-firebase-emulator:13.6.0")`.

{% include "../features/common_functional_options.md" %}

### Container Methods

The Firebase container exposes the following methods:
