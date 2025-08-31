# Promptext Performance Testing Makefile

.PHONY: bench bench-all bench-critical bench-processor bench-binary bench-token bench-filter bench-memory
.PHONY: bench-compare bench-profile bench-trace bench-stats clean-bench help

# Default benchmark run
bench: bench-critical

# Run critical performance benchmarks (fast, for CI)
bench-critical:
	@echo "Running critical performance benchmarks..."
	@go test -run=^$$ -bench=BenchmarkProcessDirectory_1000Files -benchtime=3s ./internal/processor
	@go test -run=^$$ -bench=BenchmarkBinaryRule_Match_1000Files -benchtime=3s ./internal/filter/rules  
	@go test -run=^$$ -bench=BenchmarkTokenCounter_MediumText -benchtime=3s ./internal/token
	@go test -run=^$$ -bench=BenchmarkFilter_ShouldProcess_1000Files -benchtime=3s ./internal/filter

# Run all performance benchmarks (comprehensive)
bench-all:
	@echo "Running all performance benchmarks..."
	@go test -run=^$$ -bench=. -benchtime=5s -benchmem ./internal/processor
	@go test -run=^$$ -bench=. -benchtime=5s -benchmem ./internal/filter/rules
	@go test -run=^$$ -bench=. -benchtime=5s -benchmem ./internal/token  
	@go test -run=^$$ -bench=. -benchtime=5s -benchmem ./internal/filter

# Component-specific benchmarks
bench-processor:
	@echo "Running file processing benchmarks..."
	@go test -run=^$$ -bench=BenchmarkProcessDirectory -benchmem ./internal/processor
	@go test -run=^$$ -bench=BenchmarkProcessFile -benchmem ./internal/processor
	@go test -run=^$$ -bench=BenchmarkWalkDir -benchmem ./internal/processor

bench-binary:
	@echo "Running binary detection benchmarks..."
	@go test -run=^$$ -bench=BenchmarkBinaryRule -benchmem ./internal/filter/rules

bench-token:
	@echo "Running token counting benchmarks..."
	@go test -run=^$$ -bench=BenchmarkTokenCounter -benchmem ./internal/token

bench-filter:
	@echo "Running filter performance benchmarks..."
	@go test -run=^$$ -bench=BenchmarkFilter -benchmem ./internal/filter

# Memory-focused benchmarks
bench-memory:
	@echo "Running memory usage benchmarks..."
	@go test -run=^$$ -bench=BenchmarkProcessDirectory_MemoryProfile -benchmem ./internal/processor
	@go test -run=^$$ -bench=BenchmarkTokenCounter_MemoryProfile -benchmem ./internal/token
	@go test -run=^$$ -bench=BenchmarkFilter_MemoryUsage -benchmem ./internal/filter

# Benchmark comparison (requires baseline.txt)
bench-compare: bench-current
	@echo "Comparing benchmark results..."
	@if [ -f baseline.txt ]; then \
		benchcmp baseline.txt current.txt || echo "benchcmp not installed. Run: go install golang.org/x/tools/cmd/benchcmp@latest"; \
	else \
		echo "No baseline.txt found. Run 'make bench-baseline' first."; \
	fi

# Create baseline benchmark results
bench-baseline:
	@echo "Creating baseline benchmark results..."
	@go test -bench=. -benchtime=5s ./internal/... > baseline.txt
	@echo "Baseline saved to baseline.txt"

# Create current benchmark results  
bench-current:
	@echo "Running current benchmarks..."
	@go test -bench=. -benchtime=5s ./internal/... > current.txt

# Performance profiling
bench-profile:
	@echo "Running benchmarks with CPU and memory profiling..."
	@mkdir -p profiles
	@go test -bench=BenchmarkProcessDirectory_1000Files -cpuprofile=profiles/cpu.prof -memprofile=profiles/mem.prof -benchtime=10s ./internal/processor
	@echo "Profiles saved to profiles/"
	@echo "View with: go tool pprof profiles/cpu.prof"
	@echo "View with: go tool pprof profiles/mem.prof"

# Execution tracing
bench-trace:
	@echo "Running benchmarks with execution tracing..."
	@mkdir -p profiles
	@go test -bench=BenchmarkProcessDirectory_1000Files -trace=profiles/trace.out -benchtime=5s ./internal/processor
	@echo "Trace saved to profiles/trace.out"
	@echo "View with: go tool trace profiles/trace.out"

# Statistical analysis (requires benchstat)
bench-stats:
	@echo "Running statistical benchmark analysis..."
	@go test -bench=BenchmarkProcessDirectory_1000Files -count=10 ./internal/processor | benchstat || echo "benchstat not installed. Run: go install golang.org/x/perf/cmd/benchstat@latest"

# Memory leak detection
bench-leak-check:
	@echo "Running memory leak detection..."
	@mkdir -p profiles
	@go test -bench=BenchmarkProcessDirectory_MemoryProfile -benchtime=30s -memprofile=profiles/leak.prof ./internal/processor
	@echo "Check for memory leaks with: go tool pprof -alloc_space profiles/leak.prof"

# Performance regression check (for CI)
bench-regression: bench-current
	@echo "Checking for performance regressions..."
	@if [ -f baseline.txt ]; then \
		echo "Comparing with baseline..."; \
		benchcmp baseline.txt current.txt; \
		if [ $$? -gt 0 ]; then \
			echo "Performance regression detected!"; \
			exit 1; \
		fi; \
	else \
		echo "No baseline found - saving current results as baseline"; \
		cp current.txt baseline.txt; \
	fi

# Clean up benchmark artifacts
clean-bench:
	@echo "Cleaning benchmark artifacts..."
	@rm -f current.txt baseline.txt
	@rm -rf profiles/
	@rm -f *.prof *.out

# Install required tools
install-bench-tools:
	@echo "Installing benchmark analysis tools..."
	@go install golang.org/x/tools/cmd/benchcmp@latest
	@go install golang.org/x/perf/cmd/benchstat@latest
	@echo "Tools installed successfully"

# Validate benchmark installation
bench-validate:
	@echo "Validating benchmark setup..."
	@go test -run=^$$ -bench=BenchmarkProcessDirectory_100Files -benchtime=1s ./internal/processor > /dev/null
	@if [ $$? -eq 0 ]; then \
		echo "✓ Processor benchmarks working"; \
	else \
		echo "✗ Processor benchmarks failed"; \
		exit 1; \
	fi
	@go test -run=^$$ -bench=BenchmarkBinaryRule_Match_100Files -benchtime=1s ./internal/filter/rules > /dev/null
	@if [ $$? -eq 0 ]; then \
		echo "✓ Binary detection benchmarks working"; \
	else \
		echo "✗ Binary detection benchmarks failed"; \
		exit 1; \
	fi
	@go test -run=^$$ -bench=BenchmarkTokenCounter_SmallText -benchtime=1s ./internal/token > /dev/null
	@if [ $$? -eq 0 ]; then \
		echo "✓ Token counting benchmarks working"; \
	else \
		echo "✗ Token counting benchmarks failed"; \
		exit 1; \
	fi
	@go test -run=^$$ -bench=BenchmarkFilter_Creation_SimpleRules -benchtime=1s ./internal/filter > /dev/null
	@if [ $$? -eq 0 ]; then \
		echo "✓ Filter benchmarks working"; \
	else \
		echo "✗ Filter benchmarks failed"; \
		exit 1; \
	fi
	@echo "All benchmarks validated successfully!"

# Performance report generation
bench-report:
	@echo "Generating performance report..."
	@mkdir -p reports
	@echo "# Promptext Performance Report" > reports/performance_report.md
	@echo "Generated on: $$(date)" >> reports/performance_report.md
	@echo "" >> reports/performance_report.md
	@echo "## File Processing Performance" >> reports/performance_report.md
	@go test -bench=BenchmarkProcessDirectory -benchmem ./internal/processor | grep -E '^Benchmark|^PASS' >> reports/performance_report.md
	@echo "" >> reports/performance_report.md
	@echo "## Binary Detection Performance" >> reports/performance_report.md  
	@go test -bench=BenchmarkBinaryRule -benchmem ./internal/filter/rules | grep -E '^Benchmark|^PASS' >> reports/performance_report.md
	@echo "" >> reports/performance_report.md
	@echo "## Token Counting Performance" >> reports/performance_report.md
	@go test -bench=BenchmarkTokenCounter -benchmem ./internal/token | grep -E '^Benchmark|^PASS' >> reports/performance_report.md
	@echo "" >> reports/performance_report.md
	@echo "## Filter Performance" >> reports/performance_report.md
	@go test -bench=BenchmarkFilter -benchmem ./internal/filter | grep -E '^Benchmark|^PASS' >> reports/performance_report.md
	@echo "Performance report saved to reports/performance_report.md"

# Help target
help:
	@echo "Promptext Performance Testing Commands:"
	@echo ""
	@echo "Basic Benchmarks:"
	@echo "  bench              Run critical performance benchmarks (default)"
	@echo "  bench-all          Run comprehensive performance benchmarks"
	@echo "  bench-critical     Run fast benchmarks suitable for CI"
	@echo ""
	@echo "Component Benchmarks:"
	@echo "  bench-processor    File processing pipeline benchmarks"
	@echo "  bench-binary       Binary detection performance benchmarks"
	@echo "  bench-token        Token counting performance benchmarks"
	@echo "  bench-filter       Filter performance benchmarks"
	@echo "  bench-memory       Memory usage focused benchmarks"
	@echo ""
	@echo "Analysis & Comparison:"
	@echo "  bench-baseline     Create baseline benchmark results"
	@echo "  bench-compare      Compare current vs baseline results"
	@echo "  bench-regression   Check for performance regressions (CI)"
	@echo "  bench-stats        Statistical benchmark analysis"
	@echo ""
	@echo "Profiling:"
	@echo "  bench-profile      CPU and memory profiling"
	@echo "  bench-trace        Execution trace analysis"
	@echo "  bench-leak-check   Memory leak detection"
	@echo ""
	@echo "Utilities:"
	@echo "  bench-validate     Validate benchmark setup"
	@echo "  bench-report       Generate performance report"
	@echo "  install-bench-tools Install benchcmp and benchstat"
	@echo "  clean-bench        Clean benchmark artifacts"
	@echo ""
	@echo "Examples:"
	@echo "  make bench                    # Quick performance check"
	@echo "  make bench-all               # Comprehensive testing"
	@echo "  make bench-profile           # Profile for optimization"
	@echo "  make bench-baseline          # Set performance baseline"
	@echo "  make bench-compare           # Check for regressions"
	@echo ""
	@echo "For detailed documentation, see PERFORMANCE_BENCHMARKS.md"