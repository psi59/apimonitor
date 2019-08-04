package rsdb

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	lalaerrors "github.com/lalaworks/gopkg/errors"
	"github.com/pkg/errors"
)

const ErrInvalidTransaction = lalaerrors.Error("invalid transaction")

var (
	conn *defaultConnection
)

type Connection interface {
	Begin() (Transaction, error)
	Close() error
}

type defaultConnection struct {
	db *gorm.DB
}

func (conn *defaultConnection) Begin() (Transaction, error) {
	if conn.db == nil {
		return nil, errors.New("connection is nil")
	}
	tx := conn.db.Begin().Set("gorm:table_options", "ENGINE=InnoDB charset=utf8mb4")
	return &defaultTransaction{tx: tx}, nil
}

func (conn *defaultConnection) Close() error {
	if conn.db == nil {
		return errors.New("connection is nil")
	}
	return conn.db.Close()
}

func GetConnection() Connection {
	return conn
}

type Transaction interface {
	Rollback() error
	Commit() error
	Tx() *gorm.DB
}

type defaultTransaction struct {
	tx *gorm.DB
}

func (t *defaultTransaction) Rollback() error {
	return t.Tx().Rollback().Error
}

func (t *defaultTransaction) Commit() error {
	return t.Tx().Commit().Error
}

func (t *defaultTransaction) Tx() *gorm.DB {
	return t.tx
}

func NewTransaction(tx *gorm.DB) Transaction {
	return &defaultTransaction{tx}
}

func Init(config DatabaseConfig) error {
	dbAccessInfo := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true&sql_mode=STRICT_ALL_TABLES",
		config.GetUsername(),
		config.GetPassword(),
		config.GetHost(),
		config.GetPort(),
		config.GetDatabaseName(),
	)
	db, err := gorm.Open("mysql", dbAccessInfo)
	if err != nil {
		db, err = gorm.Open("mysql", fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/?charset=utf8mb4&parseTime=true&sql_mode=STRICT_ALL_TABLES",
			config.GetUsername(),
			config.GetPassword(),
			config.GetHost(),
			config.GetPort(),
		))
		if err != nil {
			return errors.WithStack(err)
		}

		db = db.Set("gorm:table_options", "ENGINE=InnoDB charset=utf8mb4")
		_, err = db.DB().Exec("CREATE DATABASE " + config.GetDatabaseName())
		if err != nil {
			return errors.WithStack(err)
		}

		if err = db.Close(); err != nil {
			return errors.WithStack(err)
		}

		db, err = gorm.Open("mysql", dbAccessInfo)
		if err != nil {
			return errors.WithStack(err)
		}

		if err = db.DB().Ping(); err != nil {
			return errors.WithStack(err)
		}
	}
	db = db.Set("gorm:table_options", "ENGINE=InnoDB charset=utf8mb4")

	db.DB().SetMaxOpenConns(80)
	db.DB().SetMaxIdleConns(0)
	db.DB().SetConnMaxLifetime(20 * time.Second)

	if config.GetVerbose() {
		db = db.Debug()
	}
	conn = &defaultConnection{db: db}
	return nil
}
