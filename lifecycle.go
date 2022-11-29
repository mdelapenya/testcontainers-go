package testcontainers

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

func (p *DockerProvider) preCreateContainerHook(ctx context.Context, req ContainerRequest, dockerInput *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, opts ...ContainerConfigOption) error {
	// prepare mounts
	hostConfig.Mounts = mapToDockerMounts(req.Mounts)

	endpointSettings := map[string]*network.EndpointSettings{}

	// #248: Docker allows only one network to be specified during container creation
	// If there is more than one network specified in the request container should be attached to them
	// once it is created. We will take a first network if any specified in the request and use it to create container
	if len(req.Networks) > 0 {
		attachContainerTo := req.Networks[0]

		nw, err := p.GetNetwork(ctx, NetworkRequest{
			Name: attachContainerTo,
		})
		if err == nil {
			endpointSetting := network.EndpointSettings{
				Aliases:   req.NetworkAliases[attachContainerTo],
				NetworkID: nw.ID,
			}
			endpointSettings[attachContainerTo] = &endpointSetting
		}
	}

	if len(opts) == 0 {
		opts = WithDefaultOptions(req)
	}
	for _, opt := range opts {
		opt(hostConfig, networkingConfig)
	}

	networkingConfig.EndpointsConfig = endpointSettings

	exposedPorts := req.ExposedPorts
	// this check must be done after the PreCreateModifier is called, so the network mode is already set
	if len(exposedPorts) == 0 && !hostConfig.NetworkMode.IsContainer() {
		image, _, err := p.client.ImageInspectWithRaw(ctx, dockerInput.Image)
		if err != nil {
			return err
		}
		for p := range image.ContainerConfig.ExposedPorts {
			exposedPorts = append(exposedPorts, string(p))
		}
	}

	exposedPortSet, exposedPortMap, err := nat.ParsePortSpecs(exposedPorts)
	if err != nil {
		return err
	}

	dockerInput.ExposedPorts = exposedPortSet
	hostConfig.PortBindings = exposedPortMap

	return nil
}
