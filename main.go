package main

import (
	"fmt"
	"os"

	"github.com/app-nerds/configinator"
	"github.com/app-nerds/docker-compose-depends/internal/dependency"
	"github.com/app-nerds/docker-compose-depends/internal/parsing"
	"github.com/go-playground/validator/v10"
)

type configuration struct {
	InFile      string `flag:"in" description:"File name and path to docker compose file" validate:"required"`
	ServiceName string `flag:"service" description:"Name of service to model" validate:"required_if=List false"`
	List        bool   `flag:"list" description:"Lists all services in docker compose file"`
}

func main() {
	var (
		err               error
		f                 *os.File
		dockerComposeFile parsing.DockerCompose
		tree              dependency.DependencyTree
	)

	config := configuration{}
	configinator.Behold(&config)
	validate := validator.New()

	if err = validate.Struct(config); err != nil {
		fmt.Printf("error validating configuration: %s\n", err.Error())
		os.Exit(-1)
	}

	if f, err = os.Open(config.InFile); err != nil {
		fmt.Printf("error opening input file '%s': %s\n", config.InFile, err.Error())
		os.Exit(-1)
	}

	if dockerComposeFile, err = parsing.ParseDockerCompose(f); err != nil {
		fmt.Printf("error parsing docker-compose file: %s\n", err.Error())
		os.Exit(-1)
	}

	if config.List {
		printServices(dockerComposeFile)
		os.Exit(0)
	}

	if tree, err = dependency.BuildTree(dockerComposeFile, config.ServiceName); err != nil {
		fmt.Printf("error generating DOT tree: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Printf("%s\n", tree.String())
}

func printServices(dockerfile parsing.DockerCompose) {
	for serviceName := range dockerfile.Services {
		fmt.Printf("%s\n", serviceName)
	}
}
