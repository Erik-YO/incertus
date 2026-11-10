package post_processors

import (
	"incertus/internal/common"
	"iter"
	"maps"
	"slices"
)

type OrderKind int

const (
	OrderKindAsc OrderKind = iota
	OrderKindDesc
)

type SortKind int

const (
	SortKindNo SortKind = iota
	SortKindValue
	SortKindFrequency
)

type NormalizationKind int

const (
	NormalizationKindNo NormalizationKind = iota
	NormalizationKind1
	NormalizationKind100
)

type PostProcessor struct {
	Sort          SortKind
	Order         OrderKind
	Normalization NormalizationKind
}

func (p PostProcessor) Process(r common.DataRecord) common.ResultSequence {

	keys := p.keys(r)
	frequencySet := p.normalize(r)

	return func(yield func(int64, float64) bool) {

		for key := range keys {
			if !yield(key, frequencySet[key]) {
				return
			}

		}
	}

}

func (p PostProcessor) normalize(r common.DataRecord) map[int64]float64 {
	normalized := map[int64]float64{}
	for k, v := range r {
		normalized[k] = float64(v)
	}

	if p.Normalization == NormalizationKindNo {
		return normalized
	}

	var total float64
	for _, v := range r {
		total += float64(v)
	}

	var factor float64

	switch p.Normalization {
	case NormalizationKind1:
		factor = 1
	case NormalizationKind100:
		factor = 100
	default:
		panic("Unsupported normalization option")
	}

	for k, v := range normalized {
		// Multiplication before division to avoid potential loss in accuracy
		normalized[k] = v * factor / total
	}

	return normalized
}

func (p PostProcessor) keys(r common.DataRecord) iter.Seq[int64] {
	if p.Sort == SortKindNo {
		return maps.Keys(r)
	}

	sorted := make([]int64, 0, len(r))
	for k := range r {
		sorted = append(sorted, k)
	}

	switch p.Sort {

	case SortKindValue:
		slices.Sort(sorted)

	case SortKindFrequency:

		slices.SortFunc(sorted, func(aV, bV int64) int {
			aF := r[aV]
			bF := r[bV]
			if aF == bF {
				return 0
			}
			if aF < bF {
				return -1
			}
			return 1
		})

	default:
		panic("Unsupported sort option")
	}

	if p.Order == OrderKindDesc {
		slices.Reverse(sorted)
	}

	return slices.Values(sorted)
}
