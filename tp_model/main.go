package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/Razzle131/line316/tp_model/adapters/rest"
	"github.com/Razzle131/line316/tp_model/config"
	"github.com/Razzle131/line316/tp_model/core"
)

// данная программа написано криво и гексагональной архитектуре не соответствует, просьба не смотреть), модель учебная, переписывать ее полностью уже поздно

func WithoutCORS(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := mustMakeLogger(cfg.LogLevel)

	if err := run(cfg, log); err != nil {
		log.Error("run func", "error", err)
		os.Exit(1)
	}
}

func run(cfg config.Config, log *slog.Logger) error {
	log.Info("starting server")
	log.Debug("debug messages are enabled")

	service := core.NewService(log)

	mux := http.NewServeMux()

	mux.Handle("GET /tp/ping", rest.NewPingHandler())

	mux.Handle("POST /tp/puck", rest.NewStartPuck(log, service))

	mux.Handle("GET /tp/sensor/{sensor_id}", rest.NewSensorHandler(log, service))

	// gripper
	mux.Handle("POST /tp/gripper/left", rest.NewGripperLeftHandler(log, service))
	mux.Handle("POST /tp/gripper/right", rest.NewGripperRightHandler(log, service))
	mux.Handle("POST /tp/gripper/up", rest.NewGripperUpHandler(log, service))
	mux.Handle("POST /tp/gripper/down", rest.NewGripperDownHandler(log, service))

	mux.Handle("POST /tp/gripper/open", rest.NewGripperOpenHandler(log, service))
	mux.Handle("POST /tp/gripper/close", rest.NewGripperCloseHandler(log, service))

	mux.Handle("POST /tp/gripper/stop", rest.NewGripperStopHandler(log, service))

	// carousel
	mux.Handle("POST /tp/carousel/rotate", rest.NewCarouselRotateHandler(log, service))
	mux.Handle("POST /tp/carousel/inspect", rest.NewCarouselInspectHandler(log, service))
	mux.Handle("POST /tp/carousel/drill", rest.NewCarouselDrillHandler(log, service))

	// packaging
	mux.Handle("POST /tp/packaging/pack", rest.NewPackagingHandler(log, service))

	// sorting
	mux.Handle("POST /tp/sorting/sort", rest.NewSortingHandler(log, service))

	// visualisation
	mux.Handle("GET /vis/start", WithoutCORS(rest.NewStartHandler(service)))
	mux.Handle("GET /vis/gripper", WithoutCORS(rest.NewGripperHandler(service)))
	mux.Handle("GET /vis/carousel", WithoutCORS(rest.NewCarouselHandler(service)))
	mux.Handle("GET /vis/packaging", WithoutCORS(rest.NewPackagingLineHandler(service)))
	mux.Handle("GET /vis/sorting", WithoutCORS(rest.NewSortingLineHandler(service)))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server := http.Server{
		Addr:        cfg.Address,
		ReadTimeout: cfg.Timeout,
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("erroneous shutdown", "error", err)
		}
	}()

	log.Info("Running HTTP server", "address", cfg.Address)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("server closed unexpectedly", "error", err)
			return err
		}
	}

	return nil
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		slog.Error("failed parsing log level", "error", err)
		os.Exit(1)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	return slog.New(handler)
}
