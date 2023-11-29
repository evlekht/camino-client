package node

import (
	"caminoclient/internal/logger"
	"caminoclient/internal/node_client"
	"caminoclient/internal/utils"
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
)

func NewClient(nodeURI string, logger logger.Logger) (*Client, error) {
	client := node_client.NewClient(nodeURI, logger)

	nodeCfg, err := client.P.GetConfiguration(context.Background())
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// feeKeyAddrStr, err := address.Format("P", hrp, feeKey.Address().Bytes())
	// if err != nil {
	// 	logger.Error(err)
	// 	return nil, err
	// }
	// logger.Infof("feeKey: %s : %s : %s\n\n", feeKey.String(), feeKey.Address().String(), feeKeyAddrStr)

	return &Client{
		client: client,
		logger: logger,
		utils:  utils.NewUtils(logger),

		networkID: uint32(nodeCfg.NetworkID),
		hrp:       constants.GetHRP(uint32(nodeCfg.NetworkID)),
	}, nil
}

type Client struct {
	client node_client.Client
	logger logger.Logger
	utils  *utils.UtilsWithLogger

	networkID uint32
	hrp       string
}

func (c *Client) GetPTX(txID ids.ID) (*txs.Tx, error) {
	c.logger.Info("Getting P-Chain tx...")
	txBytes, err := c.client.P.GetTx(context.Background(), txID)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	return txs.Parse(txs.Codec, txBytes)
}
