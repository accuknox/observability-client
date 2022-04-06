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

var sysFilters aggregator.SystemLogsRequest
var limit int
var export struct {
	json bool
	csv  bool
}

//sysCmd represents the system command
var sysCmd = &cobra.Command{
	Use:   "system",
	Short: "System commands for get aggregated kubearmor logs",
	Long:  `System commands give all the aggregated kubearmor logs for observability. `,
	RunE: func(cmd *cobra.Command, args []string) error {
		//Fetch logs
		sysFilters.Limit = int64(limit)
		stream, err := get.GetSystemLogs(sysFilters)
		if err != nil {
			return err
		}
		var systemLogs []aggregator.SystemLog
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
				systemLogs = append(systemLogs, *res.Logs)
			} else {
				count = res.Count
				break
			}
		}
		//Show count of logs
		if sysFilters.Count {
			fmt.Println("Total Unique Logs : ", count)
			return nil
		}
		//Check Logs exist
		if len(systemLogs) == 0 {
			fmt.Println("No Log Found")
			return nil
		}

		//Convert logs into Json file
		if export.json {
			file, _ := json.MarshalIndent(systemLogs, "", "")
			fileName := "system_log" + time.Now().UTC().Format("2006-01-02 15:04:05") + ".json"
			_ = ioutil.WriteFile(fileName, file, 0644)
			return nil
		}
		//Convert logs into CSV file
		if export.csv {
			fileName := "system_log_" + time.Now().UTC().Format("2006-01-02 15:04:05") + ".csv"
			file, _ := os.Create(fileName)
			writer := csv.NewWriter(file)
			defer writer.Flush()

			header := []string{"ClusterName", "HostName", "Namespace", "PodName", "ContainerID",
				"ContainerName", "Uid", "Type", "Source", "Operation", "Resource", "Data",
				"StartTime", "LastUpdatedTime", "Result", "Total"}
			_ = writer.Write(header)

			for _, logs := range systemLogs {
				var row []string
				row = append(row, logs.ClusterName, logs.HostName, logs.Namespace, logs.PodName, logs.ContainerID,
					logs.ContainerName, fmt.Sprint(logs.Uid), logs.Type, logs.Source, logs.Operation, logs.Resource, logs.Data,
					fmt.Sprint(logs.StartTime), fmt.Sprint(logs.UpdateTime), logs.Result, fmt.Sprint(logs.Total))
				_ = writer.Write(row)
			}
			return nil
		}
		for _, log := range systemLogs {
			output, _ := json.Marshal(log)
			fmt.Println(string(output))
			fmt.Println()
		}

		return nil
	},
}

func init() {
	logCmd.AddCommand(sysCmd)
	sysCmd.Flags().IntVarP(&limit, "limit", "l", 0, "fetch limited logs")
	sysCmd.Flags().BoolVar(&sysFilters.Count, "count", false, "count number of unique logs")
	sysCmd.Flags().StringArrayVarP(&sysFilters.Namespace, "namespace", "n", nil, "fetch logs based on Namespace")
	sysCmd.Flags().StringVar(&sysFilters.Type, "type", "", "fetch log based on type ContainerLog/HostLog")
	sysCmd.Flags().StringArrayVarP(&sysFilters.Operation, "operation", "o", nil, "fetch logs based on Operation Network/Process/File")
	sysCmd.Flags().StringArrayVarP(&sysFilters.Pod, "pod", "p", nil, "fetch logs based on Pod Name")
	sysCmd.Flags().StringArrayVar(&sysFilters.Host, "host", nil, "fetch logs based on Host Name")
	sysCmd.Flags().StringVarP(&sysFilters.Source, "source", "s", "", "fetch log based on source")
	sysCmd.Flags().StringVarP(&sysFilters.Resource, "resource", "r", "", "fetch log based on resource")
	sysCmd.Flags().StringArrayVarP(&sysFilters.Container, "container", "c", nil, "fetch logs based on container Name")
	sysCmd.Flags().StringVar(&sysFilters.Since, "since", "", "fetch log based on time {1d/1h/1m/1s}")

	// sysCmd.AddCommand(exportSys)
	sysCmd.Flags().BoolVar(&export.json, "json", false, "export file in Json format")
	sysCmd.Flags().BoolVar(&export.csv, "csv", false, "export file in CSV format")
}
