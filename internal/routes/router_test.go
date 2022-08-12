package routes

import (
	"net/http"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("router.go", func() {
	Describe("NewRouter()", func() {
		var method, path string
		var result bool

		JustBeforeEach(func() {
			var _rm mux.RouteMatch
			req, _ := http.NewRequest(method, path, nil)

			req.Header.Add("Content-Type", "application/json")

			cfg := config.Configuration()

			result = NewRouter(cfg).Match(req, &_rm)
		})

		Describe("POST /signup", func() {
			BeforeEach(func() {
				method = "POST"
				path = "/signup"
			})

			It("is registered", func() {
				Expect(result).To(BeTrue())
			})
		})

		Describe("POST /login", func() {
			BeforeEach(func() {
				method = "POST"
				path = "/login"
			})

			It("is registered", func() {
				Expect(result).To(BeTrue())
			})
		})

		Describe("GET /users", func() {
			BeforeEach(func() {
				method = "GET"
				path = "/users"
			})

			It("is registered", func() {
				Expect(result).To(BeTrue())
			})
		})

		Describe("PUT /users", func() {
			BeforeEach(func() {
				method = "PUT"
				path = "/users"
			})

			It("is registered", func() {
				Expect(result).To(BeTrue())
			})
		})
	})
})
