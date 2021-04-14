package fclip

// PasteFromClipboardInfo struct
type PasteFromClipboardInfo struct {
	Name   string `json:"name"`
	Mime   string `json:"mime"`
	Base64 string `json:"base64"`
}

// Clipboard interface
type Clipboard interface {
	CopyToClipboard(filePaths []string) (err error)
	PasteFromClipboard() (files []PasteFromClipboardInfo, err error)
}

// ClipboardObj struct
type ClipboardObj struct {

}

