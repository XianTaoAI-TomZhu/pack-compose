package commands

import (
	"context"
	"os"
	"strings"

	"github.com/pack-compose/pack-compose/pkg/compose"
	"github.com/pack-compose/pack-compose/pkg/image"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull all images from Docker Compose file",
	Long:  `Parse the docker-compose file and pull all referenced images.`,
	RunE:  runPull,
}

var (
	platformsFlag []string
)

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().StringP("file", "f", "", "Path to docker-compose file (optional)")
	pullCmd.Flags().StringSliceVarP(&platformsFlag, "platform", "p", []string{}, "Target platforms (e.g., linux/amd64,linux/arm64)")
	pullCmd.Flags().StringSliceP("image-arch", "i", []string{}, "Target image architectures (e.g., amd64,arm64)")
}

func runPull(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	composeFilePath, _ := cmd.Flags().GetString("file")

	var loader *compose.Loader
	if composeFilePath != "" {
		loader = compose.NewLoaderWithFile(workDir, composeFilePath)
	} else {
		loader = compose.NewLoader(workDir)
	}

	project, err := loader.Load()
	if err != nil {
		return err
	}

	images := loader.GetImages(project)

	puller, err := image.NewPuller()
	if err != nil {
		return err
	}
	defer puller.Close()

	var platforms []string

	platformFromFlag, _ := cmd.Flags().GetStringSlice("platform")
	imageArchFromFlag, _ := cmd.Flags().GetStringSlice("image-arch")

	for _, p := range platformFromFlag {
		p = strings.TrimSpace(p)
		if p != "" {
			platforms = append(platforms, p)
		}
	}

	for _, arch := range imageArchFromFlag {
		arch = strings.TrimSpace(arch)
		if arch != "" {
			platforms = append(platforms, "linux/"+arch)
		}
	}

	return puller.PullImages(ctx, images, platforms)
}
