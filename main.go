package main

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	_ "github.com/lib/pq"
	"github.com/pdmp/retina/models"
)

func newApp() *iris.Application {

	app := iris.New()
	orm, err := xorm.NewEngine("postgres", "dbname=retina_dev user=postgres password=root sslmode=disable")
	if err != nil {
		app.Logger().Fatalf("orm failed to initialized: %v", err)
	}

	iris.RegisterOnInterrupt(func() {
		orm.Close()
	})

	err = orm.Sync2(new(models.User))

	if err != nil {
		app.Logger().Fatalf("orm failed to initialized User table: %v", err)
	}

	app.RegisterView(iris.HTML("./public/lens/dist/lens", ".html"))
	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	assetHandler := app.StaticHandler("./public/lens/dist/lens", false, false)
	app.SPA(assetHandler)

	pass, _ := models.GeneratePassword("haha")
	app.Get("/insert", func(ctx iris.Context) {
		user := &models.User{Firstname: "kataras", Username: "test", HashedPassword: pass, CreatedAt: time.Now()}
		orm.Insert(user)

		ctx.Writef("user inserted: %#v", user)
	})

	app.Get("/get", func(ctx iris.Context) {
		user := models.User{ID: 1}
		if ok, _ := orm.Get(&user); ok {
			ctx.Writef("user found: %#v", user)
		}
	})

	return app
}

func main() {
	app := newApp()

	app.Run(iris.Addr("localhost:8000"))
}
