package main

import (
	"encoding/json"
	"net/http"

	"tipatwitter/backend/database"

	"golang.org/x/crypto/bcrypt"
)

// --- User registration ---
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	_, err := database.Db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", req.Username, string(hash))
	if err != nil {
		http.Error(w, "User exists", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registered"})
}

// --- User login ---
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	row := database.Db.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", req.Username)
	var id int
	var hash string
	if err := row.Scan(&id, &hash); err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{ "user_id": id, "username": req.Username })
}

// --- Create & get posts ---
func HandlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req struct {
			AuthorID int    `json:"author_id"`
			Content  string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}
		_, err := database.Db.Exec("INSERT INTO posts (author_id, content) VALUES (?, ?)", req.AuthorID, req.Content)
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Post created"})
	case http.MethodGet:
		rows, err := database.Db.Query(`SELECT posts.id, users.username, posts.content, posts.created_at FROM posts JOIN users ON posts.author_id = users.id ORDER BY posts.created_at DESC LIMIT 50`)
		if err != nil {
			http.Error(w, "Failed to get posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var posts []map[string]interface{}
		for rows.Next() {
			var id int
			var username, content, createdAt string
			if err := rows.Scan(&id, &username, &content, &createdAt); err != nil {
				continue
			}
			posts = append(posts, map[string]interface{}{
				"id": id, "author": username, "content": content, "date": createdAt,
			})
		}
		json.NewEncoder(w).Encode(posts)
	}
}

// --- Feed (posts only from subscriptions) ---
func HandleFeed(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}
	rows, err := database.Db.Query(`SELECT posts.id, users.username, posts.content, posts.created_at FROM posts JOIN users ON posts.author_id = users.id WHERE posts.author_id IN (SELECT followee_id FROM subscriptions WHERE follower_id = ?) ORDER BY posts.created_at DESC LIMIT 50`, userID)
	if err != nil {
		http.Error(w, "Failed to get feed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var username, content, createdAt string
		if err := rows.Scan(&id, &username, &content, &createdAt); err != nil {
			continue
		}
		posts = append(posts, map[string]interface{}{
			"id": id, "author": username, "content": content, "date": createdAt,
		})
	}
	json.NewEncoder(w).Encode(posts)
}

// --- User search ---
func HandleUserSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	rows, err := database.Db.Query("SELECT id, username FROM users WHERE username LIKE ? LIMIT 20", "%"+q+"%")
	if err != nil {
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username string
		if err := rows.Scan(&id, &username); err != nil {
			continue
		}
		users = append(users, map[string]interface{}{"id": id, "username": username})
	}
	json.NewEncoder(w).Encode(users)
}

// --- Post search ---
func HandlePostSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	rows, err := database.Db.Query(`SELECT posts.id, users.username, posts.content, posts.created_at FROM posts JOIN users ON posts.author_id = users.id WHERE posts.content LIKE ? ORDER BY posts.created_at DESC LIMIT 20`, "%"+q+"%")
	if err != nil {
		http.Error(w, "Failed to search posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var username, content, createdAt string
		if err := rows.Scan(&id, &username, &content, &createdAt); err != nil {
			continue
		}
		posts = append(posts, map[string]interface{}{
			"id": id, "author": username, "content": content, "date": createdAt,
		})
	}
	json.NewEncoder(w).Encode(posts)
}
