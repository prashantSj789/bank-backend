package main

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)
type CreateAccountRequest struct{
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	PinCode string `json:"pinCode"`
	Email   string  `json"email"`
}

type transactionRequest struct{
	Number int64`json:"accountNumber"`
	PinCode string `json:"pinCode"`
    Amount int `json:"amount"`
}
type transactionresponse struct{
	Number int64 `json:"accountNumber"`
	Amount int `json:"amount"`
	Message string `json:"message"`
}
type LoginRespone struct{
	Number int64 `json:"accountNumber"`
	Token string `json:"token"`
}
type LoginAccountRequest struct{
	Number int64 `json:"accountNumber"`
	PinCode string `json:"pinCode"`
}
type Account struct{
	ID int `json:"id"`
	FirstName string `json:"firstName"`
	LastName string  `json:"lastName"`
	PinCode string `son:"pinCode"`
	Number int64  `json:"accountNumber"`
	Email string  `json:"email"`
	Balance int64 `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}
type transaction struct {
	ID string `json:"transactionId"`
	Account int `json:"fromAccount"`
	Amount int `json:"amount"`
	ToAccount int `json:"toAccount"`
	MadeAt time.Time `json:"time"`
}
 func NewAccount(firstname, lastname, pincode, email string) (*Account,error){
	loc, _ := time.LoadLocation("Asia/Kolkata")
	pin,err:=bcrypt.GenerateFromPassword([]byte(pincode),bcrypt.DefaultCost)
	if err!=nil{
		return nil,err
	}
	return &Account{
		FirstName: firstname,
		LastName:  lastname,
		Number:    int64(rand.Intn(1000000000)),
		Balance:   0,
		CreatedAt: time.Now().In(loc),
		PinCode: string(pin),
		Email: email,
	},nil
 }
 func (a *Account) ValidatePincodw(pw string)bool{
	return bcrypt.CompareHashAndPassword([]byte(a.PinCode),[]byte(pw))==nil
 }
 type DispAccount struct{
	ID int `json:"id"`
	FirstName string `json:"firstName"`
	LastName string  `json:"lastName"`
	Number int64  `json:"accountNumber"`
	CreatedAt time.Time `json:"createdAt"`
 } 
 func NewTransaction( sender, reciever  ,amount int ) (*transaction ,error){
	loc, _ := time.LoadLocation("Asia/Kolkata")
	id:=string(uuid.New().String())
	return &transaction{
		ID: id,
		Account: sender,
		ToAccount: reciever,
		Amount: amount,
		MadeAt: time.Now().In(loc),
	},nil
 }