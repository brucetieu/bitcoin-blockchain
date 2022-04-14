package handlers

import (
	"net/http"

	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type WalletHandler struct {
	walletService services.WalletService
}

func NewWalletHandler(walletService services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (wh *WalletHandler) CreateWallet(c *gin.Context) {
	// Create wallet with private / public key pair
	wallet, err := wh.walletService.CreateWallet()
	if err != nil {
		log.Error("error creating wallet: ", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"address": wallet.Address})
	}

}