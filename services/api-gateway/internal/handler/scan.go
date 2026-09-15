package handler
import(
	"net/http"

	"github.com/puneeth-grinds/cloud-platform-kit/services/api-gateway/internal/proxy"
)

func NewScanHandler(scannerProxy *proxy.ScannerProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scanHandler(w, r, scannerProxy)
	})
}