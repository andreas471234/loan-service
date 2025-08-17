# Tools

This directory contains development tools and generated files.

## Generated Files

- **coverage.out** - Go test coverage data
- **coverage.html** - HTML coverage report

## Usage

```bash
# Generate coverage report
make coverage

# View coverage report
open tools/coverage.html

# Clean up generated files
rm tools/coverage.*
```

## Notes

- Generated during development and testing
- Not committed to version control (see .gitignore)
- Coverage reports help identify untested code areas
