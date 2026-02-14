package logic

import (
	"testing"
	"time"

	apiSponsor "github.com/bamboo-services/bamboo-main/api/sponsor"
	"github.com/bamboo-services/bamboo-main/internal/entity"
)

func TestBuildRecordPublicItemResponse_Anonymous(t *testing.T) {
	redirectURL := "https://example.com"
	record := &entity.SponsorRecord{
		ID:          1001,
		Nickname:    "Visible Name",
		RedirectURL: &redirectURL,
		Amount:      199,
		IsAnonymous: true,
	}

	resp := buildRecordPublicItemResponse(record)

	if resp.Nickname != "匿名用户" {
		t.Fatalf("expected anonymous nickname, got=%q", resp.Nickname)
	}
	if resp.RedirectURL != nil {
		t.Fatal("expected redirect_url to be nil for anonymous record")
	}
}

func TestBuildRecordEntityResponse_IncludesChannelSummary(t *testing.T) {
	icon := "https://example.com/icon.png"
	message := "thanks"
	now := time.Now()
	record := &entity.SponsorRecord{
		ID:        2001,
		Nickname:  "Sponsor",
		Amount:    660,
		Message:   &message,
		SponsorAt: &now,
		ChannelFKey: &entity.SponsorChannel{
			ID:   3001,
			Name: "Alipay",
			Icon: &icon,
		},
	}

	resp := buildRecordEntityResponse(record)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Channel == nil {
		t.Fatal("expected non-nil channel summary")
	}
	if resp.Channel.ID != 3001 {
		t.Fatalf("expected channel id=3001, got=%d", resp.Channel.ID)
	}
	if resp.Channel.Name != "Alipay" {
		t.Fatalf("expected channel name=Alipay, got=%q", resp.Channel.Name)
	}
}

func TestBuildChannelEntityResponse_MapsSponsorCount(t *testing.T) {
	channel := &entity.SponsorChannel{
		ID:     4001,
		Name:   "WeChat",
		Status: true,
	}

	resp := buildChannelEntityResponse(channel, 7)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.SponsorCount != 7 {
		t.Fatalf("expected sponsor_count=7, got=%d", resp.SponsorCount)
	}
	if resp.Name != "WeChat" {
		t.Fatalf("expected channel name=WeChat, got=%q", resp.Name)
	}
}

func TestBuildRecordChannelResponse_NilInput(t *testing.T) {
	var channel *entity.SponsorChannel
	resp := buildRecordChannelResponse(channel)
	if resp != nil {
		t.Fatal("expected nil response for nil channel")
	}
}

func TestBuildRecordChannelResponse_MinimalFields(t *testing.T) {
	icon := "https://example.com/icon.png"
	channel := &entity.SponsorChannel{ID: 5001, Name: "Patreon", Icon: &icon}

	resp := buildRecordChannelResponse(channel)
	if resp == nil {
		t.Fatal("expected non-nil channel response")
	}
	want := apiSponsor.SponsorChannelSimpleResponse{ID: 5001, Name: "Patreon", Icon: &icon}
	if resp.ID != want.ID || resp.Name != want.Name || *resp.Icon != *want.Icon {
		t.Fatalf("unexpected channel response: %+v", resp)
	}
}
