package repository

import (
	"context"
	"database/sql"
	"fmt"
	"graves/internal/models"
	"graves/pkg/config"
	"graves/pkg/dto"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	mysqlDriver "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/payOSHQ/payos-lib-golang"
)

type Repository interface {
	CreateOrder(ctx context.Context, userId uint64, data payos.CheckoutResponseDataType) error
	UpdateOrderStatus(ctx context.Context, orderId int32, status string) error
	ListOrders(ctx context.Context, userId uint64, req dto.ListOrders) ([]models.Order, error)
	GetOrderById(ctx context.Context, userId uint64, orderId int32) (models.Order, error)
}

type repository struct {
	db      *sql.DB
	queries *models.Queries
}

var (
	repositoryInstance Repository
	once               sync.Once
	onceErr            error
)

func GetInstance() (Repository, error) {
	once.Do(func() {
		cfg, err := config.GetInstance()
		if err != nil || cfg.DataBase == nil {
			onceErr = err
			return
		}
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.DataBase.UserName,
			cfg.DataBase.Password,
			cfg.DataBase.Host,
			cfg.DataBase.Port,
			cfg.DataBase.Name,
		)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			onceErr = err
			return
		}

		driver, err := mysqlDriver.WithInstance(db, &mysqlDriver.Config{})
		if err != nil {
			onceErr = err
			return
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://sql/migrations",
			cfg.DataBase.Name, driver)
		if err != nil {
			onceErr = err
			return
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			onceErr = err
			return
		}

		repositoryInstance = &repository{
			db:      db,
			queries: models.New(db),
		}

	})
	return repositoryInstance, onceErr
}


