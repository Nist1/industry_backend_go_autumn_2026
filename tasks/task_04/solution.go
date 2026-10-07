package main

import "math"

type Stats struct {
	Count         int
	Sum, Min, Max int64
}

func Calc(nums []int64) Stats {
	if len(nums) < 2 {
		return Stats{}
	}

	var count int
	var sum int64
	var minDiff int64 = math.MaxInt64
	var maxDiff int64 = math.MinInt64
	for i := 1; i < len(nums); i++ {
		curDiff := nums[i] - nums[i-1]
		count++
		sum += curDiff
		minDiff = min(minDiff, curDiff)
		maxDiff = max(maxDiff, curDiff)
	}

	return Stats{count, sum, minDiff, maxDiff}
}
