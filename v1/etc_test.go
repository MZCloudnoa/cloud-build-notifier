package notifier

import "testing"

func Test_executeTemplate(t *testing.T) {
	data := map[string]interface{}{
		"str": "str value",
	}

	t.Run("simple", func(t *testing.T) {
		assert(t, "", mustExecuteTemplate("", data))
		assert(t, "1", mustExecuteTemplate("1", data))
		assert(t, data["str"], mustExecuteTemplate("{{.str}}", data))
	})
}
