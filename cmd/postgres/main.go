package postgres

import (
	"flag"
	"log"

	"github.com/warehouse/auth-service/internal/db"
)

func main() {
	var dsn, pgCertLoc string
	flag.StringVar(&dsn, "dsn", "", "")
	flag.StringVar(&pgCertLoc, "cert", "", "")
	flag.Parse()

	c, err := db.NewPostgresClient(dsn, pgCertLoc)
	if err != nil {
		log.Panicln(err)
	}

	c.Close()
}
