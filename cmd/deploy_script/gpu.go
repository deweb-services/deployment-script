package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tcli "github.com/deweb-services/terraform-provider-nodeshift/nodeshift/provider/client"
	"github.com/spf13/cobra"
)

var (
	clientCfg    tcli.NodeshiftProviderConfiguration
	gpuCreateCfg tcli.GPUConfig
	gpuDeleteCfg string

	gpuCmd = &cobra.Command{
		Use:   "gpu",
		Short: "deploy GPU instance",
	}
	gpuDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete gpu instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			success := false
			log.Debugf("delete gpu instance with uuid %s", gpuDeleteCfg)
			cli := tcli.NewClient(cmd.Context(), clientCfg, tcli.ClientOptWithURL(APIURL))
			err := cli.DeleteGPU(cmd.Context(), gpuDeleteCfg)
			if err == nil {
				success = true
			}
			res := fmt.Sprintf("success=%t\n", success)
			_ = os.WriteFile("result", []byte(res), 0644)
			log.Debug(res)

			if err != nil {
				log.Errorf("delete gpu instance with uuid %s, error: %v", gpuDeleteCfg, err)
				return err
			}
			return nil
		},
	}

	gpuCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create gpu instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debugf("create gpu instance with config %#+v", gpuCreateCfg)
			cli := tcli.NewClient(cmd.Context(), clientCfg, tcli.ClientOptWithURL(APIURL))
			respCreate, err := cli.CreateGPU(cmd.Context(), &gpuCreateCfg)
			if err != nil {
				return fmt.Errorf("create gpu error: %w", err)
			}
			time.Sleep(time.Second * 30)
			respGet := &tcli.RentedGpuInfoResponse{}
		Loop:
			for i := 0; i < maxTries; i++ {
				respGet, err = cli.GetGPU(cmd.Context(), respCreate.UUID)
				if err == nil {
					switch strings.ToLower(respGet.ActualStatus) {
					case "running":
						break Loop
					case "destroying", "exited":
						err = fmt.Errorf("failed to create gpu")
						break Loop
					default:
						log.Debugf("get gpu instance status: %s", respGet.ActualStatus)
					}
				} else {
					log.Debugf("get gpu error: %s", err)
				}
				time.Sleep(sleepTime)
			}
			// write uuid to file anyway
			res := fmt.Sprintf("uuid=%s\nhost=%s\nport=%d\n", respCreate.UUID, respGet.SshHost, respGet.SshPort)
			_ = os.WriteFile("result", []byte(res), 0644)
			// then check the error
			if err != nil {
				return fmt.Errorf("get created gpu parameters error: %s", err)
			}
			return nil
		},
	}
)
