package sub

import (
	"encoding/json"
	"fmt"
	"strings"
)

type happSubscriptionConfig struct {
	InfoColor           string `json:"infoColor"`
	InfoText            string `json:"infoText"`
	InfoButtonText      string `json:"infoButtonText"`
	InfoButtonLink      string `json:"infoButtonLink"`
	ExpireEnable        bool   `json:"expireEnable"`
	ExpireButtonLink    string `json:"expireButtonLink"`
	HideServerSettings  bool   `json:"hideServerSettings"`
	AutoUpdateOnOpen    bool   `json:"autoUpdateOnOpen"`
	PinCurrent          bool   `json:"pinCurrent"`
	AutoUpdate          bool   `json:"autoUpdate"`
	FragmentationEnable bool   `json:"fragmentationEnable"`
	FragmentationPackets string `json:"fragmentationPackets"`
	FragmentationLength string `json:"fragmentationLength"`
	FragmentationInterval string `json:"fragmentationInterval"`
	FragmentationMaxSplit string `json:"fragmentationMaxSplit"`
	NoisesEnable        bool   `json:"noisesEnable"`
	NoisesType          string `json:"noisesType"`
	NoisesPacket        string `json:"noisesPacket"`
	NoisesDelay         string `json:"noisesDelay"`
	NoisesApplyTo       string `json:"noisesApplyTo"`
	PingType            string `json:"pingType"`
	PingResult          string `json:"pingResult"`
	SubscriptionSort    string `json:"subscriptionSort"`
}

func buildHappSubscriptionDirectives(raw string) string {
	var cfg happSubscriptionConfig
	if strings.TrimSpace(raw) == "" || json.Unmarshal([]byte(raw), &cfg) != nil {
		return ""
	}
	var out strings.Builder
	line := func(key, value string) {
		value = strings.TrimSpace(value)
		if value != "" && !strings.ContainsAny(value, "\r\n") {
			fmt.Fprintf(&out, "#%s: %s\n", key, value)
		}
	}
	flag := func(key string, enabled bool) { if enabled { line(key, "1") } }
	line("sub-info-color", cfg.InfoColor)
	line("sub-info-text", cfg.InfoText)
	line("sub-info-button-text", cfg.InfoButtonText)
	line("sub-info-button-link", cfg.InfoButtonLink)
	flag("sub-expire", cfg.ExpireEnable)
	if cfg.ExpireEnable { line("sub-expire-button-link", cfg.ExpireButtonLink) }
	flag("subscription-hide-server-settings", cfg.HideServerSettings)
	flag("subscription-auto-update-open-enable", cfg.AutoUpdateOnOpen)
	flag("subscription-pin-current", cfg.PinCurrent)
	flag("subscription-auto-update-enable", cfg.AutoUpdate)
	flag("fragmentation-enable", cfg.FragmentationEnable)
	if cfg.FragmentationEnable {
		line("fragmentation-packets", cfg.FragmentationPackets); line("fragmentation-length", cfg.FragmentationLength)
		line("fragmentation-interval", cfg.FragmentationInterval); line("fragmentation-maxsplit", cfg.FragmentationMaxSplit)
	}
	flag("noises-enable", cfg.NoisesEnable)
	if cfg.NoisesEnable {
		line("noises-type", cfg.NoisesType); line("noises-packet", cfg.NoisesPacket)
		line("noises-delay", cfg.NoisesDelay); line("noises-applyto", cfg.NoisesApplyTo)
	}
	line("ping-type", cfg.PingType)
	line("ping-result", cfg.PingResult)
	line("subscription-sort", cfg.SubscriptionSort)
	return out.String()
}
