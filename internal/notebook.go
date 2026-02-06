package internal

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Notebook struct {
	Device *Device
	Footer *Footer
	Header *Header
	Pages  []*Page
	Name   string
}

func NewNotebook(input string) (*Notebook, error) {
	file, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if ok, err := isNote(file); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("Unsupported file format")
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	notebook := &Notebook{Name: fileInfo.Name()}

	if err := NewFooter(file, notebook); err != nil {
		return nil, err
	}

	if err := NewHeader(file, notebook); err != nil {
		return nil, err
	}

	NewDevice(notebook)

	notebook.Pages = make([]*Page, len(notebook.Footer.PAGES))
	var wg sync.WaitGroup
	for i, addr := range notebook.Footer.PAGES {
		wg.Go(func() {
			page, err := NewPage(file, notebook, addr)
			if err != nil {
				log.Println(err)
			} else {
				notebook.Pages[i] = page
			}
		})
	}

	wg.Wait()

	return notebook, nil
}

func (notebook *Notebook) compositePage(p *Page) *image.RGBA {
	bounds := image.Rect(0, 0, notebook.Device.PageWidth, notebook.Device.PageHeight)
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, &image.Uniform{color.White}, image.Point{}, draw.Src)
	for _, l := range p.LAYERSEQ {
		draw.Draw(canvas, bounds, l.Data, image.Point{}, draw.Over)
	}
	return canvas
}

func (notebook *Notebook) ToPNG(outputPath string) error {
	opDir := filepath.Join(outputPath, notebook.Name)
	err := os.MkdirAll(opDir, 0o755)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i, p := range notebook.Pages {
		wg.Go(func() {
			op, err := os.Create(filepath.Join(opDir, "page-"+strconv.Itoa(i)+".png"))
			if err != nil {
				log.Println(err)
				return
			}
			defer op.Close()

			err = png.Encode(op, notebook.compositePage(p))
			if err != nil {
				log.Println(err)
			}
		})
	}

	wg.Wait()
	return nil
}

func (notebook *Notebook) ToPDF(outputPath string) error {
	err := os.MkdirAll(outputPath, 0o755)
	if err != nil {
		return err
	}

	width := notebook.Device.PageWidth
	height := notebook.Device.PageHeight
	totalPages := len(notebook.Pages)

	type pdfPageChunk struct {
		pageObject     []byte
		contentsObject []byte
		imageObject    []byte
	}

	// Build page chunks in parallel
	chunks := make([]pdfPageChunk, totalPages)
	var wg sync.WaitGroup
	for i, p := range notebook.Pages {
		wg.Go(func() {
			canvas := notebook.compositePage(p)

			var compressed bytes.Buffer
			zw, _ := zlib.NewWriterLevel(&compressed, zlib.BestCompression)
			pix := canvas.Pix
			rgb := [3]byte{}
			for j := range width * height {
				rgb[0] = pix[j*4]
				rgb[1] = pix[j*4+1]
				rgb[2] = pix[j*4+2]
				zw.Write(rgb[:])
			}
			zw.Close()
			compressedBytes := compressed.Bytes()

			pageObjID := (i * 3) + 3
			contentsObjID := (i * 3) + 4
			imageObjID := (i * 3) + 5

			pageObj := fmt.Sprintf(
				"%d 0 obj\n<< /Type /Page\n   /Parent 2 0 R\n   /MediaBox [0 0 595 842]\n   /Contents %d 0 R\n   /Resources << /XObject << /Im1 %d 0 R >> >>\n>>\nendobj\n",
				pageObjID, contentsObjID, imageObjID,
			)

			contents := "q\n595 0 0 842 0 0 cm\n/Im1 Do\nQ\n"
			contentsObj := fmt.Sprintf(
				"%d 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n",
				contentsObjID, len(contents), contents,
			)

			imageHeader := fmt.Sprintf(
				"%d 0 obj\n<< /Type /XObject\n   /Subtype /Image\n   /Width %d\n   /Height %d\n   /ColorSpace /DeviceRGB\n   /BitsPerComponent 8\n   /Filter /FlateDecode\n   /Length %d >>\nstream\n",
				imageObjID, width, height, len(compressedBytes),
			)

			var imageObj bytes.Buffer
			imageObj.WriteString(imageHeader)
			imageObj.Write(compressedBytes)
			imageObj.WriteString("\nendstream\nendobj\n")

			chunks[i] = pdfPageChunk{
				pageObject:     []byte(pageObj),
				contentsObject: []byte(contentsObj),
				imageObject:    imageObj.Bytes(),
			}
		})
	}
	wg.Wait()

	// Write PDF sequentially
	name := strings.TrimSuffix(notebook.Name, filepath.Ext(notebook.Name))
	f, err := os.Create(filepath.Join(outputPath, name+".pdf"))
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	var byteOffset uint64
	xrefOffsets := make([]uint64, totalPages*3+2)

	header := []byte("%PDF-1.7\n%\xe2\xe3\xcf\xd3\n")
	writer.Write(header)
	byteOffset += uint64(len(header))

	// Object 1: Catalog
	xrefOffsets[0] = byteOffset
	catalog := []byte("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	writer.Write(catalog)
	byteOffset += uint64(len(catalog))

	// Object 2: Pages root
	xrefOffsets[1] = byteOffset
	var pageRefs strings.Builder
	for i := range totalPages {
		if i > 0 {
			pageRefs.WriteString(" ")
		}
		fmt.Fprintf(&pageRefs, "%d 0 R", (i*3)+3)
	}
	pagesRoot := fmt.Sprintf("2 0 obj\n<< /Type /Pages /Kids [ %s ] /Count %d >>\nendobj\n", pageRefs.String(), totalPages)
	writer.WriteString(pagesRoot)
	byteOffset += uint64(len(pagesRoot))

	// Write page chunks
	for i, chunk := range chunks {
		idx := (i * 3) + 2

		xrefOffsets[idx] = byteOffset
		writer.Write(chunk.pageObject)
		byteOffset += uint64(len(chunk.pageObject))

		xrefOffsets[idx+1] = byteOffset
		writer.Write(chunk.contentsObject)
		byteOffset += uint64(len(chunk.contentsObject))

		xrefOffsets[idx+2] = byteOffset
		writer.Write(chunk.imageObject)
		byteOffset += uint64(len(chunk.imageObject))
	}

	// Cross-reference table
	xrefStart := byteOffset
	writer.WriteString("xref\n")
	fmt.Fprintf(writer, "0 %d\n", len(xrefOffsets)+1)
	writer.WriteString("0000000000 65535 f \n")
	for _, offset := range xrefOffsets {
		fmt.Fprintf(writer, "%010d 00000 n \n", offset)
	}

	// Trailer
	writer.WriteString("trailer\n")
	fmt.Fprintf(writer, "<< /Size %d /Root 1 0 R >>\n", len(xrefOffsets)+1)
	writer.WriteString("startxref\n")
	fmt.Fprintf(writer, "%d\n", xrefStart)
	writer.WriteString("%%EOF\n")

	writer.Flush()
	return nil
}
