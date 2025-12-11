package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bart-jaskulski/toolbox/bm/internal/bookmarks"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new bookmark",
	Long:  `Add a new bookmark to the .bookmark file`,
	Run: func(cmd *cobra.Command, args []string) {
		var urlStr string
		if len(args) < 1 {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			urlStr = scanner.Text()
		} else {
			urlStr = args[0]
		}
		if err := bookmarks.AddBookmark(urlStr); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Bookmark added successfully")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
