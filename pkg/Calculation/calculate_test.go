package calculate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostpone(t *testing.T) {
	nums := []float64{1, 2, 3, 4}
	result := Postpone(nums, 2)
	assert.Equal(t, []float64{1, 2, 4}, result)
}

func TestPostponeStringSlice(t *testing.T) {
	strs := []string{"a", "b", "c", "d"}
	result := PostponeStringSlice(strs, 1)
	assert.Equal(t, []string{"a", "c", "d"}, result)
}

func TestCalcBasic_Valid(t *testing.T) {
	disableLogOutput()
	result, err, code := CalcBasic("2+3*4")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 14.0, result)
}

func TestCalcBasic_WithSpaces(t *testing.T) {
	disableLogOutput()
	result, err, code := CalcBasic(" 10 + 2 * 5 ")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 20.0, result)
}

func TestCalcBasic_InvalidCharacter(t *testing.T) {
	disableLogOutput()
	_, err, code := CalcBasic("3 + a")
	assert.Error(t, err)
	assert.Equal(t, 422, code)
}

func TestCalcBasic_DivideByZero(t *testing.T) {
	disableLogOutput()
	_, err, code := CalcBasic("5 / 0")
	assert.Error(t, err)
	assert.Equal(t, 422, code)
}

func TestCalcBasic_EmptyString(t *testing.T) {
	disableLogOutput()
	_, err, code := CalcBasic("   ")
	assert.Error(t, err)
	assert.Equal(t, 404, code)
}

func TestCalc_ValidExpression(t *testing.T) {
	disableLogOutput()
	result, err, code := Calc("(2+3)*(4+1)")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 25.0, result)
}

func TestCalc_MissingBrackets(t *testing.T) {
	disableLogOutput()
	_, err, code := Calc("(2+3")
	assert.Error(t, err)
	assert.Equal(t, 422, code)
}

func TestCalc_NestedBrackets(t *testing.T) {
	disableLogOutput()
	result, err, code := Calc("((1+2)*(3+4))")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 21.0, result)
}

func TestCalc_ExtraBrackets(t *testing.T) {
	disableLogOutput()
	_, err, code := Calc("((1+2)))))")
	assert.Error(t, err)
	assert.Equal(t, 422, code)
}

func disableLogOutput() {
	// prevent logging to file during tests
	_ = os.MkdirAll("../log", os.ModePerm)
	_ = os.WriteFile("../log/CalcBasicLog.txt", []byte(""), 0644)
	_ = os.WriteFile("../log/CalcLog.txt", []byte(""), 0644)
}
