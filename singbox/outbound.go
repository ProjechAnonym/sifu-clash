package singbox

import "sifu-clash/models"

func MergeOutbound(provider models.Provider,template string) {
	FetchProxies(provider.Path,provider.Name,template)
}