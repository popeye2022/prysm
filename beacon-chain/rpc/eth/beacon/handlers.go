package beacon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/api"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/transition"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/rpc/eth/helpers"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	http2 "github.com/prysmaticlabs/prysm/v4/network/http"
	ethpbv1 "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
	ethpbv2 "github.com/prysmaticlabs/prysm/v4/proto/eth/v2"
	"github.com/prysmaticlabs/prysm/v4/proto/migration"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

const (
	broadcastValidationQueryParam               = "broadcast_validation"
	broadcastValidationConsensus                = "consensus"
	broadcastValidationConsensusAndEquivocation = "consensus_and_equivocation"
)

// PublishBlindedBlockV2 instructs the beacon node to use the components of the `SignedBlindedBeaconBlock` to construct and publish a
// `SignedBeaconBlock` by swapping out the `transactions_root` for the corresponding full list of `transactions`.
// The beacon node should broadcast a newly constructed `SignedBeaconBlock` to the beacon network,
// to be included in the beacon chain. The beacon node is not required to validate the signed
// `BeaconBlock`, and a successful response (20X) only indicates that the broadcast has been
// successful. The beacon node is expected to integrate the new block into its state, and
// therefore validate the block internally, however blocks which fail the validation are still
// broadcast but a different status code is returned (202). Pre-Bellatrix, this endpoint will accept
// a `SignedBeaconBlock`. The broadcast behaviour may be adjusted via the `broadcast_validation`
// query parameter.
func (bs *Server) PublishBlindedBlockV2(w http.ResponseWriter, r *http.Request) {
	if ok := bs.checkSync(r.Context(), w); !ok {
		return
	}
	isSSZ, err := http2.SszRequested(r)
	if isSSZ && err == nil {
		publishBlindedBlockV2SSZ(bs, w, r)
	} else {
		publishBlindedBlockV2(bs, w, r)
	}
}

func publishBlindedBlockV2SSZ(bs *Server, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errJson := &http2.DefaultErrorJson{
			Message: "Could not read request body: " + err.Error(),
			Code:    http.StatusInternalServerError,
		}
		http2.WriteError(w, errJson)
		return
	}
	denebBlockContents := &ethpbv2.SignedBlindedBeaconBlockContentsDeneb{}
	if err := denebBlockContents.UnmarshalSSZ(body); err == nil {
		v1block, err := migration.BlindedDenebBlockContentsToV1Alpha1(denebBlockContents)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: "Could not decode request body into consensus block: " + err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		genericBlock := &eth.GenericSignedBeaconBlock{
			Block: &eth.GenericSignedBeaconBlock_BlindedDeneb{
				BlindedDeneb: v1block,
			},
		}
		if err = bs.validateBroadcast(r, genericBlock); err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		bs.proposeBlock(r.Context(), w, genericBlock)
		return
	}
	capellaBlock := &ethpbv2.SignedBlindedBeaconBlockCapella{}
	if err := capellaBlock.UnmarshalSSZ(body); err == nil {
		v1block, err := migration.BlindedCapellaToV1Alpha1SignedBlock(capellaBlock)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: "Could not decode request body into consensus block: " + err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		genericBlock := &eth.GenericSignedBeaconBlock{
			Block: &eth.GenericSignedBeaconBlock_BlindedCapella{
				BlindedCapella: v1block,
			},
		}
		if err = bs.validateBroadcast(r, genericBlock); err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		bs.proposeBlock(r.Context(), w, genericBlock)
		return
	}
	bellatrixBlock := &ethpbv2.SignedBlindedBeaconBlockBellatrix{}
	if err := bellatrixBlock.UnmarshalSSZ(body); err == nil {
		v1block, err := migration.BlindedBellatrixToV1Alpha1SignedBlock(bellatrixBlock)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: "Could not decode request body into consensus block: " + err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		genericBlock := &eth.GenericSignedBeaconBlock{
			Block: &eth.GenericSignedBeaconBlock_BlindedBellatrix{
				BlindedBellatrix: v1block,
			},
		}
		if err = bs.validateBroadcast(r, genericBlock); err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		bs.proposeBlock(r.Context(), w, genericBlock)
		return
	}

	// blinded is not supported before bellatrix hardfork
	altairBlock := &ethpbv2.SignedBeaconBlockAltair{}
	if err := altairBlock.UnmarshalSSZ(body); err == nil {
		v1block, err := migration.AltairToV1Alpha1SignedBlock(altairBlock)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: "Could not decode request body into consensus block: " + err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		genericBlock := &eth.GenericSignedBeaconBlock{
			Block: &eth.GenericSignedBeaconBlock_Altair{
				Altair: v1block,
			},
		}
		if err = bs.validateBroadcast(r, genericBlock); err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		bs.proposeBlock(r.Context(), w, genericBlock)
		return
	}
	phase0Block := &ethpbv1.SignedBeaconBlock{}
	if err := phase0Block.UnmarshalSSZ(body); err == nil {
		v1block, err := migration.V1ToV1Alpha1SignedBlock(phase0Block)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: "Could not decode request body into consensus block: " + err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		genericBlock := &eth.GenericSignedBeaconBlock{
			Block: &eth.GenericSignedBeaconBlock_Phase0{
				Phase0: v1block,
			},
		}
		if err = bs.validateBroadcast(r, genericBlock); err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		bs.proposeBlock(r.Context(), w, genericBlock)
		return
	}
	errJson := &http2.DefaultErrorJson{
		Message: "Body does not represent a valid block type",
		Code:    http.StatusBadRequest,
	}
	http2.WriteError(w, errJson)
}

func publishBlindedBlockV2(bs *Server, w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errJson := &http2.DefaultErrorJson{
			Message: "Could not read request body",
			Code:    http.StatusInternalServerError,
		}
		http2.WriteError(w, errJson)
		return
	}
	var denebBlockContents *SignedBlindedBeaconBlockContentsDeneb
	if err = unmarshalStrict(body, &denebBlockContents); err == nil {
		if err = validate.Struct(denebBlockContents); err == nil {
			consensusBlock, err := denebBlockContents.ToGeneric()
			if err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: "Could not decode request body into consensus block: " + err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			if err = bs.validateBroadcast(r, consensusBlock); err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			bs.proposeBlock(r.Context(), w, consensusBlock)
			return
		}
	}

	var capellaBlock *SignedBlindedBeaconBlockCapella
	if err = unmarshalStrict(body, &capellaBlock); err == nil {
		if err = validate.Struct(capellaBlock); err == nil {
			consensusBlock, err := capellaBlock.ToGeneric()
			if err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: "Could not decode request body into consensus block: " + err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			if err = bs.validateBroadcast(r, consensusBlock); err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			bs.proposeBlock(r.Context(), w, consensusBlock)
			return
		}
	}

	var bellatrixBlock *SignedBlindedBeaconBlockBellatrix
	if err = unmarshalStrict(body, &bellatrixBlock); err == nil {
		if err = validate.Struct(bellatrixBlock); err == nil {
			consensusBlock, err := bellatrixBlock.ToGeneric()
			if err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: "Could not decode request body into consensus block: " + err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			if err = bs.validateBroadcast(r, consensusBlock); err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			bs.proposeBlock(r.Context(), w, consensusBlock)
			return
		}
	}
	var altairBlock *SignedBeaconBlockAltair
	if err = unmarshalStrict(body, &altairBlock); err == nil {
		if err = validate.Struct(altairBlock); err == nil {
			consensusBlock, err := altairBlock.ToGeneric()
			if err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: "Could not decode request body into consensus block: " + err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			if err = bs.validateBroadcast(r, consensusBlock); err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			bs.proposeBlock(r.Context(), w, consensusBlock)
			return
		}
	}
	var phase0Block *SignedBeaconBlock
	if err = unmarshalStrict(body, &phase0Block); err == nil {
		if err = validate.Struct(phase0Block); err == nil {
			consensusBlock, err := phase0Block.ToGeneric()
			if err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: "Could not decode request body into consensus block: " + err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			if err = bs.validateBroadcast(r, consensusBlock); err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			bs.proposeBlock(r.Context(), w, consensusBlock)
			return
		}
	}

	errJson := &http2.DefaultErrorJson{
		Message: "Body does not represent a valid block type",
		Code:    http.StatusBadRequest,
	}
	http2.WriteError(w, errJson)
}

// PublishBlockV2 instructs the beacon node to broadcast a newly signed beacon block to the beacon network,
// to be included in the beacon chain. A success response (20x) indicates that the block
// passed gossip validation and was successfully broadcast onto the network.
// The beacon node is also expected to integrate the block into the state, but may broadcast it
// before doing so, so as to aid timely delivery of the block. Should the block fail full
// validation, a separate success response code (202) is used to indicate that the block was
// successfully broadcast but failed integration. The broadcast behaviour may be adjusted via the
// `broadcast_validation` query parameter.
func (bs *Server) PublishBlockV2(w http.ResponseWriter, r *http.Request) {
	if ok := bs.checkSync(r.Context(), w); !ok {
		return
	}
	isSSZ, err := http2.SszRequested(r)
	if isSSZ && err == nil {
		publishBlockV2SSZ(bs, w, r)
	} else {
		publishBlockV2(bs, w, r)
	}
}

func publishBlockV2SSZ(bs *Server, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errJson := &http2.DefaultErrorJson{
			Message: "Could not read request body",
			Code:    http.StatusInternalServerError,
		}
		http2.WriteError(w, errJson)
		return
	}
	denebBlockContents := &ethpbv2.SignedBeaconBlockContentsDeneb{}
	if err := denebBlockContents.UnmarshalSSZ(body); err == nil {
		v1BlockContents, err := migration.DenebBlockContentsToV1Alpha1(denebBlockContents)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: "Could not decode request body into consensus block: " + err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		genericBlock := &eth.GenericSignedBeaconBlock{
			Block: &eth.GenericSignedBeaconBlock_Deneb{
				Deneb: v1BlockContents,
			},
		}
		if err = bs.validateBroadcast(r, genericBlock); err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		bs.proposeBlock(r.Context(), w, genericBlock)
		return
	}
	capellaBlock := &ethpbv2.SignedBeaconBlockCapella{}
	if err := capellaBlock.UnmarshalSSZ(body); err == nil {
		v1block, err := migration.CapellaToV1Alpha1SignedBlock(capellaBlock)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: "Could not decode request body into consensus block: " + err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		genericBlock := &eth.GenericSignedBeaconBlock{
			Block: &eth.GenericSignedBeaconBlock_Capella{
				Capella: v1block,
			},
		}
		if err = bs.validateBroadcast(r, genericBlock); err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		bs.proposeBlock(r.Context(), w, genericBlock)
		return
	}
	bellatrixBlock := &ethpbv2.SignedBeaconBlockBellatrix{}
	if err := bellatrixBlock.UnmarshalSSZ(body); err == nil {
		v1block, err := migration.BellatrixToV1Alpha1SignedBlock(bellatrixBlock)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: "Could not decode request body into consensus block: " + err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		genericBlock := &eth.GenericSignedBeaconBlock{
			Block: &eth.GenericSignedBeaconBlock_Bellatrix{
				Bellatrix: v1block,
			},
		}
		if err = bs.validateBroadcast(r, genericBlock); err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		bs.proposeBlock(r.Context(), w, genericBlock)
		return
	}
	altairBlock := &ethpbv2.SignedBeaconBlockAltair{}
	if err := altairBlock.UnmarshalSSZ(body); err == nil {
		v1block, err := migration.AltairToV1Alpha1SignedBlock(altairBlock)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: "Could not decode request body into consensus block: " + err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		genericBlock := &eth.GenericSignedBeaconBlock{
			Block: &eth.GenericSignedBeaconBlock_Altair{
				Altair: v1block,
			},
		}
		if err = bs.validateBroadcast(r, genericBlock); err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		bs.proposeBlock(r.Context(), w, genericBlock)
		return
	}
	phase0Block := &ethpbv1.SignedBeaconBlock{}
	if err := phase0Block.UnmarshalSSZ(body); err == nil {
		v1block, err := migration.V1ToV1Alpha1SignedBlock(phase0Block)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: "Could not decode request body into consensus block: " + err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		genericBlock := &eth.GenericSignedBeaconBlock{
			Block: &eth.GenericSignedBeaconBlock_Phase0{
				Phase0: v1block,
			},
		}
		if err = bs.validateBroadcast(r, genericBlock); err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		bs.proposeBlock(r.Context(), w, genericBlock)
		return
	}
	errJson := &http2.DefaultErrorJson{
		Message: "Body does not represent a valid block type",
		Code:    http.StatusBadRequest,
	}
	http2.WriteError(w, errJson)
}

func publishBlockV2(bs *Server, w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errJson := &http2.DefaultErrorJson{
			Message: "Could not read request body",
			Code:    http.StatusInternalServerError,
		}
		http2.WriteError(w, errJson)
		return
	}
	var denebBlockContents *SignedBeaconBlockContentsDeneb
	if err = unmarshalStrict(body, &denebBlockContents); err == nil {
		if err = validate.Struct(denebBlockContents); err == nil {
			consensusBlock, err := denebBlockContents.ToGeneric()
			if err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: "Could not decode request body into consensus block: " + err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			if err = bs.validateBroadcast(r, consensusBlock); err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			bs.proposeBlock(r.Context(), w, consensusBlock)
			return
		}
	}
	var capellaBlock *SignedBeaconBlockCapella
	if err = unmarshalStrict(body, &capellaBlock); err == nil {
		if err = validate.Struct(capellaBlock); err == nil {
			consensusBlock, err := capellaBlock.ToGeneric()
			if err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: "Could not decode request body into consensus block: " + err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			if err = bs.validateBroadcast(r, consensusBlock); err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			bs.proposeBlock(r.Context(), w, consensusBlock)
			return
		}
	}
	var bellatrixBlock *SignedBeaconBlockBellatrix
	if err = unmarshalStrict(body, &bellatrixBlock); err == nil {
		if err = validate.Struct(bellatrixBlock); err == nil {
			consensusBlock, err := bellatrixBlock.ToGeneric()
			if err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: "Could not decode request body into consensus block: " + err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			if err = bs.validateBroadcast(r, consensusBlock); err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			bs.proposeBlock(r.Context(), w, consensusBlock)
			return
		}
	}
	var altairBlock *SignedBeaconBlockAltair
	if err = unmarshalStrict(body, &altairBlock); err == nil {
		if err = validate.Struct(altairBlock); err == nil {
			consensusBlock, err := altairBlock.ToGeneric()
			if err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: "Could not decode request body into consensus block: " + err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			if err = bs.validateBroadcast(r, consensusBlock); err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			bs.proposeBlock(r.Context(), w, consensusBlock)
			return
		}
	}
	var phase0Block *SignedBeaconBlock
	if err = unmarshalStrict(body, &phase0Block); err == nil {
		if err = validate.Struct(phase0Block); err == nil {
			consensusBlock, err := phase0Block.ToGeneric()
			if err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: "Could not decode request body into consensus block: " + err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			if err = bs.validateBroadcast(r, consensusBlock); err != nil {
				errJson := &http2.DefaultErrorJson{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				}
				http2.WriteError(w, errJson)
				return
			}
			bs.proposeBlock(r.Context(), w, consensusBlock)
			return
		}
	}

	errJson := &http2.DefaultErrorJson{
		Message: "Body does not represent a valid block type",
		Code:    http.StatusBadRequest,
	}
	http2.WriteError(w, errJson)
}

func (bs *Server) proposeBlock(ctx context.Context, w http.ResponseWriter, blk *eth.GenericSignedBeaconBlock) {
	_, err := bs.V1Alpha1ValidatorServer.ProposeBeaconBlock(ctx, blk)
	if err != nil {
		errJson := &http2.DefaultErrorJson{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		http2.WriteError(w, errJson)
		return
	}
}

func unmarshalStrict(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func (bs *Server) validateBroadcast(r *http.Request, blk *eth.GenericSignedBeaconBlock) error {
	switch r.URL.Query().Get(broadcastValidationQueryParam) {
	case broadcastValidationConsensus:
		b, err := blocks.NewSignedBeaconBlock(blk.Block)
		if err != nil {
			return errors.Wrapf(err, "could not create signed beacon block")
		}
		if err = bs.validateConsensus(r.Context(), b); err != nil {
			return errors.Wrap(err, "consensus validation failed")
		}
	case broadcastValidationConsensusAndEquivocation:
		b, err := blocks.NewSignedBeaconBlock(blk.Block)
		if err != nil {
			return errors.Wrapf(err, "could not create signed beacon block")
		}
		if err = bs.validateConsensus(r.Context(), b); err != nil {
			return errors.Wrap(err, "consensus validation failed")
		}
		if err = bs.validateEquivocation(b.Block()); err != nil {
			return errors.Wrap(err, "equivocation validation failed")
		}
	default:
		return nil
	}
	return nil
}

func (bs *Server) validateConsensus(ctx context.Context, blk interfaces.ReadOnlySignedBeaconBlock) error {
	parentBlockRoot := blk.Block().ParentRoot()
	parentBlock, err := bs.Blocker.Block(ctx, parentBlockRoot[:])
	if err != nil {
		return errors.Wrap(err, "could not get parent block")
	}
	parentStateRoot := parentBlock.Block().StateRoot()
	parentState, err := bs.Stater.State(ctx, parentStateRoot[:])
	if err != nil {
		return errors.Wrap(err, "could not get parent state")
	}
	_, err = transition.ExecuteStateTransition(ctx, parentState, blk)
	if err != nil {
		return errors.Wrap(err, "could not execute state transition")
	}
	return nil
}

func (bs *Server) validateEquivocation(blk interfaces.ReadOnlyBeaconBlock) error {
	if bs.ForkchoiceFetcher.HighestReceivedBlockSlot() == blk.Slot() {
		return fmt.Errorf("block for slot %d already exists in fork choice", blk.Slot())
	}
	return nil
}

func (bs *Server) checkSync(ctx context.Context, w http.ResponseWriter) bool {
	isSyncing, syncDetails, err := helpers.ValidateSyncHTTP(ctx, bs.SyncChecker, bs.HeadFetcher, bs.TimeFetcher, bs.OptimisticModeFetcher)
	if err != nil {
		errJson := &http2.DefaultErrorJson{
			Message: "Could not check if node is syncing: " + err.Error(),
			Code:    http.StatusInternalServerError,
		}
		http2.WriteError(w, errJson)
		return false
	}
	if isSyncing {
		msg := "Beacon node is currently syncing and not serving request on that endpoint"
		details, err := json.Marshal(syncDetails)
		if err == nil {
			msg += " Details: " + string(details)
		}
		errJson := &http2.DefaultErrorJson{
			Message: msg,
			Code:    http.StatusServiceUnavailable,
		}
		http2.WriteError(w, errJson)
		return false
	}
	return true
}

// ProduceBlockV3 Requests a beacon node to produce a valid block, which can then be signed by a validator. The
// returned block may be blinded or unblinded, depending on the current state of the network as
// decided by the execution and beacon nodes.
// The beacon node must return an unblinded block if it obtains the execution payload from its
// paired execution node. It must only return a blinded block if it obtains the execution payload
// header from an MEV relay.
// Metadata in the response indicates the type of block produced, and the supported types of block
// will be added to as forks progress.
func (bs *Server) ProduceBlockV3(w http.ResponseWriter, r *http.Request) {
	if ok := bs.checkSync(r.Context(), w); !ok {
		return
	}
	segments := strings.Split(r.URL.Path, "/")
	rawSlot := segments[len(segments)-1]
	rawRandaoReveal := r.URL.Query().Get("randao_reveal")
	rawGraffiti := r.URL.Query().Get("graffiti")
	rawSkipRandaoVerification := r.URL.Query().Get("skip_randao_verification")

	if rawSlot == "" {
		errJson := &http2.DefaultErrorJson{
			Message: "slot is required",
			Code:    http.StatusBadRequest,
		}
		http2.WriteError(w, errJson)
		return
	}
	slot, err := strconv.ParseUint(rawSlot, 10, 64)
	if err != nil {
		errJson := &http2.DefaultErrorJson{
			Message: "slot is invalid: " + err.Error(),
			Code:    http.StatusBadRequest,
		}
		http2.WriteError(w, errJson)
		return
	}
	var randaoReveal []byte
	if rawSkipRandaoVerification == "true" {
		randaoReveal = primitives.PointAtInfinity
	} else {
		rr, err := hexutil.Decode(rawRandaoReveal)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: errors.Wrap(err, "unable to decode randao reveal").Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		randaoReveal = rr
	}
	if len(randaoReveal) != fieldparams.BLSSignatureLength {
		errJson := &http2.DefaultErrorJson{
			Message: fmt.Sprintf("a valid randao reveal is required as a query parameter: received length %d but wanted length %d", len(randaoReveal), fieldparams.BLSSignatureLength),
			Code:    http.StatusBadRequest,
		}
		http2.WriteError(w, errJson)
		return
	}
	var graffiti []byte
	if rawGraffiti != "" {
		g, err := hexutil.Decode(rawGraffiti)
		if err != nil {
			errJson := &http2.DefaultErrorJson{
				Message: errors.Wrap(err, "unable to decode graffiti").Error(),
				Code:    http.StatusBadRequest,
			}
			http2.WriteError(w, errJson)
			return
		}
		graffiti = g
	}

	produceBlockV3(bs, w, r, &eth.BlockRequest{
		Slot:         primitives.Slot(slot),
		RandaoReveal: randaoReveal,
		Graffiti:     graffiti,
		SkipMevBoost: false,
	})

}

func produceBlockV3(bs *Server, w http.ResponseWriter, r *http.Request, v1alpha1req *eth.BlockRequest) {
	isSSZ, err := http2.SszRequested(r)
	if err != nil {
		log.WithError(err).Error("verifying ssz request failed, defaulting to non ssz.")
		isSSZ = false
	}
	validate := validator.New()
	v1alpha1resp, err := bs.V1Alpha1ValidatorServer.GetBeaconBlock(r.Context(), v1alpha1req)
	if err != nil {
		http2.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	phase0Block, ok := v1alpha1resp.Block.(*eth.GenericBeaconBlock_Phase0)
	if ok {
		handleProducePhase0V3(isSSZ, phase0Block, w, validate, v1alpha1resp.IsBlinded, v1alpha1resp.PayloadValue)
		return
	}
	altairBlock, ok := v1alpha1resp.Block.(*eth.GenericBeaconBlock_Altair)
	if ok {
		handleProduceAltairV3(isSSZ, altairBlock, w, validate, v1alpha1resp.IsBlinded, v1alpha1resp.PayloadValue)
		return
	}
	optimistic, err := bs.OptimisticModeFetcher.IsOptimistic(r.Context())
	if err != nil {
		http2.HandleError(w, errors.Wrap(err, "Could not determine if the node is a optimistic node").Error(), http.StatusInternalServerError)
		return
	}
	if optimistic {
		http2.HandleError(w, "The node is currently optimistic and cannot serve validators", http.StatusInternalServerError)
		return
	}
	blindedBellatrixBlock, ok := v1alpha1resp.Block.(*eth.GenericBeaconBlock_BlindedBellatrix)
	if ok {
		handleProduceBlindedBellatrixV3(isSSZ, blindedBellatrixBlock, w, validate, v1alpha1resp.IsBlinded, v1alpha1resp.PayloadValue)
		return
	}
	bellatrixBlock, ok := v1alpha1resp.Block.(*eth.GenericBeaconBlock_Bellatrix)
	if ok {
		handleProduceBellatrixV3(isSSZ, bellatrixBlock, w, validate, v1alpha1resp.IsBlinded, v1alpha1resp.PayloadValue)
		return
	}
	blindedCapellaBlock, ok := v1alpha1resp.Block.(*eth.GenericBeaconBlock_BlindedCapella)
	if ok {
		handleProduceBlindedCapellaV3(isSSZ, blindedCapellaBlock, w, validate, v1alpha1resp.IsBlinded, v1alpha1resp.PayloadValue)
		return
	}
	capellaBlock, ok := v1alpha1resp.Block.(*eth.GenericBeaconBlock_Capella)
	if ok {
		handleProduceCapellaV3(isSSZ, capellaBlock, w, validate, v1alpha1resp.IsBlinded, v1alpha1resp.PayloadValue)
		return
	}
	blindedDenebBlockContents, ok := v1alpha1resp.Block.(*eth.GenericBeaconBlock_BlindedDeneb)
	if ok {
		handleProduceBlindedDenebV3(isSSZ, blindedDenebBlockContents, w, validate, v1alpha1resp.IsBlinded, v1alpha1resp.PayloadValue)
		return
	}
	denebBlockContents, ok := v1alpha1resp.Block.(*eth.GenericBeaconBlock_Deneb)
	if ok {
		handleProduceDenebV3(isSSZ, denebBlockContents, w, validate, v1alpha1resp.IsBlinded, v1alpha1resp.PayloadValue)
		return
	}
}

func handleProducePhase0V3(isSSZ bool, phase0Block *eth.GenericBeaconBlock_Phase0, w http.ResponseWriter, validate *validator.Validate, isBlinded bool, payloadValue uint64) {
	w.Header().Set(api.ExecutionPayloadBlindedHeader, fmt.Sprintf("%v", isBlinded))
	w.Header().Set(api.ExecutionPayloadValueHeader, fmt.Sprintf("%d", payloadValue))
	if isSSZ {
		sszResp, err := phase0Block.Phase0.MarshalSSZ()
		if err != nil {
			http2.HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http2.WriteSsz(w, sszResp, "phase0Block.ssz")
		return
	}
	block, err := convertInternalBeaconBlock(phase0Block.Phase0)
	if err != nil {
		http2.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = validate.Struct(block); err == nil {
		http2.WriteJson(w, &Phase0ProduceBlockV3Response{
			Version:                 version.String(version.Phase0),
			ExecutionPayloadBlinded: false,
			ExeuctionPayloadValue:   fmt.Sprintf("%d", payloadValue), // mev not available at this point
			Data:                    block,
		})
		return
	}
}

func handleProduceAltairV3(isSSZ bool, altairBlock *eth.GenericBeaconBlock_Altair, w http.ResponseWriter, validate *validator.Validate, isBlinded bool, payloadValue uint64) {
	w.Header().Set(api.ExecutionPayloadBlindedHeader, fmt.Sprintf("%v", isBlinded))
	w.Header().Set(api.ExecutionPayloadValueHeader, fmt.Sprintf("%d", payloadValue))
	if isSSZ {
		sszResp, err := altairBlock.Altair.MarshalSSZ()
		if err != nil {
			http2.HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http2.WriteSsz(w, sszResp, "altairBlock.ssz")
		return
	}
	block, err := convertInternalBeaconBlockAltair(altairBlock.Altair)
	if err != nil {
		http2.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = validate.Struct(block); err == nil {
		http2.WriteJson(w, &AltairProduceBlockV3Response{
			Version:                 version.String(version.Altair),
			ExecutionPayloadBlinded: false,
			ExeuctionPayloadValue:   fmt.Sprintf("%d", payloadValue), // mev not available at this point
			Data:                    block,
		})
		return
	}
}

func handleProduceBellatrixV3(isSSZ bool, bellatrixBlock *eth.GenericBeaconBlock_Bellatrix, w http.ResponseWriter, validate *validator.Validate, isBlinded bool, payloadValue uint64) {
	w.Header().Set(api.ExecutionPayloadBlindedHeader, fmt.Sprintf("%v", isBlinded))
	w.Header().Set(api.ExecutionPayloadValueHeader, fmt.Sprintf("%d", payloadValue))
	if isSSZ {
		sszResp, err := bellatrixBlock.Bellatrix.MarshalSSZ()
		if err != nil {
			http2.HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http2.WriteSsz(w, sszResp, "bellatrixBlock.ssz")
		return
	}
	block, err := convertInternalBeaconBlockBellatrix(bellatrixBlock.Bellatrix)
	if err != nil {
		http2.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = validate.Struct(block); err == nil {
		http2.WriteJson(w, &BellatrixProduceBlockV3Response{
			Version:                 version.String(version.Bellatrix),
			ExecutionPayloadBlinded: false,
			ExeuctionPayloadValue:   fmt.Sprintf("%d", payloadValue), // mev not available at this point
			Data:                    block,
		})
		return
	}
}

func handleProduceBlindedBellatrixV3(isSSZ bool, blindedBellatrixBlock *eth.GenericBeaconBlock_BlindedBellatrix, w http.ResponseWriter, validate *validator.Validate, isBlinded bool, payloadValue uint64) {
	w.Header().Set(api.ExecutionPayloadBlindedHeader, fmt.Sprintf("%v", isBlinded))
	w.Header().Set(api.ExecutionPayloadValueHeader, fmt.Sprintf("%d", payloadValue))
	if isSSZ {
		sszResp, err := blindedBellatrixBlock.BlindedBellatrix.MarshalSSZ()
		if err != nil {
			http2.HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http2.WriteSsz(w, sszResp, "blindeBellatrixBlock.ssz")
		return
	}
	block, err := convertInternalBlindedBeaconBlockBellatrix(blindedBellatrixBlock.BlindedBellatrix)
	if err != nil {
		http2.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = validate.Struct(block); err == nil {
		http2.WriteJson(w, &BlindedBellatrixProduceBlockV3Response{
			Version:                 version.String(version.Bellatrix),
			ExecutionPayloadBlinded: true,
			ExeuctionPayloadValue:   fmt.Sprintf("%d", payloadValue),
			Data:                    block,
		})
		return
	}
}

func handleProduceBlindedCapellaV3(isSSZ bool, blindedCapellaBlock *eth.GenericBeaconBlock_BlindedCapella, w http.ResponseWriter, validate *validator.Validate, isBlinded bool, payloadValue uint64) {
	w.Header().Set(api.ExecutionPayloadBlindedHeader, fmt.Sprintf("%v", isBlinded))
	w.Header().Set(api.ExecutionPayloadValueHeader, fmt.Sprintf("%d", payloadValue))
	if isSSZ {
		sszResp, err := blindedCapellaBlock.BlindedCapella.MarshalSSZ()
		if err != nil {
			http2.HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http2.WriteSsz(w, sszResp, "blindedCapellaBlock.ssz")
		return
	}
	block, err := convertInternalBlindedBeaconBlockCapella(blindedCapellaBlock.BlindedCapella)
	if err != nil {
		http2.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = validate.Struct(block); err == nil {
		http2.WriteJson(w, &BlindedCapellaProduceBlockV3Response{
			Version:                 version.String(version.Capella),
			ExecutionPayloadBlinded: true,
			ExeuctionPayloadValue:   fmt.Sprintf("%d", payloadValue),
			Data:                    block,
		})
		return
	}
}

func handleProduceCapellaV3(isSSZ bool, capellaBlock *eth.GenericBeaconBlock_Capella, w http.ResponseWriter, validate *validator.Validate, isBlinded bool, payloadValue uint64) {
	w.Header().Set(api.ExecutionPayloadBlindedHeader, fmt.Sprintf("%v", isBlinded))
	w.Header().Set(api.ExecutionPayloadValueHeader, fmt.Sprintf("%d", payloadValue))
	if isSSZ {
		sszResp, err := capellaBlock.Capella.MarshalSSZ()
		if err != nil {
			http2.HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http2.WriteSsz(w, sszResp, "capellaBlock.ssz")
		return
	}
	block, err := convertInternalBeaconBlockCapella(capellaBlock.Capella)
	if err != nil {
		http2.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = validate.Struct(block); err == nil {
		http2.WriteJson(w, &CapellaProduceBlockV3Response{
			Version:                 version.String(version.Capella),
			ExecutionPayloadBlinded: false,
			ExeuctionPayloadValue:   fmt.Sprintf("%d", payloadValue), // mev not available at this point
			Data:                    block,
		})
		return
	}
}

func handleProduceBlindedDenebV3(isSSZ bool, blindedDenebBlockContents *eth.GenericBeaconBlock_BlindedDeneb, w http.ResponseWriter, validate *validator.Validate, isBlinded bool, payloadValue uint64) {
	w.Header().Set(api.ExecutionPayloadBlindedHeader, fmt.Sprintf("%v", isBlinded))
	w.Header().Set(api.ExecutionPayloadValueHeader, fmt.Sprintf("%d", payloadValue))
	if isSSZ {
		sszResp, err := blindedDenebBlockContents.BlindedDeneb.MarshalSSZ()
		if err != nil {
			http2.HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http2.WriteSsz(w, sszResp, "blindedDenebBlockContents.ssz")
		return
	}
	blockContents, err := convertInternalBlindedBeaconBlockContentsDeneb(blindedDenebBlockContents.BlindedDeneb)
	if err != nil {
		http2.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = validate.Struct(blockContents); err == nil {
		http2.WriteJson(w, &BlindedDenebProduceBlockV3Response{
			Version:                 version.String(version.Deneb),
			ExecutionPayloadBlinded: true,
			ExeuctionPayloadValue:   fmt.Sprintf("%d", payloadValue),
			Data:                    blockContents,
		})
		return
	}
}

func handleProduceDenebV3(isSSZ bool, denebBlockContents *eth.GenericBeaconBlock_Deneb, w http.ResponseWriter, validate *validator.Validate, isBlinded bool, payloadValue uint64) {
	w.Header().Set(api.ExecutionPayloadBlindedHeader, fmt.Sprintf("%v", isBlinded))
	w.Header().Set(api.ExecutionPayloadValueHeader, fmt.Sprintf("%d", payloadValue))
	if isSSZ {
		sszResp, err := denebBlockContents.Deneb.MarshalSSZ()
		if err != nil {
			http2.HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http2.WriteSsz(w, sszResp, "denebBlockContents.ssz")
		return
	}
	blockContents, err := convertInternalBeaconBlockContentsDeneb(denebBlockContents.Deneb)
	if err != nil {
		http2.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = validate.Struct(blockContents); err == nil {
		http2.WriteJson(w, &DenebProduceBlockV3Response{
			Version:                 version.String(version.Deneb),
			ExecutionPayloadBlinded: false,
			ExeuctionPayloadValue:   fmt.Sprintf("%d", payloadValue), // mev not available at this point
			Data:                    blockContents,
		})
		return
	}
}
