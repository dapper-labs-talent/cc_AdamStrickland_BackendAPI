package security

import (
	"time"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/bxcodec/faker/v3"
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("security/auth.go", func() {
	var (
		cfg *config.Config
	)

	BeforeEach(func() {
		cfg = config.Configuration()
	})

	Describe("TokenSubject", func() {
		var (
			token, email string
		)

		BeforeEach(func() {
			email = faker.Email()
			token, _ = newTokenWithClaims(cfg, &jwt.StandardClaims{
				Subject: email,
			})
		})

		It("returns the subject", func() {
			e, _ := TokenSubject(cfg, token)
			Expect(*e).To(Equal(email))
		})
	})

	Describe("IsValidToken", func() {
		var (
			err   error
			token string
		)

		When("the token is not a JWT", func() {
			BeforeEach(func() {
				token = "lk;jsdlkjlskfdj"
			})

			It("is false", func() {
				r, _ := IsValidToken(cfg, token)
				Expect(r).To(BeFalse())
			})

			It("reports an error", func() {
				_, e := IsValidToken(cfg, token)
				Expect(e).To(HaveOccurred())
			})
		})

		When("the token is a JWT", func() {
			var (
				claims *jwt.StandardClaims
			)

			JustBeforeEach(func() {
				token, err = newTokenWithClaims(cfg, claims)
				Expect(err).NotTo(HaveOccurred())
			})

			When("but has expired", func() {
				BeforeEach(func() {
					claims = &jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour * -24).Unix(),
					}
				})

				It("is false", func() {
					r, _ := IsValidToken(cfg, token)
					Expect(r).To(BeFalse())
				})

				It("reports an error", func() {
					_, e := IsValidToken(cfg, token)
					Expect(e).To(HaveOccurred())
				})
			})

			When("and is valid", func() {
				BeforeEach(func() {
					claims = &jwt.StandardClaims{}
				})

				It("is false", func() {
					r, _ := IsValidToken(cfg, token)
					Expect(r).To(BeTrue())
				})

				It("reports an error", func() {
					_, e := IsValidToken(cfg, token)
					Expect(e).NotTo(HaveOccurred())
				})
			})
		})
	})
})
