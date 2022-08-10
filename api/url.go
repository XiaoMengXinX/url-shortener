package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func init() {
	client, err := connectDB(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("short_url").Collection("urls")
}

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("url") == "" {
		if len(r.URL.Path) <= 1 {
			_, _ = fmt.Fprintf(w, responseJson(resData{Error: "Invalid short name or url"}))
			return
		}
		var u urlData
		err := collection.FindOne(context.TODO(), bson.M{"token": r.URL.Path[1:]}).Decode(&u)
		if err != nil {
			_, _ = fmt.Fprintf(w, responseJson(resData{Error: "Invalid short name"}))
			return
		}
		http.Redirect(w, r, u.URL, http.StatusMovedPermanently)
		return
	}

	url := r.PostFormValue("url")
	token := r.PostFormValue("token")

	if token == "" || !strRules(token) || len(token) > 15 || len(token) < 3 {
		token = randToken(5)
	}

	err := collection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&urlData{})
	for err == nil {
		token = randToken(5)
		err = collection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&urlData{})
	}

	result := re.FindAllStringSubmatch(url, -1)
	if result == nil {
		_, _ = fmt.Fprintf(w, responseJson(resData{Error: "Invalid url"}))
		return
	}

	u := urlData{
		Token:     token,
		URL:       r.PostFormValue("url"),
		CreatedAt: time.Now(),
	}

	_, err = collection.InsertOne(context.TODO(), u)
	if err != nil {
		_, _ = fmt.Fprintln(w, responseJson(resData{Error: err.Error()}))
		return
	}

	_, _ = fmt.Fprintf(w, responseJson(resData{Token: u.Token}))
}

func responseJson(r resData) string {
	b, _ := json.Marshal(r)
	return string(b)
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
	match, _ := regexp.MatchString(`^[A-Za-z0-9]+$`, str)
	return match
}
