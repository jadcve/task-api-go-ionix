package enums

type TaskStatus string

const (
	TaskStatusAssigned   TaskStatus = "ASSIGNED"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusCancelled  TaskStatus = "CANCELLED"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleExecutor UserRole = "EXECUTOR"
	UserRoleAuditor  UserRole = "AUDITOR"
)

func IsValidTaskStatus(status TaskStatus) bool {
	switch status {
	case TaskStatusAssigned, TaskStatusInProgress, TaskStatusCompleted, TaskStatusCancelled:
		return true
	default:
		return false
	}
}

func CanTransitionTaskStatus(from TaskStatus, to TaskStatus) bool {
	switch from {
	case TaskStatusAssigned:
		return to == TaskStatusInProgress || to == TaskStatusCancelled
	case TaskStatusInProgress:
		return to == TaskStatusCompleted || to == TaskStatusCancelled
	case TaskStatusCompleted, TaskStatusCancelled:
		return false
	default:
		return false
	}
}
