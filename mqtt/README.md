# mqtt

demostrate mqtt pub / sub model with golang client

## deploy env

### connect to eclipse demo mqtt broker

`tcp://iot.eclipse.org:1883`

### run from docker image

```
docker pull toke/mosquitto
docker run -ti -p 1883:1883 -p 9001:9001 toke/mosquitto
```