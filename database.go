package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB(dataSourceName string) {
	var err error
	db, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

    createTableSQL := `CREATE TABLE IF NOT EXISTS presentations (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        group_name TEXT NOT NULL,
        canva_embed_url TEXT NOT NULL,
        allowed_channel_id INTEGER NOT NULL
    );`

    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatalf("Error creating table: %v", err)
    }
}

type Presentation struct {
    ID               string
    Title            string
    GroupName        string
    CanvaEmbedURL    string
    AllowedChannelID int64
}

func getPresentation(id string) (*Presentation, error) {
    row := db.QueryRow("SELECT id, title, group_name, canva_embed_url, allowed_channel_id FROM presentations WHERE id = ?", id)
    var p Presentation
    err := row.Scan(&p.ID, &p.Title, &p.GroupName, &p.CanvaEmbedURL, &p.AllowedChannelID)
    if err != nil {
        return nil, err
    }
    return &p, nil
}

func addPresentation(p Presentation) error {
    _, err := db.Exec("INSERT INTO presentations (id, title, group_name, canva_embed_url, allowed_channel_id) VALUES (?, ?, ?, ?, ?)", p.ID, p.Title, p.GroupName, p.CanvaEmbedURL, p.AllowedChannelID)
    return err
}

func getAllPresentations() ([]Presentation, error) {
    rows, err := db.Query("SELECT id, title, group_name, canva_embed_url, allowed_channel_id FROM presentations ORDER BY group_name, title")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var presentations []Presentation
    for rows.Next() {
        var p Presentation
        if err := rows.Scan(&p.ID, &p.Title, &p.GroupName, &p.CanvaEmbedURL, &p.AllowedChannelID); err != nil {
            return nil, err
        }
        presentations = append(presentations, p)
    }
    return presentations, nil
}
