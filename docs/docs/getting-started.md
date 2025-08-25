---
sidebar_position: 2
---

# Getting Started

## Installation

### Prerequisites

#### All Platforms

- Git (for version control features)

#### Platform-Specific

- **Linux/macOS**: No additional requirements
- **Windows**:
  - PowerShell 5.1 or higher
- **Go Installation**: Go 1.22 or higher (if installing via `go install`)

### Installation Methods

1. Quick Install:

**Linux/macOS**:

Install:

```bash
curl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.sh | bash
```

Additional options:

```bash
curl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.sh | bash -h
```

The Linux/macOS installer provides:

- 🏠 User-level installation in ~/.local/bin
- 🔒 Automatic checksum verification
- 📁 Custom installation directory support
- ⚡ PATH environment configuration
- 💫 Shell alias configuration
- 🔄 Clean uninstallation

**Windows (PowerShell)**:

Install:

```powershell
irm https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.ps1 | iex
```

Uninstall:

```powershell
irm https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.ps1 | iex -Uninstall
```

Both installers provide:

- 🏠 User-level installation (in ~/.local/bin or %LOCALAPPDATA%)
- 🔒 Automatic checksum verification
- 📁 Custom installation directory support
- ⚡ PATH environment configuration
- 💫 Shell/PowerShell alias creation (prx)
- 🔄 Clean uninstallation

2. Using Go Install:

```bash
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

3. Manual Installation:

Download the appropriate binary for your platform from the [releases page](https://github.com/1broseidon/promptext/releases):

**Linux/macOS**:

- Download the appropriate binary
- Make it executable: `chmod +x promptext`
- Create directory: `mkdir -p ~/.local/bin`
- Move to PATH: `mv promptext ~/.local/bin/`

**Windows**:

- Download the Windows binary (ZIP file)
- Extract to a directory (e.g., `%LOCALAPPDATA%\promptext`)
- Add the directory to your PATH:
  - System Settings > Advanced > Environment Variables
  - Edit the User PATH variable
  - Add the installation directory

## Usage

### Basic Command Structure

```bash
promptext [flags]
```

### Available Flags

- `--version`: Show version information
- `-directory, -d string`: Directory to process (default ".")
- `-extension, -e string`: File extensions to include (comma-separated, e.g., ".go,.js")
- `-exclude, -x string`: Patterns to exclude (comma-separated)
- `-format, -f string`: Output format (markdown/xml)
- `-output, -o string`: Output file path
- `--no-copy`: Disable automatic clipboard copying
- `-info, -i`: Show only project summary with token counts
- `-verbose, -v`: Show full file contents in terminal
- `-debug, -D`: Enable debug logging
- `-gitignore, -g`: Use .gitignore patterns (default true)
- `-use-default-rules, -u`: Use default filtering rules (default true)
- `-help, -h`: Show help message

### Examples

1. Process specific file types:

```bash
promptext -extension .go,.js
```

2. Export as XML with debug info:

```bash
promptext -format xml -output project.xml -debug
```

3. Show project overview with token counts:

```bash
promptext -info
```

4. Process with exclusions and verbose output:

```bash
promptext -exclude "test/,vendor/" -v
```

5. Check version:

```bash
promptext --version  # Show version information

# Example output:
# promptext version v0.2.4 (2024-12-19)
```

6. Process without copying to clipboard:

```bash
promptext --no-copy -output project.md
```

7. Process all files including dependencies:

```bash
promptext -u=false -exclude "test/" # Disable default rules but keep test/ excluded
```
