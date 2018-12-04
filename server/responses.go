package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type respGeneric struct {
	Resp  string `json:"resp,omitempty"`
	Error string `json:"err,omitempty"`
}

type respSpin struct {
	Total int        `json:"total"`
	Spins []spin     `json:"spins"`
	JWT   userClaims `json:"jwt"`
}

type spin struct {
	Type  string `json:"type"`
	Total int    `json:"total"`
	Stops []int  `json:"stops"`
}

type userClaims struct {
	UID   string `json:"uid"`
	Chips int    `json:"chips"`
	Bet   int    `json:"bet"`
}

// UserClaims is always valid since we do not expire a JWT token
// 'Valid' method is present to adhere to the jwt.Claims interface
func (u userClaims) Valid() error { return nil }

// ----------------------------- Response Methods ---------------------------------- //

func writeSpinResponse(w http.ResponseWriter, statusCode int, resp respSpin) {
	var respEncoder *json.Encoder = json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	respEncoder.Encode(resp)
}

func respondWithError(w http.ResponseWriter, statusCode int, err error) {
	var (
		resp        respGeneric
		respEncoder *json.Encoder = json.NewEncoder(w)
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp.Error = fmt.Sprintf("%s", err)
	respEncoder.Encode(resp)
}
