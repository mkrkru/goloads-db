package main

type ExtensionIDRequest struct {
	ExtensionID string `json:"extension_id"`
}

type BannerIDRequest struct {
	ID string `json:"id"`
}

type BannerGotInteractedRequest struct {
	BannerID    string `json:"banner_id"`
	ExtensionID string `json:"extension_id"`
}

type TelegramIDRequest struct {
	TelegramID int `json:"user_id"`
}

type MoneyResponse struct {
	Money    float64 `json:"money"`
	Username string  `json:"username"`
	PhotoURL string  `json:"photo_url"`
}

type BannerRequest struct {
	URL     string   `json:"url"`
	Domains []string `json:"domains"`
}

type NewUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Hash      string `json:"hash"`
	ID        int    `json:"id"`
	PhotoUrl  string `json:"photoUrl"`
	UserName  string `json:"username"`
}

type LinkExtensionIDRequest struct {
	UserID             int    `json:"user_id"`
	ExtensionIDRequest string `json:"extension_id_request"`
}

type AdvertiserBannerResponse struct {
	BannerID string   `json:"id"`
	URL      string   `json:"redirect"`
	Domains  []string `json:"domains"`
	Image    string   `json:"image"`
}

type CookieResponse struct {
	UserCookie string `json:"userCookie"`
}
