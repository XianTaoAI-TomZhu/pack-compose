# pack-compose

一个命令行工具，用于解析 docker-compose.yaml 和 .env 文件，拉取多架构 Docker 镜像，并将其打包为 tar 文件以便离线传输。

## 功能特性

- **Compose 文件解析**：自动检测并解析 docker-compose.yaml/docker-compose.yml 文件
- **环境变量支持**：加载并处理 .env 文件
- **多架构镜像**：支持拉取多种架构的镜像（linux/amd64、linux/arm64 等）
- **镜像打包**：将拉取的镜像导出为 tar 文件，可使用 `docker load` 完整恢复
- **友好的 CLI**：提供清晰的子命令（parse、pull、bundle），支持 --help
- **自定义文件路径**：使用 `-f/--file` 指定自定义的 docker-compose 文件路径
- **简化架构名称**：使用 `-i/--image-arch` 提供简化的架构名称（amd64、arm64）

## 安装

### 前置要求

- Go 1.21 或更高版本
- Docker 守护进程正在运行

### 从源码构建

```bash
git clone https://github.com/pack-compose/pack-compose.git
cd pack-compose
go mod tidy
go build -o pack-compose ./cmd/pack-compose
```

### 交叉编译

#### Linux/macOS (bash/zsh)

```bash
# Windows 64位
GOOS=windows GOARCH=amd64 go build -o pack-compose-windows-amd64.exe ./cmd/pack-compose

# Linux 64位
GOOS=linux GOARCH=amd64 go build -o pack-compose-linux-amd64 ./cmd/pack-compose

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o pack-compose-linux-arm64 ./cmd/pack-compose

# macOS Intel (amd64)
GOOS=darwin GOARCH=amd64 go build -o pack-compose-darwin-amd64 ./cmd/pack-compose

# macOS Apple Silicon (arm64)
GOOS=darwin GOARCH=arm64 go build -o pack-compose-darwin-arm64 ./cmd/pack-compose
```

#### Windows (PowerShell)

```powershell
# Windows 64位（当前平台）
go build -o pack-compose.exe ./cmd/pack-compose

# Linux 64位
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o pack-compose-linux-amd64 ./cmd/pack-compose

# Linux ARM64
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o pack-compose-linux-arm64 ./cmd/pack-compose

# macOS Intel (amd64)
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o pack-compose-darwin-amd64 ./cmd/pack-compose

# macOS Apple Silicon (arm64)
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o pack-compose-darwin-arm64 ./cmd/pack-compose
```

**在 PowerShell 中清除环境变量：**
```powershell
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
```

## 使用方法

### 解析 Compose 文件

解析 docker-compose 文件并列出所有引用的镜像：

```bash
pack-compose parse
```

使用自定义文件：

```bash
pack-compose parse -f ./path/to/docker-compose.yml
```

### 拉取镜像

拉取 docker-compose 文件中引用的所有镜像：

```bash
pack-compose pull
```

拉取指定架构的镜像：

```bash
# 使用完整平台格式
pack-compose pull --platform linux/amd64,linux/arm64

# 使用简化架构名
pack-compose pull -i amd64
pack-compose pull -i arm64
pack-compose pull -i amd64,arm64
```

使用自定义文件：

```bash
pack-compose pull -f ./custom-compose.yml -i amd64
```

### 打包所有内容

解析、拉取（可选）并将所有内容打包为 tar 文件：

```bash
pack-compose bundle -o ./output.tar
```

跳过拉取，使用本地镜像：

```bash
pack-compose bundle --skip-pull -o ./output.tar
```

打包指定架构的镜像：

```bash
# 使用完整平台格式
pack-compose bundle --platform linux/amd64,linux/arm64 -o ./output.tar

# 使用简化架构名
pack-compose bundle -i amd64 -o amd64-bundle.tar
pack-compose bundle -i arm64 -o arm64-bundle.tar
pack-compose bundle -i amd64,arm64 -o multi-arch-bundle.tar
```

创建 gzip 压缩包：

```bash
pack-compose bundle -o ./output.tar.gz
```

使用自定义文件：

```bash
pack-compose bundle -f ./my-compose.yml -i amd64 -o output.tar
```

## 项目结构

```
pack-compose/
├── cmd/
│   └── pack-compose/
│       ├── main.go          # 入口文件
│       └── commands/        # CLI 命令
│           ├── root.go       # 根命令
│           ├── parse.go      # parse 命令
│           ├── pull.go       # pull 命令
│           └── bundle.go     # bundle 命令
├── pkg/
│   ├── compose/             # Compose 文件解析
│   │   └── loader.go
│   ├── image/               # 镜像操作
│   │   └── puller.go
│   └── bundle/              # 打包操作
│       └── bundler.go
├── go.mod
├── go.sum
├── README.md
├── README.en.md
└── README.cn.md
```

## 常见问题

### 磁盘空间不足

如果遇到 `no space left on device` 错误，请清理 Docker 资源：

```bash
# 清理未使用的镜像、容器、网络等
docker system prune -a

# 仅清理未使用的镜像
docker image prune -a
```

### PowerShell 环境变量错误

如果看到 `GOOS=windows : 无法将"GOOS=windows"项识别` 错误，请使用 PowerShell 语法：

```powershell
# 错误（bash 语法）
GOOS=windows GOARCH=amd64 go build ...

# 正确（PowerShell 语法）
$env:GOOS="windows"; $env:GOARCH="amd64"; go build ...
```

## 许可证

MIT License
