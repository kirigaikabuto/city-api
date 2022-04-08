package common

import (
	"database/sql"
	"strconv"
)

type PostgresConfig struct {
	Host             string
	Port             int
	User             string
	Password         string
	Database         string
	Params           string
	ConnectionString string
}

func GetConnString(cfg PostgresConfig) string {
	var connStr string
	if cfg.ConnectionString == "" {
		connStr = "postgres://"
		if cfg.Host == "" {
			cfg.Host = "localhost"
		}
		if cfg.Port == 0 {
			cfg.Port = 5432
		}
		connStr +=
			cfg.User + ":" +
				cfg.Password + "@" +
				cfg.Host + ":" +
				strconv.Itoa(cfg.Port) + "/" +
				cfg.Database + "?" +
				cfg.Params
	} else {
		connStr = cfg.ConnectionString
	}
	return connStr
}

func GetDbConn(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
