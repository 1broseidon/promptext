package rules

import "github.com/1broseidon/promptext/internal/filter/types"

func DefaultExcludes() []types.Rule {
    return []types.Rule{
        NewPatternRule([]string{
            ".git/",
            ".git*",
            "node_modules/",
            "vendor/",
            ".idea/",
            ".vscode/",
            "__pycache__/",
            ".pytest_cache/",
            ".aider*/",
            ".aider.*",
            "dist/",
            "build/",
            "coverage/",
            "bin/",
            ".terraform/",
        }, types.Exclude),
        NewExtensionRule([]string{
            // Binary and image files
            ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp",
            ".ico", ".icns", ".svg", ".eps", ".raw", ".cr2", ".nef",
            ".exe", ".dll", ".so", ".dylib", ".bin", ".obj",
            ".zip", ".tar", ".gz", ".7z", ".rar", ".iso",
            ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
            ".class", ".pyc", ".pyo", ".pyd", ".o", ".a", ".db",
            ".db-shm", ".db-wal", ".DS_Store",
        }, types.Exclude),
    }
}
