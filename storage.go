package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "github.com/pelletier/go-toml/query"
)
type storage interface{
	CreateAccount(*Account) error
	DeleAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts()([]*DispAccount,error)
	GetAccountById(int) (*Account,error)
	GetAccountByNumber(int) (*Account,error)
    MakeTransaction(int ,int,int) error
	CreateTransaction(*transaction) error
	CheckAccountBalance(int) (int,error)
	GetTransaction(int)([]*transaction,error)
}
type PostgressStore struct{
	db *sql.DB
} 
func NewPostgressStore() (*PostgressStore, error){
	conStr:= "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err:= sql.Open("postgres",conStr)
	if err!=nil{
		return nil,err
	}
	if err:=db.Ping();err!=nil{
		return nil,err
	}
	return &PostgressStore{
		db: db,
	},nil
}
func (s *PostgressStore) init() (error,error){
 return s.CreateAccountTable(),s.CreateTransactionTable()
 
}
func (s *PostgressStore) CreateAccountTable() error {
 query:= `create table if not exists bankaccount ( 
	id serial primary key, 
	first_name varchar(50),  
	last_name varchar(50),
	pin_code varchar(150), 
	email varchar(50),
	number serial, 
	balance serial,
	createad_at timestamp 
	)`
	_,err := s.db.Exec(query)
	return err
}
func (s *PostgressStore) CreateTransactionTable() error{
	query:= `create table if not exists transaction(
		id varchar(200) primary key,
		account serial,
		amount serial,
		to_account serial,
		made_at timestamp
	)`
	_,err:= s.db.Exec(query)
	return err
}

func (s *PostgressStore) CreateAccount(ac *Account) error{
	query:= `insert into bankaccount
	(first_name, last_name, number, balance, createad_at, pin_code, email)
	values ($1, $2, $3, $4, $5, $6, $7)`
	resp, err := s.db.Query(
		query,
		ac.FirstName,
		ac.LastName,
        ac.Number,
		ac.Balance,
		ac.CreatedAt,
		ac.PinCode,
		ac.Email,
	)
	if  err!=nil {
		return err;
	}
	fmt.Printf("%+v/n",resp)
	return nil
}
func (s *PostgressStore) CreateTransaction(tr *transaction) error{
  query:=`insert into transaction
  (id, account, amount, to_account, made_at )
  values($1, $2, $3, $4, $5)`
  resp, err:= s.db.Query(
	query,
	tr.ID,
	tr.Account,
	tr.Amount,
	tr.ToAccount,
	tr.MadeAt,
  )	
  if  err!=nil {
	return err;
}
  fmt.Printf("%+v/n",resp)
  return nil
}
func(s *PostgressStore) DeleAccount(id int) error{
	_ ,err:= s.db.Query("delete *from bankaccount where id = $1",id)
	if err!= nil {
		return err
	}
	return nil
}
func (s *PostgressStore) UpdateAccount(ac *Account) error{
	return nil
}
func (s *PostgressStore) GetAccountById(id int) (*Account,error){
	rows,err:= s.db.Query("select *from bankaccount where id = $1",id)
	if err!=nil{
		return nil, err
	}
    for rows.Next(){
		return ScanintoAccount(rows)
	}
	return nil,fmt.Errorf("account %d not found",id)
}
func (s *PostgressStore) GetAccountByNumber( accn int) (*Account,error){
	rows,err:= s.db.Query("select *from bankaccount where number = $1",accn)
	if err!=nil{
		return nil, err
	}
    for rows.Next(){
		return ScanintoAccount(rows)
	}
	return nil,fmt.Errorf("account %d not found",accn)
}
func (s *PostgressStore) GetAccounts() ([]*DispAccount,error){
	rows,err:= s.db.Query("select id, first_name, last_name, number, createad_at from bankaccount")
	if err!=nil{
		return nil, err
	}
	accounts := []*DispAccount{

	}
	for rows.Next() {
      account:= new(DispAccount)
	  if err= rows.Scan(&account.ID,&account.FirstName,&account.LastName,&account.Number,&account.CreatedAt,);err!=nil{
		return nil,err
	  }
	  accounts= append(accounts, account)
	}

	return accounts,nil
}
func ScanintoAccount(rows *sql.Rows) (*Account, error){
	account:= new(Account)
	 err:= rows.Scan(&account.ID,&account.FirstName,&account.LastName,&account.PinCode,&account.Email,&account.Number,&account.Balance,&account.CreatedAt,)
	 if err!=nil{
		return nil,err
	  }
	  return account,err
}
func (s *PostgressStore) MakeTransaction(reciever ,sender, amount int) error{
	 acc,err:=s.GetAccountByNumber(reciever)
	 if err!=nil{
		return err
	 }
	 acc.Balance=acc.Balance+int64(amount)
	er:= s.UpdateBankBalance(reciever,int(acc.Balance))
	if er!=nil{
		return er
	}
	 sacc,err:=s.GetAccountByNumber(sender)
	 if err!=nil{
		return err
	 }
     sacc.Balance=sacc.Balance-int64(amount)
	 ers:=s.UpdateBankBalance(sender,int(sacc.Balance))
	 if ers!=nil{
		return ers
	 }
	 return nil

}
func (s *PostgressStore) UpdateBankBalance(number,balance int) error{
	query:=`UPDATE bankaccount set balance = $1 where number = $2`
	_,err:=s.db.Query(
		query,
		balance,
		number,
	)
	if err!=nil {
		return err
	}
 
	return nil
}
func (s *PostgressStore) CheckAccountBalance( accn int) (int ,error){
	acc,err:=s.GetAccountByNumber(accn)
	if err!=nil{
		return 0,err
	}
	avlbal:=acc.Balance
	return int(avlbal),nil
}
func(s *PostgressStore) GetTransaction(number int) ([]*transaction,error){
	rows,err:= s.db.Query("select *from transaction where account = $1",number)
	if err!=nil{
		return nil,err
	}
	transactions := []*transaction{

	}
	for rows.Next() {
		transaction:= new(transaction)
		err:=rows.Scan(&transaction.ID,&transaction.Account,&transaction.Amount,&transaction.ToAccount,&transaction.MadeAt)
		if err!=nil{
			return nil,err
		}
		transactions=append(transactions, transaction)
	}
	return transactions,nil   
}

