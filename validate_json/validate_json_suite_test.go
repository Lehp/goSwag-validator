package validate_json_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/lehp/goswag-validator/validate_json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSchemeComparison(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "compareScheme Suite")
}

var _ = Describe("json with scheme comparison", func() {
	It("should return true on res fitting scheme", func(){
        scheme, _ := os.ReadFile("tests/test_mocks/product_swagger.json")
        res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")

        resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "kats")

		Expect(swaggy).To(Equal(true))
	})

	It("should fail on missing attribute in json", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/product_swagger_broken.json")
        res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")

		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "kats")

		Expect(swaggy).To(Equal(false))
	})

	It("should fail on missing attribute in json array", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/product_swagger_broken_arr.json")
        res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")

		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "kats")

		Expect(swaggy).To(Equal(true))
	})
})