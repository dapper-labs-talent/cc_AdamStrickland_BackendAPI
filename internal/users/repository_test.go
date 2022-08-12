package users

import (
	"github.com/adamstrickland/dapper-api/internal"
	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("users/repository.go", func() {
	var (
		cfg                  *config.Config
		email, insert, query string
		count, after, before int
		db                   *internal.Conn
	)

	BeforeEach(func() {
		cfg = config.Configuration()
		db, _ = internal.NewConnection(cfg)
		email = faker.Email()
		query = "SELECT COUNT(*) FROM users WHERE email = ?"
		insert = "INSERT INTO users (email) VALUES (?)"

		db.Raw(query, email).Scan(&before)
		Expect(before).To(BeNumerically("==", 0))
	})

	Describe("Create()", func() {
		When("a record with the same email already exists", func() {
			BeforeEach(func() {
				db.Exec(insert, email)

				db.Raw(query, email).Scan(&count)
				Expect(count).To(BeNumerically(">", 0))
			})

			It("does not create a record", func() {
				Create(cfg, &User{Email: email})

				db.Raw(query, email).Scan(&after)
				Expect(after).To(BeNumerically("==", count))
			})

			It("returns an error", func() {
				_, err := Create(cfg, &User{Email: email})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("FindByEmail()", func() {
		When("a record with the same email does not exist", func() {
			It("returns an error", func() {
				u, err := FindByEmail(cfg, email)
				Expect(err).To(HaveOccurred())
				Expect(u).To(BeNil())
			})
		})

		When("a record with the same email already exists", func() {
			BeforeEach(func() {
				db.Exec(insert, email)

				db.Raw(query, email).Scan(&count)
				Expect(count).To(BeNumerically(">", 0))
			})

			It("returns the record", func() {
				u, err := FindByEmail(cfg, email)
				Expect(err).NotTo(HaveOccurred())
				Expect(u.Email).To(Equal(email))
			})
		})
	})

	Describe("All()", func() {
		It("returns a slice containing each entry", func() {
			db.Exec(insert, email)

			db.Raw("SELECT COUNT(*) FROM users").Scan(&count)
			Expect(count).To(BeNumerically(">", 0))

			us, err := All(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(*us)).To(BeNumerically(">", 0))
			Expect(*us).To(HaveLen(count))
		})
	})

	Describe("Update()", func() {
		var (
			nu *User
		)

		JustBeforeEach(func() {
			nu = &User{
				Email:     email,
				FirstName: "Something",
				LastName:  "Different",
			}
		})

		When("a record with the same email does not exist", func() {
			BeforeEach(func() {
				db.Raw(query, email).Scan(&count)
				Expect(count).To(BeNumerically("==", 0))
			})

			It("returns an error", func() {
				_, err := Update(cfg, nu)
				Expect(err).To(HaveOccurred())
			})

			It("returns a user", func() {
				u, _ := Update(cfg, nu)
				Expect(u).To(BeNil())
			})
		})

		When("a record with the same email already exists", func() {
			BeforeEach(func() {
				db.Exec(insert, email)

				db.Raw(query, email).Scan(&count)
				Expect(count).To(BeNumerically(">", 0))
			})

			It("updates the record", func() {
				uu, _ := Update(cfg, nu)

				Expect(uu.Email).To(Equal(nu.Email))
				Expect(uu.FirstName).To(Equal(nu.FirstName))
				Expect(uu.LastName).To(Equal(nu.LastName))

				uu, _ = FindByEmail(cfg, nu.Email)

				Expect(uu.Email).To(Equal(nu.Email))
				Expect(uu.FirstName).To(Equal(nu.FirstName))
				Expect(uu.LastName).To(Equal(nu.LastName))
			})

			It("does not return an error", func() {
				_, err := Update(cfg, nu)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
