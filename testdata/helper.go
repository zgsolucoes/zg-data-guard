package testdata

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

func CompareLogs(expectedLog *entity.AccessPermissionLog) any {
	return mock.MatchedBy(func(resultLog *entity.AccessPermissionLog) bool {
		// Ignore ID and Date fields for comparison
		expectedLog.ID = resultLog.ID
		expectedLog.Date = resultLog.Date
		return assert.ObjectsAreEqual(expectedLog, resultLog)
	})
}
