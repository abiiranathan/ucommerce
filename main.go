package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/abiiranathan/gora/env"
	"github.com/abiiranathan/gora/gora"
	"github.com/abiiranathan/gora/ws"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const dotenvFileName = "app.env"

type Config struct {
	DatabaseUrl      string `name:"DATABASE_URL" required:"true"`
	RedisUrl         string `name:"REDIS_URL" required:"true"`
	PostgresUser     string `name:"POSTGRES_USER" required:"true"`
	PostgresPassword string `name:"POSTGRES_PASSWORD" required:"true"`
	PostgresDB       string `name:"POSTGRES_DB" required:"true"`
	MigrationUrl     string `name:"MIGRATION_URL" required:"true"`
}

// migrationUrl: Location of migration files. Can on local machine or remote-host.
//
// dbUrl: postgres connection string
func RunDBMigrations(migrationUrl string, db *sql.DB, databaseName string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("can not get the driver: %v\n", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(migrationUrl, databaseName, driver)
	if err != nil {
		log.Fatalf("can not create a new migrate instance: %v\n", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrate up: %v\n", err)
	}

	log.Println("db migrated successfully")
}

func main() {
	config := &Config{}
	err := env.LoadConfig(dotenvFileName, config)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("pgx", config.DatabaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	RunDBMigrations(config.MigrationUrl, db, config.PostgresDB)

	r := gora.Default(os.Stdout)

	r.GET("/", func(c *gora.Context) { c.String("Hello World") })
	r.GET("/metrics", gora.WrapH(promhttp.Handler()))

	hub, quit := ws.NewHandler()
	defer quit()
	go hub.Run()

	r.GET("/ws", gora.WrapH(hub))
	r.Run(":8000")
}
