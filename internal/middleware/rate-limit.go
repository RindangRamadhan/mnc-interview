package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/RindangRamadhan/mnc-interview/internal/helper"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type Request struct {
	AuthToken string                 `json:"auth_token"`
	Payload   map[string]interface{} `json:"payload"`
}

// Change the the map to hold values of the type visitor.
var mu sync.Mutex
var visitors = make(map[string]*Visitor)

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanupVisitors()
}

func (_m *maiMiddleware) Limit(ro *mux.Router) mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req = Request{
				AuthToken: r.URL.Query().Get("auth_token"),
				Payload:   map[string]interface{}{},
			}

			for k, v := range r.URL.Query() {
				if k != "auth_token" && len(v) > 0 {
					req.Payload[k] = v[0]
				}
			}

			// Conditional Request Method "POST"
			if r.Method == "POST" {
				req.Payload = map[string]interface{}{}

				err := json.NewDecoder(r.Body).Decode(&req.Payload)
				if err != nil {
					log.Println("Error decode request body :", err.Error())
					return
				}

				// Set a new body, with same payload
				jp, _ := json.Marshal(req.Payload)
				r.Body = ioutil.NopCloser(bytes.NewBuffer(jp))
			}

			// Convert req to byte
			reqByte, _ := json.Marshal(req)

			// Encode req byte
			reqEncode := base64.StdEncoding.EncodeToString(reqByte)

			// Limiter
			limiter := getVisitor(reqEncode)

			if !limiter.Allow() {
				helper.HttpHandler.ResponseJSON(w, &helper.ResponseJSON{
					HTTPCode: http.StatusTooManyRequests,
					Status:   "error",
					Message:  http.StatusText(http.StatusTooManyRequests),
				})
				return
			}

			next.ServeHTTP(w, r)
		})

	}

}

func getVisitor(payload string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[payload]
	if !exists {
		// Create new limter, only 1 request /second
		limiter := rate.NewLimiter(1, 1)

		// Include the current time when creating a new visitor.
		visitors[payload] = &Visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

// Every 30 seconds check the map for visitors that haven't been seen for
// more than 1 minute and delete the entries.
func cleanupVisitors() {
	for {
		time.Sleep(1 * time.Second)

		mu.Lock()
		for p, v := range visitors {
			if time.Since(v.lastSeen) > 2*time.Second {
				delete(visitors, p)
			}
		}
		mu.Unlock()
	}
}
