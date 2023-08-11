package beacon

type Phase0ProduceBlockV3Response struct {
	Version                 string       `json:"version"`
	ExecutionPayloadBlinded bool         `json:"execution_payload_blinded"`
	ExecutionPayloadValue   string       `json:"execution_payload_value"`
	Data                    *BeaconBlock `json:"data"`
}

type AltairProduceBlockV3Response struct {
	Version                 string             `json:"version"`
	ExecutionPayloadBlinded bool               `json:"execution_payload_blinded"`
	ExecutionPayloadValue   string             `json:"execution_payload_value"`
	Data                    *BeaconBlockAltair `json:"data"`
}

type BellatrixProduceBlockV3Response struct {
	Version                 string                `json:"version"`
	ExecutionPayloadBlinded bool                  `json:"execution_payload_blinded"`
	ExecutionPayloadValue   string                `json:"execution_payload_value"`
	Data                    *BeaconBlockBellatrix `json:"data"`
}

type BlindedBellatrixProduceBlockV3Response struct {
	Version                 string                       `json:"version"`
	ExecutionPayloadBlinded bool                         `json:"execution_payload_blinded"`
	ExecutionPayloadValue   string                       `json:"execution_payload_value"`
	Data                    *BlindedBeaconBlockBellatrix `json:"data"`
}

type CapellaProduceBlockV3Response struct {
	Version                 string              `json:"version"`
	ExecutionPayloadBlinded bool                `json:"execution_payload_blinded"`
	ExecutionPayloadValue   string              `json:"execution_payload_value"`
	Data                    *BeaconBlockCapella `json:"data"`
}

type BlindedCapellaProduceBlockV3Response struct {
	Version                 string                     `json:"version"`
	ExecutionPayloadBlinded bool                       `json:"execution_payload_blinded"`
	ExecutionPayloadValue   string                     `json:"execution_payload_value"`
	Data                    *BlindedBeaconBlockCapella `json:"data"`
}

type DenebProduceBlockV3Response struct {
	Version                 string                    `json:"version"`
	ExecutionPayloadBlinded bool                      `json:"execution_payload_blinded"`
	ExecutionPayloadValue   string                    `json:"execution_payload_value"`
	Data                    *BeaconBlockContentsDeneb `json:"data"`
}

type BlindedDenebProduceBlockV3Response struct {
	Version                 string                           `json:"version"`
	ExecutionPayloadBlinded bool                             `json:"execution_payload_blinded"`
	ExecutionPayloadValue   string                           `json:"execution_payload_value"`
	Data                    *BlindedBeaconBlockContentsDeneb `json:"data"`
}
