package idgen

import (
	"fmt"
	"hash/fnv"
	"net"
)

// NodeID - returns a node (machine) related identifier
func NodeID() (int64, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve network interfaces, %w", err)
	}

	var addr string
	for _, i := range ifaces {
		if i.Flags&net.FlagLoopback == 0 && i.Flags&net.FlagUp == 1 {
			addr = i.HardwareAddr.String()
			break
		}
	}

	if len(addr) == 0 {
		return 0, fmt.Errorf("failed to find network interface")
	}

	hash := fnv.New64()
	_, err = hash.Write([]byte(addr))
	if err != nil {
		return 0, fmt.Errorf("failed to build node id hash (%s), %w", addr, err)
	}

	return NormalizeNodeID(hash.Sum64()), nil
}

func NormalizeNodeID(nodeID uint64) int64 {
	for nodeID > MaxNodeID {
		nodeID = nodeID >> 1
	}

	return int64(nodeID)
}
