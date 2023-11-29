package playground

import (
	"caminoclient/internal/config"
	"caminoclient/internal/logger"
	"caminoclient/internal/node"
	"caminoclient/internal/utils"
	"context"
)

func NewPlayground(ctx context.Context, logger logger.Logger, cfg *config.Config) (*Playground, error) {
	return &Playground{
		logger: logger,
		utils:  utils.NewUtils(logger),
		// creator: LocalClient(logger),
		// issuer:  LocalClient(logger),
		// client:  LocalClient(logger),
	}, nil
}

type Playground struct {
	logger  logger.Logger
	utils   *utils.UtilsWithLogger
	creator *node.Client
	issuer  *node.Client
	client  *node.Client
}

func (p *Playground) Run(ctx context.Context) error {
	// payload := []byte("IRs0DByYNvOhlr5XUwzLlgEOGXTwPpkP")
	// signature, err := evgeniiTestKey.Sign(payload)
	// p.logger.NoError(err)
	// encodedSignature, err := formatting.Encode(formatting.Hex, signature)
	// p.logger.NoError(err)
	// p.logger.Debug(encodedSignature)
	// tx := p.utils.PTX("0x000000002002000003ea00000000000000000000000000000000000000000000000000000000000000000000000259eb48b8b3a928ca9d6b90a0f3492ab47ebf06e9edc553cfb6bcd2d3f38e319a00000007000000003b8b87c000000000000000000000000100000001312f06ca09361eb1699860f36f5ac0635afe8c5659eb48b8b3a928ca9d6b90a0f3492ab47ebf06e9edc553cfb6bcd2d3f38e319a000020010000000000000000000000000000000000000000000000000000000000000000746869732074782069640000000000000000000000000000000000000000000000000007000001d1a94a200000000000000000000000000100000001312f06ca09361eb1699860f36f5ac0635afe8c5600000001b46a7d30a289585d845ee6a0e98be0764a0933674d7325ec298b6fcaf6b1bfd90000000059eb48b8b3a928ca9d6b90a0f3492ab47ebf06e9edc553cfb6bcd2d3f38e319a00000005000001d1e4d5a7c00000000100000000000000009766dbf732231be13e6913a8e0450b9b068659610000000065607c42000000006561cdb3000001d1a94a2000000000000000000b00000000000000000000000100000001312f06ca09361eb1699860f36f5ac0635afe8c56000000000000000a00000001000000000000000200000009000000010b94da4a088c203fa73ebc89d7fa157a87e62fc9ce1b5255bae59b16a7b33c9162ba3dffbe7a6ae6bbf5413be449df8e6394c96f1d82d4e449aa7f9b541709bf0000000009000000010b94da4a088c203fa73ebc89d7fa157a87e62fc9ce1b5255bae59b16a7b33c9162ba3dffbe7a6ae6bbf5413be449df8e6394c96f1d82d4e449aa7f9b541709bf00da7958fd")
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
