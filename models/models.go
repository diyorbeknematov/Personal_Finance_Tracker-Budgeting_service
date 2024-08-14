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

type Notification struct {
	ID        string    `bson:"_id"`
	Type      string    `bson:"type"`
	Message   string    `bson:"message"`
	Status    string    `bson:"status"`
	IsRead    bool      `bson:"is_read"`
	CreatedAt time.Time `bson:"created_at"`
}

type BudgetPerformance struct {
	CategoryId string  `json:"category_id,omitempty"`
	Target     float64 `json:"target,omitempty"`
	Actual     float64 `json:"actual,omitempty"`
	Progress   float64 `json:"progress,omitempty"`
}

type GoalProgress struct {
	Id            string  `json:"id,omitempty"`
	Name          string  `json:"name,omitempty"`
	TargetAmount  float64 `json:"target_amount,omitempty"`
	CurrentAmount float64 `json:"current_amount,omitempty"`
	Progress      float64 `json:"progress,omitempty"`
}
