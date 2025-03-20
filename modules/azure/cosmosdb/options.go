package cosmosdb

import (
	"strconv"

	"github.com/testcontainers/testcontainers-go"
)

// WithPartitions sets the number of partitions for the emulator. Default is 10.
func WithPartitions(partitions int) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env["AZURE_COSMOS_EMULATOR_PARTITION_COUNT"] = strconv.Itoa(partitions)
		return nil
	}
}
