package slotmachine

/*
   All types defined here will only have concurrent reads
   Writes happen only on init.
   Therefore, thread-safety is not necessary here
*/

import (
	"fmt"
)

const (
	MAIN_SPIN string = "main"
	FREE_SPIN string = "free"
)

// TODO: Make this interface for each machine to have it's own Stringer method
type Symbol int

func GetSymbol(n int) Symbol {
	return Symbol(n)
}

func (s Symbol) String() string {
	return fmt.Sprintf("%d", s)
}

type PayTable map[Symbol]Pays
type Pays map[int]int

type Reels []ReelLine
type ReelLine []Symbol

type PayLines []PayLine
type PayLine []int

func SetAllPayLines(stripCount int) PayLines {
	// TODO: Set all possible pay lines
	return []PayLine{{}}
}

type WinLine struct {
	Index  int      `json:"index"`  //  number of the line
	Symbol Symbol   `json:"symbol"` // paid symbol, can be code or index
	Count  int      `json:"count"`  // number of symbols paid
	Payout int      `json:"payout"` // Payout for this line
	Line   []Symbol `json:"-"`      // The line of symbols
}

type SpecialSymbols struct {
	Wildcard Symbol
	Scatter  Symbol
}

type SpinResult struct {
	Type         string
	Stops        []int
	Pay          int
	WinLines     []WinLine
	ScatterCount int
	FreeSpins    int
}
