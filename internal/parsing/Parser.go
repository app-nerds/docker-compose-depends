package parsing

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

func ParseDockerCompose(r io.Reader) (DockerCompose, error) {
	var (
		err    error
		result = DockerCompose{}
	)

	decoder := yaml.NewDecoder(r)

	if err = decoder.Decode(&result); err != nil {
		return result, fmt.Errorf("error decoding docker-compose file: %w", err)
	}

	return result, nil
}
