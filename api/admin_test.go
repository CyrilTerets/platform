// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/store"
	"github.com/mattermost/platform/utils"
	"testing"
)

func TestGetLogs(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()

	if _, err := th.BasicClient.GetLogs(); err == nil {
		t.Fatal("Shouldn't have permissions")
	}

	if logs, err := th.SystemAdminClient.GetLogs(); err != nil {
		t.Fatal(err)
	} else if len(logs.Data.([]string)) <= 0 {
		t.Fatal()
	}
}

func TestGetAllAudits(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()

	if _, err := th.BasicClient.GetAllAudits(); err == nil {
		t.Fatal("Shouldn't have permissions")
	}

	if audits, err := th.SystemAdminClient.GetAllAudits(); err != nil {
		t.Fatal(err)
	} else if len(audits.Data.(model.Audits)) <= 0 {
		t.Fatal()
	}
}

func TestGetClientProperties(t *testing.T) {
	th := Setup().InitBasic()

	if result, err := th.BasicClient.GetClientProperties(); err != nil {
		t.Fatal(err)
	} else {
		props := result.Data.(map[string]string)

		if len(props["Version"]) == 0 {
			t.Fatal()
		}
	}
}

func TestGetConfig(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()

	if _, err := th.BasicClient.GetConfig(); err == nil {
		t.Fatal("Shouldn't have permissions")
	}

	if result, err := th.SystemAdminClient.GetConfig(); err != nil {
		t.Fatal(err)
	} else {
		cfg := result.Data.(*model.Config)

		if len(cfg.TeamSettings.SiteName) == 0 {
			t.Fatal()
		}
	}
}

func TestSaveConfig(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()

	if _, err := th.BasicClient.SaveConfig(utils.Cfg); err == nil {
		t.Fatal("Shouldn't have permissions")
	}

	if result, err := th.SystemAdminClient.SaveConfig(utils.Cfg); err != nil {
		t.Fatal(err)
	} else {
		cfg := result.Data.(*model.Config)

		if len(cfg.TeamSettings.SiteName) == 0 {
			t.Fatal()
		}
	}
}

func TestEmailTest(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()

	if _, err := th.BasicClient.TestEmail(utils.Cfg); err == nil {
		t.Fatal("Shouldn't have permissions")
	}

	if _, err := th.SystemAdminClient.TestEmail(utils.Cfg); err != nil {
		t.Fatal(err)
	}
}

func TestGetTeamAnalyticsStandard(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()

	if _, err := th.BasicClient.GetTeamAnalytics(th.BasicTeam.Id, "standard"); err == nil {
		t.Fatal("Shouldn't have permissions")
	}

	if result, err := th.SystemAdminClient.GetTeamAnalytics(th.BasicTeam.Id, "standard"); err != nil {
		t.Fatal(err)
	} else {
		rows := result.Data.(model.AnalyticsRows)

		if rows[0].Name != "channel_open_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[0].Value != 2 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[1].Name != "channel_private_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[1].Value != 1 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[2].Name != "post_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[2].Value != 1 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[3].Name != "unique_user_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[3].Value != 2 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[4].Name != "team_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[4].Value == 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}
	}

	if result, err := th.SystemAdminClient.GetSystemAnalytics("standard"); err != nil {
		t.Fatal(err)
	} else {
		rows := result.Data.(model.AnalyticsRows)

		if rows[0].Name != "channel_open_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[0].Value < 2 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[1].Name != "channel_private_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[1].Value == 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[2].Name != "post_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[2].Value == 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[3].Name != "unique_user_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[3].Value == 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[4].Name != "team_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[4].Value == 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}
	}
}

func TestGetPostCount(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()

	// manually update creation time, since it's always set to 0 upon saving and we only retrieve posts < today
	Srv.Store.(*store.SqlStore).GetMaster().Exec("UPDATE Posts SET CreateAt = :CreateAt WHERE ChannelId = :ChannelId",
		map[string]interface{}{"ChannelId": th.BasicChannel.Id, "CreateAt": utils.MillisFromTime(utils.Yesterday())})

	if _, err := th.BasicClient.GetTeamAnalytics(th.BasicTeam.Id, "post_counts_day"); err == nil {
		t.Fatal("Shouldn't have permissions")
	}

	if result, err := th.SystemAdminClient.GetTeamAnalytics(th.BasicTeam.Id, "post_counts_day"); err != nil {
		t.Fatal(err)
	} else {
		rows := result.Data.(model.AnalyticsRows)

		if rows[0].Value != 1 {
			t.Log(rows.ToJson())
			t.Fatal()
		}
	}
}

func TestUserCountsWithPostsByDay(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()

	// manually update creation time, since it's always set to 0 upon saving and we only retrieve posts < today
	Srv.Store.(*store.SqlStore).GetMaster().Exec("UPDATE Posts SET CreateAt = :CreateAt WHERE ChannelId = :ChannelId",
		map[string]interface{}{"ChannelId": th.BasicChannel.Id, "CreateAt": utils.MillisFromTime(utils.Yesterday())})

	if _, err := th.BasicClient.GetTeamAnalytics(th.BasicTeam.Id, "user_counts_with_posts_day"); err == nil {
		t.Fatal("Shouldn't have permissions")
	}

	if result, err := th.SystemAdminClient.GetTeamAnalytics(th.BasicTeam.Id, "user_counts_with_posts_day"); err != nil {
		t.Fatal(err)
	} else {
		rows := result.Data.(model.AnalyticsRows)

		if rows[0].Value != 1 {
			t.Log(rows.ToJson())
			t.Fatal()
		}
	}
}

func TestGetTeamAnalyticsExtra(t *testing.T) {
	th := Setup().InitBasic().InitSystemAdmin()

	th.CreatePost(th.BasicClient, th.BasicChannel)

	if _, err := th.BasicClient.GetTeamAnalytics("", "extra_counts"); err == nil {
		t.Fatal("Shouldn't have permissions")
	}

	if result, err := th.SystemAdminClient.GetTeamAnalytics(th.BasicTeam.Id, "extra_counts"); err != nil {
		t.Fatal(err)
	} else {
		rows := result.Data.(model.AnalyticsRows)

		if rows[0].Name != "file_post_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[0].Value != 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[1].Name != "hashtag_post_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[1].Value != 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[2].Name != "incoming_webhook_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[2].Value != 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[3].Name != "outgoing_webhook_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[3].Value != 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[4].Name != "command_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[4].Value != 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[5].Name != "session_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[5].Value == 0 {
			t.Log(rows.ToJson())
			t.Fatal()
		}
	}

	if result, err := th.SystemAdminClient.GetSystemAnalytics("extra_counts"); err != nil {
		t.Fatal(err)
	} else {
		rows := result.Data.(model.AnalyticsRows)

		if rows[0].Name != "file_post_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[1].Name != "hashtag_post_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[1].Value < 1 {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[2].Name != "incoming_webhook_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[3].Name != "outgoing_webhook_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[4].Name != "command_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}

		if rows[5].Name != "session_count" {
			t.Log(rows.ToJson())
			t.Fatal()
		}
	}
}
