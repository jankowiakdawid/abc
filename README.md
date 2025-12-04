# ABC Metrics Tool

A command-line tool for measuring Assignment, Branch, Condition (ABC) metrics in code.

## Overview

ABC Metrics is a code complexity measurement that quantifies a method's complexity based on:

- **A**: Number of assignments
- **B**: Number of branches (function calls)
- **C**: Number of conditions (if/else, loops, etc.)

The ABC score is calculated as `sqrt(A² + B² + C²)`.

## Understanding ABC Metrics

ABC Metrics provide a quantitative measure of code complexity:

- **Low** (< 10): Simple, well-structured code
- **Medium** (10-20): Moderately complex code
- **High** (20-40): Complex code that may need refactoring
- **Very High** (> 40): Overly complex code that should be refactored

## Installation

```bash
# Clone the repository
git clone https://github.com/your-username/abc-metrics.git

# Navigate to the project directory
cd abc-metrics

# Build the project
go build -o abc ./cmd/abc
```

## Usage

```bash
# Analyze a single file
./abc analyze -f path/to/your/file.go

# Alternative syntax
./abc analyze path/to/your/file.go

# Enable verbose output
./abc analyze -f path/to/your/file.go -v

# Show detailed breakdown of metrics
./abc analyze -f path/to/your/file.go --show
```

## Supported Languages

Currently, the tool supports:

- Go (`.go` files)
- TypeScript (coming soon)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
