package handlers

import (
	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	log.Info("CreateWallet handler called")
	// Create wallet with private / public key pair
	wallet, err := wh.walletService.CreateWallet()
	if err != nil {
		log.Error("error creating wallet: ", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"address": wallet.Address})
	}
}

func (wh *WalletHandler) GetWallet(c *gin.Context) {
	address := c.Param("address")
	log.Infof("GetWallet handler called with address: %s", address)

	wallet, err := wh.walletService.GetWallet(address)

	// Don't display private key to user
	wallet.PrivateKey = nil

	if err != nil {
		log.Errorf("error getting wallet with address: %s %s", address, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"wallet": wallet})
	}
}

func (wh *WalletHandler) GetWallets(c *gin.Context) {
	log.Info("GetWallets handler called")
	wallets, err := wh.walletService.GetWallets()

	// Don't display private key to user
	for i := 0; i < len(wallets); i++ {
		wallets[i].PrivateKey = nil
	}

	if err != nil {
		log.Error("error getting all wallets: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"wallets": wallets})
	}
}
