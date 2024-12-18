package rules

import "github.com/1broseidon/promptext/internal/filter/types"

func DefaultExcludes() []types.Rule {
    return []types.Rule{
        NewPatternRule([]string{
            // Version control
            ".git/",
            ".git*",
            ".svn/",
            ".hg/",
            
            // Dependencies and packages
            "node_modules/",
            "vendor/",
            "bower_components/",
            "jspm_packages/",
            
            // IDE and editor
            ".idea/",
            ".vscode/",
            ".vs/",
            "*.sublime-*",
            
            // Build and output
            "dist/",
            "build/",
            "out/",
            "bin/",
            "target/",
            
            // Cache directories
            "__pycache__/",
            ".pytest_cache/",
            ".sass-cache/",
            ".npm/",
            ".yarn/",
            
            // Test coverage
            "coverage/",
            ".nyc_output/",
            
            // Infrastructure
            ".terraform/",
            ".vagrant/",
            
            // Logs and temp
            "logs/",
            "*.log",
            "tmp/",
            "temp/",
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
