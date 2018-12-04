package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

const (
	_WEBSERVER_PORT      = "7070"
	_WS_SHUTDOWN_TIMEOUT = 5 * time.Second

	_ATKINS_DIET_MACHINE = "atkins-diet"

	_PARA_SPIN_MACHINE = "machine"
)

var (
	middlewareIgnoreList = map[string]struct{}{
		"/": struct{}{},
	}
)

func (s *Server) StartWebServer(stopped chan struct{}) {
	slog.Println("WebServer starting...")

	router := httprouter.New()
	router.GET("/", Home)                            // Root
	router.GET("/hello/:name", Hello)                // Hello test API
	router.GET("/api/machines/:machine/spins", Spin) // Spin the respective slot machine

	neg := negroni.Classic()
	//neg.Use(negroni.HandlerFunc(authMiddleware))
	neg.UseHandler(router)

	httpsrv := &http.Server{
		Addr:    ":" + _WEBSERVER_PORT,
		Handler: neg,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Println("WebServer: Recovered from panic", r)
				debug.PrintStack()
			}

			slog.Println("WebServer exiting...")
			stopped <- struct{}{}
		}()
		if err := httpsrv.ListenAndServe(); err != nil {
			// ErrServerClosed is returned after every call to Shutdown/Close.
			// Shutdown/Close of webserver should not be reported as an error
			if err.Error() != http.ErrServerClosed.Error() {
				slog.Printf("Error: Webserver: Http Server:ListenAndServe() [E:%s]", err)
			}
		}
	}()

	s.webserver = httpsrv
}

// Shutdown is a graceful shutdown of webserver
func (s *Server) StopWebServer() {
	ctx, _ := context.WithTimeout(context.Background(), _WS_SHUTDOWN_TIMEOUT)
	slog.Printf("Webserver: Starting Graceful Shutdown with [Timeout:%s]..", _WS_SHUTDOWN_TIMEOUT)
	if err := s.webserver.Shutdown(ctx); err != nil {
		slog.Printf("Error: Webserver Shutdown [E:%s]", err)
	} else {
		slog.Printf("Webserver Shutdown successful")
	}
}

// Close is an ungraceful shutdown of webserver
func (s *Server) CloseWebServer() {
	slog.Printf("Webserver: Closing without waiting for active connections..")
	if err := s.webserver.Close(); err != nil {
		slog.Printf("Error: Webserver Close [E:%s]", err)
	} else {
		slog.Printf("Webserver Close successful")
	}
}

/*
func authMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if _, ok := middlewareIgnoreList[r.URL.Path]; ok {
		slog.Printf("Auth Middleware: [%s] in ignore list", r.URL.Path)
		next(w, r)
		return
	}
	token := r.Header.Get("Token")
	if token == "" {
		slog.Printf("ERR: Request:[%s] blocked, no token in Request", r.URL)
		respondWithError(w, http.StatusUnauthorized, fmt.Errorf("Please provide a token"))
		return
	}
	if token != apiKey {
		slog.Printf("ERR: Request:[%s] blocked, invalid token [Token:%s]", r.URL, token)
		respondWithError(w, http.StatusUnauthorized, fmt.Errorf("Invalid token"))
		return
	}
	next(w, r)
}
*/

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to the slot machine!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, Trippy %s!\n", ps.ByName("name"))
}

func Spin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	machine := ps.ByName(_PARA_SPIN_MACHINE)
	if machine == "" {
		slog.Printf("Spin: Parameter:[%s] is empty", _PARA_SPIN_MACHINE)
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("Parameter:[%s] cannot be empty", _PARA_SPIN_MACHINE))
		return
	}

	token := r.Header.Get("token")
	if token == "" {
		slog.Println("Header:[token] is empty")
		respondWithError(w, http.StatusBadRequest, errors.New("Header:[token] cannot be empty"))
		return
	}

	user, err := parseToken(token)
	if err != nil {
		slog.Printf("Parsing token failed [Error:%s]", err)
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid token received"))
		return
	}

	response := respSpin{
		JWT: user,
	}

	if machine == _ATKINS_DIET_MACHINE {
		wager, ok := atkinsDietMachine.Wager(user.Bet, user.Chips)
		if !ok {
			slog.Printf("User:[%s] Balance not enough. Wager:[%d] Balance:[%d]", user.UID, wager, user.Chips)
			respondWithError(w, http.StatusBadRequest, fmt.Errorf("Balance insufficient. Need [%d] chips more.", wager-user.Chips))
			return
		}
		stops, pay, err := atkinsDietMachine.Spin(user.Bet)
		if err != nil {
			slog.Printf("Spin failed for User:[%s] Error:[%s]", user.UID, err)
			respondWithError(w, http.StatusInternalServerError, errors.New("Unable to spin"))
			return
		}
		response.Total = pay
		response.Spins = []spin{
			{
				Type:  "main",
				Total: pay,
				Stops: stops,
			},
		}
		response.JWT.Chips = response.JWT.Chips - wager + pay
		writeSpinResponse(w, http.StatusOK, response)
		return
	}
	respondWithError(w, http.StatusBadRequest, fmt.Errorf("Unknown machine:[%s]", machine))
}

func parseToken(tokenString string) (userClaims, error) {
	user := new(userClaims)
	token, err := jwt.ParseWithClaims(tokenString, user, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing invalid. method %v", t.Header["alg"])
		}
		return []byte(apiKey), nil
	})

	if err != nil {
		return *user, fmt.Errorf("Could not parse JWT token [Error:%s]", err)
	}

	if !token.Valid {
		return *user, errors.New("Invalid JWT token")
	}

	return *user, nil
}
