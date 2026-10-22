// Package statistics provides reusable numerical helpers: descriptive
// statistics, correlation, percentile/rank, and least-squares regression.
//
// Design rules:
//
//  1. Pure functions. No globals, no I/O, no logging. Each function is
//     deterministic and side-effect-free so callers can use them in hot
//     loops without a logger dependency.
//
//  2. NaN / Inf policy. Functions over float64 slices never silently drop
//     non-finite values; callers that want filtering use the *Finite
//     variants (e.g. MeanFinite). Silent dropping turns "the input had a
//     NaN" into "the output looks plausible" — the failure is then far
//     from its cause.
//
//  3. Population vs sample stddev. The function names make the denominator
//     explicit: PopulationStdDev divides by n, SampleStdDev divides by n-1.
//     There is no ambiguous "stddev" overload.
//
//  4. Empty / single-element slices. Return 0 (not NaN) so callers can
//     `if s == 0 { skip }` without panicking. Functions that cannot be
//     computed at all (e.g. Pearson with constant input) may return NaN to
//     signal "not computable" — see their doc comments.
//
// Scope note: domain calibers (log-returns, annualization factors) are
// deliberately NOT here — they encode how a particular caller measures
// time or growth. This package stays pure math.
package statistics
