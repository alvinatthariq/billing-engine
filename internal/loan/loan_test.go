package loan

import (
	"testing"
	"time"
)

func TestNewLoan(t *testing.T) {
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	amount := 5000000.0
	interestRate := 0.10

	loan, err := NewLoan("100", amount, interestRate, startDate)
	if err != nil {
		t.Errorf("Failed to create new loan: %v", err)
	}

	// Test total repayment calculation
	expectedTotalRepayment := amount + (amount * interestRate)
	if loan.TotalRepayment != expectedTotalRepayment {
		t.Errorf("Expected total repayment %f, got %f", expectedTotalRepayment, loan.TotalRepayment)
	}

	// Test weekly payment calculation
	expectedWeeklyPayment := expectedTotalRepayment / float64(TotalWeeks)
	if loan.WeeklyPayment != expectedWeeklyPayment {
		t.Errorf("Expected weekly payment %f, got %f", expectedWeeklyPayment, loan.WeeklyPayment)
	}
}

func TestMakePayment(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		paymentDate time.Time
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "Single payment success",
			amount:      110000, // Weekly payment amount
			paymentDate: time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
			wantErr:     true,
			errMsg:      "payment amount 110000 does not match number of pending payments (2 weeks pending)",
		},
		{
			name:        "Multiple weeks payment success",
			amount:      220000, // Two weeks payment
			paymentDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:     true,
			errMsg:      "payment amount 220000 does not match number of pending payments (3 weeks pending)",
		},
		{
			name:        "Payment less than pending weeks",
			amount:      110000, // One week when three are pending
			paymentDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:     true,
			errMsg:      "payment amount 110000 does not match number of pending payments (3 weeks pending)",
		},
		{
			name:        "Invalid payment date",
			amount:      110000,
			paymentDate: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			wantErr:     true,
			errMsg:      "invalid payment date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new loan for each test case
			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			loan, _ := NewLoan("100", 5000000.0, 0.10, startDate)

			err := loan.MakePayment(tt.amount, tt.paymentDate)

			if tt.wantErr {
				if err == nil {
					t.Errorf("MakePayment() error = nil, want error %v", tt.errMsg)
				} else if err.Error() != tt.errMsg {
					t.Errorf("MakePayment() error = %v, want %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("MakePayment() unexpected error = %v", err)
				}

				// Verify payments were recorded correctly
				expectedPaid := int(tt.amount / loan.WeeklyPayment)
				paidCount := 0
				for i := 0; i < expectedPaid; i++ {
					if loan.Payments[i].Paid {
						paidCount++
					}
				}
				if paidCount != expectedPaid {
					t.Errorf("Expected %d payments recorded, got %d", expectedPaid, paidCount)
				}

				// Verify outstanding amount
				expectedOutstanding := loan.TotalRepayment - tt.amount
				if loan.GetOutstanding() != expectedOutstanding {
					t.Errorf("Expected outstanding %f, got %f", expectedOutstanding, loan.GetOutstanding())
				}
			}
		})
	}
}

func TestIsDelinquent(t *testing.T) {
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	loan, _ := NewLoan("100", 5000000.0, 0.10, startDate)

	// Not delinquent at start
	if loan.IsDelinquent(startDate) {
		t.Error("Loan should not be delinquent at start")
	}

	// Move forward 3 weeks without payments
	currentTime := startDate.Add(24 * time.Hour * 7 * 3)
	if !loan.IsDelinquent(currentTime) {
		t.Error("Loan should be delinquent after 3 weeks of no payments")
	}

	// Make payment for all pending weeks
	err := loan.MakePayment(loan.WeeklyPayment*4, currentTime)
	if err != nil {
		t.Errorf("Failed to make payment: %v", err)
	}

	// Should no longer be delinquent after paying all pending weeks
	if loan.IsDelinquent(currentTime) {
		t.Error("Loan should not be delinquent after paying all pending weeks")
	}
}
