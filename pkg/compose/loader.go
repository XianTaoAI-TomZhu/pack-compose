package compose

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/types"
)

type Loader struct {
	workDir     string
	composeFile string
}

func NewLoader(workDir string) *Loader {
	return &Loader{
		workDir: workDir,
	}
}

func NewLoaderWithFile(workDir string, composeFile string) *Loader {
	return &Loader{
		workDir:     workDir,
		composeFile: composeFile,
	}
}

func (l *Loader) findComposeFile() (string, error) {
	candidates := []string{
		"docker-compose.yml",
		"docker-compose.yaml",
		"compose.yml",
		"compose.yaml",
	}

	for _, candidate := range candidates {
		path := filepath.Join(l.workDir, candidate)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no docker-compose file found in %s", l.workDir)
}

func (l *Loader) findEnvFile() (string, error) {
	envPath := filepath.Join(l.workDir, ".env")
	if _, err := os.Stat(envPath); err == nil {
		return envPath, nil
	}
	return "", nil
}

func (l *Loader) Load() (*types.Project, error) {
	var composeFile string
	var err error

	if l.composeFile != "" {
		composeFile = l.composeFile
		if !filepath.IsAbs(composeFile) {
			composeFile = filepath.Join(l.workDir, composeFile)
		}
		if _, err := os.Stat(composeFile); err != nil {
			return nil, fmt.Errorf("compose file not found: %s", composeFile)
		}
	} else {
		composeFile, err = l.findComposeFile()
		if err != nil {
			return nil, err
		}
	}

	envFile, err := l.findEnvFile()

	var envFiles []string
	if envFile != "" {
		envFiles = []string{envFile}
	}

	options, err := cli.NewProjectOptions(
		[]string{composeFile},
		cli.WithWorkingDirectory(l.workDir),
		cli.WithEnvFiles(envFiles...),
	)
	if err != nil {
		return nil, err
	}

	return cli.ProjectFromOptions(options)
}

func (l *Loader) GetImages(project *types.Project) []string {
	var images []string
	seen := make(map[string]bool)

	for _, service := range project.Services {
		if service.Image != "" {
			if !seen[service.Image] {
				images = append(images, service.Image)
				seen[service.Image] = true
			}
		}
	}

	return images
}
