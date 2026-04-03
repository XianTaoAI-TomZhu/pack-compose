package bundle

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/client"
)

type Bundler struct {
	cli *client.Client
}

func NewBundler() (*Bundler, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Bundler{cli: cli}, nil
}

func (b *Bundler) Close() error {
	return b.cli.Close()
}

func (b *Bundler) SaveImages(ctx context.Context, images []string, outputPath string) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	rc, err := b.cli.ImageSave(ctxWithTimeout, images)
	if err != nil {
		return fmt.Errorf("failed to save images: %w", err)
	}
	defer rc.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func (b *Bundler) Bundle(ctx context.Context, images []string, outputPath string, includeComposeFiles bool, workDir string) error {
	if filepath.Ext(outputPath) == ".gz" || filepath.Ext(outputPath) == ".tgz" {
		return b.bundleGzipped(ctx, images, outputPath, includeComposeFiles, workDir)
	}
	return b.bundlePlain(ctx, images, outputPath, includeComposeFiles, workDir)
}

func (b *Bundler) bundlePlain(ctx context.Context, images []string, outputPath string, includeComposeFiles bool, workDir string) error {
	return b.SaveImages(ctx, images, outputPath)
}

func (b *Bundler) bundleGzipped(ctx context.Context, images []string, outputPath string, includeComposeFiles bool, workDir string) error {
	tempTar, err := os.CreateTemp("", "pack-compose-*.tar")
	if err != nil {
		return err
	}
	tempTarPath := tempTar.Name()
	tempTar.Close()
	defer os.Remove(tempTarPath)

	err = b.SaveImages(ctx, images, tempTarPath)
	if err != nil {
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gw := gzip.NewWriter(outFile)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	tarFile, err := os.Open(tempTarPath)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	info, err := tarFile.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    "images.tar",
		Size:    info.Size(),
		Mode:    0644,
		ModTime: time.Now(),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := io.Copy(tw, tarFile); err != nil {
		return err
	}

	if includeComposeFiles {
		files := []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml", ".env"}
		for _, f := range files {
			path := filepath.Join(workDir, f)
			if _, err := os.Stat(path); err == nil {
				err = addFileToTar(tw, path, f)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func addFileToTar(tw *tar.Writer, filePath, arcName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = arcName

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, file); err != nil {
		return err
	}

	return nil
}
