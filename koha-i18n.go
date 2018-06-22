package main

import "bytes"
import "flag"
import "log"
import "os"
import "regexp"
import "strings"

import "golang.org/x/net/html"
import "golang.org/x/net/html/atom"

var in_a_non_translatable_block bool = false
var process_i18n_needed bool = false

func main() {
    in_place := false
    flag.BoolVar(&in_place, "in-place", false, "Modify files passed in arguments")
    flag.Parse()

    for _, filename := range flag.Args() {
        in_a_non_translatable_block = false
        process_i18n_needed = false

        log.Println(filename)
        if strings.Contains(filename, "/tables/") {
            log.Println("Skip because it's not HTML");
            continue
        }

        file, _  := os.Open(filename)
        z := html.NewTokenizer(file)
        var output bytes.Buffer

        for {
            if z.Next() == html.ErrorToken {
                break
            }

            token := z.Token()
            if (token.Type == html.TextToken) {
                writePossiblyTranslatableTextToBuffer(&output, token.Data)
            } else {
                if (token.Type == html.StartTagToken && (token.DataAtom == atom.Script || token.DataAtom == atom.Style)) {
                    in_a_non_translatable_block = true
                } else if (token.Type == html.EndTagToken && (token.DataAtom == atom.Script || token.DataAtom == atom.Style)) {
                    in_a_non_translatable_block = false
                }
                if (token.Type == html.StartTagToken || token.Type == html.SelfClosingTagToken) {
                    output.WriteString("<")
                    output.WriteString(token.Data)
                    for _, attr := range token.Attr {
                        output.WriteString(" ")
                        output.WriteString(attr.Key)
                        output.WriteString("=\"")
                        if (0 == strings.Compare(attr.Key, "alt") ||
                          0 == strings.Compare(attr.Key, "title") ||
                          0 == strings.Compare(attr.Key, "label") ||
                          0 == strings.Compare(attr.Key, "placeholder")) {
                            writePossiblyTranslatableTextToBuffer(&output, attr.Val)
                        } else {
                            output.WriteString(strings.Replace(attr.Val, "&#39;", "'", -1))
                        }
                        output.WriteString("\"")
                    }
                    if (token.Type == html.SelfClosingTagToken) {
                        output.WriteString(" /")
                    }
                    output.WriteString(">")
                } else {
                    output.WriteString(token.String())
                }
            }
        }
        file.Close()

        if in_place {
            file, _ = os.OpenFile(filename, os.O_WRONLY | os.O_TRUNC, 0644)
            if process_i18n_needed {
                file.WriteString("[% PROCESS 'i18n.inc' %]\n")
            }
            file.Write(output.Bytes())
            file.Close()
        } else {
            os.Stdout.Write(output.Bytes())
        }
    }
}

func writePossiblyTranslatableTextToBuffer(b *bytes.Buffer, text string) {
    // This ugly regexp is used to split HTML text nodes into
    // multiple chunks: TT chunks, "words" chunks (translatable),
    // white space chunks and "everything else" chunks
    r := "(\\[%(?s:.)*?%\\]|(?:[^\\s[]|\\[[^\\s%])+(?: (?:[^\\s[]|\\[[^\\s%])+)*|\\s+|(?s:.)+)"
    dataIndices := regexp.MustCompile(r).FindAllStringIndex(text, -1)
    for i := 0; i < len(dataIndices); i++ {
        tmp := text[dataIndices[i][0]:dataIndices[i][1]]

        // Never translate a TT directive, a string that contains
        // only non-word characters, or the content of script and
        // style elements
        if (regexp.MustCompile("(\\[%|^\\W+$)").MatchString(tmp) || in_a_non_translatable_block) {
            b.WriteString(tmp)
        } else {
            b.WriteString("[% t('")
            text := html.EscapeString(tmp)
            text = strings.Replace(text, "&#39;", "\\'", -1)
            b.WriteString(text)
            b.WriteString("') %]")
            process_i18n_needed = true
        }
    }
}
