package awsmeter

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/marketplacemetering"
	"github.com/aws/aws-sdk-go/service/marketplacemetering/marketplacemeteringiface"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go/utils"

	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
)

type AwsMeter struct {
	store   store.Store
	service *AWSMeterService
	config  *model_helper.Config
}

type AWSMeterService struct {
	AwsDryRun      bool
	AwsProductCode string
	AwsMeteringSvc marketplacemeteringiface.MarketplaceMeteringAPI
}

type AWSMeterReport struct {
	Dimension string    `json:"dimension"`
	Value     int64     `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func (o *AWSMeterReport) ToJSON() string {
	return model_helper.ModelToJson(o)
}

func New(store store.Store, config *model_helper.Config) *AwsMeter {
	svc := &AWSMeterService{
		AwsDryRun:      false,
		AwsProductCode: "12345", //TODO
	}

	service, err := newAWSMarketplaceMeteringService()
	if err != nil {
		slog.Debug("Could not create AWS metering service", slog.String("error", err.Error()))
		return nil
	}

	svc.AwsMeteringSvc = service
	return &AwsMeter{
		store:   store,
		service: svc,
		config:  config,
	}
}

func newAWSMarketplaceMeteringService() (*marketplacemetering.MarketplaceMetering, error) {
	region := os.Getenv("AWS_REGION")
	s, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		return nil, err
	}

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(s),
			},
		})

	_, err = creds.Get()
	if err != nil {
		return nil, errors.Wrap(err, "cannot obtain credentials")
	}

	return marketplacemetering.New(session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
	}))), nil
}

// a report entry is for all metrics
func (awsm *AwsMeter) GetUserCategoryUsage(dimensions []string, startTime time.Time, endTime time.Time) []*AWSMeterReport {
	reports := make([]*AWSMeterReport, 0)

	for _, dimension := range dimensions {
		var userCount int64
		var err error

		switch dimension {
		case model_helper.AwsMeteringDimensionUsageHrs:
			userCount, err = awsm.store.User().AnalyticsActiveCountForPeriod(utils.TimeToMicrosecondsSinceEpochInt64(startTime), utils.TimeToMicrosecondsSinceEpochInt64(endTime), model_helper.UserCountOptions{})
			if err != nil {
				slog.Warn("Failed to obtain usage data", slog.String("dimension", dimension), slog.String("start", startTime.String()), slog.Int64("count", userCount), slog.Err(err))
				continue
			}
		default:
			slog.Debug("Dimension does not exist!", slog.String("dimension", dimension))
			continue
		}

		report := &AWSMeterReport{
			Dimension: dimension,
			Value:     userCount,
			Timestamp: startTime,
		}

		reports = append(reports, report)
	}

	return reports
}

func (awsm *AwsMeter) ReportUserCategoryUsage(reports []*AWSMeterReport) error {
	for _, report := range reports {
		err := sendReportToMeteringService(awsm.service, report)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendReportToMeteringService(ams *AWSMeterService, report *AWSMeterReport) error {
	params := &marketplacemetering.MeterUsageInput{
		DryRun:         aws.Bool(ams.AwsDryRun),
		ProductCode:    aws.String(ams.AwsProductCode),
		UsageDimension: aws.String(report.Dimension),
		UsageQuantity:  aws.Int64(report.Value),
		Timestamp:      aws.Time(report.Timestamp),
	}

	resp, err := ams.AwsMeteringSvc.MeterUsage(params)
	if err != nil {
		return errors.Wrap(err, "Invalid metering service id.")
	}
	if resp.MeteringRecordId == nil {
		return errors.Wrap(err, "Invalid metering service id.")
	}

	slog.Debug("Sent record to AWS metering service", slog.String("dimension", report.Dimension), slog.Int64("value", report.Value), slog.String("timestamp", report.Timestamp.String()))

	return nil
}
