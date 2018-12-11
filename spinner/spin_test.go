package spinner

import (
	"reflect"
	"testing"

	SM "trippy/slotmachine"
)

var (
	spinSamples = []spinSample{
		{reels: SM.Reels{{1, 2, 3, 4, 5}, {5, 4, 3, 2, 1}}, err: nil},
		{reels: SM.Reels{}, err: errEmptyReel},
		{reels: SM.Reels{{1, 2}, {5, 4}, {6, 4}, {7, 6}, {4, 1}, {3, 4}, {5, 7}}, err: nil},
	}
)

type spinSample struct {
	reels SM.Reels
	err   error
}

func TestSpin(t *testing.T) {
	for _, sample := range spinSamples {
		testSpin(t, sample)
	}
}

func testSpin(t *testing.T, sample spinSample) {
	stops, err := Spin(sample.reels)
	if err != sample.err {
		t.Errorf("Expected:[%s] Got:[%s]", sample.err, err)
		return
	}
	if err != nil {
		// If test is for an error case, stop tests and return
		return
	}
	if len(sample.reels) != 0 && len(stops) != len(sample.reels[0]) {
		t.Errorf("Expected:[%d] Got:[%d]", len(sample.reels[0]), len(stops))
	}
	t.Logf("Stops:%v", stops)
}

var (
	winSamples = []winSample{
		{stops: []int{1}, reels: SM.Reels{}, payLines: []SM.PayLine{{0, 0, 0}}, special: SM.SpecialSymbols{}, err: errEmptyReel},
		{stops: []int{1, 1, 1}, reels: SM.Reels{{1, 2, 3}, {5, 4, 3}, {2, 1, 3}}, payLines: []SM.PayLine{}, special: SM.SpecialSymbols{}, err: errEmptyPayLine},
		{stops: []int{1, 1, 1}, reels: SM.Reels{{1, 2, 3}}, payLines: []SM.PayLine{}, special: SM.SpecialSymbols{}, err: errEmptyPayLine},

		// Non-error samples with Winning Line
		{stops: []int{1, 1, 1}, reels: SM.Reels{{1, 2, 3}, {5, 4, 3}, {2, 1, 3}}, payLines: []SM.PayLine{{1, 1, 1}}, special: SM.SpecialSymbols{},
			err:  nil,
			wins: []SM.WinLine{}},
		{stops: []int{1, 1, 1}, reels: SM.Reels{{1, 1, 3}, {5, 4, 3}, {2, 1, 3}}, payLines: []SM.PayLine{{1, 1, 1}}, special: SM.SpecialSymbols{},
			err:  nil,
			wins: []SM.WinLine{SM.WinLine{Index: 1, Line: []SM.Symbol{SM.Symbol(1), SM.Symbol(1)}}},
		},
		{stops: []int{1, 1, 1}, reels: SM.Reels{{1, 1, 1}, {5, 4, 3}, {2, 1, 3}}, payLines: []SM.PayLine{{1, 1, 1}}, special: SM.SpecialSymbols{},
			err:  nil,
			wins: []SM.WinLine{SM.WinLine{Index: 1, Line: []SM.Symbol{SM.Symbol(1), SM.Symbol(1), SM.Symbol(1)}}},
		},
		{stops: []int{1, 1, 1}, reels: SM.Reels{{1, 1, 1}, {5, 5, 4}, {2, 1, 3}}, payLines: []SM.PayLine{{2, 2, 2}}, special: SM.SpecialSymbols{},
			err:  nil,
			wins: []SM.WinLine{SM.WinLine{Index: 1, Line: []SM.Symbol{SM.Symbol(5), SM.Symbol(5)}}},
		},
		{stops: []int{1, 1, 1}, reels: SM.Reels{{4, 1, 4}, {5, 4, 4}, {2, 4, 3}}, payLines: []SM.PayLine{{1, 3, 2}}, special: SM.SpecialSymbols{},
			err:  nil,
			wins: []SM.WinLine{SM.WinLine{Index: 1, Line: []SM.Symbol{SM.Symbol(4), SM.Symbol(4), SM.Symbol(4)}}},
		},

		// Winning Lines with Wildcard Symbol
		{stops: []int{1, 1, 1}, reels: SM.Reels{{4, 1, 4}, {5, 4, 888}, {2, 4, 3}}, payLines: []SM.PayLine{{1, 3, 2}}, special: SM.SpecialSymbols{Wildcard: 888}, err: nil,
			wins: []SM.WinLine{SM.WinLine{Index: 1, Line: []SM.Symbol{SM.Symbol(4), SM.Symbol(4), SM.Symbol(888)}}},
		},
		{stops: []int{1, 2, 0}, reels: SM.Reels{{4, 888, 4}, {5, 4, 888}, {2, 4, 3}}, payLines: []SM.PayLine{{1, 3, 2}}, special: SM.SpecialSymbols{Wildcard: 888}, err: nil,
			wins: []SM.WinLine{SM.WinLine{Index: 1, Line: []SM.Symbol{SM.Symbol(4), SM.Symbol(888), SM.Symbol(4)}}},
		},
		{stops: []int{1, 2, 0}, reels: SM.Reels{{888, 888, 4}, {5, 4, 888}, {2, 4, 3}}, payLines: []SM.PayLine{{1, 3, 2}}, special: SM.SpecialSymbols{Wildcard: 888}, err: nil,
			wins: []SM.WinLine{SM.WinLine{Index: 1, Line: []SM.Symbol{SM.Symbol(888), SM.Symbol(888)}}},
		},

		// Multiple paylines
		{
			stops:    []int{1, 2, 0},
			reels:    SM.Reels{{888, 888, 4}, {4, 4, 888}, {2, 4, 3}},
			payLines: []SM.PayLine{{1, 3, 2}, {2, 2, 2}},
			special:  SM.SpecialSymbols{Wildcard: 888},
			err:      nil,
			wins: []SM.WinLine{
				SM.WinLine{Index: 1, Line: []SM.Symbol{SM.Symbol(888), SM.Symbol(888)}},
				SM.WinLine{Index: 2, Line: []SM.Symbol{SM.Symbol(4), SM.Symbol(4), SM.Symbol(4)}},
			},
		},
	}
)

type winSample struct {
	stops        []int
	reels        SM.Reels
	payLines     SM.PayLines
	special      SM.SpecialSymbols
	err          error
	wins         []SM.WinLine
	scatterCount int
}

func TestWin(t *testing.T) {
	for _, sample := range winSamples {
		testWin(t, sample)
	}
}

func testWin(t *testing.T, sample winSample) {
	wins, err := FindWins(sample.stops, sample.reels, sample.payLines, sample.special)
	if err != sample.err {
		t.Errorf("Expected:[%s] Got:[%s]", sample.err, err)
		return
	}
	if err != nil {
		// If test is for an error case, stop tests and return
		return
	}
	equal := reflect.DeepEqual(wins, sample.wins)
	if !equal {
		t.Errorf("Expected:[%v] Got:[%v]", sample.wins, wins)
	}
}

var (
	scatterSamples = []scatterSample{
		{
			stops:        []int{0, 0, 0, 0, 0},
			reels:        SM.Reels{{4, 777, 4, 777, 6}, {777, 4, 888, 777, 2}, {2, 777, 3, 3, 6}, {2, 9, 3, 3, 6}},
			special:      SM.SpecialSymbols{Scatter: 777},
			scatterCount: 4,
		},
		{
			stops:        []int{1, 2, 0},
			reels:        SM.Reels{{4, 888, 4}, {5, 4, 7}, {2, 4, 3}},
			special:      SM.SpecialSymbols{Scatter: 777},
			scatterCount: 0,
		},
		{
			stops:        []int{1, 2, 0},
			reels:        SM.Reels{{8, 888, 4}, {5, 4, 888}, {2, 888, 3}, {888, 4, 3}},
			special:      SM.SpecialSymbols{Scatter: 888},
			scatterCount: 2,
		},
	}
)

type scatterSample struct {
	stops        []int
	reels        SM.Reels
	special      SM.SpecialSymbols
	scatterCount int
}

func TestCountScatter(t *testing.T) {
	for _, sample := range scatterSamples {
		testCountScatter(t, sample)
	}
}

func testCountScatter(t *testing.T, sample scatterSample) {
	scatter := CountScatter(sample.stops, sample.reels, sample.special.Scatter)
	if scatter != sample.scatterCount {
		t.Errorf("Expected:[%v] Got:[%v]", sample.scatterCount, scatter)
	}
}

var (
	overflowSamples = []overflowSample{
		{max: 6, offset: 7, expected: 0},
		{max: 0, offset: 7, expected: 0},
		{max: 6, offset: -2, expected: 5},
		{max: 1, offset: -2, expected: 0},
		{max: 200, offset: 205, expected: 4},
		{max: 200, offset: -5, expected: 196},
		{max: 200, offset: 0, expected: 0},
		{max: 200, offset: 200, expected: 200},
		{max: 200, offset: 201, expected: 0},
		//{max: 1, offset: -3, expected: 1},
		//{max: 4, offset: -6, expected: 2},
	}
)

type overflowSample struct {
	max, offset, expected int
}

func TestRotateOverflow(t *testing.T) {
	for _, sample := range overflowSamples {
		testRotateOverflow(t, sample)
	}
}

func testRotateOverflow(t *testing.T, sample overflowSample) {
	n := rotateOverflow(sample.max, sample.offset)
	if n != sample.expected {
		t.Errorf("Max:[%d] Offset:[%d] Expected:[%d] Got:[%d]", sample.max, sample.offset, sample.expected, n)
	}
}

var (
	paySamples = []paySample{
		{
			wins: []SM.WinLine{
				SM.WinLine{Line: []SM.Symbol{SM.Symbol(1), SM.Symbol(3), SM.Symbol(3)}},
			},
			special: SM.SpecialSymbols{Wildcard: 1},
			err:     nil,
			pay:     30,
		},
		{
			wins: []SM.WinLine{
				SM.WinLine{Line: []SM.Symbol{SM.Symbol(1), SM.Symbol(1)}},
			},
			special: SM.SpecialSymbols{Wildcard: 1},
			err:     nil,
			pay:     5,
		},
		{
			wins: []SM.WinLine{
				SM.WinLine{Line: []SM.Symbol{SM.Symbol(2), SM.Symbol(2), SM.Symbol(2)}},
			},
			special: SM.SpecialSymbols{Wildcard: 1},
			err:     nil,
			pay:     40,
		},
		{
			wins: []SM.WinLine{
				SM.WinLine{Line: []SM.Symbol{SM.Symbol(3), SM.Symbol(1), SM.Symbol(3)}},
			},
			special: SM.SpecialSymbols{Wildcard: 1},
			err:     nil,
			pay:     30,
		},
		{
			wins: []SM.WinLine{
				SM.WinLine{Line: []SM.Symbol{SM.Symbol(2), SM.Symbol(1)}},
			},
			special: SM.SpecialSymbols{Wildcard: 1},
			err:     nil,
			pay:     3,
		},
		{
			wins: []SM.WinLine{
				SM.WinLine{Line: []SM.Symbol{SM.Symbol(1), SM.Symbol(2), SM.Symbol(1)}},
			},
			special: SM.SpecialSymbols{Wildcard: 1},
			err:     nil,
			pay:     40,
		},
	}

	samplePayTable = SM.PayTable{
		1: SM.Pays{
			5: 5000,
			4: 500,
			3: 50,
			2: 5,
		},
		2: SM.Pays{
			5: 1000,
			4: 200,
			3: 40,
			2: 3,
		},
		3: SM.Pays{
			5: 500,
			4: 150,
			3: 30,
			2: 2,
		},
	}
)

type paySample struct {
	wins    []SM.WinLine
	special SM.SpecialSymbols
	err     error
	pay     int
}

func TestCalculatePay(t *testing.T) {
	for _, sample := range paySamples {
		testCalculatePay(t, sample)
	}
}

func testCalculatePay(t *testing.T, sample paySample) {
	spinResult, err := CalculatePay(sample.wins, samplePayTable, sample.special)
	if err != sample.err {
		t.Errorf("Expected:[%s] Got:[%s]", sample.err, err)
		return
	}
	if err != nil {
		// If test is for an error case, stop tests and return
		return
	}
	if spinResult.Pay != sample.pay {
		t.Errorf("Expected:[%d] Got:[%d]", sample.pay, spinResult.Pay)
		//return
	}
	//t.Logf("Pay: %d", pay)
}
