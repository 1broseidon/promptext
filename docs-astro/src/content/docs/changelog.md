---
title: Changelog
description: Release notes and version history for promptext
---

# Changelog

All notable changes to promptext are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.3] - 2025-11-07

### Changed
- **Simplified Installation**: Removed all package manager support in favor of install scripts
  - Removed Homebrew tap configuration
  - Removed Snapcraft packaging
  - Removed Chocolatey packaging
  - Removed AUR packaging
  - Streamlined to curl|bash for Linux/macOS and PowerShell for Windows
  - Reduced maintenance burden and complexity
- **Landing Page Redesign**: Modern segmented control for OS selection
  - Sleek three-way toggle with macOS, Linux, and Windows options
  - Sliding background indicator with smooth animations
  - SVG icons for visual clarity
  - Improved layout with better visual hierarchy
- **Terminal Animation Enhancement**: More realistic workflow demonstration
  - Added directory navigation (cd command) for context
  - Shows project directory in prompt
  - Enhanced readability with additional spacing
  - Better demonstration of real-world usage patterns

### Removed
- GoReleaser configurations for Homebrew, Snap, Chocolatey, and AUR
- Package manager installation instructions from documentation
- Maintenance overhead for multiple distribution channels

---

## [0.6.2] - 2025-11-07

### Added
- **Custom Domain**: Configured promptext.sh as the official documentation domain
  - DNS records configured for GitHub Pages
  - Domain verified and secured with HTTPS
  - All documentation now served from promptext.sh
- **Landing Page**: Beautiful minimal hero section with live terminal demo
  - Split-screen layout with hero text and animated terminal
  - Real-time demonstration of token budget and relevance filtering
  - Authentic output styling matching CLI behavior
  - Feature grid highlighting 6 core capabilities
  - Dark/light theme toggle with localStorage persistence
  - Fully responsive design (desktop, tablet, mobile)
- **Terminal Animation**: Interactive demo showing real-world usage
  - Live typing animation for commands
  - Green-colored success messages matching CLI output
  - Warning-colored budget exclusion details
  - Auto-plays on page load for immediate engagement

### Changed
- Documentation site now uses custom domain instead of GitHub Pages subdomain
- Hero section redesigned with better hierarchy and visual balance
- Navigation links updated to use Starlight's root-level routing
- Feature descriptions refined for clarity and impact

---

## Documentation Updates - 2025-11-06

### Changed
- **Documentation Overhaul**: Professional README redesign with centered header, clear navigation, and better organization
- **Comparison Update**: Replaced tool comparison with manual workflow comparison for accuracy
- **Progressive Examples**: Smart Context Building section now teaches simple → complex progressively
- **Generic Terminology**: Removed specific model names, using "smaller/larger context windows" instead
- **Configuration Enhancement**: Added `--init` flag documentation for easier config file generation
- **Comprehensive Changelog**: Added complete version history from v0.1.0 to v0.5.1 with detailed release notes

### Added
- Contributing section with development setup instructions
- Professional footer with community engagement links
- Improved Quick Start with numbered steps
- Complete historical changelog following Keep a Changelog format

---

## [0.5.1] - 2025-11-05

### Added
- **--init Flag**: Automatic configuration file generation with sensible defaults
  - Creates `.promptext.yml` in current directory
  - Pre-populated with common settings for detected project type
  - Interactive prompts for customization options

### Fixed
- **Security**: Addressed path traversal vulnerability in file processing
- **Testing**: Enhanced test coverage with comprehensive edge case scenarios
- **Documentation**: Clarified configuration hierarchy and precedence rules

### Changed
- Improved error messages for invalid configuration values
- Enhanced help text for CLI flags

---

## [0.5.0] - 2025-11-01

### Added
- **PTX v2.0 Enhanced Manifest**: Comprehensive project metadata in output
  - Git information (branch, commit, status)
  - Dependency analysis and versioning
  - Language and framework detection
  - Entry point identification
- **JSONL Format**: New line-delimited JSON format for streaming and large datasets
- **Deterministic Output**: Consistent file ordering for reproducible results

### Changed
- Enhanced manifest section with richer project context
- Improved directory tree representation
- Better token counting in budget summaries

---

## [0.4.6] - 2025-11-01

### Fixed
- **Update Mechanism**: Handle cross-filesystem binary replacement
  - Resolves issues when tmp directory is on different filesystem
  - Properly handles atomic file operations
  - Improved error messages for update failures

---

## [0.4.5] - 2025-11-01

### Changed
- **PTX v2.0**: Explicit file path keys for zero ambiguity
  - File paths used directly as keys (e.g., `"cmd/main.go"` not `cmd_main_go`)
  - Preserves original path separators
  - Easier to parse and reference specific files
  - Better compatibility with AI assistants

### Fixed
- Path sanitization edge cases in PTX formatter

---

## [0.4.4] - 2025-11-01

### Added
- **Self-Update Mechanism**: Built-in update functionality
  - `prx --update` to update to latest version
  - `prx --check-update` to check for new releases
  - Automatic daily update notifications (non-intrusive)
  - Smart caching to avoid excessive GitHub API calls

### Changed
- Version checking respects network failures silently
- Update notifications only shown once per day

---

## [0.4.3] - 2025-11-01

### Fixed
- Go Report Card badge cache issues
- Code quality improvements for better report card score

### Changed
- Removed unused code and improved formatting
- Enhanced code documentation

---

## [0.4.2] - 2025-11-01

### Breaking Changes
- **Default Format Changed**: PTX format is now the default output format (previously TOON)
- **Format Options Renamed**: `toon` now maps to PTX format for backward compatibility

### Major Features

#### PTX Format v1.0
Introduced PTX (Promptext Context Format) as the new default output format:

- **Hybrid Design**: Combines TOON v1.3 metadata efficiency with readable multiline code blocks
- **Token Efficiency**: 25-30% smaller than JSON while maintaining code readability
- **Debugging-Friendly**: Preserves code formatting, indentation, and line breaks
- **Specification**: Created comprehensive PTX v1.0 specification document
- **Backward Compatible**: `toon` format option maps to PTX for existing scripts

#### TOON-Strict Format
Added full TOON v1.3 specification compliance mode:

- **Maximum Compression**: 30-60% smaller than JSON through aggressive optimization
- **Full Compliance**: Follows official TOON v1.3 specification exactly
- **Escaped Strings**: Uses escaped newlines and quotes for all string content
- **Best For**: Token-limited models and metadata-heavy projects
- **Access**: Use `-f toon-strict` or `-f toon-v1.3`

### Bug Fixes
- **Fixed Relevance Filtering**: The `-r` flag now properly excludes files with score=0 when keywords are provided
- **Fixed Format Detection**: Improved file extension-based format auto-detection

### Documentation
- Created PTX Format Specification v1.0
- Updated all examples to use PTX format
- Added format selection guidelines
- Clarified differences between PTX and TOON-strict modes

### Implementation
- Renamed internal `TOONFormatter` to `PTXFormatter`
- Implemented new `TOONStrictFormatter` for TOON v1.3 compliance
- Added proper string escaping for TOON-strict mode
- Updated CLI to support new format options

## [0.4.1] - 2025-10-29

### Critical Bug Fixes
- **Fixed Token Counting Accuracy**: Resolved critical bug where token counter initialization failure returned 0 tokens for all files
- **Fixed Token Budget Display**: Corrected display showing wrong token counts when using `--max-tokens` (was showing pre-filter totals instead of actual included tokens)

### Major Features

#### Multi-Layered Lock File Detection System
Implemented sophisticated 3-layer detection system for automatic lock file exclusion:

- **Layer 1 - Signature-Based Detection (99% confidence)**
  - Pattern matching for 15+ lock file formats across all major ecosystems
  - Supported formats: `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml`, `bun.lockb`, `composer.lock`, `poetry.lock`, `Pipfile.lock`, `Gemfile.lock`, `Cargo.lock`, `go.sum`, `packages.lock.json`, and more
  - Requires multiple signature matches to reduce false positives

- **Layer 2 - Ecosystem-Aware Detection (95% confidence)**
  - Automatically detects package managers by scanning for manifest files
  - Context-aware exclusion based on detected ecosystems (Node.js, PHP, Python, Ruby, Rust, Go, .NET, Java)
  - Prevents lock file inclusion even for uncommon formats

- **Layer 3 - Generated File Detection (85% confidence)**
  - Heuristic detection of auto-generated files via markers like "@generated", "autogenerated", "do not edit"
  - Low-entropy pattern analysis for repetitive structures
  - Catches minified files, bundles, source maps, and other build artifacts

**Impact**: Automatically excludes massive lock files saving 20-40% tokens on real-world projects (50K-100K+ tokens per file)

### Token Counting System Improvements
- **Fallback Mode**: Added sophisticated approximation when tiktoken is unavailable
  - Word-based + character-based hybrid estimation
  - Code vs prose detection (~3.5 chars/token for code, ~4 for prose)
  - Graceful degradation with user notification
- **Debug Mode**: Added detailed token breakdowns with `DebugTokenCount()` for troubleshooting
- **Comprehensive Tests**: New test suite validating accuracy across different content types

### Performance Optimizations
- **Relevance Scoring**: Eliminated redundant `strings.ToLower()` calls (3x reduction in lowercase operations)
- **String Normalization**: Pre-normalize all strings once instead of per-keyword iteration
- **Config Merging**: Added deduplication for merged exclude patterns

### Code Quality
- Updated all tests to pass with new detection rules
- Enhanced error handling for edge cases
- Improved debug logging throughout the system

### Breaking Changes
- None - maintains full backward compatibility

### Token Savings Examples
| Project Type | Lock Files Excluded | Token Reduction |
|-------------|---------------------|-----------------|
| Node.js (React) | package-lock.json (500KB) | ~80,000 tokens |
| Python (Poetry) | poetry.lock (200KB) | ~60,000 tokens |
| Rust | Cargo.lock (150KB) | ~45,000 tokens |
| PHP (Laravel) | composer.lock (400KB) | ~70,000 tokens |

---

## [0.3.0] - 2025-08-31

### Major Changes
- **New Documentation System**: Migrated from Docusaurus to Astro for better performance and maintainability
- **Modern CLI Interface**: Complete refactoring with pflag for better flag handling and POSIX compliance
- **Global Configuration Support**: Added support for global and project-level configuration files
- **Performance Improvements**: Major refactoring for better performance on large codebases

### Features
- Implemented modern CLI interface with improved help text and flag descriptions
- Added global configuration support with proper precedence (CLI flags > project config > global config)
- Enhanced .gitignore patterns with comprehensive Go CLI project support
- Improved configuration merging logic for better flexibility
- Added support for XDG_CONFIG_HOME standard for configuration paths

### Bug Fixes
- Fixed CLI flag mapping for version and verbose options
- Corrected angle bracket escaping in MDX documentation
- Fixed configuration precedence issues

### Documentation
- Migrated to Astro-based documentation system
- Updated GitHub Actions workflow for new documentation deployment
- Improved documentation structure and navigation

### Internal Improvements
- Refactored internal modules for better separation of concerns
- Enhanced test coverage for configuration and flag handling
- Improved error handling and user-friendly messages
- Better cross-platform compatibility

### Breaking Changes
- None - this release maintains backward compatibility

### Migration Notes
- The documentation has moved to a new Astro-based system
- Configuration files now support both global and project-level settings
- CLI flags have been standardized for better consistency

---

## [0.4.0] - 2025-10-28

### Added
- **TOON Format**: Token-optimized output format (30-60% smaller than JSON)
- **Relevance Filtering**: Multi-factor scoring system for smart file prioritization
  - Filename matches (10x weight)
  - Directory path matches (5x weight)
  - Import statement matches (3x weight)
  - Content matches (1x weight)
- **Token Budget Management**: Hard limits with `--max-tokens` flag
- **Budget Visualization**: Clear reporting of included/excluded files

### Changed
- Major refactoring of output generation system
- Enhanced file filtering logic

---

## [0.3.0] - 2025-08-31

### Added
- **Astro Documentation System**: Migrated from Docusaurus for better performance
- **Global Configuration Support**: `~/.config/promptext/config.yml` for user-wide defaults
- **XDG_CONFIG_HOME Support**: Respects XDG Base Directory specification

### Changed
- **Modern CLI Interface**: Complete refactoring with pflag for better flag handling
- **Configuration Precedence**: CLI flags > project config > global config
- Enhanced .gitignore pattern support
- Improved help text and flag descriptions

### Fixed
- CLI flag mapping for version and verbose options
- Angle bracket escaping in documentation
- Configuration precedence edge cases

---

## [0.2.6] - 2024-12-19

### Added
- **Enhanced File Type Detection**: Detailed categorization with size information
- Better language and framework identification

### Changed
- Improved project analysis output

---

## [0.2.5] - 2024-12-19

### Added
- **Build-Time Version Management**: Automatic version injection during builds
- Version information in CLI output

---

## [0.2.4] - 2024-12-19

### Fixed
- Documentation clarity for `UseDefaultRules` configuration option

---

## [0.2.3] - 2024-12-18

### Changed
- Reorganized project info retrieval for better performance
- Improved logging format consistency

---

## [0.2.2] - 2024-12-18

### Changed
- Code formatting cleanup across codebase
- Improved consistency in code style

---

## [0.2.1] - 2024-12-17

### Added
- Release process improvements
- GitHub Actions workflow enhancements

---

## [0.2.0] - 2024-12-17

### Added
- Initial stable release with core functionality
- Directory processing and file filtering
- Configuration file support
- Multiple output formats (Markdown, XML)

---

## [0.1.9] - 2024-12-16

### Added
- Python and Go sample projects
- Example configuration files
- Entry point detection examples

---

## [0.1.8] - 2024-12-16

### Added
- Cross-platform binary releases via GitHub Actions
- Automated build pipeline

---

## [0.1.0] - 2024-12-16

### Added
- Initial release of promptext
- Core directory processing functionality
- Basic file filtering
- Markdown output format