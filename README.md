# Imgood - Image Compression Tool

A simple command-line tool for image compression using the [bimg](https://github.com/h2non/bimg) library, which is powered by libvips.

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

```bash
./imgood --input <path-to-image> [--quality <1-100>] [--width <pixels>] [--height <pixels>]
```

You can also use short flag options:

```bash
./imgood -i <path-to-image> [-q <1-100>] [-w <pixels>] [-h <pixels>]
```

### Options

- `--input`, `-i`: Path to the input image file (required)
- `--quality`, `-q`: Quality of the compressed image (1-100, default: 80)
- `--width`, `-w`: Width of the output image (0 for original, default: 0)
- `--height`, `-h`: Height of the output image (0 for original, default: 0)

### Examples

Using long options:

```bash
./imgood --input sample.jpg --quality 70 --width 800
```

Using short options:

```bash
./imgood -i sample.jpg -q 70 -w 800
```

This will compress `sample.jpg` with 70% quality and resize it to 800px width (maintaining aspect ratio), then save it to `tmp/sample_compressed.jpg`.

## Output

The compressed image will be saved to the system's temporary directory (`/tmp` on Unix/Linux/macOS) with the naming format `{original_filename}_compressed.{extension}`.
