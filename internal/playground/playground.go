package playground

import (
	"caminoclient/internal/config"
	"caminoclient/internal/logger"
	"caminoclient/internal/matrix"
	"caminoclient/internal/node"
	"caminoclient/internal/utils"
	"context"

	"github.com/ava-labs/avalanchego/utils/constants"
)

func NewPlayground(ctx context.Context, logger logger.Logger, cfg *config.Config) (*Playground, error) {
	return &Playground{
		logger: logger,
		utils:  utils.NewUtils(logger, constants.KopernikusID),
		// creator: LocalClient(logger),
		// issuer:  LocalClient(logger),
		// client:  LocalClient(logger),
		// matrix: MatrixDebug(logger),
		// matrix: MatrixLocal(logger),
		// matrix: MatrixCamino(logger),
		matrix: MatrixChain4Travel(logger),
	}, nil
}

type Playground struct {
	logger     logger.Logger
	utils      *utils.UtilsWithLogger
	creator    *node.Client
	issuer     *node.Client
	nodeClient *node.Client
	matrix     *matrix.Client
}

func (p *Playground) Run(ctx context.Context) error {
	keys, addresses, err := p.utils.ParseKeysFromFile("unchained")
	p.logger.NoError(err)
	for i := range keys {
		p.logger.Infof("key[%d]: %s", i, keys[i].String())
		p.logger.Infof("addr[%d]: %s", i, addresses[i])
		p.logger.NoError(p.matrix.Register(keys[i], false))
	}
	// key := p.utils.PrivateKey("PrivateKey-2vBHMSNNSEtgdcG2HqVUZvGw9E3VgGmwNBbCsGEGqj3zLmHa83")
	// p.logger.Infof("key: %s", key.String())
	// p.logger.Infof("addr: %s", utils.UtilsNoLog.KeyAddress(key))
	// accessToken, _, err := p.matrix.Login(key)
	// p.logger.NoError(err)
	// p.logger.Debugf("accessToken: %s", accessToken)
	// key1 := localValidator0Key

	// lockTx, err := creator.TLockTx(key1, 1000)
	// p.logger.NoError(err)
	// p.logger.NoError(issuer.IssueTTx(lockTx.Bytes()))

	// time.Sleep(2 * time.Second)

	// key2, _, err := p.utils.ParseKey("\"PrivateKey-6m47M9z8X5Dome27PjyXs5wXjKKLWNhvjEmNPgZVQH1XgQ1C5\"")
	// key, _, err := p.utils.GenerateKey(constants.KopernikusID, true)
	// p.logger.NoError(err)
	// p.logger.Infof("key: %s", key.String())
	// p.logger.Infof("addr: %s", utils.UtilsNoLog.KeyAddress(key))
	// key3, _, err := p.utils.ParseKey("\"PrivateKey-pVPb9P91E67ii2RHT9HWFTtJyNXZhHbzGMCmhcDrFqH4ehSfx\"")
	// key3, _, err := p.utils.GenerateKey(constants.KopernikusID, true)
	// p.logger.NoError(err)
	// cheque, err := creator.CreateCheque(
	// 	key1,
	// 	key2.Address(),
	// 	key3.Address(),
	// 	100,
	// 	2,
	// 	true,
	// )
	// p.logger.NoError(err)
	// cashOutTx, err := creator.TCashOutTx(cheque)
	// p.logger.NoError(err)
	// p.logger.NoError(issuer.IssueTTx(cashOutTx.Bytes()))

	// roomID, err := matrix.CreateRoom()
	// p.logger.NoError(err)
	// p.logger.NoError(matrix.SendMessageWithCheque(roomID, "hello", cheque))
	// p.logger.NoError(p.matrix.Login(key1))

	// tx := p.utils.PTX("0x000000002006000003ea00000000000000000000000000000000000000000000000000000000000000000000000159eb48b8b3a928ca9d6b90a0f3492ab47ebf06e9edc553cfb6bcd2d3f38e319a00000007000003a351fba9800000000000000000000000010000000146a9c04f4bf783aa69daabd519dcf36978168b6600000001fb031c0f8d72aac496e29ae3800f4a892ccc7f22a74121e050a5fbf3c7d7d7d90000000059eb48b8b3a928ca9d6b90a0f3492ab47ebf06e9edc553cfb6bcd2d3f38e319a00000005000003a3520aebc0000000010000000000000000000000010000000900000001e44b0620da0250698d5cab96366e4a733f8e2cdd46a2d08a75f05c60e9d5c41a5749b5e0d66782fbac78027f1d2f6b700c90185b6c851fa0bde5488b635f3dec00ddea5f45")
	// fmt.Println(tx.SyntacticVerify(snow.DefaultContextTest()))
	// tx, err := p.client.GetPTX(p.utils.ID("2NsoB8oNWJLAWPf629soQf8YkN351C79UbNVBreMTSMG8LC7eb"))
	// p.logger.NoError(err)
	// utx, ok := tx.Unsigned.(*txs.AddProposalTx)
	// if !ok {
	// 	p.logger.Fatalf("!ok")
	// }
	// proposal, err := utx.Proposal()
	// p.logger.NoError(err)
	// fmt.Println(proposal.Verify())
	// decoded, err := p.utils.DecodeHexString("000000002010000003ea0000000000000000000000000000000000000000000000000000000000000000000000015e21ded8a9e53a62f6c48ef045b37c938c5c5e9b25a14b4987db93682ca30f7600000007000000002f472a2800000000000000000000000200000001e60e72f2c9b4ea0e13e1be92cf5b73e1cec9d5f70000000121c809bcac64f2170c45d99348cf85ae57e57ad7cd99cd9eff714939fd9a94f5000000005e21ded8a9e53a62f6c48ef045b37c938c5c5e9b25a14b4987db93682ca30f7600000005000000002f566c680000000100000000000000153c703e52532066726f6d2054616977616e3c2f703e0000002a000000002016e60e72f2c9b4ea0e13e1be92cf5b73e1cec9d5f700000000655b085d0000000065aa225de60e72f2c9b4ea0e13e1be92cf5b73e1cec9d5f70000000a0000000100000000000000020000200c000000016b21bd9a3110cd65b5e9887c856ee95d1e009e720847227d93517e81e789bb006c842334a2062ce3a4b4bc93679b65b0dd5cf774e21aab252d15d0be3d6289860000000001000000000000200c000000026b21bd9a3110cd65b5e9887c856ee95d1e009e720847227d93517e81e789bb006c842334a2062ce3a4b4bc93679b65b0dd5cf774e21aab252d15d0be3d6289860009a4872d411f2bd70a4cc2879a50a5165be50787fc147d7bb524c7a29bb55f237681d03ec0531e990e74857c1150432ccec4af1c1ee0ec3c51c857eb0e69c0cf00000000020000000100000002", false)
	// p.logger.NoError(err)
	// tx, err := txs.Parse(txs.Codec, decoded)
	// p.logger.NoError(err)
	// fmt.Printf("%+v", tx)
	// ctx1 := snow.DefaultContextTest()
	// ctx1.NetworkID = 1002
	// tx.SyntacticVerify(ctx1)

	// addressStateTx, err := p.creator.AddressStateTx(evgeniiTestKey.Address(), addrstate.AddressStateBitRoleAdmin, false, localValidator0Key, localValidator0Key)
	// p.logger.NoError(err)
	// p.logger.NoError(p.issuer.IssuePTx(addressStateTx.Bytes()))

	// time.Sleep(5 * time.Second)

	// _, err := p.creator.ProposalTx(&dac.BaseFeeProposal{
	// 	Start:   uint64(time.Now().Unix()) + 10,
	// 	End:     uint64(time.Now().Unix()) + 240,
	// 	Options: []uint64{10},
	// }, evgeniiTestKey, evgeniiTestKey)
	// _, err := p.creator.ProposalTx(
	// 	&dac.ExcludeMemberProposal{
	// 		Start:         uint64(time.Now().Unix()) + 10,
	// 		End:           uint64(time.Now().Unix()) + 10 + dac.ExcludeMemberProposalMinDuration,
	// 		MemberAddress: localValidator0Key.Address(),
	// 	},
	// 	evgeniiTestKey,
	// 	evgeniiTestKey,
	// )
	// _, err := p.creator.ProposalTx(
	// 	&dac.AddMemberProposal{
	// 		Start:            uint64(time.Now().Unix()) + 10,
	// 		End:              uint64(time.Now().Unix()) + 10 + dac.ExcludeMemberProposalMinDuration,
	// 		ApplicantAddress: localValidator0Key.Address(),
	// 	},
	// 	evgeniiTestKey,
	// 	evgeniiTestKey,
	// )
	// proposalTx, err := p.creator.ProposalTx(
	// 	&dac.AdminProposal{
	// 		Proposal: &dac.AddMemberProposal{
	// 			Start:            uint64(time.Now().Unix()) + 10,
	// 			End:              uint64(time.Now().Unix()) + 10 + dac.AddMemberProposalDuration,
	// 			ApplicantAddress: localValidator0Key.Address(),
	// 		},
	// 		OptionIndex: 0,
	// 	},
	// 	evgeniiTestKey,
	// 	evgeniiTestKey,
	// )
	// proposalTx, err := p.creator.ProposalTx(
	// 	&dac.AdminProposal{
	// 		Proposal: &dac.ExcludeMemberProposal{
	// 			Start:         uint64(time.Now().Unix()) + 10,
	// 			End:           uint64(time.Now().Unix()) + 10 + dac.ExcludeMemberProposalMinDuration,
	// 			MemberAddress: localValidator0Key.Address(),
	// 		},
	// 		OptionIndex: 0,
	// 	},
	// 	evgeniiTestKey,
	// 	evgeniiTestKey,
	// )
	// p.logger.NoError(err)
	// p.logger.NoError(p.issuer.IssuePTx(proposalTx.Bytes()))

	// time.Sleep(5 * time.Second)

	// proposalTxID := proposalTx.ID()
	// proposalTxID := p.utils.ID("mYFYkuzAV6tdGPHMUReNAnTPU5u8yC5ogR4PpYQPPxAa52BmT")
	// _, err := p.creator.VoteTx(proposalTxID, 0,
	// 	evgeniiTestKey,
	// 	evgeniiTestKey,
	// )
	// p.logger.NoError(err)
	// p.logger.NoError(p.issuer.IssuePTx(voteTx.Bytes()))

	return nil
}

func (p *Playground) Close(ctx context.Context) error {
	return nil
}
