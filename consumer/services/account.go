package services

import (
	"consumer/repositories"
	"encoding/json"
	"events"
	"log"
	"reflect"
)

type EventHandler interface {
	Handle(topic string, evenBytes []byte)
}

type accountEventHandler struct {
	accountRepo repositories.AccountRepository
}

func NewAccountEventHandler(accountRepo repositories.AccountRepository) EventHandler {
	return accountEventHandler{accountRepo}
}

func (obj accountEventHandler) Handle(topic string, evenBytes []byte) {
	switch topic {
	case reflect.TypeOf(events.OpenAccountEvent{}).Name():
		event := &events.OpenAccountEvent{}

		err := json.Unmarshal(evenBytes, event)
		if err != nil {
			log.Println(err)
		}

		bankAccount := repositories.BankAccount{
			ID:            event.ID,
			AccountHolder: event.AcoountHolder,
			AccountType:   event.AccountType,
			Balance:       event.OpeningBalance,
		}

		err = obj.accountRepo.Save(bankAccount)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%#v", event)
	case reflect.TypeOf(events.DepositFundEvent{}).Name():
		event := &events.DepositFundEvent{}

		err := json.Unmarshal(evenBytes, event)
		if err != nil {
			log.Println(err)
			return
		}

		bankAccount, err := obj.accountRepo.FindByID(event.ID)
		if err != nil {
			log.Println(err)
			return
		}

		bankAccount.Balance += event.Amount

		err = obj.accountRepo.Save(bankAccount)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%#v", event)
	case reflect.TypeOf(events.WithdrawFundEvent{}).Name():
		event := &events.WithdrawFundEvent{}

		err := json.Unmarshal(evenBytes, event)
		if err != nil {
			log.Println(err)
			return
		}

		bankAccount, err := obj.accountRepo.FindByID(event.ID)
		if err != nil {
			log.Println(err)
			return
		}

		bankAccount.Balance -= event.Amount

		err = obj.accountRepo.Save(bankAccount)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%#v", event)
	case reflect.TypeOf(events.CloseAccountEvent{}).Name():
		event := &events.CloseAccountEvent{}

		err := json.Unmarshal(evenBytes, event)
		if err != nil {
			log.Println(err)
			return
		}

		err = obj.accountRepo.Delete(event.ID)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%#v", event)
	default:
		log.Println("no event handler")
	}

}
