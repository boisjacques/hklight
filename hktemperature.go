package main

import (
	"encoding/binary"
	"flag"
	"github.com/boisjacques/hc"
	"github.com/boisjacques/hc/accessory"
	"github.com/boisjacques/hc/mqtt"
	"log"
	"math"
	"strconv"
)

func float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func refresh(acc accessory.Thermometer, c chan []byte) {
	payload := make([]byte, 0)
	for {
		payload = <-c
		stringified := string(payload)
		if temp, err := strconv.ParseFloat(stringified, 32); err == nil {
			acc.TempSensor.CurrentTemperature.SetValue(temp)
			log.Printf("Temperature set to %f\n", temp)
		}

	}

}

func main() {

	username := flag.String("username", "", "mqtt user")
	password := flag.String("password", "", "mqtt password")
	topic := flag.String("topic", "", "mqtt topic")
	clientId := flag.String("id", "", "mqtt client id")
	accessoryName := flag.String("acc-name", "", "accessory name")
	flag.Parse()
	c := make(chan []byte)
	mqtt.NewMQTTBridge(*username, *password, *topic, *clientId, c)
	info := accessory.Info{
		Name:         *accessoryName,
		Manufacturer: "HoChiMinh Flowerpower Enterprises",
	}
	acc := accessory.NewTemperatureSensor(info, 20, -50, 50, .1)

	go refresh(*acc, c)

	t, err := hc.NewIPTransport(hc.Config{Pin: "11223344"}, acc.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
