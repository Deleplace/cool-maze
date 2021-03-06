package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	hotmaze "github.com/Deleplace/hot-maze/B3"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

const (
	projectID               = "hot-maze"
	backendBaseURL          = "https://hot-maze-udtcugbenq-uc.a.run.app"
	storageServiceAccountID = "ephemeral-storage@hot-maze.iam.gserviceaccount.com"
	bucket                  = "hot-maze.appspot.com"
	fileDeleteAfter         = 9 * time.Minute
)

func main() {
	ctx := context.Background()

	log.Println("GOOGLE_APPLICATION_CREDENTIALS =", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	// log.Println(os.Environ())
	cred, errCred := google.FindDefaultCredentials(ctx)
	log.Println("errCred =", errCred)
	if cred != nil {
		log.Println("FindDefaultCredentials =", string(cred.JSON))
	}

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		// log.Fatal("Couldn't create Storage client:", err)
		log.Println("Couldn't create Storage client:", err)
	}
	sa, errSA := storageClient.ServiceAccount(ctx, "hot-maze")
	log.Println("storageClient.ServiceAccount is", sa, errSA)

	storagePrivateKey, errSecret := hotmaze.AccessSecretVersion("projects/230384242501/secrets/B3-storage-private-key/versions/latest")
	if errSecret != nil {
		// log.Fatal("Couldn't read Storage service account private key:", err)
		log.Println("Couldn't read Storage service account private key:", errSecret)
	}

	server := hotmaze.Server{
		GCPProjectID:        projectID,
		BackendBaseURL:      backendBaseURL,
		StorageClient:       storageClient,
		StorageAccountID:    storageServiceAccountID,
		StoragePrivateKey:   storagePrivateKey,
		StorageBucket:       bucket,
		StorageFileTTL:      fileDeleteAfter,
		CloudTasksQueuePath: "projects/hot-maze/locations/us-central1/queues/b3-file-expiry",
	}
	server.RegisterHandlers()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatal(err)
}
