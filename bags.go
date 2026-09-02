package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func bagsHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value(userIDKey).(int)

		if !ok {
			http.Error(w, "userId missing or wrong type", http.StatusInternalServerError)
			return
		}

		var bagRequest CreateBagRequest

		err := json.NewDecoder(r.Body).Decode(&bagRequest)
		if err != nil {
			http.Error(w, "bad json request", http.StatusBadRequest)
			return
		}

		var id int
		err = db.QueryRow(context.Background(), "INSERT INTO bags (user_id, roaster_name, bean_origin, roast_date) VALUES ($1,$2,$3,$4) RETURNING id", userId, bagRequest.RoasterName, bagRequest.BeanOrigin, bagRequest.RoastDate).Scan(&id)
		if err != nil {
			http.Error(w, "db connection error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":          id,
			"roasterName": bagRequest.RoasterName,
			"beanOrigin":  bagRequest.BeanOrigin,
			"roastDate":   bagRequest.RoastDate,
		})
	}
}
