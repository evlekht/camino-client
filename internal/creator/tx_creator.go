package creator

import (
	"caminoclient/internal/logger"
	"caminoclient/internal/node"
	"caminoclient/internal/utils"
	"context"
	"sort"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/components/multisig"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	pLocked "github.com/ava-labs/avalanchego/vms/platformvm/locked"
	pTxs "github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
)

func NewTxCreator(nodeURI string, feeKey *secp256k1.PrivateKey, logger logger.Logger) (*TxCreator, error) {
	client := node.NewClient(nodeURI, logger)

	nodeCfg, err := client.P.GetConfiguration(context.Background())
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	hrp := constants.GetHRP(uint32(nodeCfg.NetworkID))

	feeKeyAddrStr, err := address.Format("P", hrp, feeKey.Address().Bytes())
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Infof("feeKey: %s : %s : %s\n\n", feeKey.String(), feeKey.Address().String(), feeKeyAddrStr)

	return &TxCreator{
		client:    client,
		nodeCfg:   nodeCfg,
		networkID: uint32(nodeCfg.NetworkID),
		hrp:       hrp,
		logger:    logger,
		utils:     utils.NewUtils(logger),

		FeeKey: feeKey,
	}, nil
}

type TxCreator struct {
	logger    logger.Logger
	client    node.Client
	networkID uint32
	hrp       string
	nodeCfg   *platformvm.GetConfigurationReply
	utils     *utils.UtilsWithLogger

	FeeKey *secp256k1.PrivateKey
}

func (tc *TxCreator) MsigAliasTx(addrs []string, threshold uint32) (*pTxs.Tx, error) {
	tc.logger.Info("Creating P-Chain MsigAliasTx...")
	sorting := utils.NewSorting(len(addrs))
	for i, argAddrStr := range addrs {
		_, _, addrBytes, err := address.Parse(argAddrStr)
		if err != nil {
			tc.logger.Error(err)
			return nil, err
		}
		addr, err := ids.ToShortID(addrBytes)
		if err != nil {
			tc.logger.Error(err)
			return nil, err
		}
		addrStr, err := address.Format("P", tc.hrp, addr.Bytes())
		if err != nil {
			tc.logger.Error(err)
			return nil, err
		}
		sorting.Addrs[i] = addr
		tc.logger.Infof("%s : %s : %s", argAddrStr, addr, addrStr)
	}

	sort.Sort(sorting)

	ins, outs, err := tc.client.Spend(
		context.Background(),
		tc.networkID,
		tc.FeeKey.Address(),
		tc.FeeKey.Address(),
		0, getNetworkVMParams(tc.networkID).TxFee,
		pLocked.StateUnlocked,
	)
	if err != nil {
		tc.logger.Error(err)
		return nil, err
	}
	utx := &pTxs.MultisigAliasTx{
		BaseTx: pTxs.BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    tc.networkID,
			BlockchainID: constants.PlatformChainID,
			Ins:          ins,
			Outs:         outs,
		}},
		MultisigAlias: multisig.Alias{
			Owners: &secp256k1fx.OutputOwners{
				Threshold: threshold,
				Addrs:     sorting.Addrs,
			},
		},
		Auth: &secp256k1fx.Input{},
	}
	signers := make([][]*secp256k1.PrivateKey, len(utx.Ins))
	for i := range signers {
		signers[i] = []*secp256k1.PrivateKey{tc.FeeKey}
	}

	avax.SortTransferableInputsWithSigners(utx.Ins, signers)
	avax.SortTransferableOutputs(utx.Outs, pTxs.Codec)
	tx, err := pTxs.NewSigned(utx, pTxs.Codec, signers)
	if err != nil {
		tc.logger.Error(err)
		return nil, err
	}
	txEncodedBytes, err := formatting.Encode(formatting.Hex, tx.Bytes())
	if err != nil {
		tc.logger.Error(err)
		return nil, err
	}
	tc.logger.Info(txEncodedBytes)
	tc.logger.Infof("txID: %s", tx.ID())

	aliasID := multisig.ComputeAliasID(tx.ID())
	aliasAddrStr, err := address.Format("P", constants.GetHRP(constants.CaminoID), aliasID.Bytes())
	if err != nil {
		tc.logger.Error(err)
		return nil, err
	}
	owners := utx.MultisigAlias.Owners.(*secp256k1fx.OutputOwners)
	aliasAddrs, err := tc.utils.AddressesFromIDs(owners.Addrs, tc.networkID)
	if err != nil {
		tc.logger.Error(err)
		return nil, err
	}
	tc.logger.Infof("alias: %s", aliasAddrStr)
	tc.logger.Infof("alias definition: {\n    threshold: %d\n    addresses: %v\n}", owners.Threshold, aliasAddrs)
	return tx, nil
}

func (tc *TxCreator) AddressStateTx(address ids.ShortID, state pTxs.AddressStateBit, remove bool) (*pTxs.Tx, error) {
	tc.logger.Info("Creating P-Chain AddressStateTx...")
	ins, outs, err := tc.client.Spend(
		context.Background(),
		tc.networkID,
		tc.FeeKey.Address(),
		tc.FeeKey.Address(),
		0, getNetworkVMParams(tc.networkID).TxFee,
		pLocked.StateUnlocked,
	)
	if err != nil {
		tc.logger.Error(err)
		return nil, err
	}
	signers := make([][]*secp256k1.PrivateKey, len(ins))
	for i := range signers {
		signers[i] = []*secp256k1.PrivateKey{tc.FeeKey}
	}
	avax.SortTransferableInputsWithSigners(ins, signers)
	avax.SortTransferableOutputs(outs, pTxs.Codec)

	executorKey := tc.FeeKey

	tx, err := pTxs.NewSigned(&pTxs.AddressStateTx{
		BaseTx: pTxs.BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    tc.networkID,
			BlockchainID: constants.PlatformChainID,
			Ins:          ins,
			Outs:         outs,
		}},
		Address:      address,
		State:        state,
		Remove:       remove,
		Executor:     executorKey.Address(),
		ExecutorAuth: &secp256k1fx.Input{SigIndices: []uint32{0}},
	}, pTxs.Codec, append(signers, []*secp256k1.PrivateKey{executorKey}))
	if err != nil {
		tc.logger.Error(err)
		return nil, err
	}
	txEncodedBytes, err := formatting.Encode(formatting.Hex, tx.Bytes())
	if err != nil {
		tc.logger.Error(err)
		return nil, err
	}
	tc.logger.Info(txEncodedBytes)
	tc.logger.Infof("txID: %s", tx.ID())
	return tx, nil
}

func getNetworkVMParams(networkID uint32) *genesis.Params {
	switch networkID {
	case constants.CaminoID:
		return &genesis.CaminoParams
	case constants.ColumbusID:
		return &genesis.ColumbusParams
	case constants.KopernikusID:
		return &genesis.KopernikusParams
	}
	return &genesis.KopernikusParams
}
