package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/chasefleming/elem-go/styles"
	"github.com/rprtr258/fun"
	"github.com/samber/lo"
	"github.com/yuin/goldmark"
)

const (
	_dirPublic   = "public"
	_dirPosts    = "posts"
	_indexHeader = "Posts"
	_navIndex    = "Home"
)

var mainStyle = styles.Props{
	styles.Padding: "20px",
}

func pageIndex(postFilenames []string) string {
	return elem.Html(nil,
		elem.Head(nil,
			elem.Title(nil, elem.Text(_indexHeader)),
		),
		elem.Body(nil,
			elem.H1(nil,
				elem.Text(_indexHeader),
			),
			elem.Main(attrs.Props{attrs.Style: mainStyle.ToInline()},
				elem.Ul(nil,
					lo.Map(
						postFilenames,
						func(filename string, _ int) elem.Node {
							return elem.Li(nil,
								elem.A(attrs.Props{attrs.Href: filename},
									elem.Text(filename)))
						})...),
			),
		),
	).Render()
}

func page(title string, content elem.Node) string {
	return elem.Html(nil,
		elem.Head(nil,
			elem.Title(nil, elem.Text(title)),
		),
		elem.Body(nil,
			elem.Nav(nil, elem.A(attrs.Props{attrs.Href: "index.html"}, elem.Text(_navIndex))),
			elem.H1(nil,
				elem.Text(title),
			),
			elem.Main(attrs.Props{attrs.Style: mainStyle.ToInline()},
				content,
			),
		),
	).Render()
}

func writeHTML(filename, content string) error {
	log.Println(filename)
	return os.WriteFile(filepath.Join(_dirPublic, filename), []byte(content), 0o644)
}

func run() error {
	var filenames []string
	if err := filepath.Walk(_dirPosts, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			filenames = append(filenames, path)
		}
		return err
	}); err != nil {
		return err
	}

	postFilenames, err := fun.MapErr[string, string, error](
		func(path string) (string, error) {
			content, err := os.ReadFile(path)
			if err != nil {
				return "", err
			}

			var buf bytes.Buffer
			if err := goldmark.New().Convert(content, &buf); err != nil {
				return "", err
			}

			filename := filepath.Base(path)
			title := strings.TrimSuffix(filename, filepath.Ext(filename))
			postFilename := title + ".html"
			writeHTML(postFilename, page(title, elem.Raw(buf.String())))

			return postFilename, nil
		},
		filenames...)
	if err != nil {
		return err
	}

	return writeHTML("index.html", pageIndex(postFilenames))
}

func main() {
	log.SetFlags(0)

	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
