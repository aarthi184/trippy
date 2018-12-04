package slotmachine

/*
   Paytable structure
   Paytable will have only concurrent reads across goroutines
   Writes will happen only on initialization
   Therefore, thread-safety is not necessary here
*/

import (
	"fmt"
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

type SpecialSymbols struct {
	Wildcard Symbol
	Scatter  Symbol
}

/*
func (pt *PayTable) SetPay(sym Symbol, count, amount int) error {
	p, ok := pt[sym]
	if ok {
		p.SetPay(count, amount)
		return nil
	}
	return errors.New("Symbol doesn't exist")
}

func (pt *PayTable) GetPay(sym Symbol, count int) (int, bool) {
	p, ok := pt[sym]
	if ok {
		return p.GetPay(count)
	}
	return 0, false
}
func (p *Pays) SetPay(count, amount int) {
	p[count] = amount
}

func (p *Pays) GetPay(count int) (int, bool) {
	v, ok := p[count]
	return v, ok
}
*/
