package domain

import (
	"testing"
	"time"
)

func TestTaskIsExpired(t *testing.T) {
	now := time.Date(2026, time.July, 4, 12, 0, 0, 0, time.UTC)

	taskExpired := Task{DueDate: now.Add(-time.Minute)}
	if !taskExpired.IsExpired(now) {
		t.Fatal("expected task with past due date to be expired")
	}

	taskNotExpired := Task{DueDate: now.Add(time.Minute)}
	if taskNotExpired.IsExpired(now) {
		t.Fatal("expected task with future due date to not be expired")
	}

	taskZeroDate := Task{}
	if taskZeroDate.IsExpired(now) {
		t.Fatal("expected task with zero due date to not be expired")
	}
}
