package logparser

import "testing"

func TestDefaultOptions(t *testing.T) {
	o := defaultOptions()
	if o.timeField != "time" {
		t.Errorf("default timeField: want %q, got %q", "time", o.timeField)
	}
	if o.timeLayout == "" {
		t.Error("default timeLayout should not be empty")
	}
}

func TestWithTimeField(t *testing.T) {
	o := defaultOptions()
	WithTimeField("ts")(&o)
	if o.timeField != "ts" {
		t.Errorf("WithTimeField: want %q, got %q", "ts", o.timeField)
	}
}

func TestWithTimeField_Empty(t *testing.T) {
	o := defaultOptions()
	WithTimeField("")(&o)
	// Empty string should not override the default.
	if o.timeField != "time" {
		t.Errorf("WithTimeField empty: default should be preserved, got %q", o.timeField)
	}
}

func TestWithTimeLayout(t *testing.T) {
	o := defaultOptions()
	layout := "2006/01/02 15:04:05"
	WithTimeLayout(layout)(&o)
	if o.timeLayout != layout {
		t.Errorf("WithTimeLayout: want %q, got %q", layout, o.timeLayout)
	}
}

func TestWithTimeLayout_Empty(t *testing.T) {
	o := defaultOptions()
	original := o.timeLayout
	WithTimeLayout("")(&o)
	if o.timeLayout != original {
		t.Errorf("WithTimeLayout empty: default should be preserved, got %q", o.timeLayout)
	}
}

func TestOptions_MultipleOptions(t *testing.T) {
	o := defaultOptions()
	for _, opt := range []Option{
		WithTimeField("@timestamp"),
		WithTimeLayout("2006-01-02"),
	} {
		opt(&o)
	}
	if o.timeField != "@timestamp" {
		t.Errorf("timeField: want %q, got %q", "@timestamp", o.timeField)
	}
	if o.timeLayout != "2006-01-02" {
		t.Errorf("timeLayout: want %q, got %q", "2006-01-02", o.timeLayout)
	}
}
