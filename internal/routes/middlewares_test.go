package routes

import (
	"net/http"
	"net/http/httptest"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/adamstrickland/dapper-api/internal/security"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("routes/middlewares.go", func() {
	var (
		req        *http.Request
		rr         *httptest.ResponseRecorder
		middleware func(http.Handler) http.Handler
		cfg        *config.Config
		handler    func() http.Handler
	)

	BeforeEach(func() {
		cfg = config.Configuration()

		rr = httptest.NewRecorder()

		handler = func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			})
		}

		req, _ = http.NewRequest("GET", "/itdoesntmatter", nil)
	})

	Describe("ContentTypeMiddleware()", func() {
		BeforeEach(func() {
			middleware = ContentTypeMiddleware(cfg)
		})

		JustBeforeEach(func() {
			middleware(handler()).ServeHTTP(rr, req)
		})

		When("the request does not have the right content type", func() {
			It("should be rejected", func() {
				Expect(rr.Code).NotTo(Equal(http.StatusOK))
			})
		})

		When("the request does have an auth header", func() {
			BeforeEach(func() {
				req.Header.Add("Content-Type", "application/json")
			})

			It("should be rejected", func() {
				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("AuthnMiddleware()", func() {
		BeforeEach(func() {
			middleware = AuthnMiddleware(cfg)
		})

		JustBeforeEach(func() {
			req.Header.Add("Content-Type", "application/json")

			middleware(handler()).ServeHTTP(rr, req)
		})

		When("the request does not have an auth header", func() {
			It("should be rejected", func() {
				Expect(rr.Code).NotTo(Equal(http.StatusOK))
			})
		})

		When("the request does have an auth header", func() {
			When("but the token isn't valid", func() {
				BeforeEach(func() {
					req.Header.Add(cfg.GetString("tokenHeader"), "someinvalidstring")
				})

				It("should be rejected", func() {
					Expect(rr.Code).NotTo(Equal(http.StatusOK))
				})
			})

			When("and the token is valid", func() {
				BeforeEach(func() {
					t, e := security.NewTokenForSubject(cfg, "foo@bar.com")
					Expect(e).NotTo(HaveOccurred())

					req.Header.Add(cfg.GetString("tokenHeader"), t)
				})

				It("should be accepted", func() {
					Expect(rr.Code).To(Equal(http.StatusOK))
				})
			})
		})
	})
})
