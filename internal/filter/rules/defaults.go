package rules

import "github.com/1broseidon/promptext/internal/filter"

func DefaultExcludes() []filter.Rule {
    return []filter.Rule{
        NewPatternRule([]string{
            ".git/",
            "node_modules/",
            "vendor/",
            ".idea/",
            ".vscode/",
            "__pycache__/",
            ".pytest_cache/",
            "dist/",
            "build/",
            "coverage/",
            "bin/",
            ".terraform/",
        }, filter.Exclude),
        NewExtensionRule([]string{
            ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp",
            ".ico", ".icns", ".svg", ".eps", ".raw", ".cr2", ".nef",
            ".exe", ".dll", ".so", ".dylib", ".bin", ".obj",
            ".zip", ".tar", ".gz", ".7z", ".rar", ".iso",
            ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
            ".class", ".pyc", ".pyo", ".pyd", ".o", ".a",
        }, filter.Exclude),
    }
}
