package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/cookiejar"
	"strconv"
	"strings"
	"time"
)

const (
	MoneyForView  = 0.1
	MoneyForClick = 0.3
)

type AdsServer struct {
	userStorage      UserStorage
	bannerStorage    BannerStorage
	analyticsStorage AnalyticsStorage
	goloadsDB        *sql.DB
}

type test struct {
	Body string `json:"body"`
}

var Test = test{Body: "OK"}

var ClicksMap = make(map[BannerGotInteractedRequest]bool)
var ViewsMap = make(map[BannerGotInteractedRequest]bool)

var counter int = 0

var DB = InitializeDB()

func ClickerCounter(a *AdsServer) {
	for key := range ClicksMap {
		a.userStorage.addMoney(a.userStorage.returnUserIDFromExtensionID(key.ExtensionID), MoneyForClick)
	}
	ClicksMap = map[BannerGotInteractedRequest]bool{}
	time.Sleep(360000)
}

func ViewCounter(a *AdsServer) {
	for key := range ViewsMap {
		a.userStorage.addMoney(a.userStorage.returnUserIDFromExtensionID(key.ExtensionID), MoneyForView)
	}
	ViewsMap = map[BannerGotInteractedRequest]bool{}
	time.Sleep(120000)
}

func checkForError(err error, errorCode int, w http.ResponseWriter) {
	if err != nil {
		if errorCode == 0 {
			fmt.Println(err)
			if errorCode != 0 {
				http.Error(w, http.StatusText(errorCode), errorCode)
			}
			return
		}
	}
}

func returnHTTPError(errorCode int, w http.ResponseWriter) {
	http.Error(w, http.StatusText(errorCode), errorCode)
}

func PreInnitiallizeStuff(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got request with method", r.Method, counter, "URL:", r.URL)
	counter++
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (a *AdsServer) sendExtensionIDHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		returnHTTPError(http.StatusBadRequest, w)
	}

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var id_request ExtensionIDRequest
	err = json.Unmarshal(rawData, &id_request)
	checkForError(err, http.StatusBadRequest, w)

	var id_response TelegramIDRequest
	id_response.TelegramID = a.userStorage.returnUserIDFromExtensionID(id_request.ExtensionID)

	bytes, err := json.Marshal(id_response)
	w.Write(bytes)
}

func (a *AdsServer) deleteBannerHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "DELETE" {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusInternalServerError, w)
	fmt.Println(string(rawData))

	var id_request BannerIDRequest
	err = json.Unmarshal(rawData, &id_request)
	checkForError(err, http.StatusBadRequest, w)

	a.bannerStorage.deleteAdvertisement(id_request.ID)

	bytes, err := json.Marshal(Test)
	checkForError(err, http.StatusInternalServerError, w)
	_, err = w.Write(bytes)
	if err != nil {
		return
	}
	// fmt.Fprint(w, string(bytes))
}

func (a *AdsServer) sendBannerHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	ads := a.bannerStorage.getRandomBanner()

	bytes, err := json.Marshal(ads)
	checkForError(err, http.StatusInternalServerError, w)

	fmt.Fprint(w, string(bytes))
}

/*func (a *AdsServer) receivePostHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	rawBody, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(rawBody))
	checkForError(err, http.StatusBadRequest, w)
	var newBanner Banner
	err = json.Unmarshal(rawBody, &newBanner)
	if err != nil {
		fmt.Println(err)
		return
	}
	a.bannerStorage.addBanner(newBanner)

	a.bannerStorage.putBannerIntoDB(newBanner.BannerID)

	bytes, err := json.Marshal(Test)
	checkForError(err, http.StatusInternalServerError, w)

	fmt.Fprint(w, string(bytes))

}*/

func (a *AdsServer) bannerClickedHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)
	if r.Method != "POST" {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusOK)
		return
	}

	rawBody, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(rawBody))
	checkForError(err, http.StatusBadRequest, w)

	var addClick BannerGotInteractedRequest
	err = json.Unmarshal(rawBody, &addClick)
	fmt.Println(string(rawBody))

	user := a.userStorage.getUserByID(a.userStorage.returnUserIDFromExtensionID(addClick.ExtensionID))
	a.analyticsStorage.addClickToDB(addClick.BannerID, user.ID)
	ClicksMap[addClick] = true
}

func (a *AdsServer) sendAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	id := r.URL.Query().Get("id")
	analytics := a.analyticsStorage.AnalyticsMap[id]

	bytes, err := json.Marshal(analytics)
	checkForError(err, http.StatusInternalServerError, w)

	fmt.Fprint(w, string(bytes))

}

func (a *AdsServer) sendFaviconHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	http.ServeFile(w, r, "favicon.ico")
}

var newBanner Banner

func (a *AdsServer) receiveBannerFromAdmin1(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var newBannerRequest BannerRequest
	if err := json.Unmarshal(rawData, &newBannerRequest); err != nil {
		fmt.Println(err)
		return
	}

	newBanner.BannerID = RandomString(20)
	newBanner.Domains = newBannerRequest.Domains
	newBanner.DomainURL = newBannerRequest.URL
	newBanner.Image = ""
	newBanner.ImageBase64 = true

	var IDResponse BannerIDRequest
	IDResponse.ID = newBanner.BannerID
	bytes, err := json.Marshal(IDResponse)
	checkForError(err, http.StatusInternalServerError, w)

	a.bannerStorage.addBanner(newBanner)
	w.Write(bytes)
}

func (a *AdsServer) receiveBannerFromAdmin2(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	bannerID := r.Header.Get("tgId")

	imgType := r.Header.Get("Content-Type")

	filetype := strings.Split(imgType, "/")[1]
	filename := fmt.Sprintf("./images/%s.%s", bannerID, filetype)
	err = ioutil.WriteFile(filename, rawData, 0666)
	checkForError(err, http.StatusInternalServerError, w)

	a.bannerStorage.changeBannerImage(bannerID, filename)
	a.bannerStorage.putBannerIntoDB(bannerID)
	var newCookie CookieResponse
	newCookie.UserCookie = uuid.New().String()
	rawBytes, err := json.Marshal(newCookie)

	w.Write(rawBytes)
}

/*func (a *AdsServer) receiveBannerImageHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	rawData, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var newImage BannerRequest
	if err := json.Unmarshal(rawData, &newImage); err != nil {
		fmt.Println(err)
		return
	}

	var newBanner Banner
	newBanner.BannerID = RandomString(20)
	newBanner.Domains = newImage.Domains
	newBanner.DomainURL = newImage.URL
	newBanner.ImageBase64 = true

	a.bannerStorage.addBanner(newBanner)
	a.bannerStorage.putBannerIntoDB(newBanner.BannerID)

	w.WriteHeader(http.StatusOK)

}*/

func (a *AdsServer) getUserMoneyHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)
	if r.Method != "POST" {
		returnHTTPError(http.StatusBadRequest, w)
		return
	}

	var extensionRequest ExtensionIDRequest
	rawBytes, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(rawBytes))
	checkForError(err, http.StatusBadRequest, w)

	err = json.Unmarshal(rawBytes, &extensionRequest)
	checkForError(err, http.StatusBadRequest, w)
	fmt.Println(extensionRequest)

	userID := a.userStorage.returnUserIDFromExtensionID(extensionRequest.ExtensionID)
	user := a.userStorage.getUserByID(userID)

	var money MoneyResponse
	money.Money = GtToMoney(user.Gotubles, user.Gopeykis)
	money.Username = user.Username
	money.PhotoURL = user.PhotoURL
	bytes, err := json.Marshal(money)
	w.Write(bytes)
}

func (a *AdsServer) bannerWatchedHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)
	if r.Method != "POST" {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusOK)
		return
	}

	rawBody, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(rawBody))
	checkForError(err, http.StatusBadRequest, w)

	var addView BannerGotInteractedRequest
	err = json.Unmarshal(rawBody, &addView)

	user := a.userStorage.getUserByID(a.userStorage.returnUserIDFromExtensionID(addView.ExtensionID))
	a.analyticsStorage.addViewToDB(addView.BannerID, user.ID)
	ViewsMap[addView] = true
}

func (a *AdsServer) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		returnHTTPError(http.StatusBadRequest, w)
		return
	}

	rawBytes, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var NewUserRequest NewUserRequest
	err = json.Unmarshal(rawBytes, &NewUserRequest)
	checkForError(err, http.StatusBadRequest, w)

	var newUser User
	newUser.Firstname = NewUserRequest.FirstName
	newUser.Lastname = NewUserRequest.LastName
	newUser.ID = NewUserRequest.ID
	newUser.Account = NewUserRequest.ID
	newUser.Money = 0.0

	a.userStorage.addUserToDB(newUser)

}

func (a *AdsServer) sendMoneyToUserHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		returnHTTPError(http.StatusBadRequest, w)
		return
	}

	rawBytes, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var extensionRequest ExtensionIDRequest
	err = json.Unmarshal(rawBytes, &extensionRequest)
	checkForError(err, http.StatusBadRequest, w)

	userToSendMoney := a.userStorage.getUserByID(a.userStorage.returnUserIDFromExtensionID(extensionRequest.ExtensionID))

	var moneyAm = GtToMoney(
		a.userStorage.getUserByID(userToSendMoney.ID).Gotubles,
		a.userStorage.getUserByID(userToSendMoney.ID).Gopeykis)
	var statusOK = false
	response, err := sendMoneyToUser(userToSendMoney.ID, moneyAm)

	if err != nil || response.StatusCode != http.StatusOK {
		returnHTTPError(http.StatusInternalServerError, w)
		return
	} else {
		a.userStorage.resetUserMoney(userToSendMoney.ID)
		statusOK = true
	}

	Test.Body = "OK"

	a.analyticsStorage.addTransactionToDB(userToSendMoney.ID, moneyAm, statusOK)
	bytes, err := ioutil.ReadAll(response.Body)
	checkForError(err, http.StatusInternalServerError, w)
	fmt.Println(string(bytes))
	w.Write(bytes)
}

func (a *AdsServer) linkExtensionIDToUserHandler(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	rawBytes, err := ioutil.ReadAll(r.Body)
	checkForError(err, http.StatusBadRequest, w)

	var linkRequest LinkExtensionIDRequest
	err = json.Unmarshal(rawBytes, &linkRequest)
	checkForError(err, http.StatusBadRequest, w)

	a.userStorage.linkExtensionID(linkRequest.ExtensionIDRequest, linkRequest.UserID)
	w.WriteHeader(http.StatusOK)
}

func (a *AdsServer) sendAdvertiserBanners(w http.ResponseWriter, r *http.Request) {
	PreInnitiallizeStuff(w, r)

	if r.Method != "POST" {
		returnHTTPError(http.StatusBadRequest, w)
	}

	AdvertiserID, _ := strconv.Atoi(r.Header.Get("tgId"))

	banners := a.bannerStorage.getAdvertiserBanners(AdvertiserID)
	advertiserBanners := make([]AdvertiserBannerResponse, 0)
	for _, banner := range banners {
		advertiserBanners = append(advertiserBanners,
			AdvertiserBannerResponse{
				BannerID: banner.BannerID,
				URL:      banner.DomainURL,
				Domains:  banner.Domains,
				Image:    banner.Image,
			})
	}
	rawBytes, err := json.Marshal(advertiserBanners)
	checkForError(err, http.StatusInternalServerError, w)
	w.Write(rawBytes)
}

func main() {

	// initializing test objects

	TestAdvertisement1 := Banner{
		BannerID:    "nbn9ewnd",
		Image:       "https://klike.net/uploads/posts/2019-05/1556708032_1.jpg",
		DomainURL:   "yandex.ru",
		Domains:     []string{"stackoverflow.com"},
		ImageBase64: false,
	}

	TestAdvertisement2 := Banner{
		BannerID:     "dnkjfnor",
		Image:        "https://lh3.googleusercontent.com/proxy/CRGj8PbtxI-4VfWouAiAbClb0uTfRwrt6FxZhFVtigesM2xkSebu0mV2bKAw6G8Xzxsd3VwQhIuxGeUvDeS0-fz0imr7yVb6xb_UBwxg_X7gHkeMY0U",
		DomainURL:    "github.com",
		Domains:      nil,
		AdvertiserID: 0,
		ImageBase64:  false,
	}

	TestAdvertisement3 := Banner{
		BannerID:     "dniebskj",
		Image:        "https://turbaza.ru/images/bases/2051/c2c19a0c42130d967b1eb0ff376b6cf6.jpg",
		DomainURL:    "google.com",
		Domains:      nil,
		AdvertiserID: 0,
		ImageBase64:  false,
	}

	TestAdvertisementStorage := BannerStorage{map[string]Banner{
		TestAdvertisement1.BannerID: TestAdvertisement1,
		TestAdvertisement2.BannerID: TestAdvertisement2,
		TestAdvertisement3.BannerID: TestAdvertisement3}}

	arrayLength := 14
	TestAnalytics := Analytics{
		BannerID:     "nbn9ewnd",
		Clicks:       RandomArray(arrayLength),
		UniqueClicks: RandomArray(arrayLength),
		Views:        RandomArray(arrayLength),
		UniqueViews:  RandomArray(arrayLength),
	}
	TestAnalyticsStorage := AnalyticsStorage{map[string]Analytics{TestAnalytics.BannerID: TestAnalytics}}
	GoloAdsServer := AdsServer{UserStorage{}, TestAdvertisementStorage, TestAnalyticsStorage, InitializeDB()}

	// initializing http handlers

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(GoloAdsServer.sendBannerHandler))
	mux.Handle("/delete", http.HandlerFunc(GoloAdsServer.deleteBannerHandler))
	mux.Handle("/add/image", http.HandlerFunc(GoloAdsServer.receiveBannerFromAdmin2))
	mux.Handle("/favicon.ico", http.HandlerFunc(GoloAdsServer.sendFaviconHandler))
	mux.Handle("/add", http.HandlerFunc(GoloAdsServer.receiveBannerFromAdmin1))
	mux.Handle("/analytics", http.HandlerFunc(GoloAdsServer.sendAnalyticsHandler))
	mux.Handle("/clicked", http.HandlerFunc(GoloAdsServer.bannerClickedHandler))
	mux.Handle("/watched", http.HandlerFunc(GoloAdsServer.bannerWatchedHandler))
	mux.Handle("/info/get", http.HandlerFunc(GoloAdsServer.getUserMoneyHandler))
	mux.Handle("/info/withdraw", http.HandlerFunc(GoloAdsServer.sendMoneyToUserHandler))
	mux.Handle("/user", http.HandlerFunc(GoloAdsServer.sendExtensionIDHandler))
	mux.Handle("/register", http.HandlerFunc(GoloAdsServer.registerUserHandler))
	mux.Handle("/link", http.HandlerFunc(GoloAdsServer.linkExtensionIDToUserHandler))
	mux.Handle("/banners", http.HandlerFunc(GoloAdsServer.sendAdvertiserBanners))

	go ClickerCounter(&GoloAdsServer)
	go ViewCounter(&GoloAdsServer)

	log.Fatal(http.ListenAndServeTLS("doats.ml:8080", "certificate.crt", "private.key", mux))
}
