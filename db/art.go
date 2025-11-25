package db

import (
	"time"

	"github.com/google/uuid"
)

type Art struct {
	Id            string
	ImgURL        string
	ThumbURL      string
	Title         string
	Medium        string
	Width         int
	Height        int
	Year          string
	Description   string
	Sold          bool
	CreatedAt     string
	Ordering      float64
	ShowInGallery bool
}

type ArtPatch struct {
	Title         *string  `json:"title,omitempty"`
	ImgURL        *string  `json:"img_url,omitempty"`
	ThumbURL      *string  `json:"thumb_url,omitempty"`
	Medium        *string  `json:"medium,omitempty"`
	Width         *int     `json:"width,omitempty"`
	Height        *int     `json:"height,omitempty"`
	Year          *string  `json:"year,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Sold          *bool    `json:"sold,omitempty"`
	Ordering      *float64 `json:"ordering,omitempty"`
	ShowInGallery *bool    `json:"show_in_gallery,omitempty"`
}

func (db *DB) createArtTable() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS arts (
		id TEXT PRIMARY KEY,
		img_url TEXT NOT NULL,
		thumb_url TEXT NOT NULL,
		title TEXT NOT NULL,
		medium TEXT NOT NULL,
		width INTEGER NOT NULL,
		height INTEGER NOT NULL,
		year TEXT NOT NULL,
		description TEXT NOT NULL,
		sold BOOLEAN NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL,
		ordering REAL NOT NULL DEFAULT 0,
		show_in_gallery BOOLEAN NOT NULL DEFAULT 1
	);
	`)
	return err
}

func (db *DB) AddArt(art Art) error {
	art.Id = uuid.NewString()
	art.CreatedAt = time.Now().Format(time.RFC3339)

	_, err := db.Exec(`
	INSERT INTO arts (id, img_url, thumb_url, title, medium, width, height, year, description, sold, created_at, ordering)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, art.Id, art.ImgURL, art.ThumbURL, art.Title, art.Medium, art.Width, art.Height, art.Year, art.Description, art.Sold, art.CreatedAt, art.Ordering)
	return err
}

func (db *DB) DeleteArt(id string) error {
	_, err := db.Exec(`DELETE FROM arts WHERE id = ?;`, id)
	return err
}

func (db *DB) GetArtPaged(limit, offset int) ([]Art, error) {
	rows, err := db.Query(`
		SELECT 
			id, 
			img_url, 
			thumb_url, 
			title, 
			medium, 
			width, 
			height, 
			year, 
			description, 
			sold, 
			created_at, 
			ordering 
		FROM 
			arts
		WHERE
			show_in_gallery = 1
		ORDER BY 
			ordering DESC, title ASC 
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var arts []Art
	for rows.Next() {
		var art Art
		if err := rows.Scan(&art.Id, &art.ImgURL, &art.ThumbURL, &art.Title, &art.Medium, &art.Width, &art.Height, &art.Year, &art.Description, &art.Sold, &art.CreatedAt, &art.Ordering); err != nil {
			return nil, err
		}
		arts = append(arts, art)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return arts, nil
}

func (db *DB) GetArts() ([]Art, error) {
	rows, err := db.Query(`SELECT id, img_url, thumb_url, title, medium, width, height, year, description, sold, created_at, ordering, show_in_gallery FROM arts ORDER BY ordering DESC, title ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var arts []Art
	for rows.Next() {
		var art Art
		if err := rows.Scan(&art.Id, &art.ImgURL, &art.ThumbURL, &art.Title, &art.Medium, &art.Width, &art.Height, &art.Year, &art.Description, &art.Sold, &art.CreatedAt, &art.Ordering, &art.ShowInGallery); err != nil {
			return nil, err
		}
		arts = append(arts, art)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return arts, nil
}

func (db *DB) GetArtById(id string) (*Art, error) {
	row := db.QueryRow(`SELECT id, img_url, thumb_url, title, medium, width, height, year, description, sold, created_at, ordering, show_in_gallery FROM arts WHERE id = ?;`, id)

	var art Art
	if err := row.Scan(&art.Id, &art.ImgURL, &art.ThumbURL, &art.Title, &art.Medium, &art.Width, &art.Height, &art.Year, &art.Description, &art.Sold, &art.CreatedAt, &art.Ordering, &art.ShowInGallery); err != nil {
		return nil, err
	}

	return &art, nil
}

func (db *DB) UpdateArtField(id string, field string, value interface{}) error {
	query := "UPDATE arts SET " + field + " = ? WHERE id = ?;"
	_, err := db.Exec(query, value, id)
	return err
}

func (db *DB) UpdateArt(id string, artPatch ArtPatch) error {
	_, err := db.Exec(`
	UPDATE arts
	SET title = COALESCE(?, title),
		medium = COALESCE(?, medium),
		width = COALESCE(?, width),
		height = COALESCE(?, height),
		img_url = COALESCE(?, img_url),
		thumb_url = COALESCE(?, thumb_url),
		year = COALESCE(?, year),
		description = COALESCE(?, description),
		sold = COALESCE(?, sold),
		ordering = COALESCE(?, ordering),
		show_in_gallery = COALESCE(?, show_in_gallery)
	WHERE id = ?;
	`, artPatch.Title, artPatch.Medium, artPatch.Width, artPatch.Height, artPatch.ImgURL, artPatch.ThumbURL, artPatch.Year, artPatch.Description, artPatch.Sold, artPatch.Ordering, artPatch.ShowInGallery, id)
	return err
}
