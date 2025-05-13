package collector_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/F0Rt04ka/otus_project/internal/daemon/collector"
	"github.com/F0Rt04ka/otus_project/internal/daemon/collector/cpuusage"
	"github.com/stretchr/testify/assert"
)

func Test_CollectorResultMap_CheckDeleteOldDataCorrectly(t *testing.T) {
	t.Parallel()

	crm := collector.NewCollectorResultMap(5, 1*time.Second)
	startTime := time.Now().Unix()

	for i := int64(0); i < 6; i++ {
		crm.AddCPUStats(startTime-i, &cpuusage.Result{
			UserMode:   float64(i),
			SystemMode: float64(i),
			Idle:       float64(i),
		})
	}

	crm.RunClearDataHandler(startTime - 1)
	time.Sleep(1 * time.Second)

	// Проверяем, что старые данные были удалены
	for i := int64(0); i < 6; i++ {
		_, ok := crm.GetCPUStats(startTime - i)
		if i < 5 {
			assert.True(t, ok, fmt.Sprintf("Expected data for time %d to exist; Start time %d", startTime-i, startTime))
		} else {
			assert.False(t, ok, fmt.Sprintf("Expected data for time %d to be deleted; Start time %d", startTime-i, startTime))
		}
	}

	stats := crm.GetAvgCPUStats(startTime, 5)
	assert.Equal(t, 2.0, stats.UserMode)
}
