package app

import (
	"errors"
	"fmt"
	"io"

	"github.com/antchfx/xmlquery"
	"github.com/maxwellhertz/maven-project-ci-tool/internal/pkg"
)

type DefaultPomUtils struct {
	xmlRoot     *xmlquery.Node
	projectNode *xmlquery.Node
	pluginsNode *xmlquery.Node

	mavenRepoApi MavenRepositoryApi
}

type MavenRepositoryApi interface {
	GetLatestReleaseVersion(artifactId string) (string, error)
}

const (
	requiredModelVersion = "4.0.0"

	surefirePluginGroupId    = "org.apache.maven.plugins"
	surefirePluginArtifactId = "maven-surefire-plugin"

	jacocoPluginGroupId    = "org.jacoco"
	jacocoPluginArtifactId = "jacoco-maven-plugin"
)

func NewDefaultPomUitls(mavenRepoApi MavenRepositoryApi) pkg.PomUtils {
	return &DefaultPomUtils{mavenRepoApi: mavenRepoApi}
}

func (u *DefaultPomUtils) ConfigurePlugins(pomFile io.Reader) ([]byte, error) {
	// Validate the POM.
	doc, err := xmlquery.Parse(pomFile)
	if err != nil {
		return nil, fmt.Errorf("can't parse POM file: %w", err)
	}
	u.xmlRoot = doc
	projectNode, err := validatePom(doc)
	if err != nil {
		return nil, fmt.Errorf("POM is not valid: %w", err)
	}
	u.projectNode = projectNode

	// Add <build> section if needed.
	buildNode := xmlquery.FindOne(u.projectNode, "build")
	if buildNode == nil {
		buildNode = &xmlquery.Node{
			Type: xmlquery.ElementNode,
			Data: "build",
		}
		xmlquery.AddChild(u.projectNode, buildNode)
	}
	pluginsNode := xmlquery.FindOne(buildNode, "plugins")
	if pluginsNode == nil {
		pluginsNode = &xmlquery.Node{
			Type: xmlquery.ElementNode,
			Data: "plugins",
		}
		xmlquery.AddChild(buildNode, pluginsNode)
	}
	u.pluginsNode = pluginsNode

	// Add Surefire plugin.
	if _, err := u.createPluginNode(surefirePluginGroupId, surefirePluginArtifactId); err != nil {
		return nil, fmt.Errorf("can't add %v: %w", surefirePluginArtifactId, err)
	}

	// Add JaCoCo plugin.
	jacocoNode, err := u.createPluginNode(jacocoPluginGroupId, jacocoPluginArtifactId)
	if err != nil {
		return nil, fmt.Errorf("can't add %v: %w", jacocoPluginArtifactId, err)
	}
	u.configureJaCoCoPlugin(jacocoNode)

	return []byte(u.xmlRoot.OutputXML(true)), nil
}

func (u *DefaultPomUtils) createPluginNode(groupId string, artifactId string) (*xmlquery.Node, error) {
	if node := xmlquery.FindOne(u.pluginsNode, fmt.Sprintf("plugin[artifactId='%v']", artifactId)); node != nil {
		xmlquery.RemoveFromTree(node)
	}

	// Get the latest release.
	version, err := u.mavenRepoApi.GetLatestReleaseVersion(artifactId)
	if err != nil {
		return nil, fmt.Errorf("can't query the latest release of %v: %w", artifactId, err)
	}

	// Add a new <plugin> section.
	versionNode := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "version",
		FirstChild: &xmlquery.Node{
			Type: xmlquery.TextNode,
			Data: version,
		},
	}
	artifactIdNode := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "artifactId",
		FirstChild: &xmlquery.Node{
			Type: xmlquery.TextNode,
			Data: artifactId,
		},
		NextSibling: versionNode,
	}
	groupIdNode := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "groupId",
		FirstChild: &xmlquery.Node{
			Type: xmlquery.TextNode,
			Data: groupId,
		},
		NextSibling: artifactIdNode,
	}
	pluginNode := &xmlquery.Node{
		Type:       xmlquery.ElementNode,
		Data:       "plugin",
		FirstChild: groupIdNode,
		LastChild:  versionNode,
	}
	xmlquery.AddChild(u.pluginsNode, pluginNode)
	return pluginNode, nil
}

func validatePom(xmlRoot *xmlquery.Node) (*xmlquery.Node, error) {
	projectNode := xmlquery.FindOne(xmlRoot, "project")
	if projectNode == nil {
		return nil, errors.New("<project> element is missing")
	}

	modelVersionNode := xmlquery.FindOne(projectNode, "modelVersion")
	if modelVersionNode == nil {
		return nil, errors.New("<modelVersion> element is missing")
	}
	if modelVersionNode.InnerText() != requiredModelVersion {
		return nil, errors.New("modelVersion must be " + requiredModelVersion)
	}
	return projectNode, nil
}

func (u *DefaultPomUtils) configureJaCoCoPlugin(jacocoNode *xmlquery.Node) {
	executionsNode := createExecutionsNode(
		execution{
			id:   "default-prepare-agent",
			goal: "prepare-agent",
		}, execution{
			id:    "default-report",
			goal:  "report",
			phase: "test",
		})
	xmlquery.AddChild(jacocoNode, executionsNode)
}

type execution struct {
	id    string
	goal  string
	phase string
}

func createExecutionsNode(executions ...execution) *xmlquery.Node {
	executionsNode := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "executions",
	}
	for _, exec := range executions {
		goalsNode := &xmlquery.Node{
			Type: xmlquery.ElementNode,
			Data: "goals",
			FirstChild: &xmlquery.Node{
				Type: xmlquery.ElementNode,
				Data: "goal",
				FirstChild: &xmlquery.Node{
					Type: xmlquery.TextNode,
					Data: exec.goal,
				},
			},
		}
		idNode := &xmlquery.Node{
			Type: xmlquery.ElementNode,
			Data: "id",
			FirstChild: &xmlquery.Node{
				Type: xmlquery.TextNode,
				Data: exec.id,
			},
			NextSibling: goalsNode,
		}
		executionNode := &xmlquery.Node{
			Type:       xmlquery.ElementNode,
			Data:       "execution",
			FirstChild: idNode,
		}
		if exec.phase != "" {
			phaseNode := &xmlquery.Node{
				Type: xmlquery.ElementNode,
				Data: "phase",
				FirstChild: &xmlquery.Node{
					Type: xmlquery.TextNode,
					Data: exec.phase,
				},
			}
			goalsNode.NextSibling = phaseNode
		}
		xmlquery.AddChild(executionsNode, executionNode)
	}
	return executionsNode
}
