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
// requested scale is smaller than the dividend's own scale. Since
// KT-06 the extra division by 10^(-shift) happens on the denominator
// side, so the dividend's low-order digits survive until the rounding
// step: 2.5 / 1 at scale 0 answers 3 (half away from zero on the exact
// quotient), not the pre-truncated 2.
func TestDiv_ResultScaleBelowDividendScale(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		a, b  Decimal
		scale int
		want  string
	}{
		{name: "2.5_div_1_at_0", a: New(25, 1), b: New(1, 0), scale: 0, want: "3"},
		{name: "2.5_div_4_at_0", a: New(25, 1), b: New(4, 0), scale: 0, want: "1"},
		{name: "2.5_div_8_at_0", a: New(25, 1), b: New(8, 0), scale: 0, want: "0"},
		{name: "0.015_div_1_at_0", a: New(15, 3), b: New(1, 0), scale: 0, want: "0"},
		{name: "0.005_div_2_at_0", a: New(1, 2), b: New(2, 0), scale: 0, want: "0"},
		{name: "neg_2.5_div_1_at_0", a: New(-25, 1), b: New(1, 0), scale: 0, want: "-3"},
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

// TestDiv_OverflowFailsClosed pins the KT-06 overflow semantics: any
// scaling step that would leave the int64 range returns ErrOverflow
// instead of silently truncating the numerator (a precision error) or
// skipping the multiply (a magnitude error — the old code's `break`
// made the quotient smaller by that power of ten).
func TestDiv_OverflowFailsClosed(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		a, b  Decimal
		scale int
	}{
		// shift > 0: numerator one past math.MaxInt64/10 cannot take the
		// final ×10 without leaving int64.
		{name: "numerator_over_maxScaled10", a: Decimal{value: 922337203685477581, scale: 0}, b: New(1, 0), scale: 1},
		// Same boundary from below must still succeed, not error.
		{name: "numerator_at_maxScaled10_boundary", a: Decimal{value: 922337203685477580, scale: 0}, b: New(1, 0), scale: 1},
		// Negative operand at the same boundary.
		{name: "numerator_negative_overflow", a: Decimal{value: -922337203685477581, scale: 0}, b: New(1, 0), scale: 1},
		// The old silent-`break` shape: numerator crosses 1<<62 mid-loop,
		// the remaining multiplies were skipped and the quotient came back
		// wrong by a power of ten with no error.
		{name: "old_break_shape_now_errors", a: Decimal{value: 5_000_000_000, scale: 0}, b: New(1, 0), scale: 18},
		// shift < 0: the scaling moved to the denominator, so an oversized
		// divisor fails closed too (old code truncated the numerator here
		// and never overflowed).
		{name: "denominator_overflow_on_negative_shift", a: Decimal{value: 1, scale: 18}, b: Decimal{value: 922337203685477581, scale: 0}, scale: 0},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got, err := Div(c.a, c.b, c.scale)
			if c.name == "numerator_at_maxScaled10_boundary" {
				require.NoError(t, err)
				assert.Equal(t, "922337203685477580.0", got.ToString())
				return
			}
			require.ErrorIs(t, err, ErrOverflow)
			assert.Equal(t, Zero, got)
		})
	}
}
