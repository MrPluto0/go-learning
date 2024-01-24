package main

import (
	"example/gee"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		// c.Fail(500, "Internal Server Error")
		c.Abort()
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.R.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger()) // global middleware
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	{
		v1 := r.Group("/hello")
		v1.GET("/panic", func(c *gee.Context) {
			names := []string{"geektutu"}
			c.String(http.StatusOK, names[100])
		})
		v1.GET("/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	{
		v2 := r.Group("/assets")
		v2.Use(onlyForV2())
		v2.GET("/*string", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"filepath": c.Param("string")})
		})
	}

	r.Run(":9999")
}
