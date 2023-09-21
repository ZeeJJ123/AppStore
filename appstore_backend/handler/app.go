package handler

import (
	"appstore/model"
	"appstore/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

// uploadHandler handles the HTTP request for uploading an app.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse from body of request to get a json object.
    fmt.Println("Received one upload request")
    
    // Get the username from the JWT token.
    token := r.Context().Value("user")
    claims := token.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"] 
    
    // Create an App object with the provided data.
    app := model.App{
        Id:          uuid.New(),
        User:        username.(string),
        Title:       r.FormValue("title"),
        Description: r.FormValue("description"),
    }
    
    // Parse and assign the price value.
    price, err := strconv.Atoi(r.FormValue("price"))
    fmt.Printf("%v,%T", price, price)
    if err != nil {
        fmt.Println(err)
    }
    app.Price = price
    
     // Retrieve the media file from the request.
    file, _, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }

    // Save the app and media file to the backend.
    err = service.SaveApp(&app, file)
    if err != nil {
        http.Error(w, "Failed to save app to backend", http.StatusInternalServerError)
        fmt.Printf("Failed to save app to backend %v\n", err)
        return
    }

    fmt.Println("App is saved successfully.")
 
    // return
    fmt.Fprintf(w, "App is saved successfully: %s\n", app.Description)
}

// searchHandler handles the HTTP request for searching apps.
func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one search request")
    w.Header().Set("Content-Type", "application/json")

    // Extract the search parameters from the request query.
    title := r.URL.Query().Get("title")
    description := r.URL.Query().Get("description")
 
    // Search for apps based on the provided parameters.
    var apps []model.App
    var err error
    apps, err = service.SearchApps(title, description)
    if err != nil {
        http.Error(w, "Failed to read Apps from backend", http.StatusInternalServerError)
        return
    }
 
    // Convert the apps to JSON format and write the response.
    js, err := json.Marshal(apps)
    if err != nil {
        http.Error(w, "Failed to parse Apps into JSON format", http.StatusInternalServerError)
        return
    }
    w.Write(js)
 }

// checkoutHandler handles the HTTP request for the app checkout process.
func checkoutHandler(w http.ResponseWriter, r *http.Request) {
   fmt.Println("Received one checkout request")
   w.Header().Set("Content-Type", "text/plain")
   
   // Extract the app ID from the request.
   appID := r.FormValue("appID")

   // Start the checkout process and get the checkout URL.
   url, err := service.CheckoutApp(r.Header.Get("Origin"), appID)
   if err != nil {
       fmt.Println("Checkout failed.")
       w.Write([]byte(err.Error()))
       return
   }

   w.WriteHeader(http.StatusOK)
   w.Write([]byte(url))

   fmt.Println("Checkout process started!")
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for delete")

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"].(string)
    id := mux.Vars(r)["id"]

    if err := service.DeleteApp(id, username); err != nil {
        http.Error(w, "Failed to delete app from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to delete app from backend %v\n", err)
        return
    }
    fmt.Println("App is deleted successfully")
}


 