package docx

import (
	"bytes"
	"fmt"
	"generator/internal/entity"
	"strings"

	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/schema/soo/wml"
)

type Generator struct{}

const (
	WaterMark = "{table}"
)

func NewGenerator(cfg Config) (*Generator, error) {
	if err := license.SetMeteredKey(cfg.Key); err != nil {
		  return nil, err
	}

	return &Generator{}, nil
}

func (g *Generator) Generate(data []byte, table entity.Table) ([]byte, error){
	doc, err := document.Read(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to read document: %w", err)
	}

	i := -1
	builder := strings.Builder{}
	for j, p := range doc.Paragraphs(){
		for _, r := range p.Runs(){
			builder.WriteString(r.Text())
		}
		if strings.Contains(builder.String(), WaterMark){
			i = j
			break 
		}
		builder.Reset()
	}

	if i == -1 {
		return nil, fmt.Errorf("watermark not found")
	}


	temp := doc.Paragraphs()[i].Runs()
	for _, r := range temp {
		doc.Paragraphs()[i].RemoveRun(r)
	}

	p := doc.Paragraphs()[i]
	tbl := doc.InsertTableAfter(p)

	borders := tbl.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, 2)

	clmn := tbl.AddRow()
	for _, col := range table.Columns {
		clmn.AddCell().AddParagraph().AddRun().AddText(col)
	}

	for _, row := range table.Rows {
		r := tbl.AddRow()
		for _, c := range row {
			r.AddCell().AddParagraph().AddRun().AddText(c)
		}
   
	}

	buf := new(bytes.Buffer)
	if err := doc.Save(buf); err != nil {
  		return nil, fmt.Errorf("failed to save document: %w", err)
 	}

	return buf.Bytes(), nil
}