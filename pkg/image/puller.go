package image

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

type Puller struct {
	cli *client.Client
}

func NewPuller() (*Puller, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Puller{cli: cli}, nil
}

func (p *Puller) Close() error {
	return p.cli.Close()
}

func (p *Puller) PullImage(ctx context.Context, imageRef string, platform string) error {
	options := types.ImagePullOptions{}

	if platform != "" {
		options.Platform = platform
	}

	authConfig := registry.AuthConfig{}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	options.RegistryAuth = base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := p.cli.ImagePull(ctx, imageRef, options)
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", imageRef, err)
	}
	defer out.Close()

	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return fmt.Errorf("failed to read pull output: %w", err)
	}

	return nil
}

func (p *Puller) PullImages(ctx context.Context, images []string, platforms []string) error {
	if len(platforms) == 0 {
		platforms = []string{""}
	}

	for _, imageRef := range images {
		for _, platform := range platforms {
			fmt.Printf("Pulling %s", imageRef)
			if platform != "" {
				fmt.Printf(" for %s", platform)
			}
			fmt.Println("...")

			ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute)
			err := p.PullImage(ctxWithTimeout, imageRef, platform)
			cancel()

			if err != nil {
				return err
			}
		}
	}

	return nil
}
