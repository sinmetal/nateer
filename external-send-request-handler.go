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
	userStore     *UserStore
}

func NewExternalSendRequestHandler(ctx context.Context, sampleDSStore *SampleDSStore, userStore *UserStore) (*ExternalSendRequestHandler, error) {
	return &ExternalSendRequestHandler{
		sampleDSStore: sampleDSStore,
		userStore:     userStore,
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

	userRow, err := h.userStore.Select(ctx, "000000a4-fb52-4c9c-8be3-6627b9c181ce")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
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
		Note:      fmt.Sprintf("%s:%s:%s", res.Status, u.String(), userRow),
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
