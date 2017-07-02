package main

import (
    "os"

    "fmt"
    "net/http"

    "github.com/urfave/negroni"
    "github.com/gorilla/mux"
    "github.com/tobyjsullivan/event-store.v3/events"
    "encoding/base64"
    "github.com/tobyjsullivan/event-store.v3/store"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "log"
    "encoding/json"
)

var (
    eventStore *store.Store
)

func init() {
    bucket := os.Getenv("S3_BUCKET")
    if bucket == "" {
        panic("S3_BUCKET must be set.")
    }
    region := os.Getenv("AWS_REGION")
    if os.Getenv("AWS_REGION") == "" {
        panic("AWS_REGION must be set.")
    }
    if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
        panic("AWS_ACCESS_KEY_ID must be set.")
    }
    if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
        panic("AWS_SECRET_ACCESS_KEY must be set.")
    }

    logger := log.New(os.Stdout, "[svc] ", 0)

    sess := session.Must(session.NewSession(
        &aws.Config{
            Credentials: credentials.NewEnvCredentials(),
            Region: aws.String(region),
            Logger: aws.LoggerFunc(logger.Println),
            LogLevel: aws.LogLevel(aws.LogOff),
        },
    ))
    s3svc := s3.New(sess)

    eventStore = store.NewS3Store(s3svc, bucket)
}

func main() {
    r := buildRoutes()

    n := negroni.New()
    n.UseHandler(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    n.Run(":" + port)
}

func buildRoutes() http.Handler {
    r := mux.NewRouter()
    r.HandleFunc("/", statusHandler).Methods("GET")
    r.HandleFunc("/events", addEventHandler).Methods("POST")

    return r
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "The service is online!\n")
}

type jsResponse struct {
    Data interface{} `json:"data,omitempty"`
    Error error `json:"error,omitempty"`
}

type addEventResponse struct {
    EventID string `json:"eventId"`
}

func addEventHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    f := r.Form
    prevIdHex := f.Get("previous")
    eventType := f.Get("type")
    dataBase64 := f.Get("data")

    if prevIdHex == "" {
        http.Error(w, "Previous event ID is required.", http.StatusBadRequest)
        return
    }
    if eventType == "" {
        http.Error(w, "Event type is required.", http.StatusBadRequest)
        return
    }
    if dataBase64 == "" {
        http.Error(w, "Data is required.", http.StatusBadRequest)
        return
    }

    prevId := events.NewEventID()
    err = prevId.Parse(prevIdHex)
    if err != nil {
        http.Error(w, "Error parsing previous ID: "+err.Error(), http.StatusBadRequest)
        return
    }

    data, err := base64.StdEncoding.DecodeString(dataBase64)
    if err != nil {
        http.Error(w, "Error parsing data: "+err.Error(), http.StatusBadRequest)
        return
    }

    e := &events.Event{
        PreviousEvent: prevId,
        Type: eventType,
        Data: data,
    }
    err = eventStore.Save(e)
    if err != nil {
        http.Error(w, "Error saving event: "+err.Error(), http.StatusInternalServerError)
        return
    }

    id := e.ID()

    out := jsResponse{
        Data: addEventResponse{
            EventID: id.String(),
        },
    }

    encoder := json.NewEncoder(w)
    err = encoder.Encode(&out)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    return
}
