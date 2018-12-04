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
	special slotmachine.SpecialSymbols) (int, error) {

	stops, err := Spin(reels)
	if err != nil {
		return 0, err
	}

	wins, _, err := FindWins(stops, reels, payLines, special)
	if err != nil {
		return 0, err
	}

	pay, err := CalculatePay(wins, payTable, special)
	if err != nil {
		return 0, err
	}

	return pay, nil
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
	special slotmachine.SpecialSymbols) ([][]slotmachine.Symbol, int, error) {

	var (
		scatter               int
		j                     int
		prevSymbol, curSymbol slotmachine.Symbol
		firstSymbol           slotmachine.Symbol
		breakInLine           bool
		payLineSymbolsTable   = make([][]slotmachine.Symbol, 0, len(payLines))
		payLineSymbols        []slotmachine.Symbol
	)
	if len(reels) == 0 {
		return payLineSymbolsTable, 0, errEmptyReel
	}
	if len(payLines) == 0 {
		return payLineSymbolsTable, 0, errEmptyPayLine
	}
	for i, line := range payLines {
		if len(reels[i]) != len(line) {
			return payLineSymbolsTable, 0, errReelPayLineMismatch
		}

		// Keeping track of the first symbol and comparing each symbol in line with previous symbol
		firstSymbol = getSymbol(reels, stops[0], line[0], 0)
		prevSymbol = firstSymbol
		//slog.Printf("FirstSymbol:%s", firstSymbol)
		payLineSymbols = make([]slotmachine.Symbol, 0, len(stops))
		payLineSymbols = append(payLineSymbols, firstSymbol)

		// Maintaining a breaker so that we can keep counting scatter even if line is a no-win.
		breakInLine = false

		for j = 1; j < len(line); j++ {
			curSymbol = getSymbol(reels, stops[j], line[j], j)

			if !breakInLine {
				// Any wildcard symbol or a symbol equal to the firstSymbol is a win
				if curSymbol == special.Wildcard || curSymbol == prevSymbol {
					payLineSymbols = append(payLineSymbols, curSymbol)
					// If firstSymbol was a wildcard, any symbol next to the wildcard becomes the primary symbol
				} else if firstSymbol == special.Wildcard && curSymbol != special.Wildcard {
					firstSymbol = curSymbol
					payLineSymbols = append(payLineSymbols, curSymbol)
				} else {
					breakInLine = true
				}
			}
			// Counting scatter
			if curSymbol == special.Scatter {
				scatter++
			}
		}
		//slog.Printf("PayLine Symbols:%v", payLineSymbols)
		if len(payLineSymbols) != 0 {
			payLineSymbolsTable = append(payLineSymbolsTable, payLineSymbols)
		}
	}
	slog.Println("Scatter count:", scatter)
	return payLineSymbolsTable, scatter, nil
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
// TODO: Handle cases where |offset| > length for offset<0, eg: maxIndex:4 offset: -6
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

/*
var (
	errAllWildcards = errors.New("All symbols in stops are wildcards")
)

	if firstSymbol == special.Wilcard {
		j, err = findFirstNonWildcardSymbol(stops, reels[i], line, wildcard Symbol)
		if err == errAllWildcards {
			slog.Printf("All Wildcards in pay line:[%d][%v]", i, line)
			slog.Println("Amount:%d", payTable[wildcard][len(stops)]
		}
	}
func findFirstNonWildcard(stops []slotmachine.Symbol, reelLine slotmachine.ReelLine, payLine slotmachine.PayLine, wildcard slotmachine.Symbol) (int, error) {
	for i := 0; i < len(stops); i++ {
		if reelLine[i] != wildcard {
			return i, nil
		}
	}
	return -1, errAllWildcards
}
*/

func CalculatePay(wins [][]slotmachine.Symbol, payTable slotmachine.PayTable, special slotmachine.SpecialSymbols) (int, error) {
	var (
		pay       int
		pays      slotmachine.Pays
		ok        bool
		paySymbol slotmachine.Symbol
	)
	for _, winLine := range wins {
		slog.Printf("Win Line:%v CurrentPay:%d", winLine, pay)
		if len(winLine) == 0 {
			continue
		}
		if winLine[0] != special.Wildcard {
			//slog.Printf("Not wildcard:[%d] symbol:[%d]", special.Wildcard, winLine[0])
			if pays, ok = payTable[winLine[0]]; ok {
				pay = pay + pays[len(winLine)]
				//slog.Println("iNot wildcard Pay", pay)
			}
			continue
		}
		//slog.Println("First symbol is wildcard", winLine[0])
		paySymbol = special.Wildcard
		for j := 0; j < len(winLine); j++ {
			if winLine[j] != special.Wildcard {
				paySymbol = winLine[j]
			}
		}
		if pays, ok = payTable[paySymbol]; ok {
			pay = pay + pays[len(winLine)]
			//slog.Println("With wildcard Pay", pay)
		}
	}
	return pay, nil
}
