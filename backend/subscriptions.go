package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tipatwitter/backend/database"
)

type Subscription struct {
	FollowerID int `json:"follower_id"`
	FolloweeID int `json:"followee_id"`
}

func HandleSubscriptions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var sub Subscription
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		_, err := database.Db.Exec("INSERT OR IGNORE INTO subscriptions (follower_id, followee_id) VALUES (?, ?)", sub.FollowerID, sub.FolloweeID)
		if err != nil {
			http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Subscribed"})
	case http.MethodDelete:
		followerID, _ := strconv.Atoi(r.URL.Query().Get("follower_id"))
		followeeID, _ := strconv.Atoi(r.URL.Query().Get("followee_id"))
		_, err := database.Db.Exec("DELETE FROM subscriptions WHERE follower_id = ? AND followee_id = ?", followerID, followeeID)
		if err != nil {
			http.Error(w, "Failed to unsubscribe", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unsubscribed"})
	case http.MethodGet:
		userID := r.URL.Query().Get("user_id")
		rows, err := database.Db.Query("SELECT followee_id FROM subscriptions WHERE follower_id = ?", userID)
		if err != nil {
			http.Error(w, "Failed to get subscriptions", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var ids []int
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err == nil {
				ids = append(ids, id)
			}
		}
		json.NewEncoder(w).Encode(ids)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}