package subscriptions

import (
	"encoding/json"
	"net/http"

	"tipatwitter/backend/database"
)

type Subscription struct {
	UserID       int `json:"user_id"`
	SubscribedTo int `json:"subscribed_to"`
}

func HandleSubscriptions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var sub Subscription
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err := database.Db.Exec("INSERT INTO subscriptions (user_id, subscribed_to) VALUES (?, ?)", 
			sub.UserID, sub.SubscribedTo)
		if err != nil {
			http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Subscription created"})

	case http.MethodGet:
		rows, err := database.Db.Query("SELECT user_id, subscribed_to FROM subscriptions")
		if err != nil {
			http.Error(w, "Failed to get subscriptions", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var subscriptions []Subscription
		for rows.Next() {
			var sub Subscription
			if err := rows.Scan(&sub.UserID, &sub.SubscribedTo); err != nil {
				http.Error(w, "Failed to read subscription", http.StatusInternalServerError)
				return
			}
			subscriptions = append(subscriptions, sub)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(subscriptions)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}