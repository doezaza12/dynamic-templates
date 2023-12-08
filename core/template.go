package core

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func GetTemplateFunction() map[string]any {
	return map[string]any{
		"join":      strings.Join,
		"split":     strings.Split,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"base":      filepath.Base,
		"dir":       filepath.Dir,
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"toString":  strconv.Itoa,
		"quote": func(text string) string {
			return fmt.Sprintf("\"%v\"", text)
		},
	}
}

func RenderTemplate(templateName string, outputDir string, name string, data any) error {

	templateFunc := GetTemplateFunction()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if outputDir == "" {
		outputDir = filepath.Join(homeDir, APP_STAGING_DIR, DEFAULT_RENDER_DIR)
	}

	templatePath := filepath.Join(homeDir, APP_STAGING_DIR, DEFAULT_TEMPLATE_DIR, templateName)
	outputPath := filepath.Join(outputDir, name)

	if err := os.MkdirAll(outputPath, DEFAULT_FILE_PERMISSION); err != nil {
		if os.IsExist(err) {
			fmt.Printf("%v already existed", outputPath)
			os.RemoveAll(outputPath)
		} else {
			panic(err)
		}
	}

	templateFilesPath := filepath.Join(templatePath, DEFAULT_TEMPLATE_CONTENT)
	templateFiles := []string{}
	var actualDir string

	/*
		collect every template files (.tpl extension)
		this will ignore other extensions
	*/
	if err := filepath.WalkDir(templateFilesPath, func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("%s %s\n", path, d.Name())
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".tpl") {
			templateFiles = append(templateFiles, path)
		} else if d.IsDir() {
			actualDir = filepath.Join(outputPath, strings.TrimPrefix(path, templateFilesPath))
			if err := os.MkdirAll(actualDir, DEFAULT_FILE_PERMISSION); err != nil {
				if os.IsExist(err) {
					return nil
				} else {
					panic(err)
				}
			}
		}
		return err
	}); err != nil {
		panic(err)
	}

	for _, path := range templateFiles {
		templateName := filepath.Base(path)
		tmpl := template.Must(template.New(templateName).Funcs(templateFunc).ParseFiles(path))

		actualFile := strings.TrimPrefix(path, templateFilesPath)
		actualFile = strings.TrimSuffix(actualFile, ".tpl")

		out, err := os.Create(filepath.Join(outputPath, actualFile))
		if err != nil {
			panic(err)
		}
		if err = tmpl.Execute(out, data); err != nil {
			panic(err)
		}
	}

	hookPath := filepath.Join(templatePath, DEFAULT_TEMPLATE_HOOK)

	if err := filepath.WalkDir(hookPath, func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("%s %s\n", path, d.Name())
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sh") {
			templateName := filepath.Base(path)
			tmpl := template.Must(template.New(templateName).Funcs(templateFunc).ParseFiles(path))

			var rawHookCmd bytes.Buffer
			err = tmpl.Execute(&rawHookCmd, data)
			if err != nil {
				panic(err)
			}

			cmd := exec.Command(DEFAULT_SHELL, "-c", rawHookCmd.String())
			cmd.Dir = outputPath

			out, err := cmd.Output()
			if err != nil {
				panic(err)
			}
			fmt.Println(out)
		}
		return err
	}); err != nil {
		panic(err)
	}

	return nil
}
