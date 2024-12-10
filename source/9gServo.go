package source

import (
	PhoeniciaDigitalUtils "Phoenicia-Digital-Base-API/base/utils"
	PhoeniciaDigitalConfig "Phoenicia-Digital-Base-API/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cgxeiji/servo"
)

type servoMotor struct {
	Motor        *servo.Servo
	loitering    bool
	currentPos   float64
	rotateDegree int
	ctx          context.Context
	cancel       context.CancelFunc
}

var ServoMotor *servoMotor = InitializeServoMotor()

func HandleLoiter(w http.ResponseWriter, r *http.Request) PhoeniciaDigitalUtils.PhoeniciaDigitalResponse {
	if err := ServoMotor.Loiter(); err != nil {
		return PhoeniciaDigitalUtils.ApiError{Code: http.StatusInternalServerError, Quote: "Failed to Toggle Loiter"}
	}

	return PhoeniciaDigitalUtils.ApiSuccess{Code: http.StatusOK, Quote: "Loiter Toggled"}
}

func HandleRotateRight(w http.ResponseWriter, r *http.Request) PhoeniciaDigitalUtils.PhoeniciaDigitalResponse {
	if err := ServoMotor.RotateRight(); err != nil {
		return PhoeniciaDigitalUtils.ApiError{Code: http.StatusConflict, Quote: err.Error()}
	}

	return PhoeniciaDigitalUtils.ApiSuccess{Code: http.StatusOK, Quote: fmt.Sprintf("Rotated %d Degrees to the right", ServoMotor.rotateDegree)}
}

func HandleRotateLeft(w http.ResponseWriter, r *http.Request) PhoeniciaDigitalUtils.PhoeniciaDigitalResponse {
	if err := ServoMotor.RotateLeft(); err != nil {
		return PhoeniciaDigitalUtils.ApiError{Code: http.StatusConflict, Quote: err.Error()}
	}

	return PhoeniciaDigitalUtils.ApiSuccess{Code: http.StatusOK, Quote: fmt.Sprintf("Rotated %d Degrees to the left", ServoMotor.rotateDegree)}
}

func InitializeServoMotor() *servoMotor {

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

	// Set Desired Rotation Degree
	rotationdeg, err := strconv.Atoi(PhoeniciaDigitalConfig.Config.Pins.RotateDegree)
	if err != nil {
		log.Fatalf("Trigger Pin Value in .env file is an invalid pin number")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &servoMotor{
		Motor:        servo.New(motorPin),
		loitering:    false,
		currentPos:   90.0,
		rotateDegree: rotationdeg,
		ctx:          ctx,
		cancel:       cancel,
	}

}

func (s *servoMotor) Loiter() error {
	if !s.loitering {
		s.loitering = true
		go func() {
			defer s.cancel() // Cancel context on goroutine exit
			for {
				select {
				case <-s.ctx.Done():
					return
				default:
					s.Motor.SetPosition(180)

					time.Sleep(500 * time.Millisecond)

					s.Motor.SetPosition(0)

					time.Sleep(500 * time.Millisecond)
				}
			}
		}()
	} else {
		s.cancel()
		s.loitering = false
		s.Motor.SetPosition(s.currentPos)
	}

	return nil
}

func (s *servoMotor) RotateRight() error {
	if s.loitering {
		return fmt.Errorf("cannot rotate while loitering")
	}

	if s.currentPos < 180 {
		s.currentPos += float64(s.rotateDegree)
		s.Motor.SetPosition(s.currentPos)
	} else {
		return fmt.Errorf("cannot rotate max angle reached")
	}

	return nil
}

func (s *servoMotor) RotateLeft() error {
	if s.loitering {
		return fmt.Errorf("cannot rotate while loitering")
	}

	if s.currentPos > 0 {
		s.currentPos -= float64(s.rotateDegree)
		s.Motor.SetPosition(s.currentPos)
	} else {
		return fmt.Errorf("cannot rotate max angle reached")
	}

	return nil
}
