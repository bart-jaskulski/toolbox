package cmd

import (
	"fmt"
	"os"

	"github.com/bart-jaskulski/toolbox/bm/internal/bookmarks"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bm",
	Short: "Bookmark manager",
	Long:  `Bookmark manager for managing bookmarks in the ~/.bookmarks file.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			bookmarks, err := bookmarks.FindBookmarks("")
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, bookmark := range bookmarks {
				fmt.Println(bookmark)
			}
			return
		}

		if err := bookmarks.AddBookmark(args[0]); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Bookmark added successfully")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
