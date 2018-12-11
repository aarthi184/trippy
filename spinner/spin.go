package spinner

import (
	"errors"
	"fmt"
	"log"
	"os"

	"trippy/slotmachine"
)

var (
	slog *log.Logger
)

func init() {
	// Log to stdout
	slog = log.New(os.Stdout, "", 0)
	slog.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)
	slog.SetPrefix("SPIN:")
}

var (
	errEmptyReel           = errors.New("Reel strips are empty")
	errEmptyPayLine        = errors.New("Pay lines are empty")
	errReelPayLineMismatch = errors.New("Reel width and Pay line width do not match")
)

func SpinNPay(
	reels slotmachine.Reels,
	payLines slotmachine.PayLines,
	payTable slotmachine.PayTable,
	special slotmachine.SpecialSymbols) (slotmachine.SpinResult, error) {

	var spinResult slotmachine.SpinResult

	stops, err := Spin(reels)
	if err != nil {
		return spinResult, err
	}

	winLines, err := FindWins(stops, reels, payLines, special)
	if err != nil {
		return spinResult, err
	}

	spinResult, err = CalculatePay(winLines, payTable, special)
	if err != nil {
		return spinResult, err
	}

	spinResult.ScatterCount = CountScatter(stops, reels, special.Scatter)

	// Changing stops to Human-friendly numbering (starts from 1)
	for i := 0; i < len(stops); i++ {
		stops[i]++
	}
	spinResult.Stops = stops

	return spinResult, nil
}

func Spin(reels slotmachine.Reels) (stops []int, err error) {

	if len(reels) == 0 {
		return stops, errEmptyReel
	}

	stops = make([]int, len(reels[0]))
	//slog.Println("Min:", 0, "Max:", len(reels)-1)

	// Spinning the reels
	for i := range stops {
		stops[i], err = randInt(0, len(reels)-1)
		if err != nil {
			return stops, fmt.Errorf("Unable to generate random stop [Error:%s]", err)
		}
	}

	return stops, nil
}

func FindWins(
	stops []int,
	reels slotmachine.Reels,
	payLines slotmachine.PayLines,
	special slotmachine.SpecialSymbols) ([]slotmachine.WinLine, error) {

	var (
		j                      int
		primeSymbol, curSymbol slotmachine.Symbol
		payLineSymbolsTable    = make([]slotmachine.WinLine, 0, len(payLines))
		payLineSymbols         []slotmachine.Symbol
		winLine                slotmachine.WinLine
	)
	if len(reels) == 0 {
		return payLineSymbolsTable, errEmptyReel
	}
	if len(payLines) == 0 {
		return payLineSymbolsTable, errEmptyPayLine
	}
	for i, line := range payLines {
		if len(reels[i]) != len(line) {
			return payLineSymbolsTable, errReelPayLineMismatch
		}

		// Keeping track of the prime symbol and comparing each symbol in line with it
		primeSymbol = getSymbol(reels, stops[0], line[0], 0)
		//slog.Printf("Prime Symbol:%s", primeSymbol)
		payLineSymbols = make([]slotmachine.Symbol, 0, len(stops))
		payLineSymbols = append(payLineSymbols, primeSymbol)

		for j = 1; j < len(line); j++ {
			curSymbol = getSymbol(reels, stops[j], line[j], j)

			// Any wildcard symbol or a symbol equal to the firstSymbol is a win
			if curSymbol == special.Wildcard || curSymbol == primeSymbol {
				payLineSymbols = append(payLineSymbols, curSymbol)
			} else if primeSymbol == special.Wildcard && curSymbol != special.Wildcard {
				// If firstSymbol was a wildcard, any symbol next to the wildcard becomes the primary symbol
				primeSymbol = curSymbol
				payLineSymbols = append(payLineSymbols, curSymbol)
			} else {
				break
			}
		}
		//slog.Printf("PayLine Symbols:%v", payLineSymbols)
		if len(payLineSymbols) > 1 {
			winLine = slotmachine.WinLine{
				Index: i + 1, // Starting human-friendly indexing (starts from 1)
				Line:  payLineSymbols,
			}
			payLineSymbolsTable = append(payLineSymbolsTable, winLine)
		}
	}
	return payLineSymbolsTable, nil
}

func CountScatter(stops []int, reels slotmachine.Reels, scatterSymbol slotmachine.Symbol) int {
	var (
		scatter   int
		curSymbol slotmachine.Symbol
	)
	if len(reels) == 0 {
		return 0
	}
	reelStrips := len(reels[0])
	// Counting scatter in the the 3 slots 1,2 and 3
	for i := 1; i <= 3; i++ {
		for j := 0; j < reelStrips; j++ {
			curSymbol = getSymbol(reels, stops[j], i, j)
			// Counting scatter
			if curSymbol == scatterSymbol {
				scatter++
			}
		}
	}
	return scatter
}

func getSymbol(reels slotmachine.Reels, stop, payLineSpot, stripNumber int) slotmachine.Symbol {
	// payLines are numbered from 1 to n where n is the number of slots
	// Here, we assume that only 3 slots are present
	// Note: To allow multislots, subtract paylineSpot by center i.e. ((n/2) + 1)
	payLineOffset := payLineSpot - 2
	offset := rotateOverflow(len(reels)-1, stop+payLineOffset)
	//slog.Println("Got offset:", offset)
	return reels[offset][stripNumber]
}

// rotateOverflow rotates the reel to get a number within 0 and maxIndex
// TODO: Handle cases where |offset| > length for negative offsets, eg: maxIndex:4 offset: -6
// In our case, since our slotmachines will always have 3 slots, we will never hit the above mentioned scenario
func rotateOverflow(maxIndex, offset int) int {
	if maxIndex <= 0 {
		return maxIndex
	}
	length := maxIndex + 1
	if offset < 0 {
		return length + offset
	}
	return offset % length
}

// CalculatePay finds the total pay for this spin from the list of winning lines
func CalculatePay(wins []slotmachine.WinLine, payTable slotmachine.PayTable, special slotmachine.SpecialSymbols) (slotmachine.SpinResult, error) {
	var (
		totalPayout int
		linePayout  int
		pays        slotmachine.Pays
		ok          bool
		paySymbol   slotmachine.Symbol
		spinResult  slotmachine.SpinResult
		line        []slotmachine.Symbol
	)
	for i := 0; i < len(wins); i++ {
		line = wins[i].Line
		slog.Printf("Win Line:%v CurrentPay:%d", line, totalPayout)
		if len(line) == 0 {
			continue
		}
		if line[0] != special.Wildcard {
			//slog.Printf("Not wildcard:[%d] symbol:[%d]", special.Wildcard, winLine[0])
			paySymbol = line[0]
		} else {
			//slog.Println("First symbol is wildcard", winLine[0])
			paySymbol = special.Wildcard

			// If first symbol was a wildcard, check if line contains another symbol
			// Eg: WC WC 1 1
			for j := 0; j < len(line); j++ {
				if line[j] != special.Wildcard {
					paySymbol = line[j]
					break
				}
			}
		}
		if pays, ok = payTable[paySymbol]; ok {
			linePayout = pays[len(line)]
		}
		totalPayout = totalPayout + linePayout
		wins[i].Symbol = paySymbol
		wins[i].Count = len(line)
		wins[i].Payout = linePayout
	}
	spinResult.Pay = totalPayout
	spinResult.WinLines = wins
	//slog.Printf("22 %v", spinResult)
	return spinResult, nil
}
