package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/yusologia/go-core/helpers/logialog"
	logiares "github.com/yusologia/go-response"
	"log"
	"net/http"
	"os"
)

func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.Header().Set("Content-Type", "application/json")

				fmt.Fprintf(os.Stderr, "panic: %v\n", r)
				logialog.Error(r)

				var res *logiares.ResponseError
				if panicData, ok := r.(*logiares.ResponseError); ok {
					res = panicData
				} else {
					res = &logiares.ResponseError{
						Status: logiares.Status{
							Code:    http.StatusInternalServerError,
							Message: "An error Occurred.",
						},
					}
				}

				w.WriteHeader(res.Status.Code)

				jsonData, err := json.Marshal(res)
				if err != nil {
					log.Println("Failed to marshal error response:", err)
					return
				}

				w.Write(jsonData)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
