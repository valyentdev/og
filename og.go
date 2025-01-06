package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/valyentcloud/og/screenshot"
	"github.com/valyentcloud/og/storage"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	driver, err := storage.NewS3(storage.S3Config{
		EndpointURL:     os.Getenv("AWS_S3_ENDPOINT_URL"),
		Region:          os.Getenv("AWS_S3_REGION"),
		AccessKeyID:     os.Getenv("AWS_S3_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_S3_SECRET_ACCESS_KEY"),
		Bucket:          os.Getenv("AWS_S3_BUCKET_NAME"),
	})
	if err != nil {
		slog.Error("failed to intitialize storage driver", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		path := os.Getenv("BASE_WEBSITE_URL") + r.URL.Path

		if image, err := driver.Get(r.URL.Path); err == nil {
			w.Header().Set("Content-Type", "image/png")
			w.Write(image)
		} else {
			image, err := screenshot.Take(path)
			if err != nil {
				slog.Error("failed to take screenshot", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to take screenshot."))
				return
			}
			if err := driver.Upload(r.URL.Path, image); err != nil {
				slog.Error("failed to upload screenshot", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to upload screenshot."))
				return
			}
			w.Header().Set("Content-Type", "image/png")
			w.Write(image)
		}
	})

	slog.Info("server starting", "addr", os.Getenv("HTTP_ADDR"))
	if err := http.ListenAndServe(os.Getenv("HTTP_ADDR"), mux); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
