package database

import (
	"errors"
	"time"

	"github.com/JeremyCurmi/simpleBank/pkg/utils"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.Logger
	conn   *sqlx.DB
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Conn() *sqlx.DB {
	return c.conn
}

func New(logger *zap.Logger, url string, connMaxLifetime int) (*Client, error) {
	var (
		conn  *sqlx.DB
		err   error
		start = time.Now()
	)

	for {
		if time.Since(start) > time.Second*10 {
			return nil, errors.New(utils.ErrTimeout)
		}

		if conn == nil {
			conn, err = sqlx.Connect("postgres", url)
			if err != nil {
				logger.Warn(utils.WarnDBNotConnected, zap.Error(err))
				continue
			}
		}

		err = conn.Ping()
		if err == nil {
			logger.Info("connected to database ⚡️")
			break
		}

		logger.Warn(utils.WarnDBNotConnected)
		time.Sleep(time.Second * 2)
	}

	if connMaxLifetime > 0 {
		conn.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	}
	return &Client{logger: logger, conn: conn}, nil
}
