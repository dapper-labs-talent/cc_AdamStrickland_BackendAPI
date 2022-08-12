package signups_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/adamstrickland/dapper-api/internal/signups"
	"github.com/adamstrickland/dapper-api/internal/users"
	"github.com/bxcodec/faker/v3"
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("signups/handlers.go", func() {
	var rr *httptest.ResponseRecorder
	var body, secret string
	var cfg *config.Config

	Describe("NewPostHandler()", func() {
		BeforeEach(func() {
			secret = "whatever"

			cfg = config.Configuration()
			cfg.Set("secret", secret)
		})

		JustBeforeEach(func() {
			req, err := http.NewRequest("POST", "/signup", bytes.NewBufferString(body))

			Expect(err).NotTo(HaveOccurred())

			rr = httptest.NewRecorder()

			handler := http.HandlerFunc(signups.NewPostHandler(cfg))

			handler.ServeHTTP(rr, req)
		})

		When("an invalid payload is provided", func() {
			It("is not OK", func() {
				Expect(rr.Code).NotTo(Equal(http.StatusOK))
			})
		})

		When("a full payload is provided", func() {
			var email string

			BeforeEach(func() {
				email = faker.Email()

				body = fmt.Sprintf(`
				{
				  "email": "%s",
				  "password": "p@ssw0rd",
				  "firstName": "Ford",
				  "lastName": "Prefect"
				}`, email)
			})

			When("and the email is unique", func() {
				type responsePayload struct {
					Token string `json:"token"`
				}

				var p responsePayload

				It("is OK", func() {
					Expect(rr.Code).To(Equal(http.StatusOK))
				})

				It("returns a JSON reponse", func() {
					Expect(rr.Header().Get("Content-Type")).To(Equal("application/json"))
				})

				It(`returns a JWT in the "token" key`, func() {
					err := json.NewDecoder(rr.Body).Decode(&p)
					Expect(err).NotTo(HaveOccurred())

					Expect(p.Token).NotTo(BeEmpty())
				})

				It("returns a decodable JWT", func() {
					Expect(json.NewDecoder(rr.Body).Decode(&p)).NotTo(HaveOccurred())

					Expect(regexp.Match(`\w+\.\w+\.\w+`, []byte(p.Token))).To(BeTrue())

					token, err := jwt.Parse(p.Token, func(t *jwt.Token) (interface{}, error) {
						return []byte(secret), nil
					})

					Expect(err).NotTo(HaveOccurred())

					Expect(token.Valid).To(BeTrue())

					claims, ok := token.Claims.(jwt.MapClaims)

					Expect(ok).To(BeTrue())

					Expect(claims.Valid()).NotTo(HaveOccurred())
				})
			})

			When("but the email is not unique", func() {
				BeforeEach(func() {
					u, err := users.Create(cfg, &users.User{
						Email: email,
					})
					Expect(err).NotTo(HaveOccurred())
					Expect(u.ID).To(BeNumerically(">", 0))
				})

				It("is not OK", func() {
					Expect(rr.Code).NotTo(Equal(http.StatusOK))
				})
			})
		})
	})
})
