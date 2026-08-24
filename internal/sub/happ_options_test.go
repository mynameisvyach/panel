package sub

import (
	"strings"
	"testing"
)

func TestBuildHappSubscriptionDirectives(t *testing.T) {
	raw := `{"infoColor":"blue","expireEnable":true,"hideServerSettings":true,"autoUpdateOnOpen":true,"pinCurrent":true,"autoUpdate":true,"fragmentationEnable":true,"fragmentationPackets":"tlshello","fragmentationLength":"50-100","fragmentationInterval":"10-20","fragmentationMaxSplit":"100-200","noisesEnable":true,"noisesType":"base64","noisesPacket":"abc=","noisesDelay":"50","noisesApplyTo":"all","pingType":"proxy-head","pingResult":"time","subscriptionSort":"ping"}`
	got := buildHappSubscriptionDirectives(raw)
	for _, want := range []string{
		"#sub-info-color: blue", "#sub-expire: 1", "#subscription-hide-server-settings: 1",
		"#subscription-auto-update-open-enable: 1", "#subscription-pin-current: 1",
		"#subscription-auto-update-enable: 1", "#fragmentation-packets: tlshello",
		"#noises-type: base64", "#ping-type: proxy-head", "#ping-result: time",
		"#subscription-sort: ping",
	} {
		if !strings.Contains(got, want+"\n") { t.Errorf("missing %q in %q", want, got) }
	}
}

func TestBuildHappSubscriptionDirectivesRejectsInvalidJSON(t *testing.T) {
	if got := buildHappSubscriptionDirectives("{"); got != "" { t.Fatalf("got %q", got) }
}
