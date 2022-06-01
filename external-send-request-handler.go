package nateer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type ExternalSendRequestHandler struct {
	sampleDSStore *SampleDSStore
}

func NewExternalSendRequestHandler(ctx context.Context, sampleDSStore *SampleDSStore) (*ExternalSendRequestHandler, error) {
	return &ExternalSendRequestHandler{
		sampleDSStore: sampleDSStore,
	}, nil
}

func (h *ExternalSendRequestHandler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	_, err = h.sampleDSStore.Create(ctx, &SampleDSEntity{
		ID:        uuid.New().String(),
		Note:      fmt.Sprintf("%s:%s", res.Status, u.String()),
		CreatedAt: time.Now(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
}
