package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/adampresley/adamgokit/slices"
	"golang.org/x/crypto/bcrypt"
)

type RenderConfig struct {
	Base                  string
	AppName               string
	DirName               string
	GithubRepo            string
	DBName                string
	AddIdentityManagement bool
	AdminUserEmail        string
	AdminUserPassword     string
	AdminUserPasswordHash string
}

func (rc *RenderConfig) Prepare() error {
	if rc.AddIdentityManagement {
		if rc.AdminUserEmail == "" {
			return fmt.Errorf("admin user email is required when identity management is enabled")
		}

		if rc.AdminUserPassword == "" {
			return fmt.Errorf("admin user password is required when identity management is enabled")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rc.AdminUserPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error hashing admin user password: %w", err)
		}

		rc.AdminUserPasswordHash = string(hashedPassword)
	}

	return nil
}

func (rc *RenderConfig) ShouldCreateDir(path string) (bool, string) {
	base := filepath.Base(path)
	adjusted := path
	result := true

	if strings.HasPrefix(base, "_db_") {
		if rc.DBName == "" {
			result = false
		} else {
			adjusted = strings.ReplaceAll(adjusted, "_db_", "")
		}
	}

	if strings.HasPrefix(base, "_identities_") {
		if !rc.AddIdentityManagement {
			result = false
		} else {
			adjusted = strings.ReplaceAll(adjusted, "_identities_", "")
		}
	}

	return result, adjusted
}

func (rc *RenderConfig) ShouldCreateFile(path string) (bool, string) {
	base := filepath.Base(path)
	adjusted := path
	result := true

	if strings.HasPrefix(base, "_db_") {
		if rc.DBName == "" {
			result = false
		} else {
			adjusted = strings.ReplaceAll(adjusted, "_db_", "")
		}
	}

	if strings.HasPrefix(base, "_identities_") {
		if !rc.AddIdentityManagement {
			result = false
		} else {
			adjusted = strings.ReplaceAll(adjusted, "_identities_", "")
		}
	}

	return result, adjusted
}

func renderTemplates(config *RenderConfig) error {
	fs.WalkDir(_templateFS, config.Base, func(path string, d fs.DirEntry, err error) error {
		var (
			f            *os.File
			shouldCreate bool
		)

		adjustedPath := strings.TrimPrefix(path, config.Base)

		if after, ok := strings.CutPrefix(adjustedPath, string(os.PathSeparator)); ok {
			adjustedPath = after
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
			shouldCreate, adjustedPath = config.ShouldCreateDir(adjustedPath)

			if shouldCreate {
				fmt.Printf("  creating dir '%s'... \n", adjustedPath)

				newPath := filepath.Join(config.DirName, adjustedPath)

				if err = os.MkdirAll(newPath, 0755); err != nil {
					fmt.Printf("  error creating path '%s': %s\n", newPath, err.Error())
					return err
				}
			}

			return nil
		}

		/*
		 * Tis a file
		 */
		fileNameToCreate := filepath.Join(config.DirName, strings.TrimSuffix(adjustedPath, ".tmpl"))
		shouldCreate, fileNameToCreate = config.ShouldCreateFile(fileNameToCreate)

		if shouldCreate {
			fmt.Printf("  creating file '%s'... \n", fileNameToCreate)

			if shouldParseTemplate(fileNameToCreate) {
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
			} else {
				existingFile, err := _templateFS.Open(path)

				if err != nil {
					fmt.Printf("error reading file '%s': %s\n", path, err.Error())
					return err
				}

				if f, err = os.Create(fileNameToCreate); err != nil {
					fmt.Printf("error opening '%s' for creation: %s\n", adjustedPath, err.Error())
					return err
				}

				defer f.Close()

				if _, err = io.Copy(f, existingFile); err != nil {
					fmt.Printf("error copying file '%s' to '%s': %s\n", path, fileNameToCreate, err.Error())
				}
			}
		}

		return nil
	})

	return nil
}

func shouldParseTemplate(fileNameToCreate string) bool {
	dont := []string{
		".png", ".jpg", ".jpeg", ".svg", ".gif", ".ico", ".pdf",
	}

	ext := strings.ToLower(filepath.Ext(fileNameToCreate))
	return !slices.IsInSlice(ext, dont)
}

func renameCmdAppFolder(config *RenderConfig) error {
	return os.Rename(
		filepath.Join("./", config.DirName, "cmd", "renameapp"),
		filepath.Join("./", config.DirName, "cmd", config.DirName),
	)
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
