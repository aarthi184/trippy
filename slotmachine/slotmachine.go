package slotmachine

type SlotMachine interface {
	Wager(bet, balance int) (wager int, err error)
	Spin(bet int) (payout int, results []SpinResult, err error)
}
