package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_"github.com/pelletier/go-toml/query"
	_ "github.com/pelletier/go-toml/query"
)

type storage interface {
	CreateAccount(*Account) error
	DeleAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*DispAccount, error)
	GetAccountById(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
	MakeTransaction(int, int, int) error
	CreateTransaction(*transaction) error
	CheckAccountBalance(int) (int, error)
	GetTransaction(int) ([]*transaction, error)
	GetAccountByUserName(string) (*Account, error)
	CreateRequest(*Request) error
	Getrequests(string) ([]*Request, error)
	Createfriends(*Friends) error
	GetRequestbyId(string) (*Request,error)
	DeleteRequestbyId(string) error
	MakeFriendsGraph()(error,*Graph)
}
type PostgressStore struct {
	db *sql.DB

}


func NewPostgressStore() (*PostgressStore, error) {
	conStr := "user=postgres dbname=postgres host=localhost port=5432 password=mysecretpassword sslmode=disable"

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgressStore{
		db: db,
	}, nil
}
func (s *PostgressStore) init() (error, error, error, error) {
	return s.CreateAccountTable(), s.CreateTransactionTable(), s.CreateRequestTable(), s.CreateFriendsTable()

}
func (s *PostgressStore) CreateAccountTable() error {
	query := `create table if not exists bankaccount ( 
	id serial primary key, 
	first_name varchar(50),  
	last_name varchar(50),
	user_name varchar(50),
	pin_code varchar(150), 
	email varchar(50),
	number serial, 
	balance serial,
	createad_at timestamp 
	)`
	_, err := s.db.Exec(query)
	return err
}
func (s *PostgressStore) CreateRequestTable() error {
	query := `create table if not exists requests (
	id varchar(200) primary key,
	sender varchar(50),
	reciever varchar(50),
	time timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}
func (s *PostgressStore) CreateTransactionTable() error {
	query := `create table if not exists transaction(
		id varchar(200) primary key,
		account serial,
		amount serial,
		to_account serial,
		made_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}
func (s *PostgressStore) CreateFriendsTable() error {
	query := `create table if not exists friend(
	id serial primary key,
	member1 varchar(50),
	member2 varchar(50),
	made_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}
func (s *PostgressStore) MakeFriendsGraph() (error,*Graph){
   rows,err:=s.db.Query("Select member1,member2 from friend")
   if err!=nil{
	return err,nil
   }
   graph := NewGraph()
   for rows.Next() {
	var member1, member2 string
	if err := rows.Scan(&member1, &member2); err != nil {
		panic(err)
	}
	graph.AddEdge(member1, member2)
   }
   for user, friends := range graph.AdjacencyList {
	fmt.Printf("%s: %v\n", user, friends)
   }
   return nil,graph
}
func (s *PostgressStore) CreateAccount(ac *Account) error {
	query := `insert into bankaccount
	(first_name, last_name, number, user_name, balance, createad_at, pin_code, email)
	values ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.Query(
		query,
		ac.FirstName,
		ac.LastName,
		ac.Number,
		ac.UserName,
		ac.Balance,
		ac.CreatedAt,
		ac.PinCode,
		ac.Email,
	)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostgressStore) CreateTransaction(tr *transaction) error {
	query := `insert into transaction
  (id, account, amount, to_account, made_at )
  values($1, $2, $3, $4, $5)`
	resp, err := s.db.Query(
		query,
		tr.ID,
		tr.Account,
		tr.Amount,
		tr.ToAccount,
		tr.MadeAt,
	)
	if err != nil {
		return err
	}
	fmt.Printf("%+v/n", resp)
	return nil
}
func (s *PostgressStore) CreateRequest(rq *Request) error {
	query := `insert into requests 
	(id, sender, reciever, time)
	values($1,$2,$3,$4)`
	resp, err := s.db.Query(
		query,
		rq.ID,
		rq.Sender,
		rq.Reciever,
		rq.Time,
	)
	if err != nil {
		return err
	}
	fmt.Printf("%+v/n", resp)
	return nil
}

func (s *PostgressStore) Createfriends(fr *Friends) error{
	query := `insert into friend 
	(member1,member2,made_at)
	values($1,$2,$3)`
	resp, err := s.db.Query(
		query,
		fr.Member1,
		fr.Member2,
		fr.Time,
	)
	if err!=nil{
		return err
	}
	fmt.Printf("%+v/n", resp)
	return nil
}

func (s *PostgressStore) DeleAccount(id int) error {
	_, err := s.db.Query("delete from bankaccount where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostgressStore) UpdateAccount(ac *Account) error {
	return nil
}
func (s *PostgressStore) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query("select *from bankaccount where id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanintoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}
func (s *PostgressStore) GetAccountByNumber(accn int) (*Account, error) {
	rows, err := s.db.Query("select *from bankaccount where number = $1", accn)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanintoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", accn)
}
func (s *PostgressStore) GetAccounts() ([]*DispAccount, error) {
	rows, err := s.db.Query("select id, first_name, last_name, number, createad_at, user_name from bankaccount")
	if err != nil {
		return nil, err
	}
	accounts := []*DispAccount{}
	for rows.Next() {
		account := new(DispAccount)
		if err = rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.CreatedAt, &account.UserName); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgressStore) Getrequests(user string) ([]*Request, error) {
	rows, err := s.db.Query("select * from requests where reciever = $1", user)
	if err != nil {
		return nil, err
	}
	requests := []*Request{}
	for rows.Next() {
		request := new(Request)
		if err = rows.Scan(&request.ID, &request.Sender, &request.Reciever, &request.Time); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func ScanintoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.UserName, &account.PinCode, &account.Email, &account.Number, &account.Balance, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, err
}
func (s *PostgressStore) MakeTransaction(reciever, sender, amount int) error {
	acc, err := s.GetAccountByNumber(reciever)
	if err != nil {
		return err
	}
	acc.Balance = acc.Balance + int64(amount)
	er := s.UpdateBankBalance(reciever, int(acc.Balance))
	if er != nil {
		return er
	}
	sacc, err := s.GetAccountByNumber(sender)
	if err != nil {
		return err
	}
	sacc.Balance = sacc.Balance - int64(amount)
	ers := s.UpdateBankBalance(sender, int(sacc.Balance))
	if ers != nil {
		return ers
	}
	return nil

}
func (s *PostgressStore) UpdateBankBalance(number, balance int) error {
	query := `UPDATE bankaccount set balance = $1 where number = $2`
	_, err := s.db.Query(
		query,
		balance,
		number,
	)
	if err != nil {
		return err
	}

	return nil
}
func (s *PostgressStore) CheckAccountBalance(accn int) (int, error) {
	acc, err := s.GetAccountByNumber(accn)
	if err != nil {
		return 0, err
	}
	avlbal := acc.Balance
	return int(avlbal), nil
}
func (s *PostgressStore) GetTransaction(number int) ([]*transaction, error) {
	rows, err := s.db.Query("select * from transaction where account = $1", number)
	if err != nil {
		return nil, err
	}
	transactions := []*transaction{}
	for rows.Next() {
		transaction := new(transaction)
		err := rows.Scan(&transaction.ID, &transaction.Account, &transaction.Amount, &transaction.ToAccount, &transaction.MadeAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *PostgressStore) GetAccountByUserName(username string) (*Account, error) {
	rows, err := s.db.Query("select * from bankaccount where user_name = $1", username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanintoAccount(rows)
	}
	return nil, fmt.Errorf("user not found %s", username)

}
func (s *PostgressStore) GetRequestbyId(id string) (*Request,error) {
    rows, err := s.db.Query("select * from requests where id = $1",id)
	if err!=nil{
		return nil,err
	}
	for rows.Next() {
		return ScanIntoRequest(rows)
	}
	return nil,fmt.Errorf("Connection Request Not found")
}
func ScanIntoRequest(rows *sql.Rows) (*Request,error){
	request := new(Request)
	err:= rows.Scan(&request.ID,&request.Sender,&request.Reciever,&request.Time)
	if err!=nil{
		return nil,err
	}
	return request,err
}
func (s *PostgressStore) DeleteRequestbyId(id string) error{
	_, err := s.db.Query("delete from requests where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
