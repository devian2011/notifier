package storage

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
	"net/url"
	"notifications/internal/dto"
)

type SqliteStorage struct {
	dsn *url.URL
	sql *sql.DB
}

const SQLiteInit = `
CREATE TABLE IF NOT EXISTS templates (
    code TEXT PRIMARY KEY UNIQUE,
    scheme TEXT NOT NULL
)`

func NewSqliteStorage(dsn *url.URL) (*SqliteStorage, error) {
	db, dbErr := sql.Open("sqlite3", dsn.Path)
	if dbErr != nil {
		return nil, dbErr
	}
	_, initErr := db.Exec(SQLiteInit)
	if initErr != nil {
		return nil, initErr
	}

	return &SqliteStorage{
		dsn: dsn,
		sql: db,
	}, nil
}

func (s *SqliteStorage) List() ([]dto.MessageTmpl, error) {
	type record struct {
		code   string
		scheme []byte
	}

	rows, selectErr := s.sql.Query("SELECT code, scheme FROM templates")
	if selectErr != nil {
		return nil, selectErr
	}
	defer rows.Close()

	result := make([]dto.MessageTmpl, 0)
	for rows.Next() {
		r := record{}
		scanErr := rows.Scan(&r.code, &r.scheme)
		if scanErr != nil {
			continue
		}

		tmpl := dto.MessageTmpl{}
		unmarshalErr := yaml.Unmarshal(r.scheme, &tmpl)
		if unmarshalErr != nil {
			continue
		}

		result = append(result, tmpl)
	}

	return result, nil
}

func (s *SqliteStorage) Load(code string) (dto.MessageTmpl, error) {
	var record struct {
		code   string
		scheme []byte
	}

	row := s.sql.QueryRow("SELECT code, scheme FROM templates WHERE code=?", code)

	result := dto.MessageTmpl{}

	if scanErr := row.Scan(&record.code, &record.scheme); scanErr != nil {
		if errors.Is(scanErr, ErrMessageNotExists) {
			return result, ErrMessageNotExists
		}
		return result, scanErr
	}

	unmarshalErr := yaml.Unmarshal(record.scheme, &result)

	return result, unmarshalErr
}

func (s *SqliteStorage) Create(msg dto.MessageTmpl) error {
	schemeData, marshalErr := yaml.Marshal(msg)
	if marshalErr != nil {
		return marshalErr
	}

	_, err := s.sql.Exec("INSERT INTO templates (code, scheme) VALUES (?, ?)", msg.Code, schemeData)
	return err
}

func (s *SqliteStorage) Update(msg dto.MessageTmpl) error {
	schemeData, marshalErr := yaml.Marshal(msg)
	if marshalErr != nil {
		return marshalErr
	}

	_, err := s.sql.Exec("UPDATE templates SET scheme = ? WHERE code = ?", schemeData, msg.Code)
	return err
}

func (s *SqliteStorage) Rm(code string) error {
	_, err := s.sql.Exec("DELETE FROM templates WHERE code=?", code)
	return err
}
