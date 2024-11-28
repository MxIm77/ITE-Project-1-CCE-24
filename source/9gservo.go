package source

import (
	PhoeniciaDigitalConfig "Phoenicia-Digital-Base-API/config"
	"log"
	"strconv"

	"github.com/warthog618/gpio"
)

// This struct will be used to create our 9g Servo and to be reused through out the entire project if needed it will store important data such as
// MotorPin --> The Digital Pin of the 9g Servo that is responsible for moving the servo via signal
type _Servo struct {
	MotorPin *gpio.Pin
}

// The global struct for our 9g Servo (A pointer will make sure its not a copy reducing ram usage)
var ServoMotor *_Servo = InitializeServo()

func InitializeServo() *_Servo {

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
	motorPin, err := strconv.Atoi(PhoeniciaDigitalConfig.Config.Pins.MotorPin)
	if err != nil {
		log.Fatalf("Trigger Pin Value in .env file is an invalid pin number")
	}

	// Make sure the Motor pin is in the GPIO range map of a raspberry pi zero w v1
	if motorPin < 2 || motorPin > 27 {
		log.Fatalf("Motor pin: %d, out of GPIO map range [2 -> 27] | Please Change it in the ~/config/.env file", motorPin)
	}

	// Return a pointer to the initialized Servo
	return &_Servo{
		// We set a Motor pin to be used accross the code base directly from the struct fo reusability and optimization
		// this Motor pins can be chanegd in the ~/config/.env file at the very bottom lines
		MotorPin: gpio.NewPin(motorPin),
	}
}
