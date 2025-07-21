package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type userGraphHTTP struct {
	BaseURL string
}

type UserGraphService interface {
	GetFollowings(userID string) ([]string, error)
	GetFollowers(userID string) ([]string, error)
}

func NewUserGraphHTTP(baseURL string) UserGraphService {
	return &userGraphHTTP{BaseURL: baseURL}
}

func (g *userGraphHTTP) GetFollowings(userID string) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/relations/%s/followings", g.BaseURL,userID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Followings []struct {
			FollowingID string `json:"following_id"`
		} `json:"followings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var followingIDs []string
	for _, f := range result.Followings {
		followingIDs = append(followingIDs, f.FollowingID)
	}

	return followingIDs, nil
}

func (g *userGraphHTTP) GetFollowers(userID string) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/relations/%s/followers", g.BaseURL, userID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Followers []struct {
			FollowerID string `json:"follower_id"`
		} `json:"followers"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var followerIDs []string
	for _, f := range result.Followers {
		followerIDs = append(followerIDs, f.FollowerID)
	}

	return followerIDs, nil
}