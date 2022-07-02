package core

import (
	"context"
	"sync"

	"github.com/jmoiron/sqlx"
)

const (
	createTablesSQL = `
	CREATE TABLE IF NOT EXISTS item (
		item_id INTEGER NOT NULL PRIMARY KEY,
		title STRING NOT NULL,
		description STRING NOT NULL,
		url STRING NOT NULL,
		published_at DATETIME NOT NULL,
		feed STRING NOT NULL
	);

	CREATE TABLE IF NOT EXISTS article (
		item_id INTEGER NOT NULL PRIMARY KEY,
		title STRING NOT NULL,
		description STRING NOT NULL,
		keywords STRING NOT NULL,
		content STRING NOT NULL,
		image_url STRING,
		FOREIGN KEY(item_id) REFERENCES item(item_id) ON DELETE CASCADE
	);
	`
	insertItemSQL = `
	INSERT INTO item (item_id, title, description, url, published_at)
		VALUES (:item_id, :title, :description, :url, :published_at);
	`
	deleteItemSQL = `
	DELETE FROM item WHERE item_id = :item_id;
	`
)

type Store struct {
	database     *sqlx.DB
	databaseLock sync.Mutex
}

func (s *Store) db() (*sqlx.DB, error) {
	s.databaseLock.Lock()
	defer s.databaseLock.Unlock()

	if s.database != nil {
		return s.database, nil
	}

	database, err := sqlx.Connect("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}

	if _, err = database.Exec(createTablesSQL); err != nil {
		return nil, err
	}

	s.database = database
	return database, nil
}

func (s *Store) withTx(fn func(*sqlx.Tx) error) error {
	db, err := s.db()
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err == nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) Insert(items []Item) error {
	return s.withTx(func(tx *sqlx.Tx) error {
		for _, item := range items {
			if _, err := tx.NamedExec(insertItemSQL, item); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Delete(items []Item) error {
	return s.withTx(func(tx *sqlx.Tx) error {
		for _, item := range items {
			if _, err := tx.NamedExec(deleteItemSQL, item.ID); err != nil {
				return err
			}
		}
		return nil
	})
}
