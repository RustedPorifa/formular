package variantconfig

type VariantInfo struct {
	UUID          string `json:"uuid"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Class         string `json:"class"`
	Subject       string `json:"subject"`
	Solved        bool   `json:"solved"`
	PDFFilePath   string `json:"pdfFilePath"`
	VideoFilePath string `json:"videoFilePath"`
}
