package cmd

import (
	"fmt"

	"github.com/bart-jaskulski/toolbox/bm/internal/bookmarks"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:     "find",
	Aliases: []string{"f"},
	Short:   "Search for bookmarks containing the given pattern",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchPattern := ""
		if len(args) > 0 {
			searchPattern = args[0]
		}

		bookmarks, err := bookmarks.FindBookmarks(searchPattern)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(bookmarks) == 0 {
			fmt.Println("No bookmarks found")
			return
		}

		for _, bookmark := range bookmarks {
			fmt.Println(bookmark)
		}
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
