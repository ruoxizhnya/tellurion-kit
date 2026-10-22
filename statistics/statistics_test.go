package statistics

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSumMeanMeanFinite(t *testing.T) {
	assert.Equal(t, 6.0, Sum([]float64{1, 2, 3}))
	assert.Equal(t, 2.0, Mean([]float64{1, 2, 3}))
	assert.Equal(t, 0.0, Mean(nil), "empty slice → 0, not NaN")

	// Mean 不静默丢掉非有限值：一个 NaN 会把结果变成 NaN（可观察的失败）。
	assert.True(t, math.IsNaN(Mean([]float64{1, math.NaN(), 3})))
	// MeanFinite 才过滤。
	assert.Equal(t, 2.0, MeanFinite([]float64{1, math.NaN(), 3, math.Inf(1)}))
	assert.Equal(t, 0.0, MeanFinite(nil))
}

func TestVariancePopulationVsSample(t *testing.T) {
	v := []float64{1, 2, 3, 4}
	require.InDelta(t, 1.25, PopulationVariance(v), 1e-12)
	require.InDelta(t, 5.0/3.0, SampleVariance(v), 1e-12)
	require.InDelta(t, math.Sqrt(1.25), PopulationStdDev(v), 1e-12)
	require.InDelta(t, math.Sqrt(5.0/3.0), SampleStdDev(v), 1e-12)

	// 单元素：总体方差 0；样本方差无定义，返回 0（不 panic、不 NaN）。
	assert.Equal(t, 0.0, PopulationVariance([]float64{7}))
	assert.Equal(t, 0.0, SampleVariance([]float64{7}))
}

func TestPercentileLinearInterpolation(t *testing.T) {
	v := []float64{4, 1, 3, 2}
	assert.Equal(t, 1.0, Percentile(v, 0))
	assert.Equal(t, 4.0, Percentile(v, 1))
	require.InDelta(t, 2.5, Percentile(v, 0.5), 1e-12)
	// p 越界被 clamp。
	assert.Equal(t, 1.0, Percentile(v, -1))
	assert.Equal(t, 4.0, Percentile(v, 2))
	assert.Equal(t, 0.0, Percentile(nil, 0.5))
	// 输入不被就地排序。
	assert.Equal(t, []float64{4, 1, 3, 2}, v)
}

func TestMedian(t *testing.T) {
	require.InDelta(t, 2.5, Median([]float64{1, 2, 3, 4}), 1e-12)
	assert.Equal(t, 3.0, Median([]float64{1, 3, 5}))
}

func TestRankAveragesTiesAndKeepsDenseRanks(t *testing.T) {
	got := Rank([]float64{10, 20, 20, 40})
	assert.Equal(t, []float64{1, 2.5, 2.5, 4}, got)

	// NaN 得到 0，且不占名次。
	got = Rank([]float64{10, math.NaN(), 30})
	assert.Equal(t, []float64{1, 0, 2}, got)
}

func TestPercentileRank(t *testing.T) {
	got := PercentileRank([]float64{10, 20, 20, 40})
	assert.Equal(t, []float64{0, 0.5, 0.5, 1}, got)
	assert.Nil(t, PercentileRank(nil))
}

func TestPearsonAndSpearman(t *testing.T) {
	require.InDelta(t, 1.0, Pearson([]float64{1, 2, 3}, []float64{2, 4, 6}), 1e-12)
	require.InDelta(t, -1.0, Pearson([]float64{1, 2, 3}, []float64{3, 2, 1}), 1e-12)

	assert.True(t, math.IsNaN(Pearson([]float64{1, 2}, []float64{1})), "长度不等 → NaN")
	assert.True(t, math.IsNaN(Pearson([]float64{1}, []float64{2})), "少于 2 个 → NaN")
	assert.True(t, math.IsNaN(Pearson([]float64{1, 1}, []float64{1, 2})), "零方差 → NaN")

	// 严格单调但非线性的关系：Pearson < 1，Spearman = 1。
	x := []float64{1, 2, 3, 4}
	y := []float64{1, 4, 9, 16}
	assert.Less(t, Pearson(x, y), 0.99)
	require.InDelta(t, 1.0, Spearman(x, y), 1e-12)
	require.InDelta(t, 1.0, Spearman(x, []float64{5, 6, 7, 8}), 1e-12)
}

func TestSlopeAndLinearRegression(t *testing.T) {
	y := []float64{3, 5, 7, 9} // y = 2x + 3
	require.InDelta(t, 2.0, Slope(y), 1e-12)

	slope, intercept, r2 := LinearRegression(y)
	require.InDelta(t, 2.0, slope, 1e-12)
	require.InDelta(t, 3.0, intercept, 1e-12)
	require.InDelta(t, 1.0, r2, 1e-12)

	assert.Equal(t, 0.0, Slope([]float64{1}))
	slope, intercept, r2 = LinearRegression(nil)
	assert.Equal(t, 0.0, slope)
	assert.Equal(t, 0.0, intercept)
	assert.Equal(t, 0.0, r2)
}

func TestZScore(t *testing.T) {
	got := ZScore([]float64{1, 2, 3})
	require.Len(t, got, 3)
	require.InDelta(t, -math.Sqrt(1.5), got[0], 1e-12) // (1-2)/sqrt(2/3)
	require.InDelta(t, 0, got[1], 1e-12)
	require.InDelta(t, math.Sqrt(1.5), got[2], 1e-12)

	// 全部相同 → 全 0（不除零）。
	assert.Equal(t, []float64{0, 0}, ZScore([]float64{5, 5}))
	// 少于 2 个 → 全 0，长度不变。
	assert.Equal(t, []float64{0}, ZScore([]float64{5}))
}
