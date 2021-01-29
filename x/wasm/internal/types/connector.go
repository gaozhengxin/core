package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

// DefaultFeatures - Cosmwasm feature
const DefaultFeatures = "staking,terra"

// ParseEvents converts wasm EventAttribute into an sdk.Events (with 0 or 1 elements)
func ParseEvents(logs []wasmvmtypes.EventAttribute, contractAddr sdk.AccAddress) sdk.Events {
	if len(logs) == 0 {
		return nil
	}

	// we always tag with the contract address issuing this event
	attrs := []sdk.Attribute{sdk.NewAttribute(AttributeKeyContractAddress, contractAddr.String())}
	for _, l := range logs {
		// and reserve the contract_address key for our use (not contract)
		if l.Key != AttributeKeyContractAddress {
			attr := sdk.NewAttribute(l.Key, l.Value)
			attrs = append(attrs, attr)
		}
	}

	return sdk.Events{sdk.NewEvent(EventTypeFromContract, attrs...)}
}

// ParseToCoin converts wasm coin to sdk.Coin
func ParseToCoin(wasmCoin wasmvmtypes.Coin) (coin sdk.Coin, err error) {
	amount, ok := sdk.NewIntFromString(wasmCoin.Amount)
	if !ok {
		err = sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("Failed to parse %s", coin.Amount))
		return
	}

	coin = sdk.Coin{
		Denom:  wasmCoin.Denom,
		Amount: amount,
	}
	return
}

// ParseToCoins converts wasm coins to sdk.Coins
func ParseToCoins(wasmCoins []wasmvmtypes.Coin) (coins sdk.Coins, err error) {
	for _, coin := range wasmCoins {
		c, err := ParseToCoin(coin)
		if err != nil {
			return nil, err
		}

		coins = append(coins, c)
	}
	return
}

// EncodeSdkCoin - encode sdk coin to wasm coin
func EncodeSdkCoin(coin sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.String(),
	}
}

// EncodeSdkCoins - encode sdk coins to wasm coins
func EncodeSdkCoins(coins sdk.Coins) wasmvmtypes.Coins {
	encodedCoins := make(wasmvmtypes.Coins, len(coins))
	for i, c := range coins {
		encodedCoins[i] = EncodeSdkCoin(c)
	}
	return encodedCoins
}
