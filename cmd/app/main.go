package main

import (
	"fmt"
	"github.com/ndodanli/go-clean-architecture/pkg/config"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/datastore"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/router"
	"github.com/ndodanli/go-clean-architecture/pkg/registry"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)

func main() {
	config.ReadConfig()

	db := datastore.NewDB()
	db.LogMode(true)
	defer db.Close()

	r := registry.NewRegistry(db)

	e := echo.New()
	e = router.NewRouter(e, r.NewAppController())

	fmt.Println("Server listen at http://localhost" + ":" + config.C.Server.Address)
	if err := e.Start(":" + config.C.Server.Address); err != nil {
		log.Fatalln(err)
	}
}
