package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/viper"

	"github.com/cenkalti/backoff"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/sqlite3"

	"github.com/techievee/xero/xeroLog/debugcore"
)

const (
	createProductTable = `CREATE TABLE IF NOT EXISTS "Products" (
	"Id"	varchar(36) DEFAULT NULL,
	"Name"	varchar(17) DEFAULT NULL,
	"Description"	varchar(35) DEFAULT NULL,
	"Price"	decimal(6 , 2) DEFAULT NULL,
	"DeliveryPrice"	decimal(4 , 2) DEFAULT NULL,
	PRIMARY KEY("Id")
	)`

	createProductOptionsTable = `CREATE TABLE IF NOT EXISTS  "ProductOptions" (
	"Id"	varchar(36) DEFAULT NULL,
	"ProductId"	varchar(36) DEFAULT NULL,
	"Name"	varchar(9) DEFAULT NULL,
	"Description"	varchar(23) DEFAULT NULL,
	PRIMARY KEY("Id"),
	FOREIGN KEY("ProductId") REFERENCES "Products"("Id") ON DELETE CASCADE
	)`

	createProductIndex = `CREATE INDEX IF NOT EXISTS "product_id_index" ON "Products" (
	"Name"	ASC
	)`
)

type DB struct {
	// Read write connection with one DB connection open always
	RW func(ctx context.Context, label ...string) *sql.DB
	// ReadOnly connetion where multiple connection can be existing simultaneously
	RO       func(ctx context.Context) *sql.DB
	dbConfig *viper.Viper
	Logger   debugcore.Logger
}

type DBCfg struct {
	Driver   string
	FilePath string
	Database string

	Options map[string]interface{}

	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
}

type dbConn struct {
	config DBCfg   //holds config info
	db     *sql.DB //holds connection pool
}

func NewDB(config *viper.Viper, configFile string, logger debugcore.Logger) *DB {

	dbConfig := config.Sub(configFile)

	rw := InitDB(dbConfig)
	ro := func(ctx context.Context) *sql.DB {
		return rw(ctx, "readonly-db")
	}

	// Initialize the value with the default table initizlize structure
	ctx := context.Background()
	rwDB := rw(ctx)
	rwDB.Exec(createProductTable)
	rwDB.Exec(createProductIndex)
	rwDB.Exec(createProductOptionsTable)

	return &DB{
		RW:       rw,
		RO:       ro,
		dbConfig: dbConfig,
		Logger:   logger,
	}
}

var dbConnections = make(map[string]*dbConn)
var dbMutex sync.RWMutex

// labelOverride is useful for testing purposes.
func InitDB(dbConfig *viper.Viper, labelOverride ...string) func(ctx context.Context, label ...string) *sql.DB {

	glog.Info("InitDB")

	defaultLabel := dbConfig.GetString("default")
	if strings.TrimSpace(defaultLabel) == "" {
		log.Fatalf("A default database connection label is required")
	}

	allConfig := dbConfig.AllSettings()
	for key := range allConfig {
		if key == "default" {
			continue
		}

		var cfg DBCfg
		err := dbConfig.UnmarshalKey(key, &cfg)
		if err != nil {
			log.Fatalf("unable to decode database configuration, %v", err)
		}

		dbMutex.Lock()
		dbConnections[key] = &dbConn{
			config: cfg,
		}
		dbMutex.Unlock()
	}

	return func(ctx context.Context, _label ...string) *sql.DB {

		var label string
		if len(_label) == 0 {
			label = defaultLabel
		} else {
			label = _label[0]
		}

		if len(labelOverride) != 0 {
			label = labelOverride[0]
		}

		var db *sql.DB

		// Create a retryable operation
		operation := func() error {

			justCreated := false // Used to indicate if connection pool was registered for first time

			//Check if database connection pool already in dbConnections Map
			dbMutex.RLock()
			cfg := dbConnections[label].config
			pool := dbConnections[label].db
			dbMutex.RUnlock()

			if pool == nil {

				//Connection pool does not exist - make a new one
				var err error
				pool, err = apmsql.Open(cfg.Driver, cfg.ConnectionOpenString())
				if err != nil {
					return err //Retry attempt
				}

				// Create standard DN structures if it doesn't exist
				/*pool.Exec(createProductTable)
				pool.Exec(createProductIndex)
				pool.Exec(createProductOptionsTable)*/

				pool.SetConnMaxLifetime(cfg.ConnMaxLifetime * time.Minute)
				pool.SetMaxIdleConns(cfg.MaxIdleConns)
				pool.SetMaxOpenConns(cfg.MaxOpenConns)
				justCreated = true
			}

			// Call PingContext here and test for driver.ErrBadConn. If driver.ErrBadConn, don't try again. Otherwise try again.
			err := pool.PingContext(ctx)
			if err != nil {
				return err //Retry
			}

			if justCreated {
				//Store into dbConnections Map so we can reuse again
				dbMutex.Lock()
				dbConnections[label].db = pool
				dbMutex.Unlock()
			}

			db = pool
			return nil //All good
		}

		backoffAlgorithm := backoff.NewExponentialBackOff()
		backoffAlgorithm.MaxElapsedTime = time.Duration(10000) * time.Millisecond
		err := backoff.Retry(operation, backoffAlgorithm)
		if err != nil {
			glog.Errorf("DBError: %v", err)
			panic(err)
		}

		return db
	}

}

func (d *DBCfg) ConnectionOpenString() string {
	var opts string

	if len(d.Options) > 0 {
		opts = "?"
		var vs []string
		for key, value := range d.Options {
			switch v := value.(type) {
			case bool:
				if v {
					value = "true"
				} else {
					value = "false"
				}
			}
			vs = append(vs, key+"="+value.(string))
		}
		opts += strings.Join(vs, "&")
	}

	v := []interface{}{
		d.FilePath,
		d.Database,
	}
	str := fmt.Sprintf("file:%s%s.db", v...)

	conns := str + opts

	return conns
}
