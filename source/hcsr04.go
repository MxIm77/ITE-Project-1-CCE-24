package source

import (
	PhoeniciaDigitalConfig "Phoenicia-Digital-Base-API/config"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/warthog618/gpio"
)

// This struct will be used to create our HC-SR04 Ultrasonic sensor and to be reused through out the entire project if needed it will store important data such as
// Echo --> Echo pin
// Trigger --> Trigger pin
// MaxDuration --> the max ammount of time the ulrasonic sensor could use up at its theoretical maximal operation range of 400cm (WILL BE USED AS TIMEOUT FOR ECHO PIN)
type _HCSR04 struct {
	Trigger     *gpio.Pin
	Echo        *gpio.Pin
	MaxDuration time.Duration
}

// Prevents any other memory mappings
var gpioOpen bool = false

// The global struct for our HC-SR04 Ultrasonic Sensor (A pointer will make sure its not a copy reducing ram usage)
var HCSR04 *_HCSR04 = initializeHCSR04()

// Function to initialize our HCSR04 variable with the proper data needed
func initializeHCSR04() *_HCSR04 {

	// Open & Memory map the GPIO pins temporarely to ba able to assign the pins
	if !gpioOpen {
		if err := gpio.Open(); err != nil {
			log.Fatalf("Failed to Open & Memory map GPIO pins | Error: %s", err.Error())
		}
		defer gpio.Close()
		gpioOpen = true

		defer func() {
			gpioOpen = false
		}()
	}

	// Check Pin Conversion from the .env file (should be actual numbers and in range of the raspberry pi zero w pins)
	// If an issue occured with conversion the program wont run!
	trigPin, err := strconv.Atoi(PhoeniciaDigitalConfig.Config.Pins.TriggerPin)
	if err != nil {
		log.Fatalf("Trigger Pin Value in .env file is an invalid pin number")
	}

	// Check Pin Conversion from the .env file (should be actual numbers and in range of the raspberry pi zero w pins)
	// If an issue occured with conversion the program wont run!
	echoPin, err := strconv.Atoi(PhoeniciaDigitalConfig.Config.Pins.EchoPin)
	if err != nil {
		log.Fatalf("Echo Pin Value in .env file is an invalid pin number")
	}

	// Make sure the Trigger & Echo pins are not the same & make sure the Trigger & Echo pins are in the GPIO range map of a raspberry pi zero w v1
	if trigPin == echoPin {
		log.Fatalf("Echo pin: %d, & Trigger pin: %d, CANT be assigned to the same pin", echoPin, trigPin)
	} else if trigPin < 2 || trigPin > 27 {
		log.Fatalf("Trigger pin: %d, out of GPIO map range [2 -> 27] | Please Change it in the ~/config/.env file", trigPin)
	} else if echoPin < 2 || echoPin > 27 {
		log.Fatalf("Echo pin: %d, out of GPIO map range [2 -> 27] | Please Change it in the ~/config/.env file", echoPin)
	}

	// Return a pointer to the initialized HCSR04 struct
	return &_HCSR04{
		// We set a Trigger pin & echo pin to be used accross the code base directly from the struct fo reusability and optimization
		// these trig & echi pins can be chanegd in the ~/config/.env file at the very bottom lines
		Trigger: gpio.NewPin(trigPin),
		Echo:    gpio.NewPin(echoPin),

		// For this Specific Ultrasonic Sensor (2cm - 400cm) is the theoretical operation range; so we are assigning a MaxDuration
		// which acts as a timeout for the echo pin to recieve a signal otherwise return 0 since it cant operate at distances greater than 400cm
		// the calculation is straight forward "2 * 400" this is since the wave will travel bounce back so essentialy it can travel 800cm at the highest theoretical range
		// then we divide this distance of 800cm by 0.0343 which is the speed of the sound wave for cm/ns | In addition we add 1 second just in case & convert the value to nanoseconds.
		MaxDuration: 46647 + time.Second,
	}
}

func (hcsr04 *_HCSR04) MeasureDistance() (float64, error) {
	// Assign Measuring Time Variables
	var pulseStart time.Time
	var pulseEnd time.Time

	// Send a 5v signal to the sensor for 10 Milliseconds then turn it off (Minimum ammount needed for the HC-SR04 Ultrasonic Sensor to function properly)
	hcsr04.Trigger.High()
	time.Sleep(10 * time.Millisecond)
	hcsr04.Trigger.Low()

	// Measure Pulse Duration
	pulseStart = time.Now()
	for time.Since(pulseStart) < hcsr04.MaxDuration {
		if hcsr04.Echo.Read() == gpio.Low {
			pulseEnd = time.Now()
			break
		}
	}

	// Check if timeout occurred
	if pulseEnd.IsZero() {
		return 0, fmt.Errorf("echo pin timed out it's Max Duration: %s Nanoseconds| Theoretical operation range 2cm - 400cm | Effective operation range 2cm - 80cm", hcsr04.MaxDuration)
	}

	return float64(pulseEnd.Sub(pulseStart).Nanoseconds()) / 2 * 0.0343, nil
}
