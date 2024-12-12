package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	Conn *sql.DB
	mu   sync.Mutex
}

type Domain struct {
	Domain      string
	Address     string
	Redirect_to string
}

func NewSQLite(database string) (*SQLite, error) {
	// соединение с sqlite3
	conn, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, fmt.Errorf("failing to connect to sqlite3: %w", err)
	}

	// проверка соединения
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failing to check connection to sqlite3: %w", err)
	}

	return &SQLite{Conn: conn}, nil
}

// закрытие соединения
func (db *SQLite) Close() {
	db.Conn.Close()
}

// добавление домена в бд
func (db *SQLite) AddDomain(domain *Domain) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.Conn.Exec("INSERT INTO domains(domain, address, redirect_to) VALUES (?, ?, ?)", domain.Domain, domain.Address, domain.Redirect_to)
	if err != nil {
		return fmt.Errorf("failing to insert data to sqlite3: %w", err)
	}

	return nil
}

// получение адресов домена
func (db *SQLite) GetAddress(domain string) ([]string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.Conn.Query("SELECT address FROM domains WHERE domain = ? AND address IS NOT NULL", domain)
	if err != nil {
		return nil, fmt.Errorf("failing to query data from sqlite3: %w", err)
	}
	defer rows.Close()

	var addresses []string
	for rows.Next() {
		var address string
		if err := rows.Scan(&address); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return addresses, nil
}

// получение домена для переадресации запроса
func (db *SQLite) GetRedirectAddress(domain string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var redirectAddress string
	err := db.Conn.QueryRow("SELECT redirect_to FROM domains WHERE domain = ? AND redirect_to IS NOT NULL", domain).Scan(&redirectAddress)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("domain not found: %w", err)
		}
		return "", fmt.Errorf("failing to query data from sqlite3: %w", err)
	}

	return redirectAddress, nil
}

// проверка наличия домена в бд
func (db *SQLite) DomainExists(domain string) (bool, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var exists bool
	var redirectTo sql.NullString
	err := db.Conn.QueryRow("SELECT 1, redirect_to FROM domains WHERE domain = ?", domain).Scan(&exists, &redirectTo)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("failed to check domain: %v", err)
		}
		return false, false
	}

	return true, redirectTo.Valid
}
