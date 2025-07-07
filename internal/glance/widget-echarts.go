package glance

import (
	"context"
	"errors"
	"html/template"
	"io"
	"net/http"
	"time"
)

var echartsWidgetTemplate = mustParseTemplate("echarts.html", "widget-base.html")

type echartsWidget struct {
	widgetBase `yaml:",inline"`
	Height     string        `yaml:"height"`
	Theme      string        `yaml:"theme"`
	Data       template.JS   `yaml:"data"`
	DataURL    string        `yaml:"data-url"`
	cachedHTML template.HTML `yaml:"-"`
}

func (widget *echartsWidget) initialize() error {
	widget.withTitle("ECharts").withCacheDuration(10 * time.Minute).withError(nil)

	return nil
}

func (widget *echartsWidget) update(ctx context.Context) {
	if widget.DataURL == "" {
		return
	}
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(widget.DataURL)
	if !widget.canContinueUpdateAfterHandlingErr(err) {
        return
    }
	if err != nil {
		widget.withError(errors.New("failed to fetch data-url: " + err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		widget.withError(errors.New("non-200 response from data-url: " + resp.Status))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		widget.withError(errors.New("failed to read data-url body: " + err.Error()))
		return
	}

	widget.Data = template.JS(body)
}

func (widget *echartsWidget) Render() template.HTML {
	return widget.renderTemplate(widget, echartsWidgetTemplate)
}
