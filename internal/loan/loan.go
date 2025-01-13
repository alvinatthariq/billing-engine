package loan

import (
	"errors"
	"fmt"
	"math"
	"time"
)

const (
	WeeksInYear = 52
	TotalWeeks  = 50
)

type PaymentStatus struct {
	Week     int
	Amount   float64
	Paid     bool
	PaidDate *time.Time
}

type Loan struct {
	ID             string
	Amount         float64
	InterestRate   float64
	WeeklyPayment  float64
	StartDate      time.Time
	Payments       []PaymentStatus
	Outstanding    float64
	TotalRepayment float64
}

// NewLoan creates a new loan with the given parameters
func NewLoan(id string, amount float64, interestRate float64, startDate time.Time) (*Loan, error) {
	if amount <= 0 {
		return nil, errors.New("loan amount must be positive")
	}
	if interestRate <= 0 {
		return nil, errors.New("interest rate must be positive")
	}

	// Calculate total repayment with interest
	interestAmount := amount * interestRate
	totalRepayment := amount + interestAmount

	// Calculate weekly payment
	weeklyPayment := totalRepayment / float64(TotalWeeks)

	// Initialize payment schedule
	payments := make([]PaymentStatus, TotalWeeks)
	for i := 0; i < TotalWeeks; i++ {
		payments[i] = PaymentStatus{
			Week:     i + 1,
			Amount:   weeklyPayment,
			Paid:     false,
			PaidDate: nil,
		}
	}

	return &Loan{
		ID:             id,
		Amount:         amount,
		InterestRate:   interestRate,
		WeeklyPayment:  weeklyPayment,
		StartDate:      startDate,
		Payments:       payments,
		Outstanding:    totalRepayment,
		TotalRepayment: totalRepayment,
	}, nil
}

// GetOutstanding returns the current outstanding amount on the loan
func (l *Loan) GetOutstanding() float64 {
	return l.Outstanding
}

// IsDelinquent checks if the borrower has missed more than 2 consecutive payments
func (l *Loan) IsDelinquent(currentTime time.Time) bool {
	// currentTime = currentTime.Add(24 * time.Hour * 120)
	currentWeek := int(currentTime.Sub(l.StartDate).Hours() / (24 * 7))
	if currentWeek < 2 {
		return false
	}

	// Check last two due payments
	unpaidCount := 0
	for i := currentWeek - 1; i >= 0 && i >= currentWeek-2; i-- {
		if i >= len(l.Payments) {
			continue
		}
		if !l.Payments[i].Paid {
			unpaidCount++
		}
	}

	return unpaidCount >= 2
}

// MakePayment processes a payment for the loan
func (l *Loan) MakePayment(amount float64, paymentDate time.Time) error {
	if paymentDate.Before(l.StartDate) {
		return errors.New("invalid payment date")
	}

	currentWeek := int(paymentDate.Sub(l.StartDate).Hours() / (24 * 7))
	if currentWeek >= len(l.Payments) {
		return errors.New("invalid payment date")
	}

	// Find all unpaid weeks up to current week
	var unpaidWeeks []int
	for i := 0; i <= currentWeek; i++ {
		if !l.Payments[i].Paid {
			unpaidWeeks = append(unpaidWeeks, i)
		}
	}

	if len(unpaidWeeks) == 0 {
		return errors.New("no pending payments found")
	}

	// Check if payment amount matches exactly one or more weekly payments
	if math.Mod(amount, l.WeeklyPayment) != 0 {
		return errors.New("payment amount must be a multiple of the weekly payment amount")
	}

	numPayments := int(amount / l.WeeklyPayment)
	if numPayments != len(unpaidWeeks) {
		return fmt.Errorf("payment amount %v does not match number of pending payments (%v weeks pending)", amount, len(unpaidWeeks))
	}

	// Process the payments starting from earliest unpaid week
	for i := 0; i < numPayments; i++ {
		weekIndex := unpaidWeeks[i]
		l.Payments[weekIndex].Paid = true
		l.Payments[weekIndex].PaidDate = &paymentDate
	}
	l.Outstanding -= amount

	return nil
}
