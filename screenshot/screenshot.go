package screenshot

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func Take(targetURL string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(1200, 739),
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Headless,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body"),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		return nil, err
	}

	return buf, nil
}
