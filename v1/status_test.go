package notifier

import "testing"

func Test_isValidBuildStatus(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert(t, true, isValidBuildStatus("SUCCESS"))
		assert(t, true, isValidBuildStatus("INTERNAL_ERROR"))
	})

	t.Run("invalid", func(t *testing.T) {
		assert(t, false, isValidBuildStatus("SUCCESSS"))
	})
}
