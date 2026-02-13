package markdown

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/net/html"
)

func PlainTextFromMarkdown(md []byte) (string, error) {
	var buf bytes.Buffer
	mdConverter := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)

	if err := mdConverter.Convert(md, &buf); err != nil {
		return "", err
	}

	htmlContent := buf.String()
	return htmlToPlainText(htmlContent), nil
}

func htmlToPlainText(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}

	var sb strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
            text := strings.TrimSpace(n.Data)
            if text != "" {
                sb.WriteString(text)
                sb.WriteString(" ")
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)

    return strings.TrimSpace(sb.String())
}
