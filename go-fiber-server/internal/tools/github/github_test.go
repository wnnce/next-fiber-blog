package github

import (
	"fmt"
	"go-fiber-ent-web-layout/pkg/pool"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	accessToken, err := AccessToken("e2b4447c8eb73c248d2d")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(accessToken)
	ch1 := make(chan map[string]string)
	ch2 := make(chan []*Email)
	pool.Go(func() {
		profile, _ := UserProfile(accessToken)
		ch1 <- profile
		close(ch1)
	})
	pool.Go(func() {
		emails, _ := UserEmails(accessToken)
		ch2 <- emails
		close(ch2)
	})
	profile := <-ch1
	emails := <-ch2
	fmt.Printf("%v\n", profile)
	fmt.Printf("%v\n", emails)
}
