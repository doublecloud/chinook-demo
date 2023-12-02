package main

import (
	"database/sql"
	"time"
)

type Album struct {
	AlbumID  int32
	Title    string
	ArtistID int32
}

type Artist struct {
	ArtistID int32
	Name     sql.NullString
}

type Customer struct {
	CustomerID   int32
	FirstName    string
	LastName     string
	Company      sql.NullString
	Address      sql.NullString
	City         sql.NullString
	State        sql.NullString
	Country      sql.NullString
	PostalCode   sql.NullString
	Phone        sql.NullString
	Fax          sql.NullString
	Email        string
	SupportRepID sql.NullInt32
}

type Genre struct {
	GenreID int32
	Name    sql.NullString
}

type Invoice struct {
	InvoiceID         int32
	CustomerID        int32
	InvoiceDate       time.Time
	BillingAddress    sql.NullString
	BillingCity       sql.NullString
	BillingState      sql.NullString
	BillingCountry    sql.NullString
	BillingPostalCode sql.NullString
	Total             string
}

type InvoiceLine struct {
	InvoiceLineID int32
	InvoiceID     int32
	TrackID       int32
	UnitPrice     string
	Quantity      int32
}

type MediaType struct {
	MediaTypeID int32
	Name        sql.NullString
}

type Track struct {
	TrackID      int32
	Name         string
	AlbumID      sql.NullInt32
	MediaTypeID  int32
	GenreID      sql.NullInt32
	Composer     sql.NullString
	Milliseconds int32
	Bytes        sql.NullInt32
	UnitPrice    string
}
