package logins_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/adamstrickland/dapper-api/internal/logins"
	"github.com/adamstrickland/dapper-api/internal/security"
	"github.com/adamstrickland/dapper-api/internal/users"
	"github.com/bxcodec/faker/v3"
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("logins/handlers.go", func() {
	var (
		rr           *httptest.ResponseRecorder
		body, secret string
		cfg          *config.Config
	)

	Describe("NewPostHandler()", func() {
		BeforeEach(func() {
			secret = "whatever"

			cfg = config.Configuration()
			cfg.Set("secret", secret)
		})

		JustBeforeEach(func() {
			req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(body))

			Expect(err).NotTo(HaveOccurred())

			rr = httptest.NewRecorder()

			handler := http.HandlerFunc(logins.NewPostHandler(cfg))

			handler.ServeHTTP(rr, req)
		})

		When("an invalid payload is provided", func() {
			It("is not OK", func() {
				Expect(rr.Code).NotTo(Equal(http.StatusOK))
			})
		})

		When("a valid payload is provided", func() {
			var (
				password, email string
			)

			BeforeEach(func() {
				email = faker.Email()
				password = "p@ssw0rd"

				body = fmt.Sprintf(`
				{
				  "email": "%s",
				  "password": "%s"
				}`, email, password)
			})

			When("but the provided email is not found", func() {
				BeforeEach(func() {
					u, _ := users.FindByEmail(cfg, email)
					Expect(u).To(BeNil())
				})

				It("is not OK", func() {
					Expect(rr.Code).NotTo(Equal(http.StatusOK))
				})
			})

			When("the provided email is found", func() {
				When("but the password does not match", func() {
					BeforeEach(func() {
						_, err := users.Create(cfg, &users.User{
							Email:               email,
							UnencryptedPassword: "somethingthatdoesnotmatch",
						})
						Expect(err).NotTo(HaveOccurred())
					})

					It("is not OK", func() {
						Expect(rr.Code).NotTo(Equal(http.StatusOK))
					})
				})

				When("and the password matches", func() {
					var (
						p security.TokenPayload
					)

					BeforeEach(func() {
						_, err := users.Create(cfg, &users.User{
							Email:               email,
							UnencryptedPassword: password,
						})
						Expect(err).NotTo(HaveOccurred())
					})

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
			})
		})
	})
})
