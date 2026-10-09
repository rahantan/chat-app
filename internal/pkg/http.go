package pkg

import "net/http"

func GetStatusCode(err error) int {
	errKind, ok := err.(*AppErr)
	if ok {
		switch errKind.Kind {
		case KindBadRequest:
			return http.StatusBadRequest
		case KindUnauthorized:
			return http.StatusUnauthorized
		case KindNotFound:
			return http.StatusNotFound
		}
	}
	return http.StatusInternalServerError
}
