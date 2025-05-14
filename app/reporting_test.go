package app

import "testing"

func TestCalculateActiveReports(t *testing.T) {
	t.Run("EmptyReport_NotStored", func(t *testing.T) {
		reports := []Report{}
		storedReports := []Report{}

		activeReports, storedActiveReports := calculateActiveReports(reports, storedReports)

		if len(activeReports) != 0 {
			t.Error("Active reports length is not correct", activeReports)
		}

		if len(storedActiveReports) != 0 {
			t.Error("Stored reports length is not correct", storedActiveReports)
		}
	})

	t.Run("EmptyActiveReports_NotStored", func(t *testing.T) {
		reports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
		}
		storedReports := []Report{}

		activeReports, storedActiveReports := calculateActiveReports(reports, storedReports)

		if len(activeReports) != 0 {
			t.Error("Active reports length is not correct", activeReports)
		}

		if len(storedActiveReports) != 0 {
			t.Error("Stored reports length is not correct", storedActiveReports)
		}
	})

	t.Run("ActiveReports_NotStored", func(t *testing.T) {
		reports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 1, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 2, Leads: 0, Clicks: 0, Conversions: 0},
		}
		storedReports := []Report{}

		activeReports, storedActiveReports := calculateActiveReports(reports, storedReports)

		if len(activeReports) != 2 {
			t.Error("Active reports length is not correct", activeReports)
		}

		if len(storedActiveReports) != 2 {
			t.Error("Stored reports length is not correct", storedActiveReports)
		}
	})

	t.Run("NoReports_Stored", func(t *testing.T) {
		reports := []Report{}
		storedReports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 1, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 2, Leads: 0, Clicks: 0, Conversions: 0},
		}

		activeReports, storedActiveReports := calculateActiveReports(reports, storedReports)

		if len(activeReports) != 0 {
			t.Error("Active reports length is not correct", activeReports)
		}

		if len(storedActiveReports) != 0 {
			t.Error("Stored reports length is not correct", storedActiveReports)
		}
	})

	t.Run("ActiveReport_Stored_NoDeposit", func(t *testing.T) {
		reports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 1, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 2, Leads: 0, Clicks: 0, Conversions: 0},
		}
		storedReports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 1, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 2, Leads: 0, Clicks: 0, Conversions: 0},
		}

		activeReports, storedActiveReports := calculateActiveReports(reports, storedReports)

		if len(activeReports) != 0 {
			t.Error("Active reports length is not correct", activeReports)
		}

		if len(storedActiveReports) != 2 {
			t.Error("Stored reports length is not correct", storedActiveReports)
		}
	})

	t.Run("Reports_Stored_NoDeposit", func(t *testing.T) {
		reports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
		}
		storedReports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 1, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 2, Leads: 0, Clicks: 0, Conversions: 0},
		}

		activeReports, storedActiveReports := calculateActiveReports(reports, storedReports)

		if len(activeReports) != 0 {
			t.Error("Active reports length is not correct", activeReports)
		}

		if len(storedActiveReports) != 0 {
			t.Error("Stored reports length is not correct", storedActiveReports)
		}
	})
	t.Run("ActiveReports_Stored_OneDeposit", func(t *testing.T) {
		reports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 2, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 1, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 2, Leads: 0, Clicks: 0, Conversions: 0},
		}
		storedReports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 1, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 2, Leads: 0, Clicks: 0, Conversions: 0},
		}

		activeReports, storedActiveReports := calculateActiveReports(reports, storedReports)

		if len(activeReports) != 1 {
			t.Error("Active reports length is not correct", activeReports)
		}

		if len(storedActiveReports) != 3 {
			t.Error("Stored reports length is not correct", storedActiveReports)
		}
	})

	t.Run("ActiveReports_Stored_AllDeposit", func(t *testing.T) {

		reports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 1, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 2, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 3, Leads: 0, Clicks: 0, Conversions: 0},
		}
		storedReports := []Report{
			{CampaignId: 1, CampaignName: "Test 1", Sales: 0, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 2, CampaignName: "Test 2", Sales: 1, Leads: 0, Clicks: 0, Conversions: 0},
			{CampaignId: 3, CampaignName: "Test 3", Sales: 2, Leads: 0, Clicks: 0, Conversions: 0},
		}

		activeReports, storedActiveReports := calculateActiveReports(reports, storedReports)

		if len(activeReports) != 3 {
			t.Error("Active reports length is not correct", activeReports)
		}

		if len(storedActiveReports) != 3 {
			t.Error("Stored reports length is not correct", storedActiveReports)
		}
	})
}
