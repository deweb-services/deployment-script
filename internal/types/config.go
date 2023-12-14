package types

import (
	"time"
)

type DWSProviderConfiguration struct {
	Timeout         time.Duration
	AccessKey       string `mapstructure:"access_key"`
	SecretAccessKey string `mapstructure:"secret_key"`
}

type GPUCreateConfig struct {
	GPUName  string `json:"gpuName"`
	Image    string `json:"image"`
	SSHKey   string `json:"sshKey"`
	GPUCount int64  `json:"gpuCount,omitempty"`
	Region   string `json:"region"`
}

type GPUCreateConfigResponse struct {
	Region   string `json:"region,omitempty"`
	Image    string `json:"image"`
	GPUName  string `json:"gpuName"`
	GPUCount int64  `json:"gpuCount"`
	UUID     string `json:"uuid"`
}

type GPUDeleteConfig struct {
	UUID string `json:"uuid"`
}

type RentedGpuInfoResponse struct {
	BundleID          int64    `json:"bundle_id"`
	BundledResults    int64    `json:"bundled_results"`
	BwNvlink          int64    `json:"bw_nvlink"`
	ComputeCap        int64    `json:"compute_cap"`
	CpuCores          int64    `json:"cpu_cores"`
	CpuCoresEffective float64  `json:"cpu_cores_effective"`
	CpuName           string   `json:"cpu_name"`
	CpuRam            int64    `json:"cpu_ram"`
	CudaMaxGood       float64  `json:"cuda_max_good"`
	DirectPortCount   int64    `json:"direct_port_count"`
	DiskBw            float64  `json:"disk_bw"`
	DiskName          string   `json:"disk_name"`
	DiskSpace         float64  `json:"disk_space"`
	Dlperf            float64  `json:"dlperf"`
	DlperfPerDphtotal float64  `json:"dlperf_per_dphtotal"`
	DphBase           float64  `json:"dph_base"`
	DphTotal          float64  `json:"dph_total"`
	DriverVersion     string   `json:"driver_version"`
	Duration          float64  `json:"duration"`
	EndDate           float64  `json:"end_date"`
	External          bool     `json:"external"`
	FlopsPerDphtotal  float64  `json:"flops_per_dphtotal"`
	Geolocation       string   `json:"geolocation,omitempty"`
	GpuDisplayActive  bool     `json:"gpu_display_active"`
	GpuFrac           float64  `json:"gpu_frac"`
	GpuLanes          float64  `json:"gpu_lanes"`
	GpuMemBw          float64  `json:"gpu_mem_bw"`
	GpuName           string   `json:"gpu_name"`
	GpuRam            int64    `json:"gpu_ram"`
	HasAvx            int64    `json:"has_avx"`
	HostId            int64    `json:"host_id"`
	HostRunTime       int64    `json:"host_run_time"`
	HostingType       string   `json:"hosting_type,omitempty"`
	Id                int64    `json:"id"`
	InetDown          float64  `json:"inet_down"`
	InetDownBilled    float64  `json:"inet_down_billed,omitempty"`
	InetDownCost      float64  `json:"inet_down_cost"`
	InetUp            float64  `json:"inet_up"`
	InetUpBilled      float64  `json:"inet_up_billed,omitempty"`
	InetUpCost        float64  `json:"inet_up_cost"`
	IsBid             bool     `json:"is_bid"`
	MachineId         int64    `json:"machine_id"`
	MinBid            float64  `json:"min_bid"`
	MoboName          string   `json:"mobo_name"`
	NumGpus           int64    `json:"num_gpus"`
	PciGen            int64    `json:"pci_gen"`
	PcieBw            float64  `json:"pcie_bw"`
	PendingCount      int64    `json:"pending_count"`
	PublicIpaddr      string   `json:"public_ipaddr"`
	Reliability2      float64  `json:"reliability2"`
	Rentable          bool     `json:"rentable"`
	Rented            bool     `json:"rented"`
	Score             float64  `json:"score"`
	StartDate         float64  `json:"start_date"`
	StorageCost       float64  `json:"storage_cost"`
	StorageTotalCost  float64  `json:"storage_total_cost"`
	TotalFlops        float64  `json:"total_flops"`
	Verification      string   `json:"verification"`
	Webpage           string   `json:"webpage,omitempty"`
	ActualStatus      string   `json:"actual_status,omitempty"`
	CurState          string   `json:"cur_state"`
	DirectPortEnd     int64    `json:"direct_port_end"`
	DirectPortStart   int64    `json:"direct_port_start"`
	DiskUtil          float64  `json:"disk_util"`
	ExtraEnv          []string `json:"extra_env"`
	GpuTemp           float64  `json:"gpu_temp,omitempty"`
	GpuUtil           float64  `json:"gpu_util,omitempty"`
	ImageArgs         []string `json:"image_args,omitempty"`
	ImageRuntype      string   `json:"image_runtype"`
	ImageUuid         string   `json:"image_uuid"`
	IntendedStatus    string   `json:"intended_status"`
	JupyterToken      string   `json:"jupyter_token"`
	Label             string   `json:"label,omitempty"`
	LocalIpaddrs      string   `json:"local_ipaddrs"`
	Logo              string   `json:"logo"`
	MachineDirSshPort int64    `json:"machine_dir_ssh_port"`
	MemLimit          float64  `json:"mem_limit,omitempty"`
	MemUsage          float64  `json:"mem_usage,omitempty"`
	NextState         string   `json:"next_state"`
	Onstart           string   `json:"onstart"`
	SshHost           string   `json:"ssh_host"`
	SshIdx            string   `json:"ssh_idx"`
	SshPort           int64    `json:"ssh_port"`
	StatusMsg         string   `json:"status_msg"`
	VmemUsage         float64  `json:"vmem_usage,omitempty"`
}
