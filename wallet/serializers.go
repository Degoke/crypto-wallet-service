package wallet

import (
	"time"

	"github.com/Degoke/crypto-wallet-service/address"
	"github.com/google/uuid"
)

type WalletsSerializer struct {
	wallets []Wallet
}

type WalletSerializer struct {
	wallet Wallet
}

type WalletResponse struct {
	ID    uuid.UUID   `json:"id"`
	Addresses []address.AddressResponse `json:"addresses"`
	Network string `json:"network"`
	CreatedAt time.Duration `json:"created_at"`
	UpdatedAt time.Duration `json:"updated_at"`
}

func (w *WalletsSerializer) Response() []WalletResponse {
	var wallets []WalletResponse
	for _, wallet := range w.wallets {
		addressSerializer := address.AddressesSerializer{Addresses: wallet.Addresses}
		wallet := WalletResponse{
			ID:    wallet.ID,
			Addresses: addressSerializer.Response(),
			Network: wallet.Network,
			CreatedAt: time.Since(wallet.CreatedAt),
			UpdatedAt: time.Since(wallet.UpdatedAt),
		}
		wallets = append(wallets, wallet)
	}
	return wallets
}

func (w *WalletSerializer) Response() WalletResponse {
	addressSerializer := address.AddressesSerializer{Addresses: w.wallet.Addresses}
	wallet := WalletResponse{
		ID:    w.wallet.ID,
		Addresses: addressSerializer.Response(),
		Network: w.wallet.Network,
		CreatedAt: time.Since(w.wallet.CreatedAt),
		UpdatedAt: time.Since(w.wallet.UpdatedAt),
	}
	return wallet
}