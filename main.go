package main

import (
    "time"
    "github.com/labstack/echo"
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"math/rand"
)

type M map[string]interface{}
var buckets = []float64{300, 1200, 5000}

//definisiin tipe tipe metrics nya
var (
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "my_counter",
			Help:      "This is my counter",
		},
		[]string{"code", "method", "path"},
		)
		

	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "my_gauge",
			Help:      "This is my gauge",
		})
	
	gaugevec = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "golang",
				Name:      "my_gauge",
				Help:      "This is my counter",
			},
			[]string{"code", "method", "path"},
			)

	histogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "golang",
			Name:      "my_histogram",
			Help:      "This is my histogram",
		})

	histogramvec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   "golang",
			Name:      "my_histogram",
			Help:        "How long it took to process the request",
			// ConstLabels: prometheus.Labels{"service": name},
			Buckets:     buckets,
		},
			[]string{"code", "method", "path"},
		)

	summary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "golang",
			Name:      "my_summary",
			Help:      "This is my summary",
		})
)


func main() {
	// register variable metrics nya
    prometheus.MustRegister(counter)
	prometheus.MustRegister(histogramvec)
	prometheus.MustRegister(gaugevec)



    r := echo.New()
    r.GET("/", func(ctx echo.Context) error {
        data := "Hello from /index"
        return ctx.String(http.StatusOK, data)
    })
	r.GET("/get_trx", get_trx)
	r.GET("/get_report", get_report)
       



	
	//open metrics prometheus path
	r.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

    r.Start(":9000")
}

//fungsi dummy, counter
func functionCounter(code string, method string, func_name string) {
	counter.WithLabelValues(code, method, func_name).Add(1)

}

//fungsi dummy, time track
func timeTrack(start time.Time, sleeptime int) (int64) {
	time.Sleep(time.Duration(sleeptime)*time.Second)
    elapsed := time.Since(start)
    return elapsed.Nanoseconds()
}


func get_trx(c echo.Context) error {

	//add counter
	functionCounter("200", "GET", "get_trx")
	

	//add histogram
	elapsed_time  := timeTrack(time.Now(),rand.Intn(5))
	histogramvec.WithLabelValues("200", "GET", "get_report").Observe(float64(elapsed_time) / 1000000) //milisecons


	//add gauge
	gaugevec.WithLabelValues("200", "GET", "get_report").Add(rand.Float64()*15 - 5) 

    return c.String(http.StatusOK, strconv.FormatInt(int64(elapsed_time) / 1000000, 10)+ "ms")
}

func get_report(c echo.Context) error {
	//add counter
	functionCounter("200", "GET", "get_report")

	//add histogram
	elapsed_time  := timeTrack(time.Now(),rand.Intn(5))
	histogramvec.WithLabelValues("200", "GET", "get_trx").Observe(float64(elapsed_time) / 1000000) //milisecons

	// add gauge
	gaugevec.WithLabelValues("200", "GET", "get_trx").Add(rand.Float64()*15 - 5) 
	
	return c.String(http.StatusOK, "nice "+ strconv.FormatInt(int64(elapsed_time) / 1000000, 10)+ "ms")
}


func getUser(c echo.Context) error {
        // User ID from path `users/:id`
        id := c.Param("id")
  return c.String(http.StatusOK, id)
}


// ref : 
// https://github.com/zbindenren/negroni-prometheus/blob/master/middleware.go
// https://pkg.go.dev/github.com/prometheus/client_golang/prometheus
// https://prometheus.io/docs/guides/go-application/