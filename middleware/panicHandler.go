package logiamdw

import (
	"encoding/json"
	"fmt"
	logiapkg "github.com/yusologia/go-core/v2/pkg"
	"github.com/yusologia/go-core/v2/response"
	"log"
	"net/http"
	"os"
)

func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.Header().Set("Content-Type", "application/json")

				var res *logiares.ResponseError
				if panicData, ok := r.(*logiares.ResponseError); ok {
					res = panicData
				} else {
					res = &logiares.ResponseError{
						Status: logiares.StatusError{
							Bug: true,
							Status: logiares.Status{
								Code:    http.StatusInternalServerError,
								Message: "An error Occurred.",
							},
						},
					}
				}

				fmt.Fprintf(os.Stderr, "panic: %v\n", r)
				logiapkg.LogError(r)

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
