package parsing

type DockerCompose struct {
	Version  string `yaml:"version"`
	Services map[string]Service
}

type Service struct {
	DependsOn []string `yaml:"depends_on,omitempty"`
}
