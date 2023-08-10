package beacon

type Phase0ProduceBlockV3Response struct {
	Version                 string       `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool         `json:"execution_payload_blinded"`
	ExeuctionPayloadValue   string       `json:"exeuction_payload_value"`
	Data                    *BeaconBlock `json:"data" validate:"required"`
}

type AltairProduceBlockV3Response struct {
	Version                 string             `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool               `json:"execution_payload_blinded"`
	ExeuctionPayloadValue   string             `json:"exeuction_payload_value"`
	Data                    *BeaconBlockAltair `json:"data" validate:"required"`
}

type BellatrixProduceBlockV3Response struct {
	Version                 string                `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                  `json:"execution_payload_blinded"`
	ExeuctionPayloadValue   string                `json:"exeuction_payload_value"`
	Data                    *BeaconBlockBellatrix `json:"data" validate:"required"`
}

type BlindedBellatrixProduceBlockV3Response struct {
	Version                 string                       `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                         `json:"execution_payload_blinded"`
	ExeuctionPayloadValue   string                       `json:"exeuction_payload_value"`
	Data                    *BlindedBeaconBlockBellatrix `json:"data" validate:"required"`
}

type CapellaProduceBlockV3Response struct {
	Version                 string              `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                `json:"execution_payload_blinded"`
	ExeuctionPayloadValue   string              `json:"exeuction_payload_value"`
	Data                    *BeaconBlockCapella `json:"data" validate:"required"`
}

type BlindedCapellaProduceBlockV3Response struct {
	Version                 string                     `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                       `json:"execution_payload_blinded"`
	ExeuctionPayloadValue   string                     `json:"exeuction_payload_value"`
	Data                    *BlindedBeaconBlockCapella `json:"data" validate:"required"`
}

type DenebProduceBlockV3Response struct {
	Version                 string                    `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                      `json:"execution_payload_blinded"`
	ExeuctionPayloadValue   string                    `json:"exeuction_payload_value"`
	Data                    *BeaconBlockContentsDeneb `json:"data" validate:"required"`
}

type BlindedDenebProduceBlockV3Response struct {
	Version                 string                           `json:"version" validate:"required"`
	ExecutionPayloadBlinded bool                             `json:"execution_payload_blinded"`
	ExeuctionPayloadValue   string                           `json:"exeuction_payload_value"`
	Data                    *BlindedBeaconBlockContentsDeneb `json:"data" validate:"required"`
}
