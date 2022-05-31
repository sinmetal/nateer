package nateer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ExternalSendRequestHandler struct {
}

func NewExternalSendRequestHandler(ctx context.Context) (*ExternalSendRequestHandler, error) {
	return &ExternalSendRequestHandler{}, nil
}

func (h *ExternalSendRequestHandler) Handler(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("url")
	u, err := url.Parse(value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("%s is invalid URL. err=%s", value, err)
		return
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	req = req.WithContext(r.Context())

	c := http.DefaultClient
	res, err := c.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
}
