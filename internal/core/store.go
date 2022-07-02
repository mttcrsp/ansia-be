package core

import (
	"context"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	createTablesSQL = `
	PRAGMA foreign_keys=ON;

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
		keywords STRING NOT NULL,
		content STRING NOT NULL,
		image_url STRING,
		FOREIGN KEY(item_id) REFERENCES item(item_id) ON DELETE CASCADE
	);
	`
	insertItemSQL = `
	INSERT INTO item (item_id, title, description, url, published_at, feed)
		VALUES (:item_id, :title, :description, :url, :published_at, :feed);
	`
	deleteItemSQL = `
	DELETE FROM item WHERE item_id = ?;
	`
	getItemsSql = `
	SELECT * FROM item;
	`
	insertArticleSQL = `
	INSERT INTO article (item_id, keywords, content, image_url)
		VALUES (:item_id, :keywords, :content, :image_url)
	`
	getArticlesSQL = `
	SELECT * FROM article;
	`
	getFeedSQL = `
	SELECT *
	FROM item i JOIN article a ON i.item_id = a.item_id
	WHERE i.feed = ?;
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

	database, err := sqlx.Connect("sqlite3", "file::memory:?cache=shared&mode=rwc")
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

	tx, err := db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) InsertItems(items []Item) error {
	return s.withTx(func(tx *sqlx.Tx) error {
		for _, item := range items {
			if _, err := tx.NamedExec(insertItemSQL, item); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) DeleteItems(items []Item) error {
	return s.withTx(func(tx *sqlx.Tx) error {
		for _, item := range items {
			if _, err := tx.Exec(deleteItemSQL, item.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) InsertArticle(article Article) error {
	return s.withTx(func(tx *sqlx.Tx) error {
		_, err := tx.NamedExec(insertArticleSQL, article)
		return err
	})
}

func (s *Store) GetItems() ([]Item, error) {
	db, err := s.db()
	if err != nil {
		return nil, err
	}

	items := []Item{}
	err = db.Select(&items, getItemsSql)
	return items, err
}

func (s *Store) GetArticles() ([]Article, error) {
	db, err := s.db()
	if err != nil {
		return nil, err
	}

	articles := []Article{}
	err = db.Select(&articles, getArticlesSQL)
	return articles, err
}

func (s *Store) GetFeed(feed string) ([]FeedItem, error) {
	db, err := s.db()
	if err != nil {
		return nil, err
	}

	items := []FeedItem{}
	err = db.Select(&items, getFeedSQL, feed)
	return items, err
}
