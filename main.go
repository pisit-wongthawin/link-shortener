package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var db *sql.DB

type Url struct {
	ID       int    `json:"id"`
	Original string `json:"original"`
	Shorten  string `json:"shorten"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		host         = os.Getenv("DATABASE_HOST")
		port         = os.Getenv("DATABASE_PORT")
		databaseName = os.Getenv("DATABASE_NAME")
		username     = os.Getenv("DATABASE_USERNAME")
		password     = os.Getenv("DATABASE_PASSWORD")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, databaseName)

	sdb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	db = sdb

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected")

	app := fiber.New()

	app.Post("/shorten", createUrlHandler)
	app.Get("/:shorten", getUrlHandler)

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func randomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func getUrlHandler(c *fiber.Ctx) error {
	shorten := c.Params("shorten")
	originalUrl, err := getOriginalUrl(shorten)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	c.Set("Location", originalUrl)

	return c.SendStatus(307)
}

func createUrlHandler(c *fiber.Ctx) error {
	url := new(Url)

	if err := c.BodyParser(url); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	shorten := randomString(8)

	err := createUrl(&Url{Original: url.Original, Shorten: shorten})
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendString(shorten)
}
