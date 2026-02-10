package main

import (
	"log"
)

func seedDB() {
	initDB("./presentations.db")

	// Real presentation from user (Grade 5 example)
	p1 := Presentation{
		ID:               "math-g5-lesson1",
		Title:            "Бөлшектер (5-сынып)",
		GroupName:        "5-сынып",
		CanvaEmbedURL:    "https://www.canva.com/design/DAG_coZQi0U/q0EOCfdZHzs2jal89EFjqA/view?embed",
		AllowedChannelID: -1003814950604,
	}

	// Dummy presentation for Grade 6
	p2 := Presentation{
		ID:               "math-g6-lesson1",
		Title:            "Пропорция (6-сынып)",
		GroupName:        "6-сынып",
		CanvaEmbedURL:    "https://www.canva.com/design/DAG_coZQi0U/q0EOCfdZHzs2jal89EFjqA/view?embed", // Using same URL for demo
		AllowedChannelID: -1003814950604,
	}

	// Another Grade 5
	p3 := Presentation{
		ID:               "math-g5-lesson2",
		Title:            "Ондық бөлшектер (5-сынып)",
		GroupName:        "5-сынып",
		CanvaEmbedURL:    "https://www.canva.com/design/DAG_coZQi0U/q0EOCfdZHzs2jal89EFjqA/view?embed",
		AllowedChannelID: -1003814950604,
	}

	presentations := []Presentation{p1, p2, p3}

	for _, p := range presentations {
		err := addPresentation(p)
		if err != nil {
			log.Printf("Error adding %s: %v", p.Title, err)
		} else {
			log.Printf("Added: %s", p.Title)
		}
	}
	log.Println("Database seeded with grouped presentations.")
}
