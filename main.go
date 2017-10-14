// package main contains an example on how to use the ReadForm, but with the same way you can do the ReadJSON & ReadJSON
package main

import (
	"github.com/kataras/iris"
	arango "github.com/diegogub/aranGO"
	"fmt"
)

type Users struct {
	arango.Document
	Username string `json:"username"`
	Password string `json:"password"`
	Correo    string `json:"correo"`
	Permisos []string `form:"permission"`
}

type Permiso struct{
	Nombre string
}

func main() {
	session, err := arango.Connect("http://192.168.0.106:8529", "root", "canaima", false)
	if err != nil {
		panic(err)
	}

	// session.CreateDB("proyect", nil)

	// if !session.DB("proyect").ColExist("usuarios") {

	// 	nueva_coleccion := arango.NewCollectionOptions("usuarios", true)

	// 	session.DB("proyect").CreateCollection(nueva_coleccion)
	// }

	app := iris.New()
	app.RegisterView(iris.HTML("./templates", ".html").Reload(true))

	app.Get("/", func(ctx iris.Context) {
		if err := ctx.View("index.html"); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
		}
	})

	app.Post("/", func(ctx iris.Context) {
		users := Users{}
		err := ctx.ReadForm(&users)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
		}

		//Insertar registro
		err = session.DB("proyect").Col("usuarios").Save(&users)
		ctx.JSON(users)
	})

	app.Get("/users", func(ctx iris.Context) {
		_q := `
			FOR usuario in usuarios
			RETURN usuario
		`
		query := arango.NewQuery(_q)

		cursor, _ := session.DB("proyect").Execute(query)

		if err != nil {
			panic(err)
		}

		var _users []Users

		cursor.FetchBatch(&_users)

		ctx.JSON(_users)
	})

	app.Get("/users/{id:int}", func(ctx iris.Context) {

		id, _ := ctx.Params().GetInt("id")

		_q := `
			FOR usuario in usuarios
			FILTER usuario._key == "%d"
			RETURN usuario
		`
		query := arango.NewQuery(fmt.Sprintf(_q, id))

		cursor, _ := session.DB("proyect").Execute(query)

		if err != nil {
			panic(err)
		}

		var usuario_buscado Users

		cursor.FetchOne(&usuario_buscado)

		ctx.JSON(usuario_buscado)
	})

	app.Put("/users/{id:int}", func(ctx iris.Context) {

		id, _ := ctx.Params().GetInt("id")

		_q := `
			FOR usuario in usuarios
			FILTER usuario._key == "%d"
			RETURN usuario
		`
		query := arango.NewQuery(fmt.Sprintf(_q, id))

		cursor, err := session.DB("proyect").Execute(query)

		if err != nil {
			panic(err)
		}

		var usuario_act Users

		cursor.FetchOne(&usuario_act)

		fmt.Println(usuario_act)

		usuario_act.Username = "hoalsfgsd"

		err = session.DB("proyect").Col("usuarios").Replace(usuario_act.Key, usuario_act)
		
		if err != nil {
			panic(err)
		}

		ctx.JSON(usuario_act)
	})

	app.Delete("/users/{id:string}", func(ctx iris.Context) {

		id := ctx.Params().Get("id")

		//Eliminar
		err := session.DB("proyect").Col("usuarios").Delete(id)

		if err != nil {
			panic(err)
		}

	})

	app.Run(iris.Addr(":8080"))
}
