package v1

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"jonathanschmitt.com.br/go-api/internal/config"
)

func Utils_StringContainsSplit(val string, validValues string, splitChar string) bool {
	if len(validValues) > 0 {
		s := strings.Split(validValues, splitChar)

		for _, item := range s {
			if item == val {
				return true
			}
		}
	}
	return false
}

func Utils_StringSplitInt(items string, splitChar string) []int {
	result := []int{}
	if len(items) > 0 {
		s := strings.Split(items, splitChar)

		for _, item := range s {
			if val, err := strconv.Atoi(item); err == nil {
				result = append(result, val)
			}
		}
	}
	return result
}

func ConnectDatabase(config config.ServiceConfiguration) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.Database.ConnectionString)
	if err != nil {
		return db, err
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(config.Database.DatabaseMinutesIdle))
	db.SetMaxIdleConns(config.Database.DatabaseMaxIdleConns)
	db.SetMaxOpenConns(config.Database.DatabaseMaxOpenConns)

	var count int
	err = db.QueryRow("SELECT 1 AS count;").Scan(&count)
	if err != nil {
		return db, err
	}
	if count != 1 {
		return nil, errors.New("failed to return test database")
	}
	return db, err
}
