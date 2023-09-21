package backend

// backend package and it involves interacting with Google Cloud Storage (GCS)

import (
   "context"
   "fmt"
   "io"
   "appstore/util"
   "cloud.google.com/go/storage"
)

var (
   GCSBackend *GoogleCloudStorageBackend
)

// GoogleCloudStorageBackend represents the Google Cloud Storage backend.
type GoogleCloudStorageBackend struct {
   client *storage.Client
   bucket string
}

// InitGCSBackend initializes the GCS backend.
func InitGCSBackend(config *util.GCSInfo) {
   // Create a new GCS client.
   client, err := storage.NewClient(context.Background())
   if err != nil {
       panic(err)
   }
   
   // Assign the GCS client and bucket to the GCSBackend instance.
   GCSBackend = &GoogleCloudStorageBackend{
       client: client,
       bucket: config.Bucket,
   }
}

// SaveToGCS saves data to Google Cloud Storage.
func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
   // Create a new context.
   ctx := context.Background()
   // Get the GCS object handle.
   object := backend.client.Bucket(backend.bucket).Object(objectName)
   // Create a writer for the object.
   wc := object.NewWriter(ctx)
   // Copy the content from the reader to the writer.
   if _, err := io.Copy(wc, r); err != nil {
       return "", err
   }
   // Close the writer.
   if err := wc.Close(); err != nil {
       return "", err
   }
   // Set the object's ACL to allow public read access.
   if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
       return "", err
   }
   // Get the object's attributes, including the media link.
   attrs, err := object.Attrs(ctx)
   if err != nil {
       return "", err
   }
   // Print the media link and return it along with nil error.
   fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
   return attrs.MediaLink, nil
}
