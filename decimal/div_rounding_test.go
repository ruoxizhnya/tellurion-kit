package decimal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDiv_Rounding pins the rounding rule Div applies at the requested
// result scale. KT-03 requires this to be executable evidence rather
// than prose: the table covers a tie at exactly half a unit on both
// signs and both divisor signs, remainders above and below half, odd
// divisors (where a tie cannot occur), and resultScale above the
// dividend's scale. A tie steps away from zero (四舍五入), the same rule
// Round uses; truncating a tie towards zero would systematically
// undercount positive amounts.
func TestDiv_Rounding(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		a, b  Decimal
		scale int
		want  string
	}{
		// No remainder: the quotient is exact.
		{name: "exact_10_div_4_at_2", a: New(10, 0), b: New(4, 0), scale: 2, want: "2.50"},
		{name: "exact_1_div_8_at_5", a: New(1, 0), b: New(8, 0), scale: 5, want: "0.12500"},

		// Remainder below half: the quotient is kept as truncated.
		{name: "below_half_1_div_3_at_2", a: New(1, 0), b: New(3, 0), scale: 2, want: "0.33"},
		{name: "below_half_1_div_7_at_2", a: New(1, 0), b: New(7, 0), scale: 2, want: "0.14"},
		{name: "below_half_1_div_8_at_0", a: New(1, 0), b: New(8, 0), scale: 0, want: "0"},

		// Remainder above half: the quotient steps away from zero.
		{name: "above_half_2_div_3_at_2", a: New(2, 0), b: New(3, 0), scale: 2, want: "0.67"},
		{name: "above_half_2_div_7_at_2", a: New(2, 0), b: New(7, 0), scale: 2, want: "0.29"},
		{name: "above_half_1_div_6_at_1", a: New(1, 0), b: New(6, 0), scale: 1, want: "0.2"},
		{name: "above_half_9_div_8_at_0", a: New(9, 0), b: New(8, 0), scale: 0, want: "1"},
		{name: "above_half_negative_1_div_3_at_2", a: New(-1, 0), b: New(3, 0), scale: 2, want: "-0.33"},

		// Exactly half a unit (KT-03 decision point): half away from zero.
		{name: "tie_1_div_2_at_0", a: New(1, 0), b: New(2, 0), scale: 0, want: "1"},
		{name: "tie_neg_1_div_2_at_0", a: New(-1, 0), b: New(2, 0), scale: 0, want: "-1"},
		{name: "tie_3_div_2_at_0", a: New(3, 0), b: New(2, 0), scale: 0, want: "2"},
		{name: "tie_neg_3_div_2_at_0", a: New(-3, 0), b: New(2, 0), scale: 0, want: "-2"},
		{name: "tie_5_div_2_at_0", a: New(5, 0), b: New(2, 0), scale: 0, want: "3"},
		{name: "tie_1_div_8_at_2", a: New(1, 0), b: New(8, 0), scale: 2, want: "0.13"},
		{name: "tie_neg_1_div_8_at_2", a: New(-1, 0), b: New(8, 0), scale: 2, want: "-0.13"},
		{name: "tie_1_div_4_at_1", a: New(1, 0), b: New(4, 0), scale: 1, want: "0.3"},
		{name: "tie_neg_1_div_4_at_1", a: New(-1, 0), b: New(4, 0), scale: 1, want: "-0.3"},
		{name: "tie_3_div_4_at_1", a: New(3, 0), b: New(4, 0), scale: 1, want: "0.8"},
		{name: "tie_5_div_4_at_1", a: New(5, 0), b: New(4, 0), scale: 1, want: "1.3"},
		{name: "tie_1_div_20_at_1", a: New(1, 0), b: New(20, 0), scale: 1, want: "0.1"},
		{name: "tie_1_div_16_at_3", a: New(1, 0), b: New(16, 0), scale: 3, want: "0.063"},
		{name: "tie_3_div_16_at_3", a: New(3, 0), b: New(16, 0), scale: 3, want: "0.188"},
		{name: "tie_1_div_64_at_5", a: New(1, 0), b: New(64, 0), scale: 5, want: "0.01563"},
		// The sign of the quotient follows the sign of b as well.
		{name: "tie_1_div_neg_2_at_0", a: New(1, 0), b: New(-2, 0), scale: 0, want: "-1"},
		{name: "tie_neg_1_div_neg_2_at_0", a: New(-1, 0), b: New(-2, 0), scale: 0, want: "1"},

		// resultScale above the dividend's scale: extra places are exact
		// unless the quotient still has a remainder (then the tie rule
		// above applies).
		{name: "scale_above_dividend_exact", a: New(25, 1), b: New(4, 0), scale: 3, want: "0.625"},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got, err := Div(c.a, c.b, c.scale)
			require.NoError(t, err)
			assert.Equal(t, c.want, got.ToString())
			assert.Equal(t, c.scale, got.Scale())
		})
	}
}

// TestDiv_ResultScaleBelowDividendScale pins what happens when the
// requested scale is smaller than the dividend's own scale. The
// dividend is first brought to the target scale by dropping its excess
// fractional digits (towards zero), and only the resulting quotient is
// rounded. This is today's behaviour, not the ideal one: 2.5 / 1 at
// scale 0 answers 2, not 3. Making this path exact is tracked
// separately (KT-06) so that it does not ride in on the KT-03
// rounding-mode decision.
func TestDiv_ResultScaleBelowDividendScale(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		a, b  Decimal
		scale int
		want  string
	}{
		{name: "2.5_div_1_at_0", a: New(25, 1), b: New(1, 0), scale: 0, want: "2"},
		{name: "2.5_div_4_at_0", a: New(25, 1), b: New(4, 0), scale: 0, want: "1"},
		{name: "0.015_div_1_at_0", a: New(15, 3), b: New(1, 0), scale: 0, want: "0"},
		{name: "0.005_div_2_at_0", a: New(1, 2), b: New(2, 0), scale: 0, want: "0"},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got, err := Div(c.a, c.b, c.scale)
			require.NoError(t, err)
			assert.Equal(t, c.want, got.ToString())
		})
	}
}
