package playground

import (
	"caminoclient/internal/creator"
	"caminoclient/internal/issuer"
	"caminoclient/internal/logger"

	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
)

const (
	localNode      = "http://127.0.0.1:19651"
	kopernikusNode = "https://kopernikus.camino.network"
	unchainedNode  = "https://kopernikus.unchained.camino.network"
	// columbusNode   = ""
	caminoNode = "https://api.camino.network"

	caminoFeeKeyStr = "\"PrivateKey-2bMrxpyN24b6BsiTjRDw3h7w7nC75ecZ5vkSyrDksxthvxXQ8o\""

	// 7Sdex3LTEjsnswW38Eb48hQ9insctGrsN : P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3
	localValidator0KeyStr = "\"PrivateKey-vmRQiZeXEXYMyJhEiqdC2z5JhuDbxL8ix9UVvjgMu2Er1NepE\""

	// 6Y3kysjF9jnHnYkdS9yGAuoHyae2eNmeV : P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68
	localValidator1KeyStr = "\"PrivateKey-ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN\""
)

var (
	caminoFeeKey       *secp256k1.PrivateKey
	localValidator0Key *secp256k1.PrivateKey
	localValidator1Key *secp256k1.PrivateKey
)

func init() {
	caminoFeeKey = new(secp256k1.PrivateKey)
	if err := caminoFeeKey.UnmarshalText([]byte(caminoFeeKeyStr)); err != nil {
		panic(err)
	}
	localValidator0Key = new(secp256k1.PrivateKey)
	if err := localValidator0Key.UnmarshalText([]byte(localValidator0KeyStr)); err != nil {
		panic(err)
	}
	localValidator1Key = new(secp256k1.PrivateKey)
	if err := localValidator1Key.UnmarshalText([]byte(localValidator1KeyStr)); err != nil {
		panic(err)
	}
}

// TxCreator

func LocalTxCreator(logger logger.Logger) *creator.TxCreator {
	txc, err := creator.NewTxCreator(
		localNode,
		localValidator0Key,
		logger,
	)
	logger.NoError(err)
	return txc
}

func LocalKopernikusTxCreator(logger logger.Logger) *creator.TxCreator {
	txc, err := creator.NewTxCreator(
		localNode,
		nil, // TODO@
		// "\"PrivateKey-2dkKrzcDQ2JKxJw1M3W5Dc5QB2QinhgCBZqeSswZXgsTEsytar\"",
		logger,
	)
	logger.NoError(err)
	return txc
}

func LocalCaminoTxCreator(logger logger.Logger) *creator.TxCreator {
	txc, err := creator.NewTxCreator(
		localNode,
		nil, // TODO@
		// "\"PrivateKey-2dkKrzcDQ2JKxJw1M3W5Dc5QB2QinhgCBZqeSswZXgsTEsytar\"",
		logger,
	)
	logger.NoError(err)
	return txc
}

func KopernikusTxCreator(logger logger.Logger) *creator.TxCreator {
	txc, err := creator.NewTxCreator(
		kopernikusNode,
		nil, // TODO@
		// "\"PrivateKey-2dkKrzcDQ2JKxJw1M3W5Dc5QB2QinhgCBZqeSswZXgsTEsytar\"",
		logger,
	)
	logger.NoError(err)
	return txc
}

func UnchainedTxCreator(logger logger.Logger) *creator.TxCreator {
	txc, err := creator.NewTxCreator(
		unchainedNode,
		localValidator0Key,
		// "\"PrivateKey-2dkKrzcDQ2JKxJw1M3W5Dc5QB2QinhgCBZqeSswZXgsTEsytar\"",
		logger,
	)
	logger.NoError(err)
	return txc
}

// TxIssuer

func UnchainedTxIssuer(logger logger.Logger) *issuer.TxIssuer {
	txi, err := issuer.NewTxIssuer(unchainedNode, logger)
	logger.NoError(err)
	return txi
}

func KopernikusTxIssuer(logger logger.Logger) *issuer.TxIssuer {
	txi, err := issuer.NewTxIssuer(kopernikusNode, logger)
	logger.NoError(err)
	return txi
}

func CaminoTxIssuer(logger logger.Logger) *issuer.TxIssuer {
	txi, err := issuer.NewTxIssuer(caminoNode, logger)
	logger.NoError(err)
	return txi
}

func LocalTxIssuer(logger logger.Logger) *issuer.TxIssuer {
	txi, err := issuer.NewTxIssuer(localNode, logger)
	logger.NoError(err)
	return txi
}
