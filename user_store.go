package nateer

import (
	"context"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

type UserStore struct {
	spa *spanner.Client
}

func NewUserStore(ctx context.Context, spa *spanner.Client) (*UserStore, error) {
	return &UserStore{
		spa: spa,
	}, nil
}

func (s *UserStore) Select(ctx context.Context, id string) (string, error) {
	tx := s.spa.ReadOnlyTransaction()
	defer tx.Close()

	statement := spanner.NewStatement("SELECT ID FROM User10m WHERE ID = @ID")
	statement.Params = map[string]interface{}{"ID": id}
	iter := tx.QueryWithOptions(ctx, statement, spanner.QueryOptions{
		Mode:       nil,
		Options:    nil,
		Priority:   0,
		RequestTag: "nateer-userStore-select",
	})
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}
		return row.String(), nil
	}
	return "", nil
}
