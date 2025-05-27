# ImgOod - Image Optimization and S3 Management Tool

A powerful command-line tool for image optimization, format conversion, and S3 object management, built with [bimg](https://github.com/h2non/bimg) (powered by libvips) and AWS SDK for Go v2.

## Features

- Image compression and resizing
- Format conversion (WebP, JPEG, PNG)
- S3 upload with customizable paths
- S3 object copying with format conversion
- Configuration via TOML files and environment variables
- Support for custom S3-compatible storage services

## Prerequisites

- Go 1.24 or later
- libvips 8.10+ (required by bimg)

### Installing libvips

#### macOS

```bash
brew install vips
```

#### Ubuntu/Debian

```bash
apt-get install libvips-dev
```

#### CentOS/RHEL

```bash
yum install vips-devel
```

## Installation

```bash
go install github.com/mingeme/imgood@latest
```

Or build from source:

```bash
git clone https://github.com/mingeme/imgood.git
cd imgood
go build -o imgood ./cmd/imgood
```

## Commands

ImgOod provides the following commands:

- `up`: Upload images to S3 with optional compression and format conversion
- `cp`: Copy objects within S3 with optional format conversion and resizing

## Configuration

ImgOod supports configuration through:

1. Configuration files (`config.toml`)
2. Environment variables

### Configuration File Locations

The tool looks for a `config.toml` file in the following locations (in order):

1. Current directory (`./config.toml`)
2. User's home directory (`~/.imgood/config.toml`)
3. XDG config directory (`~/.config/imgood/config.toml`)

### Example Configuration

```toml
# ImgOod Configuration

# S3 Configuration
[s3]
# S3 bucket name (required)
bucket = "my-images-bucket"

# S3 endpoint URL (for non-AWS S3 services)
# Format: https://s3.example.com
# Leave empty for AWS S3
endpoint = "https://s3.bitiful.net"

# AWS region (required for AWS S3)
region = "us-east-1"

# AWS credentials
# These can be left empty if using environment variables or AWS credential files
access_key = "your-access-key"
secret_key = "your-secret-key"
```

### Setting Environment Variables

All configuration options can also be set using environment variables with the prefix `IMGOOD_`:

```bash
# S3 configuration
export IMGOOD_S3_BUCKET="my-images-bucket"
export IMGOOD_S3_ENDPOINT="https://s3.bitiful.net"
export IMGOOD_S3_REGION="us-east-1"
export IMGOOD_S3_ACCESS_KEY="your-access-key"
export IMGOOD_S3_SECRET_KEY="your-secret-key"
```

## Command Usage

### Upload Command (`up`)

Upload images to S3 with optional compression and format conversion.

```bash
imgood up [options]
```

#### Upload Options

- `-i, --input string`: Path to the input image file (required)
- `-k, --key string`: S3 object key (path in bucket), defaults to filename
- `-c, --compress`: Compress image before uploading
- `-q, --quality int`: Quality of the compressed image (1-100) (default 80)
- `-w, --width int`: Width of the output image (0 for original)
- `-h, --height int`: Height of the output image (0 for original)

#### Examples

Upload an image with default settings:

```bash
imgood up -i sample.jpg
```

Upload with compression and custom key:

```bash
imgood up -i sample.jpg -k images/2025/05/sample.jpg -c -q 85
```

Upload with resizing (converts to WebP format by default):

```bash
imgood up -i sample.jpg -c -w 800 -h 600
```

### Copy Command (`cp`)

Copy objects within S3 with optional format conversion and resizing.

```bash
imgood cp [options]
```

#### Copy Options

- `-s, --source string`: Source S3 object key to copy (required)
- `-t, --target string`: Target S3 object key (destination), defaults to source-copy
- `-f, --format string`: Convert to format (webp, jpeg, png) (default "webp")
- `-q, --quality int`: Quality of the converted image (1-100) (default 80)
- `-w, --width int`: Width of the output image (0 for original)
- `-h, --height int`: Height of the output image (0 for original)

#### Copy Command Examples

Copy an object with default settings (converts to WebP):

```bash
imgood cp -s images/original.jpg
```

Copy with specific format and target key:

```bash
imgood cp -s images/original.jpg -t images/copy.png -f png
```

Copy with resizing and quality adjustment:

```bash
imgood cp -s images/original.jpg -w 1200 -h 800 -q 90
```

## URL Format

When using custom S3 endpoints, Imgood generates URLs in the format:

```text
https://{bucket}.{endpoint}/{key}
```

For example, with bucket `imgood` and endpoint `s3.example.com`, the URL would be:

```text
https://imgood.s3.example.com/path/to/image.jpg
```

For AWS S3, the standard format is used:

```text
https://{bucket}.s3.{region}.amazonaws.com/{key}
```

## License

MIT
