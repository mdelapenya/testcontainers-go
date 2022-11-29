package testcontainers

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// ContainerConfigOption is a function that can be used to modify the container config before it is created
type ContainerConfigOption func(*container.HostConfig, *network.NetworkingConfig)

// WithDefaultOptions sets the default options for the container, at the host config level,
// extracting them from the container request.
func WithDefaultOptions(req ContainerRequest) []ContainerConfigOption {
	return []ContainerConfigOption{
		WithLegacyAutoRemove(false),
		WithBinds(req.Binds),
		WithCapAdd(req.CapAdd),
		WithCapDrop(req.CapDrop),
		WithExtraHosts(req.ExtraHosts),
		WithNetworkMode(req.NetworkMode),
		WithResources(req.Resources),
	}
}

// WithAutoRemove sets the auto remove option for the container, at the host config level
func WithAutoRemove() ContainerConfigOption {
	return func(hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig) {
		hostConfig.AutoRemove = true
	}
}

// WithLegacyAutoRemove sets the auto remove option for the container, at the host config level
// Created for the sole purpose of supporting legacy code
func WithLegacyAutoRemove(autoRemove bool) ContainerConfigOption {
	if autoRemove {
		return WithAutoRemove()
	}

	return func(hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig) {
		hostConfig.AutoRemove = false
	}
}

// WithBinds binds a volume to the container, at the host config level
func WithBinds(binds []string) ContainerConfigOption {
	return func(hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig) {
		hostConfig.Binds = binds
	}
}

// WithCapAdd adds a capability to the container, at the host config level
func WithCapAdd(capAdds []string) ContainerConfigOption {
	return func(hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig) {
		hostConfig.CapAdd = capAdds
	}
}

// WithCapDrop drops a capability from the container, at the host config level
func WithCapDrop(capDrops []string) ContainerConfigOption {
	return func(hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig) {
		hostConfig.CapDrop = capDrops
	}
}

// WithExtraHosts adds extra hosts to the container, at the host config level
func WithExtraHosts(extraHosts []string) ContainerConfigOption {
	return func(hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig) {
		hostConfig.ExtraHosts = extraHosts
	}
}

// WithNetworkMode sets the network mode for the container, at the host config level
func WithNetworkMode(networkMode container.NetworkMode) ContainerConfigOption {
	return func(hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig) {
		hostConfig.NetworkMode = networkMode
	}
}

// WithResources sets the resources for the container, at the host config level
func WithResources(resources container.Resources) ContainerConfigOption {
	return func(hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig) {
		hostConfig.Resources = resources
	}
}
