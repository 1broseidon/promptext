package processor

import (
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"

    "github.com/1broseidon/promptext/internal/filter"
)

type Config struct {
    DirPath    string
    Extensions []string
    Excludes   []string
}

func ParseCommaSeparated(input string) []string {
    if input == "" {
        return nil
    }
    return strings.Split(input, ",")
}

func ProcessDirectory(config Config) (string, error) {
    var builder strings.Builder

    err := filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return nil
        }

        if !filter.ShouldProcessFile(path, config.Extensions, config.Excludes) {
            return nil
        }

        content, err := os.ReadFile(path)
        if err != nil {
            return fmt.Errorf("error reading file %s: %w", path, err)
        }

        // Add file header
        builder.WriteString(fmt.Sprintf("\n### File: %s\n", path))
        builder.WriteString("```\n")
        builder.Write(content)
        builder.WriteString("\n```\n")

        return nil
    })

    if err != nil {
        return "", fmt.Errorf("error walking directory: %w", err)
    }

    return builder.String(), nil
}
