package jobs

import (
	"testing"

	"github.com/sitename/sitename/einterfaces/mocks"
	"github.com/sitename/sitename/modules/util/testutils"
	"github.com/sitename/sitename/store/storetest"
)

func makeJobServer(t *testing.T) (*JobServer, *storetest.Store, *mocks.MetricsInterface) {
	configService := &testutils.StaticConfigService{}
	mockStore := &storetest.Store{}

	t.Cleanup(func() {
		mockStore.AssertExpectations(t)
	})

	mockMetrics := &mocks.MetricsInterface{}
	t.Cleanup(func() {
		mockMetrics.AssertExpectations(t)
	})

	jobServer := &JobServer{
		ConfigService: configService,
		// Store:         mockStore,
		metrics: mockMetrics,
	}

	return jobServer, mockStore, mockMetrics
}
