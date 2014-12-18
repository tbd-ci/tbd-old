package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

func main() {
	var endpoint, appname, directory string
	flag.StringVar(&endpoint, "endpoint", "unix:///var/run/docker.sock", "docker daemon endpoint")
	flag.StringVar(&appname, "appname", "tbd_app_tmp", "container name for build")
	flag.StringVar(&directory, "directory", ".", "context directory for docker build")

	flag.Parse()

	client, err := docker.NewClient(endpoint)
	if err != nil {
		fmt.Println(err)
	}

	var builder docker.BuildImageOptions

	builder.Name = appname
	builder.RmTmpContainer = true
	builder.ContextDir = directory
	builder.OutputStream = os.Stderr

	err = client.BuildImage(builder)
	if err != nil {
		fmt.Println(err)
	}
}
