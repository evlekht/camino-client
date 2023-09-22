package utils

import (
	"bytes"
	"caminoclient/internal/logger"
	"encoding/hex"
	"os"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	avax_secp256k1 "github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func NewSorting(len int) Sorting {
	return Sorting{Addrs: make([]ids.ShortID, len)}
}

type Sorting struct {
	Addrs []ids.ShortID
}

func (s Sorting) Len() int {
	return len(s.Addrs)
}

func (s Sorting) Swap(i, j int) {
	s.Addrs[i], s.Addrs[j] = s.Addrs[j], s.Addrs[i]
}

func (s Sorting) Less(i, j int) bool {
	return bytes.Compare(s.Addrs[i].Bytes(), s.Addrs[j].Bytes()) < 0
}

var UtilsNoLog = UtilsWithLogger{logger: &logger.NoLog}

func NewUtils(logger logger.Logger) *UtilsWithLogger {
	return &UtilsWithLogger{logger: logger}
}

type UtilsWithLogger struct {
	logger logger.Logger
}

func (u *UtilsWithLogger) ReadAndEncodeToHex(filepath string) (string, error) {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		u.logger.Error(err)
		return "", err
	}

	hex, err := formatting.Encode(formatting.Hex, fileBytes)
	if err != nil {
		u.logger.Error(err)
		return "", err
	}

	u.logger.Infof("%s: %s\n\n", filepath, hex)
	return hex, nil
}

func (u *UtilsWithLogger) AddressesFromIDs(ids []ids.ShortID, networkID uint32) ([]string, error) {
	addrs := make([]string, len(ids))
	for i, id := range ids {
		addr, err := address.Format("P", constants.GetHRP(networkID), id[:])
		if err != nil {
			u.logger.Error(err)
			return nil, err
		}
		addrs[i] = addr
	}

	return addrs, nil
}

func (u *UtilsWithLogger) AddressFromIDStr(idStr string, networkID uint32) (string, error) {
	id, err := ids.ShortFromString(idStr)
	if err != nil {
		u.logger.Error(err)
		return "", err
	}

	addr, err := address.Format("P", constants.GetHRP(networkID), id[:])
	if err != nil {
		u.logger.Error(err)
		return "", err
	}

	u.logger.Infof("id %s, addr %s\n\n", idStr, addr)
	return addr, nil
}

func (u *UtilsWithLogger) GenerateKey(networkID uint32, print bool) (*avax_secp256k1.PrivateKey, string, error) {
	factory := avax_secp256k1.Factory{}
	key, err := factory.NewPrivateKey()
	if err != nil {
		u.logger.Error(err)
		return nil, "", err
	}
	keyAddrStr, err := address.Format("t", constants.GetHRP(networkID), key.Address().Bytes())
	if err != nil {
		u.logger.Error(err)
		return nil, "", err
	}
	if print {
		u.logger.Infof("key: %s", key.String())
		u.logger.Infof("addr-id: %s", key.Address())
		u.logger.Infof("addr: %s", keyAddrStr)
	}
	return key, keyAddrStr, nil
}

func (u *UtilsWithLogger) SignPublicKey(key *avax_secp256k1.PrivateKey, print bool) (signature string, message string, err error) {
	signatureBytes, err := key.Sign(key.PublicKey().Bytes())
	if err != nil {
		u.logger.Error(err)
		return "", "", err
	}
	signature, err = formatting.Encode(formatting.Hex, signatureBytes)
	if err != nil {
		u.logger.Error(err)
		return "", "", err
	}
	message, err = formatting.Encode(formatting.Hex, key.PublicKey().Bytes())
	if err != nil {
		u.logger.Error(err)
		return "", "", err
	}
	if print {
		u.logger.Infof("signature: %s", signature)
		u.logger.Infof("message: %s", message)
	}
	return signature, message, nil
}

func (u *UtilsWithLogger) UncompressedPublicKeyAddress(bytes []byte) error {
	pubKey, err := secp256k1.ParsePubKey(bytes)
	if err != nil {
		u.logger.Error(err)
		return err
	}

	addr, err := ids.ToShortID(hashing.PubkeyBytesToAddress(pubKey.SerializeCompressed()))
	if err != nil {
		u.logger.Error(err)
		return err
	}

	addrStr, err := address.Format("P", constants.GetHRP(constants.ColumbusID), addr[:])
	if err != nil {
		u.logger.Error(err)
		return err
	}
	u.logger.Info(addrStr)
	return nil
}

func (u *UtilsWithLogger) DecodeHexString(str string, avax bool) (decodedBytes []byte, err error) {
	if avax {
		decodedBytes, err = formatting.Decode(formatting.Hex, str)
	} else {
		decodedBytes, err = hex.DecodeString(str)
	}
	if err != nil {
		u.logger.Error(err)
	}
	return decodedBytes, err
}
