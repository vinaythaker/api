package api

import (
	"database/sql"
	"fmt"
)

//Pet ...
type Pet struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Category  Category `json:"category,omitempty"`
	Status    string   `json:"status,omitempty"`
	Tags      []Tag    `json:"tags,omitempty"`
	PhotoUrls []string `json:"photoUrls,omitempty"`
}

func (p *Pet) addPet(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO pets(name) VALUES('%s')", p.Name)
	result, err := db.Exec(statement)
	if err != nil {
		return err
	}

	p.ID, _ = result.LastInsertId()

	//if err := db.QueryRow("SELECT currval('pets_id_seq')").Scan(&p.ID); err != nil {
	//	return err
	//}

	return nil
}

func (p *Pet) getPet(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, name FROM pets WHERE id=%d", p.ID)
	return db.QueryRow(statement).Scan(&p.ID, &p.Name)
}

func (p *Pet) getPets(db *sql.DB, start, count int) ([]Pet, error) {
	statement := fmt.Sprintf("SELECT id, name FROM pets LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pets := []Pet{}

	for rows.Next() {
		var p Pet
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		pets = append(pets, p)
	}

	return pets, nil
}

func (p *Pet) updatePet(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE pets SET name='%s' WHERE id=%d", p.Name, p.ID)
	result, err := db.Exec(statement)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (p *Pet) deletePet(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM pets WHERE id=%d", p.ID)
	result, err := db.Exec(statement)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}

	return nil
}
