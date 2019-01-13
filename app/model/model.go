package model

import (
	// base core imports
	"database/sql"
	"errors"
	"time"

	// external imports
	mgo "gopkg.in/mgo.v2"
)

var (
	// ErrCode is a config or an internal error
	ErrCode = errors.New("Case statement in code is not correct.")
	// ErrNoResult is a not results error
	ErrNoResult = errors.New("Result not found.")
	// ErrUnavailable is a database not available error
	ErrUnavailable = errors.New("Database is unavailable.")
	// ErrUnauthorized is a permissions violation
	ErrUnauthorized = errors.New("User does not have permission to perform this operation.")
)

// standardizeErrors returns the same error regardless of the database used
func standardizeError(err error) error {
	if err == sql.ErrNoRows || err == mgo.ErrNotFound {
		return ErrNoResult
	}

	return err
}

// getDateTime return current
func getDateTime() string {
	now := time.Now().UTC()
	return now.Format("2006-01-02 15:04:05")
}

// prepareTime formats time output to custom timezone
func prepareTime(tt string) time.Time {
	loc, _ := time.LoadLocation("Europe/Budapest")
	fintime, _ := time.ParseInLocation("2006-01-02 15:04:05", tt, loc)
	return fintime
}
