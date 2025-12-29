package core

import (
	"bytes"
	"errors"
	"io"
	"os"
	"slices"

	"github.com/edlingao/internal/pkg/parser"
	meta "github.com/yuin/goldmark-meta"
	metaParser "github.com/yuin/goldmark/parser"
)

const (
	FailedToReadMDFile    = "failed to read markdown file"
	FailedToConvertMDFile = "failed to convert markdown to HTML"
)

type Blog struct {
	ID          string   `json:"id" db:"id"`
	Title       string   `json:"title" db:"title"`
	URL         string   `json:"url" db:"url"`
	MDURL       string   `json:"md_url" db:"md_url"`
	Description string   `json:"description" db:"description"`
	Tags        []Tag    `json:"tags" db:"tags"`
	Reactions   string   `json:"reactions" db:"reactions"`
	tags        []string `json:"-" db:"-"`
	CreatedAt   string   `json:"created_at" db:"created_at"`
}

func NewBlog(title string) *Blog {
	return &Blog{
		Title:       title,
		Description: "",
		Reactions:   "",
		URL:         "assets/processed/" + title + ".html",
		MDURL:       "assets/blogs/" + title + ".md",
	}
}

func (blog *Blog) SetTags(tags []string) {
	blog.tags = tags
}

func (blog *Blog) GetTags() []string {
	return blog.tags
}

func (blog *Blog) RemoveTag(tag string) {
	blog.tags = slices.DeleteFunc(blog.tags, func(t string) bool {
		return t == tag
	})
}

func (blog *Blog) SetDescription(description string) {
	blog.Description = description
}

func (blog *Blog) GetContent() string {
	content, err := os.ReadFile(blog.URL)
	if err != nil {
		return ""
	}
	return string(content)
}

// This step assumes that the markdown file already exists in the specified path.
// That the metadata and other content is already saved on the DB
// This will only process the markdown file and parsed it as HTML
// Also save the metadata like description and tags to the struct
func (blog *Blog) ProcessFileAndSave() error {
	data, err := blog.getMDFile()
	if err != nil {
		return err
	}

	processedHTML, metaData, err := blog.processMD(data)
	if err != nil {
		return err
	}

	description := metaData["description"]
	tags := metaData["tags"]

	if descriptionStr, ok := description.(string); ok {
		blog.SetDescription(descriptionStr)
	}

	if tagsSlice, ok := tags.([]any); ok {
		var tagsArr []string
		for _, tag := range tagsSlice {
			if tagStr, ok := tag.(string); ok {
				tagsArr = append(tagsArr, tagStr)
			}
		}
		blog.SetTags(tagsArr)
	}

	htmlBuffer := bytes.NewBuffer(processedHTML)
	err = blog.saveProcessedHTML(htmlBuffer)
	if err != nil {
		return err
	}

	return nil
}

func (blog *Blog) getMDFile() ([]byte, error) {
	mdFile := "./assets/blogs/" + blog.Title + ".md"
	data, err := os.ReadFile(mdFile)
	if err != nil {
		return []byte{}, errors.New(FailedToReadMDFile + err.Error())
	}

	return data, nil
}

func (blog *Blog) processMD(data []byte) ([]byte, map[string]any, error) {
	var buff io.Writer = &bytes.Buffer{}
	context := metaParser.NewContext()
	if err := parser.MDParser.Convert(data, buff, metaParser.WithContext(context)); err != nil {
		return []byte{}, nil, errors.New(FailedToConvertMDFile + err.Error())
	}

	metaData := meta.Get(context)
	return buff.(*bytes.Buffer).Bytes(), metaData, nil
}

func (blog *Blog) saveProcessedHTML(htmlContent *bytes.Buffer) error {
	htmlFileDIR := "./assets/processed/" + blog.Title + ".html"
	htmlFile, err := os.Create(htmlFileDIR)
	defer htmlFile.Close()
	if err != nil {
		return errors.New("failed to create HTML file: " + err.Error())
	}

	_, err = htmlFile.Write(htmlContent.Bytes())
	if err != nil {
		return errors.New("failed to write to HTML file: " + err.Error())
	}

	return nil
}
