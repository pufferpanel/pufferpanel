package utils

import (
	"net/netip"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
)

var dockerPortBindingRegex = regexp.MustCompile(`^((?P<ip>([0-9]{0,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})|(\[([0-9a-f:]+)\])):)?((?P<hostport>\d+):)?(?P<port>\d+(\/((tcp)|(udp)))?)$`)
var dockerPortBindingRegexComponents map[string]int = make(map[string]int)

func init() {
	for _, v := range dockerPortBindingRegex.SubexpNames() {
		dockerPortBindingRegexComponents[v] = dockerPortBindingRegex.SubexpIndex(v)
	}
}

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

func ParsePortMap(str string) (network.PortMap, error) {
	parts := dockerPortBindingRegex.FindStringSubmatch(str)

	res := make(network.PortMap)

	portPart := parts[dockerPortBindingRegexComponents["port"]]
	port, err := network.ParsePort(portPart)
	if err != nil {
		return nil, err
	}

	ipPart := parts[dockerPortBindingRegexComponents["ip"]]
	if ipPart == "" {
		ipPart = "0.0.0.0"
	}
	ip, err := netip.ParseAddr(strings.Trim(ipPart, "[]"))
	if err != nil {
		return nil, err
	}

	hostPortPart := parts[dockerPortBindingRegexComponents["hostport"]]
	if hostPortPart == "" {
		hostPortPart = port.Port()
	}

	res[port] = []network.PortBinding{{
		HostIP:   ip,
		HostPort: hostPortPart,
	}}

	return res, nil
}
