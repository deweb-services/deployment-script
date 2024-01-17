package main

import (
	"os"
	"time"

	"github.com/deweb-services/deployment-script/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	maxTries  = 100
	sleepTime = 10 * time.Second
	APIURL    = "https://app.nodeshift.com"
)

var (
	log     = logger.Logger().Sugar().Named("deployment-script")
	rootCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy instance",
	}
)

func init() {
	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
	})
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

	gpuDeleteCmd.Flags().StringVar(&gpuDeleteCfg, "uuid", "uuid", "uuid of deployment to delete")
	_ = viper.BindPFlag("uuid", rootCmd.Flags().Lookup("uuid"))

	gpuCmd.AddCommand(gpuDeleteCmd)

	cpuCreateCmd.Flags().StringVar(&cpuCreateCfg.Region, "region", "region", "Region where you want to deploy like 'USA'")
	cpuCreateCmd.Flags().StringVar(&cpuCreateCfg.ImageVersion, "image", "Ubuntu-v22.04", "OS Image used to install on the target Vitrual Machine Deployment like 'Ubuntu-v22.04'")
	cpuCreateCmd.Flags().IntVar(&cpuCreateCfg.CPU, "cpu", 1, "number of CPU cores for your deployment")
	cpuCreateCmd.Flags().IntVar(&cpuCreateCfg.RAM, "ram", 1, "Amount of RAM for your Deployment in GB")
	cpuCreateCmd.Flags().IntVar(&cpuCreateCfg.Hdd, "disk_size", 10, "Disk size for your Deployment in GB")
	cpuCreateCmd.Flags().StringVar(&cpuCreateCfg.HddType, "disk_type", "hdd", "Disk type for your Deployment. Available options: hdd, ssd")
	cpuCreateCmd.Flags().BoolVar(&cpuCreateCfg.Ipv4, "assign_public_ipv4", false, "If true assigns a public ipv4 address for your Deployment")
	cpuCreateCmd.Flags().BoolVar(&cpuCreateCfg.Ipv4, "assign_public_ipv6", false, "If true assigns a public ipv6 address for your Deployment")
	cpuCreateCmd.Flags().BoolVar(&cpuCreateCfg.Ygg, "assign_ygg_ip", false, "If true assigns a yggdrasil address for your Deployment")
	cpuCreateCmd.Flags().StringVar(&cpuCreateCfg.NetworkUUID, "vpc_id", "", "ID of the vpc to deploy your VM into")
	cpuCreateCmd.Flags().StringVar(&cpuCreateCfg.SSHKey, "ssh_key", "", "SSH key to add to the target VM to make it possible to connect to your VM")
	cpuCreateCmd.Flags().StringVar(&cpuCreateCfg.HostName, "host_name", "", "Host name for your Deployment")
	_ = viper.BindPFlag("region", rootCmd.Flags().Lookup("region"))
	_ = viper.BindPFlag("image", rootCmd.Flags().Lookup("image"))
	_ = viper.BindPFlag("cpu", rootCmd.Flags().Lookup("cpu"))
	_ = viper.BindPFlag("ram", rootCmd.Flags().Lookup("ram"))
	_ = viper.BindPFlag("disk_size", rootCmd.Flags().Lookup("disk_size"))
	_ = viper.BindPFlag("disk_type", rootCmd.Flags().Lookup("disk_type"))
	_ = viper.BindPFlag("assign_public_ipv4", rootCmd.Flags().Lookup("assign_public_ipv4"))
	_ = viper.BindPFlag("assign_public_ipv6", rootCmd.Flags().Lookup("assign_public_ipv6"))
	_ = viper.BindPFlag("assign_ygg_ip", rootCmd.Flags().Lookup("assign_ygg_ip"))
	_ = viper.BindPFlag("vpc_id", rootCmd.Flags().Lookup("vpc_id"))
	_ = viper.BindPFlag("ssh_key", rootCmd.Flags().Lookup("ssh_key"))
	_ = viper.BindPFlag("host_name", rootCmd.Flags().Lookup("host_name"))
	cpuCmd.AddCommand(cpuCreateCmd)

	cpuDeleteCmd.Flags().StringVar(&cpuDeleteCfg, "uuid", "uuid", "uuid of deployment to delete")
	_ = viper.BindPFlag("uuid", rootCmd.Flags().Lookup("uuid"))
	cpuCmd.AddCommand(cpuDeleteCmd)

	rootCmd.AddCommand(cpuCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("command exited with error: %s", err)
		os.Exit(1)
	}
}
