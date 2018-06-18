package main

import (
	"database/sql"
)

// PhonesTableCreationQuery SQL to create the table phones
const PhonesTableCreationQuery = `CREATE TABLE IF NOT EXISTS phones
(
phone TEXT NOT NULL,
company TEXT NOT NULL,
phoneType TEXT NOT NULL,
userId TEXT NOT NULL,
CONSTRAINT phones_pkey PRIMARY KEY (phone)
)`

type phone struct {
	Phone     string `json:"phone,omitempty"`
	Company   string `json:"company,omitempty"`
	PhoneType string `json:"phone_type,omitempty"`
	UserID    string `json:"userId,omitempty"`
}

func (p *phone) createPhone(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO phones(phone, company, phoneType, userId) VALUES($1, $2, $3, $4) RETURNING phone",
		p.Phone, p.Company, p.PhoneType, p.UserID).Scan(&p.Phone)

	if err != nil {
		return err
	}

	return nil
}

func (p *phone) deletePhone(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM phones WHERE phone=$1", p.Phone)

	return err
}

func (p *phone) getPhone(db *sql.DB) error {
	return db.QueryRow("SELECT phone, company, phoneType, userId FROM phones WHERE phone=$1",
		p.Phone).Scan(&p.Phone, &p.Company, &p.PhoneType, &p.UserID)
}

func getPhones(db *sql.DB, start, count int) ([]phone, error) {
	rows, err := db.Query(
		"SELECT phone, company, phoneType, userId FROM phones LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	phones := []phone{}

	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.Phone, &p.Company, &p.PhoneType, &p.UserID); err != nil {
			return nil, err
		}
		phones = append(phones, p)
	}

	return phones, nil
}

func (p *phone) updatePhone(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE phones SET phoneNumber=$1, company=$2, phoneType=$3, userId=$4 WHERE phone=$5",
			p.Phone, p.Company, p.PhoneType, p.UserID, p.Phone)

	return err
}
