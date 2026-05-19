package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/diagnoxix/core/ports"
	"golang.org/x/net/html"
)

type TesseractOCR struct{}

var _ ports.OCR = (*TesseractOCR)(nil)

func (t *TesseractOCR) Parse(ctx context.Context, imgURL []byte) ([]ports.OCRWord, error) {

	return parseHOCR(""), nil
}

// --- HOCR parsing helper ---
func parseHOCR(hocr string) []ports.OCRWord {
	var words []ports.OCRWord
	doc, err := html.Parse(strings.NewReader(hocr))
	if err != nil {
		return words
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			var (
				title string
				text  string
			)
			for _, attr := range n.Attr {
				if attr.Key == "title" {
					title = attr.Val
				}
			}
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				text = strings.TrimSpace(n.FirstChild.Data)
			}
			if strings.Contains(title, "bbox") {
				var left, top, right, bottom, conf int
				fmt.Sscanf(title, "bbox %d %d %d %d; x_wconf %d", &left, &top, &right, &bottom, &conf)

				if text != "" && conf >= 80 { // filter low confidence
					words = append(words, ports.OCRWord{
						Text:       text,
						Confidence: conf,
						BBox: [4]int{
							left, top, right, bottom,
						},
					})
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return words
}

func NewTesseractOCR() *TesseractOCR {
	return &TesseractOCR{}
}
