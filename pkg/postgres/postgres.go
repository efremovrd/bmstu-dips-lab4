package postgres

import (
	"bmstu-dips-lab2/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second

	PermDenied = "42501"
)

type Postgres struct {
	ConnAttempts uint
	ConnTimeout  time.Duration
	Builder      squirrel.StatementBuilderType
	Pool         pgxpoolmock.PgxPool
	//Pool         *pgxpool.Pool
}

func New(c *config.Config) (*Postgres, error) {
	pg := new(Postgres)

	pg.ConnAttempts = _defaultConnAttempts
	pg.ConnTimeout = _defaultConnTimeout

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlDbname,
		c.Postgres.PostgresqlPassword,
	)

	var err error

	for pg.ConnAttempts > 0 {
		pool, err := pgxpool.Connect(context.Background(), dataSourceName)
		if err == nil {
			pg.Pool = pool
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.ConnAttempts)

		time.Sleep(pg.ConnTimeout)

		pg.ConnAttempts--
	}

	if err != nil {
		return nil, err
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
