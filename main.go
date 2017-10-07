package main

import (
  "database/sql"
  "fmt"
  "encoding/json"
  "os"
  "net/http"
  "html/template"
  "log"
  "io/ioutil"
  _ "github.com/lib/pq"
)

// Global configuration contants
var Conf Configuration = Configuration{}

type AnimalImage struct {
  Location string
}

type Configuration struct {
  Host     string
  Port     int
  User     string
  Password string
  Dbname   string
}


func main() {
  // Get consts
  file, _ := os.Open("conf.json")
  decoder := json.NewDecoder(file)
  err := decoder.Decode(&Conf)
  if err != nil {
    fmt.Println("error:", err)
  }

  // Training page
  http.HandleFunc("/train/", home)

  // Next image
  http.HandleFunc("/next_image/", nextImage)

  // Preference event
  http.HandleFunc("/preference_event/", preferenceEvent)

  // Static files
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  // Server
  http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("home.html")
  // TODO: execute without param
  t.Execute(w, "https://i.imgur.com/WQBMaDU.jpg")
}

// Response object for the next_image api endpoint
type NextImageResponse struct {
  ImageUrl string
}

// Return a new image for the user
func nextImage(w http.ResponseWriter, r *http.Request) {
  // TODO: Query database for next image
  var nextImage string = "https://s-media-cache-ak0.pinimg.com/originals/86/09/a2/8609a24d00485b9d0463c24e2a4aec37.jpg"
  resp := NextImageResponse{nextImage}
  js, err := json.Marshal(resp)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

// Struct for decoding preference request parameters
type PreferenceParams struct {
  Liked string
}

// When a user likes or dislikes an image, we write the event to a database
func preferenceEvent(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Received prefernce event")
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err)
  }

  var p PreferenceParams
  err = json.Unmarshal(body, &p)
  if err != nil {
    panic(err)
  }

  // Connect to database
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    Conf.Host, Conf.Port, Conf.User, Conf.Password, Conf.Dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }
  fmt.Println("Successfully connected to pg")

  sqlStatement := `  
  INSERT INTO preference_events (user_id, image_id, liked)
  VALUES ($1, $2, $3)`
  _, err = db.Exec(sqlStatement, 1, 1, p.Liked) // TODO: get image id from client
  if err != nil {
    panic(err)
  }

  log.Println(p.Liked)
  // TODO: Write result to database
}
