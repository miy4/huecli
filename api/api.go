package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Hue struct {
	IpAddress string `json:"ipAddress"`
	UserName  string `json:"userName"`
}

type Light struct {
	Id string
}

type LightState struct {
	On          bool
	Color       string
	Brightness  int
	Hue         int
	Saturation  int
	Temperature int
	Effect      string
}

type Bridge struct {
	Serial     string `json:"id"`
	IpAddrress string `json:"internalipaddress"`
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, 2*time.Second)
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(2 * time.Second))
			return conn, nil
		},
	},
}

func nupnpSearch() (*[]Bridge, error) {
	response, err := http.Get("https://www.meethue.com/api/nupnp")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var bridges []Bridge
	err = json.NewDecoder(response.Body).Decode(&bridges)
	if err != nil {
		return nil, err
	}

	return &bridges, nil
}

func (hue *Hue) GetBridgeIP() (string, error) {
	b, err := nupnpSearch()
	if err != nil {
		return "", err
	}

	bridges := *b
	if len(bridges) <= 0 {
		return "", errors.New("bridges not found")
	}

	return bridges[0].IpAddrress, nil
}

func (hue *Hue) RegisterUser() (string, error) {
	params := map[string]string{"devicetype": "Golang API: huecli"}
	requestBody, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	uri := fmt.Sprintf("http://%s/api", hue.IpAddress)
	response, err := httpClient.Post(uri, "text/json", bytes.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var results []map[string]map[string]string
	json.NewDecoder(response.Body).Decode(&results)
	username := results[0]["success"]["username"]

	return username, nil
}

//func (hue *Hue) SetLightState(state LightState) error {
func (hue *Hue) SetLightState(light *Light, params map[string]interface{}) error {
	requestBody, err := json.Marshal(params)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("http://%s/api/%s/lights/%s/state", hue.IpAddress, hue.UserName, light.Id)
	httpRequest, err := http.NewRequest("PUT", uri, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	response, err := httpClient.Do(httpRequest)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var results []map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&results)
	return err
}

func (hue *Hue) GetAllLights() ([]*Light, error) {
	uri := fmt.Sprintf("http://%s/api/%s/lights", hue.IpAddress, hue.UserName)
	response, err := httpClient.Get(uri)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var results map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	var lights []*Light
	for id := range results {
		light := Light{Id: id}
		lights = append(lights, &light)
	}

	return lights, nil
}
