package repository

import (
	"fmt"

	"synapsis-online-store/services/cart-service/internal/proto"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CartRepository handles database operations for cart-related actions
type CartRepository struct {
	db *pgxpool.Pool
}

// NewCartRepository initializes a new CartRepository
func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

// AddToCart adds a product to the user's cart
func (repo *CartRepository) AddToCart(userID, productID string, quantity int32) error {
	// Check if the product already exists in the user's cart, then update or insert
	_, err := repo.db.Exec(
		"INSERT INTO cart (user_id, product_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = quantity + $3",
		userID, productID, quantity,
	)
	if err != nil {
		return fmt.Errorf("could not add to cart: %w", err)
	}
	return nil
}

// RemoveFromCart removes a product from the user's cart
func (repo *CartRepository) RemoveFromCart(userID, productID string) error {
	_, err := repo.db.Exec("DELETE FROM cart WHERE user_id = $1 AND product_id = $2", userID, productID)
	if err != nil {
		return fmt.Errorf("could not remove from cart: %w", err)
	}
	return nil
}

// GetCart retrieves all items in a user's cart
func (repo *CartRepository) GetCart(userID string) ([]*proto.CartItem, error) {
	rows, err := repo.db.Query("SELECT product_id, quantity FROM cart WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("could not get cart: %w", err)
	}
	defer rows.Close()

	var items []*proto.CartItem
	for rows.Next() {
		var item proto.CartItem
		if err := rows.Scan(&item.ProductId, &item.Quantity); err != nil {
			return nil, fmt.Errorf("could not scan cart item: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}
