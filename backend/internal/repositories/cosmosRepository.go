package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"kubastach.pl/backend/internal/config"
	"kubastach.pl/backend/internal/models"
)

// CosmosRepository provides data access to Azure Cosmos DB.
type CosmosRepository struct {
	client   *azcosmos.Client
	dbName   string
	collName string
}

// NewRepository creates a new cosmos repository
func NewCosmosRepository(ctx context.Context, cfg *config.CosmosConfig) (*CosmosRepository, error) {
	cred, err := azcosmos.NewKeyCredential(cfg.Key)
	if err != nil {
		return nil, fmt.Errorf("creating credential: %w", err)
	}
	client, err := azcosmos.NewClientWithKey(cfg.Endpoint, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("creating cosmos client: %w", err)
	}
	return &CosmosRepository{client: client, dbName: cfg.DatabaseName, collName: cfg.ContainerName}, nil
}

func (r *CosmosRepository) GetJunkFoodInDateRange(ctx context.Context, minDate, maxDate time.Time) ([]models.JunkFood, error) {
	q := "SELECT c.date, c.type, c.elapsedTime FROM c WHERE c.category = 'junkFood' AND c.date >= @minDate AND c.date <= @maxDate"
	params := []azcosmos.QueryParameter{
		{Name: "@minDate", Value: minDate.Format("2006-01-02")},
		{Name: "@maxDate", Value: maxDate.Format("2006-01-02")},
	}

	raws, err := r.queryRaw(ctx, q, params)
	if err != nil {
		return nil, err
	}

	return decodeAll[models.JunkFood](raws)
}

func (r *CosmosRepository) GetActivitiesInDateRange(ctx context.Context, minDate, maxDate time.Time) ([]models.Activity, error) {
	q := "SELECT c.date, c.type, c.elapsedTime FROM c WHERE c.category = 'activity' AND c.date >= @minDate AND c.date <= @maxDate"
	params := []azcosmos.QueryParameter{
		{Name: "@minDate", Value: minDate.Format("2006-01-02")},
		{Name: "@maxDate", Value: maxDate.Format("2006-01-02")},
	}

	raws, err := r.queryRaw(ctx, q, params)
	if err != nil {
		return nil, err
	}

	return decodeAll[models.Activity](raws)
}

// ListPosts returns all posts with category post.
func (r *CosmosRepository) ListPosts(ctx context.Context) ([]models.Post, error) {
	q := "SELECT c.content, c.date, c.media FROM c WHERE c.category = 'post'"

	raws, err := r.queryRaw(ctx, q, nil)
	if err != nil {
		return nil, err
	}

	return decodeAll[models.Post](raws)
}

// GetPostByDate returns a post for a given date (YYYY-MM-DD) by matching date portion of date.
func (r *CosmosRepository) GetPostByDate(ctx context.Context, date string) (*models.Post, error) {
	q := "SELECT c.content, c.date, c.media FROM c WHERE c.category = 'post' AND c.date = @date"
	params := []azcosmos.QueryParameter{{Name: "@date", Value: date}}

	raw, err := r.firstRaw(ctx, q, params)
	if err != nil {
		return nil, err
	}

	return decodeOne[models.Post](raw)
}

// ListDates returns list of distinct dates for posts.
func (r *CosmosRepository) ListDates(ctx context.Context) ([]string, error) {
	type item struct {
		Date string `json:"date"`
	}

	q := "SELECT c.date FROM c WHERE c.category = 'post'"
	raws, err := r.queryRaw(ctx, q, nil)
	if err != nil {
		return nil, err
	}

	items, err := decodeAll[item](raws)
	if err != nil {
		return nil, err
	}

	dates := make([]string, 0, len(items))
	for _, it := range items {
		dates = append(dates, it.Date)
	}

	slices.Sort(dates)
	slices.Reverse(dates)
	return dates, nil
}

func (r *CosmosRepository) container() (*azcosmos.ContainerClient, error) {
	return r.client.NewContainer(r.dbName, r.collName)
}

func (r *CosmosRepository) getPager(ctx context.Context, query string, params []azcosmos.QueryParameter) (*runtime.Pager[azcosmos.QueryItemsResponse], error) {
	var pkr azcosmos.PartitionKey
	container, err := r.container()
	if err != nil {
		return nil, err
	}

	enableCrossPartitionQuery := true
	opts := &azcosmos.QueryOptions{QueryParameters: params, EnableCrossPartitionQuery: &enableCrossPartitionQuery}
	pager := container.NewQueryItemsPager(query, pkr, opts)
	return pager, nil
}

// queryRaw returns raw JSON documents (caller decides target type).
func (r *CosmosRepository) queryRaw(ctx context.Context, query string, params []azcosmos.QueryParameter) ([]json.RawMessage, error) {
	pager, err := r.getPager(ctx, query, params)
	if err != nil {
		return nil, err
	}
	var items []json.RawMessage
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, b := range page.Items {
			items = append(items, json.RawMessage(b))
		}
	}
	return items, nil
}

// firstRaw returns the first raw JSON document or nil.
func (r *CosmosRepository) firstRaw(ctx context.Context, query string, params []azcosmos.QueryParameter) (json.RawMessage, error) {
	pager, err := r.getPager(ctx, query, params)
	if err != nil {
		return nil, err
	}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, b := range page.Items {
			return json.RawMessage(b), nil
		}
	}
	return nil, nil
}

// Generic decode helpers (optional usage).
func decodeAll[T any](raws []json.RawMessage) ([]T, error) {
	out := make([]T, 0, len(raws))
	for _, r := range raws {
		var v T
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, fmt.Errorf("decode item: %w", err)
		}
		out = append(out, v)
	}
	return out, nil
}

func decodeOne[T any](raw json.RawMessage) (*T, error) {
	if raw == nil {
		return nil, nil
	}
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("decode item: %w", err)
	}
	return &v, nil
}
