package arango_test

import (
	"context"
	"os"
	"testing"

	arango "github.com/joselitofilho/gorm-arango/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
	Age   uint
}

var gormDB *gorm.DB

func newArangoDBTestConfig() *arango.Config {
	arangodbUri := os.Getenv("ARANGODB_URI")
	if arangodbUri == "" {
		arangodbUri = "http://localhost:8529"
	}
	return &arango.Config{
		URI:                  arangodbUri,
		User:                 "user",
		Password:             "password",
		Database:             "gorm-arango-test",
		Timeout:              120,
		MaxConnectionRetries: 10,
	}
}

var _ = BeforeSuite(func() {
	arangodbUri := os.Getenv("ARANGODB_URI")
	if arangodbUri == "" {
		arangodbUri = "http://localhost:8529"
	}
	arangodbConfig := &arango.Config{
		URI:                  arangodbUri,
		User:                 "user",
		Password:             "password",
		Database:             "gorm-arango-test",
		Timeout:              120,
		MaxConnectionRetries: 10,
	}

	By("Connecting to the ArangoDB", func() {
		db, err := gorm.Open(arango.Open(arangodbConfig), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())
		gormDB = db
	})
})

var _ = AfterSuite(func() {
	dialector := gormDB.Dialector.(arango.Dialector)
	Expect(dialector).NotTo(BeNil())
	err := dialector.Database.Remove(context.Background())
	Expect(err).NotTo(HaveOccurred())
})

func TestArangoSuite(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "ArangoDB Suite")
}
