package functions

import (
	"fmt"
	"incertus/internal/common"
)

type DataRecord = common.DataRecord

func ensureArgumentCount(expected uint, got []DataRecord) error {
	if expected == uint(len(got)) {
		return nil
	}
	return fmt.Errorf("Wrong number of arguments, %d expected, got %d", expected, len(got))
}
