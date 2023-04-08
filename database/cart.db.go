package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("sorry, we can't find the products")
	ErrCantDecodeProducts = errors.New("sorry, we can't find the products")
	ErrUserIdIsNotValid   = errors.New("sorry, this user is not valid")
	ErrCantUpdateUser     = errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("can't remove this item from the cart")
	ErrCantGetItems       = errors.New("unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("can't update the purchase")
)

func AddProductToCart()

func RemoveCartItem()

func BuyItemFromCart()

func InstantBuyer()
