package node

import (
	"context"
)

func (c *Client) IssuePTx(txBytes []byte) error {
	c.logger.Info("Issuing P-Chain tx...")
	txID, err := c.client.P.IssueTx(context.Background(), txBytes)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	c.logger.Infof("\ntx %s issued!\n\n", txID)
	return nil
}
