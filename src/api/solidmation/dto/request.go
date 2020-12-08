package dto

type Token struct {
	Token string `json:"Token"`
}

type GetPacketRequest struct {
	Token   Token   `json:"token"`
	HomeID  uint64  `json:"homeID"`
	Serials Serials `json:"serials"`
}

type Serials struct {
	Home           uint64 `json:"Home"`
	Groups         uint64 `json:"Groups"`
	Devices        uint64 `json:"Devices"`
	EndPoints      uint64 `json:"Endpoints"`
	EndPointValues uint64 `json:"EndpointValues"`
	Scenes         uint64 `json:"Scenes"`
	Macros         uint64 `json:"Macros"`
	Alarms         uint64 `json:"Alarms"`
}
