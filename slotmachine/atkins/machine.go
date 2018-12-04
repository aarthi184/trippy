package atkins

import (
	"errors"
	"log"
	"os"

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

func (ad *AtkinsDietMachine) Spin(bet int) ([]int, int, error) {
	stops, err := spinner.Spin(ad.Reels)
	if err != nil {
		return stops, 0, err
	}

	wins, _, err := spinner.FindWins(stops, ad.Reels, ad.PayLines, ad.SpecialSymbols)
	if err != nil {
		return stops, 0, err
	}

	pay, err := spinner.CalculatePay(wins, ad.PayTable, ad.SpecialSymbols)
	if err != nil {
		return stops, 0, err
	}

	pay = pay * bet
	return stops, pay, nil
}
