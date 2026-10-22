package errors

import (
	stderrors "errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- Error() 文案 --------------------------------------------------------

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		want string
	}{
		{
			name: "with cause includes cause verbatim",
			err:  New(ErrCodeInvalidInput, "bad field").WithOperation("parse").WithCause(stderrors.New("EOF")),
			want: "[INVALID_INPUT] bad field: EOF (op=parse)",
		},
		{
			name: "without cause omits the cause segment",
			err:  New(ErrCodeNotFound, "strategy not found: s1").WithOperation("lookup"),
			want: "[NOT_FOUND] strategy not found: s1 (op=lookup)",
		},
		{
			name: "empty operation still renders the op slot",
			err:  New(ErrCodeInternal, "boom"),
			want: "[INTERNAL] boom (op=)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}

// ---- Is 按分类码相等（不是实例身份）-------------------------------------

func TestAppError_Is(t *testing.T) {
	tests := []struct {
		name   string
		err    *AppError
		target error
		want   bool
	}{
		{
			name:   "same code different message matches",
			err:    New(ErrCodeNotFound, "strategy not found: s1"),
			target: New(ErrCodeNotFound, "another message entirely"),
			want:   true,
		},
		{
			name:   "different code does not match",
			err:    New(ErrCodeNotFound, "x"),
			target: New(ErrCodeConflict, "x"),
			want:   false,
		},
		{
			name:   "nil target does not match",
			err:    New(ErrCodeInternal, "x"),
			target: nil,
			want:   false,
		},
		{
			name:   "non-AppError target does not match",
			err:    New(ErrCodeInternal, "x"),
			target: stderrors.New("INTERNAL"),
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Is(tt.target))
		})
	}
}

func TestAppError_IsThroughStdlibErrors(t *testing.T) {
	err := NotFound("strategy", "s1")
	// 同码不同 message 通过 stdlib errors.Is 也应命中 —— 这正是 usage 的
	// 「按分类码分支」语义，见 docs/api-schema.md §2。
	assert.True(t, stderrors.Is(err, New(ErrCodeNotFound, "whatever")))
	assert.False(t, stderrors.Is(err, New(ErrCodeConflict, "whatever")))
}

// ---- Unwrap 与 %w 链 ----------------------------------------------------

func TestAppError_Unwrap(t *testing.T) {
	cause := stderrors.New("db down")
	assert.Same(t, cause, Wrap(cause, ErrCodeUnavailable, "db unavailable", "connect").Unwrap())
	assert.Nil(t, New(ErrCodeInternal, "no cause").Unwrap())
}

func TestErrorChain_PassesThroughToSentinel(t *testing.T) {
	sentinel := stderrors.New("sentinel")
	inner := Wrap(sentinel, ErrCodeInternal, "inner", "step1")
	outer := fmt.Errorf("outer: %w", inner)

	// stdlib errors.Is 必须能穿过 AppError.Unwrap 找到哨兵。
	assert.True(t, stderrors.Is(outer, sentinel))
	// errors.As 必须能从非 AppError 外层剥出 *AppError。
	var appErr *AppError
	require.True(t, stderrors.As(outer, &appErr))
	assert.Equal(t, ErrCodeInternal, appErr.Code)
	assert.Equal(t, "step1", appErr.Operation)
	// 再往下穿一层仍可达到哨兵。
	assert.Same(t, sentinel, stderrors.Unwrap(appErr))
}

// ---- WithOperation / WithCause 的返回值语义 ------------------------------

func TestAppError_WithOperation_MutatesAndReturnsSamePointer(t *testing.T) {
	err := New(ErrCodeInvalidInput, "bad")
	same := err.WithOperation("parse")
	assert.Same(t, err, same)
	assert.Equal(t, "parse", err.Operation)
}

func TestAppError_WithCause_MutatesAndReturnsSamePointer(t *testing.T) {
	err := New(ErrCodeInternal, "boom")
	cause := stderrors.New("root")
	same := err.WithCause(cause)
	assert.Same(t, err, same)
	assert.Same(t, cause, err.Cause)
	// 原地修改必须反映到同一实例的文案上。
	assert.Contains(t, err.Error(), "boom: root")
}

// ---- ToMap --------------------------------------------------------------

func TestAppError_ToMap(t *testing.T) {
	err := New(ErrCodeDataQuality, "stale").WithOperation("etl")

	got := err.ToMap()
	assert.Len(t, got, 3, "ToMap must expose exactly code/message/op")
	assert.Equal(t, ErrCodeDataQuality, got["code"])
	assert.Equal(t, "stale", got["message"])
	assert.Equal(t, "etl", got["op"])
	assert.NotContains(t, got, "stack_trace", "stack trace must not leak into ToMap")
}

// ---- 构造器 --------------------------------------------------------------

func TestConstructors(t *testing.T) {
	cause := stderrors.New("root")
	tests := []struct {
		name      string
		err       *AppError
		wantCode  ErrorCode
		wantMsg   string
		wantOp    string
		wantCause error
	}{
		{
			name:     "InvalidInput",
			err:      InvalidInput("bad field", "parse"),
			wantCode: ErrCodeInvalidInput,
			wantMsg:  "bad field",
			wantOp:   "parse",
		},
		{
			name:     "NotFound",
			err:      NotFound("strategy", "s1"),
			wantCode: ErrCodeNotFound,
			wantMsg:  "strategy not found: s1",
		},
		{
			name:      "Internal",
			err:       Internal(cause, "run"),
			wantCode:  ErrCodeInternal,
			wantMsg:   "internal error",
			wantOp:    "run",
			wantCause: cause,
		},
		{
			name:      "Unavailable",
			err:       Unavailable("tushare", cause),
			wantCode:  ErrCodeUnavailable,
			wantMsg:   "tushare service unavailable",
			wantOp:    "external_call",
			wantCause: cause,
		},
		{
			name:      "Timeout",
			err:       Timeout("fetch", cause),
			wantCode:  ErrCodeTimeout,
			wantMsg:   "operation timeout: fetch",
			wantOp:    "fetch",
			wantCause: cause,
		},
		{
			name:     "DataQuality",
			err:      DataQuality("stale quote", "etl"),
			wantCode: ErrCodeDataQuality,
			wantMsg:  "stale quote",
			wantOp:   "etl",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.err)
			assert.Equal(t, tt.wantCode, tt.err.Code)
			assert.Equal(t, tt.wantMsg, tt.err.Message)
			assert.Equal(t, tt.wantOp, tt.err.Operation)
			if tt.wantCause != nil {
				assert.Same(t, tt.wantCause, tt.err.Cause)
			} else {
				assert.Nil(t, tt.err.Cause)
			}
			assert.NotEmpty(t, tt.err.StackTrace)
		})
	}
}

func TestInternal_NilCauseReturnsNil(t *testing.T) {
	// Internal 复用 Wrap 的 nil 语义：没有 cause 就没有错误。
	assert.Nil(t, Internal(nil, "run"))
	assert.Nil(t, Unavailable("tushare", nil))
	assert.Nil(t, Timeout("fetch", nil))
}

// ---- Wrap / Wrapf 的 nil 语义与格式化 ------------------------------------

func TestWrap_NilErrorReturnsNil(t *testing.T) {
	assert.Nil(t, Wrap(nil, ErrCodeInternal, "msg", "op"))
}

func TestWrapf_NilErrorReturnsNil(t *testing.T) {
	assert.Nil(t, Wrapf(nil, ErrCodeInternal, "op", "msg %d", 1))
}

func TestWrapf_FormatsMessage(t *testing.T) {
	cause := stderrors.New("root")
	err := Wrapf(cause, ErrCodeInvalidInput, "parse", "bad %s: %d", "field", 3)
	require.NotNil(t, err)
	assert.Equal(t, "bad field: 3", err.Message)
	assert.Same(t, cause, err.Cause)
	assert.Equal(t, "parse", err.Operation)
	assert.Equal(t, ErrCodeInvalidInput, err.Code)
}

// ---- 堆栈捕获 ------------------------------------------------------------

func TestNew_CapturesStack(t *testing.T) {
	stack := New(ErrCodeInternal, "boom").StackTrace
	require.NotEmpty(t, stack)
	assert.True(t, strings.HasPrefix(stack, "goroutine "), "stack must start with the goroutine header, got %q", stack[:min(40, len(stack))])
	assert.Contains(t, stack, "captureStack", "stack must include the capture site")
}

func TestCaptureStack_NonEmpty(t *testing.T) {
	assert.NotEmpty(t, captureStack())
}
