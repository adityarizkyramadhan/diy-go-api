package env

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq" // pq driver calls for blank identifier
)

// Datastore struct stores common environment related items
type Datastore struct {
	MainDb *sql.DB
	LogDb  *sql.DB
}

// NewDatastore initializes the datastore struct
// NOTE: I have chosen to use the same database for logging as
// my "main" app database. I'd recommend having a separate db and
// would have a separate method to start that connection pool up and
// pass it, but since this is just an example....
func NewDatastore() (*Datastore, error) {
	mdb, err := newMainDB()
	if err != nil {
		log.Error().Err(err).Msg("Error returned from newMainDB")
		return nil, err
	}

	return &Datastore{MainDb: mdb, LogDb: mdb}, nil
}

// NewMainDB returns an open database handle of 0 or more underlying connections
func newMainDB() (*sql.DB, error) {

	// Get Database connection credentials from environment variables
	dbName := os.Getenv("PG_DBNAME_TEST")
	dbUser := os.Getenv("PG_USERNAME_TEST")
	dbPassword := os.Getenv("PG_PASSWORD_TEST")
	dbHost := os.Getenv("PG_HOST_TEST")
	dbPort, err := strconv.Atoi(os.Getenv("PG_PORT_TEST"))
	if err != nil {
		log.Error().Err(err).Msg("Unable to complete string to int conversion for dbPort")
		return nil, err
	}

	// Craft string for database connection
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open the postgres database using the postgres driver (pq)
	// func Open(driverName, dataSourceName string) (*DB, error)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Error().Err(err).Msg("Error returned from sql.Open")
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error().Err(err).Msg("Error returned from db.Ping")
		return nil, err
	}
	return db, nil
}

// Tx opens a database connection and starts a database transaction using the
// BeginTx method which allows for rollback of the transaction if the context
// is cancelled
func (ds Datastore) Tx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {

	// Calls the BeginTx method of the MainDb opened database
	tx, err := ds.MainDb.BeginTx(ctx, opts)
	if err != nil {
		log.Error().Err(err).Msg("Error returned from BeginTx")
		return nil, err
	}

	return tx, nil
}
