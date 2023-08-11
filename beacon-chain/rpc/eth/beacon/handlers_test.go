package beacon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/prysmaticlabs/prysm/v4/api"
	blockchainTesting "github.com/prysmaticlabs/prysm/v4/beacon-chain/blockchain/testing"
	testing2 "github.com/prysmaticlabs/prysm/v4/beacon-chain/blockchain/testing"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/transition"
	doublylinkedtree "github.com/prysmaticlabs/prysm/v4/beacon-chain/forkchoice/doubly-linked-tree"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/rpc/testutil"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	mockSync "github.com/prysmaticlabs/prysm/v4/beacon-chain/sync/initial-sync/testing"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/v4/proto/migration"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/testing/assert"
	mock2 "github.com/prysmaticlabs/prysm/v4/testing/mock"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	"github.com/prysmaticlabs/prysm/v4/testing/util"
	"github.com/stretchr/testify/mock"
)

func TestPublishBlockV2(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("Phase 0", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Phase0)
			converted, err := convertInternalBeaconBlock(block.Phase0.Block)
			require.NoError(t, err)
			var signedblock *SignedBeaconBlock
			err = json.Unmarshal([]byte(phase0Block), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(phase0Block)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Altair", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Altair)
			converted, err := convertInternalBeaconBlockAltair(block.Altair.Block)
			require.NoError(t, err)
			var signedblock *SignedBeaconBlockAltair
			err = json.Unmarshal([]byte(altairBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(altairBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Bellatrix", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Bellatrix)
			converted, err := convertInternalBeaconBlockBellatrix(block.Bellatrix.Block)
			require.NoError(t, err)
			var signedblock *SignedBeaconBlockBellatrix
			err = json.Unmarshal([]byte(bellatrixBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(bellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Capella", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Capella)
			converted, err := convertInternalBeaconBlockCapella(block.Capella.Block)
			require.NoError(t, err)
			var signedblock *SignedBeaconBlockCapella
			err = json.Unmarshal([]byte(capellaBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(capellaBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Deneb", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Deneb)
			converted, err := convertInternalToDenebBlock(block.Deneb.Block.Block)
			require.NoError(t, err)
			var signedblock *SignedBeaconBlockContentsDeneb
			err = json.Unmarshal([]byte(denebBlockContents), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.SignedBlock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(denebBlockContents)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("invalid block", func(t *testing.T) {
		server := &Server{
			SyncChecker: &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(blindedBellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Body does not represent a valid block type"))
	})
	t.Run("syncing", func(t *testing.T) {
		chainService := &testing2.ChainService{}
		server := &Server{
			SyncChecker:           &mockSync.Sync{IsSyncing: true},
			HeadFetcher:           chainService,
			TimeFetcher:           chainService,
			OptimisticModeFetcher: chainService,
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte("foo")))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusServiceUnavailable, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Beacon node is currently syncing and not serving request on that endpoint"))
	})
}

func TestPublishBlockV2SSZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Run("Bellatrix", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_Bellatrix)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}
		var bellablock SignedBeaconBlockBellatrix
		err := json.Unmarshal([]byte(bellatrixBlock), &bellablock)
		require.NoError(t, err)
		genericBlock, err := bellablock.ToGeneric()
		require.NoError(t, err)
		sszvalue, err := genericBlock.GetBellatrix().MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Capella", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_Capella)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		var cblock SignedBeaconBlockCapella
		err := json.Unmarshal([]byte(capellaBlock), &cblock)
		require.NoError(t, err)
		genericBlock, err := cblock.ToGeneric()
		require.NoError(t, err)
		sszvalue, err := genericBlock.GetCapella().MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Deneb", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_Deneb)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		var dblock SignedBeaconBlockContentsDeneb
		err := json.Unmarshal([]byte(denebBlockContents), &dblock)
		require.NoError(t, err)
		genericBlock, err := dblock.ToGeneric()
		require.NoError(t, err)
		v2block, err := migration.V1Alpha1SignedBeaconBlockDenebAndBlobsToV2(genericBlock.GetDeneb())
		require.NoError(t, err)
		sszvalue, err := v2block.MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("invalid block", func(t *testing.T) {
		server := &Server{
			SyncChecker: &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(blindedBellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Body does not represent a valid block type"))
	})
}

func TestPublishBlindedBlockV2(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Run("Phase 0", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Phase0)
			converted, err := convertInternalBeaconBlock(block.Phase0.Block)
			require.NoError(t, err)
			var signedblock *SignedBeaconBlock
			err = json.Unmarshal([]byte(phase0Block), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(phase0Block)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Altair", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Altair)
			converted, err := convertInternalBeaconBlockAltair(block.Altair.Block)
			require.NoError(t, err)
			var signedblock *SignedBeaconBlockAltair
			err = json.Unmarshal([]byte(altairBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(altairBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Blinded Bellatrix", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedBellatrix)
			converted, err := convertInternalBlindedBeaconBlockBellatrix(block.BlindedBellatrix.Block)
			require.NoError(t, err)
			var signedblock *SignedBlindedBeaconBlockBellatrix
			err = json.Unmarshal([]byte(blindedBellatrixBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(blindedBellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Blinded Capella", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedCapella)
			converted, err := convertInternalBlindedBeaconBlockCapella(block.BlindedCapella.Block)
			require.NoError(t, err)
			var signedblock *SignedBlindedBeaconBlockCapella
			err = json.Unmarshal([]byte(blindedCapellaBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(blindedCapellaBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Blinded Deneb", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedDeneb)
			converted, err := convertInternalToBlindedDenebBlock(block.BlindedDeneb.Block.Block)
			require.NoError(t, err)
			var signedblock *SignedBlindedBeaconBlockContentsDeneb
			err = json.Unmarshal([]byte(blindedDenebBlockContents), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.SignedBlindedBlock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(blindedDenebBlockContents)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("invalid block", func(t *testing.T) {
		server := &Server{
			SyncChecker: &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(bellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Body does not represent a valid block type"))
	})
	t.Run("syncing", func(t *testing.T) {
		chainService := &testing2.ChainService{}
		server := &Server{
			SyncChecker:           &mockSync.Sync{IsSyncing: true},
			HeadFetcher:           chainService,
			TimeFetcher:           chainService,
			OptimisticModeFetcher: chainService,
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte("foo")))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusServiceUnavailable, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Beacon node is currently syncing and not serving request on that endpoint"))
	})
}

func TestPublishBlindedBlockV2SSZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Run("Bellatrix", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedBellatrix)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		var bellablock SignedBlindedBeaconBlockBellatrix
		err := json.Unmarshal([]byte(blindedBellatrixBlock), &bellablock)
		require.NoError(t, err)
		genericBlock, err := bellablock.ToGeneric()
		require.NoError(t, err)
		sszvalue, err := genericBlock.GetBlindedBellatrix().MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Capella", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedCapella)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		var cblock SignedBlindedBeaconBlockCapella
		err := json.Unmarshal([]byte(blindedCapellaBlock), &cblock)
		require.NoError(t, err)
		genericBlock, err := cblock.ToGeneric()
		require.NoError(t, err)
		sszvalue, err := genericBlock.GetBlindedCapella().MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Deneb", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedDeneb)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		var cblock SignedBlindedBeaconBlockContentsDeneb
		err := json.Unmarshal([]byte(blindedDenebBlockContents), &cblock)
		require.NoError(t, err)
		genericBlock, err := cblock.ToGeneric()
		require.NoError(t, err)
		v1block, err := migration.V1Alpha1SignedBlindedBlockAndBlobsDenebToV2Blinded(genericBlock.GetBlindedDeneb())
		require.NoError(t, err)
		sszvalue, err := v1block.MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("invalid block", func(t *testing.T) {
		server := &Server{
			SyncChecker: &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(bellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Body does not represent a valid block type"))
	})
}

func TestProduceBlockV3(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("Phase 0", func(t *testing.T) {
		var block *SignedBeaconBlock
		err := json.Unmarshal([]byte(phase0Block), &block)
		require.NoError(t, err)
		jsonBytes, err := json.Marshal(block.Message)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.Message.ToGeneric()
			}())
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		want := fmt.Sprintf(`{"version":"phase0","execution_payload_blinded":false,"execution_payload_value":"0","data":%s}`, string(jsonBytes))
		body := strings.ReplaceAll(writer.Body.String(), "\n", "")
		require.Equal(t, want, body)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "false", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("Altair", func(t *testing.T) {
		var block *SignedBeaconBlockAltair
		err := json.Unmarshal([]byte(altairBlock), &block)
		require.NoError(t, err)
		jsonBytes, err := json.Marshal(block.Message)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.Message.ToGeneric()
			}())
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		want := fmt.Sprintf(`{"version":"altair","execution_payload_blinded":false,"execution_payload_value":"0","data":%s}`, string(jsonBytes))
		body := strings.ReplaceAll(writer.Body.String(), "\n", "")
		require.Equal(t, want, body)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "false", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("Bellatrix", func(t *testing.T) {
		var block *SignedBeaconBlockBellatrix
		err := json.Unmarshal([]byte(bellatrixBlock), &block)
		require.NoError(t, err)
		jsonBytes, err := json.Marshal(block.Message)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.Message.ToGeneric()
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		want := fmt.Sprintf(`{"version":"bellatrix","execution_payload_blinded":false,"execution_payload_value":"0","data":%s}`, string(jsonBytes))
		body := strings.ReplaceAll(writer.Body.String(), "\n", "")
		require.Equal(t, want, body)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "false", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("BlindedBellatrix", func(t *testing.T) {
		var block *SignedBlindedBeaconBlockBellatrix
		err := json.Unmarshal([]byte(blindedBellatrixBlock), &block)
		require.NoError(t, err)
		jsonBytes, err := json.Marshal(block.Message)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.Message.ToGeneric()
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		want := fmt.Sprintf(`{"version":"bellatrix","execution_payload_blinded":true,"execution_payload_value":"0","data":%s}`, string(jsonBytes))
		body := strings.ReplaceAll(writer.Body.String(), "\n", "")
		require.Equal(t, want, body)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "true", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("Capella", func(t *testing.T) {
		var block *SignedBeaconBlockCapella
		err := json.Unmarshal([]byte(capellaBlock), &block)
		require.NoError(t, err)
		jsonBytes, err := json.Marshal(block.Message)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.Message.ToGeneric()
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		want := fmt.Sprintf(`{"version":"capella","execution_payload_blinded":false,"execution_payload_value":"0","data":%s}`, string(jsonBytes))
		body := strings.ReplaceAll(writer.Body.String(), "\n", "")
		require.Equal(t, want, body)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "false", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("Blinded Capella", func(t *testing.T) {
		var block *SignedBlindedBeaconBlockCapella
		err := json.Unmarshal([]byte(blindedCapellaBlock), &block)
		require.NoError(t, err)
		jsonBytes, err := json.Marshal(block.Message)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				g, err := block.Message.ToGeneric()
				require.NoError(t, err)
				g.PayloadValue = 2000 //some fake value
				return g, err
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		want := fmt.Sprintf(`{"version":"capella","execution_payload_blinded":true,"execution_payload_value":"2000","data":%s}`, string(jsonBytes))
		body := strings.ReplaceAll(writer.Body.String(), "\n", "")
		require.Equal(t, want, body)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "true", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "2000", true)
	})
	t.Run("Deneb", func(t *testing.T) {
		var block *SignedBeaconBlockContentsDeneb
		err := json.Unmarshal([]byte(denebBlockContents), &block)
		require.NoError(t, err)
		jsonBytes, err := json.Marshal(block.ToUnsigned())
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.ToUnsigned().ToGeneric()
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		want := fmt.Sprintf(`{"version":"deneb","execution_payload_blinded":false,"execution_payload_value":"0","data":%s}`, string(jsonBytes))
		body := strings.ReplaceAll(writer.Body.String(), "\n", "")
		require.Equal(t, want, body)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "false", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("Blinded Deneb", func(t *testing.T) {
		var block *SignedBlindedBeaconBlockContentsDeneb
		err := json.Unmarshal([]byte(blindedDenebBlockContents), &block)
		require.NoError(t, err)
		jsonBytes, err := json.Marshal(block.ToUnsigned())
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.ToUnsigned().ToGeneric()
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		want := fmt.Sprintf(`{"version":"deneb","execution_payload_blinded":true,"execution_payload_value":"0","data":%s}`, string(jsonBytes))
		body := strings.ReplaceAll(writer.Body.String(), "\n", "")
		require.Equal(t, want, body)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "true", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("invalid query parameter slot empty", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}
		request := httptest.NewRequest(http.MethodGet, "http://foo.example/eth/v3/validator/blocks/", nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "slot is required"))
	})
	t.Run("invalid query parameter slot invalid", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}
		request := httptest.NewRequest(http.MethodGet, "http://foo.example/eth/v3/validator/blocks/asdfsad", nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "slot is invalid"))
	})
	t.Run("invalid query parameter randao_reveal invalid", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}
		request := httptest.NewRequest(http.MethodGet, "http://foo.example/eth/v3/validator/blocks/1?randao_reveal=0x213123", nil)
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "a valid randao reveal is required as a query parameter"))
	})
	t.Run("syncing", func(t *testing.T) {
		chainService := &testing2.ChainService{}
		server := &Server{
			SyncChecker:           &mockSync.Sync{IsSyncing: true},
			HeadFetcher:           chainService,
			TimeFetcher:           chainService,
			OptimisticModeFetcher: chainService,
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte("foo")))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusServiceUnavailable, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Beacon node is currently syncing and not serving request on that endpoint"))
	})
}

func TestProduceBlockV3SSZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Run("Phase 0", func(t *testing.T) {
		var block *SignedBeaconBlock
		err := json.Unmarshal([]byte(phase0Block), &block)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.Message.ToGeneric()
			}())
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		g, err := block.ToGeneric()
		require.NoError(t, err)
		bl, ok := g.Block.(*eth.GenericSignedBeaconBlock_Phase0)
		require.Equal(t, true, ok)
		ssz, err := bl.Phase0.Block.MarshalSSZ()
		require.Equal(t, string(ssz), writer.Body.String())
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "false", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("Altair", func(t *testing.T) {
		var block *SignedBeaconBlockAltair
		err := json.Unmarshal([]byte(altairBlock), &block)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.Message.ToGeneric()
			}())
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		g, err := block.ToGeneric()
		require.NoError(t, err)
		bl, ok := g.Block.(*eth.GenericSignedBeaconBlock_Altair)
		require.Equal(t, true, ok)
		ssz, err := bl.Altair.Block.MarshalSSZ()
		require.Equal(t, string(ssz), writer.Body.String())
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "false", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("Bellatrix", func(t *testing.T) {
		var block *SignedBeaconBlockBellatrix
		err := json.Unmarshal([]byte(bellatrixBlock), &block)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.Message.ToGeneric()
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		g, err := block.ToGeneric()
		require.NoError(t, err)
		bl, ok := g.Block.(*eth.GenericSignedBeaconBlock_Bellatrix)
		require.Equal(t, true, ok)
		ssz, err := bl.Bellatrix.Block.MarshalSSZ()
		require.Equal(t, string(ssz), writer.Body.String())
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "false", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("BlindedBellatrix", func(t *testing.T) {
		var block *SignedBlindedBeaconBlockBellatrix
		err := json.Unmarshal([]byte(blindedBellatrixBlock), &block)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.Message.ToGeneric()
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		g, err := block.ToGeneric()
		require.NoError(t, err)
		bl, ok := g.Block.(*eth.GenericSignedBeaconBlock_BlindedBellatrix)
		require.Equal(t, true, ok)
		ssz, err := bl.BlindedBellatrix.Block.MarshalSSZ()
		require.Equal(t, string(ssz), writer.Body.String())
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "true", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("Capella", func(t *testing.T) {
		var block *SignedBeaconBlockCapella
		err := json.Unmarshal([]byte(capellaBlock), &block)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.Message.ToGeneric()
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		g, err := block.ToGeneric()
		require.NoError(t, err)
		bl, ok := g.Block.(*eth.GenericSignedBeaconBlock_Capella)
		require.Equal(t, true, ok)
		ssz, err := bl.Capella.Block.MarshalSSZ()
		require.Equal(t, string(ssz), writer.Body.String())
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "false", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("Blinded Capella", func(t *testing.T) {
		var block *SignedBlindedBeaconBlockCapella
		err := json.Unmarshal([]byte(blindedCapellaBlock), &block)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				g, err := block.Message.ToGeneric()
				require.NoError(t, err)
				g.PayloadValue = 2000 //some fake value
				return g, err
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		g, err := block.ToGeneric()
		require.NoError(t, err)
		bl, ok := g.Block.(*eth.GenericSignedBeaconBlock_BlindedCapella)
		require.Equal(t, true, ok)
		ssz, err := bl.BlindedCapella.Block.MarshalSSZ()
		require.Equal(t, string(ssz), writer.Body.String())
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "true", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "2000", true)
	})
	t.Run("Deneb", func(t *testing.T) {
		var block *SignedBeaconBlockContentsDeneb
		err := json.Unmarshal([]byte(denebBlockContents), &block)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.ToUnsigned().ToGeneric()
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		g, err := block.ToUnsigned().ToGeneric()
		require.NoError(t, err)
		bl, ok := g.Block.(*eth.GenericBeaconBlock_Deneb)
		require.Equal(t, true, ok)
		ssz, err := bl.Deneb.MarshalSSZ()
		require.Equal(t, string(ssz), writer.Body.String())
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "false", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
	t.Run("Blinded Deneb", func(t *testing.T) {
		var block *SignedBlindedBeaconBlockContentsDeneb
		err := json.Unmarshal([]byte(blindedDenebBlockContents), &block)
		require.NoError(t, err)
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().GetBeaconBlock(gomock.Any(), gomock.Any()).Return(
			func() (*eth.GenericBeaconBlock, error) {
				return block.ToUnsigned().ToGeneric()
			}())
		mockChainService := &blockchainTesting.ChainService{}
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
			OptimisticModeFetcher:   mockChainService,
		}
		rr := "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505" +
			"&graffiti=0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://foo.example/eth/v3/validator/blocks/1?randao_reveal=%s", rr), nil)
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.ProduceBlockV3(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		g, err := block.ToUnsigned().ToGeneric()
		require.NoError(t, err)
		bl, ok := g.Block.(*eth.GenericBeaconBlock_BlindedDeneb)
		require.Equal(t, true, ok)
		ssz, err := bl.BlindedDeneb.MarshalSSZ()
		require.Equal(t, string(ssz), writer.Body.String())
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadBlindedHeader) == "true", true)
		require.Equal(t, writer.Header().Get(api.ExecutionPayloadValueHeader) == "0", true)
	})
}

func TestValidateConsensus(t *testing.T) {
	ctx := context.Background()

	parentState, privs := util.DeterministicGenesisState(t, params.MinimalSpecConfig().MinGenesisActiveValidatorCount)
	parentBlock, err := util.GenerateFullBlock(parentState, privs, util.DefaultBlockGenConfig(), parentState.Slot())
	require.NoError(t, err)
	parentSbb, err := blocks.NewSignedBeaconBlock(parentBlock)
	require.NoError(t, err)
	st, err := transition.ExecuteStateTransition(ctx, parentState, parentSbb)
	require.NoError(t, err)
	block, err := util.GenerateFullBlock(st, privs, util.DefaultBlockGenConfig(), st.Slot())
	require.NoError(t, err)
	sbb, err := blocks.NewSignedBeaconBlock(block)
	require.NoError(t, err)
	parentRoot, err := parentSbb.Block().HashTreeRoot()
	require.NoError(t, err)
	server := &Server{
		Blocker: &testutil.MockBlocker{RootBlockMap: map[[32]byte]interfaces.ReadOnlySignedBeaconBlock{parentRoot: parentSbb}},
		Stater:  &testutil.MockStater{StatesByRoot: map[[32]byte]state.BeaconState{bytesutil.ToBytes32(parentBlock.Block.StateRoot): parentState}},
	}

	require.NoError(t, server.validateConsensus(ctx, sbb))
}

func TestValidateEquivocation(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		st, err := util.NewBeaconState()
		require.NoError(t, err)
		require.NoError(t, st.SetSlot(10))
		fc := doublylinkedtree.New()
		require.NoError(t, fc.InsertNode(context.Background(), st, bytesutil.ToBytes32([]byte("root"))))
		server := &Server{
			ForkchoiceFetcher: &testing2.ChainService{ForkChoiceStore: fc},
		}
		blk, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlock())
		require.NoError(t, err)
		blk.SetSlot(st.Slot() + 1)

		require.NoError(t, server.validateEquivocation(blk.Block()))
	})
	t.Run("block already exists", func(t *testing.T) {
		st, err := util.NewBeaconState()
		require.NoError(t, err)
		require.NoError(t, st.SetSlot(10))
		fc := doublylinkedtree.New()
		require.NoError(t, fc.InsertNode(context.Background(), st, bytesutil.ToBytes32([]byte("root"))))
		server := &Server{
			ForkchoiceFetcher: &testing2.ChainService{ForkChoiceStore: fc},
		}
		blk, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlock())
		require.NoError(t, err)
		blk.SetSlot(st.Slot())

		assert.ErrorContains(t, "already exists", server.validateEquivocation(blk.Block()))
	})
}

const (
	phase0Block = `{
  "message": {
    "slot": "1",
    "proposer_index": "1",
    "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "body": {
      "randao_reveal": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
      "eth1_data": {
        "deposit_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "deposit_count": "1",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
      },
      "graffiti": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
      "proposer_slashings": [
        {
          "signed_header_1": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          },
          "signed_header_2": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "attester_slashings": [
        {
          "attestation_1": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          },
          "attestation_2": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          }
        }
      ],
      "attestations": [
        {
          "aggregation_bits": "0xffffffffffffffffffffffffffffffffff3f",
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
          "data": {
            "slot": "1",
            "index": "1",
            "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "source": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "target": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            }
          }
        }
      ],
      "deposits": [
        {
          "proof": [
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
			"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
          ],
          "data": {
            "pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
            "withdrawal_credentials": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "amount": "1",
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "voluntary_exits": [
        {
          "message": {
            "epoch": "1",
            "validator_index": "1"
          },
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
        }
      ]
    }
  },
  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
}`
	altairBlock = `{
  "message": {
    "slot": "1",
    "proposer_index": "1",
    "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "body": {
      "randao_reveal": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
      "eth1_data": {
        "deposit_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "deposit_count": "1",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
      },
      "graffiti": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
      "proposer_slashings": [
        {
          "signed_header_1": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          },
          "signed_header_2": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "attester_slashings": [
        {
          "attestation_1": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          },
          "attestation_2": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          }
        }
      ],
      "attestations": [
        {
          "aggregation_bits": "0xffffffffffffffffffffffffffffffffff3f",
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
          "data": {
            "slot": "1",
            "index": "1",
            "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "source": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "target": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            }
          }
        }
      ],
      "deposits": [
        {
          "proof": [
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
			"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
          ],
          "data": {
            "pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
            "withdrawal_credentials": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "amount": "1",
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "voluntary_exits": [
        {
          "message": {
            "epoch": "1",
            "validator_index": "1"
          },
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
        }
      ],
      "sync_aggregate": {
        "sync_committee_bits": "0x6451e9f951ebf05edc01de67e593484b672877054f055903ff0df1a1a945cf30ca26bb4d4b154f94a1bc776bcf5d0efb3603e1f9b8ee2499ccdcfe2a18cef458",
        "sync_committee_signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
      }
    }
  },
  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
}`
	bellatrixBlock = `{
  "message": {
    "slot": "1",
    "proposer_index": "1",
    "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "body": {
      "randao_reveal": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
      "eth1_data": {
        "deposit_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "deposit_count": "1",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
      },
      "graffiti": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
      "proposer_slashings": [
        {
          "signed_header_1": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          },
          "signed_header_2": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "attester_slashings": [
        {
          "attestation_1": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          },
          "attestation_2": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          }
        }
      ],
      "attestations": [
        {
          "aggregation_bits": "0xffffffffffffffffffffffffffffffffff3f",
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
          "data": {
            "slot": "1",
            "index": "1",
            "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "source": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "target": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            }
          }
        }
      ],
      "deposits": [
        {
          "proof": [
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
 			"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
          ],
          "data": {
            "pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
            "withdrawal_credentials": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "amount": "1",
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "voluntary_exits": [
        {
          "message": {
            "epoch": "1",
            "validator_index": "1"
          },
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
        }
      ],
      "sync_aggregate": {
        "sync_committee_bits": "0x6451e9f951ebf05edc01de67e593484b672877054f055903ff0df1a1a945cf30ca26bb4d4b154f94a1bc776bcf5d0efb3603e1f9b8ee2499ccdcfe2a18cef458",
        "sync_committee_signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
      },
      "execution_payload": {
        "parent_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "fee_recipient": "0xabcf8e0d4e9587369b2301d0790347320302cc09",
        "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "receipts_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "logs_bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "prev_randao": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "block_number": "1",
        "gas_limit": "1",
        "gas_used": "1",
        "timestamp": "1",
        "extra_data": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "base_fee_per_gas": "14074904626401341155369551180448584754667373453244490859944217516317499064576",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "transactions": [
          "0x02f878831469668303f51d843b9ac9f9843b9aca0082520894c93269b73096998db66be0441e836d873535cb9c8894a19041886f000080c001a031cc29234036afbf9a1fb9476b463367cb1f957ac0b919b69bbc798436e604aaa018c4e9c3914eb27aadd0b91e10b18655739fcf8c1fc398763a9f1beecb8ddc86"
        ]
      }
    }
  },
  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
}`
	blindedBellatrixBlock = `{
  "message": {
    "slot": "1",
    "proposer_index": "1",
    "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "body": {
      "randao_reveal": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
      "eth1_data": {
        "deposit_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "deposit_count": "1",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
      },
      "graffiti": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
      "proposer_slashings": [
        {
          "signed_header_1": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          },
          "signed_header_2": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "attester_slashings": [
        {
          "attestation_1": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          },
          "attestation_2": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          }
        }
      ],
      "attestations": [
        {
          "aggregation_bits": "0xffffffffffffffffffffffffffffffffff3f",
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
          "data": {
            "slot": "1",
            "index": "1",
            "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "source": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "target": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            }
          }
        }
      ],
      "deposits": [
        {
          "proof": [
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
			"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
          ],
          "data": {
            "pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
            "withdrawal_credentials": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "amount": "1",
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "voluntary_exits": [
        {
          "message": {
            "epoch": "1",
            "validator_index": "1"
          },
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
        }
      ],
      "sync_aggregate": {
        "sync_committee_bits": "0x6451e9f951ebf05edc01de67e593484b672877054f055903ff0df1a1a945cf30ca26bb4d4b154f94a1bc776bcf5d0efb3603e1f9b8ee2499ccdcfe2a18cef458",
        "sync_committee_signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
      },
      "execution_payload_header": {
        "parent_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "fee_recipient": "0xabcf8e0d4e9587369b2301d0790347320302cc09",
        "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "receipts_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "logs_bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "prev_randao": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "block_number": "1",
        "gas_limit": "1",
        "gas_used": "1",
        "timestamp": "1",
        "extra_data": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "base_fee_per_gas": "1",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "transactions_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
      }
    }
  },
  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
}`
	capellaBlock = `{
  "message": {
    "slot": "1",
    "proposer_index": "1",
    "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "body": {
      "randao_reveal": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
      "eth1_data": {
        "deposit_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "deposit_count": "1",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
      },
      "graffiti": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
      "proposer_slashings": [
        {
          "signed_header_1": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          },
          "signed_header_2": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "attester_slashings": [
        {
          "attestation_1": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          },
          "attestation_2": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          }
        }
      ],
      "attestations": [
        {
          "aggregation_bits": "0xffffffffffffffffffffffffffffffffff3f",
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
          "data": {
            "slot": "1",
            "index": "1",
            "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "source": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "target": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            }
          }
        }
      ],
      "deposits": [
        {
          "proof": [
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
			"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
          ],
          "data": {
            "pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
            "withdrawal_credentials": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "amount": "1",
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "voluntary_exits": [
        {
          "message": {
            "epoch": "1",
            "validator_index": "1"
          },
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
        }
      ],
      "sync_aggregate": {
        "sync_committee_bits": "0x6451e9f951ebf05edc01de67e593484b672877054f055903ff0df1a1a945cf30ca26bb4d4b154f94a1bc776bcf5d0efb3603e1f9b8ee2499ccdcfe2a18cef458",
        "sync_committee_signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
      },
      "execution_payload": {
        "parent_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "fee_recipient": "0xabcf8e0d4e9587369b2301d0790347320302cc09",
        "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "receipts_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "logs_bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "prev_randao": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "block_number": "1",
        "gas_limit": "1",
        "gas_used": "1",
        "timestamp": "1",
        "extra_data": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "base_fee_per_gas": "14074904626401341155369551180448584754667373453244490859944217516317499064576",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "transactions": [
          "0x02f878831469668303f51d843b9ac9f9843b9aca0082520894c93269b73096998db66be0441e836d873535cb9c8894a19041886f000080c001a031cc29234036afbf9a1fb9476b463367cb1f957ac0b919b69bbc798436e604aaa018c4e9c3914eb27aadd0b91e10b18655739fcf8c1fc398763a9f1beecb8ddc86"
        ],
        "withdrawals": [
          {
            "index": "1",
            "validator_index": "1",
            "address": "0xabcf8e0d4e9587369b2301d0790347320302cc09",
            "amount": "1"
          }
        ]
      },
      "bls_to_execution_changes": [
        {
          "message": {
            "validator_index": "1",
            "from_bls_pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
            "to_execution_address": "0xabcf8e0d4e9587369b2301d0790347320302cc09"
          },
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
        }
      ]
    }
  },
  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
}`
	blindedCapellaBlock = `{
  "message": {
    "slot": "1",
    "proposer_index": "1",
    "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "body": {
      "randao_reveal": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
      "eth1_data": {
        "deposit_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "deposit_count": "1",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
      },
      "graffiti": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
      "proposer_slashings": [
        {
          "signed_header_1": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          },
          "signed_header_2": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "attester_slashings": [
        {
          "attestation_1": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          },
          "attestation_2": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          }
        }
      ],
      "attestations": [
        {
          "aggregation_bits": "0xffffffffffffffffffffffffffffffffff3f",
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
          "data": {
            "slot": "1",
            "index": "1",
            "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "source": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "target": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            }
          }
        }
      ],
      "deposits": [
        {
          "proof": [
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
			"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
          ],
          "data": {
            "pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
            "withdrawal_credentials": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "amount": "1",
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "voluntary_exits": [
        {
          "message": {
            "epoch": "1",
            "validator_index": "1"
          },
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
        }
      ],
      "sync_aggregate": {
        "sync_committee_bits": "0x6451e9f951ebf05edc01de67e593484b672877054f055903ff0df1a1a945cf30ca26bb4d4b154f94a1bc776bcf5d0efb3603e1f9b8ee2499ccdcfe2a18cef458",
        "sync_committee_signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
      },
      "execution_payload_header": {
        "parent_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "fee_recipient": "0xabcf8e0d4e9587369b2301d0790347320302cc09",
        "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "receipts_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "logs_bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "prev_randao": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "block_number": "1",
        "gas_limit": "1",
        "gas_used": "1",
        "timestamp": "1",
        "extra_data": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "base_fee_per_gas": "14074904626401341155369551180448584754667373453244490859944217516317499064576",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "transactions_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "withdrawals_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
      },
      "bls_to_execution_changes": [
        {
          "message": {
            "validator_index": "1",
            "from_bls_pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
            "to_execution_address": "0xabcf8e0d4e9587369b2301d0790347320302cc09"
          },
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
        }
      ]
    }
  },
  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
}`
	blindedDenebBlockContents = `{
	"signed_blinded_block":{
  "message": {
    "slot": "1",
    "proposer_index": "1",
    "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
    "body": {
      "randao_reveal": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
      "eth1_data": {
        "deposit_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "deposit_count": "1",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
      },
      "graffiti": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
      "proposer_slashings": [
        {
          "signed_header_1": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          },
          "signed_header_2": {
            "message": {
              "slot": "1",
              "proposer_index": "1",
              "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "attester_slashings": [
        {
          "attestation_1": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          },
          "attestation_2": {
            "attesting_indices": [
              "1"
            ],
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
            "data": {
              "slot": "1",
              "index": "1",
              "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
              "source": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              },
              "target": {
                "epoch": "1",
                "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
              }
            }
          }
        }
      ],
      "attestations": [
        {
          "aggregation_bits": "0xffffffffffffffffffffffffffffffffff3f",
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
          "data": {
            "slot": "1",
            "index": "1",
            "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "source": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            },
            "target": {
              "epoch": "1",
              "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
            }
          }
        }
      ],
      "deposits": [
        {
          "proof": [
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
			"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
          ],
          "data": {
            "pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
            "withdrawal_credentials": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
            "amount": "1",
            "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
          }
        }
      ],
      "voluntary_exits": [
        {
          "message": {
            "epoch": "1",
            "validator_index": "1"
          },
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
        }
      ],
      "sync_aggregate": {
        "sync_committee_bits": "0x6451e9f951ebf05edc01de67e593484b672877054f055903ff0df1a1a945cf30ca26bb4d4b154f94a1bc776bcf5d0efb3603e1f9b8ee2499ccdcfe2a18cef458",
        "sync_committee_signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
      },
      "execution_payload_header": {
        "parent_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "fee_recipient": "0xabcf8e0d4e9587369b2301d0790347320302cc09",
        "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "receipts_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "logs_bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "prev_randao": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "block_number": "1",
        "gas_limit": "1",
        "gas_used": "1",
        "timestamp": "1",
        "extra_data": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "base_fee_per_gas": "14074904626401341155369551180448584754667373453244490859944217516317499064576",
        "block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
		"blob_gas_used": "1",
		"excess_blob_gas": "2",
        "transactions_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
        "withdrawals_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
      },
      "bls_to_execution_changes": [
        {
          "message": {
            "validator_index": "1",
            "from_bls_pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
            "to_execution_address": "0xabcf8e0d4e9587369b2301d0790347320302cc09"
          },
          "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
        }
      ],
 	  "blob_kzg_commitments":["0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8000"]
    }
  },
  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
},
"signed_blinded_blob_sidecars":[{
			"message":{
				"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"index":"1",
				"slot":"1",
				"block_parent_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"proposer_index":"1",
				"blob_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8000",
				"kzg_proof":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8000"
			},
			"signature":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
		}]
}`
)

var denebBlockContents = `{
 	"signed_block":{
		 "message": {
			"slot": "1",
			"proposer_index": "1",
			"parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
			"state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
			"body": {
			  "randao_reveal": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
			  "eth1_data": {
				"deposit_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"deposit_count": "1",
				"block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
			  },
			  "graffiti": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
			  "proposer_slashings": [
				{
				  "signed_header_1": {
					"message": {
					  "slot": "1",
					  "proposer_index": "1",
					  "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					  "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					  "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
					},
					"signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
				  },
				  "signed_header_2": {
					"message": {
					  "slot": "1",
					  "proposer_index": "1",
					  "parent_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					  "state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					  "body_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
					},
					"signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
				  }
				}
			  ],
			  "attester_slashings": [
				{
				  "attestation_1": {
					"attesting_indices": [
					  "1"
					],
					"signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
					"data": {
					  "slot": "1",
					  "index": "1",
					  "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					  "source": {
						"epoch": "1",
						"root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
					  },
					  "target": {
						"epoch": "1",
						"root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
					  }
					}
				  },
				  "attestation_2": {
					"attesting_indices": [
					  "1"
					],
					"signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
					"data": {
					  "slot": "1",
					  "index": "1",
					  "beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					  "source": {
						"epoch": "1",
						"root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
					  },
					  "target": {
						"epoch": "1",
						"root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
					  }
					}
				  }
				}
			  ],
			  "attestations": [
				{
				  "aggregation_bits": "0x01",
				  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",
				  "data": {
					"slot": "1",
					"index": "1",
					"beacon_block_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"source": {
					  "epoch": "1",
					  "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
					},
					"target": {
					  "epoch": "1",
					  "root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
					}
				  }
				}
			  ],
			  "deposits": [
				{
				  "proof": [
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"
				  ],
				  "data": {
					"pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
					"withdrawal_credentials": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
					"amount": "1",
					"signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
				  }
				}
			  ],
			  "voluntary_exits": [
				{
				  "message": {
					"epoch": "1",
					"validator_index": "1"
				  },
				  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
				}
			  ],
			  "sync_aggregate": {
				"sync_committee_bits": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"sync_committee_signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
			  },
			  "execution_payload": {
				"parent_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"fee_recipient": "0xabcf8e0d4e9587369b2301d0790347320302cc09",
				"state_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"receipts_root": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"logs_bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"prev_randao": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"block_number": "1",
				"gas_limit": "1",
				"gas_used": "1",
				"timestamp": "1",
				"extra_data": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"base_fee_per_gas": "1",
				"block_hash": "0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"blob_gas_used": "1",
				"excess_blob_gas": "2",
				"transactions": [
				  "0x02f878831469668303f51d843b9ac9f9843b9aca0082520894c93269b73096998db66be0441e836d873535cb9c8894a19041886f000080c001a031cc29234036afbf9a1fb9476b463367cb1f957ac0b919b69bbc798436e604aaa018c4e9c3914eb27aadd0b91e10b18655739fcf8c1fc398763a9f1beecb8ddc86"
				],
				"withdrawals": [
				  {
					"index": "1",
					"validator_index": "1",
					"address": "0xabcf8e0d4e9587369b2301d0790347320302cc09",
					"amount": "1"
				  }
				]
			  },
			  "bls_to_execution_changes": [
				{
				  "message": {
					"validator_index": "1",
					"from_bls_pubkey": "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
					"to_execution_address": "0xabcf8e0d4e9587369b2301d0790347320302cc09"
				  },
				  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
				}
			  ],
			  "blob_kzg_commitments":["0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8000"]
			}
		  },
		  "signature": "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
		},
		"signed_blob_sidecars":[{
			"message":{
				"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"index":"1",
				"slot":"1",
				"block_parent_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2",
				"proposer_index":"1",
				"blob":"` + hexutil.Encode(make([]byte, 131072)) + `" ,
				"kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8000",
				"kzg_proof":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8000"
			},
			"signature":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"
		}]
}`
