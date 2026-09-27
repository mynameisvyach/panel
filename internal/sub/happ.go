package sub

import (
	"encoding/base64"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var happUserAgentRegex = regexp.MustCompile(`(?i)\bhapp\b`)

// HappConfig holds all Happ client customization parameters.
type HappConfig struct {
	SubscriptionSort      string
	PingResult            string
	NoisesApplyTo         string
	NoisesDelay           string
	NoisesPacket          string
	NoisesType            string
	NoisesEnable          bool
	FragmentationMaxSplit string
	FragmentationInterval string
	FragmentationLength   string
	FragmentationPackets  string
	FragmentationEnable   bool
	AutoUpdate            bool
	PinCurrent            bool
	AutoUpdateOnOpen      bool
	HideServerSettings    bool
	AutoDetect            bool
	ProviderId            string
	NewUrl                string
	FallbackUrl           string
	SubInfoColor          string
	SubInfoText           string
	SubInfoButtonText     string
	SubInfoButtonLink     string
	SubExpire             bool
	SubExpireButtonLink   string
	NotificationExpire    bool
	NoLimit               bool
	AlwaysHwid            bool
	TunMode               string
	TunType               string
	ExcludeRoutes         string
	ExcludeApns           bool
	ColorProfile          string
	PingType              string
	AutoConnect           bool
	AutoConnectType       string
	PerAppMode            string
	PerAppList            string
}

// IsHappClient checks if the client user-agent identifies as Happ.
func IsHappClient(userAgent string) bool {
	return happUserAgentRegex.MatchString(userAgent)
}

func sanitizeHeaderValue(v string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(v), "\r", ""), "\n", "")
}

// ApplyHappHeaders sets standard and advanced Happ subscription headers.
func ApplyHappHeaders(c *gin.Context, cfg HappConfig, isHapp bool) {
	if c == nil || c.Writer == nil || !cfg.AutoDetect || !isHapp {
		return
	}
	if cfg.ProviderId != "" {
		c.Writer.Header().Set("ProviderID", strings.TrimSpace(cfg.ProviderId))
	}
	if cfg.NewUrl != "" {
		c.Writer.Header().Set("New-Url", strings.TrimSpace(cfg.NewUrl))
	}
	if cfg.FallbackUrl != "" {
		c.Writer.Header().Set("Fallback-Url", strings.TrimSpace(cfg.FallbackUrl))
	}
	if text := sanitizeHeaderValue(cfg.SubInfoText); text != "" {
		color := strings.TrimSpace(cfg.SubInfoColor)
		switch strings.ToLower(color) {
		case "primary", "info":
			color = "blue"
		case "success":
			color = "green"
		case "warning", "danger":
			color = "red"
		case "":
			color = "blue"
		}
		c.Writer.Header().Set("Sub-Info-Color", color)
		c.Writer.Header().Set("Sub-Info-Text", text)
		if btnText := sanitizeHeaderValue(cfg.SubInfoButtonText); btnText != "" {
			c.Writer.Header().Set("Sub-Info-Button-Text", btnText)
		}
		if btnLink := sanitizeHeaderValue(cfg.SubInfoButtonLink); btnLink != "" {
			c.Writer.Header().Set("Sub-Info-Button-Link", btnLink)
		}
	}
	if cfg.SubExpire {
		c.Writer.Header().Set("Sub-Expire", "1")
		if link := strings.TrimSpace(cfg.SubExpireButtonLink); link != "" {
			c.Writer.Header().Set("Sub-Expire-Button-Link", link)
		}
	}
	if cfg.NotificationExpire {
		c.Writer.Header().Set("Notification-Subs-Expire", "1")
	}
	if cfg.NoLimit {
		c.Writer.Header().Set("No-Limit-Enabled", "1")
	}
	if cfg.AlwaysHwid {
		c.Writer.Header().Set("Subscription-Always-Hwid-Enable", "1")
	}
	if cfg.TunMode != "" {
		c.Writer.Header().Set("Tun-Mode", cfg.TunMode)
	}
	if cfg.TunType != "" {
		c.Writer.Header().Set("Tun-Type", cfg.TunType)
	}
	if routes := strings.TrimSpace(cfg.ExcludeRoutes); routes != "" {
		c.Writer.Header().Set("Exclude-Routes", routes)
	}
	if cfg.ExcludeApns {
		c.Writer.Header().Set("Exclude-Apns-Enable", "true")
	}
	if profile := sanitizeHeaderValue(cfg.ColorProfile); profile != "" {
		c.Writer.Header().Set("Color-Profile", profile)
	}
	if ping := strings.TrimSpace(cfg.PingType); ping != "" {
		if strings.EqualFold(ping, "http") {
			ping = "proxy"
		}
		c.Writer.Header().Set("Ping-Type", ping)
	}
	if cfg.AutoConnect {
		c.Writer.Header().Set("Subscription-Autoconnect", "1")
		autoType := strings.TrimSpace(cfg.AutoConnectType)
		switch strings.ToLower(autoType) {
		case "fastest":
			autoType = "lowestdelay"
		case "last":
			autoType = "lastused"
		}
		if autoType != "" {
			c.Writer.Header().Set("Subscription-Autoconnect-Type", autoType)
		}
	}
	if mode := strings.TrimSpace(cfg.PerAppMode); mode != "" && mode != "off" {
		switch strings.ToLower(mode) {
		case "include":
			mode = "on"
		case "exclude":
			mode = "bypass"
		}
		c.Writer.Header().Set("Per-App-Proxy-Mode", mode)
		if list := strings.TrimSpace(cfg.PerAppList); list != "" {
			c.Writer.Header().Set("Per-App-Proxy-List", list)
		}
	}
	// Extended Happ subscription controls. Omit disabled flags like the stock options.
	if cfg.HideServerSettings {
		c.Writer.Header().Set("subscription-hide-server-settings", "1")
	}
	if cfg.AutoUpdateOnOpen {
		c.Writer.Header().Set("subscription-auto-update-open-enable", "1")
	}
	if cfg.PinCurrent {
		c.Writer.Header().Set("subscription-pin-current", "1")
	}
	if cfg.AutoUpdate {
		c.Writer.Header().Set("subscription-auto-update-enable", "1")
	}
	if cfg.FragmentationEnable {
		c.Writer.Header().Set("fragmentation-enable", "1")
	}
	if value := sanitizeHeaderValue(cfg.FragmentationPackets); cfg.FragmentationEnable && value != "" {
		c.Writer.Header().Set("fragmentation-packets", value)
	}
	if value := sanitizeHeaderValue(cfg.FragmentationLength); cfg.FragmentationEnable && value != "" {
		c.Writer.Header().Set("fragmentation-length", value)
	}
	if value := sanitizeHeaderValue(cfg.FragmentationInterval); cfg.FragmentationEnable && value != "" {
		c.Writer.Header().Set("fragmentation-interval", value)
	}
	if value := sanitizeHeaderValue(cfg.FragmentationMaxSplit); cfg.FragmentationEnable && value != "" {
		c.Writer.Header().Set("fragmentation-maxsplit", value)
	}
	if cfg.NoisesEnable {
		c.Writer.Header().Set("noises-enable", "1")
	}
	if value := sanitizeHeaderValue(cfg.NoisesType); cfg.NoisesEnable && value != "" {
		c.Writer.Header().Set("noises-type", value)
	}
	if value := sanitizeHeaderValue(cfg.NoisesPacket); cfg.NoisesEnable && value != "" {
		c.Writer.Header().Set("noises-packet", value)
	}
	if value := sanitizeHeaderValue(cfg.NoisesDelay); cfg.NoisesEnable && value != "" {
		c.Writer.Header().Set("noises-delay", value)
	}
	if value := sanitizeHeaderValue(cfg.NoisesApplyTo); cfg.NoisesEnable && value != "" {
		c.Writer.Header().Set("noises-applyto", value)
	}
	if value := sanitizeHeaderValue(cfg.PingResult); value != "" {
		c.Writer.Header().Set("ping-result", value)
	}
	if value := sanitizeHeaderValue(cfg.SubscriptionSort); value != "" {
		c.Writer.Header().Set("subscription-sort", value)
	}

}

// applyHappInboundDescription applies panel metadata only to share links.
// Empty inbound descriptions preserve the upstream host-description behavior.
func applyHappInboundDescription(links, description string) string {
	description = strings.TrimSpace(description)
	if description == "" || links == "" {
		return links
	}
	lines := strings.Split(links, "\n")
	for i, link := range lines {
		scheme, rest, ok := strings.Cut(link, "://")
		if !ok {
			continue
		}
		switch scheme {
		case "vmess":
			var payload []byte
			var err error
			for _, encoding := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
				payload, err = encoding.DecodeString(rest)
				if err == nil {
					break
				}
			}
			if err != nil {
				continue
			}
			var obj map[string]json.RawMessage
			if json.Unmarshal(payload, &obj) != nil || obj == nil {
				continue
			}
			value, _ := json.Marshal(description)
			obj["serverDescription"] = value
			encoded, err := json.Marshal(obj)
			if err == nil {
				lines[i] = "vmess://" + base64.StdEncoding.EncodeToString(encoded)
			}
		case "vless", "trojan", "ss", "hysteria", "hysteria2", "hy2", "tuic":
			base, fragment, _ := strings.Cut(link, "#")
			fragment, _, _ = strings.Cut(fragment, "?serverDescription=")
			lines[i] = base + "#" + fragment + "?serverDescription=" + base64.StdEncoding.EncodeToString([]byte(description))
		}
	}
	return strings.Join(lines, "\n")
}
