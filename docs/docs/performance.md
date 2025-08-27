---
sidebar_position: 8
---

# Performance

## Optimization Features

### High-Performance File Processing

- **Fast binary detection**: Significantly improved binary file detection speeds up large codebase processing
- **Concurrent file processing**: Parallel processing of multiple files
- **Optimized token counting**: Efficient GPT token estimation using tiktoken
- **Smart caching system**: Reduces redundant operations

### Memory Management

- **Efficient file reading**: Streamlined file I/O operations
- **Memory-conscious processing**: Optimized memory usage for large codebases
- **Binary file filtering**: Fast exclusion of non-text files saves processing time

## Performance Monitoring

### Debug Mode

Enable with `-debug` or `-D` flag for detailed performance information:

```bash
promptext -debug
# or
promptext -D
```

Debug output includes:

- **File processing times**: See how long each stage takes
- **Binary detection performance**: Monitor the speed of file type detection
- **Token counting metrics**: Detailed tiktoken processing information
- **Memory usage tracking**: Monitor resource consumption
- **Filter operation timing**: See time spent on file filtering

### Performance Benefits You'll Notice

- **Faster startup**: Improved binary detection means quicker processing of large directories
- **Better responsiveness**: Enhanced file filtering reduces time spent on irrelevant files
- **Efficient token counting**: Optimized tiktoken usage for accurate GPT token estimates

## Performance Tips

### For Large Codebases

1. **Use specific file extensions**: `-extension .go,.js` processes only relevant files
2. **Configure smart exclusions**: `-exclude "vendor/,node_modules/,dist/"` skips dependency directories
3. **Leverage .gitignore**: Default `.gitignore` integration automatically excludes build artifacts
4. **Use info mode for quick overview**: `-info` flag shows project summary without processing all content

### Typical Performance Expectations

- **Small projects** (&lt;100 files): Near-instantaneous processing
- **Medium projects** (100-1000 files): Processing in seconds
- **Large codebases** (1000+ files): Optimized filtering and binary detection provide significant speed improvements

### Debug Mode for Optimization

Use debug mode to identify performance bottlenecks in your specific project:

```bash
promptext -D -info  # Quick performance overview
promptext -D        # Full processing with timing details
```
