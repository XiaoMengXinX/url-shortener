package api

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

func CreateShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	url := r.PostFormValue("url")
	token := r.PostFormValue("token")

	if url == "" {
		_, _ = fmt.Fprintf(w, responseJson(resData{Error: "URL is required"}))
		return
	}

	if token == "" || !strRules(token) || len(token) > 50 || len(token) < 3 {
		token = randToken(5)
	}

	err := collection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&urlData{})
	for err == nil {
		token = randToken(5)
		err = collection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&urlData{})
	}

	result := re.FindAllStringSubmatch(url, -1)
	if result == nil {
		_, _ = fmt.Fprintf(w, responseJson(resData{Error: "Invalid URL"}))
		return
	}

	u := urlData{
		Token:     token,
		URL:       url,
		CreatedAt: time.Now(),
	}

	_, err = collection.InsertOne(context.TODO(), u)
	if err != nil {
		_, _ = fmt.Fprintln(w, responseJson(resData{Error: err.Error()}))
		return
	}

	_, _ = fmt.Fprintf(w, responseJson(resData{Token: u.Token}))
}

func randToken(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func strRules(str string) bool {
	match, _ := regexp.MatchString(`^[A-Za-z0-9_-]+$`, str)
	return match
}
