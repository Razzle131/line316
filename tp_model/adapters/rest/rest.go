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
		w.WriteHeader(http.StatusOK)
	}
}

func NewGripperEnableHandler(log *slog.Logger, s *core.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
