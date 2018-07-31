package datacollect

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func GetMemInfo() string {
	v, _ := mem.VirtualMemory()
	return fmt.Sprintf("Mem: %v MB  Free: %v MB Used:%v Usage:%f%%", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent)
}

func GetMemValue() string {
	v, _ := mem.VirtualMemory()
	return fmt.Sprintf("%v,%v,%v,%f%%", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent)
}

func GetCPUInfo() string {
	cc, _ := cpu.Percent(time.Second, true)
	var total float64 = 0
	for _, c := range cc {
		total = total + c
	}
	return fmt.Sprintf("CPU Used: used %f%%", total)
}

func GetCPUValue() string {
	cc, _ := cpu.Percent(time.Second, true)
	var total float64 = 0
	for _, c := range cc {
		total = total + c
	}
	return fmt.Sprintf("%f%%", total)
}

func GetNetworkInfo() string {
	nv, _ := net.IOCounters(true)
	return fmt.Sprintf("Network: %v bytes / %v bytes", nv[0].BytesRecv, nv[0].BytesSent)
}

func GetNetworkValue() string {
	nv, _ := net.IOCounters(true)
	return fmt.Sprintf("%v / %v", nv[0].BytesRecv, nv[0].BytesSent)
}

func GetAllSysInfos() []string {
	memInfo := GetMemInfo()
	cpuInfo := GetCPUInfo()
	netInfo := GetNetworkInfo()
	resultInfo := []string{memInfo, cpuInfo, netInfo}
	return resultInfo
}

func Collect() {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")
	n, _ := host.Info()
	nv, _ := net.IOCounters(true)
	boottime, _ := host.BootTime()
	btime := time.Unix(int64(boottime), 0).Format("2006-01-02 15:04:05")

	fmt.Println("	Mem	: %v MB  Free: %v MB  Used:%v Usage:%f%%\n", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent)
	if len(c) > 1 {
		for _, sub_cpu := range c {
			modelname := sub_cpu.ModelName
			cores := sub_cpu.Cores
			fmt.Printf("	CPU	: %v  %v cores\n", modelname, cores)
		}
	} else {
		sub_cpu := c[0]
		modelname := sub_cpu.ModelName
		cores := sub_cpu.Cores
		fmt.Printf("	CPU	: %v  %v cores\n", modelname, cores)
	}
	fmt.Printf("	Network: %v bytes / %v bytes\n", nv[0].BytesRecv, nv[0].BytesSent)
	fmt.Printf("	SystemBoot:%v\n", btime)
	fmt.Printf("	CPU Used	: used %f%% \n", cc[0])
	fmt.Printf("	HD	: %v GB  Free: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
	fmt.Printf("	OS	: %v(%v)  %v \n", n.Platform, n.PlatformFamily, n.PlatformVersion)
	fmt.Printf("	Hostname	: %v \n", n.Hostname)
}
