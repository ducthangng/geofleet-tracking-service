package singleton

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

type DBConnection struct {
	DB *pgxpool.Pool
}

var (
	PGconn *DBConnection
	PGonce sync.Once
	PGmu   sync.Mutex
)

func GetConn() *DBConnection {
	PGmu.Lock()
	defer PGmu.Unlock()

	return PGconn
}

func ConnectPostgre(ctx context.Context) {
	PGonce.Do(func() {
		for {
			PGconn = &DBConnection{}

			if err := PGconn.Connect(ctx); err != nil {
				time.Sleep(time.Duration(1 * time.Second))
				log.Println("retry ...")
			} else {

				log.Println("connect postgres successfully!")
				break
			}
		}
	})
}

func (db *DBConnection) Connect(ctx context.Context) error {
	dbCfg := GetGlobalConfig()

	connectionString := fmt.Sprintf(`postgres://%s:%s@%s:%s/%s`,
		dbCfg.User, dbCfg.Password, dbCfg.DatabaseHost, dbCfg.DatabasePort, dbCfg.Name)

	// 2. Add Pooling Parameters to the Query String
	// We use url.Values to safely encode these parameters
	q := url.Values{}
	q.Set("sslmode", "disable") // or "require" in production

	// Set Max Connections (pool_max_conns)
	if dbCfg.MaxOpenConns > 0 {
		q.Set("pool_max_conns", fmt.Sprintf("%d", dbCfg.MaxOpenConns))
	}

	// Set Max Connection Lifetime (pool_max_conn_lifetime)
	// Using .String() converts time.Duration (e.g. 5m0s) to string "5m0s" which pgx parses
	if dbCfg.MaxConnLifeTime > 0 {
		// convert seconds to time.Duration
		maxConnLifeTimeDuration := time.Duration(dbCfg.MaxConnLifeTime * int(time.Second))
		q.Set("pool_max_conn_lifetime", maxConnLifeTimeDuration.String())
	}

	// Set Max Connection Idle Time (pool_max_conn_idle_time)
	if dbCfg.MaxConnIdleTime > 0 {
		maxConnIdleTimeDuration := time.Duration(dbCfg.MaxConnIdleTime * int(time.Second))
		q.Set("pool_max_conn_idle_time", maxConnIdleTimeDuration.String())
	}

	// Append query params to DSN
	fullDSN := connectionString + "?" + q.Encode()

	pool, err := pgxpool.New(ctx, fullDSN)
	if err != nil {
		return (errors.New("error connect to pool with error: " + err.Error()))
	}

	db.DB = pool
	return nil
}
