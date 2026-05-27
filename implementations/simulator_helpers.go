package implementations

import "bos/types"

func networkRPC(network types.Network) string {
	if network.Rpc != nil {
		return *network.Rpc
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
