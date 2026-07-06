package enums

import "testing"

func TestIsValidTaskStatus(t *testing.T) {
	if !IsValidTaskStatus(TaskStatusAssigned) {
		t.Fatal("expected ASSIGNED to be valid")
	}
	if !IsValidTaskStatus(TaskStatusStarted) {
		t.Fatal("expected STARTED to be valid")
	}
	if !IsValidTaskStatus(TaskStatusWaiting) {
		t.Fatal("expected WAITING to be valid")
	}
	if !IsValidTaskStatus(TaskStatusCompletedSuccess) {
		t.Fatal("expected COMPLETED_SUCCESS to be valid")
	}
	if !IsValidTaskStatus(TaskStatusCompletedError) {
		t.Fatal("expected COMPLETED_ERROR to be valid")
	}
	if IsValidTaskStatus(TaskStatus("IN_PROGRESS")) {
		t.Fatal("expected IN_PROGRESS to be invalid")
	}
	if IsValidTaskStatus(TaskStatus("COMPLETED")) {
		t.Fatal("expected COMPLETED to be invalid")
	}
	if IsValidTaskStatus(TaskStatus("CANCELLED")) {
		t.Fatal("expected CANCELLED to be invalid")
	}
	if IsValidTaskStatus(TaskStatus("UNKNOWN")) {
		t.Fatal("expected UNKNOWN to be invalid")
	}
}

func TestCanTransitionTaskStatus(t *testing.T) {
	if !CanTransitionTaskStatus(TaskStatusAssigned, TaskStatusStarted) {
		t.Fatal("expected ASSIGNED -> STARTED to be valid")
	}
	if !CanTransitionTaskStatus(TaskStatusStarted, TaskStatusWaiting) {
		t.Fatal("expected STARTED -> WAITING to be valid")
	}
	if !CanTransitionTaskStatus(TaskStatusWaiting, TaskStatusStarted) {
		t.Fatal("expected WAITING -> STARTED to be valid")
	}
	if !CanTransitionTaskStatus(TaskStatusStarted, TaskStatusCompletedSuccess) {
		t.Fatal("expected STARTED -> COMPLETED_SUCCESS to be valid")
	}
	if !CanTransitionTaskStatus(TaskStatusStarted, TaskStatusCompletedError) {
		t.Fatal("expected STARTED -> COMPLETED_ERROR to be valid")
	}

	if CanTransitionTaskStatus(TaskStatusAssigned, TaskStatusWaiting) {
		t.Fatal("expected ASSIGNED -> WAITING to be invalid")
	}
	if CanTransitionTaskStatus(TaskStatusWaiting, TaskStatusCompletedSuccess) {
		t.Fatal("expected WAITING -> COMPLETED_SUCCESS to be invalid")
	}
	if CanTransitionTaskStatus(TaskStatusCompletedSuccess, TaskStatusStarted) {
		t.Fatal("expected COMPLETED_SUCCESS to have no outgoing transitions")
	}
	if CanTransitionTaskStatus(TaskStatusCompletedError, TaskStatusStarted) {
		t.Fatal("expected COMPLETED_ERROR to have no outgoing transitions")
	}
	if CanTransitionTaskStatus(TaskStatus("UNKNOWN"), TaskStatusAssigned) {
		t.Fatal("expected UNKNOWN origin to be invalid")
	}
}
