package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/cors"
	"github.com/gorilla/mux"
)



type apiFunc func(http.ResponseWriter, *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string
}

func makemuxhandlefunc(f apiFunc) http.HandlerFunc {
	
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	store      storage
}

func NewApiServer(listenAddr string, store storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makemuxhandlefunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makemuxhandlefunc(s.handleGetAccountbyId))
	router.HandleFunc("/login", makemuxhandlefunc(s.handleLoginAccount))
	router.HandleFunc("/transaction", makemuxhandlefunc(s.handleTransferAccount))
	router.HandleFunc("/transaction/history",makemuxhandlefunc(s.HandleGetTransactionHistory))
	router.HandleFunc("/balance",makemuxhandlefunc(s.HandleCheckBalance))

	log.Println("JSON Api Running on port:", s.listenAddr)

    // Setup CORS middleware options
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"}, // Allow specific origin
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
    })

    // Use the CORS middleware with the router
    handler := c.Handler(router)
  
	http.ListenAndServe(s.listenAddr,handler)

}
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}
func (s *APIServer) handleGetAccountbyId(w http.ResponseWriter, r *http.Request) error {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return fmt.Errorf("wrong id entered:%s", idstr)
	}
	
	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	CreateAccountRequest := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(CreateAccountRequest); err != nil {
		return err
	}
	account, err := NewAccount(CreateAccountRequest.FirstName, CreateAccountRequest.LastName, CreateAccountRequest.PinCode, CreateAccountRequest.Email)
	if err != nil {
		return err
	}
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)

}

func (s *APIServer) handleLoginAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed :%s", r.Method)
	}
	req := new(LoginAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	acc, err := s.store.GetAccountByNumber(int(req.Number))
	if err != nil {
		return fmt.Errorf("failed to login")
	}
	if !acc.ValidatePincodw(req.PinCode) {
		return fmt.Errorf("Not Authenticated")
	}
	token, err := CreateJWT(acc)
	if err != nil {
		return err
	}
	resp := LoginRespone{
		Number: acc.Number,
		Token:  token,
	}

	return WriteJSON(w, http.StatusOK, resp)
}
func (s *APIServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed :%s", r.Method)
	}

	req:=new(transactionRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
    err,aac:=validateToken(w,r,req.PinCode)
	if err!=nil{
		return err
	}
	acc,err:=s.store.GetAccountByNumber(aac)
	if err != nil {
		return err
	}
	if acc.ValidatePincodw(req.PinCode)==false{
		return fmt.Errorf("wrong pin")
	}
	if(acc.Balance<int64(req.Amount)){
		return fmt.Errorf("Insufficient Balance")
	}
	rev,err:=s.store.GetAccountByNumber(int(req.Number))
	if err!=nil{
      return err
	}
	println(rev.Number)
    transaction,err:= NewTransaction(aac,int(req.Number),req.Amount)
	if err != nil {
		return err
	}
	er:=s.store.MakeTransaction(int(rev.Number),int(acc.Number),req.Amount)
	if er!=nil{
		return er
	}
	err= s.store.CreateTransaction(transaction)
    if err!=nil{
		return err
	}
	return WriteJSON(w, http.StatusOK, transaction)	
	}
func (s *APIServer) HandleGetTransactionHistory(w http.ResponseWriter,r *http.Request) error{
	if r.Method !="GET"{
		return fmt.Errorf("Method not allowed:%s",r.Method)
	}
	pin:=r.Header.Get("PinCode")
	if pin == ""{
		return fmt.Errorf("No Pin Code Entered")
	}
	
	err,acc:=validateToken(w,r,pin)
	if err!=nil {
		return err
	}
	ac,err:=s.store.GetAccountByNumber(acc)
	if err!=nil{
		return err
	}
	if ac.ValidatePincodw(string(pin))==false{
		return fmt.Errorf("wrong pin")
	}
	trns,err:=s.store.GetTransaction(acc)
	if err!=nil{
		return err
	}
	return WriteJSON(w,http.StatusOK,trns)
}
func (s *APIServer) HandleCheckBalance(w http.ResponseWriter,r *http.Request) error{
	if r.Method !="GET"{
		return fmt.Errorf("Method not allowed:%s",r.Method)
	}
	pin:=r.Header.Get("PinCode")
	if pin == ""{
		return fmt.Errorf("No Pin Code Entered")
	}
	err,acc:=validateToken(w,r,pin)
	if err!=nil {
		return err
	}
	account,err:=s.store.GetAccountByNumber(acc)
	if err!=nil{
		return err
	}
	if account.ValidatePincodw(string(pin))==false{
		return fmt.Errorf("wrong pin")
	}
	bal,err:=s.store.CheckAccountBalance(acc)
	if err!=nil{
		return err
	}
   return WriteJSON(w,http.StatusOK,bal)
}
func CreateJWT(account *Account) (string, error) {
	claims := jwt.MapClaims{
		"expiresAt":     jwt.NewNumericDate(time.Now().Local().Add(time.Minute * 2)),
		"accountNumber": account.Number,
	}
	secret := os.Getenv("SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
func validateToken(w http.ResponseWriter, r *http.Request,pin string) (error,int) {

	if r.Header["Token"] == nil {
		fmt.Fprintf(w, "can not find token in header")
		return fmt.Errorf( "can not find token in header %s",w),0
	}

	token,_  := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing:",)
		}
		return os.Getenv("SECRET"), nil
	})


	if token == nil {
		fmt.Fprintf(w, "invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Fprintf(w, "couldn't parse claims")
		return errors.New("Token error"),0
	}

	exp := claims["expiresAt"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		fmt.Fprintf(w, "token expired")
		return errors.New("Token error"),0
	}
	accn:= claims["accountNumber"].(float64)
	return nil,int(accn)
}
