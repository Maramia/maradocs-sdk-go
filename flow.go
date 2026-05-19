package maradocs

import (
	"context"
	"io"
)

// Flow composes multi-step document pipelines (mirrors TS Flow).
type Flow struct {
	data *DataEp
	img  *ImgEp
	pdf  *PdfEp
}

func newFlow(data *DataEp, img *ImgEp, pdf *PdfEp) *Flow {
	return &Flow{data: data, img: img, pdf: pdf}
}

func requestOptsFromMs(timeout *int) *RequestOptions {
	if timeout == nil {
		return nil
	}
	ms := *timeout
	return &RequestOptions{Timeout: &ms}
}

// OcrImgOptions configures the OCR image flow.
type OcrImgOptions struct {
	ExtractDocument *bool
	OnProgress      func(float64)
	PDFOptions      *ImgToPdfOptions
	Timeout         *int
}

// OcrPdfOptions configures the OCR PDF flow.
type OcrPdfOptions struct {
	OnProgress func(float64)
	Password   *string
	Timeout    *int
}

func extractDocuments(opt OcrImgOptions) bool {
	if opt.ExtractDocument == nil {
		return true
	}
	return *opt.ExtractDocument
}

// OcrImg uploads an image and returns a searchable optimized PDF handle.
func (f *Flow) OcrImg(ctx context.Context, fileName string, size int64, r io.Reader, opt OcrImgOptions) (PdfHandle, error) {
	uploaded, err := f.data.Upload(ctx, fileName, size, r, opt.OnProgress)
	if err != nil {
		return PdfHandle{}, err
	}
	valid, err := f.img.Validate(ctx, ImgValidateRequest{UnvalidatedFileHandle: uploaded.UnvalidatedFileHandle}, requestOptsFromMs(opt.Timeout))
	if err != nil {
		return PdfHandle{}, err
	}
	imgHandle, err := OkImg(*valid)
	if err != nil {
		return PdfHandle{}, err
	}
	return f.OcrImgHandle(ctx, imgHandle, opt)
}

// OcrImgHandle runs the OCR pipeline starting from a validated image handle.
func (f *Flow) OcrImgHandle(ctx context.Context, imgHandle ImgHandle, opt OcrImgOptions) (PdfHandle, error) {
	ro := requestOptsFromMs(opt.Timeout)
	var imgHandles []ImgHandle
	if extractDocuments(opt) {
		docs, err := f.img.FindDocuments(ctx, ImgFindDocumentsRequest{ImgHandle: imgHandle}, ro)
		if err != nil {
			return PdfHandle{}, err
		}
		for _, doc := range docs.Documents {
			ex, err := f.img.ExtractQuadrilateral(ctx, ImgExtractQuadrilateralRequest{
				ImgHandle:     imgHandle,
				Quadrilateral: doc.Quadrilateral,
			}, ro)
			if err != nil {
				return PdfHandle{}, err
			}
			imgHandles = append(imgHandles, ex.ImgHandle)
		}
	}
	if len(imgHandles) > 0 {
		for i := range imgHandles {
			oriented, err := f.img.Orientation(ctx, ImgOrientationRequest{ImgHandle: imgHandles[i]}, ro)
			if err != nil {
				return PdfHandle{}, err
			}
			imgHandles[i] = oriented.RotatedImgHandle
			enhanced, err := f.img.ImproveContrast(ctx, ImgImproveContrastRequest{ImgHandle: imgHandles[i]}, ro)
			if err != nil {
				return PdfHandle{}, err
			}
			imgHandles[i] = enhanced.ImgHandle
		}
	} else {
		imgHandles = []ImgHandle{imgHandle}
	}
	var pdfHandles []PdfHandle
	for _, handle := range imgHandles {
		ocrp, err := f.img.OcrToPdf(ctx, ImgOcrToPdfRequest{ImgHandle: handle, Options: opt.PDFOptions}, ro)
		if err != nil {
			return PdfHandle{}, err
		}
		pdfHandles = append(pdfHandles, ocrp.PdfHandle)
	}
	pdfs := make([]PdfComposePdf, len(pdfHandles))
	for i, ph := range pdfHandles {
		pdfs[i] = PdfComposePdf{PdfHandle: ph}
	}
	combined, err := f.pdf.Compose(ctx, PdfComposeRequest{Pdfs: pdfs}, ro)
	if err != nil {
		return PdfHandle{}, err
	}
	ocr, err := f.pdf.OcrToPdf(ctx, PdfOcrToPdfRequest{PdfHandle: combined.PdfHandle}, ro)
	if err != nil {
		return PdfHandle{}, err
	}
	optimized, err := f.pdf.Optimize(ctx, PdfOptimizeRequest{PdfHandle: ocr.PdfHandle}, ro)
	if err != nil {
		return PdfHandle{}, err
	}
	return optimized.PdfHandle, nil
}

// OcrPdf uploads a PDF and returns an OCR'd optimized PDF handle.
func (f *Flow) OcrPdf(ctx context.Context, fileName string, size int64, r io.Reader, opt OcrPdfOptions) (PdfHandle, error) {
	uploaded, err := f.data.Upload(ctx, fileName, size, r, opt.OnProgress)
	if err != nil {
		return PdfHandle{}, err
	}
	valid, err := f.pdf.Validate(ctx, PdfValidateRequest{
		UnvalidatedFileHandle: uploaded.UnvalidatedFileHandle,
		Password:              opt.Password,
	}, requestOptsFromMs(opt.Timeout))
	if err != nil {
		return PdfHandle{}, err
	}
	pdfHandle, err := OkPdf(*valid)
	if err != nil {
		return PdfHandle{}, err
	}
	return f.OcrPdfHandle(ctx, pdfHandle, opt)
}

// OcrPdfHandle runs orientation + OCR + optimize on an existing PDF handle.
func (f *Flow) OcrPdfHandle(ctx context.Context, pdfHandle PdfHandle, opt OcrPdfOptions) (PdfHandle, error) {
	ro := requestOptsFromMs(opt.Timeout)
	oriented, err := f.pdf.Orientation(ctx, PdfOrientationRequest{PdfHandle: pdfHandle}, ro)
	if err != nil {
		return PdfHandle{}, err
	}
	pdfHandle = oriented.RotatedPdfHandle
	ocr, err := f.pdf.OcrToPdf(ctx, PdfOcrToPdfRequest{PdfHandle: pdfHandle}, ro)
	if err != nil {
		return PdfHandle{}, err
	}
	pdfHandle = ocr.PdfHandle
	optimized, err := f.pdf.Optimize(ctx, PdfOptimizeRequest{PdfHandle: pdfHandle}, ro)
	if err != nil {
		return PdfHandle{}, err
	}
	return optimized.PdfHandle, nil
}
