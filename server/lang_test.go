package server

import (
	"flag"
	"testing"
)

func TestWebURL(t *testing.T) {
	original := *httpPort
	t.Cleanup(func() { _ = flag.Set("http", original) })

	for _, test := range []struct {
		address string
		want    string
		wantErr bool
	}{
		{address: ":8118", want: "http://localhost:8118"},
		{address: ":9123", want: "http://localhost:9123"},
		{address: "127.0.0.1:9191", want: "http://127.0.0.1:9191"},
		{address: "0.0.0.0:9191", want: "http://localhost:9191"},
		{address: "[::]:9191", want: "http://localhost:9191"},
		{address: "", wantErr: true},
	} {
		t.Run(test.address, func(t *testing.T) {
			if err := flag.Set("http", test.address); err != nil {
				t.Fatal(err)
			}
			url, err := WebURL()
			if (err != nil) != test.wantErr || url != test.want {
				t.Errorf("WebURL() = %q, %v; want %q, error: %v", url, err, test.want, test.wantErr)
			}
		})
	}
}
