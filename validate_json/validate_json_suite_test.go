package validate_json_test

import (
	"bytes"
	"fmt"
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
        scheme, schemeErr := os.ReadFile("tests/test_mocks/product_swagger.json")
        res, resErr := os.ReadFile("tests/test_mocks/product_res_body.json")

        if (len(scheme) == 0) {
            fmt.Println(schemeErr.Error())
            return
        }
        if (len(res) == 0) {
            fmt.Println(resErr.Error())
            return
        }

        resBuf := bytes.NewBuffer(res)

		var fits bool = validate_json.EquivalentToScheme(resBuf, scheme, "kats")

		Expect(fits).To(Equal(true))

		
	})
})