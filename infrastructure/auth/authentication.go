package auth

import (
	"context"
	"errors"
	"fmt"
	"food-api/infrastructure/auth/model"
	"github.com/go-redis/redis/v8"
	"time"
)

type ClientData struct {
	client *redis.Client
}

type InterfaceAuth interface {
	CreateAuth(ctx context.Context, userId string, details *model.TokenDetails) error
	FetchAuth(ctx context.Context, tokenUuid string) (string, error)
	DeleteRefresh(ctx context.Context, refreshUuid string) error
	DeleteTokens(ctx context.Context, details *model.AccessDetails) error
}

func NewAuth(client *redis.Client) *ClientData {
	return &ClientData{client: client}
}

var _ InterfaceAuth = &ClientData{}

//CreateAuth Save token metadata to Redis
func (cl *ClientData) CreateAuth(ctx context.Context, userId string, details *model.TokenDetails) error {
	atExpires := time.Unix(details.AtExpires, 0) //converting Unix to UTC(to Time object)
	rtExpires := time.Unix(details.RtExpires, 0)
	now := time.Now()

	atCreated, err := cl.client.Set(ctx, details.TokenUuid, userId, atExpires.Sub(now)).Result()
	if err != nil {
		return err
	}

	rtCreated, err := cl.client.Set(ctx, details.RefreshUuid, userId, rtExpires.Sub(now)).Result()
	if err != nil {
		return err
	}

	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}

	return nil
}

// FetchAuth Get authentication
func (cl *ClientData) FetchAuth(ctx context.Context, tokenUuid string) (string, error) {
	userId, err := cl.client.Get(ctx, tokenUuid).Result()
	if err != nil {
		return "", err
	}

	return userId, nil
}

// DeleteRefresh delete refresh token
func (cl *ClientData) DeleteRefresh(ctx context.Context, refreshUuid string) error {
	deleted, err := cl.client.Del(ctx, refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}

	return nil
}

// DeleteTokens Once a user row in the token table
func (cl *ClientData) DeleteTokens(ctx context.Context, details *model.AccessDetails) error {
	// Get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%s", details.TokenUuid, details.UserId)

	// Delete access token
	deletedAt, err := cl.client.Del(ctx, details.TokenUuid).Result()
	if err != nil {
		return err
	}

	// Delete refresh token
	deletedRt, err := cl.client.Del(ctx, refreshUuid).Result()
	if err != nil {
		return err
	}

	// When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}

	return nil
}