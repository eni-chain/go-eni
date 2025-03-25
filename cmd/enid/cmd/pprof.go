package cmd

import "net/http"

func init() {
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			return
		}
	}()
}
