package core

import (
	"errors"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/Razzle131/line316/tp_model/adapters/sensor"
)

type Service struct {
	logger *slog.Logger

	// actuators map[string]Actuator // addr -> obj
	sensors map[string]Sensor // addr -> obj

	gripper       Gripper
	start         Start
	carousel      Carousel
	packagingLine PackagingLine
	sortingLine   SortingLine
}

func NewService(logger *slog.Logger) *Service {
	s := Service{
		logger:        logger,
		sensors:       make(map[string]Sensor),
		start:         NewStart(),
		carousel:      NewCarousel(),
		gripper:       NewGripper(),
		packagingLine: NewPackagingLine(),
		sortingLine:   NewSortingLine(),
	}

	// // Processing station PLC sensors
	// s.sensors["ns:4, i:5"] = sensor.New("processing_input_4_workpiece_detected", "ns:4, i:5")
	// s.sensors["ns:4, i:7"] = sensor.New("processing_input_2_workpiece_silver", "ns:4, i:7")
	// s.sensors["ns:4, i:3"] = sensor.New("processing_input_5_carousel_init", "ns:4, i:3")
	// s.sensors["ns:4, i:4"] = sensor.New("processing_input_6_hole_detected", "ns:4, i:4")
	// s.sensors["ns:4, i:6"] = sensor.New("processing_input_7_workpiece_not_black", "ns:4, i:6")

	// // Handling and Packing PLC sensors
	// s.sensors["ns:4, i:29"] = sensor.New("handling_input_0_workpiece_pushed", "ns:4, i:29")
	// s.sensors["ns:4, i:32"] = sensor.New("handling_input_1_grippe_at_right", "ns:4, i:32")
	// s.sensors["ns:4, i:31"] = sensor.New("handling_input_2_gripper_at_start", "ns:4, i:31")
	// s.sensors["ns:4, i:33"] = sensor.New("handling_input_3_gripper_down_pack_lvl", "ns:4, i:33")
	// s.sensors["ns:4, i:42"] = sensor.New("packing_input_7_pack_turned_on", "ns:4, i:42")

	// // Sorting station PLC sensors
	// s.sensors["ns:4, i:9"] = sensor.New("sorting_input_3_box_on_conveyor", "ns:4, i:9")
	// s.sensors["ns:4, i:10"] = sensor.New("sorting_input_4_box_is_down", "ns:4, i:10")

	// // Processing station PLC actuators
	// s.actuators["ns:4, i:12"] = actuator.New("processing_output_0_drill", "ns:4, i:12")
	// s.actuators["ns:4, i:13"] = actuator.New("processing_output_1_rotate_carousel", "ns:4, i:13")
	// s.actuators["ns:4, i:14"] = actuator.New("processing_output_2_drill_down", "ns:4, i:14")
	// s.actuators["ns:4, i:15"] = actuator.New("processing_output_3_drill_up", "ns:4, i:15")
	// s.actuators["ns:4, i:16"] = actuator.New("processing_output_4_fix_workpiece", "ns:4, i:16")
	// s.actuators["ns:4, i:17"] = actuator.New("processing_output_5_detect_hole", "ns:4, i:17")

	// // Handling and Packing PLC actuators
	// s.actuators["ns:4, i:34"] = actuator.New("handling_output_0_to_green", "ns:4, i:34")
	// s.actuators["ns:4, i:35"] = actuator.New("handling_output_1_to_yellow", "ns:4, i:35")
	// s.actuators["ns:4, i:36"] = actuator.New("handling_output_2_to_red", "ns:4, i:36")
	// s.actuators["ns:4, i:37"] = actuator.New("handling_output_3_gripper_to_right", "ns:4, i:37")
	// s.actuators["ns:4, i:38"] = actuator.New("handling_output_4_gripper_to_left", "ns:4, i:38")
	// s.actuators["ns:4, i:39"] = actuator.New("handling_output_5_gripper_to_down", "ns:4, i:39")
	// s.actuators["ns:4, i:40"] = actuator.New("handling_output_6_gripper_to_open", "ns:4, i:40")
	// s.actuators["ns:4, i:41"] = actuator.New("handling_output_7_gripper_push_workpiece", "ns:4, i:41")
	// s.actuators["ns:4, i:43"] = actuator.New("packing_output_4_push_box", "ns:4, i:43")
	// s.actuators["ns:4, i:44"] = actuator.New("packing_output_5_fix_box_upper_side", "ns:4, i:44")
	// s.actuators["ns:4, i:45"] = actuator.New("packing_output_6_fix_box_tongue", "ns:4, i:45")
	// s.actuators["ns:4, i:46"] = actuator.New("packing_output_7_pack_box", "ns:4, i:46")

	// // Sorting station PLC actuators
	// s.actuators["ns:4, i:19"] = actuator.New("sorting_output_0_move_conveyor_right", "ns:4, i:19")
	// s.actuators["ns:4, i:20"] = actuator.New("sorting_output_1_move_conveyor_left", "ns:4, i:20")
	// s.actuators["ns:4, i:21"] = actuator.New("sorting_output_2_push_silver_workpiece", "ns:4, i:21")
	// s.actuators["ns:4, i:22"] = actuator.New("sorting_output_3_push_red_workpiece", "ns:4, i:22")

	s.sensors["ns:1, i:1"] = sensor.New("gripper carousel position", "ns:1, i:1")
	s.sensors["ns:1, i:2"] = sensor.New("gripper start position", "ns:1, i:2")
	s.sensors["ns:1, i:3"] = sensor.New("gripper packaging position", "ns:1, i:3")
	s.sensors["ns:1, i:4"] = sensor.New("gripper sorting position", "ns:1, i:4")

	go s.updateSensors()
	go s.printGripperPos()

	return &s
}

func (s *Service) printGripperPos() {
	ticker := time.NewTicker(time.Millisecond * 100)
	for range ticker.C {
		s.logger.Debug("gripper pos", "x", s.gripper.CurHorizontalPosition, "y", s.gripper.CurVerticalPosition)
	}
}

func (s *Service) updateSensors() {
	ticker := time.NewTicker(time.Second / tickrate)
	for range ticker.C {
		curGripperPos := s.gripper.CurHorizontalPosition
		if math.Abs(gripperCarouselPos-curGripperPos) <= gripperAbleMiss {
			s.sensors["ns:1, i:1"].WriteValue(true)
		} else {
			s.sensors["ns:1, i:1"].WriteValue(false)
		}

		if math.Abs(gripperStartPos-curGripperPos) <= gripperAbleMiss {
			s.sensors["ns:1, i:2"].WriteValue(true)
		} else {
			s.sensors["ns:1, i:2"].WriteValue(false)
		}

		if math.Abs(gripperPackagingPos-curGripperPos) <= gripperAbleMiss {
			s.sensors["ns:1, i:3"].WriteValue(true)
		} else {
			s.sensors["ns:1, i:3"].WriteValue(false)
		}

		if math.Abs(gripperSortingPos-curGripperPos) <= gripperAbleMiss {
			s.sensors["ns:1, i:4"].WriteValue(true)
		} else {
			s.sensors["ns:1, i:4"].WriteValue(false)
		}
	}
}

func (s *Service) GetSensorValue(sensorId string) (bool, error) {
	sensor, found := s.sensors[sensorId]
	if !found {
		return false, errors.New("sensor not found")
	}

	return sensor.GetValue(), nil
}

func (s *Service) PlaceNewStartPuck() error {
	colors := []string{"red", "silver", "black"}

	puck := NewPuck(colors[rand.Intn(len(colors))])
	err := s.start.PlacePuck(puck)

	return err
}

func (s *Service) MoveGripperLeft() error {
	err := s.gripper.MoveLeft()
	if err != nil {
		s.logger.Error("move left", "err", err)
	}
	return err
}

func (s *Service) MoveGripperRight() error {
	err := s.gripper.MoveRight()
	if err != nil {
		s.logger.Error("move right", "err", err)
	}
	return err
}

func (s *Service) MoveGripperUp() error {
	err := s.gripper.MoveUp()
	if err != nil {
		s.logger.Error("move up", "err", err)
	}
	return err
}

func (s *Service) MoveGripperDown() error {
	err := s.gripper.MoveDown()
	if err != nil {
		s.logger.Error("move down", "err", err)
	}
	return err
}

func (s *Service) OpenGripper() error {
	s.gripper.Open()

	var err error
	if s.gripper.PuckSlot != nil {
		if s.gripper.CurVerticalPosition > gripperDownPos {
			err = errors.New("opening gripper with puck in higher than lower position")
		} else {
			err = s.placePuck()
			if err != nil {
				s.logger.Error("place puck", "error", err)
			}
		}
	}

	return err
}

func (s *Service) CloseGripper() error {
	var err error
	if s.gripper.PuckSlot == nil && s.gripper.CurVerticalPosition <= gripperDownPos {
		err = s.takePuck()
		if err != nil {
			s.logger.Error("take puck", "error", err)
		}
	}

	s.gripper.Close()

	return err
}

func (s *Service) StopGripper() {
	s.gripper.Stop()
}

func (s *Service) EnableMovingGripper() {
	s.gripper.EnableMoving()
}

func (s *Service) takePuck() error {
	curGripperPos := s.gripper.CurHorizontalPosition
	var pucker Pucker
	if math.Abs(gripperCarouselPos-curGripperPos) <= gripperAbleMiss {
		pucker = &s.carousel
	} else if math.Abs(gripperStartPos-curGripperPos) <= gripperAbleMiss {
		pucker = &s.start
	} else if math.Abs(gripperPackagingPos-curGripperPos) <= gripperAbleMiss {
		pucker = &s.packagingLine
	} else if math.Abs(gripperSortingPos-curGripperPos) <= gripperAbleMiss {
		return errors.New("should not take puck from sorting line")
	} else {
		return errors.New("no position matched for gripper")
	}

	puck, err := pucker.TakePuck()
	if err != nil {
		s.logger.Error("take puck", "error", err)
		return err
	}

	err = s.gripper.TakePuck(puck)
	if err != nil {
		pucker.PlacePuck(puck)
		s.logger.Error("take puck", "error", err)
		return err
	}

	return nil
}

func (s *Service) placePuck() error {
	puck, err := s.gripper.PlacePuck()
	if err != nil {
		s.logger.Error("place puck", "error", err)
		return err
	}

	curGripperPos := s.gripper.CurHorizontalPosition
	var pucker Pucker
	if math.Abs(gripperCarouselPos-curGripperPos) <= gripperAbleMiss {
		pucker = &s.carousel
	} else if math.Abs(gripperStartPos-curGripperPos) <= gripperAbleMiss {
		pucker = &s.start
	} else if math.Abs(gripperPackagingPos-curGripperPos) <= gripperAbleMiss {
		pucker = &s.packagingLine
	} else if math.Abs(gripperSortingPos-curGripperPos) <= gripperAbleMiss {
		pucker = &s.sortingLine
	} else {
		return errors.New("no position matched for gripper")
	}

	err = pucker.PlacePuck(puck)
	if err != nil {
		s.gripper.TakePuck(puck)
		s.logger.Error("place puck", "error", err)
		return err
	}

	return nil
}

func (s *Service) RotateCarousel() {
	s.carousel.RotateOnce()
}

func (s *Service) InspectPuck() (Puck, error) {
	puck, err := s.carousel.InspectPuck()
	if err != nil {
		s.logger.Error("inspect puck", "error", err)
	}
	return puck, err
}

func (s *Service) DrillPuck() error {
	err := s.carousel.DrillPuck()
	if err != nil {
		s.logger.Error("drill puck", "error", err)
	}
	return err
}

func (s *Service) PackagePuck() error {
	err := s.packagingLine.PackagePuck()
	if err != nil {
		s.logger.Error("package puck", "error", err)
	}
	return err
}

func (s *Service) SortPuck() error {
	err := s.sortingLine.SortPuck()
	if err != nil {
		s.logger.Error("package puck", "error", err)
	}
	return err
}
