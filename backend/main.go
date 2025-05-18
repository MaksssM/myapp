package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"tipatwitter/backend/recommendations"
	"tipatwitter/backend/subscriptions"
	"tipatwitter/database"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Post struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Media   string `json:"media,omitempty"`
	Date    string `json:"date"`
}

var (
	PostgresPool *pgxpool.Pool
	Neo4jDriver  neo4j.Driver
)

func init() {
	if !isDeviceConnected() {
		log.Println("ВНИМАНИЕ: Не обнаружено подключённого мобильного устройства или эмулятора. Проверьте подключение перед запуском клиента.")
	}
}

func initDatabaseConnections() {
	// PostgreSQL connection
	pgURL := os.Getenv("POSTGRES_URL")
	var err error
	PostgresPool, err = pgxpool.Connect(context.Background(), pgURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Neo4j connection
	neo4jURL := os.Getenv("NEO4J_URL")
	neo4jUser := os.Getenv("NEO4J_USER")
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")
	Neo4jDriver, err = neo4j.NewDriver(neo4jURL, neo4j.BasicAuth(neo4jUser, neo4jPassword, ""))
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
}

// --- ПРОВЕРКА ПОДКЛЮЧЕНИЯ УСТРОЙСТВА ПО USB (Android) ---
func isDeviceConnected() bool {
	cmd := exec.Command("adb", "devices")
	output, err := cmd.Output()
	if err != nil {
		log.Println("adb не найден или не установлен. Проверьте наличие adb в PATH.")
		return false
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of devices attached") {
			continue
		}
		// device	device, device	offline, device	unauthorized
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "device" {
			return true
		}
		if len(fields) >= 2 && (fields[1] == "unauthorized" || fields[1] == "offline") {
			log.Printf("adb: устройство найдено, но статус: %s. Разрешите отладку по USB и подтвердите ключ на устройстве.", fields[1])
		}
	}
	return false
}

func handleRecommendations(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "user_id is required"})
		return
	}
	recs := recommendations.GenerateRecommendations(userID)
	c.JSON(200, gin.H{"recommendations": recs})
}

func main() {
	initDatabaseConnections()
	defer PostgresPool.Close()
	defer Neo4jDriver.Close()

	database.InitDatabase()
	defer database.Db.Close()

	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
	})

	r.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": "dummy-token"})
	})

	r.GET("/posts", func(c *gin.Context) {
		rows, err := database.Db.Query("SELECT id, author, content, media, date FROM posts ORDER BY id DESC")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer rows.Close()
		var posts []Post
		for rows.Next() {
			var post Post
			if err := rows.Scan(&post.ID, &post.Author, &post.Content, &post.Media, &post.Date); err != nil {
				continue
			}
			posts = append(posts, post)
		}
		c.JSON(http.StatusOK, posts)
	})

	r.POST("/posts", func(c *gin.Context) {
		var post Post
		if err := c.ShouldBindJSON(&post); err != nil {
			log.Printf("Ошибка декодирования post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post data"})
			return
		}
		if post.Author == "" || post.Content == "" {
			log.Printf("Пустой author или content при публикации: author='%s', content='%s'", post.Author, post.Content)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Author and content required"})
			return
		}
		_, err := database.Db.Exec("INSERT INTO posts (author, content, media, date) VALUES (?, ?, ?, datetime('now'))", post.Author, post.Content, post.Media)
		if err != nil {
			log.Printf("Ошибка записи поста в БД: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		log.Printf("Пост успешно опубликован: author='%s', content='%s'", post.Author, post.Content)
		c.JSON(http.StatusCreated, gin.H{"message": "Post created"})
	})

	r.GET("/feed", func(c *gin.Context) {
		c.JSON(http.StatusOK, []Post{})
	})

	r.GET("/users/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter required"})
			return
		}
		c.JSON(http.StatusOK, []User{})
	})

	r.GET("/posts/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter required"})
			return
		}
		c.JSON(http.StatusOK, []Post{})
	})

	r.GET("/subscriptions", func(c *gin.Context) {
		subscriptions.HandleSubscriptions(c.Writer, c.Request)
	})

	r.GET("/recommendations", handleRecommendations)

	log.Println("Server is running on port 8080")
	log.Printf("Server is running on http://0.0.0.0:8080")
	r.Run()
}
