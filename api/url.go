package api

import (
	"context"
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
	ShortName string    `bson:"short"`
	URL       string    `bson:"url"`
	CreatedAt time.Time `bson:"created_at"`
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

func init() {
	client, err := connectDB(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("test").Collection("urls")
}

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("short_name") == "" || r.PostFormValue("url") == "" {
		fmt.Fprintf(w, r.RequestURI)
		return

		if len(r.RequestURI) <= 1 {
			_, _ = fmt.Fprintf(w, "Invaid short name or url")
			return
		}
		var u urlData
		err := collection.FindOne(context.TODO(), bson.M{"short": r.RequestURI[1:]}).Decode(&u)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Invaid short name")
			return
		}
		http.Redirect(w, r, u.URL, http.StatusMovedPermanently)
		return
	}

	url := r.PostFormValue("url")

	result := re.FindAllStringSubmatch(url, -1)
	if result == nil {
		_, _ = fmt.Fprintf(w, "Invaid url")
		return
	}

	u := urlData{
		ShortName: r.PostFormValue("short_name"),
		URL:       r.PostFormValue("url"),
		CreatedAt: time.Now(),
	}

	insertResult, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		_, _ = fmt.Fprintln(w, err)
		return
	}

	_, _ = fmt.Fprintf(w, "Success! %s", insertResult.InsertedID)
}
