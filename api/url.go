package api

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) <= 1 {
		_, _ = fmt.Fprintf(w, responseJson(resData{Error: "Invalid short name"}))
		return
	}

	var u urlData
	token := r.URL.Path[1:]
	err := collection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&u)
	if err != nil {
		_, _ = fmt.Fprintf(w, responseJson(resData{Error: "Invalid short name"}))
		return
	}

	http.Redirect(w, r, u.URL, http.StatusMovedPermanently)
}
