// Package constants holds small shared types used across packages (e.g. context keys).
package constants

// ContextKey avoids collisions when storing values in context.Context.
type ContextKey string

const (
	OpNormalPurchase      = 1
	OpInstallmentPurchase = 2
	OpWithdrawal          = 3
	OpCreditVoucher       = 4
)

var OperationNames = map[int]string{
	OpNormalPurchase:      "Normal Purchase",
	OpInstallmentPurchase: "Purchase with Installments",
	OpWithdrawal:          "Withdrawal",
	OpCreditVoucher:       "Credit Voucher",
}

var OperationSign = map[int]int{
	OpNormalPurchase:      -1,
	OpWithdrawal:          -1,
	OpInstallmentPurchase: -1,
	OpCreditVoucher:       1,
}
