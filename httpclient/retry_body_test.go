package httpclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unicode/utf8"

	apperrors "github.com/ruoxizhnya/tellurion-kit/errors"
)

// TestDo_RetriesWithBody_PreservesBodyOnEveryAttempt pins KT-04's first
// defect: a request body handed to the first attempt used to be consumed to
// EOF, so every retry sent an empty body. The server records what it actually
// received on each attempt.
func TestDo_RetriesWithBody_PreservesBodyOnEveryAttempt(t *testing.T) {
	var mu sync.Mutex
	var received [][]byte
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		mu.Lock()
		received = append(received, data)
		mu.Unlock()

		if atomic.AddInt32(&attempts, 1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	c := New(server.URL, 10*time.Second, 1) // two attempts, one 1s backoff
	payload := map[string]int{"n": 1}
	resp, err := c.Do(context.Background(), Request{Method: http.MethodPost, Path: "/submit", Body: payload})
	if err != nil {
		t.Fatalf("expected no error after retry, got: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Fatalf("expected 2 attempts, got %d", len(received))
	}
	want := []byte(`{"n":1}`)
	for i, got := range received {
		if !bytes.Equal(got, want) {
			t.Errorf("attempt %d body = %q, want %q", i+1, got, want)
		}
	}
}

// TestDo_4xxReturnsAppError pins KT-04's second defect: a 4xx used to come
// back with err == nil, so `if err != nil` callers treated it as success.
func TestDo_4xxReturnsAppError(t *testing.T) {
	const upstreamBody = `{"error":"bad things"}`

	cases := []struct {
		status int
		want   apperrors.ErrorCode
	}{
		{http.StatusBadRequest, apperrors.ErrCodeInvalidInput},
		{http.StatusUnauthorized, apperrors.ErrCodePermission},
		{http.StatusForbidden, apperrors.ErrCodePermission},
		{http.StatusNotFound, apperrors.ErrCodeNotFound},
		{http.StatusRequestTimeout, apperrors.ErrCodeTimeout},
		{http.StatusConflict, apperrors.ErrCodeConflict},
		{http.StatusUnprocessableEntity, apperrors.ErrCodeInvalidInput},
		{http.StatusTeapot, apperrors.ErrCodeInvalidInput}, // unmapped 4xx falls back to invalid input
	}

	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.status), func(t *testing.T) {
			var attempts int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&attempts, 1)
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(upstreamBody))
			}))
			defer server.Close()

			c := New(server.URL, 10*time.Second, 2)
			resp, err := c.Do(context.Background(), Request{Method: http.MethodGet, Path: "/test"})

			if resp != nil {
				t.Errorf("expected nil response for %d, got %+v", tc.status, resp)
			}
			if err == nil {
				t.Fatalf("expected an error for status %d, got nil", tc.status)
			}

			var appErr *apperrors.AppError
			if !errors.As(err, &appErr) {
				t.Fatalf("expected *errors.AppError, got %T: %v", err, err)
			}
			if appErr.Code != tc.want {
				t.Errorf("status %d: expected code %s, got %s", tc.status, tc.want, appErr.Code)
			}
			if !errors.Is(err, apperrors.New(tc.want, "does not matter")) {
				t.Errorf("status %d: expected errors.Is to match by code %s", tc.status, tc.want)
			}
			if appErr.Operation != "GET /test" {
				t.Errorf("status %d: expected operation %q, got %q", tc.status, "GET /test", appErr.Operation)
			}
			if !strings.Contains(appErr.Message, fmt.Sprintf("unexpected HTTP status %d", tc.status)) {
				t.Errorf("status %d: expected message to name the status, got %q", tc.status, appErr.Message)
			}
			if !strings.Contains(appErr.Message, "bad things") {
				t.Errorf("status %d: expected message to quote the upstream body, got %q", tc.status, appErr.Message)
			}
			if got := atomic.LoadInt32(&attempts); got != 1 {
				t.Errorf("status %d: expected 1 attempt (4xx is not retried), got %d", tc.status, got)
			}
		})
	}
}

// TestDo_4xxTruncatesUpstreamBody pins the bounded diagnostic quote: the
// upstream body is capped at 200 characters (counted in runes) so a chatty
// upstream cannot inflate the error.
func TestDo_4xxTruncatesUpstreamBody(t *testing.T) {
	longBody := strings.Repeat("错", 250)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(longBody))
	}))
	defer server.Close()

	c := New(server.URL, 10*time.Second, 0)
	_, err := c.Do(context.Background(), Request{Method: http.MethodPost, Path: "/submit"})
	if err == nil {
		t.Fatal("expected an error for 400")
	}
	msg := err.Error()
	if !utf8.ValidString(msg) {
		t.Errorf("error message is not valid UTF-8: %q", msg)
	}
	if !strings.Contains(msg, strings.Repeat("错", 200)+"...") {
		t.Errorf("expected the body to be truncated to 200 runes with an ellipsis, got %q", msg)
	}
	if strings.Contains(msg, longBody) {
		t.Error("expected the full upstream body not to appear in the message")
	}
}

// TestDo_2xxIsStillSuccess makes sure the new fail-closed 4xx handling does
// not accidentally turn successful responses into errors.
func TestDo_2xxIsStillSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"created":true}`))
	}))
	defer server.Close()

	c := New(server.URL, 10*time.Second, 2)
	resp, err := c.Do(context.Background(), Request{Method: http.MethodPost, Path: "/create", Body: map[string]int{"n": 1}})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}
