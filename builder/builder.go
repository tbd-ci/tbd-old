package builder

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

func Build() error {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return fmt.Errorf("docker.Newclient: %s", err)
		return err
	}

	var builder docker.BuildImageOptions

	builder.Name = "testapp"
	builder.RmTmpContainer = true
	builder.ContextDir = "."
	builder.OutputStream = os.Stderr

	err = client.BuildImage(builder)
	if err != nil {
		return fmt.Errorf("client.BuildImage: %s", err)
	}

	return err
}
