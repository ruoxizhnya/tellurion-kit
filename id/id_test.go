package id

import (
	"regexp"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
)

// uuidV7Pattern is the standard 8-4-4-4-12 form. We re-validate
// via a regex rather than the package's own parser so the
// test pin doesn't depend on uuid package internals.
var uuidV7Pattern = regexp.MustCompile(
	`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestJobID_IsValidV7(t *testing.T) {
	id := JobID()
	if !uuidV7Pattern.MatchString(id) {
		t.Errorf("JobID must be canonical UUIDv7 form; got %q", id)
	}

	u, err := uuid.Parse(id)
	if err != nil {
		t.Fatalf("uuid.Parse on JobID output: %v", err)
	}
	if v := u.Version(); v != 7 {
		t.Errorf("JobID version = %d, want 7 (UUIDv7)", v)
	}
}

func TestSubscriptionID_IsValidV7(t *testing.T) {
	id := SubscriptionID()
	u, err := uuid.Parse(id)
	if err != nil {
		t.Fatalf("uuid.Parse on SubscriptionID output: %v", err)
	}
	if v := u.Version(); v != 7 {
		t.Errorf("SubscriptionID version = %d, want 7", v)
	}
}

// TestJobIDsAreTimeOrdered 是 v7 的核心性质：字典序 == 时间序。
func TestJobIDsAreTimeOrdered(t *testing.T) {
	// 生成 5 个 ID，间隔 2ms 让时间戳位真正变化；同一毫秒内的 ID
	// 共享时间戳前缀，排序断言会退化。
	ids := make([]string, 5)
	for i := range ids {
		ids[i] = JobID()
		if i < len(ids)-1 {
			time.Sleep(2 * time.Millisecond)
		}
	}
	clone := append([]string(nil), ids...)
	sort.Strings(clone)
	for i := range ids {
		if ids[i] != clone[i] {
			t.Errorf("UUIDv7 ordering broken at index %d: generated=%v sorted=%v",
				i, ids, clone)
			break
		}
	}
}

// TestJobID_Unique：1 万个 ID 无碰撞。
func TestJobID_Unique(t *testing.T) {
	seen := make(map[string]struct{}, 10_000)
	for i := 0; i < 10_000; i++ {
		id := JobID()
		if _, dup := seen[id]; dup {
			t.Fatalf("UUIDv7 collision after %d draws: %s", i, id)
		}
		seen[id] = struct{}{}
	}
}

// TestJobID_ContainsTimestamp：前 48 bit 必须是 unix 毫秒时间戳
// （RFC 9562 §5.7）。手工解码字节，不依赖 uuid 包的 Time helper。
func TestJobID_ContainsTimestamp(t *testing.T) {
	before := time.Now()
	id := JobID()
	after := time.Now()

	u, err := uuid.Parse(id)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	ts := int64(u[0])<<40 | int64(u[1])<<32 |
		int64(u[2])<<24 | int64(u[3])<<16 |
		int64(u[4])<<8 | int64(u[5])
	tsTime := time.UnixMilli(ts)
	if tsTime.Before(before.Add(-time.Second)) || tsTime.After(after.Add(time.Second)) {
		t.Errorf("UUIDv7 timestamp %v outside [%v, %v]", tsTime, before, after)
	}
}

// TestHelpers_UseTheSamePrimitive：两个 helper 必须都产出 v7，
// 否则同一个包里的 ID 会分版本（排序收益只对其中一半成立）。
func TestHelpers_UseTheSamePrimitive(t *testing.T) {
	for name, fn := range map[string]func() string{
		"JobID":          JobID,
		"SubscriptionID": SubscriptionID,
	} {
		id := fn()
		u, err := uuid.Parse(id)
		if err != nil {
			t.Errorf("%s: parse: %v", name, err)
			continue
		}
		if v := u.Version(); v != 7 {
			t.Errorf("%s version = %d, want 7 (all helpers must agree on v7)", name, v)
		}
	}
}
