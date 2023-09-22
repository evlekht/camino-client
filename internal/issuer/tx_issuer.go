package issuer

import (
	"caminoclient/internal/logger"
	"caminoclient/internal/node"
	"context"
)

func NewTxIssuer(nodeURI string, logger logger.Logger) (*TxIssuer, error) {
	return &TxIssuer{
		client: node.NewClient(nodeURI, logger),
		logger: logger,
	}, nil
}

type TxIssuer struct {
	client node.Client
	logger logger.Logger
}

func (ti *TxIssuer) IssuePTx(txBytes []byte) error {
	ti.logger.Info("Issuing P-Chain tx...")
	txID, err := ti.client.P.IssueTx(context.Background(), txBytes)
	if err != nil {
		ti.logger.Error(err)
		return err
	}
	ti.logger.Infof("\ntx %s issued!\n\n", txID)
	return nil
}
