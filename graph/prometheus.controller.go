package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

var (
	addressPrometheus   = "http://localhost:9090"
	logEnabled = false
	inited = false
)

type Config struct {
	Spec struct {
		PrometheusAddress string `yaml:"prometheusAddress"`
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
		addressPrometheus = c.Spec.PrometheusAddress
	}
	logEnabled = c.Spec.EnableLog
}


func Init(ctx context.Context) (context.Context, context.CancelFunc, v1.API) {
	if !inited {
		loadConfig()
	}
	inited = true

	client, err := api.NewClient(api.Config{
		Address: addressPrometheus,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx0, cancel := context.WithTimeout(ctx, 10*time.Second)
	return ctx0, cancel, v1api
}

func queryRange(ctx context.Context, query string, stepInMin int, historyInMin int) model.Matrix {
	ctx0, cancel, v1api := Init(ctx)
	defer cancel()

	r := v1.Range{
		Start: time.Now().Add(-time.Minute * time.Duration(historyInMin)),
		End:   time.Now(),
		Step:  time.Minute * time.Duration(stepInMin),
	}

	result, warnings, err := v1api.QueryRange(ctx0, query, r)
	logError(warnings, err)
	logResult(result)

	return result.(model.Matrix)
}

func querySimple(ctx context.Context, query string) model.Vector {
	ctx0, cancel, v1api := Init(ctx)
	defer cancel()


	result, warnings, err := v1api.Query(ctx0, query, time.Now())
	logError(warnings, err)
	logResult(result)

	return result.(model.Vector)
}

func labelValues(ctx context.Context, label string) model.LabelValues {
	ctx0, cancel, v1api := Init(ctx)
	defer cancel()

	result, warnings, err := v1api.LabelValues(ctx0, label, []string{}, time.Now().Add(-time.Hour), time.Now())
	logError(warnings, err)
	logResult(result)
	return result
}

func labels(ctx context.Context) []string {
	ctx0, cancel, v1api := Init(ctx)
	defer cancel()


	result, warnings, err := v1api.LabelNames(ctx0, []string{}, time.Now().Add(-time.Hour), time.Now())
	logError(warnings, err)
	logResult(result)
	return result
}

func series(ctx context.Context, match []*string) []model.LabelSet {
	ctx0, cancel, v1api := Init(ctx)
	defer cancel()

	result, warnings, err := v1api.Series(ctx0, CreateStringArrayFromPointers(match), time.Now().Add(-time.Hour), time.Now())
	logError(warnings, err)
	logResult(result)
	return result
}

func targets(ctx context.Context,) v1.TargetsResult {
	ctx0, cancel, v1api := Init(ctx)
	defer cancel()

	result, _ := v1api.Targets(ctx0)
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

func createMapFromMetric(obj model.Metric) (map[string]interface{}, error){
	result := map[string]interface{}{}
	for k, v := range obj {
		result[string(k)] = v
	}
	return result, nil
}

func createMapFromLabelSet(obj []model.LabelSet) ([]map[string]interface{}, error){
	result := []map[string]interface{}{}
	for _, labelSet := range obj {
		m := map[string]interface{}{}
		for k, v := range labelSet {
			m[string(k)] = v
		}
		result = append(result, m)
	}
	return result, nil
}

func CreateStringArrayFromPointers(strPtr []*string) []string {
	rtn := []string{}
	for _, v := range strPtr {
		rtn = append(rtn, *v)
	}
	return rtn
}