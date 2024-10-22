package config

import (
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/splitio/go-client/v6/splitio/client"
	"github.com/splitio/go-client/v6/splitio/conf"

	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
)

type FeatureFlagSyncParams struct {
	ServerID            string
	SplitKey            string
	SyncIntervalSeconds int
	Log                 *slog.Logger
	Attributes          map[string]any
}

type FeatureFlagSynchronizer struct {
	FeatureFlagSyncParams

	client  *client.SplitClient
	stop    chan struct{}
	stopped chan struct{}
}

var featureNames = getStructFields(model_helper.FeatureFlags{})

func NewFeatureFlagSynchronizer(params FeatureFlagSyncParams) (*FeatureFlagSynchronizer, error) {
	cfg := conf.Default()
	if params.Log != nil {
		cfg.Logger = &splitLogger{wrappedLog: params.Log.With(slog.String("service", "split"))}
	} else {
		cfg.LoggerConfig.LogLevel = math.MinInt32
	}
	factory, err := client.NewSplitFactory(params.SplitKey, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create split factory")
	}

	return &FeatureFlagSynchronizer{
		FeatureFlagSyncParams: params,
		client:                factory.Client(),
		stop:                  make(chan struct{}),
		stopped:               make(chan struct{}),
	}, nil
}

// EnsureReady blocks until the syncronizer is ready to update feature flag values
func (f *FeatureFlagSynchronizer) EnsureReady() error {
	if err := f.client.BlockUntilReady(10); err != nil {
		return errors.Wrap(err, "split.io client could not initialize")
	}

	return nil
}

func (f *FeatureFlagSynchronizer) UpdateFeatureFlagValues(base model_helper.FeatureFlags) model_helper.FeatureFlags {
	featuresMap := f.client.Treatments(f.ServerID, featureNames, f.Attributes)
	ffm := featureFlagsFromMap(featuresMap, base)
	return ffm
}

func (f *FeatureFlagSynchronizer) Close() {
	f.client.Destroy()
}

// featureFlagsFromMap sets the feature flags from a map[string]string.
// It starts with baseFeatureFlags and only sets values that are
// given by the upstream management system.
// Makes the assumption that all feature flags are strings or booleans.
// Strings are converted to booleans by considering case insensitive "on" or any value considered by strconv.ParseBool as true and any other value as false.
func featureFlagsFromMap(featuresMap map[string]string, baseFeatureFlags model_helper.FeatureFlags) model_helper.FeatureFlags {
	refStruct := reflect.ValueOf(&baseFeatureFlags).Elem()
	for fieldName, fieldValue := range featuresMap {
		refField := refStruct.FieldByName(fieldName)
		// "control" is returned by split.io if the treatment is not found, in this case we should use the default value.
		if !refField.IsValid() || !refField.CanSet() || fieldValue == "control" {
			continue
		}

		switch refField.Type().Kind() {
		case reflect.Bool:
			parsedBoolValue, _ := strconv.ParseBool(fieldValue)
			refField.Set(reflect.ValueOf(strings.ToLower(fieldValue) == "on" || parsedBoolValue))
		default:
			refField.Set(reflect.ValueOf(fieldValue))
		}

	}
	return baseFeatureFlags
}

// featureFlagsToMap returns the feature flags as a map[string]string
// Supports boolean and string feature flags.
func featureFlagsToMap(featureFlags *model_helper.FeatureFlags) map[string]string {
	refStructVal := reflect.ValueOf(*featureFlags)
	refStructType := reflect.TypeOf(*featureFlags)
	ret := make(map[string]string)
	for i := 0; i < refStructVal.NumField(); i++ {
		refFieldVal := refStructVal.Field(i)
		refFieldType := refStructType.Field(i)
		if !refFieldVal.IsValid() {
			continue
		}
		switch refFieldType.Type.Kind() {
		case reflect.Bool:
			ret[refFieldType.Name] = strconv.FormatBool(refFieldVal.Bool())
		default:
			ret[refFieldType.Name] = refFieldVal.String()
		}
	}

	return ret
}

func getStructFields(s any) []string {
	structType := reflect.TypeOf(s)
	fieldNames := make([]string, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		fieldNames = append(fieldNames, structType.Field(i).Name)
	}

	return fieldNames
}
