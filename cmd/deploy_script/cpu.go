package main

import (
	"fmt"
	"os"

	tcli "github.com/deweb-services/terraform-provider-dws/dws/provider/client"
	"github.com/spf13/cobra"
)

var (
	cpuCreateCfg tcli.DeploymentConfig
	cpuDeleteCfg string

	cpuCmd = &cobra.Command{
		Use:   "cpu",
		Short: "deploy CPU instance",
	}

	cpuDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete cpu instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			success := false
			log.Debugf("delete cpu instance with uuid %s", cpuDeleteCfg)
			cli := tcli.NewClient(cmd.Context(), clientCfg, tcli.ClientOptWithURL(APIURL))
			err := cli.DeleteDeployment(cmd.Context(), cpuDeleteCfg)
			if err == nil {
				success = true
			}
			res := fmt.Sprintf("success=%t\n", success)
			_ = os.WriteFile("result", []byte(res), 0644)
			log.Debug(res)

			if err != nil {
				log.Errorf("delete cpu instance with uuid %s, error: %v", cpuDeleteCfg, err)
				return err
			}
			return nil
		},
	}

	cpuCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create cpu instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debugf("create cpu instance with config %v", cpuCreateCfg)
			cli := tcli.NewClient(cmd.Context(), clientCfg, tcli.ClientOptWithURL(APIURL))
			resp, err := cli.CreateDeployment(cmd.Context(), &cpuCreateCfg)
			if err != nil {
				log.Debugf("get cpu error: %s", err)
				return err
			}
			res := fmt.Sprintf("id=%s\n", resp.ID)
			if resp.Data != nil {
				res += fmt.Sprintf("ip=%s\nipv6=%s\nygg=%s\n", resp.Data.IP, resp.Data.IPv6, resp.Data.Ygg)
			}
			if resp.EndTime != nil {
				res += fmt.Sprintf("host=%d\n", *resp.EndTime)
			}
			_ = os.WriteFile("result", []byte(res), 0644)
			// then check the error
			return nil
		},
	}
)
