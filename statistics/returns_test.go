package statistics

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 这一组测试钉住与 core pkg/statistics 的**逐位一致**（KT-10，裁-18）。
// 黄金值来自 core 同名实现的实跑输出（oracle，seed=20261025）：
// 有限值比 math.Float64bits（不靠 InDelta 掩盖最后一位），
// NaN 比 math.IsNaN（NaN 的 payload 不属于语义）。

// assertParity 逐元素比对：want 为 NaN 时只要求 got 也是 NaN，其余逐位相等。
func assertParity(t *testing.T, want, got []float64) {
	t.Helper()
	require.Len(t, got, len(want))
	for i := range want {
		if math.IsNaN(want[i]) {
			assert.Truef(t, math.IsNaN(got[i]),
				"index %d: want NaN, got bits %d", i, math.Float64bits(got[i]))
			continue
		}
		assert.Equalf(t, math.Float64bits(want[i]), math.Float64bits(got[i]),
			"index %d: want %v, got %v", i, want[i], got[i])
	}
}

func TestLogReturnsParity(t *testing.T) {
	cases := []struct {
		name string
		in   []float64
		want []float64
	}{
		{"nil", nil, nil},
		{"empty", []float64{}, nil},
		{"single", []float64{100}, nil},
		{"two up", []float64{100, 110}, []float64{0.09531017980432493}},
		{"compound", []float64{100, 110, 121}, []float64{0.09531017980432493, 0.09531017980432493}},
		{"zero price", []float64{0, 1, 2}, []float64{math.NaN(), 0.6931471805599453}},
		{"negative price", []float64{-1, 1}, []float64{math.NaN()}},
		{"nan first", []float64{math.NaN(), 1}, []float64{math.NaN()}},
		{"inf first", []float64{math.Inf(1), 1}, []float64{math.NaN()}},
		{"inf second", []float64{1, math.Inf(1)}, []float64{math.NaN()}},
		{"mixed", []float64{100, 90, 95, 100, 100}, []float64{
			-0.10536051565782628, 0.05406722127027579, 0.05129329438755048, 0,
		}},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := LogReturns(c.in)
			if c.want == nil {
				assert.Nil(t, got)
				return
			}
			assertParity(t, c.want, got)
		})
	}
}

// TestLogReturnsRandomParity 用 3 组随机价格向量（seed=20261025）对拍 core：
// 输入是 core 侧 FormatFloat('g', -1, 64) 的往返精确表示，输出比 bits。
func TestLogReturnsRandomParity(t *testing.T) {
	cases := []struct {
		in       []float64
		wantBits []uint64
	}{
		{
			[]float64{
				138.96504978862447, 837.9839651115572, 257.4899944056782,
				984.0910100261065, 408.2280718736615, 440.67730803787214,
				961.0447807379749, 960.8350615329102,
			},
			[]uint64{
				4610770781168522724, 13831365185116002926, 4608716964167973676,
				13829472622251872433, 4590175879906161718, 4605198205998040602,
				13775555836915947491,
			},
		},
		{
			[]float64{
				626.9095913924791, 749.7828039767754, 465.3282947160993,
				190.06644042283094, 957.2824908957931, 149.270816019313,
				646.0364825291645, 4.172443468407794, 238.18689379172008,
				814.7374069625789, 17.010006772590057, 701.7901792250567,
				801.2540891657993,
			},
			[]uint64{
				4595616498768913222, 13825637252702320418, 13829612027615308752,
				4609959900533616759, 13834420058778183741, 4609277014899231923,
				13840735241348613933, 4616239781180439081, 4608217392056011016,
				13839266813376843105, 4615558737301543909, 4593943408276860654,
			},
		},
		{
			[]float64{
				745.6668149406959, 729.2641756047663, 899.5421247994929,
				180.7237038787528, 761.1703780400576, 951.51622782869,
				627.0658572788868, 393.99352438957516, 813.9429762105061,
				935.0253884358257, 592.5408532937156, 918.4280936478734,
				119.22480332461859, 817.4969976523292, 887.9739409570172,
				268.6364648731108, 993.539616705247, 135.26638835347435,
			},
			[]uint64{
				13805440315936335857, 4596728657453445296, 13833278757453411765,
				4609154490370594825, 4597209630554057466, 13824555754102458652,
				13825415255753041150, 4604710445531917843, 4594164615370914053,
				13825260994206164741, 4601566318258895579, 13835151849050125871,
				4611349313176545320, 4590623235408679314, 13831435284815053641,
				4608569144243328320, 13835031159136170468,
			},
		},
	}
	for k, c := range cases {
		got := LogReturns(c.in)
		require.Len(t, got, len(c.wantBits), "vector %d", k)
		for i := range got {
			assert.Equalf(t, c.wantBits[i], math.Float64bits(got[i]),
				"vector %d index %d", k, i)
		}
	}
}

func TestAnnualizeVolatilityParity(t *testing.T) {
	cases := []struct {
		vol      float64
		ppy      int
		wantBits uint64
	}{
		{0, 0, 0}, {0, -1, 0}, {0, 1, 0}, {0, 252, 0}, {0, 365, 0},
		{0.01, 0, 0}, {0.01, -1, 0},
		{0.01, 1, 4576918229304087675},
		{0.01, 252, 4594887414507475660},
		{0.01, 365, 4596051312296118966},
		{0.2, 0, 0}, {0.2, -1, 0},
		{0.2, 1, 4596373779694328218},
		{0.2, 252, 4614331661571192704},
		{0.2, 365, 4615786533806996836},
		{math.Inf(1), 0, 0}, {math.Inf(1), -1, 0},
		{math.Inf(1), 1, 9218868437227405312},
		{math.Inf(1), 252, 9218868437227405312},
		{math.Inf(1), 365, 9218868437227405312},
		{-0.1, 0, 0}, {-0.1, -1, 0},
		{-0.1, 1, 13815242216921733530},
		{-0.1, 252, 13833200098798598016},
		{-0.1, 365, 13834654971034402148},
	}
	for _, c := range cases {
		assert.Equalf(t, c.wantBits, math.Float64bits(AnnualizeVolatility(c.vol, c.ppy)),
			"vol=%v ppy=%d", c.vol, c.ppy)
	}

	// NaN 的 payload 不属于语义：只要求仍是 NaN（非有限输入必须原样传播）。
	for _, ppy := range []int{0, -1, 1, 252, 365} {
		got := AnnualizeVolatility(math.NaN(), ppy)
		if ppy <= 0 {
			assert.Equal(t, 0.0, got, "ppy=%d 时不放大，NaN 前先判非正", ppy)
			continue
		}
		assert.Truef(t, math.IsNaN(got), "ppy=%d: want NaN, got %v", ppy, got)
	}
}
