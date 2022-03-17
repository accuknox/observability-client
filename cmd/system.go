package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/accuknox/observability-client/get"
	"github.com/accuknox/observability/src/proto/aggregator"
	"github.com/spf13/cobra"
)

var sysFilters aggregator.SystemLogsRequest
var limit int

//sysCmd represents the system command
var sysCmd = &cobra.Command{
	Use:   "system",
	Short: "System commends for get aggregated kubearmor logs",
	Long:  `System commands give all the aggregated kubearmor logs for observability. `,
	RunE: func(cm *cobra.Command, args []string) error {
		//Fetch logs
		sysFilters.Limit = int64(limit)
		response, err := get.GetSystemLogs(sysFilters)
		if err != nil {
			return err
		}

		//Check Logs exist
		if len(response.Logs) == 0 && response.Count == 0 {
			fmt.Println("No Log Found")
			return nil
		}
		//Show count of logs
		if sysFilters.Count {
			fmt.Println("Total Unique Logs : ", response.Count)
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
}
