package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/accuknox/observability-client/output"
	sum "github.com/accuknox/observability-client/summary"
	"github.com/accuknox/observability/src/proto/summary"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var options summary.LogsRequest

var sumCmd = &cobra.Command{
	Use:   "summary",
	Short: "To Get the Summary log on Pod Level",
	Long:  `To Get the Summary log on Network(Hubble Relay) and System(Kubearmor Relay) together in Pod Level.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		options.Namespace = "default"
		n := len(args)
		if n > 0 {
			options.Namespace = args[0]
		}
		if n > 1 {
			options.Label = args[1]
		}
		stream, err := sum.GetSummaryLogs(options)
		if err != nil {
			return err
		}
		headerFmt := color.New(color.Underline).SprintfFunc()
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			fmt.Println("\n\n**********************************************************************")
			fmt.Println("\nPod Name : ", res.PodDetail)
			fmt.Println("\nNamespace : ", res.Namespace)
			//Print List of Processes
			fmt.Println("\nList of Processes (" + fmt.Sprint(len(res.ListOfProcess)) + ") :\n")
			tbl := output.New("SOURCE", "DESTINATION", "COUNT", "LAST UPDATED TIME", "STATUS")
			tbl.WithHeaderFormatter(headerFmt)
			for _, process := range res.ListOfProcess {
				for _, source := range process.ListOfDestination {
					tbl.AddRow(process.Source, source.Destination, source.Count, time.Unix(source.LastUpdatedTime, 0).Format("1-02-2006 15:04:05"), source.Status)
				}
			}
			tbl.Print()

			//Print List of File System
			fmt.Println("\nList of File-system accesses (" + fmt.Sprint(len(res.ListOfFile)) + ") :\n")
			tbl = output.New("SOURCE", "DESTINATION", "COUNT", "LAST UPDATED TIME", "STATUS")
			tbl.WithHeaderFormatter(headerFmt)
			for _, file := range res.ListOfFile {
				for _, source := range file.ListOfDestination {
					tbl.AddRow(file.Source, source.Destination, source.Count, time.Unix(source.LastUpdatedTime, 0).Format("1-02-2006 15:04:05"), source.Status)
				}
			}
			tbl.Print()

			//Print List of Network Connection
			fmt.Println("\nList of Network connections (" + fmt.Sprint(len(res.ListOfNetwork)) + ") :\n")
			tbl = output.New("SOURCE", "Protocol", "COUNT", "LAST UPDATED TIME", "STATUS")
			tbl.WithHeaderFormatter(headerFmt)
			for _, network := range res.ListOfNetwork {
				for _, source := range network.ListOfDestination {
					tbl.AddRow(network.Source, source.Destination, source.Count, time.Unix(source.LastUpdatedTime, 0).Format("1-02-2006 15:04:05"), source.Status)
				}
			}
			tbl.Print()

			//Print Ingress Connections
			fmt.Printf("\nIngress Connections :\n\n")
			tbl = output.New("DESTINATION LABEL", "DESTINATION NAMESPACE", "PROTOCOL", "PORT", "COUNT", "LAST UPDATED TIME", "STATUS")
			tbl.WithHeaderFormatter(headerFmt)
			for _, ingress := range res.Ingress {
				tbl.AddRow(ingress.DestinationLabels, ingress.DestinationNamespace, ingress.Protocol, ingress.Port, ingress.Count, time.Unix(ingress.LastUpdatedTime, 0).Format("1-02-2006 15:04:05"), ingress.Status)
			}
			tbl.Print()

			//Print Egress Connections
			fmt.Printf("\nEgress Connections : \n\n")
			tbl = output.New("DESTINATION LABEL", "DESTINATION NAMESPACE", "PROTOCOL", "PORT", "COUNT", "LAST UPDATED TIME", "STATUS")
			tbl.WithHeaderFormatter(headerFmt)
			for _, egress := range res.Egress {
				tbl.AddRow(egress.DestinationLabels, egress.DestinationNamespace, egress.Protocol, egress.Port, egress.Count, time.Unix(egress.LastUpdatedTime, 0).Format("1-02-2006 15:04:05"), egress.Status)
			}
			tbl.Print()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sumCmd)
}
