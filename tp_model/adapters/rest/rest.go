package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Razzle131/line316/tp_model/core"
)

func NewPingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
}

func NewStartPuck(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.PlaceNewStartPuck()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
	}
}

const sensorPathName = "sensor_id"

type GetSenorResponse struct {
	Value bool `json:"value"`
}

func NewSensorHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sensorId := r.PathValue(sensorPathName)
		if sensorId == "" {
			http.Error(w, "missing sensor id field", http.StatusBadRequest)
			return
		}

		val, err := s.GetSensorValue(sensorId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(GetSenorResponse{val}); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

func NewGripperLeftHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.MoveGripperLeft()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewGripperRightHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.MoveGripperRight()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewGripperUpHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.MoveGripperUp()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewGripperDownHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.MoveGripperDown()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewGripperOpenHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.OpenGripper()
		if err != nil {
			log.Error("open gripper", "error", err)
			http.Error(w, fmt.Sprintf("suspicious opening of gripper: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewGripperCloseHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.CloseGripper()
		if err != nil {
			log.Error("close gripper", "error", err)
			http.Error(w, fmt.Sprintf("suspicious closing of gripper: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewGripperStopHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.StopGripper()
		time.Sleep(20 * time.Millisecond)
		s.EnableMovingGripper()
		w.WriteHeader(http.StatusOK)
	}
}

func NewCarouselRotateHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.RotateCarousel()

		w.WriteHeader(http.StatusOK)
	}
}

type InspectResponse struct {
	Color string `json:"color"`
}

func NewCarouselInspectHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		puck, err := s.InspectPuck()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(InspectResponse{puck.Color}); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

func NewCarouselDrillHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.DrillPuck()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewPackagingHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.PackagePuck()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewSortingHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.SortPuck()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// data for visualisation
type Start struct {
	PuckSlot *core.Puck `json:"puckSlot"`
}
type Gripper struct {
	IsOpen                bool       `json:"isOpen"`
	PuckSlot              *core.Puck `json:"puckSlot"`
	CurHorizontalPosition float64    `json:"curHorizontalPosition"`
	CurVerticalPosition   float64    `json:"curVerticalPosition"`
}
type Carousel struct {
	Slots []*core.Puck `json:"slots"`
}
type PackagingLine struct {
	PuckSlot *core.Puck `json:"puckSlot"`
}
type SortingLine struct {
	PuckSlot *core.Puck             `json:"puckSlot"`
	Produced map[string][]core.Puck `json:"produced"`
}

func NewStartHandler(s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := s.Start

		resp := Start{
			PuckSlot: model.PuckSlot,
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func NewGripperHandler(s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := &s.Gripper

		resp := Gripper{
			IsOpen:                model.IsOpen,
			PuckSlot:              model.PuckSlot,
			CurHorizontalPosition: model.CurHorizontalPosition,
			CurVerticalPosition:   model.CurVerticalPosition,
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func NewCarouselHandler(s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := s.Carousel

		resp := Carousel{
			Slots: model.Slots,
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func NewPackagingLineHandler(s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := s.PackagingLine

		resp := PackagingLine{
			PuckSlot: model.PuckSlot,
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func NewSortingLineHandler(s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := s.SortingLine

		resp := SortingLine{
			PuckSlot: model.PuckSlot,
			Produced: model.Produced,
		}

		json.NewEncoder(w).Encode(resp)
	}
}
