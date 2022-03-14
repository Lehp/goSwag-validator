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
	It("should return true on res fitting scheme", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/product_swagger.json")
		res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")
		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "kats", []string{})

		Expect(swaggy).To(Equal(true))
	})

	It("should fail on missing attribute in json", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/product_swagger_broken.json")
		res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")
		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "kats", []string{})

		Expect(swaggy).To(Equal(false))
	})

	It("should fail on missing attribute in json array", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/product_swagger_broken_arr.json")
		res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")
		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "kats", []string{})

		Expect(swaggy).To(Equal(false))
	})

	It("should fail on wrong path", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/product_swagger.json")
		res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")
		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "???", []string{})

		Expect(swaggy).To(Equal(false))
	})

	It("should ignore field", Focus, func() {
		scheme, _ := os.ReadFile("tests/test_mocks/ignore_field.json")
		res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")
		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "kats", []string{"ignore"})

		Expect(swaggy).To(Equal(true))
	})
})

var _ = Describe("json with scheme comparison gomega", func() {
	It("should fail on missing attribute in json", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/product_swagger_broken.json")
		res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")
		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "kats", []string{})

		Expect(swaggy).To(Equal(false))
	})

	It("should fail on missing attribute in json array", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/product_swagger_broken_arr.json")
		res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")

		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "kats", []string{})

		Expect(swaggy).To(Equal(false))
	})

	It("should fail on wrong path", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/product_swagger.json")
		res, _ := os.ReadFile("tests/test_mocks/product_res_body.json")

		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "???", []string{})

		Expect(swaggy).To(Equal(false))
	})
})

var _ = Describe("json with scheme comparison", func() {
	It("should compare nested swagfile", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/nested/cartService.json")
		res, _ := os.ReadFile("tests/test_mocks/nested/updateCartRes.json")

		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "com.UpdateQuantityDto", []string{})

		Expect(swaggy).To(Equal(true))
	})

	It("should compare array swagfile", func() {
		scheme, _ := os.ReadFile("tests/test_mocks/size/sizeSwagger.json")
		res, _ := os.ReadFile("tests/test_mocks/size/sizeRes.json")

		resBuf := bytes.NewBuffer(res)

		var swaggy bool = validate_json.EquivalentToScheme(resBuf, scheme, "com.ProductSizeV2", []string{})

		Expect(swaggy).To(Equal(true))
	})
})
