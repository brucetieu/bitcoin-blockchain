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

// CreateWallet ... Create a wallet to store an address and public/private key information
// @Summary      Create a wallet
// @Description  Create a wallet to store an address and public / private key information
// @Tags         Wallets
// @Success      201  {string}  string     "address"
// @Failure      404  {object}  HTTPError
// @Router       /blockchain/wallets [post]
func (wh *WalletHandler) CreateWallet(ctx *gin.Context) {
	log.Info("CreateWallet handler called")
	// Create wallet with private / public key pair
	wallet, err := wh.walletService.CreateWallet()
	if err != nil {
		log.Error("error creating wallet: ", err.Error())
		NewError(ctx, http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusCreated, gin.H{"address": wallet.Address})
	}
}

// GetWallet ... Get a wallet by address
// @Summary      Get a wallet
// @Description  Get a wallet by address
// @Tags         Wallets
// @Success      200  {object}  representations.Wallet
// @Failure      404  {object}  HTTPError
// @Router       /blockchain/wallets/{address} [get]
func (wh *WalletHandler) GetWallet(ctx *gin.Context) {
	address := ctx.Param("address")
	log.Infof("GetWallet handler called with address: %s", address)

	wallet, err := wh.walletService.GetWallet(address)

	// Don't display private key to user
	wallet.PrivateKey = nil

	if err != nil {
		log.Errorf("error getting wallet with address: %s %s", address, err.Error())
		NewError(ctx, http.StatusNotFound, err)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"wallet": wallet})
	}
}

// GetWallets ... Get all walletes
// @Summary      Get all wallets
// @Description  Get all wallets
// @Tags         Wallets
// @Success      200  {array}   representations.Wallet
// @Failure      500  {object}  HTTPError
// @Router       /blockchain/wallets [get]
func (wh *WalletHandler) GetWallets(ctx *gin.Context) {
	log.Info("GetWallets handler called")
	wallets, err := wh.walletService.GetWallets()

	// Don't display private key to user
	for i := 0; i < len(wallets); i++ {
		wallets[i].PrivateKey = nil
	}

	if err != nil {
		log.Error("error getting all wallets: ", err.Error())
		NewError(ctx, http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"wallets": wallets})
	}
}
