package cmd

import (
	"fmt"
	"io"

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
		columnFmt := color.New(color.FgGreen).SprintfFunc()
		count := 1
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			fmt.Println("\n\n***************** Row : ", count, " *****************")
			fmt.Println("\nPod Name : ", res.PodDetail)
			fmt.Println("\nNamespace : ", res.Namespace)
			//Print List of Processes
			fmt.Println("\nList of Processes (" + fmt.Sprint(len(res.ListOfProcess)) + ") :\n")
			tbl := output.New("SOURCE", "DESTINATION", "COUNT", "STATUS")
			tbl.WithHeaderFormatter(headerFmt)
			tbl.WithLastColumnFormatter(columnFmt)
			for _, process := range res.ListOfProcess {
				for _, source := range process.ListOfDestination {
					tbl.AddRow(process.Source, source.Destination, source.Count, source.Status)
				}
			}
			tbl.Print()

			//Print List of File System
			fmt.Println("\nList of File-system accesses (" + fmt.Sprint(len(res.ListOfFile)) + ") :\n")
			tbl = output.New("SOURCE", "DESTINATION", "COUNT", "STATUS")
			tbl.WithHeaderFormatter(headerFmt).WithLastColumnFormatter(columnFmt)
			for _, file := range res.ListOfFile {
				for _, source := range file.ListOfDestination {
					tbl.AddRow(file.Source, source.Destination, source.Count, source.Status)
				}
			}
			tbl.Print()

			//Print List of Network Connection
			fmt.Println("\nList of Network connections (" + fmt.Sprint(len(res.ListOfNetwork)) + ") :\n")
			tbl = output.New("SOURCE", "Protocol", "COUNT", "STATUS")
			tbl.WithHeaderFormatter(headerFmt).WithLastColumnFormatter(columnFmt)
			for _, network := range res.ListOfNetwork {
				for _, source := range network.ListOfDestination {
					tbl.AddRow(network.Source, source.Destination, source.Count, source.Status)
				}
			}
			tbl.Print()

			//Print Ingress Connection
			fmt.Printf("\nIngress Connection :\n\n")
			tbl = output.New("VISIBILITY", "COUNT", "STATUS")
			tbl.WithHeaderFormatter(headerFmt).WithLastColumnFormatter(columnFmt)
			tbl.AddRow("Within Cluster", res.Ingress.In, "ALLOW")
			tbl.AddRow("Outside Cluster", res.Ingress.Out, "ALLOW")
			tbl.Print()

			//Print Egress Connection
			fmt.Printf("\nEgress Connection : \n\n")
			tbl = output.New("VISIBILITY", "COUNT", "STATUS")
			tbl.WithHeaderFormatter(headerFmt).WithLastColumnFormatter(columnFmt)
			tbl.AddRow("Within Cluster", res.Egress.In, "ALLOW")
			tbl.AddRow("Outside Cluster", res.Egress.Out, "ALLOW")
			tbl.Print()
			count++
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sumCmd)
}
