package notifier

// BuildStatus BuildStatus
type BuildStatus string

const (
	// BuildStatusUnknown Status of the build is unknown.
	BuildStatusUnknown BuildStatus = "STATUS_UNKNOWN"

	// BuildStatusQueued Build or step is queued; work has not yet begun.
	BuildStatusQueued BuildStatus = "QUEUED"

	// BuildStatusWorking Build or step is being executed.
	BuildStatusWorking BuildStatus = "WORKING"

	// BuildStatusSuccess Build or step finished successfully.
	BuildStatusSuccess BuildStatus = "SUCCESS"

	// BuildStatusFailure Build or step failed to complete successfully.
	BuildStatusFailure BuildStatus = "FAILURE"

	// BuildStatusInternalError Build or step failed due to an internal cause.
	BuildStatusInternalError BuildStatus = "INTERNAL_ERROR"

	// BuildStatusTimeout Build or step took longer than was allowed.
	BuildStatusTimeout BuildStatus = "TIMEOUT"

	// BuildStatusCancelled Build or step was canceled by a user.
	BuildStatusCancelled BuildStatus = "CANCELLED"
)

var statuses = map[BuildStatus]bool{
	BuildStatusUnknown:       true,
	BuildStatusQueued:        true,
	BuildStatusWorking:       true,
	BuildStatusSuccess:       true,
	BuildStatusFailure:       true,
	BuildStatusInternalError: true,
	BuildStatusTimeout:       true,
	BuildStatusCancelled:     true,
}

func isValidBuildStatus(status string) bool {
	if value, found := statuses[BuildStatus(status)]; found {
		return value
	}

	return false
}
