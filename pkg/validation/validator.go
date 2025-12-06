package validation

import (
	"regexp"
	"unicode"
)

// Remove caracteres não numéricos
var reOnlyNumbers = regexp.MustCompile(`[^0-9]`)

func IsPasswordStrong(password string) bool {
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasNumber && hasSpecial
}

func IsCPFValid(cpf string) bool {
	cpf = reOnlyNumbers.ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return false
	}

	// Verifica se todos os dígitos são iguais (ex: 111.111.111-11)
	allEquals := true
	for i := 1; i < 11; i++ {
		if cpf[i] != cpf[0] {
			allEquals = false
			break
		}
	}
	if allEquals {
		return false
	}

	// Validação do 1º Dígito
	sum := 0
	for i := 0; i < 9; i++ {
		digit := int(cpf[i] - '0')
		sum += digit * (10 - i)
	}
	remainder := (sum * 10) % 11
	if remainder == 10 {
		remainder = 0
	}
	if remainder != int(cpf[9]-'0') {
		return false
	}

	// Validação do 2º Dígito
	sum = 0
	for i := 0; i < 10; i++ {
		digit := int(cpf[i] - '0')
		sum += digit * (11 - i)
	}
	remainder = (sum * 10) % 11
	if remainder == 10 {
		remainder = 0
	}
	if remainder != int(cpf[10]-'0') {
		return false
	}

	return true
}
