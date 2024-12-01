package api

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

var collection *mongo.Collection
var re = regexp.MustCompile("(http|https)://[\\w\\-_]+(\\.[\\w\\-_]+)+([\\w\\-.,@?^=%&:/~+#]*[\\w\\-@?^=%&/~+#])?")

type urlData struct {
	Token     string    `bson:"token"`
	URL       string    `bson:"url"`
	CreatedAt time.Time `bson:"created_at"`
}

type resData struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

func connectDB(uri string) (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return mongo.Connect(ctx, clientOptions)
}

func responseJson(r resData) string {
	b, _ := json.Marshal(r)
	return string(b)
}

func init() {
	client, err := connectDB(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("short_url").Collection("urls")
}

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
