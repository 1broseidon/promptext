# Performance Benchmarks & Thresholds

This document defines performance benchmarks, acceptable thresholds, and regression testing guidelines for the promptext CLI tool.

## Overview

The performance benchmarks validate five critical areas:
1. **File Processing Pipeline** - End-to-end directory processing
2. **Binary Detection** - Optimized binary file identification  
3. **Token Counting** - tiktoken-based estimation performance
4. **Memory Usage** - Memory allocation patterns and limits
5. **Filter Performance** - File filtering with complex patterns

## Benchmark Organization

### Location
- **Processor**: `/internal/processor/processor_bench_test.go`
- **Binary Detection**: `/internal/filter/rules/binary_bench_test.go` 
- **Token Counting**: `/internal/token/tiktoken_bench_test.go`
- **Filter Performance**: `/internal/filter/filter_bench_test.go`

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./internal/...

# Run specific component benchmarks
go test -bench=BenchmarkProcessDirectory ./internal/processor
go test -bench=BenchmarkBinaryRule ./internal/filter/rules
go test -bench=BenchmarkTokenCounter ./internal/token
go test -bench=BenchmarkFilter ./internal/filter

# Run with memory profiling
go test -bench=. -benchmem ./internal/processor

# Generate CPU and memory profiles
go test -bench=BenchmarkProcessDirectory_1000Files -cpuprofile=cpu.prof -memprofile=mem.prof ./internal/processor
```

## Performance Thresholds

### 1. File Processing Pipeline

#### Acceptable Thresholds
| Metric | 100 Files | 1,000 Files | 10,000 Files |
|--------|-----------|-------------|--------------|
| **Processing Time** | < 50ms | < 500ms | < 5s |
| **Memory Usage** | < 10MB | < 50MB | < 200MB |
| **Files/Second** | > 2,000 | > 2,000 | > 2,000 |
| **Memory/File** | < 100KB | < 50KB | < 20KB |

#### Critical Thresholds (Failure Points)
- **Processing Time**: > 2x acceptable threshold
- **Memory Usage**: > 3x acceptable threshold  
- **Memory Leaks**: > 10MB growth per iteration

#### Benchmark Coverage
```bash
BenchmarkProcessDirectory_100Files     # Small projects
BenchmarkProcessDirectory_1000Files    # Medium projects  
BenchmarkProcessDirectory_10000Files   # Large codebases
BenchmarkProcessDirectory_ManySmallFiles  # 5k files @ 500B each
BenchmarkProcessDirectory_FewLargeFiles   # 100 files @ 50KB each
```

### 2. Binary Detection Performance

#### Acceptable Thresholds
| Metric | 100 Files | 1,000 Files | 10,000 Files |
|--------|-----------|-------------|--------------|
| **Detection Time** | < 10ms | < 50ms | < 200ms |
| **Extension Check** | < 1ms | < 5ms | < 20ms |
| **Content Analysis** | < 20ms | < 100ms | < 500ms |
| **Memory Usage** | < 1MB | < 5MB | < 20MB |

#### Three-Stage Performance Validation
1. **Extension Check** (fastest): O(1) map lookup, no I/O
2. **Size Check** (fast): Single stat call, no content read
3. **Content Analysis** (slowest): 512-byte content read as last resort

#### Benchmark Coverage
```bash
BenchmarkBinaryRule_Match_*Files       # Scale testing
BenchmarkBinaryRule_ExtensionOnly       # Stage 1 performance
BenchmarkBinaryRule_ContentAnalysis     # Stage 3 performance
BenchmarkBinaryRule_LargeFiles          # 10MB+ file handling
BenchmarkBinaryRule_OptimizedVsNaive    # Performance comparison
```

### 3. Token Counting Performance

#### Acceptable Thresholds
| Metric | 1KB Text | 10KB Text | 100KB Text | 1MB Text |
|--------|----------|-----------|------------|----------|
| **Counting Time** | < 1ms | < 5ms | < 25ms | < 100ms |
| **Tokens/Second** | > 1M | > 500K | > 100K | > 50K |
| **Memory Usage** | < 1MB | < 5MB | < 20MB | < 50MB |
| **Memory/Token** | < 10B | < 20B | < 50B | < 100B |

#### Content Type Performance
- **Code**: Highest token density, complex syntax
- **Markdown**: Medium density, formatting tokens
- **JSON**: Structured data, predictable patterns
- **Plain Text**: Lowest token density, natural language

#### Benchmark Coverage
```bash
BenchmarkTokenCounter_*Text             # Size scaling
BenchmarkTokenCounter_*Content          # Content type variation
BenchmarkTokenCounter_ManySmallTexts    # 1k files @ 500B
BenchmarkTokenCounter_FewLargeTexts     # 10 files @ 50KB
BenchmarkTokenCounter_Concurrent        # Thread safety
```

### 4. Memory Usage Patterns

#### Acceptable Memory Patterns
- **Peak Memory**: < 3x working set size
- **Memory Growth**: < 1% per 1k files processed
- **GC Pressure**: < 10MB/s allocation rate
- **Memory Efficiency**: > 80% useful allocations

#### Memory Profile Monitoring
```bash
# Profile memory allocation patterns
go test -bench=BenchmarkProcessDirectory_MemoryProfile -memprofile=mem.prof

# Analyze with pprof
go tool pprof mem.prof
(pprof) top10
(pprof) list ProcessDirectory
```

#### Critical Memory Issues
- **Memory Leaks**: Growing memory usage across iterations
- **Excessive Allocations**: > 1MB allocations for < 100KB files
- **GC Thrashing**: > 50% time spent in garbage collection

### 5. Filter Performance

#### Acceptable Thresholds
| Metric | Simple Rules | Medium Rules | Complex Rules |
|--------|--------------|--------------|---------------|
| **Filter Creation** | < 1ms | < 5ms | < 20ms |
| **File Check Rate** | > 100K/s | > 50K/s | > 10K/s |
| **Memory Usage** | < 1MB | < 5MB | < 20MB |
| **Pattern Match** | < 10μs | < 50μs | < 200μs |

#### Rule Complexity Definitions
- **Simple**: 2-5 patterns, basic extensions
- **Medium**: 5-15 patterns, directory exclusions  
- **Complex**: 15+ patterns, nested paths, gitignore

#### Benchmark Coverage
```bash
BenchmarkFilter_Creation_*Rules         # Rule compilation cost
BenchmarkFilter_ShouldProcess_*Files    # Scale testing
BenchmarkFilter_PatternMatching_*       # Pattern complexity
BenchmarkFilter_Concurrent             # Thread safety
```

## Regression Testing Guidelines

### Continuous Integration Checks

#### Performance Gate (CI Pipeline)
```bash
# Run performance regression checks
go test -bench=. -benchtime=5s -count=3 ./internal/... > current_bench.txt

# Compare with baseline (stored in repo)
benchcmp baseline_bench.txt current_bench.txt

# Fail if performance degrades > 25%
```

#### Automated Performance Monitoring
- **Baseline Updates**: Monthly or after major changes
- **Regression Alerts**: > 15% performance degradation
- **Memory Leak Detection**: > 5MB growth per benchmark iteration

### Performance Testing Checklist

#### Before Release
- [ ] All benchmarks pass acceptable thresholds
- [ ] No memory leaks detected across 100 iterations
- [ ] Performance comparison with previous release < 25% degradation
- [ ] Large codebase testing (>10k files) validates scalability
- [ ] Memory profile shows reasonable allocation patterns

#### After Major Changes
- [ ] Affected component benchmarks run and pass
- [ ] End-to-end performance testing on realistic projects
- [ ] Memory usage analysis with profiling tools
- [ ] Concurrent usage validation (thread safety)

## Performance Optimization Guidelines

### Optimization Priorities
1. **I/O Minimization**: Reduce file system operations
2. **Memory Efficiency**: Minimize allocations per file
3. **Algorithm Complexity**: Prefer O(1) and O(log n) operations  
4. **Caching**: Cache expensive computations (binary detection, tiktoken)
5. **Concurrency**: Safe parallel processing where beneficial

### Code Review Performance Checklist
- [ ] File operations use efficient patterns (stat before read)
- [ ] Memory allocations are proportional to work done
- [ ] Loops avoid N² complexity with large file sets
- [ ] Binary detection uses three-stage approach
- [ ] Token counting batches operations efficiently
- [ ] Filter rules compile once and reuse

## Benchmark Data Interpretation

### Key Performance Indicators
- **Throughput**: Files processed per second
- **Latency**: Time per file operation
- **Memory Efficiency**: Memory used per file processed
- **Scalability**: Performance degradation with size increase

### Warning Signs
- **Sub-linear Scaling**: Performance degrades faster than O(n)
- **Memory Growth**: Memory usage grows without bound
- **High Allocation Rate**: Excessive garbage collection pressure
- **Thread Contention**: Poor concurrent performance

### Acceptable Performance Ranges

#### Small Projects (< 100 files)
- Interactive response time: < 100ms total
- Memory usage: < 10MB peak
- No noticeable delay for user

#### Medium Projects (100-1000 files)  
- CLI tool response: < 1s total
- Memory usage: < 50MB peak
- Suitable for development workflow

#### Large Projects (1000-10000 files)
- Batch processing time: < 10s total  
- Memory usage: < 200MB peak
- Acceptable for CI/CD pipelines

#### Enterprise Projects (10000+ files)
- Extended processing: < 60s total
- Memory usage: < 500MB peak  
- Resource usage scales reasonably

## Tools and Commands

### Performance Analysis Tools
```bash
# CPU profiling
go test -bench=BenchmarkProcessDirectory_1000Files -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Memory profiling  
go test -bench=BenchmarkProcessDirectory_1000Files -memprofile=mem.prof
go tool pprof mem.prof

# Block profiling (concurrency)
go test -bench=BenchmarkProcessDirectory_1000Files -blockprofile=block.prof
go tool pprof block.prof

# Trace analysis
go test -bench=BenchmarkProcessDirectory_1000Files -trace=trace.out
go tool trace trace.out
```

### Benchmark Comparison
```bash
# Generate baseline
go test -bench=. ./internal/... > baseline.txt

# Compare after changes
go test -bench=. ./internal/... > current.txt
benchcmp baseline.txt current.txt

# Statistical analysis
go test -bench=BenchmarkProcessDirectory_1000Files -count=10 | benchstat
```

### Memory Leak Detection
```bash
# Run extended test for leaks
go test -bench=BenchmarkProcessDirectory_MemoryProfile -benchtime=30s -memprofile=mem.prof

# Check for growing memory usage
go tool pprof -alloc_space mem.prof
go tool pprof -inuse_space mem.prof
```

## Integration with Development Workflow

### Pre-commit Hooks
```bash
# Quick performance check
make bench-critical

# Full performance validation  
make bench-full
```

### Development Guidelines
- Run relevant benchmarks after performance-related changes
- Profile before optimizing to identify actual bottlenecks
- Validate performance improvements with benchmarks
- Document performance implications in pull requests

This benchmark suite ensures promptext maintains production-ready performance across all supported use cases while preventing performance regressions during development.