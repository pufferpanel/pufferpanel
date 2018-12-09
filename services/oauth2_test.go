package services

import (
	"github.com/pufferpanel/pufferpanel/models"
	"testing"
)

func Test(t *testing.T) {
	service, err := GetOAuthService()
	if err != nil {
		t.Error(err)
		return
	}

	ci := &models.ClientInfo{ClientID: "test"}
	server := &models.Server{Name: "test"}

	err = service.UpdateScopes(ci, server, "newscope2")
	if err != nil {
		t.Error(err)
		return
	}
}
