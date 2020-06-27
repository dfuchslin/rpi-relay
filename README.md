# rpi-relay

Original code and idea from Adam Argo's [traincon](https://github.com/ap0/traincon), inspired by this [Medium post](https://medium.com/@adamargo/how-i-use-raspberry-pis-for-homekit-automation-ed69df8d8be7).

HTTP API to control RaspberryPi GPIO-based relays. Supports on, off, status.

```
GOOS=linux GOARM=7 GOARCH=arm go build
./rpi-relay config.yml
```

Example config file:
```
port: 8080
relays:
  1:
    on_off_pin: 12
  2:
    on_off_pin: 16
  3:
    on_off_pin: 18
  4:
    on_off_pin: 22
```

API routes:
```
GET /relay/{id}/on
GET /relay/{id}/off
GET /relay/{id}/status
```

It can be used with the [Homebridge](https://github.com/nfarina/homebridge) plugin [homebridge-http-switch](https://github.com/Supereg/homebridge-http-switch).

Sample Homebridge config:

```json
{
    "accessories": [
        {
            "accessory": "HTTP-SWITCH",
            "name": "Turnout Switch 1",
            "switchType": "stateful",
            "onUrl": "http://rpi:8080/relay/1/on",
            "offUrl": "http://rpi:8080/relay/1/off",
            "statusUrl": "http://rpi:8080/relay/1/status"
        }
    ]
}
```
