package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

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
	RunE: func(cmd *cobra.Command, args []string) error {

		//Fetch logs
		netFilters.Limit = int64(limit)
		stream, err := get.GetNetworkLogs(netFilters)
		if err != nil {
			return err
		}
		var networkLogs []aggregator.NetworkLog
		var count int64
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if res.Logs != nil {
				networkLogs = append(networkLogs, *res.Logs)
			} else {
				count = res.Count
				break
			}
		}
		//Show count of logs
		if netFilters.Count {
			fmt.Println("Total Unique Logs : ", count)
			return nil
		}
		//Check Logs exist
		if len(networkLogs) == 0 {
			fmt.Println("No Log Found")
			return nil
		}
		//Convert logs into Json file
		if export.json {
			file, _ := json.MarshalIndent(networkLogs, "", "")
			fileName := "network_log_" + time.Now().UTC().Format("2006-01-02 15:04:05") + ".json"
			_ = ioutil.WriteFile(fileName, file, 0644)
			return nil
		}
		//Convert logs into CSV file
		if export.csv {
			fileName := "network_log_" + time.Now().UTC().Format("2006-01-02 15:04:05") + ".csv"
			file, _ := os.Create(fileName)
			writer := csv.NewWriter(file)
			defer writer.Flush()

			header := []string{"Verdict", "IpSource", "IpDestination", "IpVersion", "IpEncrypted",
				"L4TcpSourcePort", "L4TcpDestinationPort", "L4UdpSourcePort", "L4UdpDestinationPort",
				"L4Icmpv4Type", "L4Icmpv4Code", "L4Icmpv6Type", "L4Icmpv6Code",
				"SourceNamespace", "SourceLabels", "SourcePodName",
				"DestinationNamespace", "DestinationLabels", "DestinationPodName",
				"Type", "NodeName", "L7Type", "L7DnsCnames", "L7DnsObservationSource",
				"L7HttpCode", "L7HttpMethod", "L7HttpUrl", "L7HttpProtocol", "L7HttpHeaders",
				"EventTypeType", "EventTypeSubType", "SourceServiceName", "SourceServiceNamespace",
				"DestinationServiceName", "DestinationServiceNamespace", "TrafficDirection", "TraceObservationPoint",
				"DropReasonDesc", "IsReply", "StartTime", "LastUpdatedTime", "Total"}
			_ = writer.Write(header)

			for _, logs := range networkLogs {
				var row []string

				row = append(row, logs.Verdict, logs.IpSource, logs.IpDestination, logs.IpVersion, fmt.Sprint(logs.IpEncrypted),
					fmt.Sprint(logs.L4TcpSourcePort), fmt.Sprint(logs.L4TcpDestinationPort), fmt.Sprint(logs.L4UdpSourcePort), fmt.Sprint(logs.L4UdpDestinationPort),
					fmt.Sprint(logs.L4Icmpv4Type), fmt.Sprint(logs.L4Icmpv4Code), fmt.Sprint(logs.L4Icmpv6Type), fmt.Sprint(logs.L4Icmpv6Code),
					logs.SourceNamespace, logs.SourceLabels, logs.SourcePodName,
					logs.DestinationNamespace, logs.DestinationLabels, logs.DestinationPodName,
					logs.Type, logs.NodeName, logs.L7Type, logs.L7DnsCnames, logs.L7DnsObservationSource,
					fmt.Sprint(logs.L7HttpCode), logs.L7HttpMethod, logs.L7HttpUrl, logs.L7HttpProtocol, logs.L7HttpHeaders,
					fmt.Sprint(logs.EventTypeType), fmt.Sprint(logs.EventTypeSubType), logs.SourceServiceName, logs.SourceServiceNamespace,
					logs.DestinationServiceName, logs.DestinationServiceNamespace, logs.TrafficDirection, logs.TraceObservationPoint,
					logs.DropReasonDesc, fmt.Sprint(logs.StartTime), fmt.Sprint(logs.UpdatedTime), fmt.Sprint(logs.Total))
				_ = writer.Write(row)
			}
			return nil
		}

		for _, log := range networkLogs {
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
	netCmd.Flags().StringArrayVarP(&netFilters.Verdict, "verdict", "v", nil, "fetch logs based on Verdict {Forwarded/Dropped/Error/Audit}")
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
	netCmd.Flags().BoolVar(&export.json, "json", false, "export file in Json format")
	netCmd.Flags().BoolVar(&export.csv, "csv", false, "export file in CSV format")
}
