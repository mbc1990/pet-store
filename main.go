package main

import (
  "encoding/json"
  "net/http"
  "html/template"
  "log"
  "io/ioutil"
)

type AnimalImage struct {
  Location string
}

func main() {
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
  // img := AnimalImage{Location: "fakeurl"}
  t, _ := template.ParseFiles("home.html")
  t.Execute(w, "https://i.imgur.com/WQBMaDU.jpg")
}

// Response object for the next_image api endpoint
type NextImageResponse struct {
  ImageUrl string
}

// Return a new image for the user
func nextImage(w http.ResponseWriter, r *http.Request) {
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
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err)
  }

  var p PreferenceParams
  err = json.Unmarshal(body, &p)
  if err != nil {
    panic(err)
  }
  log.Println(p.Liked)
}
