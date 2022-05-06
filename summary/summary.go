package summary

import (
	"context"
	"os"

	"github.com/accuknox/observability/src/proto/summary"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

//connectClient - To connect summary Client
func connectClient() (summary.SummaryClient, error) {
	gRPC := "localhost:8089"

	if val, ok := os.LookupEnv("OBSERVABILITY_SERVICE"); ok {
		gRPC = val
	}
	var client summary.SummaryClient
	connection, err := grpc.Dial(gRPC, grpc.WithInsecure())
	if err != nil {
		log.Error().Msg("Error while connecting to grpc " + err.Error())
		return client, err
	}
	client = summary.NewSummaryClient(connection)
	return client, nil
}

func GetSummaryLogs(options summary.LogsRequest) (summary.Summary_FetchLogsClient, error) {
	//Connect grpc summary Client
	client, err := connectClient()
	if err != nil {
		return nil, err
	}
	//Fetch Summary Logs
	response, err := client.FetchLogs(context.Background(), &options)
	if err != nil {
		return nil, err
	}
	return response, nil
}
