package validate

func CheckLuhn(num string) bool {
	var sum int
	parity := (len(num)) % 2

	for i := 0; i < len(num); i++ {
		digit := int(num[i] - '0')

		if i%2 == parity {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return sum%10 == 0
}
