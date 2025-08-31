---
title: Changelog
description: Release notes and version history for promptext
---

All notable changes to promptext will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2025-08-31

### Major Changes
- **New Documentation System**: Migrated from Docusaurus to Astro for better performance and maintainability
- **Modern CLI Interface**: Complete refactoring with pflag for better flag handling and POSIX compliance
- **Global Configuration Support**: Added support for global and project-level configuration files
- **Performance Improvements**: Major refactoring for better performance on large codebases
- **Code Quality Improvements**: Achieved A+ Go Report Card rating with gofmt formatting fixes and cyclomatic complexity reductions

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
- Migrated to Astro-based documentation system with improved performance
- Updated GitHub Actions workflow for new documentation deployment
- Improved documentation structure and navigation

### Internal Improvements
- Refactored internal modules for better separation of concerns
- Enhanced test coverage for configuration and flag handling
- Improved error handling and user-friendly messages
- Better cross-platform compatibility
- Reduced cyclomatic complexity across core functions
- Applied consistent code formatting with gofmt

### Breaking Changes
- None - this release maintains backward compatibility

## [Previous Releases]

### Added
- Modern CLI interface with pflag integration
- Global configuration support via YAML files
- Comprehensive file filtering engine
- Project analysis with language detection
- Token counting with tiktoken integration
- Multiple output formats (Markdown, XML)
- Performance monitoring and debug logging
- Cross-platform binary distribution

### Features
- **Smart File Filtering**: .gitignore support with intelligent defaults
- **Token Analysis**: Accurate GPT-compatible token counting
- **Project Detection**: Automatic language and framework identification
- **Configuration Management**: YAML config files with CLI flag overrides
- **Performance Optimization**: Efficient processing for large codebases
- **Multiple Output Formats**: Structured Markdown and XML output

### Technical Improvements
- Major refactoring for maintainability
- Performance improvements for large codebases
- Enhanced error handling and user feedback
- Improved cross-platform compatibility

---

## Version History

For detailed version history and specific commit information, see the [GitHub Releases](https://github.com/1broseidon/promptext/releases) page.

## Contributing

When contributing to promptext, please:

1. Follow [Conventional Commits](https://www.conventionalcommits.org/) format
2. Update this changelog for significant changes
3. Include tests for new functionality
4. Update documentation as needed

## Support

- **Issues**: [GitHub Issues](https://github.com/1broseidon/promptext/issues)
- **Discussions**: [GitHub Discussions](https://github.com/1broseidon/promptext/discussions)
- **Documentation**: [promptext.dev](https://1broseidon.github.io/promptext)