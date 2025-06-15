package containers

import (
	"context"
	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"net/http"
)

type MinecraftContainerService struct {
}

func (r *MinecraftContainerService) provideDockerClient() {
	helper, err := connhelper.GetConnectionHelper("ssh://drone@135.181.91.130")

	if err != nil {
		return
	}

	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: helper.Dialer,
		},
	}

	cli, err := client.NewClientWithOpts(
		client.WithHTTPClient(httpClient),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	val, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}
	println(val)
}

func (r *MinecraftContainerService) Start(ctx context.Context) {
	r.provideDockerClient()
}
