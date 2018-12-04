package server

import (
	"testing"
)

var (
	tokenSamples = []tokenSample{
		{token: "abc", secret: "secret", errExpected: true},
		{token: "", secret: "secret", errExpected: true},
		{
			token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxMjMiLCJjaGlwcyI6MTAwMCwiYmV0IjoxMH0.PSPp1k2raKRhrKMFOrINIOkd8bK-w2vX1SULL3IKp4Q",
			secret: "secret",
			user: userClaims{
				UID:   "123",
				Chips: 1000,
				Bet:   10,
			},
		},
	}
)

type tokenSample struct {
	token, secret string
	user          userClaims
	errExpected   bool
}

func TestParseToken(t *testing.T) {
	for _, sample := range tokenSamples {
		testParseToken(t, sample)
	}
}

func testParseToken(t *testing.T, sample tokenSample) {
	user, err := parseToken(sample.token, []byte(sample.secret))
	if sample.errExpected && err == nil {
		t.Errorf("Expected:[error] Got:[%s]", err)
	}
	if err != nil {
		return
	}
	if user.UID != sample.user.UID {
		t.Errorf("Expected:[%s] Got:[%s]", sample.user.UID, user.UID)
	}
	if user.Bet != sample.user.Bet {
		t.Errorf("Expected:[%d] Got:[%d]", sample.user.Bet, user.Bet)
	}
	if user.Chips != sample.user.Chips {
		t.Errorf("Expected:[%d] Got:[%d]", sample.user.Chips, user.Chips)
	}
}
