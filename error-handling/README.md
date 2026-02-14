# Go Error Interface - Complete Reference

A comprehensive, single-file reference demonstrating all features of Go's error interface.

## What's Included

This `main.go` file demonstrates:

### Basic Error Handling
1. **Basic Error Creation** - `errors.New()` and `fmt.Errorf()`
2. **Sentinel Errors** - Package-level predefined errors
3. **Custom Error Types** - Implementing the `error` interface
4. **Error Wrapping** - Using `%w` to preserve error chains

### Error Inspection
5. **errors.Is()** - Check error identity through wrapping
6. **errors.As()** - Type assertion for error chains
7. **Error Unwrapping** - Implementing `Unwrap()` method
8. **Error Chain Inspection** - Walking through error chains

### Advanced Features
9. **Multiple Error Handling** - `errors.Join()` (Go 1.20+)
10. **Custom Multi-Error Types** - Custom types with multiple wrapped errors
11. **Panic and Recover** - Extreme error handling
12. **Defer with Errors** - Modifying errors in defer blocks

### Practical Patterns
13. **HTTP Errors** - Errors with status codes and context
14. **Database Errors** - Custom methods on error types
15. **Validation Errors** - Field-specific error information
16. **File Operation Errors** - Handling `os.ErrNotExist`, `os.ErrPermission`
17. **Conversion Errors** - Working with `strconv.NumError`
18. **EOF Handling** - Special error conditions

### Key Concepts
- **%v vs %w**: Using `%v` loses the error chain, `%w` preserves it
- **nil checks**: Always check `if err != nil`
- **Error context**: Add context as errors bubble up
- **Type assertions**: Use `errors.As()` to extract specific error types
- **Identity checks**: Use `errors.Is()` to check for specific errors

## Running the Code

```bash
cd error-handling
go run main.go
```

## Usage as Reference

This file is designed to be:
- **Self-contained** - Everything in one file for easy reference
- **Well-commented** - Each section explains what it demonstrates
- **Runnable** - Execute it to see all concepts in action
- **Copy-paste friendly** - Use sections as templates for your code

## Go Version

Requires Go 1.20+ for `errors.Join()` functionality.
