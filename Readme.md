# led-schedule

This is a simple daylight scheduler daemon for LEDs driven by a [rpi-ws2812-server](https://github.com/tom-2015/rpi-ws2812-server).

## command

```
Usage: led-schedule [--before BEFORE] [--power POWER] [--lat LAT] [--lon LON] [--fade-interval FADE-INTERVAL] [--address ADDRESS]

Options:
  --before BEFORE        duration before sunset to turn on [default: 10m]
  --power POWER          power level [default: 140]
  --lat LAT              latitude
  --lon LON              longitude
  --fade-interval FADE-INTERVAL
                         fade interval [default: 25ms]
  --address ADDRESS      address of the light controller [default: 10.3.1.140:9999]
  --help, -h             display this help and exit
```
