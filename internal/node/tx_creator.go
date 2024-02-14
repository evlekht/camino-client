package node

import (
	"caminoclient/internal/utils"
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/utils/math"
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/components/multisig"
	as "github.com/ava-labs/avalanchego/vms/platformvm/addrstate"
	"github.com/ava-labs/avalanchego/vms/platformvm/dac"
	pLocked "github.com/ava-labs/avalanchego/vms/platformvm/locked"
	pTxs "github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/coreth/plugin/evm"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) MsigAliasTx(addrs []string, threshold uint32, fundsKey *secp256k1.PrivateKey) (*pTxs.Tx, error) {
	c.logger.Info("Creating P-Chain MsigAliasTx...")
	sorting := utils.NewSorting(len(addrs))
	for i, argAddrStr := range addrs {
		_, _, addrBytes, err := address.Parse(argAddrStr)
		if err != nil {
			c.logger.Error(err)
			return nil, err
		}
		addr, err := ids.ToShortID(addrBytes)
		if err != nil {
			c.logger.Error(err)
			return nil, err
		}
		addrStr, err := address.Format("P", c.hrp, addr.Bytes())
		if err != nil {
			c.logger.Error(err)
			return nil, err
		}
		sorting.Addrs[i] = addr
		c.logger.Infof("%s : %s : %s", argAddrStr, addr, addrStr)
	}

	sort.Sort(sorting)

	ins, outs, err := c.client.SpendP(
		context.Background(),
		c.networkID,
		fundsKey.Address(),
		fundsKey.Address(),
		0, getNetworkVMParams(c.networkID).TxFee,
		pLocked.StateUnlocked,
	)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	utx := &pTxs.MultisigAliasTx{
		BaseTx: pTxs.BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    c.networkID,
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
		signers[i] = []*secp256k1.PrivateKey{fundsKey}
	}

	avax.SortTransferableInputsWithSigners(utx.Ins, signers)
	avax.SortTransferableOutputs(utx.Outs, pTxs.Codec)
	tx, err := pTxs.NewSigned(utx, pTxs.Codec, signers)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	txEncodedBytes, err := formatting.Encode(formatting.Hex, tx.Bytes())
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	c.logger.Info(txEncodedBytes)
	c.logger.Infof("txID: %s", tx.ID())

	aliasID := multisig.ComputeAliasID(tx.ID())
	aliasAddrStr, err := address.Format("P", constants.GetHRP(constants.CaminoID), aliasID.Bytes())
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	owners := utx.MultisigAlias.Owners.(*secp256k1fx.OutputOwners)
	aliasAddrs, err := c.utils.AddressesFromIDs(owners.Addrs, c.networkID)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	c.logger.Infof("alias: %s", aliasAddrStr)
	c.logger.Infof("alias definition: {\n    threshold: %d\n    addresses: %v\n}", owners.Threshold, aliasAddrs)
	return tx, nil
}

func (c *Client) AddressStateTx(address ids.ShortID, state as.AddressStateBit, remove bool, fundsKey, executorKey *secp256k1.PrivateKey) (*pTxs.Tx, error) {
	c.logger.Info("Creating P-Chain AddressStateTx...")
	ins, outs, err := c.client.SpendP(
		context.Background(),
		c.networkID,
		fundsKey.Address(),
		fundsKey.Address(),
		0, getNetworkVMParams(c.networkID).TxFee,
		pLocked.StateUnlocked,
	)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	signers := make([][]*secp256k1.PrivateKey, len(ins))
	for i := range signers {
		signers[i] = []*secp256k1.PrivateKey{fundsKey}
	}
	avax.SortTransferableInputsWithSigners(ins, signers)
	avax.SortTransferableOutputs(outs, pTxs.Codec)

	tx, err := pTxs.NewSigned(&pTxs.AddressStateTx{
		BaseTx: pTxs.BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    c.networkID,
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
		c.logger.Error(err)
		return nil, err
	}
	txEncodedBytes, err := formatting.Encode(formatting.Hex, tx.Bytes())
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	c.logger.Info(txEncodedBytes)
	c.logger.Infof("txID: %s", tx.ID())
	return tx, nil
}

func (c *Client) ProposalTx(
	proposal dac.Proposal,
	fundsKey *secp256k1.PrivateKey,
	proposerKey *secp256k1.PrivateKey,
) (*pTxs.Tx, error) {
	c.logger.Info("Creating P-Chain AddProposalTx...")
	vmParams := getNetworkVMParams(c.networkID)
	ins, outs, err := c.client.SpendP(
		context.Background(),
		c.networkID,
		fundsKey.Address(),
		fundsKey.Address(),
		vmParams.CaminoConfig.DACProposalBondAmount, vmParams.TxFee,
		pLocked.StateUnlocked,
	)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	signers := make([][]*secp256k1.PrivateKey, len(ins))
	for i := range signers {
		signers[i] = []*secp256k1.PrivateKey{fundsKey}
	}
	avax.SortTransferableInputsWithSigners(ins, signers)
	avax.SortTransferableOutputs(outs, pTxs.Codec)

	wrappedProposal := &pTxs.ProposalWrapper{Proposal: proposal}
	proposalBytes, err := pTxs.Codec.Marshal(pTxs.Version, wrappedProposal)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	tx, err := pTxs.NewSigned(&pTxs.AddProposalTx{
		BaseTx: pTxs.BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    c.networkID,
			BlockchainID: constants.PlatformChainID,
			Ins:          ins,
			Outs:         outs,
		}},
		ProposalPayload: proposalBytes,
		ProposerAddress: proposerKey.Address(),
		ProposerAuth:    &secp256k1fx.Input{SigIndices: []uint32{0}},
	}, pTxs.Codec, append(signers, []*secp256k1.PrivateKey{proposerKey}))
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	txEncodedBytes, err := formatting.Encode(formatting.Hex, tx.Bytes())
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	fmt.Printf("%s\n", txEncodedBytes)
	fmt.Printf("txID: %s\n", tx.ID())
	return tx, nil
}

func (c *Client) VoteTx(
	proposalID ids.ID,
	optionIndex uint32,
	fundsKey *secp256k1.PrivateKey,
	voterKey *secp256k1.PrivateKey,
) (*pTxs.Tx, error) {
	c.logger.Info("Creating P-Chain AddVoteTx...")
	ins, outs, err := c.client.SpendP(
		context.Background(),
		c.networkID,
		fundsKey.Address(),
		fundsKey.Address(),
		0, getNetworkVMParams(c.networkID).TxFee,
		pLocked.StateUnlocked,
	)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	signers := make([][]*secp256k1.PrivateKey, len(ins))
	for i := range signers {
		signers[i] = []*secp256k1.PrivateKey{fundsKey}
	}
	avax.SortTransferableInputsWithSigners(ins, signers)
	avax.SortTransferableOutputs(outs, pTxs.Codec)

	vote := &pTxs.VoteWrapper{Vote: &dac.SimpleVote{OptionIndex: optionIndex}}
	voteBytes, err := pTxs.Codec.Marshal(pTxs.Version, vote)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	tx, err := pTxs.NewSigned(&pTxs.AddVoteTx{
		BaseTx: pTxs.BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    c.networkID,
			BlockchainID: constants.PlatformChainID,
			Ins:          ins,
			Outs:         outs,
		}},
		ProposalID:   proposalID,
		VotePayload:  voteBytes,
		VoterAddress: voterKey.Address(),
		VoterAuth:    &secp256k1fx.Input{SigIndices: []uint32{0}},
	}, pTxs.Codec, append(signers, []*secp256k1.PrivateKey{voterKey}))
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	txEncodedBytes, err := formatting.Encode(formatting.Hex, tx.Bytes())
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	fmt.Printf("%s\n", txEncodedBytes)
	fmt.Printf("txID: %s\n", tx.ID())
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

func (c *Client) EVMTx(amountToExport uint64, recipientAddr ids.ShortID, fundsKey *secp256k1.PrivateKey, targetChain string) (*evm.Tx, error) {
	c.logger.Info("Creating C-Chain exportTx...")

	destinationChainID, err := c.getChainID(targetChain)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	senderAddr := evm.GetEthAddress(fundsKey)
	nonce, err := c.client.CETH.NonceAt(context.Background(), senderAddr, nil)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	amountToConsume := amountToExport
	outs := []*avax.TransferableOutput{{
		Asset: avax.Asset{ID: c.avaxAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: amountToExport,
			OutputOwners: secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{recipientAddr},
			},
		},
	}}

	// calculate fee

	utx := &evm.UnsignedExportTx{
		NetworkID:        c.networkID,
		BlockchainID:     c.cChainID,
		DestinationChain: destinationChainID,
		ExportedOutputs:  outs,
	}
	tx := &evm.Tx{UnsignedAtomicTx: utx}
	if err := tx.Sign(evm.Codec, nil); err != nil {
		c.logger.Error(err)
		return nil, err
	}

	txGasUsed, err := tx.GasUsed(true)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	baseFee, err := c.client.CETH.EstimateBaseFee(context.Background())
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	newTxGasUsed, err := math.Add64(txGasUsed, EVMInputGas)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	txGasUsed = newTxGasUsed

	fee, err := calculateEVMDynamicFee(txGasUsed, baseFee)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	newAmount, err := math.Add64(amountToConsume, fee)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	amountToConsume = newAmount

	// create tx

	utx = &evm.UnsignedExportTx{
		NetworkID:        c.networkID,
		BlockchainID:     c.cChainID,
		DestinationChain: destinationChainID,
		Ins: []evm.EVMInput{{
			Address: senderAddr,
			Amount:  amountToConsume,
			AssetID: c.avaxAssetID,
			Nonce:   nonce,
		}},
		ExportedOutputs: outs,
	}
	tx = &evm.Tx{UnsignedAtomicTx: utx}
	if err := tx.Sign(evm.Codec, [][]*secp256k1.PrivateKey{{fundsKey}}); err != nil {
		c.logger.Error(err)
		return nil, err
	}

	txEncodedBytes, err := formatting.Encode(formatting.Hex, tx.Bytes())
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	c.logger.Info(txEncodedBytes)
	c.logger.Infof("txID: %s", tx.ID())
	return tx, nil
}

// copy-paste from evm
const (
	x2cRateInt64       int64 = 1_000_000_000
	x2cRateMinus1Int64 int64 = x2cRateInt64 - 1
)

// copy-paste from evm
var (
	// x2cRate is the conversion rate between the smallest denomination on the X-Chain
	// 1 nAVAX and the smallest denomination on the C-Chain 1 wei. Where 1 nAVAX = 1 gWei.
	// This is only required for AVAX because the denomination of 1 AVAX is 9 decimal
	// places on the X and P chains, but is 18 decimal places within the EVM.
	x2cRate       = big.NewInt(x2cRateInt64)
	x2cRateMinus1 = big.NewInt(x2cRateMinus1Int64)

	errNilBaseFee         = errors.New("cannot calculate dynamic fee with nil baseFee")
	errFeeOverflow        = errors.New("overflow occurred while calculating the fee")
	TxBytesGas     uint64 = 1
	EVMInputGas    uint64 = (common.AddressLength+wrappers.LongLen+hashing.HashLen+wrappers.LongLen)*TxBytesGas + secp256k1fx.CostPerSignature
)

// copy-paste from evm
//
// calculates the amount of AVAX that must be burned by an atomic transaction
// that consumes [cost] at [baseFee].
func calculateEVMDynamicFee(cost uint64, baseFee *big.Int) (uint64, error) {
	if baseFee == nil {
		return 0, errNilBaseFee
	}
	bigCost := new(big.Int).SetUint64(cost)
	fee := new(big.Int).Mul(bigCost, baseFee)
	feeToRoundUp := new(big.Int).Add(fee, x2cRateMinus1)
	feeInNAVAX := new(big.Int).Div(feeToRoundUp, x2cRate)
	if !feeInNAVAX.IsUint64() {
		// the fee is more than can fit in a uint64
		return 0, errFeeOverflow
	}
	return feeInNAVAX.Uint64(), nil
}
