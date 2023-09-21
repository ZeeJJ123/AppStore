package main

import (
	"appstore/backend"
	"appstore/handler"
	"appstore/util"
	"fmt"
	"log"
	"net/http"
)

func main() {
    fmt.Println("started-service")
	config, err := util.LoadApplicationConfig("conf", "deploy.yml")
    if err != nil {
        panic(err)
    }

    backend.InitElasticsearchBackend(config.ElasticsearchConfig)
    backend.InitGCSBackend(config.GCSConfig)

	log.Fatal(http.ListenAndServe(":8080",       handler.InitRouter(config.TokenConfig)))
}