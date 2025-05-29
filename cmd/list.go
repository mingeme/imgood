package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mingeme/imgood/internal/config"
	"github.com/mingeme/imgood/internal/s3"
)

var (
	listPrefix    string
	listLimit     int32
	listSortBy    string
	listDescending bool
	listShowURLs   bool
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List objects in S3 bucket with filtering and sorting options",
	Long: `List objects in S3 bucket with filtering and sorting options.
	
Example:
  imgood ls -p images/ -l 50 -s size -d -u`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get S3 configuration and create client
		s3Config := config.GetS3Config()
		s3Client, err := s3.NewClient(s3Config)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			fmt.Println("Check your S3 configuration in config.toml or environment variables")
			os.Exit(1)
		}

		// List objects from S3
		fmt.Printf("Listing objects in bucket '%s'", s3Config.Bucket)
		if listPrefix != "" {
			fmt.Printf(" with prefix '%s'", listPrefix)
		}
		fmt.Println()

		objects, err := s3Client.ListObjects(listPrefix, listLimit)
		if err != nil {
			fmt.Printf("Error listing objects: %s\n", err)
			os.Exit(1)
		}

		// Sort objects
		sortObjects(objects, listSortBy, listDescending)

		// Display results
		if len(objects) == 0 {
			fmt.Println("No objects found.")
			return
		}

		// Print header
		fmt.Printf("%-40s %-15s %-20s", "KEY", "SIZE", "LAST MODIFIED")
		if listShowURLs {
			fmt.Printf(" %-60s", "URL")
		}
		fmt.Println()
		fmt.Println(strings.Repeat("-", 80))

		// Print objects
		for _, obj := range objects {
			// Format the key for display (truncate if too long)
			displayKey := obj.Key
			if len(displayKey) > 38 {
				displayKey = "..." + displayKey[len(displayKey)-35:]
			}

			// Format size
			sizeStr := formatBytes(obj.Size)

			// Format date
			dateStr := obj.LastModified.Format("2006-01-02 15:04:05")

			// Print the object info
			fmt.Printf("%-40s %-15s %-20s", displayKey, sizeStr, dateStr)
			if listShowURLs {
				fmt.Printf(" %s", obj.URL)
			}
			fmt.Println()
		}

		fmt.Printf("\nTotal: %d objects\n", len(objects))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Define command line flags for listing
	listCmd.Flags().StringVarP(&listPrefix, "prefix", "p", "", "Prefix filter for S3 objects")
	listCmd.Flags().Int32VarP(&listLimit, "limit", "l", 100, "Maximum number of objects to list")
	listCmd.Flags().StringVarP(&listSortBy, "sort", "s", "name", "Sort by: name, size, date")
	listCmd.Flags().BoolVarP(&listDescending, "desc", "d", false, "Sort in descending order")
	listCmd.Flags().BoolVarP(&listShowURLs, "urls", "u", false, "Show full URLs")

	// Add shell completion for sort flag
	_ = listCmd.RegisterFlagCompletionFunc("sort", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"name", "size", "date"}, cobra.ShellCompDirectiveNoFileComp
	})
}

// formatBytes formats bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func sortObjects(objects []s3.S3Object, sortBy string, descending bool) {
	switch strings.ToLower(sortBy) {
	case "size":
		if descending {
			sort.Slice(objects, func(i, j int) bool {
				return objects[i].Size > objects[j].Size
			})
		} else {
			sort.Slice(objects, func(i, j int) bool {
				return objects[i].Size < objects[j].Size
			})
		}
	case "date":
		if descending {
			sort.Slice(objects, func(i, j int) bool {
				return objects[i].LastModified.After(objects[j].LastModified)
			})
		} else {
			sort.Slice(objects, func(i, j int) bool {
				return objects[i].LastModified.Before(objects[j].LastModified)
			})
		}
	default: // "name"
		if descending {
			sort.Slice(objects, func(i, j int) bool {
				return objects[i].Key > objects[j].Key
			})
		} else {
			sort.Slice(objects, func(i, j int) bool {
				return objects[i].Key < objects[j].Key
			})
		}
	}
}
