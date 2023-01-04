package dependency

import (
	"fmt"
	"strings"

	"github.com/app-nerds/docker-compose-depends/internal/parsing"
)

type DependencyTree struct {
	Root *DependencyNode
}

type DependencyNode struct {
	Service     parsing.Service
	ServiceName string
	Children    []*DependencyNode
}

func BuildTree(dockerfile parsing.DockerCompose, service string) (DependencyTree, error) {
	var (
		ok            bool
		dockerService parsing.Service
		rootNode      *DependencyNode
		result        = DependencyTree{}
	)

	if dockerService, ok = dockerfile.Services[service]; !ok {
		return result, fmt.Errorf("service '%s' not found", service)
	}

	rootNode = &DependencyNode{
		Service:     dockerService,
		ServiceName: service,
		Children:    []*DependencyNode{},
	}

	childNodes := buildNode(dockerfile, rootNode)
	rootNode.Children = childNodes

	result.Root = rootNode
	return result, nil
}

func buildNode(dockerfile parsing.DockerCompose, node *DependencyNode) []*DependencyNode {
	var (
		newNode  *DependencyNode
		children = []*DependencyNode{}
	)

	for serviceName, service := range dockerfile.Services {
		if serviceName == node.ServiceName {
			for _, dependsOnService := range service.DependsOn {
				childServiceName, childService := findServiceByName(dockerfile, dependsOnService)

				newNode = &DependencyNode{
					Service:     childService,
					ServiceName: childServiceName,
					Children:    []*DependencyNode{},
				}

				childNodes := buildNode(dockerfile, newNode)
				newNode.Children = childNodes

				children = append(children, newNode)
			}
		}
	}

	return children
}

func findServiceByName(dockerfile parsing.DockerCompose, serviceToFind string) (string, parsing.Service) {
	var (
		serviceName string
		service     parsing.Service
	)

	for serviceName, service = range dockerfile.Services {
		if serviceName == serviceToFind {
			return serviceName, service
		}
	}

	return serviceName, service
}

func (dt DependencyTree) String() string {
	lines := []string{}
	out := "digraph G {\n"

	lines = dt.draw(dt.Root, lines)
	lines = dt.dedupe(lines)

	out += strings.Join(lines, "\n")
	out += "\n}"

	return out
}

func (dt DependencyTree) draw(node *DependencyNode, lines []string) []string {
	for _, n := range node.Children {
		lines = append(lines, fmt.Sprintf("\t\"%s\" -> \"%s\"", node.ServiceName, n.ServiceName))

		if len(n.Children) > 0 {
			lines = dt.draw(n, lines)
		}
	}

	return lines
}

func (dt DependencyTree) dedupe(lines []string) []string {
	result := []string{}
	m := map[string]int{}

	for _, l := range lines {
		if _, ok := m[l]; ok {
			continue
		}

		result = append(result, l)
		m[l] = 1
	}

	return result
}
