package solidmation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/matiasgarcia78/bgh-home-bridge/src/api/solidmation/dto"
)

type Auth struct {
	User     string `json:"user"`
	Password string `json:"password"`
	token    string `json:"-"`
}

type SolidmationApi struct {
	auth Auth
}

var endPointIDCache map[DeviceID]EndPointID
var homesID []HomeID

type DeviceID uint64
type EndPointID uint64
type HomeID uint64

const baseUrl = "https://bgh-services.solidmation.com%s"

func NewSolidmationApi(pAuth Auth) SolidmationApi {
	endPointIDCache = make(map[DeviceID]EndPointID)
	return SolidmationApi{pAuth}
}

func (api *SolidmationApi) login() error {
	request, err := json.Marshal(api.auth)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(baseUrl, "/control/LoginPage.aspx/DoStandardLogin")
	r, err := http.Post(url, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return err
	}
	data, _ := ioutil.ReadAll(r.Body)

	response := struct {
		Token string `json:"d"`
	}{}

	if err = json.Unmarshal(data, &response); err != nil {
		return nil
	}

	api.auth.token = response.Token

	return api.setEndPoints()
}

func (api *SolidmationApi) GetDeviceStatus(deviceID DeviceID, valueType uint64) (string, error) {
	if len(api.auth.token) == 0 {
		if err := api.login(); err != nil {
			return "", err
		}
	}

	epID := endPointIDCache[deviceID]

	var value string
	for _, h := range homesID {
		req := dto.GetPacketRequest{
			Token:  dto.Token{Token: api.auth.token},
			HomeID: uint64(h),
		}

		request, err := json.Marshal(req)
		if err != nil {
			return value, err
		}

		url := fmt.Sprintf(baseUrl, "/1.0/HomeCloudService.svc/GetDataPacket")

		r, err := http.Post(url, "application/json", bytes.NewBuffer(request))
		if err != nil {
			return value, err
		}
		data, _ := ioutil.ReadAll(r.Body)

		var resp *dto.GetDataPacketResponse
		if err = json.Unmarshal(data, &resp); err != nil {
			return value, nil
		}

		var ok bool
		for _, ep := range resp.GetDataPacketResult.EndPointValues {
			if ok = EndPointID(ep.EndPointID) == epID; ok {
				for _, vt := range ep.Values {
					if vt.ValueType == valueType {
						value = vt.Value
						break
					}
				}
				break
			}
		}
		if ok {
			break
		}
	}
	return value, nil
}

func (api *SolidmationApi) setEndPoints() error {
	url := fmt.Sprintf(baseUrl, "/1.0/HomeCloudService.svc/EnumHomes")

	req := struct {
		Token dto.Token `json:"token"`
	}{Token: dto.Token{Token: api.auth.token}}

	request, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := http.Post(url, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return err
	}

	data, _ := ioutil.ReadAll(r.Body)

	var response *dto.EnumHomeResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return nil
	}

	url = fmt.Sprintf(baseUrl, "/1.0/HomeCloudService.svc/GetDataPacket")

	for _, h := range response.EnumHomesResult.Homes {
		req := dto.GetPacketRequest{
			Token:  dto.Token{Token: api.auth.token},
			HomeID: h.HomeID,
		}

		homesID = append(homesID, HomeID(h.HomeID))

		request, err := json.Marshal(req)
		if err != nil {
			return err
		}

		r, err := http.Post(url, "application/json", bytes.NewBuffer(request))
		if err != nil {
			return err
		}
		data, _ := ioutil.ReadAll(r.Body)

		var resp *dto.GetDataPacketResponse
		if err = json.Unmarshal(data, &resp); err != nil {
			return nil
		}

		for _, ep := range resp.GetDataPacketResult.EndPoints {
			endPointIDCache[DeviceID(ep.DeviceID)] = EndPointID(ep.EndPointID)
		}

	}

	return nil
}

func (api *SolidmationApi) SetDeviceStatus(deviceID DeviceID, temperature uint64, mode string) error {
	if len(api.auth.token) == 0 {
		if err := api.login(); err != nil {
			return err
		}
	}

	req := struct {
		Token       dto.Token  `json:"token"`
		EndpointID  EndPointID `json:"endpointID"`
		Mode        string     `json:"mode"`
		Temperature uint64     `json:"desiredTempC"`
	}{
		Token:       dto.Token{Token: api.auth.token},
		EndpointID:  endPointIDCache[deviceID],
		Mode:        mode,
		Temperature: 24,
	}

	if temperature > 0 {
		req.Temperature = temperature
	}

	request, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(baseUrl, "/1.0/HomeCloudCommandService.svc/HVACSetModes")
	if _, err = http.Post(url, "application/json", bytes.NewBuffer(request)); err != nil {
		return err
	}

	return nil
}

func (api *SolidmationApi) GetStatus() string {
	url := fmt.Sprintf(baseUrl, "/1.0/HomeCloudCommandService.svc/Ping")
	if _, err := http.Post(url, "application/json", nil); err != nil {
		return "offline"
	}

	return "online"
}
