package main

import (
	"context"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/sync/errgroup"
)

var (
	dsnString = flag.String("dsn", "", "connection config to postgres")
	speed     = flag.Int("speed", 3, "speed of generator")

	//go:embed schema.sql
	schema []byte
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dsn, err := pgx.ParseConfig(*dsnString)
	if err != nil {
		panic(err)
	}
	if err := initDB(ctx, dsn); err != nil {
		panic(err)
	}
	gr := errgroup.Group{}
	gr.Go(func() error {
		defer cancel()
		return customerGenerator(ctx, dsn)
	})
	gr.Go(func() error {
		defer cancel()
		return contentGenerator(ctx, dsn)
	})
	if err := gr.Wait(); err != nil {
		panic(err)
	}
}

func initDB(ctx context.Context, dsn *pgx.ConnConfig) error {
	conn := stdlib.OpenDB(*dsn)
	db := New(conn)
	genre, err := db.GetRandomGenre(ctx)
	if err == nil && genre.GenreID > 0 {
		return nil
	}
	if _, err := db.db.ExecContext(ctx, string(schema)); err != nil {
		return err
	}
	return nil
}

func contentGenerator(ctx context.Context, dsn *pgx.ConnConfig) error {
	conn := stdlib.OpenDB(*dsn)
	db := New(conn)
	for {
		genre, err := db.GetRandomGenre(ctx)
		if err != nil {
			return err
		}
		format, err := db.GetRandomMediaType(ctx)
		if err != nil {
			return err
		}
		var artist Artist
		if IntnRange(1, 100) > 70 {
			artist, err = db.CreateArtist(
				ctx,
				sql.NullString{String: fmt.Sprintf("%s %s %s", gofakeit.HipsterWord(), gofakeit.NounAbstract(), gofakeit.Noun())},
			)
		} else {
			artist, err = db.GetRandomArtist(ctx)
		}
		if err != nil {
			return err
		}
		album, err := db.CreateAlbum(ctx, CreateAlbumParams{
			Title:    gofakeit.HipsterSentence(IntnRange(1, 5)),
			ArtistID: artist.ArtistID,
		})
		if err != nil {
			return err
		}
		var tracks []Track
		for i := 0; i <= IntnRange(1, 15); i++ {
			track, err := db.CreateTrack(ctx, CreateTrackParams{
				Name:         gofakeit.HipsterSentence(IntnRange(1, 7)),
				AlbumID:      sql.NullInt32{Int32: album.AlbumID},
				MediaTypeID:  format.MediaTypeID,
				GenreID:      sql.NullInt32{Int32: genre.GenreID},
				Composer:     sql.NullString{String: gofakeit.Name()},
				Milliseconds: int32(IntnRange(120, 750)),
				Bytes:        sql.NullInt32{Int32: int32(IntnRange(1*1024*1024, 15*1024*1024))},
				UnitPrice:    fmt.Sprintf("%v", IntnRange(10, 599)/100),
			})
			if err != nil {
				return err
			}
			tracks = append(tracks, track)
		}
		log.Printf("done inserting album: %v, with %v tracks for: %v", album.Title, len(tracks), artist.Name.String)
		time.Sleep(time.Duration(IntnRange(1000, 2500)) * time.Millisecond * time.Duration(*speed))
	}
}

func customerGenerator(ctx context.Context, dsn *pgx.ConnConfig) error {
	conn := stdlib.OpenDB(*dsn)
	db := New(conn)
	for {
		addr := gofakeit.Address()
		customer, err := db.CreateCustomer(ctx, CreateCustomerParams{
			FirstName:    gofakeit.FirstName(),
			LastName:     gofakeit.LastName(),
			Company:      sql.NullString{String: gofakeit.Company()},
			Address:      sql.NullString{String: addr.Address},
			City:         sql.NullString{String: addr.City},
			State:        sql.NullString{String: addr.State},
			Country:      sql.NullString{String: addr.Country},
			PostalCode:   sql.NullString{String: addr.Zip},
			Phone:        sql.NullString{String: gofakeit.Phone()},
			Fax:          sql.NullString{String: ""},
			Email:        gofakeit.Email(),
			SupportRepID: sql.NullInt32{},
		})
		if err != nil {
			return err
		}
		invoice, err := db.CreateInvoice(ctx, CreateInvoiceParams{
			CustomerID:        customer.CustomerID,
			InvoiceDate:       time.Time{},
			BillingAddress:    sql.NullString{String: addr.Address},
			BillingCity:       sql.NullString{String: addr.City},
			BillingState:      sql.NullString{String: addr.State},
			BillingCountry:    sql.NullString{String: addr.Country},
			BillingPostalCode: sql.NullString{String: gofakeit.Phone()},
			Total:             fmt.Sprintf("%v", int32(IntnRange(10, 599))),
		})
		if err != nil {
			return err
		}
		tracks, err := db.GetRandomTracks(ctx, int32(IntnRange(3, 25)))
		if err != nil {
			return err
		}
		for _, track := range tracks {
			if err := db.CreateInvoiceLine(ctx, CreateInvoiceLineParams{
				InvoiceID: invoice.InvoiceID,
				TrackID:   track.TrackID,
				UnitPrice: track.UnitPrice,
				Quantity:  int32(IntnRange(1, 3)),
			}); err != nil {
				return err
			}
		}
		log.Printf("done inserting customer: %v, with %v tracks, total: %v", customer.CustomerID, len(tracks), invoice.Total)
		time.Sleep(time.Duration(IntnRange(1, 2500)) * time.Millisecond * time.Duration(*speed))
	}
}

func IntnRange(min, max int) int {
	return rand.Intn(max-min) + min
}
