package api

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

func init() {
	client, err := connectDB(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("short_url").Collection("urls")
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
