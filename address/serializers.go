package address

import (
	"time"

	"github.com/google/uuid"
)

type AddressSerializer struct {
	Address Address
}

type AddressesSerializer struct {
	Addresses []Address
}

type AddressResponse struct {
	ID       uuid.UUID   `json:"id"`
	Address  string      `json:"address"`
	WalletID uuid.UUID   `json:"wallet_id"`
	Currency string      `json:"currency"`
	CreatedAt time.Duration `json:"created_at"`
	UpdatedAt time.Duration `json:"updated_at"`
}

func (a *AddressSerializer) Response() AddressResponse {
	address := AddressResponse{
		ID:       a.Address.ID,
		Address:  a.Address.Address,
		WalletID: a.Address.WalletID,
		Currency: a.Address.Currency,
		CreatedAt: time.Since(a.Address.CreatedAt),
		UpdatedAt: time.Since(a.Address.UpdatedAt),
	}
	return address
}

func (a *AddressesSerializer) Response() []AddressResponse {
	var addresses []AddressResponse
	for _, address := range a.Addresses {
		address := AddressResponse{
			ID:       address.ID,
			Address:  address.Address,
			WalletID: address.WalletID,
			Currency: address.Currency,
			CreatedAt: time.Since(address.CreatedAt),
			UpdatedAt: time.Since(address.UpdatedAt),
		}
		addresses = append(addresses, address)
	}
	return addresses
}