package test

import (
	"testing"
)

func TestSendgridMail(t *testing.T) {
	defer setupTest()()

	htmlBasicContent := "<h1>Test</h1>"

	err := appServices.SendgridService.SendEmail("ndodanli14@gmail.com", "Test Subject", htmlBasicContent)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

}
