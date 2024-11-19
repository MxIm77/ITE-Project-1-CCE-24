package source

import (
	PhoeniciaDigitalConfig "Phoenicia-Digital-Base-API/config"
	"log"
	"strconv"
	"time"

	"github.com/warthog618/gpio"
)

type _HCSR04 struct {
	Trigger *gpio.Pin
	Echo    *gpio.Pin
}

var gpioOpen bool

var HCSR04 *_HCSR04 = initializeHCSR04()

func initializeHCSR04() *_HCSR04 {

	//	Try OPEN & Memory map the GPIO pins temporarely to ba able to assign the pins
	if !gpioOpen {
		if err := gpio.Open(); err != nil {
			log.Fatalf("Failed to Open & Memory map GPIO pins | Error: %s", err.Error())
		}
		defer gpio.Close()
		gpioOpen = true
	}

	// Check Pin Conversion
	trigPin, err := strconv.Atoi(PhoeniciaDigitalConfig.Config.Pins.TriggerPin)
	if err != nil {
		log.Fatalf("Trigger Pin Value in .env file is an invalid pin number")
	}

	echoPin, err := strconv.Atoi(PhoeniciaDigitalConfig.Config.Pins.EchoPin)
	if err != nil {
		log.Fatalf("Echo Pin Value in .env file is an invalid pin number")
	}

	// Return a pointer to the initialized HCSR04 struct
	return &_HCSR04{
		Trigger: gpio.NewPin(trigPin),
		Echo:    gpio.NewPin(echoPin),
	}
}

func (hcsr04 *_HCSR04) MeasureDistance() (float64, error) {

	//	Asign Measuring Time Variables
	var pulseStart time.Time
	var pulseEnd time.Time
	var pulseDuration time.Duration

	//	Trigger the sensor
	hcsr04.Trigger.High()
	time.Sleep(10 * time.Millisecond)
	hcsr04.Trigger.Low()

	// Measure Pulse Duration
	pulseStart = time.Now()
	for hcsr04.Echo.Read() != gpio.Low {
		// wait for the Echo Pin to go Low
	}
	pulseEnd = time.Now()

	pulseDuration = pulseEnd.Sub(pulseStart)

	return float64((pulseDuration.Nanoseconds() / 2 / 5800) / 10000), nil

}
