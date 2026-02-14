package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ==============================================================================
// 1. BASIC ERROR CREATION
// ==============================================================================

// errors.New creates a simple error with a static string
func basicErrorCreation() {
	err := errors.New("this is a basic error")
	fmt.Printf("1. Basic error: %v\n", err)
	fmt.Printf("   Error type: %T\n\n", err)
}

// fmt.Errorf creates formatted errors
func formattedErrorCreation() {
	name := "Alice"
	age := 25
	err := fmt.Errorf("user %s with age %d is invalid", name, age)
	fmt.Printf("2. Formatted error: %v\n\n", err)
}

// ==============================================================================
// 2. SENTINEL ERRORS (predefined package-level errors)
// ==============================================================================

var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrInvalidInput  = errors.New("invalid input")
	ErrTimeout       = errors.New("operation timed out")
)

func sentinelErrors() {
	err := ErrNotFound

	// Check using == for sentinel errors
	if err == ErrNotFound {
		fmt.Printf("3. Sentinel error detected: %v\n", err)
	}

	// errors.Is is better for wrapped errors
	if errors.Is(err, ErrNotFound) {
		fmt.Printf("   errors.Is also works: true\n\n")
	}
}

// ==============================================================================
// 3. CUSTOM ERROR TYPES (implementing error interface)
// ==============================================================================

// Simple custom error type
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

// Implement error interface (requires Error() string method)
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s' (value: %v): %s",
		e.Field, e.Value, e.Message)
}

// Custom error with additional methods
type DatabaseError struct {
	Operation string
	Table     string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s on table %s: %v",
		e.Operation, e.Table, e.Err)
}

// Unwrap method allows errors.Is and errors.As to work with wrapped errors
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// Custom method on error type
func (e *DatabaseError) IsRetryable() bool {
	return e.Operation == "select" || e.Operation == "read"
}

func customErrorTypes() {
	// Create custom validation error
	err1 := &ValidationError{
		Field:   "email",
		Value:   "invalid-email",
		Message: "must be a valid email address",
	}
	fmt.Printf("4. Custom validation error: %v\n", err1)

	// Create custom database error
	err2 := &DatabaseError{
		Operation: "insert",
		Table:     "users",
		Err:       errors.New("duplicate key violation"),
	}
	fmt.Printf("   Custom database error: %v\n", err2)
	fmt.Printf("   Is retryable: %v\n\n", err2.IsRetryable())
}

// ==============================================================================
// 4. ERROR WRAPPING (adding context to errors)
// ==============================================================================

func innerOperation() error {
	return errors.New("low-level file system error")
}

func middleOperation() error {
	err := innerOperation()
	if err != nil {
		// %w wraps the error, preserving the error chain
		return fmt.Errorf("middle layer failed: %w", err)
	}
	return nil
}

func outerOperation() error {
	err := middleOperation()
	if err != nil {
		return fmt.Errorf("outer layer failed: %w", err)
	}
	return nil
}

func errorWrapping() {
	err := outerOperation()
	fmt.Printf("5. Wrapped error: %v\n", err)

	// Unwrap to get the underlying error
	unwrapped := errors.Unwrap(err)
	fmt.Printf("   First unwrap: %v\n", unwrapped)

	unwrapped = errors.Unwrap(unwrapped)
	fmt.Printf("   Second unwrap: %v\n\n", unwrapped)
}

// ==============================================================================
// 5. ERROR INSPECTION: errors.Is
// ==============================================================================

func errorsIsExample() {
	baseErr := ErrNotFound
	wrappedErr := fmt.Errorf("failed to get user: %w", baseErr)
	doubleWrappedErr := fmt.Errorf("API call failed: %w", wrappedErr)

	fmt.Printf("6. errors.Is - checking error identity through wrapping:\n")
	fmt.Printf("   errors.Is(doubleWrappedErr, ErrNotFound): %v\n",
		errors.Is(doubleWrappedErr, ErrNotFound))
	fmt.Printf("   Direct comparison (doubleWrappedErr == ErrNotFound): %v\n\n",
		doubleWrappedErr == ErrNotFound)
}

// ==============================================================================
// 6. ERROR INSPECTION: errors.As (type assertion for error chains)
// ==============================================================================

func processUser() error {
	return &ValidationError{
		Field:   "age",
		Value:   -5,
		Message: "must be positive",
	}
}

func handleUser() error {
	err := processUser()
	if err != nil {
		return fmt.Errorf("user processing failed: %w", err)
	}
	return nil
}

func errorsAsExample() {
	err := handleUser()

	fmt.Printf("7. errors.As - type assertion through wrapping:\n")

	// Extract specific error type from chain
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		fmt.Printf("   Found ValidationError in chain\n")
		fmt.Printf("   Field: %s, Value: %v, Message: %s\n\n",
			validationErr.Field, validationErr.Value, validationErr.Message)
	}
}

// ==============================================================================
// 7. MULTIPLE ERROR WRAPPING (Go 1.20+)
// ==============================================================================

func multipleErrorsJoin() {
	err1 := errors.New("error one")
	err2 := errors.New("error two")
	err3 := errors.New("error three")

	// Join multiple errors
	combined := errors.Join(err1, err2, err3)

	fmt.Printf("8. Multiple errors (errors.Join):\n")
	fmt.Printf("   Combined: %v\n", combined)

	// Check if any specific error is in the joined errors
	fmt.Printf("   Contains err1: %v\n", errors.Is(combined, err1))
	fmt.Printf("   Contains err2: %v\n\n", errors.Is(combined, err2))
}

// ==============================================================================
// 8. CUSTOM ERROR WITH MULTIPLE WRAPPED ERRORS
// ==============================================================================

type MultiError struct {
	Errors []error
}

func (m *MultiError) Error() string {
	return fmt.Sprintf("multiple errors occurred: %d errors", len(m.Errors))
}

// Unwrap returns multiple errors (Go 1.20+)
func (m *MultiError) Unwrap() []error {
	return m.Errors
}

func multiErrorCustomType() {
	multiErr := &MultiError{
		Errors: []error{
			errors.New("first problem"),
			errors.New("second problem"),
			ErrTimeout,
		},
	}

	fmt.Printf("9. Custom multi-error type: %v\n", multiErr)
	fmt.Printf("   Contains ErrTimeout: %v\n\n", errors.Is(multiErr, ErrTimeout))
}

// ==============================================================================
// 9. ERROR COMPARISONS AND NIL CHECKS
// ==============================================================================

func errorComparisons() {
	var err error

	fmt.Printf("10. Error comparisons:\n")
	fmt.Printf("    nil error == nil: %v\n", err == nil)

	err = errors.New("some error")
	fmt.Printf("    non-nil error == nil: %v\n", err == nil)
	fmt.Printf("    non-nil error != nil: %v\n\n", err != nil)
}

// ==============================================================================
// 10. PANIC AND RECOVER (error handling extreme cases)
// ==============================================================================

func panicExample() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("11. Recovered from panic: %v\n", r)

			// Convert panic to error
			err, ok := r.(error)
			if ok {
				fmt.Printf("    Panic was an error type: %v\n\n", err)
			} else {
				fmt.Printf("    Panic was not an error type: %v\n\n", r)
			}
		}
	}()

	panic(errors.New("something went terribly wrong"))
}

// ==============================================================================
// 11. ERROR WRAPPING WITHOUT %w (loses error chain)
// ==============================================================================

func wrappingWithoutW() {
	baseErr := ErrNotFound

	// Using %v instead of %w - does NOT preserve error chain
	wrappedWithV := fmt.Errorf("failed to find user: %v", baseErr)

	// Using %w - preserves error chain
	wrappedWithW := fmt.Errorf("failed to find user: %w", baseErr)

	fmt.Printf("12. Wrapping with %%v vs %%w:\n")
	fmt.Printf("    errors.Is(wrappedWithV, ErrNotFound): %v\n",
		errors.Is(wrappedWithV, ErrNotFound))
	fmt.Printf("    errors.Is(wrappedWithW, ErrNotFound): %v\n\n",
		errors.Is(wrappedWithW, ErrNotFound))
}

// ==============================================================================
// 13. PRACTICAL ERROR HANDLING PATTERNS
// ==============================================================================

func readFileExample(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		// Check for specific error types
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file does not exist: %w", err)
		}
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("permission denied: %w", err)
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return nil
}

func practicalPatterns() {
	err := readFileExample("/nonexistent/file.txt")

	fmt.Printf("13. Practical error handling:\n")
	if err != nil {
		fmt.Printf("    Error: %v\n", err)

		// Check specific error types
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("    Specific handling: file not found\n")
		}
	}
	fmt.Println()
}

// ==============================================================================
// 14. ERROR CHAIN INSPECTION
// ==============================================================================

func inspectErrorChain(err error) {
	fmt.Printf("14. Error chain inspection:\n")
	fmt.Printf("    Full error: %v\n", err)

	current := err
	depth := 0
	for current != nil {
		fmt.Printf("    [%d] %v (type: %T)\n", depth, current, current)
		current = errors.Unwrap(current)
		depth++
	}
	fmt.Println()
}

// ==============================================================================
// 15. ERRORS WITH ADDITIONAL CONTEXT
// ==============================================================================

type HTTPError struct {
	StatusCode int
	Method     string
	URL        string
	Err        error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d error on %s %s: %v",
		e.StatusCode, e.Method, e.URL, e.Err)
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

func httpErrorExample() {
	err := &HTTPError{
		StatusCode: 404,
		Method:     "GET",
		URL:        "/api/users/123",
		Err:        ErrNotFound,
	}

	fmt.Printf("15. HTTP error with context: %v\n", err)
	fmt.Printf("    Contains ErrNotFound: %v\n\n", errors.Is(err, ErrNotFound))
}

// ==============================================================================
// 16. ERROR CREATION FROM STRING CONVERSION
// ==============================================================================

func errorFromConversion() {
	_, err := strconv.Atoi("not-a-number")

	fmt.Printf("16. Error from standard library (strconv):\n")
	fmt.Printf("    Error: %v\n", err)
	fmt.Printf("    Type: %T\n", err)

	// Extract specific error type
	var numError *strconv.NumError
	if errors.As(err, &numError) {
		fmt.Printf("    Function: %s\n", numError.Func)
		fmt.Printf("    Invalid value: %s\n", numError.Num)
		fmt.Printf("    Underlying error: %v\n\n", numError.Err)
	}
}

// ==============================================================================
// 17. IMPLEMENTING io.EOF AND SPECIAL ERRORS
// ==============================================================================

func eofHandling() {
	err := io.EOF

	fmt.Printf("17. Special errors (io.EOF):\n")
	fmt.Printf("    EOF error: %v\n", err)
	fmt.Printf("    Is EOF: %v\n", errors.Is(err, io.EOF))
	fmt.Printf("    EOF is often not an error but an expected condition\n\n")
}

// ==============================================================================
// 18. ERROR WRAPPING CHAINS WITH MULTIPLE TYPES
// ==============================================================================

func complexErrorChain() {
	// Create a complex error chain
	baseErr := errors.New("network timeout")
	dbErr := &DatabaseError{
		Operation: "select",
		Table:     "users",
		Err:       baseErr,
	}
	httpErr := &HTTPError{
		StatusCode: 500,
		Method:     "POST",
		URL:        "/api/data",
		Err:        dbErr,
	}
	topErr := fmt.Errorf("request failed: %w", httpErr)

	fmt.Printf("18. Complex error chain:\n")
	inspectErrorChain(topErr)

	// Check if we can find specific types
	var dbError *DatabaseError
	var httpError *HTTPError
	fmt.Printf("    Contains DatabaseError: %v\n", errors.As(topErr, &dbError))
	fmt.Printf("    Contains HTTPError: %v\n", errors.As(topErr, &httpError))
	if dbError != nil {
		fmt.Printf("    Database operation: %s\n", dbError.Operation)
		fmt.Printf("    Is retryable: %v\n\n", dbError.IsRetryable())
	}
}

// ==============================================================================
// 19. DEFER AND ERROR HANDLING
// ==============================================================================

func deferWithErrors() (err error) {
	defer func() {
		if err != nil {
			// Modify error before returning
			err = fmt.Errorf("deferWithErrors: %w", err)
		}
	}()

	return errors.New("original error")
}

func deferErrorExample() {
	err := deferWithErrors()
	fmt.Printf("19. Defer with error modification: %v\n\n", err)
}

// ==============================================================================
// 20. ERROR AGGREGATION PATTERN
// ==============================================================================

func collectErrors() error {
	var errs []error

	// Simulate multiple operations
	errs = append(errs, errors.New("failed to validate input"))
	errs = append(errs, errors.New("failed to connect to database"))
	errs = append(errs, nil) // Some operations succeed
	errs = append(errs, errors.New("failed to send notification"))

	// Filter out nil errors
	var actualErrors []error
	for _, err := range errs {
		if err != nil {
			actualErrors = append(actualErrors, err)
		}
	}

	if len(actualErrors) == 0 {
		return nil
	}

	return errors.Join(actualErrors...)
}

func errorAggregation() {
	err := collectErrors()
	fmt.Printf("20. Error aggregation:\n")
	fmt.Printf("    Aggregated errors: %v\n\n", err)
}

// ==============================================================================
// MAIN FUNCTION - Demonstrates all error features
// ==============================================================================

func main() {
	fmt.Println("=== Go Error Interface - Complete Reference ===\n")

	basicErrorCreation()
	formattedErrorCreation()
	sentinelErrors()
	customErrorTypes()
	errorWrapping()
	errorsIsExample()
	errorsAsExample()
	multipleErrorsJoin()
	multiErrorCustomType()
	errorComparisons()
	panicExample()
	wrappingWithoutW()
	practicalPatterns()

	// Create an error chain for inspection
	baseErr := errors.New("base error")
	wrappedOnce := fmt.Errorf("level 1: %w", baseErr)
	wrappedTwice := fmt.Errorf("level 2: %w", wrappedOnce)
	inspectErrorChain(wrappedTwice)

	httpErrorExample()
	errorFromConversion()
	eofHandling()
	complexErrorChain()
	deferErrorExample()
	errorAggregation()

	fmt.Println("=== Summary of Key Concepts ===")
	fmt.Println("1. Basic creation: errors.New(), fmt.Errorf()")
	fmt.Println("2. Sentinel errors: package-level var errors")
	fmt.Println("3. Custom types: implement Error() string method")
	fmt.Println("4. Wrapping: use %w in fmt.Errorf() to preserve chain")
	fmt.Println("5. Unwrapping: implement Unwrap() error method")
	fmt.Println("6. Inspection: errors.Is() for identity, errors.As() for types")
	fmt.Println("7. Multiple errors: errors.Join() (Go 1.20+)")
	fmt.Println("8. Panic/Recover: extreme error handling")
	fmt.Println("9. Always check: if err != nil")
	fmt.Println("10. Add context: wrap errors as they bubble up")
}
