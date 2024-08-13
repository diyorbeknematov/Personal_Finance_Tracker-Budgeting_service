package redis

import (
	"budgeting-service/models"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type BalanceRepository interface {
}

type accountBalanceImpl struct {
	client *redis.Client
}

func NewAccountBalance(rdb *redis.Client) BalanceRepository {
	return &accountBalanceImpl{client: rdb}
}

func (repo *accountBalanceImpl) SetBalance(ctx context.Context, balance models.Balance) error {
	return repo.client.Set(ctx, "balance:"+balance.AccountId, balance.Balance, 10*time.Minute).Err()
}

func (repo *accountBalanceImpl) GetBalance(ctx context.Context, accountId string) (*models.Balance, error) {
	balance, err := repo.client.Get(ctx, "balance:"+accountId).Float64()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &models.Balance{
		AccountId: accountId,
		Balance:   balance,
	}, nil
}
