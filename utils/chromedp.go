package utils

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

//ColumnPrintToPDF print pdf
func ColumnPrintToPDF(aid string, filename string, cookies map[string]string) error {
	var buf []byte
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Emulate(device.IPhone7),
			enableLifeCycleEvents(),
			setCookies(cookies),
			navigateAndWaitFor(`https://www.dedao.cn/article/`+aid, "networkIdle"),
			chromedp.ActionFunc(func(ctx context.Context) error {
				s := `
					document.querySelector('.iget-header-main').style.display='none';
					document.querySelector('.article-body-wrap').style.margin="0px 50px";
				`
				_, exp, err := runtime.Evaluate(s).Do(ctx)
				if err != nil {
					return err
				}

				if exp != nil {
					return exp
				}

				return nil
			}),
			chromedp.ActionFunc(func(ctx context.Context) error {
				s := `
					var buttons = document.getElementsByTagName('button');
					for (var i = 0; i < buttons.length; ++i){
						if(buttons[i].innerText === "展开侧边栏" || buttons[i].innerText === "设置文本"){
							buttons[i].parentNode.style.display="none";
							break;
						}
					}
					var asides = document.getElementsByTagName('aside');
					for (var i = 0; i < asides.length; ++i){
						asides[i].style.display="none";
					}
				`
				_, exp, err := runtime.Evaluate(s).Do(ctx)
				if err != nil {
					return err
				}

				if exp != nil {
					return exp
				}

				return nil
			}),

			chromedp.ActionFunc(func(ctx context.Context) error {
				// time.Sleep(time.Second * 5)
				var err error
				buf, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
				return err
			}),
		},
	)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buf, 0644)
}

func setCookies(cookies map[string]string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))

		for key, value := range cookies {
			err := network.SetCookie(key, value).WithExpires(&expr).WithDomain(".dedao.cn").WithHTTPOnly(true).Do(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func enableLifeCycleEvents() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		err := page.Enable().Do(ctx)
		if err != nil {
			return err
		}

		return page.SetLifecycleEventsEnabled(true).Do(ctx)
	}
}

func navigateAndWaitFor(url string, eventName string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		_, _, _, err := page.Navigate(url).Do(ctx)
		if err != nil {
			return err
		}

		return waitFor(ctx, eventName)
	}
}

// waitFor blocks until eventName is received.
// Examples of events you can wait for:
//     init, DOMContentLoaded, firstPaint,
//     firstContentfulPaint, firstImagePaint,
//     firstMeaningfulPaintCandidate,
//     load, networkAlmostIdle, firstMeaningfulPaint, networkIdle
//
// This is not super reliable, I've already found incidental cases where
// networkIdle was sent before load. It's probably smart to see how
// puppeteer implements this exactly.
func waitFor(ctx context.Context, eventName string) error {
	ch := make(chan struct{})
	cctx, cancel := context.WithCancel(ctx)
	chromedp.ListenTarget(cctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *page.EventLifecycleEvent:
			if e.Name == eventName {
				cancel()
				close(ch)
			}
		}
	})

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
