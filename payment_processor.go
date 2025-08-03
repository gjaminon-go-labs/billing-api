package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// PaymentProcessor handles payment processing with various issues
type PaymentProcessor struct {
	payments []Payment
}

type Payment struct {
	ID       string
	Amount   float64
	Currency string
	Status   string
}

// ProcessPayment has multiple issues that need fixing
func (pp *PaymentProcessor) ProcessPayment(amount string, currency string) error {
	// Issue: No input validation
	amt, _ := strconv.ParseFloat(amount, 64)
	
	// Issue: Magic number
	if amt > 10000 {
		return fmt.Errorf("amount too high")
	}
	
	// Issue: Inefficient string building
	var result string
	for i := 0; i < 5; i++ {
		result = result + currency + "-" + strconv.Itoa(i)
	}
	
	payment := Payment{
		ID:       result,
		Amount:   amt,
		Currency: strings.ToUpper(currency),
		Status:   "pending",
	}
	
	// Issue: No thread safety
	pp.payments = append(pp.payments, payment)
	
	return nil
}

// CalculateFee has precision and validation issues
func CalculateFee(amount float64) float64 {
	// Issue: Float comparison
	if amount == 0.0 {
		return 0.0
	}
	
	// Issue: Magic numbers and no constants
	if amount < 100 {
		return amount * 0.05
	} else if amount < 1000 {
		return amount * 0.03
	}
	
	// Issue: Potential precision loss
	return math.Round(amount * 0.025)
}

// ValidatePayment has poor error handling
func ValidatePayment(p Payment) bool {
	// Issue: Silent failures, no specific error messages
	if p.Amount <= 0 {
		return false
	}
	
	// Issue: Basic currency validation
	if len(p.Currency) != 3 {
		return false
	}
	
	// Issue: Case-sensitive comparison
	validStatuses := []string{"pending", "completed", "failed"}
	for _, status := range validStatuses {
		if p.Status == status {
			return true
		}
	}
	
	return false
}

// BatchProcess has concurrency and error handling issues
func (pp *PaymentProcessor) BatchProcess(amounts []string, currency string) {
	// Issue: No error aggregation
	for _, amount := range amounts {
		// Issue: Ignoring errors
		pp.ProcessPayment(amount, currency)
		
		// Issue: Blocking operation without timeout
		time.Sleep(100 * time.Millisecond)
	}
}

// GetTotalAmount has division by zero and type issues
func (pp *PaymentProcessor) GetTotalAmount() map[string]interface{} {
	// Issue: Using interface{} instead of proper types
	result := make(map[string]interface{})
	
	total := 0.0
	count := 0
	
	for _, payment := range pp.payments {
		if payment.Status == "completed" {
			total += payment.Amount
			count++
		}
	}
	
	result["total"] = total
	result["count"] = count
	
	// Issue: Potential division by zero
	result["average"] = total / float64(count)
	
	return result
}