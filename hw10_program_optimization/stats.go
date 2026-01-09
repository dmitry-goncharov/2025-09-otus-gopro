package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	dotdomain := "." + domain
	for scanner.Scan() {
		user := &User{}
		if err := user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return nil, fmt.Errorf("get user error: %w", err)
		}
		if strings.HasSuffix(user.Email, dotdomain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return result, nil
}
