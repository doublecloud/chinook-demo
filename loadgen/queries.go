package main

import (
	"context"
	"database/sql"
	"time"
)

const createAlbum = `-- name: CreateAlbum :one
insert into album (title, artist_id)
values ($1, $2)
RETURNING album_id, title, artist_id
`

type CreateAlbumParams struct {
	Title    string
	ArtistID int32
}

func (q *Queries) CreateAlbum(ctx context.Context, arg CreateAlbumParams) (Album, error) {
	row := q.db.QueryRowContext(ctx, createAlbum, arg.Title, arg.ArtistID)
	var i Album
	err := row.Scan(&i.AlbumID, &i.Title, &i.ArtistID)
	return i, err
}

const createArtist = `-- name: CreateArtist :one
insert into artist (name)
values ($1)
RETURNING artist_id, name
`

func (q *Queries) CreateArtist(ctx context.Context, name sql.NullString) (Artist, error) {
	row := q.db.QueryRowContext(ctx, createArtist, name)
	var i Artist
	err := row.Scan(&i.ArtistID, &i.Name)
	return i, err
}

const createCustomer = `-- name: CreateCustomer :one
insert into customer (first_name,
                      last_name, company, address,
                      city, state, country, postal_code, phone,
                      fax, email, support_rep_id)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING customer_id, first_name, last_name, company, address, city, state, country, postal_code, phone, fax, email, support_rep_id
`

type CreateCustomerParams struct {
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

func (q *Queries) CreateCustomer(ctx context.Context, arg CreateCustomerParams) (Customer, error) {
	row := q.db.QueryRowContext(ctx, createCustomer,
		arg.FirstName,
		arg.LastName,
		arg.Company,
		arg.Address,
		arg.City,
		arg.State,
		arg.Country,
		arg.PostalCode,
		arg.Phone,
		arg.Fax,
		arg.Email,
		arg.SupportRepID,
	)
	var i Customer
	err := row.Scan(
		&i.CustomerID,
		&i.FirstName,
		&i.LastName,
		&i.Company,
		&i.Address,
		&i.City,
		&i.State,
		&i.Country,
		&i.PostalCode,
		&i.Phone,
		&i.Fax,
		&i.Email,
		&i.SupportRepID,
	)
	return i, err
}

const createInvoice = `-- name: CreateInvoice :one
insert into invoice (customer_id, invoice_date,
                     billing_address, billing_city, billing_state,
                     billing_country, billing_postal_code, total)
values ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING invoice_id, customer_id, invoice_date, billing_address, billing_city, billing_state, billing_country, billing_postal_code, total
`

type CreateInvoiceParams struct {
	CustomerID        int32
	InvoiceDate       time.Time
	BillingAddress    sql.NullString
	BillingCity       sql.NullString
	BillingState      sql.NullString
	BillingCountry    sql.NullString
	BillingPostalCode sql.NullString
	Total             string
}

func (q *Queries) CreateInvoice(ctx context.Context, arg CreateInvoiceParams) (Invoice, error) {
	row := q.db.QueryRowContext(ctx, createInvoice,
		arg.CustomerID,
		arg.InvoiceDate,
		arg.BillingAddress,
		arg.BillingCity,
		arg.BillingState,
		arg.BillingCountry,
		arg.BillingPostalCode,
		arg.Total,
	)
	var i Invoice
	err := row.Scan(
		&i.InvoiceID,
		&i.CustomerID,
		&i.InvoiceDate,
		&i.BillingAddress,
		&i.BillingCity,
		&i.BillingState,
		&i.BillingCountry,
		&i.BillingPostalCode,
		&i.Total,
	)
	return i, err
}

const createInvoiceLine = `-- name: CreateInvoiceLine :exec
insert into invoice_line (invoice_id, track_id, unit_price, quantity)
values ($1, $2, $3, $4)
`

type CreateInvoiceLineParams struct {
	InvoiceID int32
	TrackID   int32
	UnitPrice string
	Quantity  int32
}

func (q *Queries) CreateInvoiceLine(ctx context.Context, arg CreateInvoiceLineParams) error {
	_, err := q.db.ExecContext(ctx, createInvoiceLine,
		arg.InvoiceID,
		arg.TrackID,
		arg.UnitPrice,
		arg.Quantity,
	)
	return err
}

const createTrack = `-- name: CreateTrack :one
insert into track (
                   name, album_id,
                   media_type_id, genre_id,
                   composer, milliseconds,
                   bytes, unit_price
)
values ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING track_id, name, album_id, media_type_id, genre_id, composer, milliseconds, bytes, unit_price
`

type CreateTrackParams struct {
	Name         string
	AlbumID      sql.NullInt32
	MediaTypeID  int32
	GenreID      sql.NullInt32
	Composer     sql.NullString
	Milliseconds int32
	Bytes        sql.NullInt32
	UnitPrice    string
}

func (q *Queries) CreateTrack(ctx context.Context, arg CreateTrackParams) (Track, error) {
	row := q.db.QueryRowContext(ctx, createTrack,
		arg.Name,
		arg.AlbumID,
		arg.MediaTypeID,
		arg.GenreID,
		arg.Composer,
		arg.Milliseconds,
		arg.Bytes,
		arg.UnitPrice,
	)
	var i Track
	err := row.Scan(
		&i.TrackID,
		&i.Name,
		&i.AlbumID,
		&i.MediaTypeID,
		&i.GenreID,
		&i.Composer,
		&i.Milliseconds,
		&i.Bytes,
		&i.UnitPrice,
	)
	return i, err
}

const getRandomGenre = `-- name: GetRandomGenre :one
select genre_id, name from genre order by random() limit 1
`

func (q *Queries) GetRandomGenre(ctx context.Context) (Genre, error) {
	row := q.db.QueryRowContext(ctx, getRandomGenre)
	var i Genre
	err := row.Scan(&i.GenreID, &i.Name)
	return i, err
}

const getRandomMediaType = `-- name: GetRandomMediaType :one
select media_type_id, name from media_type order by random() limit 1
`

func (q *Queries) GetRandomMediaType(ctx context.Context) (MediaType, error) {
	row := q.db.QueryRowContext(ctx, getRandomMediaType)
	var i MediaType
	err := row.Scan(&i.MediaTypeID, &i.Name)
	return i, err
}

const getRandomTracks = `-- name: GetRandomTracks :many
select track_id, name, album_id, media_type_id, genre_id, composer, milliseconds, bytes, unit_price from track order by random() limit $1
`

func (q *Queries) GetRandomTracks(ctx context.Context, limit int32) ([]Track, error) {
	rows, err := q.db.QueryContext(ctx, getRandomTracks, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Track
	for rows.Next() {
		var i Track
		if err := rows.Scan(
			&i.TrackID,
			&i.Name,
			&i.AlbumID,
			&i.MediaTypeID,
			&i.GenreID,
			&i.Composer,
			&i.Milliseconds,
			&i.Bytes,
			&i.UnitPrice,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRandomArtist = `-- name: GetRandomArtist :one
select artist_id, name from artist order by random() limit 1
`

func (q *Queries) GetRandomArtist(ctx context.Context) (Artist, error) {
	row := q.db.QueryRowContext(ctx, getRandomArtist)
	var i Artist
	err := row.Scan(&i.ArtistID, &i.Name)
	return i, err
}
