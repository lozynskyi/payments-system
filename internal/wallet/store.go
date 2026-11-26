// Package wallet is the wallet service: account, balance, and the ledger
package wallet

import (
	"github.com/lozynskyi/payments-system/internal/wallet/walletdb"
)

// The generated models doubles our domain types for now.The aliases mean
// callers only ever import package wallet newer walletdb
type (
	Account        = walletdb.Account
	AccountBalance = walletdb.AccountBalance
	Entry          = walletdb.Entry
	NewAccount     = walletdb.CreateAccountParams
	Direction      = walletdb.EntryDirection
)
