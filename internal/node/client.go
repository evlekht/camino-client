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

	tChainID := ids.ID{}
	for _, blockchain := range nodeCfg.Blockchains {
		if blockchain.Name == "T-Chain" {
			tChainID = blockchain.ID
			break
		}
	}

	return &Client{
		client: client,
		logger: logger,
		utils:  utils.NewUtils(logger, uint32(nodeCfg.NetworkID)),

		tChainID:  tChainID,
		networkID: uint32(nodeCfg.NetworkID),
		hrp:       constants.GetHRP(uint32(nodeCfg.NetworkID)),
	}, nil
}

type Client struct {
	client node_client.Client
	logger logger.Logger
	utils  *utils.UtilsWithLogger

	tChainID  ids.ID
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
