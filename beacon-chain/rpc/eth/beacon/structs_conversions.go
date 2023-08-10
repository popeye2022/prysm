package beacon

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	bytesutil2 "github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/v4/math"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/wealdtech/go-bytesutil"
)

func (b *SignedBeaconBlock) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}

	bl, err := b.Message.ToGeneric()
	if err != nil {
		return nil, err
	}
	ph, ok := bl.Block.(*eth.GenericBeaconBlock_Phase0)
	if !ok {
		return nil, errors.New("wrong block type")
	}
	block := &eth.SignedBeaconBlock{
		Block:     ph.Phase0,
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_Phase0{Phase0: block}}, nil
}

func (b *BeaconBlock) ToGeneric() (*eth.GenericBeaconBlock, error) {
	slot, err := strconv.ParseUint(b.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}

	block := &eth.BeaconBlock{
		Slot:          primitives.Slot(slot),
		ProposerIndex: primitives.ValidatorIndex(proposerIndex),
		ParentRoot:    parentRoot,
		StateRoot:     stateRoot,
		Body: &eth.BeaconBlockBody{
			RandaoReveal: randaoReveal,
			Eth1Data: &eth.Eth1Data{
				DepositRoot:  depositRoot,
				DepositCount: depositCount,
				BlockHash:    blockHash,
			},
			Graffiti:          graffiti,
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
		},
	}
	return &eth.GenericBeaconBlock{Block: &eth.GenericBeaconBlock_Phase0{Phase0: block}}, nil
}

func (b *SignedBeaconBlockAltair) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	bl, err := b.Message.ToGeneric()
	if err != nil {
		return nil, err
	}
	altair, ok := bl.Block.(*eth.GenericBeaconBlock_Altair)
	if !ok {
		return nil, errors.New("wrong block type")
	}
	block := &eth.SignedBeaconBlockAltair{
		Block:     altair.Altair,
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_Altair{Altair: block}}, nil
}

func (b *BeaconBlockAltair) ToGeneric() (*eth.GenericBeaconBlock, error) {
	slot, err := strconv.ParseUint(b.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := bytesutil.FromHexString(b.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(b.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.SyncAggregate.SyncCommitteeSignature")
	}

	block := &eth.BeaconBlockAltair{
		Slot:          primitives.Slot(slot),
		ProposerIndex: primitives.ValidatorIndex(proposerIndex),
		ParentRoot:    parentRoot,
		StateRoot:     stateRoot,
		Body: &eth.BeaconBlockBodyAltair{
			RandaoReveal: randaoReveal,
			Eth1Data: &eth.Eth1Data{
				DepositRoot:  depositRoot,
				DepositCount: depositCount,
				BlockHash:    blockHash,
			},
			Graffiti:          graffiti,
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &eth.SyncAggregate{
				SyncCommitteeBits:      syncCommitteeBits,
				SyncCommitteeSignature: syncCommitteeSig,
			},
		},
	}
	return &eth.GenericBeaconBlock{Block: &eth.GenericBeaconBlock_Altair{Altair: block}}, nil
}

func (b *SignedBeaconBlockBellatrix) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	bl, err := b.Message.ToGeneric()
	if err != nil {
		return nil, err
	}
	bellatrix, ok := bl.Block.(*eth.GenericBeaconBlock_Bellatrix)
	if !ok {
		return nil, errors.New("wrong block type")
	}

	block := &eth.SignedBeaconBlockBellatrix{
		Block:     bellatrix.Bellatrix,
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_Bellatrix{Bellatrix: block}}, nil
}

func (b *BeaconBlockBellatrix) ToGeneric() (*eth.GenericBeaconBlock, error) {
	slot, err := strconv.ParseUint(b.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := hexutil.Decode(b.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(b.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(b.Body.ExecutionPayload.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(b.Body.ExecutionPayload.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(b.Body.ExecutionPayload.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(b.Body.ExecutionPayload.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(b.Body.ExecutionPayload.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(b.Body.ExecutionPayload.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(b.Body.ExecutionPayload.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(b.Body.ExecutionPayload.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(b.Body.ExecutionPayload.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(b.Body.ExecutionPayload.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(b.Body.ExecutionPayload.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToSSZBytes(b.Body.ExecutionPayload.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(b.Body.ExecutionPayload.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.BlockHash")
	}
	payloadTxs := make([][]byte, len(b.Body.ExecutionPayload.Transactions))
	for i, tx := range b.Body.ExecutionPayload.Transactions {
		payloadTxs[i], err = hexutil.Decode(tx)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Body.ExecutionPayload.Transactions[%d]", i)
		}
	}

	block := &eth.BeaconBlockBellatrix{
		Slot:          primitives.Slot(slot),
		ProposerIndex: primitives.ValidatorIndex(proposerIndex),
		ParentRoot:    parentRoot,
		StateRoot:     stateRoot,
		Body: &eth.BeaconBlockBodyBellatrix{
			RandaoReveal: randaoReveal,
			Eth1Data: &eth.Eth1Data{
				DepositRoot:  depositRoot,
				DepositCount: depositCount,
				BlockHash:    blockHash,
			},
			Graffiti:          graffiti,
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &eth.SyncAggregate{
				SyncCommitteeBits:      syncCommitteeBits,
				SyncCommitteeSignature: syncCommitteeSig,
			},
			ExecutionPayload: &enginev1.ExecutionPayload{
				ParentHash:    payloadParentHash,
				FeeRecipient:  payloadFeeRecipient,
				StateRoot:     payloadStateRoot,
				ReceiptsRoot:  payloadReceiptsRoot,
				LogsBloom:     payloadLogsBloom,
				PrevRandao:    payloadPrevRandao,
				BlockNumber:   payloadBlockNumber,
				GasLimit:      payloadGasLimit,
				GasUsed:       payloadGasUsed,
				Timestamp:     payloadTimestamp,
				ExtraData:     payloadExtraData,
				BaseFeePerGas: payloadBaseFeePerGas,
				BlockHash:     payloadBlockHash,
				Transactions:  payloadTxs,
			},
		},
	}
	return &eth.GenericBeaconBlock{Block: &eth.GenericBeaconBlock_Bellatrix{Bellatrix: block}}, nil
}

func (b *SignedBlindedBeaconBlockBellatrix) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	bl, err := b.Message.ToGeneric()
	if err != nil {
		return nil, err
	}
	bellatrix, ok := bl.Block.(*eth.GenericBeaconBlock_BlindedBellatrix)
	if !ok {
		return nil, errors.New("wrong block type")
	}
	block := &eth.SignedBlindedBeaconBlockBellatrix{
		Block:     bellatrix.BlindedBellatrix,
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_BlindedBellatrix{BlindedBellatrix: block}, IsBlinded: true, PayloadValue: 0 /* can't get payload value from blinded block */}, nil
}

func (b *BlindedBeaconBlockBellatrix) ToGeneric() (*eth.GenericBeaconBlock, error) {
	slot, err := strconv.ParseUint(b.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := hexutil.Decode(b.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(b.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(b.Body.ExecutionPayloadHeader.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(b.Body.ExecutionPayloadHeader.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(b.Body.ExecutionPayloadHeader.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(b.Body.ExecutionPayloadHeader.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToSSZBytes(b.Body.ExecutionPayloadHeader.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.BlockHash")
	}
	payloadTxsRoot, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.TransactionsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.TransactionsRoot")
	}

	block := &eth.BlindedBeaconBlockBellatrix{
		Slot:          primitives.Slot(slot),
		ProposerIndex: primitives.ValidatorIndex(proposerIndex),
		ParentRoot:    parentRoot,
		StateRoot:     stateRoot,
		Body: &eth.BlindedBeaconBlockBodyBellatrix{
			RandaoReveal: randaoReveal,
			Eth1Data: &eth.Eth1Data{
				DepositRoot:  depositRoot,
				DepositCount: depositCount,
				BlockHash:    blockHash,
			},
			Graffiti:          graffiti,
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &eth.SyncAggregate{
				SyncCommitteeBits:      syncCommitteeBits,
				SyncCommitteeSignature: syncCommitteeSig,
			},
			ExecutionPayloadHeader: &enginev1.ExecutionPayloadHeader{
				ParentHash:       payloadParentHash,
				FeeRecipient:     payloadFeeRecipient,
				StateRoot:        payloadStateRoot,
				ReceiptsRoot:     payloadReceiptsRoot,
				LogsBloom:        payloadLogsBloom,
				PrevRandao:       payloadPrevRandao,
				BlockNumber:      payloadBlockNumber,
				GasLimit:         payloadGasLimit,
				GasUsed:          payloadGasUsed,
				Timestamp:        payloadTimestamp,
				ExtraData:        payloadExtraData,
				BaseFeePerGas:    payloadBaseFeePerGas,
				BlockHash:        payloadBlockHash,
				TransactionsRoot: payloadTxsRoot,
			},
		},
	}
	return &eth.GenericBeaconBlock{Block: &eth.GenericBeaconBlock_BlindedBellatrix{BlindedBellatrix: block}, IsBlinded: true, PayloadValue: 0 /* can't get payload value from blinded block */}, nil
}

func (b *SignedBeaconBlockCapella) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	bl, err := b.Message.ToGeneric()
	if err != nil {
		return nil, err
	}
	capella, ok := bl.Block.(*eth.GenericBeaconBlock_Capella)
	if !ok {
		return nil, errors.New("wrong block type")
	}

	block := &eth.SignedBeaconBlockCapella{
		Block:     capella.Capella,
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_Capella{Capella: block}}, nil
}

func (b *BeaconBlockCapella) ToGeneric() (*eth.GenericBeaconBlock, error) {
	slot, err := strconv.ParseUint(b.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := hexutil.Decode(b.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(b.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(b.Body.ExecutionPayload.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(b.Body.ExecutionPayload.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(b.Body.ExecutionPayload.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(b.Body.ExecutionPayload.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(b.Body.ExecutionPayload.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(b.Body.ExecutionPayload.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(b.Body.ExecutionPayload.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(b.Body.ExecutionPayload.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(b.Body.ExecutionPayload.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(b.Body.ExecutionPayload.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(b.Body.ExecutionPayload.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToSSZBytes(b.Body.ExecutionPayload.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(b.Body.ExecutionPayload.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayload.BlockHash")
	}
	txs := make([][]byte, len(b.Body.ExecutionPayload.Transactions))
	for i, tx := range b.Body.ExecutionPayload.Transactions {
		txs[i], err = hexutil.Decode(tx)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Body.ExecutionPayload.Transactions[%d]", i)
		}
	}
	withdrawals := make([]*enginev1.Withdrawal, len(b.Body.ExecutionPayload.Withdrawals))
	for i, w := range b.Body.ExecutionPayload.Withdrawals {
		withdrawalIndex, err := strconv.ParseUint(w.WithdrawalIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Body.ExecutionPayload.Withdrawals[%d].WithdrawalIndex", i)
		}
		validatorIndex, err := strconv.ParseUint(w.ValidatorIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Body.ExecutionPayload.Withdrawals[%d].ValidatorIndex", i)
		}
		address, err := hexutil.Decode(w.ExecutionAddress)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Body.ExecutionPayload.Withdrawals[%d].ExecutionAddress", i)
		}
		amount, err := strconv.ParseUint(w.Amount, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Body.ExecutionPayload.Withdrawals[%d].Amount", i)
		}
		withdrawals[i] = &enginev1.Withdrawal{
			Index:          withdrawalIndex,
			ValidatorIndex: primitives.ValidatorIndex(validatorIndex),
			Address:        address,
			Amount:         amount,
		}
	}
	blsChanges, err := convertBlsChanges(b.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}

	block := &eth.BeaconBlockCapella{
		Slot:          primitives.Slot(slot),
		ProposerIndex: primitives.ValidatorIndex(proposerIndex),
		ParentRoot:    parentRoot,
		StateRoot:     stateRoot,
		Body: &eth.BeaconBlockBodyCapella{
			RandaoReveal: randaoReveal,
			Eth1Data: &eth.Eth1Data{
				DepositRoot:  depositRoot,
				DepositCount: depositCount,
				BlockHash:    blockHash,
			},
			Graffiti:          graffiti,
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &eth.SyncAggregate{
				SyncCommitteeBits:      syncCommitteeBits,
				SyncCommitteeSignature: syncCommitteeSig,
			},
			ExecutionPayload: &enginev1.ExecutionPayloadCapella{
				ParentHash:    payloadParentHash,
				FeeRecipient:  payloadFeeRecipient,
				StateRoot:     payloadStateRoot,
				ReceiptsRoot:  payloadReceiptsRoot,
				LogsBloom:     payloadLogsBloom,
				PrevRandao:    payloadPrevRandao,
				BlockNumber:   payloadBlockNumber,
				GasLimit:      payloadGasLimit,
				GasUsed:       payloadGasUsed,
				Timestamp:     payloadTimestamp,
				ExtraData:     payloadExtraData,
				BaseFeePerGas: payloadBaseFeePerGas,
				BlockHash:     payloadBlockHash,
				Transactions:  txs,
				Withdrawals:   withdrawals,
			},
			BlsToExecutionChanges: blsChanges,
		},
	}
	return &eth.GenericBeaconBlock{Block: &eth.GenericBeaconBlock_Capella{Capella: block}}, nil
}

func (b *SignedBlindedBeaconBlockCapella) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	bl, err := b.Message.ToGeneric()
	if err != nil {
		return nil, err
	}
	capella, ok := bl.Block.(*eth.GenericBeaconBlock_BlindedCapella)
	if !ok {
		return nil, errors.New("wrong block type")
	}
	block := &eth.SignedBlindedBeaconBlockCapella{
		Block:     capella.BlindedCapella,
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_BlindedCapella{BlindedCapella: block}, IsBlinded: true, PayloadValue: 0 /* can't get payload value from blinded block */}, nil
}

func (b *BlindedBeaconBlockCapella) ToGeneric() (*eth.GenericBeaconBlock, error) {
	slot, err := strconv.ParseUint(b.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := hexutil.Decode(b.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(b.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(b.Body.ExecutionPayloadHeader.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(b.Body.ExecutionPayloadHeader.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(b.Body.ExecutionPayloadHeader.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(b.Body.ExecutionPayloadHeader.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToSSZBytes(b.Body.ExecutionPayloadHeader.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.BlockHash")
	}
	payloadTxsRoot, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.TransactionsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.TransactionsRoot")
	}
	payloadWithdrawalsRoot, err := hexutil.Decode(b.Body.ExecutionPayloadHeader.WithdrawalsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Body.ExecutionPayloadHeader.WithdrawalsRoot")
	}
	blsChanges, err := convertBlsChanges(b.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}

	block := &eth.BlindedBeaconBlockCapella{
		Slot:          primitives.Slot(slot),
		ProposerIndex: primitives.ValidatorIndex(proposerIndex),
		ParentRoot:    parentRoot,
		StateRoot:     stateRoot,
		Body: &eth.BlindedBeaconBlockBodyCapella{
			RandaoReveal: randaoReveal,
			Eth1Data: &eth.Eth1Data{
				DepositRoot:  depositRoot,
				DepositCount: depositCount,
				BlockHash:    blockHash,
			},
			Graffiti:          graffiti,
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &eth.SyncAggregate{
				SyncCommitteeBits:      syncCommitteeBits,
				SyncCommitteeSignature: syncCommitteeSig,
			},
			ExecutionPayloadHeader: &enginev1.ExecutionPayloadHeaderCapella{
				ParentHash:       payloadParentHash,
				FeeRecipient:     payloadFeeRecipient,
				StateRoot:        payloadStateRoot,
				ReceiptsRoot:     payloadReceiptsRoot,
				LogsBloom:        payloadLogsBloom,
				PrevRandao:       payloadPrevRandao,
				BlockNumber:      payloadBlockNumber,
				GasLimit:         payloadGasLimit,
				GasUsed:          payloadGasUsed,
				Timestamp:        payloadTimestamp,
				ExtraData:        payloadExtraData,
				BaseFeePerGas:    payloadBaseFeePerGas,
				BlockHash:        payloadBlockHash,
				TransactionsRoot: payloadTxsRoot,
				WithdrawalsRoot:  payloadWithdrawalsRoot,
			},
			BlsToExecutionChanges: blsChanges,
		},
	}
	return &eth.GenericBeaconBlock{Block: &eth.GenericBeaconBlock_BlindedCapella{BlindedCapella: block}, IsBlinded: true, PayloadValue: 0 /* can't get payload value from blinded block */}, nil
}

func (b *SignedBeaconBlockContentsDeneb) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	var signedBlobSidecars []*eth.SignedBlobSidecar
	if len(b.SignedBlobSidecars) != 0 {
		signedBlobSidecars = make([]*eth.SignedBlobSidecar, len(b.SignedBlobSidecars))
		for i := range b.SignedBlobSidecars {
			signedBlob, err := convertToSignedBlobSidecar(i, b.SignedBlobSidecars[i])
			if err != nil {
				return nil, err
			}
			signedBlobSidecars[i] = signedBlob
		}
	}
	signedDenebBlock, err := convertToSignedDenebBlock(b.SignedBlock)
	if err != nil {
		return nil, err
	}
	block := &eth.SignedBeaconBlockAndBlobsDeneb{
		Block: signedDenebBlock,
		Blobs: signedBlobSidecars,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_Deneb{Deneb: block}}, nil
}

func (b *SignedBeaconBlockContentsDeneb) ToUnsigned() *BeaconBlockContentsDeneb {
	var blobSidecars []*BlobSidecar
	if len(b.SignedBlobSidecars) != 0 {
		blobSidecars = make([]*BlobSidecar, len(b.SignedBlobSidecars))
		for i, s := range b.SignedBlobSidecars {
			blobSidecars[i] = s.Message
		}
	}
	return &BeaconBlockContentsDeneb{
		Block:        b.SignedBlock.Message,
		BlobSidecars: blobSidecars,
	}
}

func (b *BeaconBlockContentsDeneb) ToGeneric() (*eth.GenericBeaconBlock, error) {
	var blobSidecars []*eth.BlobSidecar
	if len(b.BlobSidecars) != 0 {
		blobSidecars = make([]*eth.BlobSidecar, len(b.BlobSidecars))
		for i := range b.BlobSidecars {
			blob, err := convertToBlobSidecar(i, b.BlobSidecars[i])
			if err != nil {
				return nil, err
			}
			blobSidecars[i] = blob
		}
	}
	denebBlock, err := convertToDenebBlock(b.Block)
	if err != nil {
		return nil, err
	}
	block := &eth.BeaconBlockAndBlobsDeneb{
		Block: denebBlock,
		Blobs: blobSidecars,
	}
	return &eth.GenericBeaconBlock{Block: &eth.GenericBeaconBlock_Deneb{Deneb: block}}, nil
}

func (b *SignedBlindedBeaconBlockContentsDeneb) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	var signedBlindedBlobSidecars []*eth.SignedBlindedBlobSidecar
	if len(b.SignedBlindedBlobSidecars) != 0 {
		signedBlindedBlobSidecars = make([]*eth.SignedBlindedBlobSidecar, len(b.SignedBlindedBlobSidecars))
		for i := range b.SignedBlindedBlobSidecars {
			signedBlob, err := convertToSignedBlindedBlobSidecar(i, b.SignedBlindedBlobSidecars[i])
			if err != nil {
				return nil, err
			}
			signedBlindedBlobSidecars[i] = signedBlob
		}
	}
	signedBlindedBlock, err := convertToSignedBlindedDenebBlock(b.SignedBlindedBlock)
	if err != nil {
		return nil, err
	}
	block := &eth.SignedBlindedBeaconBlockAndBlobsDeneb{
		Block: signedBlindedBlock,
		Blobs: signedBlindedBlobSidecars,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_BlindedDeneb{BlindedDeneb: block}, IsBlinded: true, PayloadValue: 0 /* can't get payload value from blinded block */}, nil
}

func (b *SignedBlindedBeaconBlockContentsDeneb) ToUnsigned() *BlindedBeaconBlockContentsDeneb {
	var blobSidecars []*BlindedBlobSidecar
	if len(b.SignedBlindedBlobSidecars) != 0 {
		blobSidecars = make([]*BlindedBlobSidecar, len(b.SignedBlindedBlobSidecars))
		for i := range b.SignedBlindedBlobSidecars {
			blobSidecars[i] = b.SignedBlindedBlobSidecars[i].Message
		}
	}
	return &BlindedBeaconBlockContentsDeneb{
		BlindedBlock:        b.SignedBlindedBlock.Message,
		BlindedBlobSidecars: blobSidecars,
	}
}

func (b *BlindedBeaconBlockContentsDeneb) ToGeneric() (*eth.GenericBeaconBlock, error) {
	var blindedBlobSidecars []*eth.BlindedBlobSidecar
	if len(b.BlindedBlobSidecars) != 0 {
		blindedBlobSidecars = make([]*eth.BlindedBlobSidecar, len(b.BlindedBlobSidecars))
		for i := range b.BlindedBlobSidecars {
			blob, err := convertToBlindedBlobSidecar(i, b.BlindedBlobSidecars[i])
			if err != nil {
				return nil, err
			}
			blindedBlobSidecars[i] = blob
		}
	}
	blindedBlock, err := convertToBlindedDenebBlock(b.BlindedBlock)
	if err != nil {
		return nil, err
	}
	block := &eth.BlindedBeaconBlockAndBlobsDeneb{
		Block: blindedBlock,
		Blobs: blindedBlobSidecars,
	}
	return &eth.GenericBeaconBlock{Block: &eth.GenericBeaconBlock_BlindedDeneb{BlindedDeneb: block}, IsBlinded: true, PayloadValue: 0 /* can't get payload value from blinded block */}, nil
}

func convertInternalBeaconBlock(b *eth.BeaconBlock) (*BeaconBlock, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	return &BeaconBlock{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BeaconBlockBody{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
		},
	}, nil
}

func convertInternalBeaconBlockAltair(b *eth.BeaconBlockAltair) (*BeaconBlockAltair, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}

	return &BeaconBlockAltair{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BeaconBlockBodyAltair{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
		},
	}, nil
}

func convertInternalBlindedBeaconBlockBellatrix(b *eth.BlindedBeaconBlockBellatrix) (*BlindedBeaconBlockBellatrix, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	baseFeePerGas, err := sszBytesToUint256String(b.Body.ExecutionPayloadHeader.BaseFeePerGas)
	if err != nil {
		return nil, err
	}
	return &BlindedBeaconBlockBellatrix{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BlindedBeaconBlockBodyBellatrix{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayloadHeader: &ExecutionPayloadHeader{
				ParentHash:       hexutil.Encode(b.Body.ExecutionPayloadHeader.ParentHash),
				FeeRecipient:     hexutil.Encode(b.Body.ExecutionPayloadHeader.FeeRecipient),
				StateRoot:        hexutil.Encode(b.Body.ExecutionPayloadHeader.StateRoot),
				ReceiptsRoot:     hexutil.Encode(b.Body.ExecutionPayloadHeader.ReceiptsRoot),
				LogsBloom:        hexutil.Encode(b.Body.ExecutionPayloadHeader.LogsBloom),
				PrevRandao:       hexutil.Encode(b.Body.ExecutionPayloadHeader.PrevRandao),
				BlockNumber:      fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.BlockNumber),
				GasLimit:         fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasLimit),
				GasUsed:          fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasUsed),
				Timestamp:        fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.Timestamp),
				ExtraData:        hexutil.Encode(b.Body.ExecutionPayloadHeader.ExtraData),
				BaseFeePerGas:    baseFeePerGas,
				BlockHash:        hexutil.Encode(b.Body.ExecutionPayloadHeader.BlockHash),
				TransactionsRoot: hexutil.Encode(b.Body.ExecutionPayloadHeader.TransactionsRoot),
			},
		},
	}, nil
}

func convertInternalBeaconBlockBellatrix(b *eth.BeaconBlockBellatrix) (*BeaconBlockBellatrix, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	baseFeePerGas, err := sszBytesToUint256String(b.Body.ExecutionPayload.BaseFeePerGas)
	if err != nil {
		return nil, err
	}
	transactions := make([]string, len(b.Body.ExecutionPayload.Transactions))
	for i, tx := range b.Body.ExecutionPayload.Transactions {
		transactions[i] = hexutil.Encode(tx)
	}
	return &BeaconBlockBellatrix{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BeaconBlockBodyBellatrix{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayload: &ExecutionPayload{
				ParentHash:    hexutil.Encode(b.Body.ExecutionPayload.ParentHash),
				FeeRecipient:  hexutil.Encode(b.Body.ExecutionPayload.FeeRecipient),
				StateRoot:     hexutil.Encode(b.Body.ExecutionPayload.StateRoot),
				ReceiptsRoot:  hexutil.Encode(b.Body.ExecutionPayload.ReceiptsRoot),
				LogsBloom:     hexutil.Encode(b.Body.ExecutionPayload.LogsBloom),
				PrevRandao:    hexutil.Encode(b.Body.ExecutionPayload.PrevRandao),
				BlockNumber:   fmt.Sprintf("%d", b.Body.ExecutionPayload.BlockNumber),
				GasLimit:      fmt.Sprintf("%d", b.Body.ExecutionPayload.GasLimit),
				GasUsed:       fmt.Sprintf("%d", b.Body.ExecutionPayload.GasUsed),
				Timestamp:     fmt.Sprintf("%d", b.Body.ExecutionPayload.Timestamp),
				ExtraData:     hexutil.Encode(b.Body.ExecutionPayload.ExtraData),
				BaseFeePerGas: baseFeePerGas,
				BlockHash:     hexutil.Encode(b.Body.ExecutionPayload.BlockHash),
				Transactions:  transactions,
			},
		},
	}, nil
}

func convertInternalBlindedBeaconBlockCapella(b *eth.BlindedBeaconBlockCapella) (*BlindedBeaconBlockCapella, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	baseFeePerGas, err := sszBytesToUint256String(b.Body.ExecutionPayloadHeader.BaseFeePerGas)
	if err != nil {
		return nil, err
	}
	blsChanges, err := convertInternalBlsChanges(b.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}

	return &BlindedBeaconBlockCapella{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BlindedBeaconBlockBodyCapella{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayloadHeader: &ExecutionPayloadHeaderCapella{
				ParentHash:       hexutil.Encode(b.Body.ExecutionPayloadHeader.ParentHash),
				FeeRecipient:     hexutil.Encode(b.Body.ExecutionPayloadHeader.FeeRecipient),
				StateRoot:        hexutil.Encode(b.Body.ExecutionPayloadHeader.StateRoot),
				ReceiptsRoot:     hexutil.Encode(b.Body.ExecutionPayloadHeader.ReceiptsRoot),
				LogsBloom:        hexutil.Encode(b.Body.ExecutionPayloadHeader.LogsBloom),
				PrevRandao:       hexutil.Encode(b.Body.ExecutionPayloadHeader.PrevRandao),
				BlockNumber:      fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.BlockNumber),
				GasLimit:         fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasLimit),
				GasUsed:          fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasUsed),
				Timestamp:        fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.Timestamp),
				ExtraData:        hexutil.Encode(b.Body.ExecutionPayloadHeader.ExtraData),
				BaseFeePerGas:    baseFeePerGas,
				BlockHash:        hexutil.Encode(b.Body.ExecutionPayloadHeader.BlockHash),
				TransactionsRoot: hexutil.Encode(b.Body.ExecutionPayloadHeader.TransactionsRoot),
				WithdrawalsRoot:  hexutil.Encode(b.Body.ExecutionPayloadHeader.WithdrawalsRoot), // new in capella
			},
			BlsToExecutionChanges: blsChanges, // new in capella
		},
	}, nil
}

func convertInternalBeaconBlockCapella(b *eth.BeaconBlockCapella) (*BeaconBlockCapella, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	baseFeePerGas, err := sszBytesToUint256String(b.Body.ExecutionPayload.BaseFeePerGas)
	if err != nil {
		return nil, err
	}
	transactions := make([]string, len(b.Body.ExecutionPayload.Transactions))
	for i, tx := range b.Body.ExecutionPayload.Transactions {
		transactions[i] = hexutil.Encode(tx)
	}
	withdrawals := make([]*Withdrawal, len(b.Body.ExecutionPayload.Withdrawals))
	for i, w := range b.Body.ExecutionPayload.Withdrawals {
		withdrawals[i] = &Withdrawal{
			WithdrawalIndex:  fmt.Sprintf("%d", w.Index),
			ValidatorIndex:   fmt.Sprintf("%d", w.ValidatorIndex),
			ExecutionAddress: hexutil.Encode(w.Address),
			Amount:           fmt.Sprintf("%d", w.Amount),
		}
	}
	blsChanges, err := convertInternalBlsChanges(b.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}
	return &BeaconBlockCapella{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BeaconBlockBodyCapella{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayload: &ExecutionPayloadCapella{
				ParentHash:    hexutil.Encode(b.Body.ExecutionPayload.ParentHash),
				FeeRecipient:  hexutil.Encode(b.Body.ExecutionPayload.FeeRecipient),
				StateRoot:     hexutil.Encode(b.Body.ExecutionPayload.StateRoot),
				ReceiptsRoot:  hexutil.Encode(b.Body.ExecutionPayload.ReceiptsRoot),
				LogsBloom:     hexutil.Encode(b.Body.ExecutionPayload.LogsBloom),
				PrevRandao:    hexutil.Encode(b.Body.ExecutionPayload.PrevRandao),
				BlockNumber:   fmt.Sprintf("%d", b.Body.ExecutionPayload.BlockNumber),
				GasLimit:      fmt.Sprintf("%d", b.Body.ExecutionPayload.GasLimit),
				GasUsed:       fmt.Sprintf("%d", b.Body.ExecutionPayload.GasUsed),
				Timestamp:     fmt.Sprintf("%d", b.Body.ExecutionPayload.Timestamp),
				ExtraData:     hexutil.Encode(b.Body.ExecutionPayload.ExtraData),
				BaseFeePerGas: baseFeePerGas,
				BlockHash:     hexutil.Encode(b.Body.ExecutionPayload.BlockHash),
				Transactions:  transactions,
				Withdrawals:   withdrawals, // new in capella
			},
			BlsToExecutionChanges: blsChanges, // new in capella
		},
	}, nil
}

func convertInternalBlindedBeaconBlockContentsDeneb(b *eth.BlindedBeaconBlockAndBlobsDeneb) (*BlindedBeaconBlockContentsDeneb, error) {
	if b == nil || b.Block == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	var blindedBlobSidecars []*BlindedBlobSidecar
	if len(b.Blobs) != 0 {
		blindedBlobSidecars = make([]*BlindedBlobSidecar, len(b.Blobs))
		for i, s := range b.Blobs {
			signedBlob, err := convertInternalToBlindedBlobSidecar(s)
			if err != nil {
				return nil, err
			}
			blindedBlobSidecars[i] = signedBlob
		}
	}
	blindedBlock, err := convertInternalToBlindedDenebBlock(b.Block)
	if err != nil {
		return nil, err
	}
	return &BlindedBeaconBlockContentsDeneb{
		BlindedBlock:        blindedBlock,
		BlindedBlobSidecars: blindedBlobSidecars,
	}, nil
}

func convertInternalBeaconBlockContentsDeneb(b *eth.BeaconBlockAndBlobsDeneb) (*BeaconBlockContentsDeneb, error) {
	if b == nil || b.Block == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	var blobSidecars []*BlobSidecar
	if len(b.Blobs) != 0 {
		blobSidecars = make([]*BlobSidecar, len(b.Blobs))
		for i, s := range b.Blobs {
			blob, err := convertInternalToBlobSidecar(s)
			if err != nil {
				return nil, err
			}
			blobSidecars[i] = blob
		}
	}
	block, err := convertInternalToDenebBlock(b.Block)
	if err != nil {
		return nil, err
	}
	return &BeaconBlockContentsDeneb{
		Block:        block,
		BlobSidecars: blobSidecars,
	}, nil
}

func convertToSignedDenebBlock(signedBlock *SignedBeaconBlockDeneb) (*eth.SignedBeaconBlockDeneb, error) {
	sig, err := hexutil.Decode(signedBlock.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock .Signature")
	}
	block, err := convertToDenebBlock(signedBlock.Message)
	if err != nil {
		return nil, err
	}
	return &eth.SignedBeaconBlockDeneb{
		Block:     block,
		Signature: sig,
	}, nil
}

func convertToDenebBlock(signedBlock *BeaconBlockDeneb) (*eth.BeaconBlockDeneb, error) {
	slot, err := strconv.ParseUint(signedBlock.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Slot")
	}
	proposerIndex, err := strconv.ParseUint(signedBlock.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(signedBlock.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(signedBlock.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(signedBlock.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(signedBlock.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(signedBlock.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(signedBlock.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(signedBlock.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(signedBlock.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(signedBlock.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(signedBlock.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(signedBlock.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(signedBlock.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := hexutil.Decode(signedBlock.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(signedBlock.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(signedBlock.Body.ExecutionPayload.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(signedBlock.Body.ExecutionPayload.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(signedBlock.Body.ExecutionPayload.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(signedBlock.Body.ExecutionPayload.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(signedBlock.Body.ExecutionPayload.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(signedBlock.Body.ExecutionPayload.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(signedBlock.Body.ExecutionPayload.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(signedBlock.Body.ExecutionPayload.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(signedBlock.Body.ExecutionPayload.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(signedBlock.Body.ExecutionPayload.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayloadHeader.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(signedBlock.Body.ExecutionPayload.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToSSZBytes(signedBlock.Body.ExecutionPayload.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(signedBlock.Body.ExecutionPayload.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.BlockHash")
	}
	txs := make([][]byte, len(signedBlock.Body.ExecutionPayload.Transactions))
	for i, tx := range signedBlock.Body.ExecutionPayload.Transactions {
		txs[i], err = hexutil.Decode(tx)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode signedBlock.Body.ExecutionPayload.Transactions[%d]", i)
		}
	}
	withdrawals := make([]*enginev1.Withdrawal, len(signedBlock.Body.ExecutionPayload.Withdrawals))
	for i, w := range signedBlock.Body.ExecutionPayload.Withdrawals {
		withdrawalIndex, err := strconv.ParseUint(w.WithdrawalIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode signedBlock.Body.ExecutionPayload.Withdrawals[%d].WithdrawalIndex", i)
		}
		validatorIndex, err := strconv.ParseUint(w.ValidatorIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode signedBlock.Body.ExecutionPayload.Withdrawals[%d].ValidatorIndex", i)
		}
		address, err := hexutil.Decode(w.ExecutionAddress)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Body.ExecutionPayload.Withdrawals[%d].ExecutionAddress", i)
		}
		amount, err := strconv.ParseUint(w.Amount, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Body.ExecutionPayload.Withdrawals[%d].Amount", i)
		}
		withdrawals[i] = &enginev1.Withdrawal{
			Index:          withdrawalIndex,
			ValidatorIndex: primitives.ValidatorIndex(validatorIndex),
			Address:        address,
			Amount:         amount,
		}
	}

	payloadBlobGasUsed, err := strconv.ParseUint(signedBlock.Body.ExecutionPayload.BlobGasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.BlobGasUsed")
	}
	payloadExcessBlobGas, err := strconv.ParseUint(signedBlock.Body.ExecutionPayload.ExcessBlobGas, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Body.ExecutionPayload.ExcessBlobGas")
	}
	blsChanges, err := convertBlsChanges(signedBlock.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}
	blobKzgCommitments := make([][]byte, len(signedBlock.Body.BlobKzgCommitments))
	for i, b := range signedBlock.Body.BlobKzgCommitments {
		kzg, err := hexutil.Decode(b)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("blob kzg commitment at index %d failed to decode", i))
		}
		blobKzgCommitments[i] = kzg
	}
	return &eth.BeaconBlockDeneb{
		Slot:          primitives.Slot(slot),
		ProposerIndex: primitives.ValidatorIndex(proposerIndex),
		ParentRoot:    parentRoot,
		StateRoot:     stateRoot,
		Body: &eth.BeaconBlockBodyDeneb{
			RandaoReveal: randaoReveal,
			Eth1Data: &eth.Eth1Data{
				DepositRoot:  depositRoot,
				DepositCount: depositCount,
				BlockHash:    blockHash,
			},
			Graffiti:          graffiti,
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &eth.SyncAggregate{
				SyncCommitteeBits:      syncCommitteeBits,
				SyncCommitteeSignature: syncCommitteeSig,
			},
			ExecutionPayload: &enginev1.ExecutionPayloadDeneb{
				ParentHash:    payloadParentHash,
				FeeRecipient:  payloadFeeRecipient,
				StateRoot:     payloadStateRoot,
				ReceiptsRoot:  payloadReceiptsRoot,
				LogsBloom:     payloadLogsBloom,
				PrevRandao:    payloadPrevRandao,
				BlockNumber:   payloadBlockNumber,
				GasLimit:      payloadGasLimit,
				GasUsed:       payloadGasUsed,
				Timestamp:     payloadTimestamp,
				ExtraData:     payloadExtraData,
				BaseFeePerGas: payloadBaseFeePerGas,
				BlockHash:     payloadBlockHash,
				Transactions:  txs,
				Withdrawals:   withdrawals,
				BlobGasUsed:   payloadBlobGasUsed,
				ExcessBlobGas: payloadExcessBlobGas,
			},
			BlsToExecutionChanges: blsChanges,
			BlobKzgCommitments:    blobKzgCommitments,
		},
	}, nil
}

func convertToSignedBlobSidecar(i int, signedBlob *SignedBlobSidecar) (*eth.SignedBlobSidecar, error) {
	blobSig, err := hexutil.Decode(signedBlob.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlob.Signature")
	}
	if signedBlob.Message == nil {
		return nil, fmt.Errorf("blobsidecar message was empty at index %d", i)
	}
	blockRoot, err := hexutil.Decode(signedBlob.Message.BlockRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.BlockRoot at index %d", i))
	}
	index, err := strconv.ParseUint(signedBlob.Message.Index, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.Index at index %d", i))
	}
	slot, err := strconv.ParseUint(signedBlob.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.Index at index %d", i))
	}
	blockParentRoot, err := hexutil.Decode(signedBlob.Message.BlockParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.BlockParentRoot at index %d", i))
	}
	proposerIndex, err := strconv.ParseUint(signedBlob.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.ProposerIndex at index %d", i))
	}
	blob, err := hexutil.Decode(signedBlob.Message.Blob)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.Blob at index %d", i))
	}
	kzgCommitment, err := hexutil.Decode(signedBlob.Message.KzgCommitment)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.KzgCommitment at index %d", i))
	}
	kzgProof, err := hexutil.Decode(signedBlob.Message.KzgProof)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.KzgProof at index %d", i))
	}
	bsc := &eth.BlobSidecar{
		BlockRoot:       blockRoot,
		Index:           index,
		Slot:            primitives.Slot(slot),
		BlockParentRoot: blockParentRoot,
		ProposerIndex:   primitives.ValidatorIndex(proposerIndex),
		Blob:            blob,
		KzgCommitment:   kzgCommitment,
		KzgProof:        kzgProof,
	}
	return &eth.SignedBlobSidecar{
		Message:   bsc,
		Signature: blobSig,
	}, nil
}

func convertToBlobSidecar(i int, b *BlobSidecar) (*eth.BlobSidecar, error) {
	if b == nil {
		return nil, fmt.Errorf("blobsidecar message was empty at index %d", i)
	}
	blockRoot, err := hexutil.Decode(b.BlockRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.BlockRoot at index %d", i))
	}
	index, err := strconv.ParseUint(b.Index, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.Index at index %d", i))
	}
	slot, err := strconv.ParseUint(b.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.Index at index %d", i))
	}
	blockParentRoot, err := hexutil.Decode(b.BlockParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.BlockParentRoot at index %d", i))
	}
	proposerIndex, err := strconv.ParseUint(b.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.ProposerIndex at index %d", i))
	}
	blob, err := hexutil.Decode(b.Blob)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.Blob at index %d", i))
	}
	kzgCommitment, err := hexutil.Decode(b.KzgCommitment)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.KzgCommitment at index %d", i))
	}
	kzgProof, err := hexutil.Decode(b.KzgProof)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.KzgProof at index %d", i))
	}
	bsc := &eth.BlobSidecar{
		BlockRoot:       blockRoot,
		Index:           index,
		Slot:            primitives.Slot(slot),
		BlockParentRoot: blockParentRoot,
		ProposerIndex:   primitives.ValidatorIndex(proposerIndex),
		Blob:            blob,
		KzgCommitment:   kzgCommitment,
		KzgProof:        kzgProof,
	}
	return bsc, nil
}

func convertToSignedBlindedDenebBlock(signedBlindedBlock *SignedBlindedBeaconBlockDeneb) (*eth.SignedBlindedBeaconBlockDeneb, error) {
	if signedBlindedBlock == nil {
		return nil, errors.New("signed blinded block is empty")
	}
	sig, err := hexutil.Decode(signedBlindedBlock.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Signature")
	}
	blindedBlock, err := convertToBlindedDenebBlock(signedBlindedBlock.Message)
	if err != nil {
		return nil, err
	}
	return &eth.SignedBlindedBeaconBlockDeneb{
		Block:     blindedBlock,
		Signature: sig,
	}, nil
}

func convertToBlindedDenebBlock(signedBlindedBlock *BlindedBeaconBlockDeneb) (*eth.BlindedBeaconBlockDeneb, error) {
	if signedBlindedBlock == nil {
		return nil, errors.New("blinded block is empty")
	}
	slot, err := strconv.ParseUint(signedBlindedBlock.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Slot")
	}
	proposerIndex, err := strconv.ParseUint(signedBlindedBlock.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(signedBlindedBlock.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(signedBlindedBlock.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(signedBlindedBlock.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(signedBlindedBlock.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(signedBlindedBlock.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(signedBlindedBlock.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(signedBlindedBlock.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(signedBlindedBlock.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(signedBlindedBlock.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(signedBlindedBlock.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(signedBlindedBlock.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(signedBlindedBlock.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := hexutil.Decode(signedBlindedBlock.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(signedBlindedBlock.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(signedBlindedBlock.Body.ExecutionPayloadHeader.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(signedBlindedBlock.Body.ExecutionPayloadHeader.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(signedBlindedBlock.Body.ExecutionPayloadHeader.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(signedBlindedBlock.Body.ExecutionPayloadHeader.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(signedBlindedBlock.Body.ExecutionPayloadHeader.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(signedBlindedBlock.Body.ExecutionPayloadHeader.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(signedBlindedBlock.Body.ExecutionPayloadHeader.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(signedBlindedBlock.Body.ExecutionPayloadHeader.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(signedBlindedBlock.Body.ExecutionPayloadHeader.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(signedBlindedBlock.Body.ExecutionPayloadHeader.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(signedBlindedBlock.Body.ExecutionPayloadHeader.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToSSZBytes(signedBlindedBlock.Body.ExecutionPayloadHeader.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(signedBlindedBlock.Body.ExecutionPayloadHeader.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.BlockHash")
	}
	payloadTxsRoot, err := hexutil.Decode(signedBlindedBlock.Body.ExecutionPayloadHeader.TransactionsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.TransactionsRoot")
	}
	payloadWithdrawalsRoot, err := hexutil.Decode(signedBlindedBlock.Body.ExecutionPayloadHeader.WithdrawalsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayloadHeader.WithdrawalsRoot")
	}

	payloadBlobGasUsed, err := strconv.ParseUint(signedBlindedBlock.Body.ExecutionPayloadHeader.BlobGasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayload.BlobGasUsed")
	}
	payloadExcessBlobGas, err := strconv.ParseUint(signedBlindedBlock.Body.ExecutionPayloadHeader.ExcessBlobGas, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Body.ExecutionPayload.ExcessBlobGas")
	}

	blsChanges, err := convertBlsChanges(signedBlindedBlock.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}

	blobKzgCommitments := make([][]byte, len(signedBlindedBlock.Body.BlobKzgCommitments))
	for i, b := range signedBlindedBlock.Body.BlobKzgCommitments {
		kzg, err := hexutil.Decode(b)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("blob kzg commitment at index %d failed to decode", i))
		}
		blobKzgCommitments[i] = kzg
	}

	return &eth.BlindedBeaconBlockDeneb{
		Slot:          primitives.Slot(slot),
		ProposerIndex: primitives.ValidatorIndex(proposerIndex),
		ParentRoot:    parentRoot,
		StateRoot:     stateRoot,
		Body: &eth.BlindedBeaconBlockBodyDeneb{
			RandaoReveal: randaoReveal,
			Eth1Data: &eth.Eth1Data{
				DepositRoot:  depositRoot,
				DepositCount: depositCount,
				BlockHash:    blockHash,
			},
			Graffiti:          graffiti,
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &eth.SyncAggregate{
				SyncCommitteeBits:      syncCommitteeBits,
				SyncCommitteeSignature: syncCommitteeSig,
			},
			ExecutionPayloadHeader: &enginev1.ExecutionPayloadHeaderDeneb{
				ParentHash:       payloadParentHash,
				FeeRecipient:     payloadFeeRecipient,
				StateRoot:        payloadStateRoot,
				ReceiptsRoot:     payloadReceiptsRoot,
				LogsBloom:        payloadLogsBloom,
				PrevRandao:       payloadPrevRandao,
				BlockNumber:      payloadBlockNumber,
				GasLimit:         payloadGasLimit,
				GasUsed:          payloadGasUsed,
				Timestamp:        payloadTimestamp,
				ExtraData:        payloadExtraData,
				BaseFeePerGas:    payloadBaseFeePerGas,
				BlockHash:        payloadBlockHash,
				TransactionsRoot: payloadTxsRoot,
				WithdrawalsRoot:  payloadWithdrawalsRoot,
				BlobGasUsed:      payloadBlobGasUsed,
				ExcessBlobGas:    payloadExcessBlobGas,
			},
			BlsToExecutionChanges: blsChanges,
			BlobKzgCommitments:    blobKzgCommitments,
		},
	}, nil
}

func convertInternalToBlindedDenebBlock(b *eth.BlindedBeaconBlockDeneb) (*BlindedBeaconBlockDeneb, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	baseFeePerGas, err := sszBytesToUint256String(b.Body.ExecutionPayloadHeader.BaseFeePerGas)
	if err != nil {
		return nil, err
	}
	blsChanges, err := convertInternalBlsChanges(b.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}

	blobKzgCommitments := make([]string, len(b.Body.BlobKzgCommitments))
	for i := range b.Body.BlobKzgCommitments {
		blobKzgCommitments[i] = hexutil.Encode(b.Body.BlobKzgCommitments[i])
	}

	return &BlindedBeaconBlockDeneb{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BlindedBeaconBlockBodyDeneb{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayloadHeader: &ExecutionPayloadHeaderDeneb{
				ParentHash:       hexutil.Encode(b.Body.ExecutionPayloadHeader.ParentHash),
				FeeRecipient:     hexutil.Encode(b.Body.ExecutionPayloadHeader.FeeRecipient),
				StateRoot:        hexutil.Encode(b.Body.ExecutionPayloadHeader.StateRoot),
				ReceiptsRoot:     hexutil.Encode(b.Body.ExecutionPayloadHeader.ReceiptsRoot),
				LogsBloom:        hexutil.Encode(b.Body.ExecutionPayloadHeader.LogsBloom),
				PrevRandao:       hexutil.Encode(b.Body.ExecutionPayloadHeader.PrevRandao),
				BlockNumber:      fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.BlockNumber),
				GasLimit:         fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasLimit),
				GasUsed:          fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasUsed),
				Timestamp:        fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.Timestamp),
				ExtraData:        hexutil.Encode(b.Body.ExecutionPayloadHeader.ExtraData),
				BaseFeePerGas:    baseFeePerGas,
				BlockHash:        hexutil.Encode(b.Body.ExecutionPayloadHeader.BlockHash),
				TransactionsRoot: hexutil.Encode(b.Body.ExecutionPayloadHeader.TransactionsRoot),
				WithdrawalsRoot:  hexutil.Encode(b.Body.ExecutionPayloadHeader.WithdrawalsRoot),
				BlobGasUsed:      fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.BlobGasUsed),   // new in deneb TODO: rename to blob
				ExcessBlobGas:    fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.ExcessBlobGas), // new in deneb TODO: rename to blob
			},
			BlsToExecutionChanges: blsChanges,         // new in capella
			BlobKzgCommitments:    blobKzgCommitments, // new in deneb
		},
	}, nil
}

func convertInternalToDenebBlock(b *eth.BeaconBlockDeneb) (*BeaconBlockDeneb, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	baseFeePerGas, err := sszBytesToUint256String(b.Body.ExecutionPayload.BaseFeePerGas)
	if err != nil {
		return nil, err
	}
	transactions := make([]string, len(b.Body.ExecutionPayload.Transactions))
	for i, tx := range b.Body.ExecutionPayload.Transactions {
		transactions[i] = hexutil.Encode(tx)
	}
	withdrawals := make([]*Withdrawal, len(b.Body.ExecutionPayload.Withdrawals))
	for i, w := range b.Body.ExecutionPayload.Withdrawals {
		withdrawals[i] = &Withdrawal{
			WithdrawalIndex:  fmt.Sprintf("%d", w.Index),
			ValidatorIndex:   fmt.Sprintf("%d", w.ValidatorIndex),
			ExecutionAddress: hexutil.Encode(w.Address),
			Amount:           fmt.Sprintf("%d", w.Amount),
		}
	}
	blsChanges, err := convertInternalBlsChanges(b.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}
	blobKzgCommitments := make([]string, len(b.Body.BlobKzgCommitments))
	for i := range b.Body.BlobKzgCommitments {
		blobKzgCommitments[i] = hexutil.Encode(b.Body.BlobKzgCommitments[i])
	}

	return &BeaconBlockDeneb{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BeaconBlockBodyDeneb{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayload: &ExecutionPayloadDeneb{
				ParentHash:    hexutil.Encode(b.Body.ExecutionPayload.ParentHash),
				FeeRecipient:  hexutil.Encode(b.Body.ExecutionPayload.FeeRecipient),
				StateRoot:     hexutil.Encode(b.Body.ExecutionPayload.StateRoot),
				ReceiptsRoot:  hexutil.Encode(b.Body.ExecutionPayload.ReceiptsRoot),
				LogsBloom:     hexutil.Encode(b.Body.ExecutionPayload.LogsBloom),
				PrevRandao:    hexutil.Encode(b.Body.ExecutionPayload.PrevRandao),
				BlockNumber:   fmt.Sprintf("%d", b.Body.ExecutionPayload.BlockNumber),
				GasLimit:      fmt.Sprintf("%d", b.Body.ExecutionPayload.GasLimit),
				GasUsed:       fmt.Sprintf("%d", b.Body.ExecutionPayload.GasUsed),
				Timestamp:     fmt.Sprintf("%d", b.Body.ExecutionPayload.Timestamp),
				ExtraData:     hexutil.Encode(b.Body.ExecutionPayload.ExtraData),
				BaseFeePerGas: baseFeePerGas,
				BlockHash:     hexutil.Encode(b.Body.ExecutionPayload.BlockHash),
				Transactions:  transactions,
				Withdrawals:   withdrawals,
				BlobGasUsed:   fmt.Sprintf("%d", b.Body.ExecutionPayload.BlobGasUsed),   // new in deneb TODO: rename to blob
				ExcessBlobGas: fmt.Sprintf("%d", b.Body.ExecutionPayload.ExcessBlobGas), // new in deneb TODO: rename to blob
			},
			BlsToExecutionChanges: blsChanges,         // new in capella
			BlobKzgCommitments:    blobKzgCommitments, // new in deneb
		},
	}, nil
}

func convertToSignedBlindedBlobSidecar(i int, signedBlob *SignedBlindedBlobSidecar) (*eth.SignedBlindedBlobSidecar, error) {
	blobSig, err := hexutil.Decode(signedBlob.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlob.Signature")
	}
	if signedBlob.Message == nil {
		return nil, fmt.Errorf("blobsidecar message was empty at index %d", i)
	}
	blockRoot, err := hexutil.Decode(signedBlob.Message.BlockRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.BlockRoot at index %d", i))
	}
	index, err := strconv.ParseUint(signedBlob.Message.Index, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.Index at index %d", i))
	}
	denebSlot, err := strconv.ParseUint(signedBlob.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.Index at index %d", i))
	}
	blockParentRoot, err := hexutil.Decode(signedBlob.Message.BlockParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.BlockParentRoot at index %d", i))
	}
	proposerIndex, err := strconv.ParseUint(signedBlob.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.ProposerIndex at index %d", i))
	}
	blobRoot, err := hexutil.Decode(signedBlob.Message.BlobRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.BlobRoot at index %d", i))
	}
	kzgCommitment, err := hexutil.Decode(signedBlob.Message.KzgCommitment)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.KzgCommitment at index %d", i))
	}
	kzgProof, err := hexutil.Decode(signedBlob.Message.KzgProof)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.KzgProof at index %d", i))
	}
	bsc := &eth.BlindedBlobSidecar{
		BlockRoot:       blockRoot,
		Index:           index,
		Slot:            primitives.Slot(denebSlot),
		BlockParentRoot: blockParentRoot,
		ProposerIndex:   primitives.ValidatorIndex(proposerIndex),
		BlobRoot:        blobRoot,
		KzgCommitment:   kzgCommitment,
		KzgProof:        kzgProof,
	}
	return &eth.SignedBlindedBlobSidecar{
		Message:   bsc,
		Signature: blobSig,
	}, nil
}

func convertToBlindedBlobSidecar(i int, b *BlindedBlobSidecar) (*eth.BlindedBlobSidecar, error) {
	if b == nil {
		return nil, fmt.Errorf("blobsidecar message was empty at index %d", i)
	}
	blockRoot, err := hexutil.Decode(b.BlockRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.BlockRoot at index %d", i))
	}
	index, err := strconv.ParseUint(b.Index, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.Index at index %d", i))
	}
	denebSlot, err := strconv.ParseUint(b.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.Index at index %d", i))
	}
	blockParentRoot, err := hexutil.Decode(b.BlockParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.BlockParentRoot at index %d", i))
	}
	proposerIndex, err := strconv.ParseUint(b.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.ProposerIndex at index %d", i))
	}
	blobRoot, err := hexutil.Decode(b.BlobRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.BlobRoot at index %d", i))
	}
	kzgCommitment, err := hexutil.Decode(b.KzgCommitment)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.KzgCommitment at index %d", i))
	}
	kzgProof, err := hexutil.Decode(b.KzgProof)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode b.KzgProof at index %d", i))
	}
	bsc := &eth.BlindedBlobSidecar{
		BlockRoot:       blockRoot,
		Index:           index,
		Slot:            primitives.Slot(denebSlot),
		BlockParentRoot: blockParentRoot,
		ProposerIndex:   primitives.ValidatorIndex(proposerIndex),
		BlobRoot:        blobRoot,
		KzgCommitment:   kzgCommitment,
		KzgProof:        kzgProof,
	}
	return bsc, nil
}

func convertInternalToBlindedBlobSidecar(b *eth.BlindedBlobSidecar) (*BlindedBlobSidecar, error) {
	if b == nil {
		return nil, errors.New("BlindedBlobSidecar is empty, nothing to convert.")
	}
	return &BlindedBlobSidecar{
		BlockRoot:       hexutil.Encode(b.BlockRoot),
		Index:           fmt.Sprintf("%d", b.Index),
		Slot:            fmt.Sprintf("%d", b.Slot),
		BlockParentRoot: hexutil.Encode(b.BlockParentRoot),
		ProposerIndex:   fmt.Sprintf("%d", b.ProposerIndex),
		BlobRoot:        hexutil.Encode(b.BlobRoot),
		KzgCommitment:   hexutil.Encode(b.KzgCommitment),
		KzgProof:        hexutil.Encode(b.KzgProof),
	}, nil
}

func convertInternalToBlobSidecar(b *eth.BlobSidecar) (*BlobSidecar, error) {
	if b == nil {
		return nil, errors.New("BlobSidecar is empty, nothing to convert.")
	}
	return &BlobSidecar{
		BlockRoot:       hexutil.Encode(b.BlockRoot),
		Index:           fmt.Sprintf("%d", b.Index),
		Slot:            fmt.Sprintf("%d", b.Slot),
		BlockParentRoot: hexutil.Encode(b.BlockParentRoot),
		ProposerIndex:   fmt.Sprintf("%d", b.ProposerIndex),
		Blob:            hexutil.Encode(b.Blob),
		KzgCommitment:   hexutil.Encode(b.KzgCommitment),
		KzgProof:        hexutil.Encode(b.KzgProof),
	}, nil
}

func convertProposerSlashings(src []*ProposerSlashing) ([]*eth.ProposerSlashing, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.ProposerSlashings")
	}
	proposerSlashings := make([]*eth.ProposerSlashing, len(src))
	for i, s := range src {
		h1Sig, err := hexutil.Decode(s.SignedHeader1.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Signature", i)
		}
		h1Slot, err := strconv.ParseUint(s.SignedHeader1.Message.Slot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Message.Slot", i)
		}
		h1ProposerIndex, err := strconv.ParseUint(s.SignedHeader1.Message.ProposerIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Message.ProposerIndex", i)
		}
		h1ParentRoot, err := hexutil.Decode(s.SignedHeader1.Message.ParentRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Message.ParentRoot", i)
		}
		h1StateRoot, err := hexutil.Decode(s.SignedHeader1.Message.StateRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Message.StateRoot", i)
		}
		h1BodyRoot, err := hexutil.Decode(s.SignedHeader1.Message.BodyRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Message.BodyRoot", i)
		}
		h2Sig, err := hexutil.Decode(s.SignedHeader2.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Signature", i)
		}
		h2Slot, err := strconv.ParseUint(s.SignedHeader2.Message.Slot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Message.Slot", i)
		}
		h2ProposerIndex, err := strconv.ParseUint(s.SignedHeader2.Message.ProposerIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Message.ProposerIndex", i)
		}
		h2ParentRoot, err := hexutil.Decode(s.SignedHeader2.Message.ParentRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Message.ParentRoot", i)
		}
		h2StateRoot, err := hexutil.Decode(s.SignedHeader2.Message.StateRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Message.StateRoot", i)
		}
		h2BodyRoot, err := hexutil.Decode(s.SignedHeader2.Message.BodyRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Message.BodyRoot", i)
		}
		proposerSlashings[i] = &eth.ProposerSlashing{
			Header_1: &eth.SignedBeaconBlockHeader{
				Header: &eth.BeaconBlockHeader{
					Slot:          primitives.Slot(h1Slot),
					ProposerIndex: primitives.ValidatorIndex(h1ProposerIndex),
					ParentRoot:    h1ParentRoot,
					StateRoot:     h1StateRoot,
					BodyRoot:      h1BodyRoot,
				},
				Signature: h1Sig,
			},
			Header_2: &eth.SignedBeaconBlockHeader{
				Header: &eth.BeaconBlockHeader{
					Slot:          primitives.Slot(h2Slot),
					ProposerIndex: primitives.ValidatorIndex(h2ProposerIndex),
					ParentRoot:    h2ParentRoot,
					StateRoot:     h2StateRoot,
					BodyRoot:      h2BodyRoot,
				},
				Signature: h2Sig,
			},
		}
	}
	return proposerSlashings, nil
}

func convertInternalProposerSlashings(src []*eth.ProposerSlashing) ([]*ProposerSlashing, error) {
	if src == nil {
		return nil, errors.New("proposer slashings are empty, nothing to convert.")
	}
	proposerSlashings := make([]*ProposerSlashing, len(src))
	for i, s := range src {
		proposerSlashings[i] = &ProposerSlashing{
			SignedHeader1: &SignedBeaconBlockHeader{
				Message: &BeaconBlockHeader{
					Slot:          fmt.Sprintf("%d", s.Header_1.Header.Slot),
					ProposerIndex: fmt.Sprintf("%d", s.Header_1.Header.ProposerIndex),
					ParentRoot:    hexutil.Encode(s.Header_1.Header.ParentRoot),
					StateRoot:     hexutil.Encode(s.Header_1.Header.StateRoot),
					BodyRoot:      hexutil.Encode(s.Header_1.Header.BodyRoot),
				},
				Signature: hexutil.Encode(s.Header_1.Signature),
			},
			SignedHeader2: &SignedBeaconBlockHeader{
				Message: &BeaconBlockHeader{
					Slot:          fmt.Sprintf("%d", s.Header_2.Header.Slot),
					ProposerIndex: fmt.Sprintf("%d", s.Header_2.Header.ProposerIndex),
					ParentRoot:    hexutil.Encode(s.Header_2.Header.ParentRoot),
					StateRoot:     hexutil.Encode(s.Header_2.Header.StateRoot),
					BodyRoot:      hexutil.Encode(s.Header_2.Header.BodyRoot),
				},
				Signature: hexutil.Encode(s.Header_2.Signature),
			},
		}
	}
	return proposerSlashings, nil
}

func convertAttesterSlashings(src []*AttesterSlashing) ([]*eth.AttesterSlashing, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.AttesterSlashings")
	}
	attesterSlashings := make([]*eth.AttesterSlashing, len(src))
	for i, s := range src {
		a1Sig, err := hexutil.Decode(s.Attestation1.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Signature", i)
		}
		a1AttestingIndices := make([]uint64, len(s.Attestation1.AttestingIndices))
		for j, ix := range s.Attestation1.AttestingIndices {
			attestingIndex, err := strconv.ParseUint(ix, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.AttestingIndices[%d]", i, j)
			}
			a1AttestingIndices[j] = attestingIndex
		}
		a1Slot, err := strconv.ParseUint(s.Attestation1.Data.Slot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Slot", i)
		}
		a1CommitteeIndex, err := strconv.ParseUint(s.Attestation1.Data.Index, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Index", i)
		}
		a1BeaconBlockRoot, err := hexutil.Decode(s.Attestation1.Data.BeaconBlockRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.BeaconBlockRoot", i)
		}
		a1SourceEpoch, err := strconv.ParseUint(s.Attestation1.Data.Source.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Source.Epoch", i)
		}
		a1SourceRoot, err := hexutil.Decode(s.Attestation1.Data.Source.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Source.Root", i)
		}
		a1TargetEpoch, err := strconv.ParseUint(s.Attestation1.Data.Target.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Target.Epoch", i)
		}
		a1TargetRoot, err := hexutil.Decode(s.Attestation1.Data.Target.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Target.Root", i)
		}
		a2Sig, err := hexutil.Decode(s.Attestation2.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Signature", i)
		}
		a2AttestingIndices := make([]uint64, len(s.Attestation2.AttestingIndices))
		for j, ix := range s.Attestation2.AttestingIndices {
			attestingIndex, err := strconv.ParseUint(ix, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.AttestingIndices[%d]", i, j)
			}
			a2AttestingIndices[j] = attestingIndex
		}
		a2Slot, err := strconv.ParseUint(s.Attestation2.Data.Slot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Slot", i)
		}
		a2CommitteeIndex, err := strconv.ParseUint(s.Attestation2.Data.Index, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Index", i)
		}
		a2BeaconBlockRoot, err := hexutil.Decode(s.Attestation2.Data.BeaconBlockRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.BeaconBlockRoot", i)
		}
		a2SourceEpoch, err := strconv.ParseUint(s.Attestation2.Data.Source.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Source.Epoch", i)
		}
		a2SourceRoot, err := hexutil.Decode(s.Attestation2.Data.Source.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Source.Root", i)
		}
		a2TargetEpoch, err := strconv.ParseUint(s.Attestation2.Data.Target.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Target.Epoch", i)
		}
		a2TargetRoot, err := hexutil.Decode(s.Attestation2.Data.Target.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Target.Root", i)
		}
		attesterSlashings[i] = &eth.AttesterSlashing{
			Attestation_1: &eth.IndexedAttestation{
				AttestingIndices: a1AttestingIndices,
				Data: &eth.AttestationData{
					Slot:            primitives.Slot(a1Slot),
					CommitteeIndex:  primitives.CommitteeIndex(a1CommitteeIndex),
					BeaconBlockRoot: a1BeaconBlockRoot,
					Source: &eth.Checkpoint{
						Epoch: primitives.Epoch(a1SourceEpoch),
						Root:  a1SourceRoot,
					},
					Target: &eth.Checkpoint{
						Epoch: primitives.Epoch(a1TargetEpoch),
						Root:  a1TargetRoot,
					},
				},
				Signature: a1Sig,
			},
			Attestation_2: &eth.IndexedAttestation{
				AttestingIndices: a2AttestingIndices,
				Data: &eth.AttestationData{
					Slot:            primitives.Slot(a2Slot),
					CommitteeIndex:  primitives.CommitteeIndex(a2CommitteeIndex),
					BeaconBlockRoot: a2BeaconBlockRoot,
					Source: &eth.Checkpoint{
						Epoch: primitives.Epoch(a2SourceEpoch),
						Root:  a2SourceRoot,
					},
					Target: &eth.Checkpoint{
						Epoch: primitives.Epoch(a2TargetEpoch),
						Root:  a2TargetRoot,
					},
				},
				Signature: a2Sig,
			},
		}
	}
	return attesterSlashings, nil
}

func convertInternalAttesterSlashings(src []*eth.AttesterSlashing) ([]*AttesterSlashing, error) {
	if src == nil {
		return nil, errors.New("AttesterSlashings is empty, nothing to convert.")
	}
	attesterSlashings := make([]*AttesterSlashing, len(src))
	for i, s := range src {
		a1AttestingIndices := make([]string, len(s.Attestation_1.AttestingIndices))
		for j, ix := range s.Attestation_1.AttestingIndices {
			a1AttestingIndices[j] = fmt.Sprintf("%d", ix)
		}
		a2AttestingIndices := make([]string, len(s.Attestation_2.AttestingIndices))
		for j, ix := range s.Attestation_2.AttestingIndices {
			a2AttestingIndices[j] = fmt.Sprintf("%d", ix)
		}
		attesterSlashings[i] = &AttesterSlashing{
			Attestation1: &IndexedAttestation{
				AttestingIndices: a1AttestingIndices,
				Data: &AttestationData{
					Slot:            fmt.Sprintf("%d", s.Attestation_1.Data.Slot),
					Index:           fmt.Sprintf("%d", s.Attestation_1.Data.CommitteeIndex),
					BeaconBlockRoot: hexutil.Encode(s.Attestation_1.Data.BeaconBlockRoot),
					Source: &Checkpoint{
						Epoch: fmt.Sprintf("%d", s.Attestation_1.Data.Source.Epoch),
						Root:  hexutil.Encode(s.Attestation_1.Data.Source.Root),
					},
					Target: &Checkpoint{
						Epoch: fmt.Sprintf("%d", s.Attestation_1.Data.Target.Epoch),
						Root:  hexutil.Encode(s.Attestation_1.Data.Target.Root),
					},
				},
				Signature: hexutil.Encode(s.Attestation_1.Signature),
			},
			Attestation2: &IndexedAttestation{
				AttestingIndices: a2AttestingIndices,
				Data: &AttestationData{
					Slot:            fmt.Sprintf("%d", s.Attestation_2.Data.Slot),
					Index:           fmt.Sprintf("%d", s.Attestation_2.Data.CommitteeIndex),
					BeaconBlockRoot: hexutil.Encode(s.Attestation_2.Data.BeaconBlockRoot),
					Source: &Checkpoint{
						Epoch: fmt.Sprintf("%d", s.Attestation_2.Data.Source.Epoch),
						Root:  hexutil.Encode(s.Attestation_2.Data.Source.Root),
					},
					Target: &Checkpoint{
						Epoch: fmt.Sprintf("%d", s.Attestation_2.Data.Target.Epoch),
						Root:  hexutil.Encode(s.Attestation_2.Data.Target.Root),
					},
				},
				Signature: hexutil.Encode(s.Attestation_2.Signature),
			},
		}
	}
	return attesterSlashings, nil
}

func convertAtts(src []*Attestation) ([]*eth.Attestation, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.Attestations")
	}
	atts := make([]*eth.Attestation, len(src))
	for i, a := range src {
		sig, err := hexutil.Decode(a.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Signature", i)
		}
		slot, err := strconv.ParseUint(a.Data.Slot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Slot", i)
		}
		committeeIndex, err := strconv.ParseUint(a.Data.Index, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Index", i)
		}
		beaconBlockRoot, err := hexutil.Decode(a.Data.BeaconBlockRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.BeaconBlockRoot", i)
		}
		sourceEpoch, err := strconv.ParseUint(a.Data.Source.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Source.Epoch", i)
		}
		sourceRoot, err := hexutil.Decode(a.Data.Source.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Source.Root", i)
		}
		targetEpoch, err := strconv.ParseUint(a.Data.Target.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Target.Epoch", i)
		}
		targetRoot, err := hexutil.Decode(a.Data.Target.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Target.Root", i)
		}
		aggregateBits, err := hexutil.Decode(a.AggregationBits)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].AggregationBits", i)
		}
		atts[i] = &eth.Attestation{
			AggregationBits: aggregateBits,
			Data: &eth.AttestationData{
				Slot:            primitives.Slot(slot),
				CommitteeIndex:  primitives.CommitteeIndex(committeeIndex),
				BeaconBlockRoot: beaconBlockRoot,
				Source: &eth.Checkpoint{
					Epoch: primitives.Epoch(sourceEpoch),
					Root:  sourceRoot,
				},
				Target: &eth.Checkpoint{
					Epoch: primitives.Epoch(targetEpoch),
					Root:  targetRoot,
				},
			},
			Signature: sig,
		}
	}
	return atts, nil
}

func convertInternalAtts(src []*eth.Attestation) ([]*Attestation, error) {
	if src == nil {
		return nil, errors.New("Attestations are empty, nothing to convert.")
	}
	atts := make([]*Attestation, len(src))
	for i, a := range src {
		atts[i] = &Attestation{
			AggregationBits: hexutil.Encode(a.AggregationBits),
			Data: &AttestationData{
				Slot:            fmt.Sprintf("%d", a.Data.Slot),
				Index:           fmt.Sprintf("%d", a.Data.CommitteeIndex),
				BeaconBlockRoot: hexutil.Encode(a.Data.BeaconBlockRoot),
				Source: &Checkpoint{
					Epoch: fmt.Sprintf("%d", a.Data.Source.Epoch),
					Root:  hexutil.Encode(a.Data.Source.Root),
				},
				Target: &Checkpoint{
					Epoch: fmt.Sprintf("%d", a.Data.Target.Epoch),
					Root:  hexutil.Encode(a.Data.Target.Root),
				},
			},
			Signature: hexutil.Encode(a.Signature),
		}
	}
	return atts, nil
}

func convertDeposits(src []*Deposit) ([]*eth.Deposit, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.Deposits")
	}
	deposits := make([]*eth.Deposit, len(src))
	for i, d := range src {
		proof := make([][]byte, len(d.Proof))
		for j, p := range d.Proof {
			var err error
			proof[j], err = hexutil.Decode(p)
			if err != nil {
				return nil, errors.Wrapf(err, "could not decode b.Message.Body.Deposits[%d].Proof[%d]", i, j)
			}
		}
		pubkey, err := hexutil.Decode(d.Data.Pubkey)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Deposits[%d].Pubkey", i)
		}
		withdrawalCreds, err := hexutil.Decode(d.Data.WithdrawalCredentials)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Deposits[%d].WithdrawalCredentials", i)
		}
		amount, err := strconv.ParseUint(d.Data.Amount, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Deposits[%d].Amount", i)
		}
		sig, err := hexutil.Decode(d.Data.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Deposits[%d].Signature", i)
		}
		deposits[i] = &eth.Deposit{
			Proof: proof,
			Data: &eth.Deposit_Data{
				PublicKey:             pubkey,
				WithdrawalCredentials: withdrawalCreds,
				Amount:                amount,
				Signature:             sig,
			},
		}
	}
	return deposits, nil
}

func convertInternalDeposits(src []*eth.Deposit) ([]*Deposit, error) {
	if src == nil {
		return nil, errors.New("deposits are empty, nothing to convert.")
	}
	deposits := make([]*Deposit, len(src))
	for i, d := range src {
		proof := make([]string, len(d.Proof))
		for j, p := range d.Proof {
			proof[j] = hexutil.Encode(p)
		}
		deposits[i] = &Deposit{
			Proof: proof,
			Data: &DepositData{
				Pubkey:                hexutil.Encode(d.Data.PublicKey),
				WithdrawalCredentials: hexutil.Encode(d.Data.WithdrawalCredentials),
				Amount:                fmt.Sprintf("%d", d.Data.Amount),
				Signature:             hexutil.Encode(d.Data.Signature),
			},
		}
	}
	return deposits, nil
}

func convertExits(src []*SignedVoluntaryExit) ([]*eth.SignedVoluntaryExit, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.VoluntaryExits")
	}
	exits := make([]*eth.SignedVoluntaryExit, len(src))
	for i, e := range src {
		sig, err := hexutil.Decode(e.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.VoluntaryExits[%d].Signature", i)
		}
		epoch, err := strconv.ParseUint(e.Message.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.VoluntaryExits[%d].Epoch", i)
		}
		validatorIndex, err := strconv.ParseUint(e.Message.ValidatorIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.VoluntaryExits[%d].ValidatorIndex", i)
		}
		exits[i] = &eth.SignedVoluntaryExit{
			Exit: &eth.VoluntaryExit{
				Epoch:          primitives.Epoch(epoch),
				ValidatorIndex: primitives.ValidatorIndex(validatorIndex),
			},
			Signature: sig,
		}
	}
	return exits, nil
}

func convertInternalExits(src []*eth.SignedVoluntaryExit) ([]*SignedVoluntaryExit, error) {
	if src == nil {
		return nil, errors.New("VoluntaryExits are empty, nothing to convert.")
	}
	exits := make([]*SignedVoluntaryExit, len(src))
	for i, e := range src {
		exits[i] = &SignedVoluntaryExit{
			Message: &VoluntaryExit{
				Epoch:          fmt.Sprintf("%d", e.Exit.Epoch),
				ValidatorIndex: fmt.Sprintf("%d", e.Exit.ValidatorIndex),
			},
			Signature: hexutil.Encode(e.Signature),
		}
	}
	return exits, nil
}

func convertBlsChanges(src []*SignedBlsToExecutionChange) ([]*eth.SignedBLSToExecutionChange, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.BlsToExecutionChanges")
	}
	changes := make([]*eth.SignedBLSToExecutionChange, len(src))
	for i, ch := range src {
		sig, err := hexutil.Decode(ch.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.BlsToExecutionChanges[%d].Signature", i)
		}
		index, err := strconv.ParseUint(ch.Message.ValidatorIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.BlsToExecutionChanges[%d].Message.ValidatorIndex", i)
		}
		pubkey, err := hexutil.Decode(ch.Message.FromBlsPubkey)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.BlsToExecutionChanges[%d].Message.FromBlsPubkey", i)
		}
		address, err := hexutil.Decode(ch.Message.ToExecutionAddress)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.BlsToExecutionChanges[%d].Message.ToExecutionAddress", i)
		}
		changes[i] = &eth.SignedBLSToExecutionChange{
			Message: &eth.BLSToExecutionChange{
				ValidatorIndex:     primitives.ValidatorIndex(index),
				FromBlsPubkey:      pubkey,
				ToExecutionAddress: address,
			},
			Signature: sig,
		}
	}
	return changes, nil
}

func convertInternalBlsChanges(src []*eth.SignedBLSToExecutionChange) ([]*SignedBlsToExecutionChange, error) {
	if src == nil {
		return nil, errors.New("BlsToExecutionChanges are empty, nothing to convert.")
	}
	changes := make([]*SignedBlsToExecutionChange, len(src))
	for i, ch := range src {
		changes[i] = &SignedBlsToExecutionChange{
			Message: &BlsToExecutionChange{
				ValidatorIndex:     fmt.Sprintf("%d", ch.Message.ValidatorIndex),
				FromBlsPubkey:      hexutil.Encode(ch.Message.FromBlsPubkey),
				ToExecutionAddress: hexutil.Encode(ch.Message.ToExecutionAddress),
			},
			Signature: hexutil.Encode(ch.Signature),
		}
	}
	return changes, nil
}

func uint256ToSSZBytes(num string) ([]byte, error) {
	uint256, ok := new(big.Int).SetString(num, 10)
	if !ok {
		return nil, errors.New("could not parse Uint256")
	}
	if !math.IsValidUint256(uint256) {
		return nil, errors.New("is")
	}
	return bytesutil2.PadTo(bytesutil2.ReverseByteOrder(uint256.Bytes()), 32), nil
}

func sszBytesToUint256String(b []byte) (string, error) {
	bi := bytesutil2.LittleEndianBytesToBigInt(b)
	if !math.IsValidUint256(bi) {
		return "", fmt.Errorf("%s is not a valid uint356", bi.String())
	}
	return string([]byte(bi.String())), nil
}
