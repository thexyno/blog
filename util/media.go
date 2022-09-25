package util

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	log "github.com/sirupsen/logrus"
)

var addedNow = make(map[string]string)

func convertImage(link string, mediaDir, postPath string) (string, error) {
	var buffer bytes.Buffer
	file, err := os.Open(postPath + "/" + link)
	if err != nil {
		return "", err
	}
	defer file.Close()
	imData, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 80)
	if err := webp.Encode(&buffer, imData, options); err != nil {
		log.Fatalln(err)
	}
	sha := sha256.Sum256(buffer.Bytes())
	name := strings.Split(path.Base(link), ".")
	filename := fmt.Sprintf("%x-%s.webp", sha[0:15], name[len(name)-2])
	writefile, err := os.Create(mediaDir + "/" + filename)
	if err != nil {
		return "", err
	}
	defer writefile.Close()
	_, err = buffer.WriteTo(writefile)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func dontConvertImage(link, mediaDir, postPath string) (string, error) {
	file, err := os.Open(postPath + "/" + link)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New() // Use the Hash in crypto/sha256
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	sum := hash.Sum(nil)[0:15]
	name := strings.Split(path.Base(link), ".")
	filename := fmt.Sprintf("%x-%s.%s", sum, name[len(name)-2], name[len(name)-1])
	writefile, err := os.Create(mediaDir + "/" + filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(writefile, file); err != nil {
		return "", err
	}
	return filename, nil

}

func addImage(link string, mediaDir, postPath string) {
	suffixspl := strings.Split(link, ".")
	suffix := suffixspl[len(suffixspl)-1]
	if suffix != "svg" {
		str, err := convertImage(link, mediaDir, postPath)
		if err != nil {
			log.Panic(err)
		}
		addedNow[link] = str
	} else {
		str, err := dontConvertImage(link, mediaDir, postPath)
		if err != nil {
			log.Panic(err)
		}
		addedNow[link] = str
	}
}

func addImageHook(mediaDir, postPath string) func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		n, ok := node.(*ast.Image)

		if !ok {
			return ast.GoToNext, false
		}
		addImage(string(n.Destination), mediaDir, postPath)

		return ast.GoToNext, true
	}
}

func AddMedia(content, mediaDir, postPath string) (string, error) {
	log.Print(mediaDir)
	info, err := os.Stat(mediaDir)
	if os.IsNotExist(err) {
		os.Mkdir(mediaDir, 0755)
	} else if !info.IsDir() || err != nil {
		return "", err
	}
	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	opts := html.RendererOptions{Flags: html.CommonFlags, RenderNodeHook: addImageHook(mediaDir, postPath)}
	renderer := html.NewRenderer(opts)
	markdown.ToHTML([]byte(content), parser, renderer)
	for i, v := range addedNow {
		content = strings.ReplaceAll(content, i, v)
	}

	return content, nil
}
