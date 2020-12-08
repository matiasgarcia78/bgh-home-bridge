package dto

type EnumHomeResponse struct {
	EnumHomesResult EnumHomesResult `json:"EnumHomesResult"`
}

type EnumHomesResult struct {
	Homes []Home `json:"Homes"`
}

type Home struct {
	HomeID uint64 `json:"HomeID"`
}

type GetDataPacketResponse struct {
	GetDataPacketResult GetDataPacketResult `json:"GetDataPacketResult"`
}

type GetDataPacketResult struct {
	EndPoints      []EndPoint      `json:"Endpoints"`
	EndPointValues []EndPointValue `json:"EndpointValues"`
}

type EndPoint struct {
	DeviceID   uint64 `json:"DeviceID"`
	EndPointID uint64 `json:"EndpointID"`
}

type EndPointValue struct {
	EndPointID uint64      `json:"EndpointID"`
	Values     []TypeValue `json:"Values"`
}

type TypeValue struct {
	Value     string `json:"Value"`
	ValueType uint64 `json:"ValueType"`
}
