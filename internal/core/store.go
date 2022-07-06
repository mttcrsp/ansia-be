package core

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mttcrsp/ansiabe/internal/articles"
	"github.com/mttcrsp/ansiabe/internal/feeds"
	"github.com/mttcrsp/ansiabe/internal/rss"
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

	CREATE TABLE IF NOT EXISTS item_feed (
		item_id INTEGER NOT NULL,
		feed STRING NOT NULL,
		PRIMARY KEY (item_id, feed)
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
	INSERT OR REPLACE INTO item (item_id, title, description, url, published_at)
		VALUES (:item_id, :title, :description, :url, :published_at);
	`
	insertItemFeedSQL = `
	INSERT OR REPLACE INTO item_feed (item_id, feed)
		VALUES (:item_id, :feed);
	`
	insertArticleSQL = `
	INSERT INTO article (item_id, keywords, content, image_url)
		VALUES (:item_id, :keywords, :content, :image_url)
	`
	getFeedSQL = `
	SELECT *
	FROM item i 
		JOIN article a ON i.item_id = a.item_id
		JOIN item_feed if ON i.item_id = if.item_id
	WHERE if.feed = ?
	ORDER BY i.published_at
	LIMIT 30;
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

func (s *Store) InsertFeedItems(feed feeds.Feed, rss rss.RSS) error {
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

		for _, item := range rss.Channel.Items {
			itemRow, err := newItemRow(item)
			if err != nil {
				return err
			}
			if _, err := tx.NamedExec(insertItemSQL, itemRow); err != nil {
				return err
			}

			itemFeedRow := newItemFeedRow(item, feed)
			if _, err := tx.NamedExec(insertItemFeedSQL, itemFeedRow); err != nil {
				return err
			}
		}

		return tx.Commit()
	})
}

func (s *Store) InsertArticle(item rss.Item, article articles.Article) error {
	return s.withDB(func(db *sqlx.DB) error {
		articleRow := *newArticleRow(item, article)
		_, err := db.NamedExec(insertArticleSQL, articleRow)
		return err
	})
}

func (s *Store) GetFeed(feed string) ([]FeedItem, error) {
	items := []FeedItem{}
	err := s.withDB(func(db *sqlx.DB) error {
		return db.Select(&items, getFeedSQL, feed)
	})
	return items, err
}

type itemRow struct {
	ID          int64     `db:"item_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	URL         string    `db:"url"`
	PublishedAt time.Time `db:"published_at"`
}

type itemFeedRow struct {
	ItemID int64  `db:"item_id"`
	Feed   string `db:"feed"`
}

type articleRow struct {
	ItemID   int64  `db:"item_id"`
	Keywords string `db:"keywords"`
	Content  string `db:"content"`
	ImageURL string `db:"image_url"`
}

func newItemRow(item rss.Item) (*itemRow, error) {
	publishedAt, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDateRaw)
	if err != nil {
		return nil, err
	}

	return &itemRow{
		ID:          item.ID(),
		Title:       item.Title,
		Description: item.Description,
		URL:         item.Link,
		PublishedAt: publishedAt,
	}, nil
}

func newItemFeedRow(item rss.Item, feed feeds.Feed) *itemFeedRow {
	return &itemFeedRow{
		ItemID: item.ID(),
		Feed:   feed.Slug(),
	}
}

func newArticleRow(item rss.Item, article articles.Article) *articleRow {
	return &articleRow{
		ItemID:   item.ID(),
		Keywords: article.Keywords,
		Content:  article.Content,
		ImageURL: article.ImageURL,
	}
}
