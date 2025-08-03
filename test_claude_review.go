package main

import (
	"fmt"
	"strings"
)

// TestCodeForReview demonstrates some code patterns for Claude to review
// This file contains intentional issues to test Claude's review capabilities
func TestCodeForReview(input string) string {
	// Issues for Claude to potentially identify:
	// 1. No input validation
	// 2. Inefficient string concatenation
	// 3. Missing error handling
	// 4. Unused variable
	
	var result string
	unused := "this variable is never used"
	_ = unused // silence linter temporarily
	
	// Inefficient string concatenation in loop
	for i := 0; i < 5; i++ {
		result = result + input + fmt.Sprintf("_%d", i)
	}
	
	// Unnecessary string conversion
	result = strings.ToUpper(result)
	
	return result
}

// ProcessData has potential improvements
func ProcessData(data []string) {
	// Missing nil check
	// No return value or error handling
	for i := 0; i < len(data); i++ {
		fmt.Println(data[i]) // Direct printing instead of logging
	}
}

// UserInfo represents user data
type UserInfo struct {
	Name  string // No validation
	Email string // No email format validation
	Age   int    // No range validation
}

// SaveUser saves user without proper validation or error handling
func SaveUser(user UserInfo) {
	// Missing validation
	// No database error handling
	// Direct printing instead of structured logging
	fmt.Printf("Saving user: %s (%s)\n", user.Name, user.Email)
}