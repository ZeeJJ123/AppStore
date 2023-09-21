package backend

import (
    "context"
    "appstore/util"
    "fmt"
    "appstore/constants"
    "github.com/olivere/elastic/v7"
)

var (
    ESBackend *ElasticsearchBackend
)

type ElasticsearchBackend struct {
    client *elastic.Client
}

// InitElasticsearchBackend initializes the Elasticsearch backend.
func InitElasticsearchBackend(config *util.ElasticsearchInfo) {
    // Create an Elasticsearch client instance.
    client, err := elastic.NewClient(
        elastic.SetURL(config.Address),
        elastic.SetBasicAuth(config.Username, config.Password))
    if err != nil {
        panic(err)
    }

    // Check if the app index exists.
    exists, err := client.IndexExists(constants.APP_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }
    
    // Create the app index and its mapping if it doesn't exist.
    if !exists {
        mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "title":      { "type": "text"},
                    "description":  { "type": "text" },
                    "price":      { "type": "keyword", "index": false },
                    "url":     { "type": "keyword", "index": false }
                }
            }
        }`
        _, err := client.CreateIndex(constants.APP_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }

    // Check if the user index exists.
    exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    // Create the user index and its mapping if it doesn't exist.
    if !exists {
        mapping := `{
                     "mappings": {
                         "properties": {
                            "username": {"type": "keyword"},
                            "password": {"type": "keyword"},
                            "age": {"type": "long", "index": false},
                            "gender": {"type": "keyword", "index": false}
                         }
                    }
                }`
        _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }

    // Print a success message for index creation.
    fmt.Println("Indexes are created.")

    // assign value to the Elasticsearch backend.
    ESBackend = &ElasticsearchBackend{client: client}
}

func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    // Perform a search operation in Elasticsearch using the provided query and index.
    searchResult, err := backend.client.Search().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())

    // If an error occurs during the search operation, return nil and the error.
    if err != nil {
        return nil, err
    }
    // Return the search result and nil error indicating a successful search operation.
    return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    // Perform an indexing operation using the Elasticsearch client.
    _, err := backend.client.Index().
        Index(index).
        Id(id).
        BodyJson(i).
        Do(context.Background())// Execute the indexing operation in the background.
    return err
}

func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
    _, err := backend.client.DeleteByQuery().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())

    return err
}