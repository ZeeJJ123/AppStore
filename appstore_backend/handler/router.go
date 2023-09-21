package handler

import (
	"appstore/util"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var mySigningKey []byte

// InitRouter initializes and configures the router for handling HTTP requests.
func InitRouter(config *util.TokenInfo) http.Handler {
    mySigningKey = []byte(config.Secret)

    // Create JWT middleware for authentication and authorization.
    jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
        ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
            return []byte(mySigningKey), nil
        },
        SigningMethod: jwt.SigningMethodHS256,
    })
    
    // Create a new router instance.
    router := mux.NewRouter()
    
    // Define routes and their corresponding handlers.
    router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
    router.Handle("/checkout", jwtMiddleware.Handler(http.HandlerFunc(checkoutHandler))).Methods("POST")    
    router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET")
    router.Handle("/app/{id}", jwtMiddleware.Handler(http.HandlerFunc(deleteHandler))).Methods("DELETE")
    router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
    router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

    originsOk := handlers.AllowedOrigins([]string{"*"})
    headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

    return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}
