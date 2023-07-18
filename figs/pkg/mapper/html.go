package mapper

type Html struct{}

func (h *Html) Map(source []byte) ([]byte, error) {
	// source + template (mustarche) = html
	return source, nil
}

func NewHtml() *Html {
	return &Html{}
}
