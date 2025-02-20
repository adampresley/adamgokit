package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type RenderConfig struct {
	Base       string
	AppName    string
	GoVersion  string
	DirName    string
	GithubRepo string
	HasDB      bool
}

func renderTemplates(config *RenderConfig) error {
	fs.WalkDir(_templateFS, config.Base, func(path string, d fs.DirEntry, err error) error {
		var (
			f *os.File
		)

		adjustedPath := strings.TrimPrefix(path, config.Base)

		if strings.HasPrefix(adjustedPath, string(os.PathSeparator)) {
			adjustedPath = strings.TrimPrefix(adjustedPath, string(os.PathSeparator))
		}

		fmt.Printf("Path: %s\n", adjustedPath)

		if adjustedPath == "" {
			if err = os.Mkdir(config.DirName, 0755); err != nil {
				fmt.Printf("error creating base directory '%s': %s\n", config.DirName, err.Error())
				return err
			}

			return nil
		}

		/*
		 * A directory
		 */
		if adjustedPath != "" && d.IsDir() {
			fmt.Printf("  creating dir '%s'... \n", adjustedPath)

			newPath := filepath.Join(config.DirName, adjustedPath)

			if err = os.MkdirAll(newPath, 0755); err != nil {
				fmt.Printf("  error creating path '%s': %s\n", newPath, err.Error())
				return err
			}

			return nil
		}

		/*
		 * Tis a file
		 */
		fileNameToCreate := filepath.Join(config.DirName, strings.TrimSuffix(adjustedPath, ".tmpl"))
		templateContent, err := _templateFS.ReadFile(path)

		if err != nil {
			fmt.Printf("error reading template content for '%s': %s\n", adjustedPath, err.Error())
			return err
		}

		t, err := template.New(filepath.Base(path)).Parse(string(templateContent))

		if err != nil {
			fmt.Printf("error parsing template '%s': %s\n", path, err.Error())
			return err
		}

		if f, err = os.Create(fileNameToCreate); err != nil {
			fmt.Printf("error opening '%s' for creation: %s\n", adjustedPath, err.Error())
			return err
		}

		defer f.Close()

		if err = t.Execute(f, config); err != nil {
			fmt.Printf("error rendering template to file '%s': %s\n", adjustedPath, err.Error())
			return err
		}
		return nil
	})

	return nil
}

func renameCmdAppFolder(config *RenderConfig) error {
	return os.Rename(
		filepath.Join("./", config.DirName, "cmd", "renameapp"),
		filepath.Join("./", config.DirName, "cmd", config.DirName),
	)
}

func goModInit(config *RenderConfig) error {
	cmd := exec.Command("go", "mod", "init", config.GithubRepo)
	cmd.Dir = filepath.Join("./", config.DirName)

	b, err := cmd.Output()

	fmt.Printf("%s\n", string(b))

	if err != nil {
		return fmt.Errorf("error running go mod init: %w", err)
	}

	return nil
}

func goModTidy(config *RenderConfig) error {
	var (
		err error
		cmd *exec.Cmd
		b   []byte
	)

	cmd = exec.Command("go", "mod", "download")
	cmd.Dir = filepath.Join("./", config.DirName)

	if b, err = cmd.Output(); err != nil {
		fmt.Printf("%s\n", string(b))
		return fmt.Errorf("error running go mod download: %w", err)
	}

	fmt.Printf("%s\n", string(b))

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = filepath.Join("./", config.DirName)

	b, err = cmd.Output()
	fmt.Printf("%s\n", string(b))

	if err != nil {
		return fmt.Errorf("error running go mod tidy: %w", err)
	}

	return nil
}
