package gintonic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/unifyi/creme-brulee/logging"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func MakeRequest(ctx context.Context, method string, url string, obj interface{}) ([]byte, error) {
	ctx, span := logging.StartSpan(ctx, "MakeRequest")
	defer span.End()
	span.SetAttributes(
		attribute.String("method", method),
		attribute.String("URL", url),
	)
	log := ctxlogrus.Extract(ctx)

	jsonBody, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	log.Debugf("%v sending payload to %v: %v", method, url, string(jsonBody))
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	span.AddEvent("beforeRemoteCall")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("response payload %v", string(body))

	span.AddEvent("afterRemoteCall", trace.WithAttributes(attribute.Int("httpStatus", res.StatusCode)))
	if isHttp, status := isHttpError(res); isHttp {
		return body, fmt.Errorf("error %v", status)
	}
	return body, nil
}

func isHttpError(res *http.Response) (bool, int) {
	return res.StatusCode >= 400, res.StatusCode
}
