package models

import "time"

type GetAccount struct {
	ID       string  `bson:"_id"`
	UserId   string  `bson:"user_id"`
	Name     string  `bson:"name"`
	Type     string  `bson:"type"`
	Balance  float64 `bson:"balance"`
	Currency string  `bson:"currency"`
}


type GetTransaction struct {
	Id          string    `bson:"_id"`
	AccountId   string    `bson:"account_id"`
	UserId      string    `bson:"user_id"`
	CategoryId  string    `bson:"category_id"`
	Type        string    `bson:"type"`
	Amount      float64   `bson:"amount"`
	Description string    `bson:"description"`
	Date        time.Time `bson:"date"`
}

type GetCategory struct {
	ID     string `bson:"_id"`
	UserId string `bson:"user_id"`
	Name   string `bson:"name"`
	Type   string `bson:"type"`
}

type GetBudget struct {
	ID         string    `bson:"_id"`
	UserId     string    `bson:"user_id"`
	CategoryId string    `bson:"category_id"`
	Amount     float64   `bson:"amount"`
	Period     string    `bson:"period"`
	StartDate  time.Time `bson:"start_date"`
	EndDate    time.Time `bson:"end_date"`
}

type GetGoal struct {
	ID            string    `bson:"_id"`
	UserId        string    `bson:"user_id"`
	Name          string    `bson:"name"`
	TargetAmount  float64   `bson:"target_amount"`
	CurrentAmount float64   `bson:"current_amount"`
	Deadline      time.Time `bson:"deadline"`
	Status        string    `bson:"status"`
}

type Balance struct {
	AccountId string  `bson:"account_id"`
	Balance   float64 `bson:"balance"`
}

type Transaction struct {
	AccountId   string    `bson:"account_id"`
	UserId      string    `bson:"user_id"`
	CategoryId  string    `bson:"category_id"`
	Type        string    `bson:"type"`
	Amount      float64   `bson:"amount"`
	Description string    `bson:"description"`
	Date        time.Time `bson:"date"`
}
