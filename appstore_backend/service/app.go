package service

import (
	"appstore/backend"
	"appstore/constants"
	"appstore/gateway/stripe"
	"appstore/model"
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"github.com/olivere/elastic/v7"
)

// SearchApps searches for apps based on the title and description.
// If the title is empty, it searches by description.
// If the description is empty, it searches by title.
func SearchApps(title string, description string) ([]model.App, error) {
   if title == "" {
       return SearchAppsByDescription(description)
   }
   if description == "" {
       return SearchAppsByTitle(title)
   }


   query1 := elastic.NewMatchQuery("title", title)
   query2 := elastic.NewMatchQuery("description", description)
   query := elastic.NewBoolQuery().Must(query1, query2)
   searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
   if err != nil {
       return nil, err
   }


   return getAppFromSearchResult(searchResult), nil
}

// SearchAppsByTitle searches for apps by title.
func SearchAppsByTitle(title string) ([]model.App, error) {
   query := elastic.NewMatchQuery("title", title)
   query.Operator("AND")
   if title == "" {
       query.ZeroTermsQuery("all")
   }
   searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
   if err != nil {
       return nil, err
   }

   return getAppFromSearchResult(searchResult), nil
}

// SearchAppsByDescription searches for apps by description.
func SearchAppsByDescription(description string) ([]model.App, error) {
   query := elastic.NewMatchQuery("description", description)
   query.Operator("AND")
   if description == "" {
       query.ZeroTermsQuery("all")
   }
   searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
   if err != nil {
       return nil, err
   }


   return getAppFromSearchResult(searchResult), nil
}

// SearchAppsByID searches for an app by its ID.
func SearchAppsByID(appID string) (*model.App, error) {
   // construct search query
   query := elastic.NewTermQuery("id", appID)

   // call backend
   searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
   if err != nil {
       return nil, err
   }

   // process result
   results := getAppFromSearchResult(searchResult)
   if len(results) == 1 {
       return &results[0], nil
   }
   
   // could optimize better error handling
   return nil, nil
}

// getAppFromSearchResult converts the search result to an array of model.App.
func getAppFromSearchResult(searchResult *elastic.SearchResult) []model.App {
   var ptype model.App
   var apps []model.App
   for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
       p := item.(model.App)
       apps = append(apps, p)
   }
   return apps
}

// SaveApp saves an app to the database, creates a product and price in Stripe, and saves the file to Google Cloud Storage.
func SaveApp(app *model.App, file multipart.File) error {
   // 优化 roll back or retry
   // Create a product and price using the Stripe SDK.
   productID, priceID, err := stripe.CreateProductWithPrice(app.Title, app.Description, int64(app.Price*100))
   if err != nil {
       fmt.Printf("Failed to create Product and Price using Stripe SDK %v\n", err)
       return err
   }
   app.ProductID = productID
   	   app.PriceID = priceID
   fmt.Printf("Product %s with price %s is successfully created", productID, priceID)
   
   // Save the file to Google Cloud Storage (GCS).
   medialink, err := backend.GCSBackend.SaveToGCS(file, app.Id)
   if err != nil {
       return err
   }
   app.Url = medialink
  
   // Save the app to Elasticsearch.
  err = backend.ESBackend.SaveToES(app, constants.APP_INDEX, app.Id)
   if err != nil {
       fmt.Printf("Failed to save app to elastic search with app index %v\n", err)
       return err
   }
   fmt.Println("App is saved successfully to ES app index.")

   return nil
}

func CheckoutApp(domain string, appID string) (string, error) {
    // appId to priceID
    // create a checkout seesion by passing in priceID
    app, err := SearchAppsByID(appID)
    if err != nil {
        return "", err
    }
    if app == nil {
        return "", errors.New("unable to find app in elasticsearch")
    }
    return stripe.CreateCheckoutSession(domain, app.PriceID)
 }

func DeleteApp(id string, user string) error {
    query := elastic.NewBoolQuery()
    query.Must(elastic.NewTermQuery("id", id))
    query.Must(elastic.NewTermQuery("user", user))

    return backend.ESBackend.DeleteFromES(query, constants.APP_INDEX)
}
 

