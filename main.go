package main

import (
	"my-keep-backend/db"
	"my-keep-backend/models"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"xorm.io/builder"
)

/*
type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	//temp := v.X*v.X + v.Y*v.Y

	return math.Sqrt(v.X*v.X + v.Y*v.Y)
	//(30 * 30) + (40 * 40)
}

func (v *Vertex) Scale(f float64) {
	fmt.Println("X=", v.X)
	fmt.Println("Y=", v.X)
	v.X = v.X * f // X= 3 * 10
	v.Y = v.Y * f // Y = 4 * 10
}

func main() {
	v := Vertex{3, 4}
	fmt.Println(v.X)
	fmt.Println(v.Y)
	v.Scale(10)
	fmt.Println(v.Abs())
}
*/
var dbClient db.Db

func CORSMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin, err := url.PathUnescape(c.Query("origin"))
		if err == nil && origin != "" {
			allowOrigin = origin
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Max-Age, Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func main() {
	dbClient = db.Db{}
	dbClient.Connect("postgres", "M1ft4hul", "localhost", 5432, "postgres", "postgres")
	route := gin.Default()
	route.Use(CORSMiddleware(""))
	route.GET("/ping", func(c *gin.Context) {
		username := "tanto"
		print(username)
		c.JSON(200, gin.H{
			"message": "pong",
			"time":    time.Now(),
		})
	})

	route.POST("/notes", func(c *gin.Context) {
		var note models.Note
		var err error
		newId, _ := uuid.NewUUID()
		note.Id = newId.String()
		c.Bind(&note)
		_, err = dbClient.Conn.Insert(&note)
		if err != nil {
			c.JSON(503, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, note)
		}
	})

	route.PUT("/notes/:id", func(c *gin.Context) {
		id := c.Param("id")
		var note models.Note
		var err error
		c.Bind(&note)
		_, err = dbClient.Conn.Where(builder.Eq{"id": id}).Update(&note)
		if err != nil {
			c.JSON(503, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, note)
		}
	})
	route.DELETE("/notes/:id", func(c *gin.Context) {
		id := c.Param("id")
		var note models.Note
		var err error
		c.Bind(&note)
		_, err = dbClient.Conn.Where(builder.Eq{"id": id}).Delete(&note)
		if err != nil {
			c.JSON(503, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, note)
		}
	})

	route.GET("/notes/:id", func(c *gin.Context) {
		id := c.Param("id")
		var note models.Note
		_, err := dbClient.Conn.Where(builder.Eq{"id": id}).Get(&note)
		if err != nil {
			c.JSON(503, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, note)
		}

	})
	route.GET("/notes", func(c *gin.Context) {
		//TODO: pagination
		var notes []models.Note
		err := dbClient.Conn.Find(&notes)
		if err != nil {
			c.JSON(503, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, notes)
		}

	})
	route.POST("/sign-in", SignIn)

	route.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
