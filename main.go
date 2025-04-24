
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/gorilla/mux"
    "github.com/jackc/pgx/v4"
    "golang.org/x/net/context"
)

var db *pgx.Conn
var rdb *redis.Client

type Client struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Slug        string    `json:"slug"`
    IsProject   bool      `json:"is_project"`
    SelfCapture bool      `json:"self_capture"`
    ClientLogo  string    `json:"client_logo"`
    Address     string    `json:"address"`
    PhoneNumber string    `json:"phone_number"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   time.Time `json:"deleted_at"`
}

func init() {
    var err error
    // Connect to PostgreSQL
    db, err = pgx.Connect(context.Background(), "postgres://username:password@localhost:5432/mydb")
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }

    // Connect to Redis
    rdb = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
}

func CreateClient(w http.ResponseWriter, r *http.Request) {
    var client Client
    err := json.NewDecoder(r.Body).Decode(&client)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.Exec(context.Background(), "INSERT INTO my_client (name, slug, is_project, self_capture, client_logo) VALUES ($1, $2, $3, $4, $5)", client.Name, client.Slug, client.IsProject, client.SelfCapture, client.ClientLogo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    clientJSON, _ := json.Marshal(client)
    rdb.Set(context.Background(), client.Slug, clientJSON, 0)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(client)
}

func GetClient(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    slug := vars["slug"]

    clientJSON, err := rdb.Get(context.Background(), slug).Result()
    if err == redis.Nil {
        var client Client
        err = db.QueryRow(context.Background(), "SELECT * FROM my_client WHERE slug = $1", slug).Scan(&client.ID, &client.Name, &client.Slug, &client.IsProject, &client.SelfCapture, &client.ClientLogo, &client.Address, &client.PhoneNumber, &client.CreatedAt, &client.UpdatedAt, &client.DeletedAt)
        if err != nil {
            http.Error(w, "Client not found", http.StatusNotFound)
            return
        }

        clientJSON, _ = json.Marshal(client)
        rdb.Set(context.Background(), slug, clientJSON, 0)
    }

    var client Client
    err = json.Unmarshal([]byte(clientJSON), &client)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(client)
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/client", CreateClient).Methods("POST")
    r.HandleFunc("/client/{slug}", GetClient).Methods("GET")

    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
    