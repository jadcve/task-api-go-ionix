package enums

import "testing"

func TestIsValidTaskStatus(t *testing.T) {
	if !IsValidTaskStatus(TaskStatusAssigned) {
		t.Fatal("expected ASSIGNED to be valid")
	}
	if !IsValidTaskStatus(TaskStatusInProgress) {
		t.Fatal("expected IN_PROGRESS to be valid")
	}
	if !IsValidTaskStatus(TaskStatusCompleted) {
		t.Fatal("expected COMPLETED to be valid")
	}
	if !IsValidTaskStatus(TaskStatusCancelled) {
		t.Fatal("expected CANCELLED to be valid")
	}
	if IsValidTaskStatus(TaskStatus("UNKNOWN")) {
		t.Fatal("expected UNKNOWN to be invalid")
	}
}

func TestCanTransitionTaskStatus(t *testing.T) {
	if !CanTransitionTaskStatus(TaskStatusAssigned, TaskStatusInProgress) {
		t.Fatal("expected ASSIGNED -> IN_PROGRESS to be valid")
	}
	if !CanTransitionTaskStatus(TaskStatusAssigned, TaskStatusCancelled) {
		t.Fatal("expected ASSIGNED -> CANCELLED to be valid")
	}
	if !CanTransitionTaskStatus(TaskStatusInProgress, TaskStatusCompleted) {
		t.Fatal("expected IN_PROGRESS -> COMPLETED to be valid")
	}
	if !CanTransitionTaskStatus(TaskStatusInProgress, TaskStatusCancelled) {
		t.Fatal("expected IN_PROGRESS -> CANCELLED to be valid")
	}

	if CanTransitionTaskStatus(TaskStatusCompleted, TaskStatusCancelled) {
		t.Fatal("expected COMPLETED to have no outgoing transitions")
	}
	if CanTransitionTaskStatus(TaskStatusCancelled, TaskStatusAssigned) {
		t.Fatal("expected CANCELLED to have no outgoing transitions")
	}
	if CanTransitionTaskStatus(TaskStatusAssigned, TaskStatusCompleted) {
		t.Fatal("expected ASSIGNED -> COMPLETED to be invalid")
	}
	if CanTransitionTaskStatus(TaskStatus("UNKNOWN"), TaskStatusAssigned) {
		t.Fatal("expected UNKNOWN origin to be invalid")
	}
}
