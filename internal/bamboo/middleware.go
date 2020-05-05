package bamboo

import (
	"log"
	"net/http"
)

func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Printf("recovering panic: %s", err)

				nErr := newError(http.StatusInternalServerError, err.(error))
				RespondError(w, nErr)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
