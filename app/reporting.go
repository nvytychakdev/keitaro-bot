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

func trackCampaigns(b *gotgbot.Bot) error {

	readActiveReports()

	reports, err := fetchAllReports()
	if err != nil {
		return err
	}

	slog.Debug("Stored reports list before calculate", "Reports", storedActiveReports)
	slog.Info("Retrieved list of reports", "ReportsCount", len(reports))

	activeReports, storedActiveReports := calculateActiveReports(reports, storedActiveReports)

	slog.Info("Collected active reports", "ActiveReportsCount", len(activeReports))
	slog.Info("Active reports list", "Reports", activeReports)

	slog.Info("Stored active reports from previous run", "Stored reports count", len(storedActiveReports))
	slog.Info("Stored reports list", "Reports", storedActiveReports)

	storeActiveReport(&storedActiveReports)

	for _, report := range activeReports {
		for _, sub := range client.GetAllSubscribers() {
			logger := createTelegramLogger(sub)

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

func calculateActiveReports(reports []Report, storedReports []Report) ([]Report, []Report) {
	// clean existing stored reports
	var remainingStoredReports []Report
	for _, storedReport := range storedReports {
		reportIndex := slices.IndexFunc(reports, func(r Report) bool { return r.CampaignId == storedReport.CampaignId })

		// stored reports that are no longer part of the reports should be removed
		if reportIndex == -1 {
			slog.Debug("Stored report cleaned up", "report", storedReport)
			continue
		}

		// stored reports that has no sales should be removed
		if storedReport.Sales == 0 {
			slog.Debug("Stored report removed due to 0 sales", "report", storedReport)
			continue
		}

		remainingStoredReports = append(remainingStoredReports, storedReport)
	}

	var activeReports []Report
	// form list of active reports and stored reports
	for _, report := range reports {
		storedReportIndex := slices.IndexFunc(remainingStoredReports, func(r Report) bool { return r.CampaignId == report.CampaignId })
		if storedReportIndex == -1 {
			if report.Sales > 0 {
				activeReports = append(activeReports, report)
				remainingStoredReports = append(remainingStoredReports, report)
				slog.Debug("Report was not found, added new active report", "report", report)
			}
			continue
		}

		storedReport := remainingStoredReports[storedReportIndex]
		if report.Sales > storedReport.Sales {
			remainingStoredReports = deleteFromSlice(remainingStoredReports, storedReportIndex)
			remainingStoredReports = append(remainingStoredReports, report)
			activeReports = append(activeReports, report)
			slog.Debug("Report was found, added new active report (report > stored)", "report", report)
			continue
		}

		if report.Sales == storedReport.Sales {
			continue
		}

		if report.Sales < storedReport.Sales {
			remainingStoredReports = deleteFromSlice(remainingStoredReports, storedReportIndex)
			if report.Sales > 0 {
				remainingStoredReports = append(remainingStoredReports, report)
				activeReports = append(activeReports, report)
				slog.Debug("Report was found, added new active report (report < stored && report > 0)", "report", report)
			}
		}
	}

	return activeReports, remainingStoredReports
}

func storeActiveReport(report *[]Report) {
	data, err := json.Marshal(&report)
	if err != nil {
		slog.Error("Not able to marshal active reports", "Err", err.Error())
		return
	}

	err = client.Redis.Set(ctx, "activeReports", data, 0).Err()
	if err != nil {
		slog.Error("Not able to store active reports")
	}
}

func readActiveReports() {
	if client.Redis.Exists(ctx, "activeReports").Val() == 0 {
		slog.Info("No active reports stored")
		return
	}

	data, err := client.Redis.Get(ctx, "activeReports").Bytes()
	if err != nil {
		slog.Error("Not able to retrieve active reports", "Err", err.Error())
		return
	}

	err = json.Unmarshal(data, &storedActiveReports)
	if err != nil {
		slog.Error("Not able to unmarshal stored reports", "Err", err.Error())
		return
	}

	slog.Info("Read active reports from database", "Reports", storedActiveReports)
}
