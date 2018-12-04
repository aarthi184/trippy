package atkins

import (
	"testing"
)

var (
	adm          = NewAtkinsDietMachine()
	wagerSamples = []wagerSample{
		{bet: 10, chips: 100, err: ErrChipsInsufficient, wager: 10 * len(adm.PayLines)},
		{bet: 10, chips: 200, err: nil, wager: 10 * len(adm.PayLines)},
		{bet: 0, chips: 200, err: ErrInvalidBet, wager: 0},
		{bet: 200, chips: 100, err: ErrChipsInsufficient, wager: 200 * len(adm.PayLines)},
		{bet: -2, chips: 100, err: ErrInvalidBet, wager: 0},
	}
)

type wagerSample struct {
	bet, chips, wager int
	err               error
}

func TestWager(t *testing.T) {
	for _, sample := range wagerSamples {
		testWager(t, sample)
	}
}

func testWager(t *testing.T, sample wagerSample) {
	wager, err := adm.Wager(sample.bet, sample.chips)
	if err != sample.err {
		t.Errorf("Wager sufficient or not. Bet:[%d] Chips:[%d] Expected:[%s] Got:[%s]",
			sample.bet, sample.chips, sample.err, err)
		return
	}
	if err != nil {
		return
	}
	if wager != sample.wager {
		t.Errorf("Bet:[%d] Chips:[%d] Expected:[%d] Got:[%d]",
			sample.bet, sample.chips, sample.wager, wager)
	}
}
