# huecli
A simple CLI for Philips Hue. Just to flip it on and off.

## Install

```bash
$ go get github.com/miy4/huecli
```

## Commands

### register

Register new user.

```bash
$ huecli register
```

You must go to the bridge, press the button and then, run `huecli register` within 30 seconds.  
Registered user information will be saved to `~/.hue.json`.

### on

Turn on all the lights.

```bash
$ huecli on [--color COLOR] [--bri BRIGHTNESS] [--hue HUE] [--sat SATURATION] [--ct TEMPERATURE]
```

- `COLOR` the color, where value from #000000 to #FFFFFF.
- `BRIGHTNESS` the brightness, where value from 0 to 255.
- `HUE` the hue, where value from 0 to 65535.
- `SATURATION` the saturation where value from 0 to 255.
- `TEMPERATURE` the color temperature where value from 153 to 500.

### off

Turn off all of the lights.

```bash
$ huecli off
```
