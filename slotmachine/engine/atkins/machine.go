package atkins

import (
	"errors"
	"log"
	"os"
	"time"

	"trippy/slotmachine"
	"trippy/spinner"
)

var (
	slog *log.Logger
)

func init() {
	// Log to stdout
	slog = log.New(os.Stdout, "", 0)
	slog.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)
	slog.SetPrefix("ATKINS:")
}

type AtkinsDietMachine struct {
	PayTable slotmachine.PayTable
	Reels    slotmachine.Reels
	PayLines slotmachine.PayLines

	slotmachine.SpecialSymbols
}

func NewAtkinsDietMachine() *AtkinsDietMachine {
	return &AtkinsDietMachine{
		PayTable: PayTable,
		Reels:    Reels,
		PayLines: PayLines,
		SpecialSymbols: slotmachine.SpecialSymbols{
			Wildcard: _ATKINS,
			Scatter:  _SCALE,
		},
	}
}

var (
	ErrChipsInsufficient = errors.New("Chips insufficient")
	ErrInvalidBet        = errors.New("Bet is not greater than 0")
)

func (ad *AtkinsDietMachine) Wager(bet, chips int) (int, error) {
	if bet <= 0 {
		return 0, ErrInvalidBet
	}
	wager := bet * len(ad.PayLines)
	if wager > chips {
		return wager, ErrChipsInsufficient
	}
	return wager, nil
}

func (ad *AtkinsDietMachine) Spin(bet int) (int, []slotmachine.SpinResult, error) {

	var (
		spinResults []slotmachine.SpinResult
		freeSpin    = true
		mainSpin    = false
	)

	// Main Spin
	spinResult, err := ad.spin(bet, mainSpin)
	if err != nil {
		return 0, spinResults, err
	}

	spinResult.Type = slotmachine.MAIN_SPIN
	spinResults = append(spinResults, spinResult)

	if spinResult.FreeSpins == 0 {
		return spinResult.Pay, spinResults, nil
	}

	// Free Spins - if any
	var (
		freeSpins   = spinResult.FreeSpins
		totalPayout = spinResult.Pay
	)

	for i := 0; freeSpins > 0; i++ {
		freeSpins--
		slog.Println("Remaining free spins:", freeSpins)

		// Incase we are stuck in an infinite loop of free spins
		// We slow down the free spins to allow other goroutines to work
		// Sleeping on every 16th iteration
		if i&16 == 16 {
			i = 1
			time.Sleep(500 * time.Millisecond)
		}

		spinResult, err = ad.spin(bet, freeSpin)
		if err != nil {
			return totalPayout, spinResults, err
		}

		spinResult.Type = slotmachine.FREE_SPIN
		spinResults = append(spinResults, spinResult)
		freeSpins = freeSpins + spinResult.FreeSpins
		totalPayout = totalPayout + spinResult.Pay
	}

	return totalPayout, spinResults, nil
}

func (ad *AtkinsDietMachine) spin(bet int, freeSpin bool) (slotmachine.SpinResult, error) {
	spinResult, err := spinner.SpinNPay(
		ad.Reels,
		ad.PayLines,
		ad.PayTable,
		ad.SpecialSymbols,
	)
	if err != nil {
		return spinResult, err
	}

	spinResult.FreeSpins = ad.getFreeSpins(spinResult.ScatterCount)
	slog.Println("Got FreeSpins:", spinResult.FreeSpins)

	// Multiplying payout by wager
	if freeSpin {
		bet = bet * 3
	}
	for i := 0; i < len(spinResult.WinLines); i++ {
		spinResult.WinLines[i].Payout = spinResult.WinLines[i].Payout * bet
	}
	spinResult.Pay = spinResult.Pay * bet
	slog.Printf("Spin Result %#v [Bet:%d]", spinResult, bet)
	return spinResult, nil
}

func (ad *AtkinsDietMachine) getFreeSpins(scatterCount int) int {
	if scatterCount >= _SCATTER_COUNT_FOR_FREE_SPIN {
		return _FREE_SPINS
	}
	return 0
}
