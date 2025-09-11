package database

import (
	"log"
)
//миграция бд 
func Migrate() {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT,
        completed BOOLEAN DEFAULT FALSE,
        priority VARCHAR(10) DEFAULT 'medium',
        due_date TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to migrate tasks table:", err)
	}

	log.Println("Tasks table migrated successfully")
}
