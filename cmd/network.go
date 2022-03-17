package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/accuknox/observability-client/get"
	"github.com/accuknox/observability/src/proto/aggregator"
	"github.com/spf13/cobra"
)

var netFilters aggregator.NetworkLogsRequest

//netCmd represents the Network command
var netCmd = &cobra.Command{
	Use:   "network",
	Short: "Network commends for get aggregated cilium logs",
	Long:  `Network commands give all the aggregated cilium logs for observability. `,
	RunE: func(cm *cobra.Command, args []string) error {

		//Fetch logs
		netFilters.Limit = int64(limit)
		response, err := get.GetNetworkLogs(netFilters)
		if err != nil {
			return err
		}
		//Show count of logs
		if netFilters.Count {
			fmt.Println("Total Unique Logs : ", response.Count)
			return nil
		}
		//Check Logs exist
		if len(response.Logs) == 0 {
			fmt.Println("No Log Found")
			return nil
		}
		for _, log := range response.Logs {
			output, _ := json.Marshal(log)
			fmt.Println(string(output))
			fmt.Println()
		}

		return nil
	},
}

func init() {
	logCmd.AddCommand(netCmd)
	netCmd.Flags().IntVarP(&limit, "limit", "l", 0, "fetch limited logs")
	netCmd.Flags().BoolVar(&netFilters.Count, "count", false, "count number of unique logs")
	netCmd.Flags().StringVar(&netFilters.Direction, "direction", "", "fetch log based on Direction {Ingress/Egress}")
	netCmd.Flags().StringVar(&netFilters.Type, "type", "", "fetch log based on Type L3_L4/L7")
	netCmd.Flags().StringArrayVarP(&netFilters.Verdict, "verdict", "v", nil, "fetch logs based on Verdict {Forward/Dropped/Error/Audit}")
	netCmd.Flags().StringVar(&netFilters.Protocol, "protocol", "", "fetch log based on L4 Protocol {TCP/UDP/ICMPv4/ICMPv6}")
	netCmd.Flags().StringVar(&netFilters.L7, "l7", "", "fetch log based on L7 Protocol {DNS/Kafka/HTTP}")
	netCmd.Flags().StringArrayVar(&netFilters.SourcePod, "source-pod", nil, "fetch logs based on Source Pod Name")
	netCmd.Flags().StringArrayVar(&netFilters.SourceNamespace, "source-namespace", nil, "fetch logs based on Source Namespace")
	netCmd.Flags().StringArrayVar(&netFilters.DestinationPod, "destination-pod", nil, "fetch logs based on Destination Pod Name")
	netCmd.Flags().StringArrayVar(&netFilters.DestinationNamespace, "destination-namespace", nil, "fetch logs based on Destination Namespace")
	netCmd.Flags().StringArrayVar(&netFilters.Node, "node", nil, "fetch logs based on Node")
	netCmd.Flags().StringVar(&netFilters.SourceLabel, "source-label", "", "fetch log based on Source Label")
	netCmd.Flags().StringVar(&netFilters.DestinationLabel, "destination-label", "", "fetch log based on Destination Label")
	netCmd.Flags().StringVar(&netFilters.Since, "since", "", "fetch log based on time {1d/1h/1m/1s}")
}
