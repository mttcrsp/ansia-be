package core

import (
	"context"

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
		published_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS article (
		item_id INTEGER NOT NULL PRIMARY KEY,
		keywords STRING NOT NULL,
		content STRING NOT NULL,
		image_url STRING,
		FOREIGN KEY(item_id) REFERENCES item(item_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS item_feed (
		item_id INTEGER NOT NULL PRIMARY KEY,
		feed STRING NOT NULL,
		FOREIGN KEY(item_id) REFERENCES item(item_id) ON DELETE CASCADE
	);
	`
	insertItemSQL = `
	INSERT INTO item (item_id, title, description, url, published_at)
		VALUES (:item_id, :title, :description, :url, :published_at);
	`
	deleteItemSQL = `
	DELETE FROM item WHERE item_id = ?;
	`
	getItemsSql = `
	SELECT * 
	FROM item i JOIN item_feed if ON i.item_id = if.item_id;
	`
	insertItemFeedSQL = `
	INSERT INTO item_feed (item_id, feed)
		VALUES (:item_id, :feed);
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
	FROM item i 
		JOIN article a ON i.item_id = a.item_id
		JOIN item_feed if ON i.item_id = if.item_id
	WHERE if.feed = ?
	ORDER BY i.published_at;
	`
)

type Store struct {
	db *sqlx.DB
}

func (s *Store) withDB(fn func(db *sqlx.DB) error) error {
	if s.db != nil {
		return fn(s.db)
	}

	db, err := sqlx.Connect("sqlite3", "file::memory:?cache=shared&mode=rwc")
	// db, err := sqlx.Connect("sqlite3", "../../db.sqlite")
	if err != nil {
		return err
	}

	if _, err = db.Exec(createTablesSQL); err != nil {
		return err
	}

	s.db = db
	return fn(s.db)
}

func (s *Store) InsertItems(items []Item) error {
	return s.withDB(func(db *sqlx.DB) error {
		tx, err := db.BeginTxx(context.Background(), nil)
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		if err != nil {
			return err
		}

		for _, item := range items {
			if _, err := tx.NamedExec(insertItemSQL, item); err != nil {
				return err
			}

			if _, err := tx.NamedExec(insertItemFeedSQL, item); err != nil {
				return err
			}
		}

		return tx.Commit()
	})
}

func (s *Store) DeleteItems(items []Item) error {
	return s.withDB(func(db *sqlx.DB) error {
		for _, item := range items {
			if _, err := db.Exec(deleteItemSQL, item.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) InsertArticle(article Article) error {
	return s.withDB(func(db *sqlx.DB) error {
		_, err := db.NamedExec(insertArticleSQL, article)
		return err
	})
}

func (s *Store) GetItems() ([]Item, error) {
	items := []Item{}
	err := s.withDB(func(db *sqlx.DB) error {
		return db.Select(&items, getItemsSql)
	})
	return items, err
}

func (s *Store) GetArticles() ([]Article, error) {
	articles := []Article{}
	err := s.withDB(func(db *sqlx.DB) error {
		return db.Select(&articles, getArticlesSQL)
	})
	return articles, err
}

func (s *Store) GetFeed(feed string) ([]FeedItem, error) {
	items := []FeedItem{}
	err := s.withDB(func(db *sqlx.DB) error {
		return db.Select(&items, getFeedSQL, feed)
	})
	return items, err
}
