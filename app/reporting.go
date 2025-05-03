package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Report struct {
	CampaignId   int    `json:"campaign_id"`
	CampaignName string `json:"campaign"`
	Clicks       int    `json:"clicks"`
	Leads        int    `json:"leads"`
	Conversions  int    `json:"conversions"`
	Sales        int    `json:"sales"`
}
type ReportRequestRange struct {
	Timezone string `json:"timezone"`
	Interval string `json:"interval"`
}

type ReportRequest struct {
	Range   ReportRequestRange `json:"range"`
	Columns []string           `json:"columns"`
}

type ReportResponse struct {
	Reports []Report `json:"rows"`
	Total   int
}

var storedActiveReports []Report

func trackCampaigns(b *gotgbot.Bot, ctx *ext.Context) error {

	reports, err := fetchAllReports()
	if err != nil {
		return err
	}

	slog.Info("Retrieved list of reports", "ReportsCount", len(reports))

	activeReports := getActiveReports(reports)
	activeReports, storedActiveReports = compareStoredReports(activeReports, storedActiveReports)

	slog.Info("Collected active reports", "ActiveReportsCount", len(activeReports))
	slog.Info("Active reports list", "Reports", activeReports)

	slog.Info("Stored active reports from previous run", "Stored reports count", len(storedActiveReports))
	slog.Info("Stored reports list", "Reports", storedActiveReports)

	for _, report := range activeReports {
		for sub := range client.GetAllSubscribers() {
			logger := createTelegramLogger(ctx)

			message := fmt.Sprintf(
				"Campaign: %s (id: %d)\n```Details:\nClicks: %v\nSales: %v\nLeads: %v\nConversions: %v```",
				report.CampaignName,
				report.CampaignId,
				report.Clicks,
				report.Sales,
				report.Leads,
				report.Conversions,
			)

			logger.Info("Sent report to the user")
			b.SendMessage(sub.EffectiveSender.ChatId, message, &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
		}
	}

	return nil
}

func fetchAllReports() ([]Report, error) {
	requestBody := ReportRequest{
		Range: ReportRequestRange{
			Timezone: "Europe/Kyiv",
			Interval: "today",
		},
		Columns: []string{"campaign_id", "campaign", "clicks", "leads", "conversions", "sales"},
	}

	jsonBody, err := json.Marshal(requestBody)

	if err != nil {
		return nil, err
	}
	requestBodyReader := bytes.NewReader(jsonBody)
	request, err := http.NewRequest("POST", os.Getenv("KEITARO_HOST_URL")+"/admin_api/v1/report/build", requestBodyReader)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Api-Key", os.Getenv("KEITARO_API_KEY"))
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var reportResponse ReportResponse
	err = json.Unmarshal(body, &reportResponse)
	if err != nil {
		return nil, err
	}

	return reportResponse.Reports, nil
}

func getActiveReports(reports []Report) []Report {
	var activeReports []Report

	for _, report := range reports {
		if report.Sales <= 0 {
			continue
		}

		storedActiveReportIndex := slices.IndexFunc(storedActiveReports, func(r Report) bool { return r.CampaignId == report.CampaignId && r.Sales == report.Sales })
		if storedActiveReportIndex != -1 {
			continue
		}

		activeReports = append(activeReports, report)
	}

	return activeReports
}

func compareStoredReports(reports []Report, storedReports []Report) ([]Report, []Report) {
	if len(storedReports) == 0 {
		return reports, reports
	}

	for _, report := range reports {
		storedReportIndex := slices.IndexFunc(storedReports, func(r Report) bool { return r.CampaignId == report.CampaignId })
		if storedReportIndex == -1 {
			storedReports = append(storedReports, report)
			continue
		}

		storedReport := storedReports[storedReportIndex]
		if report.Sales > storedReport.Sales {
			storedReports = deleteFromSlice(storedReports, storedReportIndex)
			storedReports = append(storedReports, report)
			continue
		}

		if report.Sales == storedReport.Sales {
			continue
		}

		if report.Sales < storedReport.Sales {
			storedReports = deleteFromSlice(storedReports, storedReportIndex)
		}
	}

	return reports, storedReports
}
