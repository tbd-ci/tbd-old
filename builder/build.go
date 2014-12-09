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

	//dockerfile, err := os.Open("./Dockerfile")
	//if err != nil {
	//  return fmt.Errorf("os.Open: %s", err)
	//}
	//defer dockerfile.Close()

	var builder docker.BuildImageOptions

	builder.Name = "testapp"
	builder.RmTmpContainer = true
	builder.OutputStream = os.Stderr
	builder.Remote = "foo/bar/baz"

	err = client.BuildImage(builder)
	if err != nil {
		return fmt.Errorf("client.BuildImage: %s", err)
	}

	return err
}
