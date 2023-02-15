package trx

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/the-tinderbox/tron-grpc/address"
	"github.com/the-tinderbox/tron-grpc/api"
	"github.com/the-tinderbox/tron-grpc/client"
	"github.com/the-tinderbox/tron-grpc/core"
	"github.com/the-tinderbox/tron-grpc/tx"
)

type Client struct {
	client *client.Client
}

func New(client *client.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) getSignerAddress() address.Address {
	return c.client.Signer.Address()
}

func (c *Client) newTxAndSend(ctx context.Context, tx_ *core.Transaction) (*tx.Transaction, error) {
	t := tx.New(c.client, tx_)
	return t, t.Send(ctx, c.client.Signer)
}

func (c *Client) GetAccount(ctx context.Context, account string) (*core.Account, error) {
	addr, err := address.FromBase58(account)
	in := core.Account{Address: addr}
	ret, err := c.client.GetAccount(ctx, &in)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) GetBalance(ctx context.Context, account string) (int64, error) {
	acc, err := c.GetAccount(ctx, account)
	if err != nil {
		return 0, err
	}
	return acc.Balance, nil
}

func (c *Client) GetAccountResource(ctx context.Context, addr string) (*api.AccountResourceMessage, error) {
	addr_, err := address.FromBase58(addr)
	if err != nil {
		return nil, err
	}
	account := &core.Account{
		Address: addr_,
	}

	return c.client.GetAccountResource(ctx, account)
}

func (c *Client) CreateAccount(ctx context.Context, account string) (*tx.Transaction, error) {
	toAddr, err := address.FromBase58(account)
	if err != nil {
		return nil, err
	}

	contract := &core.AccountCreateContract{
		OwnerAddress:   c.getSignerAddress(),
		AccountAddress: toAddr,
	}

	tx_, err := c.client.CreateAccount2(ctx, contract)
	if err != nil {
		return nil, err
	}
	if proto.Size(tx_) == 0 {
		return nil, fmt.Errorf("bad transaction")
	}
	if tx_.GetResult().GetCode() != 0 {
		return nil, fmt.Errorf("%s", tx_.GetResult().GetMessage())
	}
	return c.newTxAndSend(ctx, tx_.Transaction)
}

func (c *Client) Transfer(ctx context.Context, to string, amount int64) (*tx.Transaction, error) {
	toAddr, err := address.FromBase58(to)
	if err != nil {
		return nil, err
	}

	contract := &core.TransferContract{
		OwnerAddress: c.getSignerAddress(),
		ToAddress:    toAddr,
		Amount:       amount,
	}

	tx_, err := c.client.CreateTransaction2(ctx, contract)
	if err != nil {
		return nil, err
	}

	return c.newTxAndSend(ctx, tx_.Transaction)
}
