package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/bart-jaskulski/toolbox/bm/internal/bookmarks"
	"github.com/bart-jaskulski/toolbox/bm/internal/utils"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:     "open [bookmark]",
	Aliases: []string{"o"},
	Short:   "Open a bookmark in the default web browser",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookmark := args[0]
		var url string

		if !utils.IsValidURL(bookmark) {
			chosenBookmark, err := bookmarks.SearchAndChooseBookmark(bookmark)
			if err != nil {
				return err
			}
			bookmark = chosenBookmark
		}

		url = strings.Split(bookmark, " ")[0]

		var openCmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			openCmd = exec.Command("open", url)
		case "windows":
			openCmd = exec.Command("cmd", "/c", "start", url)
		default: // Assume Linux
			openCmd = exec.Command("xdg-open", url)
		}
		err := openCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to open URL: %v", err)
		}
		fmt.Println("URL opened in browser:", url)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
