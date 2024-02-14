package node

import (
	"caminoclient/internal/utils"
	"context"
	"fmt"
	"sort"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/components/multisig"
	as "github.com/ava-labs/avalanchego/vms/platformvm/addrstate"
	"github.com/ava-labs/avalanchego/vms/platformvm/dac"
	pLocked "github.com/ava-labs/avalanchego/vms/platformvm/locked"
	pTxs "github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
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

	ins, outs, err := c.client.Spend(
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
	ins, outs, err := c.client.Spend(
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
	ins, outs, err := c.client.Spend(
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
	ins, outs, err := c.client.Spend(
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

// func (tc *TxCreator) TLockTx(key *secp256k1.PrivateKey, amtToLock uint64) (*tTxs.Tx, error) {
// 	tc.logger.Info("Creating T-Chain LockTx...")

// 	ins, outs, _, _, err := tc.client.T.SpendWithWrapper(
// 		context.Background(),
// 		key.Address(),
// 		ids.ShortEmpty,
// 		key.Address(),
// 		amtToLock, getNetworkVMParams(tc.networkID).TxFee,
// 		tLocked.StateLocked,
// 		pAPI.Owner{},
// 	)
// 	if err != nil {
// 		tc.logger.Error(err)
// 		return nil, err
// 	}

// 	utx := &tTxs.LockMessengerFundsTx{
// 		BaseTx: tTxs.BaseTx{BaseTx: avax.BaseTx{
// 			NetworkID:    tc.networkID,
// 			BlockchainID: tc.tChainID,
// 			Ins:          ins,
// 			Outs:         outs,
// 		}},
// 	}
// 	signers := make([][]*secp256k1.PrivateKey, len(utx.Ins))
// 	for i := range signers {
// 		signers[i] = []*secp256k1.PrivateKey{key}
// 	}

// 	avax.SortTransferableInputsWithSigners(utx.Ins, signers)
// 	avax.SortTransferableOutputs(utx.Outs, tTxs.Codec)
// 	tx, err := tTxs.NewSigned(utx, tTxs.Codec, signers)
// 	if err != nil {
// 		tc.logger.Error(err)
// 		return nil, err
// 	}
// 	txEncodedBytes, err := formatting.Encode(formatting.Hex, tx.Bytes())
// 	if err != nil {
// 		tc.logger.Error(err)
// 		return nil, err
// 	}
// 	tc.logger.Info(txEncodedBytes)
// 	tc.logger.Infof("txID: %s", tx.ID())
// 	return tx, nil
// }

// func (tc *TxCreator) TCashOutTx(cheque *tTxs.SignedCheque) (*tTxs.Tx, error) {
// 	tc.logger.Info("Creating T-Chain TCashOutTx...")

// 	ins, outs, _, _, err := tc.client.T.SpendWithWrapper(
// 		context.Background(),
// 		cheque.Issuer,
// 		cheque.Agent,
// 		cheque.Beneficiary,
// 		cheque.Amount, 0,
// 		tLocked.StateUnlocked,
// 		pAPI.Owner{},
// 	)
// 	if err != nil {
// 		tc.logger.Error(err)
// 		return nil, err
// 	}

// 	utx := &tTxs.CashoutChequeTx{
// 		BaseTx: tTxs.BaseTx{BaseTx: avax.BaseTx{
// 			NetworkID:    tc.networkID,
// 			BlockchainID: tc.tChainID,
// 			Ins:          ins,
// 			Outs:         outs,
// 		}},
// 		Cheque: *cheque,
// 	}

// 	avax.SortTransferableInputs(utx.Ins)
// 	avax.SortTransferableOutputs(utx.Outs, tTxs.Codec)
// 	tx, err := tTxs.NewSigned(utx, tTxs.Codec, nil)
// 	if err != nil {
// 		tc.logger.Error(err)
// 		return nil, err
// 	}

// 	txEncodedBytes, err := formatting.Encode(formatting.Hex, tx.Bytes())
// 	if err != nil {
// 		tc.logger.Error(err)
// 		return nil, err
// 	}
// 	tc.logger.Info(txEncodedBytes)
// 	tc.logger.Infof("txID: %s", tx.ID())
// 	return tx, nil
// }

// func (tc *TxCreator) CreateCheque(issuerKey *secp256k1.PrivateKey, agent, beneficiary ids.ShortID, amount, serialID uint64, print bool) (*tTxs.SignedCheque, error) {
// 	unsignedCheque := tTxs.Cheque{
// 		Issuer:      issuerKey.Address(),
// 		Agent:       agent,
// 		Beneficiary: beneficiary,
// 		Amount:      amount,
// 		SerialID:    serialID,
// 	}
// 	signature, err := issuerKey.Sign(unsignedCheque.BuildMsgToSign())
// 	if err != nil {
// 		tc.logger.Error(err)
// 		return nil, err
// 	}
// 	credential := &secp256k1fx.Credential{
// 		Sigs: make([][65]byte, 1),
// 	}
// 	copy(credential.Sigs[0][:], signature)
// 	if print {
// 		encodedSignature, err := formatting.Encode(formatting.Hex, signature)
// 		if err != nil {
// 			tc.logger.Error(err)
// 			return nil, err
// 		}
// 		issuer, err := address.Format("T", tc.hrp, unsignedCheque.Issuer.Bytes())
// 		if err != nil {
// 			tc.logger.Error(err)
// 			return nil, err
// 		}
// 		beneficiary, err := address.Format("T", tc.hrp, unsignedCheque.Beneficiary.Bytes())
// 		if err != nil {
// 			tc.logger.Error(err)
// 			return nil, err
// 		}

// 		tc.logger.Infof(`cheque: {
// 			issuer: %s,
// 			agent: %s,
// 			beneficiary: %s,
// 			amount: %d,
// 			serialID: %d,
// 			signature: %s,
// 		}`,
// 			issuer,
// 			unsignedCheque.Agent.String(),
// 			beneficiary,
// 			unsignedCheque.Amount,
// 			unsignedCheque.SerialID,
// 			encodedSignature,
// 		)
// 	}
// 	return &tTxs.SignedCheque{
// 		Cheque: unsignedCheque,
// 		Auth:   credential,
// 	}, nil
// }

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
