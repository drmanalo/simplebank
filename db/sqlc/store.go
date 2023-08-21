package db

import "github.com/jackc/pgx/v5/pgxpool"

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
}

// SQLStore provides all functions to execute SQL queries and transactions
type SLQStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SLQStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
