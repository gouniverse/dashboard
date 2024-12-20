package dashboard

import (
	"strings"
	"testing"
)

func TestNotFoundController(t *testing.T) {
	dashboard := NewDashboard(Config{})
	html := dashboard.ToHTML()

	expectedExpressions := []string{
		"<!DOCTYPE html>",
		"<html>",
		"</html>",
		"<head>",
		"</head>",
		"<body>",
		"</body>",
		`<link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.css" rel="stylesheet" />`,
		`<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" />`,
		`<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>`,
		`<link href="data:image/x-icon;base64,`,
		`html, body{ height: 100%; }`,
	}

	for _, expected := range expectedExpressions {
		if !strings.Contains(html, expected) {
			t.Fatal(`Response MUST contain: `+expected, ` but found `, html)
		}
	}
}
