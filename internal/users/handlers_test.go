package users_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/adamstrickland/dapper-api/internal"
	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/adamstrickland/dapper-api/internal/security"
	"github.com/adamstrickland/dapper-api/internal/users"
	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("users/handlers.go", func() {
	var (
		rr      *httptest.ResponseRecorder
		cfg     *config.Config
		handler http.HandlerFunc
		email   string
		r       *http.Request
		err     error
		user    *users.User
		result  map[string]interface{}
	)

	BeforeEach(func() {
		cfg = config.Configuration()
		rr = httptest.NewRecorder()
		email = faker.Email()

		user = &users.User{
			Email:     email,
			FirstName: "Zaphod",
			LastName:  "Beeblebrox",
		}

		users.Create(cfg, user)
	})

	JustBeforeEach(func() {
		Expect(err).NotTo(HaveOccurred())
		handler.ServeHTTP(rr, r)
	})

	Describe("NewPutHandler()", func() {
		var (
			newfn, newln string
		)

		BeforeEach(func() {
			handler = http.HandlerFunc(users.NewPutHandler(cfg))

			newfn = "Arthur"
			newln = "Dent"

			up := &users.UserPayload{
				Email:     email,
				FirstName: newfn,
				LastName:  newln,
			}

			var data bytes.Buffer

			Expect(json.NewEncoder(&data).Encode(up)).NotTo(HaveOccurred())

			r, err = http.NewRequest("PUT", "/users", &data)
		})

		When("the request is authenticated", func() {
			var (
				token string
			)

			When("but the JWT's subject does not match the payload's email", func() {
				BeforeEach(func() {
					e := faker.Email()

					Expect(e).NotTo(Equal(email))

					token, _ = security.NewTokenForSubject(cfg, e)
					r.Header.Set(cfg.GetString("tokenHeader"), token)
				})

				It("the request is rejected", func() {
					Expect(rr.Code).NotTo(Equal(http.StatusOK))
				})

				It("the record remains unchanged", func() {
					u, _ := users.FindByEmail(cfg, email)

					Expect(u.FirstName).NotTo(Equal(newfn))
					Expect(u.LastName).NotTo(Equal(newln))
				})
			})

			When("and the JWT's subject matches the payload's email", func() {
				BeforeEach(func() {
					token, _ = security.NewTokenForSubject(cfg, email)
					r.Header.Set(cfg.GetString("tokenHeader"), token)
				})

				It("the request is accepted", func() {
					Expect(rr.Code).To(Equal(http.StatusOK))
				})

				It("modifies the record", func() {
					u, _ := users.FindByEmail(cfg, email)

					Expect(u.FirstName).To(Equal(newfn))
					Expect(u.LastName).To(Equal(newln))
				})

				It("returns the modified record", func() {
					json.Unmarshal(rr.Body.Bytes(), &result)
					Expect(result["email"]).To(Equal(email))
					Expect(result["firstName"]).To(Equal(newfn))
					Expect(result["lastName"]).To(Equal(newln))
				})
			})
		})
	})

	Describe("NewGetHandler()", func() {
		BeforeEach(func() {
			handler = http.HandlerFunc(users.NewGetHandler(cfg))
			r, err = http.NewRequest("GET", "/users", nil)
		})

		It("is OK", func() {
			Expect(rr.Code).To(Equal(http.StatusOK))
		})

		It("returns a JSON payload", func() {
			Expect(rr.Result().Header.Get("Content-Type")).To(Equal("application/json"))
		})

		It("has a top-level 'users' item", func() {
			json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(result).To(HaveKey("users"))
		})

		It("has an array of user items", func() {
			json.Unmarshal(rr.Body.Bytes(), &result)
			Expect(result["users"]).To(ContainElement(HaveKeyWithValue("email", email)))
		})

		Context("accounting for drift", func() {
			var (
				count int
			)

			BeforeEach(func() {
				db, _ := internal.NewConnection(cfg)
				db.Exec("DELETE FROM users")

				expected := 3

				for i := 0; i < expected; i++ {
					users.Create(cfg, &users.User{Email: faker.Email()})
				}

				us, _ := users.All(cfg)
				count = len(*us)

				Expect(count).To(Equal(expected))
			})

			It("has one user item for each record", func() {
				json.Unmarshal(rr.Body.Bytes(), &result)
				Expect(result["users"]).To(HaveLen(count))
			})
		})
	})
})
