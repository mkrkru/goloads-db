package main

import (
	"math/rand"
)

type User struct {
	Firstname   string  `json:"firstname"`
	Lastname    string  `json:"lastname"`
	ID          int     `json:"id"`
	Account     int     `json:"account"`
	Money       float64 `json:"money"`
	PhotoURL    string  `json:"photo_url"`
	Username    string  `json:"username"`
	Hash        string  `json:"hash"`
	Gotubles    int     `json:"gotubles"`
	Gopeykis    int     `json:"gopeykis"`
	ExtensionID string  `json:"extension_id"`
}

type UserStorage struct {
	UserMap map[string]User
}

type Banner struct {
	BannerID     string   `json:"id"`
	Image        string   `json:"image"`
	DomainURL    string   `json:"url"`
	Domains      []string `json:"domains"`
	AdvertiserID int      `json:"advertiser_id"`
	ImageBase64  bool     `json:"image-base64"`
}

type Analytics struct {
	BannerID     string `json:"id"`
	Clicks       []int  `json:"clicks"`
	UniqueClicks []int  `json:"unique_clicks"`
	Views        []int  `json:"views"`
	UniqueViews  []int  `json:"unique_views"`
}

type BannerStorage struct {
	BannerMap map[string]Banner
}

type AnalyticsStorage struct {
	AnalyticsMap map[string]Analytics
}

func (a *BannerStorage) addBanner(ad Banner) {
	a.BannerMap[ad.BannerID] = ad
}

func (a *BannerStorage) getRandomBanner() Banner {
	var ads []Banner
	for _, ad := range a.BannerMap {
		ads = append(ads, ad)
	}
	return ads[rand.Intn(len(ads))]
}

func (a *BannerStorage) deleteAdvertisement(id string) {
	delete(a.BannerMap, id)
}

func (b BannerStorage) sendBanner(id string) Banner {
	return b.BannerMap[id]
}

func (b *BannerStorage) changeBannerImage(id string, image string) {
	tempBanner := Banner{
		BannerID:    b.BannerMap[id].BannerID,
		Image:       image,
		DomainURL:   b.BannerMap[id].DomainURL,
		Domains:     b.BannerMap[id].Domains,
		ImageBase64: b.BannerMap[id].ImageBase64,
	}

	b.BannerMap[id] = tempBanner
}

func (a AnalyticsStorage) getAnalytics(id string) Analytics {
	return a.AnalyticsMap[id]
}

func (a *AnalyticsStorage) addClick(id string) {
	a.AnalyticsMap[id].Clicks[0]++
}

func (a *AnalyticsStorage) addWatch(id string) {
	a.AnalyticsMap[id].Views[0]++
}
