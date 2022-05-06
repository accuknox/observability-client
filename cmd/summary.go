package cmd

import (
	"fmt"
	"io"
	"strings"

	sum "github.com/accuknox/observability-client/summary"
	"github.com/accuknox/observability/src/proto/summary"
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

		// tbl := output.New("POD-DETAILS", "NAMESPACE", "LIST-OF-PROCESS", "LIST-OF-FILE", "LIST-OF-NETWORK","INGRESS", "EGRESS")
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
			fmt.Println("\nList of Processes (" + fmt.Sprint(len(res.ListOfProcess)) + ") :")
			for _, value := range res.ListOfProcess {
				fmt.Print("", value.Source)
				fmt.Println(" ->", strings.Join(value.Resource, ", "))
			}
			fmt.Println("\nList of File-system accesses (" + fmt.Sprint(len(res.ListOfFile)) + ") :")
			for _, value := range res.ListOfFile {
				fmt.Print("", value.Source)
				fmt.Println(" ->", strings.Join(value.Resource, ", "))
			}
			fmt.Println("\nList of Network connections (" + fmt.Sprint(len(res.ListOfNetwork)) + ") :")
			for _, value := range res.ListOfNetwork {
				fmt.Print("", value.Source)
				fmt.Println(" ->", strings.Join(value.Resource, ", "))
			}
			fmt.Println("\nIngress Connection : ")
			fmt.Println("Connections within domain : ", res.Ingress.In)
			fmt.Println("Connections outside cluster : ", res.Ingress.Out)
			fmt.Println("\nEgress Connection : ")
			fmt.Println("Connection requests from within domain : ", res.Egress.In)
			fmt.Println("Connection requests from outside cluster : ", res.Egress.Out)
			count++
			// tbl.AddRow(res.PodDetail, res.Namespace, res.ListOfProcess, res.ListOfFile, res.ListOfNetwork, res.Ingress, res.Egress)
		}
		// tbl.Print()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sumCmd)
}
