package types

import (
	"github.com/rivine/rivine/types"
)

// RegisterTransactionTypesForTestNetwork registers the transaction controllers
// for all transaction versions supported on the test network.
func RegisterTransactionTypesForTestNetwork() {
	// Explicitly remove version 0, as it was deprecated in rivine before HEChain was created
	types.RegisterTransactionVersion(types.TransactionVersionZero, nil)

}

// RegisterTransactionTypesForDevNetwork registers he transaction controllers
// for all transaction versions supported on the dev network.
func RegisterTransactionTypesForDevNetwork() {
	// Explicitly remove version 0, as it was deprecated in rivine before HEChain was created
	types.RegisterTransactionVersion(types.TransactionVersionZero, nil)
}
