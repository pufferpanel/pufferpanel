package utils

import (
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
)

func CalculateDockerCPUPercent(v *container.StatsResponse) float64 {
	//this math is from https://docs.docker.com/reference/api/engine/version/v1.45/#tag/Container/operation/ContainerStats
	cpuDelta := v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage
	systemCpuDelta := v.CPUStats.SystemUsage - v.PreCPUStats.SystemUsage
	numCpus := int(v.CPUStats.OnlineCPUs)
	if numCpus == 0 {
		numCpus = len(v.CPUStats.CPUUsage.PercpuUsage)
	}
	return (float64(cpuDelta) / float64(systemCpuDelta)) * float64(numCpus) * 100.0
}

func CalculateDockerMemoryPercent(v *container.StatsResponse) float64 {
	return float64(v.MemoryStats.Usage)
}

func ConvertToDockerBind(source string) string {
	fullPath, err := filepath.Abs(source)
	if err != nil {
		panic(err)
	}

	fullPath = strings.ReplaceAll(fullPath, "\\", "/")
	fullPath = strings.ReplaceAll(fullPath, ":", "")
	//lowercase first character as that's the drive
	fullPath = strings.ToLower(string(fullPath[0])) + fullPath[1:]
	fullPath = "/" + fullPath
	return fullPath
}
