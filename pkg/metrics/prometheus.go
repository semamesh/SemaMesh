package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    TokenCounter = promauto.NewCounterVec(
       prometheus.CounterOpts{
          Name: "semamesh_llm_tokens_total",
          Help: "Total number of LLM tokens processed by SemaMesh",
       },
       []string{"type", "model", "namespace"},
    )

    CostCounter = promauto.NewCounterVec(
       prometheus.CounterOpts{
          Name: "semamesh_llm_cost_est_total",
          Help: "Estimated cost of LLM traffic in USD",
       },
       []string{"model", "namespace"},
    )

    RequestsTotal = promauto.NewCounterVec(
           prometheus.CounterOpts{
              Name: "semamesh_http_requests_total",
              Help: "Total number of HTTP requests proxied",
           },
           []string{"namespace", "status"},
     )
)