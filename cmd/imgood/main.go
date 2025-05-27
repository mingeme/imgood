package imgood

import (
	"fmt"
	"os"

	"github.com/mingeme/imgood/internal/config"
)

// PrintUsage prints the usage instructions for the tool
func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  imgood [command] [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  up\tUpload an image to S3 with optional compression")
	fmt.Println("  cp\tCopy an object in S3 with optional format conversion")
	fmt.Println("\nFor more information about a command, run:")
	fmt.Println("  imgood [command] --help")
}

// PrintUploadUsage prints the usage instructions for the upload command
func PrintUploadUsage() {
	fmt.Println("Usage:")
	fmt.Println("  imgood up [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -i, --input string\tPath to the input image file (required)")
	fmt.Println("  -k, --key string\tS3 object key (path in bucket), defaults to filename")
	fmt.Println("  -c, --compress\t\tCompress image before uploading")
	fmt.Println("  -q, --quality int\tQuality of the compressed image (1-100) (default 80)")
	fmt.Println("  -w, --width int\t\tWidth of the output image (0 for original)")
	fmt.Println("  -h, --height int\tHeight of the output image (0 for original)")
	fmt.Println("\nS3 Configuration:")
	fmt.Println("  Configure S3 settings in config.toml or using environment variables:")
	fmt.Println("  - IMGOOD_S3_BUCKET\t\tS3 bucket name")
	fmt.Println("  - IMGOOD_S3_ENDPOINT\t\tS3 endpoint URL (for non-AWS S3 services)")
	fmt.Println("  - IMGOOD_S3_REGION\t\tAWS region")
	fmt.Println("  - IMGOOD_S3_ACCESS_KEY\tAWS access key ID")
	fmt.Println("  - IMGOOD_S3_SECRET_KEY\tAWS secret access key")
}

// PrintCopyUsage prints the usage instructions for the copy command
func PrintCopyUsage() {
	fmt.Println("Usage:")
	fmt.Println("  imgood cp [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -s, --source string\tSource S3 object key to copy (required)")
	fmt.Println("  -t, --target string\tTarget S3 object key (destination), defaults to source-copy")
	fmt.Println("  -f, --format string\tConvert to format (webp, jpeg, png)")
	fmt.Println("  -q, --quality int\tQuality of the converted image (1-100) (default 80)")
	fmt.Println("  -w, --width int\t\tWidth of the output image (0 for original)")
	fmt.Println("  -h, --height int\tHeight of the output image (0 for original)")
	fmt.Println("\nS3 Configuration:")
	fmt.Println("  Configure S3 settings in config.toml or using environment variables:")
	fmt.Println("  - IMGOOD_S3_BUCKET\t\tS3 bucket name")
	fmt.Println("  - IMGOOD_S3_ENDPOINT\t\tS3 endpoint URL (for non-AWS S3 services)")
	fmt.Println("  - IMGOOD_S3_REGION\t\tAWS region")
	fmt.Println("  - IMGOOD_S3_ACCESS_KEY\tAWS access key ID")
	fmt.Println("  - IMGOOD_S3_SECRET_KEY\tAWS secret access key")
}

// Run executes the imgood command
func Run() {
	// Initialize configuration
	if err := config.Init(); err != nil {
		fmt.Printf("Warning: %s\n", err)
	}

	// Check if no arguments provided
	if len(os.Args) == 1 {
		PrintUsage()
		return
	}

	// Handle subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "up":
			// Check if help flag is provided
			if len(os.Args) > 2 && (os.Args[2] == "--help" || os.Args[2] == "-h") {
				PrintUploadUsage()
				return
			}
			os.Args = append(os.Args[:1], os.Args[2:]...)
			ExecuteUpload()
			return
		case "cp":
			// Check if help flag is provided
			if len(os.Args) > 2 && (os.Args[2] == "--help" || os.Args[2] == "-h") {
				PrintCopyUsage()
				return
			}
			os.Args = append(os.Args[:1], os.Args[2:]...)
			ExecuteCopy()
			return
		case "--help", "-h":
			PrintUsage()
			return
		default:
			fmt.Printf("Unknown command: %s\n\n", os.Args[1])
			PrintUsage()
			os.Exit(1)
		}
	}
}
