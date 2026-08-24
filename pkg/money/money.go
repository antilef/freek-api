package money

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type CLP int64

func (m CLP) String() string {
	return "$" + FormatDots(int64(m))
}

func (m CLP) Int64() int64 {
	return int64(m)
}

func FormatDots(n int64) string {
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}

	str := strconv.FormatInt(n, 10)
	if len(str) <= 3 {
		if negative {
			return "-" + str
		}
		return str
	}

	var result []string
	remaining := str
	for len(remaining) > 3 {
		result = append([]string{remaining[len(remaining)-3:]}, result...)
		remaining = remaining[:len(remaining)-3]
	}
	result = append([]string{remaining}, result...)

	joined := strings.Join(result, ".")
	if negative {
		return "-" + joined
	}
	return joined
}

func Parse(s string) (CLP, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "$")
	s = strings.ReplaceAll(s, ".", "")

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid CLP value: %s", s)
	}

	return CLP(n), nil
}

func DivideRoundUp(amount CLP, installments int) CLP {
	if installments <= 0 {
		return amount
	}
	perInstallment := int64(amount) / int64(installments)
	remainder := int64(amount) % int64(installments)
	if remainder > 0 {
		perInstallment++
	}
	return CLP(perInstallment)
}

func TotalWithInterest(amountCLP int64, interestRate float64) CLP {
	if interestRate <= 0 {
		return CLP(amountCLP)
	}
	return CLP(math.Round(float64(amountCLP) * (1 + interestRate)))
}

func MonthlyInstallment(total CLP, installments int) CLP {
	return DivideRoundUp(total, installments)
}
