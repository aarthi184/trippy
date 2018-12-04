package slotmachine

type SlotMachine interface {
	Wager(bet, balance int) (wager int, err error)
	Spin(bet int) (stops []int, pay int, err error)
}
