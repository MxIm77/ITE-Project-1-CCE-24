package source

import (
	PhoeniciaDigitalUtils "Phoenicia-Digital-Base-API/base/utils"
	PhoeniciaDigitalConfig "Phoenicia-Digital-Base-API/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cgxeiji/servo"
)

type servoMotor struct {
	Motor        *servo.Servo
	loitering    bool
	loiterSpeed  float32
	currentPos   float64
	rotateDegree int
	ctx          context.Context
	cancel       context.CancelFunc
}

type servoResponse struct {
	Message string `json:"message"`
	Degree  int    `json:"degree"`
}

var ServoMotor *servoMotor = &servoMotor{}

func HandleLoiter(w http.ResponseWriter, r *http.Request) PhoeniciaDigitalUtils.PhoeniciaDigitalResponse {
	if err := ServoMotor.Loiter(); err != nil {
		return PhoeniciaDigitalUtils.ApiError{Code: http.StatusInternalServerError, Quote: "Failed to Toggle Loiter"}
	}

	return PhoeniciaDigitalUtils.ApiSuccess{Code: http.StatusOK, Quote: servoResponse{Message: "Loiter Toggled", Degree: int(ServoMotor.currentPos)}}
}

func HandleRotateRight(w http.ResponseWriter, r *http.Request) PhoeniciaDigitalUtils.PhoeniciaDigitalResponse {
	if err := ServoMotor.RotateRight(); err != nil {
		return PhoeniciaDigitalUtils.ApiError{Code: http.StatusConflict, Quote: err.Error()}
	}

	return PhoeniciaDigitalUtils.ApiSuccess{Code: http.StatusOK, Quote: servoResponse{Message: fmt.Sprintf("Rotated %d Degrees to the Right", ServoMotor.rotateDegree), Degree: int(ServoMotor.currentPos)}}
}

func HandleRotateLeft(w http.ResponseWriter, r *http.Request) PhoeniciaDigitalUtils.PhoeniciaDigitalResponse {
	if err := ServoMotor.RotateLeft(); err != nil {
		return PhoeniciaDigitalUtils.ApiError{Code: http.StatusConflict, Quote: err.Error()}
	}

	return PhoeniciaDigitalUtils.ApiSuccess{Code: http.StatusOK, Quote: servoResponse{Message: fmt.Sprintf("Rotated %d Degrees to the Left", ServoMotor.rotateDegree), Degree: int(ServoMotor.currentPos)}}
}

func (s *servoMotor) InitializeServoMotor() {

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

	loitspeed, err := strconv.ParseFloat(PhoeniciaDigitalConfig.Config.Pins.LoiterSpeed, 32)
	if err != nil {
		log.Fatalf("Failed to COnvert Loiter Speed to float32")
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.Motor = servo.New(motorPin)
	s.loitering = false
	s.loiterSpeed = float32(loitspeed)
	s.currentPos = 90.0
	s.rotateDegree = rotationdeg

	if err := s.Motor.Connect(); err != nil {
		log.Fatalf("Failed to connect to Servo Motor | Error: %s", err.Error())
	}

	s.Motor.MoveTo(s.currentPos).Wait()

}

func (s *servoMotor) Loiter() error {
	if !s.loitering {
		s.loitering = true
		s.ctx, s.cancel = context.WithCancel(context.Background())

		go func() {
			// defer s.wg.Done()
			for {
				select {
				case <-s.ctx.Done():
					return
				default:
					s.Motor.SetSpeed(0.15)
					s.Motor.MoveTo(180).Wait()
					s.Motor.MoveTo(0).Wait()
				}
			}
		}()
	} else {
		s.Motor.SetSpeed(0)

		s.cancel()
		s.Motor.SetPosition(s.currentPos)

		s.loitering = false
	}

	return nil
}

func (s *servoMotor) RotateRight() error {
	if s.loitering {
		return fmt.Errorf("cannot rotate while loitering")
	}

	if s.currentPos < 180 {
		s.currentPos += float64(s.rotateDegree)
		s.Motor.SetSpeed(0.15)
		s.Motor.MoveTo(s.currentPos).Wait()

	} else {
		return fmt.Errorf("cannot rotate max angle reached %f Degrees", s.currentPos)
	}

	return nil
}

func (s *servoMotor) RotateLeft() error {
	if s.loitering {
		return fmt.Errorf("cannot rotate while loitering")
	}

	if s.currentPos > 0 {
		s.currentPos -= float64(s.rotateDegree)
		s.Motor.SetSpeed(0.15)
		s.Motor.MoveTo(s.currentPos).Wait()
	} else {
		return fmt.Errorf("cannot rotate max angle reached %f Degrees", s.currentPos)
	}

	return nil
}

func init() {
	ServoMotor.InitializeServoMotor()
}
