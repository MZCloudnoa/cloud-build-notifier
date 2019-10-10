package notifier

import (
	"io/ioutil"
	"strings"
	"testing"
)

func assert(t *testing.T, wanted, result interface{}) {
	if wanted != result {
		t.Fatal("assert fail :", result)
	}
}

func fatalIf(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_NewNotifier(t *testing.T) {
	t.Run("call with nil", func(t *testing.T) {
		assert(t, true, NewNotifier(nil) != nil)
	})

	t.Run("call with options", func(t *testing.T) {
		options := &Options{}
		notifier := NewNotifier(options)
		assert(t, true, notifier != nil)
		assert(t, true, notifier.options == options)
	})
}

func Test_Notifier_HandlePubSub(t *testing.T) {
	notifier := NewNotifier(nil)

	dir := "../misc/example/"
	ext := ".json"

	t.Run("status", func(t *testing.T) {
		files, err := ioutil.ReadDir(dir)
		assert(t, nil, err)

		for _, file := range files {
			filename := file.Name()
			if !strings.HasSuffix(filename, ext) {
				continue
			}

			status := filename[:len(filename)-len(ext)]
			if status != "QUEUED" {
				continue
			}

			t.Run(status, func(t *testing.T) {
				buf, err := ioutil.ReadFile(dir + filename)
				assert(t, nil, err)

				err = notifier.HandlePubSub(buf)
				assert(t, nil, err)
			})
		}
	})
}
