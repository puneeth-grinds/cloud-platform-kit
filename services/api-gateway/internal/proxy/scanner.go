package proxy
import (
	"net/http"
)
type ScannerProxy struct {
	baseURL string `json:"baseurl"`
	httpClient *http.Client  
}
