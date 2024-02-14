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

// func (ti *TxIssuer) IssueTTx(txBytes []byte) error {
// 	ti.logger.Info("Issuing T-Chain tx...")
// 	txID, err := ti.client.T.IssueTx(context.Background(), txBytes)
// 	if err != nil {
// 		ti.logger.Error(err)
// 		return err
// 	}
// 	ti.logger.Infof("\ntx %s issued!\n\n", txID)
// 	return nil
// }
