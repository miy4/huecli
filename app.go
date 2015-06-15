package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jawher/mow.cli"
	"github.com/miy4/huecli/api"
	"github.com/miy4/huecli/color"
)

const (
	ExitCodeOK = iota
	ExitCodeRegisterCommandError
	ExitCodeOnCommandError
	ExitCodeOffCommandError
)

const (
	configFile = ".hue.json"
)

type Command struct {
	Name string
	Desc string
	Init cli.CmdInitializer
}

var Commands = []Command{
	commandRegister,
	commandOn,
	commandOff,
}

var commandRegister = Command{
	"register",
	"register new user",
	func(cmd *cli.Cmd) {
		cmd.Action = func() {
			if err := registerUser(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(ExitCodeRegisterCommandError)
			}
		}
	},
}

var commandOn = Command{
	"on",
	"turn on the lights",
	func(cmd *cli.Cmd) {
		color := cmd.StringOpt("c color", "", "the color, where value from #000000 to #FFFFFF.")
		brightness := cmd.IntOpt("b bri", -1, "the brightness, where value from 0 to 255.")
		hue := cmd.IntOpt("h hue", -1, "the hue, where value from 0 to 65535.")
		saturation := cmd.IntOpt("s sat", -1, "the saturation, where value from 0 to 255.")
		temperature := cmd.IntOpt("t ct", -1, "the color temperature, where value from 153 to 500.")

		cmd.Action = func() {
			state := api.LightState{
				Color:       *color,
				Brightness:  *brightness,
				Hue:         *hue,
				Saturation:  *saturation,
				Temperature: *temperature,
			}

			if err := turnOnLights(&state); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(ExitCodeOnCommandError)
			}
		}
	},
}

var commandOff = Command{
	"off",
	"turn off the lights",
	func(cmd *cli.Cmd) {
		cmd.Action = func() {
			if err := turnOffLights(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(ExitCodeOffCommandError)
			}
		}
	},
}

func exportConfig(hue *api.Hue, configPath string) error {
	file, err := json.Marshal(*hue)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(configPath, file, 0600); err != nil {
		return err
	}

	return nil
}

func importConfig(configPath string) (*api.Hue, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var hue api.Hue
	if err := json.Unmarshal(file, &hue); err != nil {
		return nil, err
	}

	return &hue, nil
}

func registerUser() error {
	var hue api.Hue
	ipAddr, err := hue.GetBridgeIP()
	if err != nil {
		return err
	}

	hue.IpAddress = ipAddr
	user, err := hue.RegisterUser()
	if err != nil {
		return err
	}

	hue.UserName = user
	configPath := filepath.Join(os.Getenv("HOME"), configFile)
	err = exportConfig(&hue, configPath)
	if err != nil {
		return err
	}

	return nil
}

func lightStateToParams(state *api.LightState) (map[string]interface{}, error) {
	params := make(map[string]interface{})

	params["on"] = state.On

	if state.Color != "" {
		x, y, err := color.RGBToXY(state.Color)
		if err != nil {
			return nil, err
		}
		params["xy"] = [2]float64{x, y}
	}

	if state.Brightness != -1 {
		if state.Brightness < 0 || state.Brightness > 255 {
			return nil, fmt.Errorf("brightness out of range: %s", state.Brightness)
		}
		params["bri"] = state.Brightness
	}

	if state.Hue != -1 {
		if state.Hue < 0 || state.Hue > 65535 {
			return nil, fmt.Errorf("hue out of range: %s", state.Hue)
		}
		params["hue"] = state.Hue
	}

	if state.Saturation != -1 {
		if state.Saturation < 0 || state.Saturation > 255 {
			return nil, fmt.Errorf("saturation out of range: %s", state.Saturation)
		}
		params["sat"] = state.Saturation
	}

	if state.Temperature != -1 {
		if state.Temperature < 153 || state.Temperature > 500 {
			return nil, fmt.Errorf("color temperature out of range: %s", state.Temperature)
		}
		params["ct"] = state.Temperature
	}

	return params, nil
}

func turnOnLights(state *api.LightState) error {
	configPath := filepath.Join(os.Getenv("HOME"), configFile)
	hue, err := importConfig(configPath)
	if err != nil {
		return err
	}

	lights, err := hue.GetAllLights()
	if err != nil {
		return err
	}

	state.On = true
	params, err := lightStateToParams(state)
	if err != nil {
		return nil
	}

	for _, light := range lights {
		err = hue.SetLightState(light, params)
		if err != nil {
			return err
		}
	}

	return nil
}

func turnOffLights() error {
	configPath := filepath.Join(os.Getenv("HOME"), configFile)
	hue, err := importConfig(configPath)
	if err != nil {
		return err
	}

	lights, err := hue.GetAllLights()
	if err != nil {
		return err
	}

	params := map[string]interface{}{"on": false}
	for _, light := range lights {
		err = hue.SetLightState(light, params)
		if err != nil {
			return err
		}
	}

	return nil
}
