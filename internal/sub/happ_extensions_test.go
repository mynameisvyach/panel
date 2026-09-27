package sub

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHappExtendedHeaders(t *testing.T) {
	cfg := HappConfig{AutoDetect: true, HideServerSettings: true, AutoUpdate: true, AutoUpdateOnOpen: true, PinCurrent: true,
		FragmentationEnable: true, FragmentationPackets: "tlshello", FragmentationLength: "50-100", FragmentationInterval: "10-20", FragmentationMaxSplit: "100-200",
		NoisesEnable: true, NoisesType: "hex", NoisesPacket: "abcd", NoisesDelay: "50", NoisesApplyTo: "ip", PingResult: "time", SubscriptionSort: "ping"}
	for _, enabled := range []bool{false, true} {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ApplyHappHeaders(c, cfg, enabled)
		expected := map[string]string{"subscription-hide-server-settings": "1", "subscription-auto-update-enable": "1", "subscription-auto-update-open-enable": "1", "subscription-pin-current": "1", "fragmentation-enable": "1", "fragmentation-packets": "tlshello", "fragmentation-length": "50-100", "fragmentation-interval": "10-20", "fragmentation-maxsplit": "100-200", "noises-enable": "1", "noises-type": "hex", "noises-packet": "abcd", "noises-delay": "50", "noises-applyto": "ip", "ping-result": "time", "subscription-sort": "ping"}
		for key, value := range expected {
			if !enabled {
				value = ""
			}
			if got := w.Header().Get(key); got != value {
				t.Errorf("%s = %q want %q", key, got, value)
			}
		}
	}
	cfg.AutoDetect = false
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ApplyHappHeaders(c, cfg, true)
	if len(w.Header()) != 0 {
		t.Fatal("headers emitted when management disabled")
	}
	cfg = HappConfig{AutoDetect: true, FragmentationPackets: "tlshello", NoisesPacket: "test", PingResult: "time\r\nX-Fake: 1"}
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	ApplyHappHeaders(c, cfg, true)
	if w.Header().Get("fragmentation-packets") != "" || w.Header().Get("noises-packet") != "" {
		t.Fatal("disabled feature parameters emitted")
	}
	if strings.ContainsAny(w.Header().Get("ping-result"), "\r\n") {
		t.Fatal("header injection")
	}
}

func TestHappInboundDescription(t *testing.T) {
	original := "vless://id@example.org:443?type=xhttp#Node%20One"
	if got := applyHappInboundDescription(original, " "); got != original {
		t.Fatal("empty description changed link")
	}
	got := applyHappInboundDescription(original, "  Описание  ")
	want := original + "?serverDescription=" + base64.StdEncoding.EncodeToString([]byte("Описание"))
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if got := applyHappInboundDescription(want, "New"); strings.Count(got, "?serverDescription=") != 1 {
		t.Fatal("duplicate description")
	}
	payload := base64.StdEncoding.EncodeToString([]byte(`{"add":"example.org","id":"id","ps":"Node","serverDescription":"host"}`))
	got = applyHappInboundDescription("vmess://"+payload, "Inbound")
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(got, "vmess://"))
	if err != nil {
		t.Fatal(err)
	}
	var obj map[string]any
	if err = json.Unmarshal(decoded, &obj); err != nil {
		t.Fatal(err)
	}
	if obj["serverDescription"] != "Inbound" || obj["add"] != "example.org" || obj["ps"] != "Node" {
		t.Fatal("bad VMess metadata")
	}
	for _, invalid := range []string{"vmess://broken", "vmess://bnVsbA==", "tg://proxy?server=example.org"} {
		if applyHappInboundDescription(invalid, "test") != invalid {
			t.Fatal("unsupported link modified")
		}
	}
}
