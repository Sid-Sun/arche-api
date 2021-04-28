package initializers

import (
	_ "github.com/denisenkom/go-mssqldb"

	"context"
	"database/sql"
	"fmt"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
)

func InitDBClient(cfg *config.DBConfig, lgr *zap.Logger) (*sql.DB, error) {
	var err error

	// Create connection pool
	db, err := sql.Open("sqlserver", cfg.GetConn())
	if err != nil {
		lgr.Fatal(fmt.Sprintf("[Initializers] [InitDBClient] [Open] %s", err.Error()))
		return nil, err
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		lgr.Fatal(fmt.Sprintf("[Initializers] [InitDBClient] [PingContext] %s", err.Error()))
		return nil, err
	}

	return db, nil
}
