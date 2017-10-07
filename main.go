package main

import (
  "database/sql"
  "fmt"
  "encoding/json"
  "os"
  "time"
  "net/http"
  "html/template"
  "strconv"
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
  ImageId string
}

// Return a new image for the user
func nextImage(w http.ResponseWriter, r *http.Request) {

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

  // TODO: Join against preference_events to only get unseen images
  sqlStatement := `  
  SELECT * FROM images LIMIT 1
  `
  rows, err := db.Query(sqlStatement);
  if err != nil {
    panic(err)
  }

  var nextImage string
  var imageIdStr string;
  for rows.Next() {
    var imageId int
    var original_url string
    var timestamp time.Time
    err = rows.Scan(&imageId, &nextImage, &original_url, &timestamp)
    if err != nil {
      panic(err)
    }
    imageIdStr = strconv.Itoa(imageId)
    if err != nil {
      panic(err)
    }
  }

  resp := NextImageResponse{nextImage, imageIdStr}
  js, err := json.Marshal(resp)
  fmt.Println(resp)
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
}
