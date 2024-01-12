package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

func createFileByTemplate(templatePath, filePath string, data interface{}) error {
	t, err := parseFiles(template.FuncMap{"add": add}, templatePath)
	if err != nil {
		return fmt.Errorf("unable to parse template: %w", err)
	}

	f, err := createFile(filePath, err)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return fmt.Errorf("unable to execute template: %w", err)
	}

	return nil
}

func parseFiles(funcs template.FuncMap, filenames ...string) (*template.Template, error) {
	return template.New(filepath.Base(filenames[0])).Funcs(funcs).ParseFiles(filenames...)
}

func add(a int, b int) int {
	return a + b
}

func createJsonFile(filePath string, obj interface{}) error {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal json: %w", err)
	}

	f, err := createFile(filePath, err)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = ioutil.WriteFile(filePath, data, 0o777); err != nil {
		return fmt.Errorf("unable to write file: %w", err)
	}

	return nil
}

func createFile(filePath string, err error) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o777); err != nil {
		return nil, fmt.Errorf("unable to create directory: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to create file: %w", err)
	}
	return f, nil
}
