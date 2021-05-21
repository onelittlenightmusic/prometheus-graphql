package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/hiroyukiosaki/graphql-prometheus/graph/generated"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	model1 "github.com/prometheus/common/model"
)

func (r *activeTargetResolver) DiscoveredLabels(ctx context.Context, obj *v1.ActiveTarget) (map[string]interface{}, error) {
	rtn := map[string]interface{}{}
	for k, v := range obj.DiscoveredLabels {
		rtn[k] = v
	}
	return rtn, nil
}

func (r *droppedTargetResolver) DiscoveredLabels(ctx context.Context, obj *v1.DroppedTarget) (map[string]interface{}, error) {
	rtn := map[string]interface{}{}
	for k, v := range obj.DiscoveredLabels {
		rtn[k] = v
	}
	return rtn, nil
}

func (r *queryResolver) QueryRange(ctx context.Context, query *string) ([]*model1.SampleStream, error) {
	result := queryRange(ctx, *query, 2, 30)

	return result, nil
}

func (r *queryResolver) Query(ctx context.Context, query *string) ([]*model1.Sample, error) {
	result := querySimple(ctx, *query)

	return result, nil
}

func (r *queryResolver) LabelValues(ctx context.Context, label string) ([]*string, error) {
	result := labelValues(ctx, label)

	var labelValues []*string
	for _, v := range result {
		str := string(v)
		labelValues = append(labelValues, &str)
	}
	return labelValues, nil
}

func (r *queryResolver) NameValues(ctx context.Context) ([]*string, error) {
	return r.LabelValues(ctx, "__name__")
}

func (r *queryResolver) Labels(ctx context.Context) ([]*string, error) {
	result := labels(ctx)

	var labels []*string
	for _, v := range result {
		str := string(v)
		labels = append(labels, &str)
	}
	return labels, nil
}

func (r *queryResolver) Series(ctx context.Context, match []*string) ([]map[string]interface{}, error) {
	result := series(ctx, match)
	return createMapFromLabelSet(result)
}

func (r *queryResolver) Targets(ctx context.Context) (*v1.TargetsResult, error) {
	result := targets(ctx)
	return &result, nil
}

func (r *sampleResolver) Timestamp(ctx context.Context, obj *model1.Sample) (*int, error) {
	time := int(obj.Timestamp)
	return &time, nil
}

func (r *sampleResolver) Value(ctx context.Context, obj *model1.Sample) (*float64, error) {
	return (*float64)(&obj.Value), nil
}

func (r *sampleResolver) Metric(ctx context.Context, obj *model1.Sample) (map[string]interface{}, error) {
	return createMapFromMetric(obj.Metric)
}

func (r *samplePairResolver) Timestamp(ctx context.Context, obj *model1.SamplePair) (*int, error) {
	time := int(obj.Timestamp)
	return &time, nil
}

func (r *samplePairResolver) Value(ctx context.Context, obj *model1.SamplePair) (*float64, error) {
	return (*float64)(&obj.Value), nil
}

func (r *sampleStreamResolver) Metric(ctx context.Context, obj *model1.SampleStream) (map[string]interface{}, error) {
	return createMapFromMetric(obj.Metric)
}

// ActiveTarget returns generated.ActiveTargetResolver implementation.
func (r *Resolver) ActiveTarget() generated.ActiveTargetResolver { return &activeTargetResolver{r} }

// DroppedTarget returns generated.DroppedTargetResolver implementation.
func (r *Resolver) DroppedTarget() generated.DroppedTargetResolver { return &droppedTargetResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Sample returns generated.SampleResolver implementation.
func (r *Resolver) Sample() generated.SampleResolver { return &sampleResolver{r} }

// SamplePair returns generated.SamplePairResolver implementation.
func (r *Resolver) SamplePair() generated.SamplePairResolver { return &samplePairResolver{r} }

// SampleStream returns generated.SampleStreamResolver implementation.
func (r *Resolver) SampleStream() generated.SampleStreamResolver { return &sampleStreamResolver{r} }

type activeTargetResolver struct{ *Resolver }
type droppedTargetResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type sampleResolver struct{ *Resolver }
type samplePairResolver struct{ *Resolver }
type sampleStreamResolver struct{ *Resolver }
