/*
Copyright © 2024 Adam Presley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

const (
	GO_VERSION = "1.23.3"
)

var _templateFS embed.FS

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "adamgokit",
	Short: "A generator tool for creating app using this toolkit",
	Long: `This tool generates applications that can use 'adamgokit'.
It provides flags for additional add-ons, like CSS and JS frameworks.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err              error
			appName          string
			githubPath       string
			selectedTemplate string
			wantDB           bool
			dbName           string
		)

		appTemplateOptions := []string{
			"Basic",
			"Web",
		}

		pterm.DefaultBigText.WithLetters(putils.LettersFromString("Go Kit")).Render()
		pterm.Println()

		pterm.Printf("Enter application name (no spaces, lowercase):\n")
		appName, _ = pterm.DefaultInteractiveTextInput.Show()

		pterm.Printf("\nEnter Github path (i.e. github.com/username/appname):\n")
		githubPath, _ = pterm.DefaultInteractiveTextInput.Show()

		pterm.Printf("\nChoose application template:\n")
		selectedTemplate, _ = pterm.DefaultInteractiveSelect.WithOptions(appTemplateOptions).Show()

		if selectedTemplate == "Web" {
			wantDB, _ = pterm.DefaultInteractiveConfirm.Show("\nDo you want to use a database?\n")

			if wantDB {
				dbName = "sqlite"
			}
		}

		pterm.Printf("\nCreating '%s' app '%s' with repository '%s'...\n", selectedTemplate, appName, githubPath)

		config := &RenderConfig{
			AppName:    appName,
			GoVersion:  GO_VERSION,
			Base:       fmt.Sprintf("templates/%s", strings.ToLower(selectedTemplate)),
			DirName:    appName,
			GithubRepo: githubPath,
			DBName:     dbName,
		}

		if err = buildApp(config); err != nil {
			pterm.DefaultLogger.Error("There was an error building the application", pterm.DefaultLogger.Args("error", err))
			return
		}
	},
}

func buildApp(config *RenderConfig) error {
	var (
		err error
	)

	if err = renderTemplates(config); err != nil {
		return fmt.Errorf("error rendering application templates: %w", err)
	}

	if err = renameCmdAppFolder(config); err != nil {
		return fmt.Errorf("error renaming cmd app folder '%s': %w", config.DirName, err)
	}

	if config.DBName != "" {
		if err = os.MkdirAll(filepath.Join("./", config.DirName, "cmd", config.DirName, "data"), 0755); err != nil {
			return fmt.Errorf("error making data directory: %w", err)
		}
	}

	if err = goModTidy(config); err != nil {
		return fmt.Errorf("error downloading dependencies: %w", err)
	}

	return nil
}

func Execute(templateFS embed.FS) {
	_templateFS = templateFS

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
