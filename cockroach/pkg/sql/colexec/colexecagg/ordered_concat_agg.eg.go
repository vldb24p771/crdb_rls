// Code generated by execgen; DO NOT EDIT.
// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package colexecagg

import (
	"unsafe"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/colmem"
)

func newConcatOrderedAggAlloc(allocator *colmem.Allocator, allocSize int64) aggregateFuncAlloc {
	return &concatOrderedAggAlloc{aggAllocBase: aggAllocBase{
		allocator: allocator,
		allocSize: allocSize,
	}}
}

type concatOrderedAgg struct {
	orderedAggregateFuncBase
	// col points to the output vector we are updating.
	col *coldata.Bytes
	// curAgg holds the running total.
	curAgg []byte
	// foundNonNullForCurrentGroup tracks if we have seen any non-null values
	// for the group that is currently being aggregated.
	foundNonNullForCurrentGroup bool
}

func (a *concatOrderedAgg) SetOutput(vec coldata.Vec) {
	a.orderedAggregateFuncBase.SetOutput(vec)
	a.col = vec.Bytes()
}

func (a *concatOrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, startIdx, endIdx int, sel []int,
) {
	oldCurAggSize := len(a.curAgg)
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Bytes(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture groups to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		groups := a.groups
		if sel == nil {
			_, _ = groups[endIdx-1], groups[startIdx]
			if nulls.MaybeHasNulls() {
				for i := startIdx; i < endIdx; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col.Set(a.curIdx, a.curAgg)
							}
							a.curIdx++
							a.curAgg = zeroBytesValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						a.curAgg = append(a.curAgg, col.Get(i)...)
						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for i := startIdx; i < endIdx; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col.Set(a.curIdx, a.curAgg)
							}
							a.curIdx++
							a.curAgg = zeroBytesValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						a.curAgg = append(a.curAgg, col.Get(i)...)
						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		} else {
			sel = sel[startIdx:endIdx]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col.Set(a.curIdx, a.curAgg)
							}
							a.curIdx++
							a.curAgg = zeroBytesValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						a.curAgg = append(a.curAgg, col.Get(i)...)
						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col.Set(a.curIdx, a.curAgg)
							}
							a.curIdx++
							a.curAgg = zeroBytesValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						a.curAgg = append(a.curAgg, col.Get(i)...)
						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		}
	},
	)
	newCurAggSize := len(a.curAgg)
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsageAfterAllocation(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *concatOrderedAgg) Flush(outputIdx int) {
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	col := a.col
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		col.Set(outputIdx, a.curAgg)
	}
	// Release the reference to curAgg eagerly.
	a.allocator.AdjustMemoryUsage(-int64(len(a.curAgg)))
	a.curAgg = nil
}

func (a *concatOrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.curAgg = nil
	a.foundNonNullForCurrentGroup = false
}

type concatOrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []concatOrderedAgg
}

var _ aggregateFuncAlloc = &concatOrderedAggAlloc{}

const sizeOfConcatOrderedAgg = int64(unsafe.Sizeof(concatOrderedAgg{}))
const concatOrderedAggSliceOverhead = int64(unsafe.Sizeof([]concatOrderedAgg{}))

func (a *concatOrderedAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(concatOrderedAggSliceOverhead + sizeOfConcatOrderedAgg*a.allocSize)
		a.aggFuncs = make([]concatOrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}
