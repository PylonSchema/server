package origin

import "net/mail"

func vaildEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
