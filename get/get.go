package get

import (
	"context"
	"os"

	"github.com/accuknox/observability/src/proto/aggregator"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Options struct {
	GRPC      string
	MsgPath   string
	LogPath   string
	LogFilter string
	JSON      bool
}

func connectClient() (aggregator.AggregatorClient, error) {
	gRPC := "localhost:9089"

	if val, ok := os.LookupEnv("KUBEARMOR_SERVICE"); ok {
		gRPC = val
	}
	var client aggregator.AggregatorClient
	// fmt.Println("gRPC server: " + gRPC)
	connection, err := grpc.Dial(gRPC, grpc.WithInsecure())
	if err != nil {
		log.Error().Msg("Error while connecting to grpc " + err.Error())
		return client, err
	}
	client = aggregator.NewAggregatorClient(connection)
	return client, nil
}

//GetSystemLogs - Fetch system logs
func GetSystemLogs(options aggregator.SystemLogsRequest) (aggregator.Aggregator_FetchSystemLogsClient, error) {

	//Connect grpc client
	client, err := connectClient()
	if err != nil {
		return nil, err
	}
	//Fetch System Logs
	response, err := client.FetchSystemLogs(context.Background(), &options)
	if err != nil {
		return nil, err
	}
	return response, nil
}

//GetNetworkLogs - Fetch network logs
func GetNetworkLogs(options aggregator.NetworkLogsRequest) (aggregator.Aggregator_FetchNetworkLogsClient, error) {

	//Connect grpc client
	client, err := connectClient()
	if err != nil {
		return nil, err
	}
	//Fetch Network Logs
	response, err := client.FetchNetworkLogs(context.Background(), &options)
	if err != nil {
		return nil, err
	}
	return response, nil
}
