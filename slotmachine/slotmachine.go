package slotmachine

type SlotMachine interface {
	Wager(bet, balance int) (wager int, sufficiency bool)
	Spin(bet int) (stops []int, pay int, err error)
}
