package main

import (
	"fmt"
	"github.com/deweb-services/deployment-script/internal/dws"
	"github.com/deweb-services/deployment-script/internal/types"
	"github.com/deweb-services/deployment-script/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	clientCfg    types.DWSProviderConfiguration
	gpuCreateCfg types.GPUCreateConfig
	gpuDeleteCfg types.GPUDeleteConfig
	log          = logger.Logger().Sugar().Named("deployment-script")
	rootCmd      = &cobra.Command{
		Use:   "deploy",
		Short: "deploy instance",
	}

	gpuCmd = &cobra.Command{
		Use:   "gpu",
		Short: "deploy GPU instance",
	}
	gpuDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete gpu instance",
		Run: func(cmd *cobra.Command, args []string) {
			success := false
			log.Debugf("delete gpu instance with uuid %s", gpuDeleteCfg.UUID)
			cli := dws.NewClient(cmd.Context(), &clientCfg, log, dws.ClientOptWithURL(types.APIURL))
			if err := cli.DeleteGPU(cmd.Context(), gpuDeleteCfg.UUID); err == nil {
				success = true
			} else {
				log.Errorf("delete gpu instance with uuid %s, error: %v", gpuDeleteCfg.UUID, err)
			}
			res := fmt.Sprintf("success=%t", success)
			_ = os.WriteFile("result", []byte(res), 0644)
			log.Debug(res)
		},
	}

	gpuCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create gpu instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debugf("create gpu instance with config %v", gpuCreateCfg)
			cli := dws.NewClient(cmd.Context(), &clientCfg, log, dws.ClientOptWithURL(types.APIURL))
			respCreate, err := cli.CreateGPU(cmd.Context(), &gpuCreateCfg)
			if err != nil {
				return fmt.Errorf("create gpu error: %w", err)
			}
			respGet, err := cli.GetGPU(cmd.Context(), respCreate.UUID)
			if err != nil {
				return fmt.Errorf("get gpu error: %w", err)
			}
			res := fmt.Sprintf("uuid=%s\nhost=%s\nport=%d\n", respCreate.UUID, respGet.SshHost, respGet.SshPort)
			_ = os.WriteFile("result", []byte(res), 0644)
			return nil
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&clientCfg.AccessKey, "access_key", "access_key", "access key for dws platform")
	rootCmd.PersistentFlags().StringVar(&clientCfg.SecretAccessKey, "secret_key", "secret_key", "secret key for dws platform")

	_ = viper.BindPFlag("access_key", rootCmd.PersistentFlags().Lookup("access_key"))
	_ = viper.BindPFlag("secret_key", rootCmd.PersistentFlags().Lookup("secret_key"))

	gpuCreateCmd.Flags().StringVar(&gpuCreateCfg.GPUName, "name", "name", "gpu name like 'RTX_2080'")
	gpuCreateCmd.Flags().StringVar(&gpuCreateCfg.Image, "image", "ubuntu:latest", "image to deploy, like image:version")
	gpuCreateCmd.Flags().StringVar(&gpuCreateCfg.SSHKey, "ssh_key", "ssh_key", "ssh key to connect to instance")
	gpuCreateCmd.Flags().Int64Var(&gpuCreateCfg.GPUCount, "count", 1, "count of gpus to deploy")
	gpuCreateCmd.Flags().StringVar(&gpuCreateCfg.Region, "region", "Europe", "region to deploy instance")
	_ = viper.BindPFlag("name", rootCmd.Flags().Lookup("name"))
	_ = viper.BindPFlag("image", rootCmd.Flags().Lookup("image"))
	_ = viper.BindPFlag("ssh_key", rootCmd.Flags().Lookup("ssh_key"))
	_ = viper.BindPFlag("count", rootCmd.Flags().Lookup("count"))
	_ = viper.BindPFlag("region", rootCmd.Flags().Lookup("region"))
	gpuCmd.AddCommand(gpuCreateCmd)

	gpuDeleteCmd.Flags().StringVar(&gpuDeleteCfg.UUID, "uuid", "uuid", "uuid of deployment to delete")
	_ = viper.BindPFlag("uuid", rootCmd.Flags().Lookup("uuid"))

	gpuCmd.AddCommand(gpuDeleteCmd)
	rootCmd.AddCommand(gpuCmd)
}

func initConfig() {
	viper.AutomaticEnv()
}

func main() {
	_ = rootCmd.Execute()
}
