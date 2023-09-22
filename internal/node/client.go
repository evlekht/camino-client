package node

import (
	"caminoclient/internal/logger"
	"context"

	"github.com/ava-labs/avalanchego/api"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/json"
	"github.com/ava-labs/avalanchego/utils/rpc"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	platformapi "github.com/ava-labs/avalanchego/vms/platformvm/api"
	"github.com/ava-labs/avalanchego/vms/platformvm/locked"
	pTxs "github.com/ava-labs/avalanchego/vms/platformvm/txs"
)

// newClient returns a Client for interacting with the P Chain endpoint
func NewClient(uri string, logger logger.Logger) Client {
	return Client{
		P:          platformvm.NewClient(uri),
		pRequester: rpc.NewEndpointRequester(uri + "/ext/P"),
		logger:     logger,
	}
}

// Client implementation for interacting with the P Chain endpoint
type Client struct {
	P          platformvm.Client
	pRequester rpc.EndpointRequester
	logger     logger.Logger
}

// TODO update caminogo p-spend
// TODO get spendT from app-service
// TODO probably remove package after that
func (c *Client) Spend(
	ctx context.Context,
	networkID uint32,
	from ids.ShortID,
	to ids.ShortID,
	amountToLock uint64,
	amountToBurn uint64,
	lockMode locked.State,
	options ...rpc.Option,
) ([]*avax.TransferableInput, []*avax.TransferableOutput, error) {
	fromAddr, err := address.Format("P", constants.GetHRP(networkID), from[:])
	if err != nil {
		c.logger.Error(err)
		return nil, nil, err
	}
	toAddr, err := address.Format("P", constants.GetHRP(networkID), to[:])
	if err != nil {
		c.logger.Error(err)
		return nil, nil, err
	}

	type Spend2Reply struct {
		BaseTx string `json:"baseTx"`
	}
	res := &Spend2Reply{}
	if err := c.pRequester.SendRequest(ctx, "platform.spend2", &platformvm.SpendArgs{
		JSONFromAddrs: api.JSONFromAddrs{
			From: []string{fromAddr},
		},
		To: platformapi.Owner{
			Threshold: 1,
			Addresses: []string{toAddr},
		},
		AmountToLock: json.Uint64(amountToLock),
		AmountToBurn: json.Uint64(amountToBurn),
		LockMode:     byte(lockMode),
		Encoding:     formatting.Hex,
	}, res, options...); err != nil {
		c.logger.Error(err)
		return nil, nil, err
	}

	baseTxBytes, err := formatting.Decode(formatting.Hex, res.BaseTx)
	if err != nil {
		c.logger.Error(err)
		return nil, nil, err
	}
	baseTx := &pTxs.BaseTx{}
	if _, err := pTxs.Codec.Unmarshal(baseTxBytes, baseTx); err != nil {
		c.logger.Error(err)
		return nil, nil, err
	}

	return baseTx.Ins, baseTx.Outs, err
}
