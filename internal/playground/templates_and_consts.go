package playground

import (
	"caminoclient/internal/logger"
	"caminoclient/internal/matrix"
	"caminoclient/internal/node"

	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
)

const (
	localNode      = "http://127.0.0.1:19651"
	kopernikusNode = "https://kopernikus.camino.network"
	unchainedNode  = "https://kopernikus.unchained.camino.network"
	// columbusNode   = ""
	caminoNode = "https://api.camino.network"

	matrixDebug  = "http://localhost:8008"
	matrixLocal  = "http://localhost:8008"
	matrixC4T    = "https://matrix.chain4travel.com"
	matrixCamino = "https://matrix.camino.network"

	appServiceAccessToken = "wfghWEGh3wgWHEf3478sHFWE"

	caminoFeeKeyStr = "\"PrivateKey-2bMrxpyN24b6BsiTjRDw3h7w7nC75ecZ5vkSyrDksxthvxXQ8o\""

	// 7Sdex3LTEjsnswW38Eb48hQ9insctGrsN : P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3
	localValidator0KeyStr = "\"PrivateKey-vmRQiZeXEXYMyJhEiqdC2z5JhuDbxL8ix9UVvjgMu2Er1NepE\""

	// 6Y3kysjF9jnHnYkdS9yGAuoHyae2eNmeV : P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68
	localValidator1KeyStr = "\"PrivateKey-ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN\""

	//  : P-kopernikus1ftrh6sly2fh4k8rz4wwp60jj4dfdtg2xv3unrj
	//  : P-columbus1ftrh6sly2fh4k8rz4wwp60jj4dfdtg2xxkj3e0
	evgeniiTestKeyStr = "\"PrivateKey-2bMrxpyN24b6BsiTjRDw3h7w7nC75ecZ5vkSyrDksxthvxXQ8o\""
)

var (
	caminoFeeKey       *secp256k1.PrivateKey
	localValidator0Key *secp256k1.PrivateKey
	localValidator1Key *secp256k1.PrivateKey
	evgeniiTestKey     *secp256k1.PrivateKey
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
	evgeniiTestKey = new(secp256k1.PrivateKey)
	if err := evgeniiTestKey.UnmarshalText([]byte(evgeniiTestKeyStr)); err != nil {
		panic(err)
	}
}

// Node

func LocalClient(logger logger.Logger) *node.Client {
	txc, err := node.NewClient(
		localNode,
		logger,
	)
	logger.NoError(err)
	return txc
}

func KopernikusClient(logger logger.Logger) *node.Client {
	txc, err := node.NewClient(
		kopernikusNode,
		logger,
	)
	logger.NoError(err)
	return txc
}

func UnchainedClient(logger logger.Logger) *node.Client {
	txc, err := node.NewClient(
		unchainedNode,
		logger,
	)
	logger.NoError(err)
	return txc
}

func CaminoClient(logger logger.Logger) *node.Client {
	txc, err := node.NewClient(
		caminoNode,
		logger,
	)
	logger.NoError(err)
	return txc
}

// Matrix client

func MatrixDebug(logger logger.Logger) *matrix.Client {
	client, err := matrix.NewClient(matrixLocal, appServiceAccessToken, logger, constants.KopernikusID)
	logger.NoError(err)
	return client
}

func MatrixLocal(logger logger.Logger) *matrix.Client {
	client, err := matrix.NewClient(matrixLocal, appServiceAccessToken, logger, constants.KopernikusID)
	logger.NoError(err)
	return client
}

func MatrixChain4Travel(logger logger.Logger) *matrix.Client {
	client, err := matrix.NewClient(matrixC4T, "", logger, constants.KopernikusID)
	logger.NoError(err)
	return client
}

func MatrixCamino(logger logger.Logger) *matrix.Client {
	client, err := matrix.NewClient(matrixCamino, "v44O3RcIglDKi9BwYI8sNLtrpxHyvLDxXBCqH8CRWpXEC4c3rf", logger, constants.KopernikusID)
	logger.NoError(err)
	return client
}
