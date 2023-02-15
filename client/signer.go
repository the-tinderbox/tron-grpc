package client

import (
	"github.com/the-tinderbox/tron-grpc/address"
	"github.com/the-tinderbox/tron-grpc/core"
)

type Signer interface {
	Address() address.Address
	PublicKey() []byte
	SignTransaction(tx *core.Transaction) error
	SignMessage(msg string) ([]byte, error)
}
