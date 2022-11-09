package formaterror

import (
	"errors"
	"fmt"
	"strings"
)


func FormatError(err string) error{
	if strings.Contains(err, "nickname") {
		return errors.New("Nickname Already Taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}

	if strings.Contains(err, "title"){
		return errors.New("Title Already Taken")
	}

	if strings.Contains(err, "hashPassword"){
		return errors.New("Incorrect Password")
	}

	fmt.Println(err)

	return errors.New("Incorrect Details")
}