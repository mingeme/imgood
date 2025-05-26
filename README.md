# Imgood - Image Compression Tool

A simple command-line tool for image compression and S3 uploading using the [bimg](https://github.com/h2non/bimg) library, which is powered by libvips.

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
go mod download
go build -o imgood
```

## Usage

Imgood has two main commands:

1. Default command (compression)
2. Upload command (`up`) for S3 uploads

### Compression Command

```bash
./imgood --input <path-to-image> [--quality <1-100>] [--width <pixels>] [--height <pixels>]
```

You can also use short flag options:

```bash
./imgood -i <path-to-image> [-q <1-100>] [-w <pixels>] [-h <pixels>]
```

#### Compression Options

- `--input`, `-i`: Path to the input image file (required)
- `--quality`, `-q`: Quality of the compressed image (1-100, default: 80)
- `--width`, `-w`: Width of the output image (0 for original, default: 0)
- `--height`, `-h`: Height of the output image (0 for original, default: 0)

#### Compression Examples

Using long options:

```bash
./imgood --input sample.jpg --quality 70 --width 800
```

Using short options:

```bash
./imgood -i sample.jpg -q 70 -w 800
```

This will compress `sample.jpg` with 70% quality and resize it to 800px width (maintaining aspect ratio), then save it to the system's temporary directory.

### Upload Command (S3)

```bash
./imgood up --input <path-to-image> --bucket <bucket-name> [--key <object-key>] [--compress] [other options]
```

You can also use short flag options:

```bash
./imgood up -i <path-to-image> -b <bucket-name> [-k <object-key>] [-c] [other options]
```

#### Upload Options

- `--input`, `-i`: Path to the input image file (required)
- `--bucket`, `-b`: S3 bucket name (required)
- `--key`, `-k`: S3 object key/path (default: original filename)
- `--endpoint`, `-e`: S3 endpoint URL for non-AWS S3 services
- `--region`, `-r`: AWS region (default: us-east-1)
- `--access-key`, `-a`: AWS access key ID
- `--secret-key`, `-s`: AWS secret access key
- `--compress`, `-c`: Compress image before uploading (boolean flag)

When `--compress` is specified, you can also use the compression options:

- `--quality`, `-q`: Quality of the compressed image (1-100, default: 80)
- `--width`, `-w`: Width of the output image (0 for original, default: 0)
- `--height`, `-h`: Height of the output image (0 for original, default: 0)

#### Upload Examples

Upload an image to S3 without compression:

```bash
./imgood up -i sample.jpg -b my-images-bucket
```

Upload with compression and resizing:

```bash
./imgood up -i sample.jpg -b my-images-bucket -k images/sample-resized.jpg -c -q 75 -w 800
```

Use with MinIO or other S3-compatible storage:

```bash
./imgood up -i sample.jpg -b my-bucket -e http://minio-server:9000 -a myAccessKey -s mySecretKey
```

## Configuration

Imgood supports configuration through:

1. Command-line flags
2. Configuration file (`config.yaml`)
3. Environment variables

### Configuration File

The tool looks for a `config.toml` file in the following locations:

- Current directory
- `~/.imgood/` directory

Example `config.toml`:

```toml
# Imgood Configuration

# S3 Configuration
[s3]
# Default S3 bucket name
bucket = "my-images-bucket"

# S3 endpoint URL (for non-AWS S3 services like MinIO)
# Leave empty for AWS S3
endpoint = "https://minio.example.com"

# AWS region
region = "us-east-1"

# AWS credentials
# These can be left empty if using environment variables or AWS credential files
access_key = "your-access-key"
secret_key = "your-secret-key"
```

### Environment Variables

All configuration options can also be set using environment variables with the prefix `IMGOOD_`. For nested configuration options, use underscores to separate the levels.

Examples:

```bash
# Set S3 bucket
export IMGOOD_S3_BUCKET="my-images-bucket"

# Set AWS region
export IMGOOD_S3_REGION="eu-west-1"

# Set S3 endpoint for MinIO
export IMGOOD_S3_ENDPOINT="https://minio.example.com"

# Set AWS credentials
export IMGOOD_S3_ACCESS_KEY="your-access-key"
export IMGOOD_S3_SECRET_KEY="your-secret-key"
```

## Output

The compressed image will be saved to the system's temporary directory (`/tmp` on Unix/Linux/macOS) with the naming format `{original_filename}_compressed.{extension}`.
