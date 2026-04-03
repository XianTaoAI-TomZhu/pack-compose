package commands

import (
	"fmt"
	"os"

	"github.com/pack-compose/pack-compose/pkg/compose"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse Docker Compose file and list images",
	Long:  `Parse the docker-compose file in the current directory and list all referenced images.`,
	RunE:  runParse,
}

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringP("file", "f", "", "Path to docker-compose file (optional)")
}

func runParse(cmd *cobra.Command, args []string) error {
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

	fmt.Printf("Project: %s\n", project.Name)
	fmt.Printf("Services: %d\n", len(project.Services))
	fmt.Printf("\nImages:\n")
	for i, img := range images {
		fmt.Printf("  %d. %s\n", i+1, img)
	}

	return nil
}
