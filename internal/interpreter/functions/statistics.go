package functions

import (
	"maps"
	"slices"
)

func StatMean(r []DataRecord) (DataRecord, error) {

	if err := ensureArgumentCount(1, r); err != nil {
		return DataRecord{}, err
	}
	a := r[0]

	var valueSum int64 = 0
	var freqSum int64 = 0
	for nV, nF := range a {
		valueSum += nV * int64(nF)
		freqSum += int64(nF)
	}

	// Integer division: rounds down
	mean := valueSum / freqSum

	return DataRecord{mean: 1}, nil
}

func StatMedian(r []DataRecord) (DataRecord, error) {

	if err := ensureArgumentCount(1, r); err != nil {
		return DataRecord{}, err
	}
	a := r[0]

	if len(a) < 1 {
		return DataRecord{}, nil
	}

	sorted := slices.Sorted(maps.Keys(a))

	var freqSum uint64 = 0
	for _, nF := range a {
		freqSum += nF
	}

	// Middle position
	freqSum = freqSum / 2

	for _, value := range sorted {
		freq, found := a[value]
		if !found {
			panic("Should be in map")
		}
		if freq >= freqSum {
			return DataRecord{value: 1}, nil
		}
		freqSum -= freq
	}

	panic("Median not found")
}
