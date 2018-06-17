package main

import (
	"database/sql"
	"errors"
)

type phone struct {
	Phone     string `json:"phone,omitempty"`
	Company   string `json:"company,omitempty"`
	PhoneType string `json:"phone_type,omitempty"`
	UserID    string `json:"userId,omitempty"`
}

func (p *phone) createPhone(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *phone) deletePhone(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *phone) getPhone(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getPhones(db *sql.DB, start, count int) ([]phone, error) {
	return nil, errors.New("Not implemented")
}

func (p *phone) updatePhone(db *sql.DB) error {
	return errors.New("Not implemented")
}
