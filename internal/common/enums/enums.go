package enums

type TaskStatus string

const (
	TaskStatusAssigned         TaskStatus = "ASSIGNED"
	TaskStatusStarted          TaskStatus = "STARTED"
	TaskStatusWaiting          TaskStatus = "WAITING"
	TaskStatusCompletedSuccess TaskStatus = "COMPLETED_SUCCESS"
	TaskStatusCompletedError   TaskStatus = "COMPLETED_ERROR"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleExecutor UserRole = "EXECUTOR"
	UserRoleAuditor  UserRole = "AUDITOR"
)

func IsValidTaskStatus(status TaskStatus) bool {
	switch status {
	case TaskStatusAssigned, TaskStatusStarted, TaskStatusWaiting, TaskStatusCompletedSuccess, TaskStatusCompletedError:
		return true
	default:
		return false
	}
}

func CanTransitionTaskStatus(from TaskStatus, to TaskStatus) bool {
	switch from {
	case TaskStatusAssigned:
		return to == TaskStatusStarted
	case TaskStatusStarted:
		return to == TaskStatusWaiting || to == TaskStatusCompletedSuccess || to == TaskStatusCompletedError
	case TaskStatusWaiting:
		return to == TaskStatusStarted
	case TaskStatusCompletedSuccess, TaskStatusCompletedError:
		return false
	default:
		return false
	}
}
