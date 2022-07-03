package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ericchiang/css"
	"github.com/tomnomnom/gahttp"
	"golang.org/x/net/html"
)

func extractSelector(r io.Reader, selector string) ([]string, error) {

	out := []string{}

	sel, err := css.Compile(selector)
	if err != nil {
		return out, err
	}

	node, err := html.Parse(r)
	if err != nil {
		return out, err
	}

	// it's kind of tricky to actually know what to output
	// if the resulting tags contain more than just a text node
	for _, ele := range sel.Select(node) {
		if ele.FirstChild == nil {
			continue
		}
		out = append(out, ele.FirstChild.Data)
	}

	return out, nil
}

func extractComments(r io.Reader) []string {

	z := html.NewTokenizer(r)

	out := []string{}
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		t := z.Token()

		if t.Type == html.CommentToken {
			d := strings.Replace(t.Data, "\n", " ", -1)
			d = strings.TrimSpace(d)
			if d == "" {
				continue
			}
			out = append(out, d)
		}

	}
	return out
}

func extractAttribs(r io.Reader, attribs []string) []string {
	z := html.NewTokenizer(r)

	out := []string{}

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		t := z.Token()

		for _, a := range t.Attr {

			if a.Val == "" {
				continue
			}

			for _, attrib := range attribs {
				if attrib == a.Key {
					out = append(out, a.Val)
				}
			}
		}
	}
	return out
}

func extractTags(r io.Reader, tags []string) []string {
	z := html.NewTokenizer(r)

	out := []string{}

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		t := z.Token()

		if t.Type == html.StartTagToken {

			for _, tag := range tags {
				if t.Data == tag {
					if z.Next() == html.TextToken {
						text := strings.TrimSpace(z.Token().Data)
						if text == "" {
							continue
						}
						out = append(out, text)
					}
				}
			}
		}
	}
	return out
}

type target struct {
	location string
	r        io.ReadCloser
}

func main() {
	// TODO: support quiet mode (no errors)
	// TODO: option to output file or url as context
	// TODO: add concurrency flag

	flag.Parse()

	// TODO: check mode is valid
	mode := flag.Arg(0)
	if mode == "" {
		fmt.Println("A ferramenta aceita URLs ou nomes de arquivos para documentos HTML em stdin e extrai partes deles..")
		fmt.Println("")
		fmt.Println("html-tool <comandos> [<args>]")
		fmt.Println("")
		fmt.Println("comandos:")
		fmt.Println("	tags <nome-da-tag>        Extrair texto contido em tags")
		fmt.Println("	attribs <nome-dos-atributos>  Extrair valores dos atributos")
		fmt.Println("	comments                Extrai comentários")
		fmt.Println("")
		fmt.Println("Exemplos:")
		fmt.Println("	echo site.com.br | /root/go/bin/html-tool comments ") 
			    //Filtra os comentarios da pagina/URL
		fmt.Println("	echo site.com.br | /root/go/bin/html-tool tags title ") 
			    //Filtra as tags Titulo da pagina/URL
		fmt.Println("	echo site.com.br| /root/go/bin/html-tool attribs src ")
			    //Filtra todos os atributos src da pagina/URL
		fmt.Println("	cat urls.txt | html-tool tags title a strong")
			    //Filtra as tags Titulo da pagina/URL
		fmt.Println("	find . -type f -name \"*.html\" | html-tool attribs src href")
			    //Filtra todos os atributos src e href da pagina/URL
		fmt.Println("	cat urls.txt | html-tool comments")
			     //Filtra os comentarios da pagina/URL
		return
	}

	args := flag.Args()[1:]

	targets := make(chan *target)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for t := range targets {
			vals := []string{}

			switch mode {
			case "tags":
				vals = extractTags(t.r, args)
			case "attribs":
				vals = extractAttribs(t.r, args)
			case "comments":
				vals = extractComments(t.r)
			case "query":
				var err error
				vals, err = extractSelector(t.r, flag.Arg(1))
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to parse CSS selector: %s\n", err)
					break
				}

			default:
				fmt.Fprintf(os.Stderr, "unsupported mode '%s'\n", mode)
				break
			}

			for _, v := range vals {
				fmt.Println(v)
			}

			// don't forget to close the reader when we're done with it!
			t.r.Close()
		}
		wg.Done()
	}()

	p := gahttp.NewPipeline()
	p.SetClient(gahttp.NewClient(gahttp.SkipVerify))
	p.SetConcurrency(20)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		// location can be a filename or a URL
		location := strings.TrimSpace(sc.Text())

		// if it's a URL request it with gahttp
		nl := strings.ToLower(location)
		if strings.HasPrefix(nl, "http:") || strings.HasPrefix(nl, "https:") {
			p.Get(location, func(req *http.Request, resp *http.Response, err error) {
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to fetch URL: %s\n", err)
				}
				if resp != nil && resp.Body != nil {
					targets <- &target{req.URL.String(), resp.Body}
				}
			})
			continue
		}

		// if it's a file just open it
		f, err := os.Open(location)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err)
			continue
		}

		targets <- &target{location, f}
	}
	p.Done()
	p.Wait()

	close(targets)
	wg.Wait()
}
