package accounts

import (
	"errors"
	"fmt"
)

// Account struct
type Account struct {
	owner   string
	balance int
}

// NewAccount is created
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

// Deposit 'ammount' into account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

// Balance of your account
func (a Account) Balance() int {
	return a.balance
}

var errNoMoney = errors.New("You don't got enough money")

// Withdraw 'amount' from account
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil
}

// ChangeOwner to new person
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// Owner is returned
func (a *Account) Owner() string {
	return a.owner
}

func (a *Account) String() string {
	return fmt.Sprint(a.Owner(), "'s account. \nHas: ", a.Balance())
}
