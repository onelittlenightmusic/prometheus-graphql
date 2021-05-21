package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

var (
	address   = "http://localhost:9090"
	port      = 2112
	logEnabled = false
)

type Config struct {
	Spec struct {
		PrometheusAddress string `yaml:"prometheusAddress"`
		GraphqlPort       int    `yaml:"graphqlPort"`
		EnableLog         bool   `yaml:"logEnabled"`
	}
}

func loadConfig() {
	c := Config{}
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if c.Spec.PrometheusAddress != "" {
		address = c.Spec.PrometheusAddress
	}
	if c.Spec.GraphqlPort != 0 {
		port = c.Spec.GraphqlPort
	}
	logEnabled = c.Spec.EnableLog
}

var prometheusModelType = graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
		Name:        "SampleStream",
		Description: "",
		Fields: graphql.Fields{
			"metric": &graphql.Field{
				Type:        graphql.String,
				Description: "Prometheus metrics",
				Resolve: func(rp graphql.ResolveParams) (interface{}, error) {
					if sampleStream, ok := rp.Source.(model.SampleStream); ok {
						jsonResult, _ := json.Marshal(sampleStream.Metric)
						return string(jsonResult), nil
					}
					return nil, nil
				},
			},
			"values": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
					Name:        "SamplePair",
					Description: "Prometheus sample pair",
					Fields: graphql.Fields{
						"value": &graphql.Field{
							Type:        graphql.String,
							Description: "Prometheus metrics",
							Resolve: func(rp graphql.ResolveParams) (interface{}, error) {
								if samplePair, ok := rp.Source.(model.SamplePair); ok {
									return samplePair.Value, nil
								}
								return nil, nil
							},
						},
						"timestamp": &graphql.Field{
							Type:        graphql.String,
							Description: "Prometheus timestamp",
							Resolve: func(rp graphql.ResolveParams) (interface{}, error) {
								if samplePair, ok := rp.Source.(model.SamplePair); ok {
									return samplePair.Timestamp, nil
								}
								return nil, nil
							},
						}},
				})),
				Description: "Prometheus values",
				Resolve: func(rp graphql.ResolveParams) (interface{}, error) {
					if sampleStream, ok := rp.Source.(model.SampleStream); ok {
						return sampleStream.Values, nil
					}
					return nil, nil
				},
			},
		},
	}))

func main() {
	loadConfig()

	fields := graphql.Fields{
		"query_range": &graphql.Field{
			Type: prometheusModelType,
			Resolve: func(rp graphql.ResolveParams) (interface{}, error) {
				result := queryRange(rp.Args["query"].(string), 1, 100)
				var resultsArray []model.SampleStream
				for _, v := range result {
					resultsArray = append(resultsArray, *v)
				}

				return resultsArray, nil
			},
			Args: graphql.FieldConfigArgument{
				"query": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
		},
		"query": &graphql.Field{
			Type: prometheusModelType,
			Resolve: func(rp graphql.ResolveParams) (interface{}, error) {
				result := query(rp.Args["query"].(string))
				if result == nil {
					return nil, nil
				}
				var resultsArray []model.SampleStream
				for _, v := range result {
					resultsArray = append(resultsArray, *v)
				}

				return resultsArray, nil
			},
			Args: graphql.FieldConfigArgument{
				"query": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
		},
	}
	// var names []string
	labels := labelValues()
	for _, v := range labels {
		fmt.Printf("label: %s", v)
		label := string(v)
		fields[label] = createQueryField(label)
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func createQueryField(label string) *graphql.Field {
	return &graphql.Field{
		Type: prometheusModelType,
		Resolve: func(rp graphql.ResolveParams) (interface{}, error) {
			stepInMin := rp.Args["step_in_min"].(int)
			historyInMin := rp.Args["history_in_min"].(int)
			fmt.Printf("Parameters history_in_min %d, step_in_min: %d", stepInMin, historyInMin)
			result := queryRange(fmt.Sprintf("avg_over_time(%s[%dm])", label, stepInMin), stepInMin, historyInMin)
			var resultsArray []model.SampleStream
			for _, v := range result {
				resultsArray = append(resultsArray, *v)
			}

			return resultsArray, nil
		},
		Args: graphql.FieldConfigArgument{
			"step_in_min": &graphql.ArgumentConfig{
				Description:  "Step interval length in minutes",
				Type:         graphql.Int,
				DefaultValue: 1,
			},
			"history_in_min": &graphql.ArgumentConfig{
				Description:  "Step interval length in minutes",
				Type:         graphql.Int,
				DefaultValue: 100,
			},
		},
	}
}

func queryRange(query string, stepInMin int, historyInMin int) model.Matrix {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r := v1.Range{
		Start: time.Now().Add(-time.Minute * time.Duration(historyInMin)),
		End:   time.Now(),
		Step:  time.Minute * time.Duration(stepInMin),
	}

	result, warnings, err := v1api.QueryRange(ctx, query, r)
	logError(warnings, err)
	logResult(result)

	return result.(model.Matrix)
}

func query(query string) model.Matrix {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := v1api.Query(ctx, query, time.Now())
	logError(warnings, err)
	logResult(result)

	return result.(model.Matrix)
}

func labelValues() model.LabelValues {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := v1api.LabelValues(ctx, "__name__", time.Now().Add(-time.Hour), time.Now())
	logError(warnings, err)
	logResult(result)
	return result
}

func logError(warnings v1.Warnings, err error) {
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
}

func logResult(result interface{}) {
	if logEnabled {
		fmt.Printf("Result:\n%v\n", result)

		jsonResult, err := json.Marshal(result)
		if err != nil {
			fmt.Printf("Error json marshal")
		}
		fmt.Printf("Json: %s", string(jsonResult))
	}
}