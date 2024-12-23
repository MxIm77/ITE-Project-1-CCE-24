package source

import (
	PhoeniciaDigitalUtils "Phoenicia-Digital-Base-API/base/utils"
	PhoeniciaDigitalConfig "Phoenicia-Digital-Base-API/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stianeikeland/go-rpio/v4"
)

type hcsr04 struct {
	Trigger     rpio.Pin
	Echo        rpio.Pin
	SpeedOfWave float32
	pulseWidth  time.Duration
}

// SensorData represents the data structure for the sensor's output
type SensorData struct {
	Distance float64 `json:"distance"`
	Status   string  `json:"status"`
}

// WebSocket upgrader for handling HTTP requests to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var HCSR04 *hcsr04 = &hcsr04{}

func (h *hcsr04) InitializeUltrasonicSensor() {

	// Initialize GPIO
	err := rpio.Open()
	if err != nil {
		PhoeniciaDigitalUtils.Log(fmt.Sprintf("Failed to open GPIO: %v", err))
		log.Fatalf("Failed to open GPIO: %v", err)
	}
	// defer rpio.Close()

	// Check Pin Conversion from the .env file (should be actual numbers and in range of the raspberry pi zero w pins)
	// If an issue occured with conversion the program wont run!
	trigPin, err := strconv.Atoi(PhoeniciaDigitalConfig.Config.Pins.TriggerPin)
	if err != nil {
		PhoeniciaDigitalUtils.Log("Trigger Pin Value in .env file is an invalid pin number")
		log.Fatalf("Trigger Pin Value in .env file is an invalid pin number")
	}

	// Check Pin Conversion from the .env file (should be actual numbers and in range of the raspberry pi zero w pins)
	// If an issue occured with conversion the program wont run!
	echoPin, err := strconv.Atoi(PhoeniciaDigitalConfig.Config.Pins.EchoPin)
	if err != nil {
		PhoeniciaDigitalUtils.Log("Echo Pin Value in .env file is an invalid pin number")
		log.Fatalf("Echo Pin Value in .env file is an invalid pin number")
	}

	// Make sure the Trigger & Echo pins are not the same & make sure the Trigger & Echo pins are in the GPIO range map of a raspberry pi zero w v1
	if trigPin == echoPin {
		PhoeniciaDigitalUtils.Log(fmt.Sprintf("Echo pin: %d, & Trigger pin: %d, CANT be assigned to the same pin", echoPin, trigPin))
		log.Fatalf("Echo pin: %d, & Trigger pin: %d, CANT be assigned to the same pin", echoPin, trigPin)
	} else if trigPin < 2 || trigPin > 27 {
		PhoeniciaDigitalUtils.Log(fmt.Sprintf("Trigger pin: %d, out of GPIO map range [2 -> 27] | Please Change it in the ~/config/.env file", trigPin))
		log.Fatalf("Trigger pin: %d, out of GPIO map range [2 -> 27] | Please Change it in the ~/config/.env file", trigPin)
	} else if echoPin < 2 || echoPin > 27 {
		PhoeniciaDigitalUtils.Log(fmt.Sprintf("Echo pin: %d, out of GPIO map range [2 -> 27] | Please Change it in the ~/config/.env file", echoPin))
		log.Fatalf("Echo pin: %d, out of GPIO map range [2 -> 27] | Please Change it in the ~/config/.env file", echoPin)
	}

	// Set the proper Trigger pin map to the struct HCSR04 & Make the Trigger pin an output Pin
	h.Trigger = rpio.Pin(trigPin)
	h.Trigger.Mode(rpio.Output)

	// Set the proper Echo pin map to the struct HCSR04 & Make the Echo pin an output Pin
	h.Echo = rpio.Pin(echoPin)
	h.Echo.Mode(rpio.Input)

	// Assign other variables that will be linked to the hc-sr04
	h.SpeedOfWave = 0.0343
	h.pulseWidth = 10 * time.Microsecond
	PhoeniciaDigitalUtils.Log(fmt.Sprintf("Initialized With Trigger Pin: %d, Echo Pin: %d", trigPin, echoPin))
	log.Printf("Initialized With Trigger Pin: %d, Echo Pin: %d", trigPin, echoPin)

}

// Function to measure distance in centimeters
func (h *hcsr04) MeasureDistance() float64 {
	// Send a pulse to the trigger pin
	h.Trigger.Low()
	time.Sleep(h.pulseWidth) // Delay to ensure pulse width is valid
	h.Trigger.High()
	time.Sleep(h.pulseWidth) // Trigger pulse duration
	h.Trigger.Low()

	// Wait for the echo pulse to start
	for h.Echo.Read() == rpio.Low {
	}

	// Record the start time
	start := time.Now()

	// Wait for the echo pulse to end
	for h.Echo.Read() == rpio.High {
	}

	duration := time.Since(start)

	// Calculate distance in cm
	distance := (float64(duration) * float64(h.SpeedOfWave)) / 2 // Convert to cm
	return distance
}

// WebSocket handler for handling connections and sending data to clients
func HandleMeasureDistance(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
		// return PhoeniciaDigitalUtils.ApiError{Code: http.StatusInternalServerError, Quote: err.Error()}
	}
	defer conn.Close()

	log.Println("New WebSocket client connected")

	// Start measuring and sending data to the WebSocket client in a goroutine
	for {
		// Measure the distance
		distance := HCSR04.MeasureDistance()

		// Prepare the response struct
		sensorData := SensorData{
			Distance: distance,
		}

		// If there's an error in the measurement, set the status as an error
		if distance == -1 {
			sensorData.Status = "Error measuring distance"
		} else {
			sensorData.Status = "Success"
		}

		// Marshal the struct to JSON
		jsonData, err := json.Marshal(sensorData)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
			break
		}

		// Send the JSON data to the WebSocket client
		err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("Error sending data:", err)
			break
		}

		// Sleep before taking the next measurement
		time.Sleep(1 * time.Second) // Measure every 1 second
	}

	// return PhoeniciaDigitalUtils.ApiSuccess{Code: http.StatusOK, Quote: "Websocket Closed"}
}

func init() {
	HCSR04.InitializeUltrasonicSensor()
}
