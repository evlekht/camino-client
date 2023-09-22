package playground

import (
	"caminoclient/internal/config"
	"caminoclient/internal/logger"
	"caminoclient/internal/utils"
	"context"
)

func NewPlayground(ctx context.Context, logger logger.Logger, cfg *config.Config) (*Playground, error) {
	return &Playground{
		logger: logger,
		utils:  utils.NewUtils(logger),
	}, nil
}

type Playground struct {
	logger logger.Logger
	utils  *utils.UtilsWithLogger
}

func (p *Playground) Run(ctx context.Context) error {
	// creator := LocalTxCreator(p.logger)
	// issuer := LocalTxIssuer(p.logger)

	return nil
}

func (p *Playground) Close(ctx context.Context) error {
	return nil
}
