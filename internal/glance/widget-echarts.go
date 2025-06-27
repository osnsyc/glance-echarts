package glance

import (
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
	widget.withTitle("ECharts").withError(nil)

	if widget.DataURL != "" && widget.Data == "" {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(widget.DataURL)
		if err != nil {
			return errors.New("failed to fetch data-url: " + err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return errors.New("non-200 response from data-url: " + resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.New("failed to read data-url body: " + err.Error())
		}

		widget.Data = template.JS(body)
	}

	widget.cachedHTML = widget.renderTemplate(widget, echartsWidgetTemplate)

	return nil
}

func (widget *echartsWidget) Render() template.HTML {
	return widget.cachedHTML
}
