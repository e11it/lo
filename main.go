package main

import (
	"g.e11it.ru/go/lo/builder"
	"g.e11it.ru/go/lo/loForms"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"html/template"
	"log"
	"net/http"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout:    "layout",
		Funcs: []template.FuncMap{
			{
				"unescaped": func(args ...interface{}) template.HTML {
					return template.HTML(args[0].(string))
				},
			},
		},
	}))

	m.Get("/", func(r render.Render) {
		fd := &loForms.MyForm{
			Age:   22,
			Token: "7238456923847619874612398746374",
		}
		html, _ := builder.FormCreate(fd)
		r.HTML(200, "form", map[string]interface{}{
			"FormBody": html,
		})
	})
	m.Post("/", func(r render.Render, request *http.Request) {
		fd := &loForms.MyForm{}
		if err := builder.FormRead(fd, request); err != nil {
			log.Println("Error:", err.Error())
		}
		builder.DumpForm(fd)
	})

	m.Run()
}
