package middleware

import (
	"fmt"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// パニックが発生した場合の処理を記述
				fmt.Println("Recovered:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// 次のハンドラまたはミドルウェアを呼び出す
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
