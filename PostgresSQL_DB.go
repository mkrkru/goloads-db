package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"time"
)

/*const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)*/

const (
	host     = "doats.ml"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "goloads"
)

func GtToMoney(gt int, gp int) float64 {
	return float64(gt) + (float64(gp) / 100)
}

func MoneyToGT(money float64) (int, int) {
	gt := int(money)
	gp := int((money - float64(gt)) * 100)
	return gt, gp
}

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
	host, port, user, password, dbname)

func InitializeDB() *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func (b *BannerStorage) putBannerIntoDB(id string) {
	db := DB

	var banner Banner
	banner = b.BannerMap[id]
	_, err := db.Query(`INSERT INTO "Banners" 
					VALUES ($1, $2, $3, $4, $5);`,
		banner.BannerID,
		banner.DomainURL,
		banner.Image,
		pq.Array(banner.Domains),
		banner.ImageBase64,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Query(`INSERT INTO "Analytics" 
					VALUES ($1, $2, $3, $4, $5);`,
		banner.BannerID,
		pq.Array([]int{}),
		pq.Array([]int{}),
		pq.Array([]int{}),
		pq.Array([]int{}),
		8080,
	)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func (b *BannerStorage) getBannersFromDB() []Banner {
	db := DB

	var banners []Banner
	rows, err := db.Query(`SELECT * FROM "Banners"`)
	if err != nil {
		fmt.Println(err)
		return []Banner{}
	}
	i := 0
	for rows.Next() {
		err = rows.Scan(&banners[i].BannerID, &banners[i].Image, &banners[i].DomainURL, &banners[i].Domains)
		if err != nil {
			fmt.Println(err)
			return []Banner{}
		}
		i++
	}

	return banners
}

func (a *AnalyticsStorage) addClickToDB(banner_id string, user_id int) {
	db := DB

	_, err := db.Query(`INSERT INTO "Clicks" VALUES ($1, $2, $3)`, banner_id, user_id, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (a *AnalyticsStorage) addViewToDB(banner_id string, user_id int) {
	db := DB

	_, err := db.Query(`INSERT INTO "Views" VALUES ($1, $2, $3)`, banner_id, user_id, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (a AnalyticsStorage) addTransactionToDB(user_id int, money float64, statusOK bool) {
	db := DB

	_, err := db.Query(`INSERT INTO "Transactions" VALUES ($1, $2, $3, $4)`, RandomString(30), money, user_id, statusOK)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (a *UserStorage) addUserToDB(user User) {
	db := DB

	_, err := db.Query(`INSERT INTO "User" ("Firstname", "Lastname", ID, "Account") VALUES ($1, $2, $3, $4)`,
		user.Firstname,
		user.Lastname,
		user.ID,
		user.Account,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (a *UserStorage) getUserByID(telegramID int) User {
	db := DB

	var user User
	row := db.QueryRow(`SELECT * FROM "Users" WHERE ID=$1`, telegramID)
	if err := row.Scan(
		&user.Firstname,
		&user.Lastname,
		&user.ID,
		&user.Account,
		&user.Money,
		&user.PhotoURL,
		&user.Username,
		&user.Hash,
		&user.Gotubles,
		&user.Gopeykis,
		&user.ExtensionID,
	); err != nil {
		fmt.Println("error", err)
		return User{}
	}

	return user
}

func (u *UserStorage) resetUserMoney(telegramID int) {
	db := DB

	_, _ = db.Query(`UPDATE "Users"
			SET "Gotubles"=0, "Gopeykis"=0
			WHERE ID=$1`, telegramID)

}

func (u *UserStorage) addMoney(telegramID int, moneyAmount float64) {
	db := DB

	gt, gp := MoneyToGT(moneyAmount)
	_, _ = db.Query(`UPDATE "Users"
			SET "Gotubles"="Gotubles"+$1,
			    "Gopeykis"="Gopeykis"+$2
			WHERE ID=$3`,
		gt, gp, telegramID)

}

func (u *UserStorage) linkExtensionID(extension_id string, user_id int) {
	db := DB

	_, _ = db.Query(`UPDATE "Users" 
							  SET "ExtensionID"=$1 
							  WHERE ID=$2`, extension_id, user_id)

}

func (u *UserStorage) returnUserIDFromExtensionID(extensionID string) int {
	db := DB

	var telegramID int
	row := db.QueryRow(`SELECT ID FROM "Users" WHERE "ExtensionID"=$1`, extensionID)
	err := row.Scan(&telegramID)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return telegramID

}

func (u *UserStorage) replaceHash(newHash string, user_id int){
	db := DB

	_ = db.QueryRow(`UPDATE "Users"
							  SET "Hash"=$1
							  WHERE ID=$2`, newHash, user_id)
}

func (b *BannerStorage) getAdvertiserBanners(advertiser_id int) []Banner{
	db := DB

	rows, _ := db.Query(`SELECT * FROM "Banners" WHERE "AdvertiserID"=$1`, advertiser_id)
	banners := make([]Banner, 0)
	for rows.Next() {
		banner := Banner{}
		rows.Scan(&banner.BannerID, &banner.DomainURL, &banner.Image, &banner.Domains, &banner.ImageBase64, &banner.AdvertiserID)
		banners = append(banners, banner)
	}

	return banners
}

// func (a *BannerStorage) getAdvertisementFromDB (id string)
