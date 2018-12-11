package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"trippy/slotmachine"
	"trippy/slotmachine/engine/atkins"
)

type Service interface {
	Initialize() error
	Start() error
	Stop() error
}

type Server struct {
	Id        string
	webserver *http.Server
}

var (
	slog              *log.Logger             // Stdout Logger
	apiKey            string                  // API key used for encrypting teh JWT
	atkinsDietMachine slotmachine.SlotMachine // Slot machine engine
)

const (
	// Env variable for API Key used to encrypt the JWT token
	_API_KEY_PATH = "TRIPPY_API_KEY_PATH"
)

func (s *Server) Initialize() error {

	// Initialize server
	s.Id, _ = os.Hostname()

	// Log to stdout
	slog = log.New(os.Stdout, "", 0)
	slog.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)

	// Checking the required environment variables
	apiKeyFile := os.Getenv(_API_KEY_PATH)
	if apiKeyFile == "" {
		slog.Printf("Webserver API key file [Env:%s] not initialized", _API_KEY_PATH)
		return fmt.Errorf("Webserver API key file not initialized")
	} else {
		if key, err := ioutil.ReadFile(apiKeyFile); err != nil {
			slog.Printf("WARN: Unable to read Authentication Key from [File:%s] [E:%s]", apiKeyFile, err)
			return fmt.Errorf("Unable to read Webserver API key from [File:%s] [E:%s]", apiKeyFile, err)
		} else {
			apiKey = strings.TrimSpace(string(key))
		}
	}

	// Initializing slot machines
	atkinsDietMachine = atkins.NewAtkinsDietMachine()

	return nil
}

func (s *Server) Start() (err error) {
	slog.Printf("Trippy: [%s] starting up...", s.Id)

	var (
		wsServerStopped = make(chan struct{}, 1)
	)
	// handler for signals. Capture ctrl-C and KILL signals
	sigchan := make(chan os.Signal, 1)
	// Reset undos the effect of any prior calls to Notify for the provided signals. If no signals are provided, all signal handlers will be reset.
	signal.Reset()
	signal.Notify(sigchan,
		os.Interrupt,
		os.Kill,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer func() {
		if r := recover(); r != nil {
			slog.Println("Trippy >> Recovered from panic: ", r)
			slog.Printf("%s \n", debug.Stack())
			debug.PrintStack()
			//Ungraceful shutdown of webserver
			s.CloseWebServer()
		}
		slog.Printf("Trippy: [%s] shutting down...", s.Id)
	}()

	s.StartWebServer(wsServerStopped)

LOOP:

	for {
		select {
		case <-wsServerStopped:
			slog.Printf("Trippy:[%s] WebServer stopped. Restarting...", s.Id)
			s.StartWebServer(wsServerStopped)

		case sig := <-sigchan:
			slog.Printf("Trippy received [Signal:%s]. Stopping [Server:%s]...", sig, s.Id)
			s.Stop()
			break LOOP
		}
	}
	return nil
}

func (s *Server) Stop() (err error) {

	// Graceful Shutdown of webserver
	slog.Printf("Stopping Webserver...")
	s.StopWebServer()

	return nil
}
