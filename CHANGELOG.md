# Changelog

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

## [0.2.6] - Previous Release
- Last stable release before major refactoring