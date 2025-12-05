package db

import (
	"time"

	"github.com/google/uuid"
)

type Print struct {
	Id           string
	ImgURL       string
	ThumbURL     string
	Title        string
	Medium       string
	Width        int
	Height       int
	Year         string
	Description  string
	Price        float64
	QuantityLeft int
	CreatedAt    string
	Ordering     float64
	ShowInStore  bool
}

type PrintPatch struct {
	Title        *string  `json:"title,omitempty"`
	ImgURL       *string  `json:"img_url,omitempty"`
	ThumbURL     *string  `json:"thumb_url,omitempty"`
	Medium       *string  `json:"medium,omitempty"`
	Width        *int     `json:"width,omitempty"`
	Height       *int     `json:"height,omitempty"`
	Year         *string  `json:"year,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	QuantityLeft *int     `json:"quantity_left,omitempty"`
	Ordering     *float64 `json:"ordering,omitempty"`
	ShowInStore  *bool    `json:"show_in_store,omitempty"`
}

func (db *DB) createPrintTable() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS prints (
		id TEXT PRIMARY KEY,
		img_url TEXT NOT NULL,
		thumb_url TEXT NOT NULL,
		title TEXT NOT NULL,
		medium TEXT NOT NULL,
		width INTEGER NOT NULL,
		height INTEGER NOT NULL,
		year TEXT NOT NULL,
		description TEXT NOT NULL,
		price REAL NOT NULL,
		quantity_left INTEGER NOT NULL,
		created_at TEXT NOT NULL,
		ordering REAL NOT NULL DEFAULT 0,
		show_in_store BOOLEAN NOT NULL DEFAULT 1
	);
	`)
	return err
}

func (db *DB) AddPrint(print Print) error {
	print.Id = uuid.NewString()
	print.CreatedAt = time.Now().Format(time.RFC3339)

	_, err := db.Exec(`
	INSERT INTO prints (id, img_url, thumb_url, title, medium, width, height, year, description, price, quantity_left, created_at, ordering, show_in_store)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, print.Id, print.ImgURL, print.ThumbURL, print.Title, print.Medium, print.Width, print.Height, print.Year, print.Description, print.Price, print.QuantityLeft, print.CreatedAt, print.Ordering, print.ShowInStore)
	return err
}

func (db *DB) DeletePrint(id string) error {
	_, err := db.Exec(`DELETE FROM prints WHERE id = ?;`, id)
	return err
}

func (db *DB) GetPrintPaged(limit, offset int) ([]Print, error) {
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
			price, 
			quantity_left, 
			created_at, 
			ordering, 
			show_in_store 
		FROM 
			prints
		WHERE
			show_in_store = 1
		ORDER BY 
			ordering DESC, title ASC 
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prints []Print
	for rows.Next() {
		var print Print
		if err := rows.Scan(&print.Id, &print.ImgURL, &print.ThumbURL, &print.Title, &print.Medium, &print.Width, &print.Height, &print.Year, &print.Description, &print.Price, &print.QuantityLeft, &print.CreatedAt, &print.Ordering, &print.ShowInStore); err != nil {
			return nil, err
		}
		prints = append(prints, print)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prints, nil
}

func (db *DB) GetAllPrints() ([]Print, error) {
	rows, err := db.Query(`SELECT id, img_url, thumb_url, title, medium, width, height, year, description, price, quantity_left, created_at, ordering, show_in_store FROM prints ORDER BY ordering DESC, title ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prints []Print
	for rows.Next() {
		var print Print
		if err := rows.Scan(&print.Id, &print.ImgURL, &print.ThumbURL, &print.Title, &print.Medium, &print.Width, &print.Height, &print.Year, &print.Description, &print.Price, &print.QuantityLeft, &print.CreatedAt, &print.Ordering, &print.ShowInStore); err != nil {
			return nil, err
		}
		prints = append(prints, print)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prints, nil
}

func (db *DB) GetPrintsForStore() ([]Print, error) {
	rows, err := db.Query(`SELECT id, img_url, thumb_url, title, medium, width, height, year, description, price, quantity_left, created_at, ordering, show_in_store FROM prints WHERE show_in_store = 1 ORDER BY ordering DESC, title ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prints []Print
	for rows.Next() {
		var print Print
		if err := rows.Scan(&print.Id, &print.ImgURL, &print.ThumbURL, &print.Title, &print.Medium, &print.Width, &print.Height, &print.Year, &print.Description, &print.Price, &print.QuantityLeft, &print.CreatedAt, &print.Ordering, &print.ShowInStore); err != nil {
			return nil, err
		}
		prints = append(prints, print)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prints, nil
}

func (db *DB) GetPrintById(id string) (*Print, error) {
	row := db.QueryRow(`SELECT id, img_url, thumb_url, title, medium, width, height, year, description, price, quantity_left, created_at, ordering, show_in_store FROM prints WHERE id = ?;`, id)
	var print Print
	if err := row.Scan(&print.Id, &print.ImgURL, &print.ThumbURL, &print.Title, &print.Medium, &print.Width, &print.Height, &print.Year, &print.Description, &print.Price, &print.QuantityLeft, &print.CreatedAt, &print.Ordering, &print.ShowInStore); err != nil {
		return nil, err
	}

	return &print, nil
}

func (db *DB) UpdatePrintField(id string, field string, value any) error {
	query := "UPDATE prints SET " + field + " = ? WHERE id = ?;"
	_, err := db.Exec(query, value, id)
	return err
}

func (db *DB) UpdatePrint(id string, printPatch PrintPatch) error {
	_, err := db.Exec(`
	UPDATE prints
	SET title = COALESCE(?, title),
		medium = COALESCE(?, medium),
		width = COALESCE(?, width),
		height = COALESCE(?, height),
		img_url = COALESCE(?, img_url),
		thumb_url = COALESCE(?, thumb_url),
		year = COALESCE(?, year),
		description = COALESCE(?, description),
		quantity_left = COALESCE(?, quantity_left),
		price = COALESCE(?, price),
		ordering = COALESCE(?, ordering),
		show_in_store = COALESCE(?, show_in_store)
	WHERE id = ?;
	`, printPatch.Title, printPatch.Medium, printPatch.Width, printPatch.Height, printPatch.ImgURL, printPatch.ThumbURL, printPatch.Year, printPatch.Description, printPatch.QuantityLeft, printPatch.Price, printPatch.Ordering, printPatch.ShowInStore, id)
	return err
}
