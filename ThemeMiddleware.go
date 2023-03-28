package dashboard

import (
	"context"
	"net/http"
)

func ThemeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		themeName := ThemeNameRetrieveFromCookie(r)

		themeName = themeNameVerifyAndFix(themeName)

		ctx := context.WithValue(r.Context(), ThemeNameContextKey{}, themeName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
