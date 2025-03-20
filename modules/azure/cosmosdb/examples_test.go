package cosmosdb_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/azure/cosmosdb"
)

func ExampleRun() {
	ctx := context.Background()

	cosmosDBCtr, err := cosmosdb.Run(
		ctx, "mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:latest",
		cosmosdb.WithPartitions(4),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(cosmosDBCtr); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	state, err := cosmosDBCtr.State(ctx)
	if err != nil {
		log.Printf("failed to get container state: %s", err)
		return
	}

	fmt.Println(state.Running)
	fmt.Println(len(cosmosDBCtr.Certificate()) > 0)

	// Output:
	// true
	// true
}

func ExampleRun_useClient() {
	ctx := context.Background()

	cosmosDBCtr, err := cosmosdb.Run(
		ctx, "mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:vnext-preview",
		cosmosdb.WithPartitions(4),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(cosmosDBCtr); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	// getCosmosDBCLient {
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cosmosDBCtr.Certificate()); !ok {
		log.Printf("failed to append certificate to pool")
		return
	}

	// create a client with the certificate provided by the container
	clientOptions := azcosmos.ClientOptions{
		EnableContentResponseOnWrite: true,
		ClientOptions: azcore.ClientOptions{
			Transport: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs:            certPool,
						InsecureSkipVerify: true,
					},
				},
			},
		},
	}

	emulatorURL, err := cosmosDBCtr.EmulatorURL(ctx)
	if err != nil {
		log.Printf("failed to get emulator URL: %s", err)
		return
	}

	// The emulator uses a well-known master key
	credential, err := azcosmos.NewKeyCredential(cosmosdb.EmulatorCredentials)
	if err != nil {
		log.Printf("failed to create credential: %s", err)
		return
	}

	cosmosClient, err := azcosmos.NewClientWithKey(emulatorURL, credential, &clientOptions)
	if err != nil {
		log.Printf("failed to create client: %s", err)
		return
	}
	//}

	// Create database if it doesn't exist
	databaseProperties := azcosmos.DatabaseProperties{ID: "cosmicworks"}
	dbResponse, err := cosmosClient.CreateDatabase(ctx, databaseProperties, &azcosmos.CreateDatabaseOptions{})
	if err != nil {
		log.Printf("failed to create database: %s", err)
		return
	}

	databaseClient, err := cosmosClient.NewDatabase(dbResponse.DatabaseProperties.ID)
	if err != nil {
		log.Printf("failed to create database: %s", err)
		return
	}

	// Create container if it doesn't exist
	containerProperties := azcosmos.ContainerProperties{
		ID: "products",
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/category"},
		},
	}
	ctrResponse, err := databaseClient.CreateContainer(ctx, containerProperties, nil)
	if err != nil {
		log.Printf("failed to create container: %s", err)
		return
	}

	// Get reference to container
	containerClient, err := databaseClient.NewContainer(ctrResponse.ContainerProperties.ID)
	if err != nil {
		log.Printf("failed to create container: %s", err)
		return
	}

	// itemDefinition{
	type Item struct {
		Id        string  `json:"id"`
		Category  string  `json:"category"`
		Name      string  `json:"name"`
		Quantity  int     `json:"quantity"`
		Price     float32 `json:"price"`
		Clearance bool    `json:"clearance"`
	}
	//}

	// createItem {
	itemId := "aaaaaaaa-0000-1111-2222-bbbbbbbbbbbb"

	item := Item{
		Id:        itemId,
		Category:  "gear-surf-surfboards",
		Name:      "Yamba Surfboard",
		Quantity:  12,
		Price:     850.00,
		Clearance: false,
	}

	partitionKey := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	// Create a context with timeout for database operations
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	bytes, err := json.Marshal(item)
	if err != nil {
		log.Printf("failed to marshal item: %s", err)
		return
	}

	response, err := containerClient.UpsertItem(dbCtx, partitionKey, bytes, nil)
	if err != nil {
		log.Printf("failed to upsert item: %s", err)
		return
	}
	//}

	// readItem {
	response, err = containerClient.ReadItem(dbCtx, partitionKey, itemId, nil)
	if err != nil {
		log.Printf("failed to read item: %s", err)
		return
	}

	if response.RawResponse.StatusCode == 200 {
		read_item := Item{}
		err := json.Unmarshal(response.Value, &read_item)
		if err != nil {
			log.Printf("failed to unmarshal item: %s", err)
			return
		}
	}
	// }

	// Output:
	// true
	// true
}
