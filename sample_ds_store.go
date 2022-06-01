package nateer

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

// SampleDSStore
// `--vpc-egress=all-traffic` になっていても、Cloud Datastoreにアクセスできるのかを試すためのもの
type SampleDSStore struct {
	ds *datastore.Client
}

type SampleDSEntity struct {
	ID        string `datastore:"-"`
	Note      string
	CreatedAt time.Time
}

func NewSampleDSStore(ctx context.Context, ds *datastore.Client) (*SampleDSStore, error) {
	return &SampleDSStore{
		ds: ds,
	}, nil
}

func (s *SampleDSStore) Kind(ctx context.Context) string {
	return "nateer-sample"
}

func (s *SampleDSStore) Create(ctx context.Context, entity *SampleDSEntity) (*SampleDSEntity, error) {
	key := datastore.NameKey(s.Kind(ctx), entity.ID, nil)
	m := datastore.NewInsert(key, entity)
	_, err := s.ds.Mutate(ctx, m)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
