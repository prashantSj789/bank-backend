package main

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ConnectionRequest struct {
  UserName string `json:"userName"`
}

type Request struct {
	ID string `json:"id"`
	Sender string `json:"sender"`
	Reciever string `json:"reciever"`
	Time time.Time `json:"madeAt"`
}

type CreateAccountRequest struct{
	UserName string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	PinCode string `json:"pinCode"`
	Email   string  `json"email"`
}
type simpleQRCode struct {
    Content string
    Size        int
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
	UserName string `json:userName`
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

type Friends struct {
	Id string `json:"id"`
	Member1 string `json:"member1"`
	Member2 string `json:"member2'`
	Time time.Time `json:"madeAt"`
}

type Graph struct{
	AdjacencyList map[string][]string
}



 func NewAccount(firstname, lastname, pincode, username, email string) (*Account,error){
	loc, _ := time.LoadLocation("Asia/Kolkata")
	pin,err:=bcrypt.GenerateFromPassword([]byte(pincode),bcrypt.DefaultCost)
	if err!=nil{
		return nil,err
	}
	return &Account{
		FirstName: firstname,
		LastName:  lastname,
		UserName:   username,
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
	UserName string `userName`
	FirstName string `json:"firstName"`
	LastName string  `json:"lastName"`
	Number int64  `json:"accountNumber"`
	CreatedAt time.Time `json:"createdAt"`
 } 
 func NewRequest( sender, reciever string) (*Request, error){
	loc, _ :=time.LoadLocation("Asia/kolkata")
	id:=string(uuid.New().String())
	return &Request{
     ID: id,
	 Sender: sender,
	 Reciever: reciever,
	 Time: time.Now().In(loc),
	},nil
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

 func NewFriend(m1,m2 string) (*Friends,error){
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return &Friends{
		Member1: m1,
		Member2: m2,
        Time: time.Now().In(loc),
	},nil
 }
 func NewGraph() *Graph {
    return &Graph{AdjacencyList: make(map[string][]string)}
}
 func (g *Graph) AddEdge(user1, user2 string) {
    g.AdjacencyList[user1] = append(g.AdjacencyList[user1], user2)
    g.AdjacencyList[user2] = append(g.AdjacencyList[user2], user1)
}

func FindMutualFriends(graph *Graph, user1, user2 string) []string {
    mutuals := []string{}
    friends1 := graph.AdjacencyList[user1]
    friends2 := graph.AdjacencyList[user2]

    friendSet := make(map[string]bool)
    for _, friend := range friends1 {
        friendSet[friend] = true
    }
    for _, friend := range friends2 {
        if friendSet[friend] {
            mutuals = append(mutuals, friend)
        }
    }
    return mutuals
}
func BFS(graph *Graph, start, target string) int {
    visited := make(map[string]bool)
    queue := []string{start}
    distance := 0

    for len(queue) > 0 {
        size := len(queue)
        for i := 0; i < size; i++ {
            node := queue[0]
            queue = queue[1:]
            if node == target {
                return distance
            }
            if visited[node] {
                continue
            }
            visited[node] = true
            queue = append(queue, graph.AdjacencyList[node]...)
        }
        distance++
    }
    return -1 // Not connected
}
func SuggestFriends(graph *Graph, user string) map[string]int {
    suggestions := make(map[string]int)
    
    for friend := range graph.AdjacencyList {
        if friend == user {
            continue
        }
        distance := BFS(graph, user, friend)
        if  distance>1 && distance<=4 { // Suggest only friends-of-friends
            suggestions[friend] = distance
        }
    }

    return suggestions
}

