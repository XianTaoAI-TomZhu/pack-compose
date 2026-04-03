# pack-compose

A command-line tool to parse docker-compose.yaml and .env files, pull multi-architecture Docker images, and package them into tar files for offline transfer.

## Features

- **Compose File Parsing**: Automatically detects and parses docker-compose.yaml/docker-compose.yml files
- **Environment Variable Support**: Loads and processes .env files
- **Multi-Architecture Images**: Supports pulling images for multiple architectures (linux/amd64, linux/arm64, etc.)
- **Image Bundling**: Exports pulled images to tar files that can be loaded with `docker load`
- **CLI Friendly**: Provides clear subcommands (parse, pull, bundle) with --help support

## Installation

### Prerequisites

- Go 1.21 or later
- Docker daemon running

### Build from Source

```bash
git clone https://github.com/pack-compose/pack-compose.git
cd pack-compose
go mod tidy
go build -o pack-compose ./cmd/pack-compose
```

## Usage

### Parse Compose File

Parse the docker-compose file and list all referenced images:

```bash
pack-compose parse
```

Use custom file:

```bash
pack-compose parse -f ./path/to/docker-compose.yml
```

### Pull Images

Pull all images referenced in the docker-compose file:

```bash
pack-compose pull
```

Pull images for specific architectures:

```bash
# Using full platform format
pack-compose pull --platform linux/amd64,linux/arm64

# Using simplified architecture name
pack-compose pull -i amd64
pack-compose pull -i arm64
pack-compose pull -i amd64,arm64
```

Use custom file:

```bash
pack-compose pull -f ./custom-compose.yml -i amd64
```

### Bundle Everything

Parse, pull (optional), and bundle everything into a tar file:

```bash
pack-compose bundle -o ./output.tar
```

Skip pulling and use local images:

```bash
pack-compose bundle --skip-pull -o ./output.tar
```

Bundle with specific architectures:

```bash
# Using full platform format
pack-compose bundle --platform linux/amd64,linux/arm64 -o ./output.tar

# Using simplified architecture name
pack-compose bundle -i amd64 -o amd64-bundle.tar
pack-compose bundle -i arm64 -o arm64-bundle.tar
pack-compose bundle -i amd64,arm64 -o multi-arch-bundle.tar
```

Create a gzipped bundle:

```bash
pack-compose bundle -o ./output.tar.gz
```

Use custom file:

```bash
pack-compose bundle -f ./my-compose.yml -i amd64 -o output.tar
```

## Project Structure

```
pack-compose/
├── cmd/
│   └── pack-compose/
│       ├── main.go          # Entry point
│       └── commands/        # CLI commands
│           ├── root.go       # Root command
│           ├── parse.go      # Parse command
│           ├── pull.go       # Pull command
│           └── bundle.go     # Bundle command
├── pkg/
│   ├── compose/             # Compose file parsing
│   │   └── loader.go
│   ├── image/               # Image operations
│   │   └── puller.go
│   └── bundle/              # Bundling operations
│       └── bundler.go
├── go.mod
├── go.sum
└── README.md
```

## FAQ

### No Space Left on Device

If you encounter `no space left on device` error, clean up Docker resources:

```bash
# Clean up unused images, containers, networks, etc.
docker system prune -a

# Only clean up unused images
docker image prune -a
```

## License

MIT License
