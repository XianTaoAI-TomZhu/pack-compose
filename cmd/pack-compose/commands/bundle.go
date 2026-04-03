package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pack-compose/pack-compose/pkg/bundle"
	"github.com/pack-compose/pack-compose/pkg/compose"
	"github.com/pack-compose/pack-compose/pkg/image"
	"github.com/spf13/cobra"
)

var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Bundle all images and files into a tar archive",
	Long:  `Parse the docker-compose file, pull images, and bundle everything into a tar archive.`,
	RunE:  runBundle,
}

var (
	outputFlag         string
	includeComposeFlag bool
	skipPullFlag       bool
)

func init() {
	rootCmd.AddCommand(bundleCmd)
	bundleCmd.Flags().StringP("file", "f", "", "Path to docker-compose file (optional)")
	bundleCmd.Flags().StringVarP(&outputFlag, "output", "o", "pack-compose-bundle.tar", "Output file path")
	bundleCmd.Flags().BoolVar(&includeComposeFlag, "include-compose", true, "Include compose and .env files in the bundle")
	bundleCmd.Flags().BoolVar(&skipPullFlag, "skip-pull", false, "Skip pulling images (use local images)")
	bundleCmd.Flags().StringSliceVarP(&platformsFlag, "platform", "p", []string{}, "Target platforms (e.g., linux/amd64,linux/arm64)")
	bundleCmd.Flags().StringSliceP("image-arch", "i", []string{}, "Target image architectures (e.g., amd64,arm64)")
}

func runBundle(cmd *cobra.Command, args []string) error {
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

	if !skipPullFlag {
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

		if err := puller.PullImages(ctx, images, platforms); err != nil {
			return err
		}
	}

	bundler, err := bundle.NewBundler()
	if err != nil {
		return err
	}
	defer bundler.Close()

	if outputFlag == "" {
		outputFlag = "pack-compose-bundle.tar"
	}

	absOutput, err := filepath.Abs(outputFlag)
	if err != nil {
		return err
	}

	fmt.Printf("Bundling %d images to %s...\n", len(images), absOutput)
	if includeComposeFlag {
		fmt.Println("Including compose and .env files")
	}

	err = bundler.Bundle(ctx, images, absOutput, includeComposeFlag, workDir)
	if err != nil {
		return err
	}

	fmt.Printf("Bundle created successfully: %s\n", absOutput)
	return nil
}
