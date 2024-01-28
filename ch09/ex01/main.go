package main

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

var withdraws = make(chan int)
var withdraws_status = make(chan bool)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func Withdrawn(amount int) bool {
	withdraws <- amount
	return <-withdraws_status
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-withdraws:
			if amount > balance {
				withdraws_status <- false
			} else {
				balance -= amount
				withdraws_status <- false
			}
		case amount := <-deposits:
			balance += amount
		case balances <- balance:

		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}
