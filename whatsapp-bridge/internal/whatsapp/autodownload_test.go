package whatsapp

import (
	"testing"

	"go.mau.fi/whatsmeow/types"
)

func TestShouldAutoDownload(t *testing.T) {
	contact := types.NewJID("27820000000", types.DefaultUserServer)
	tests := []struct {
		name       string
		chat       types.JID
		mediaType  string
		url        string
		directPath string
		want       bool
	}{
		{"chat image with url", contact, "image", "https://mmg.whatsapp.net/x", "", true},
		{"chat image with direct path only", contact, "image", "", "/v/t62/x", true},
		{"status story is skipped", types.StatusBroadcastJID, "video", "https://mmg.whatsapp.net/x", "/v/t62/x", false},
		{"text message", contact, "", "", "", false},
		{"media without location", contact, "image", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldAutoDownload(tt.chat, tt.mediaType, tt.url, tt.directPath); got != tt.want {
				t.Errorf("shouldAutoDownload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAutoDownloadPath(t *testing.T) {
	tests := []struct {
		name     string
		chatJID  string
		msgID    string
		filename string
		want     string
	}{
		{"same-second images get distinct names", "120363305435246290@g.us", "3EB0A1", "image_20261016_120814.jpg", "store/120363305435246290_g.us/3EB0A1_image_20261016_120814.jpg"},
		{"second image same second", "120363305435246290@g.us", "3EB0B2", "image_20261016_120814.jpg", "store/120363305435246290_g.us/3EB0B2_image_20261016_120814.jpg"},
		{"missing filename falls back to message id", "27820000000@s.whatsapp.net", "3EB0C3", "", "store/27820000000_s.whatsapp.net/3EB0C3.bin"},
		{"unsafe characters are replaced", "1:2@lid", "ID/1", "a b.pdf", "store/1_2_lid/ID_1_a_b.pdf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := autoDownloadPath(tt.chatJID, tt.msgID, tt.filename); got != tt.want {
				t.Errorf("autoDownloadPath() = %q, want %q", got, tt.want)
			}
		})
	}
}
