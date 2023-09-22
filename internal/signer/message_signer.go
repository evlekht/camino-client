package signer

import (
	"caminoclient/internal/logger"

	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting"
)

func NewMsgSigner(keyStr string, logger logger.Logger) (*msgSigner, error) {
	key := new(secp256k1.PrivateKey)
	if err := key.UnmarshalText([]byte(keyStr)); err != nil {
		logger.Error(err)
		return nil, err
	}

	return &msgSigner{key: key, logger: logger}, nil
}

type msgSigner struct {
	key    *secp256k1.PrivateKey
	logger logger.Logger
}

func (ms *msgSigner) SignStr(msg string) ([]byte, error) {
	return ms.Sign([]byte(msg))
}

func (ms *msgSigner) Sign(msg []byte) ([]byte, error) {
	sig, err := ms.key.Sign(msg)
	if err != nil {
		ms.logger.Error(err)
		return nil, err
	}

	messageStr, err := formatting.Encode(formatting.Hex, msg)
	if err != nil {
		ms.logger.Error(err)
		return nil, err
	}

	sigStr, err := formatting.Encode(formatting.Hex, sig)
	if err != nil {
		ms.logger.Error(err)
		return nil, err
	}

	ms.logger.Infof("message (hex): %s", messageStr)
	ms.logger.Infof("signature (hex): %s", sigStr)

	return sig, nil
}
