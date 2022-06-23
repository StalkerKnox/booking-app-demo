package database

import (
	"betting-app/models"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func ConnectDB() *sqlx.DB {
	db, err := sqlx.Connect("mysql", "root:OvjAcbmOh4E@(localhost:3306)/betting_app")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var DB = ConnectDB()

func InsertToDB() error {

	for _, singleOffer := range models.Offers {
		_, insertErrOffers := DB.NamedExec(`INSERT INTO offers (number, tv_channel, offer_id, title, has_statistics, time) VALUES (:number, :tv_channel, :offer_id, :title, :has_statistics, :time)`, singleOffer)
		if insertErrOffers != nil {
			return insertErrOffers
		}

		for _, singleRate := range singleOffer.Rates {
			singleRate.OfferID = singleOffer.ID
			_, insertErrRates := DB.NamedExec(`INSERT INTO rates (offer_id, name, rate) VALUES (:offer_id, :name, :rate)`, singleRate)
			if insertErrRates != nil {
				return insertErrRates
			}
		}

	}
	return nil
}

func GetOffersFromDB() []models.Offer {

	rows, _ := DB.Queryx("SELECT number, tv_channel, offer_id, title, has_statistics, time FROM offers")
	for rows.Next() {
		err := rows.StructScan(&models.SingleOffer)
		models.IterationStorage = append(models.IterationStorage, models.SingleOffer)
		if err != nil {
			log.Fatalln(err)
		}

	}

	for _, singleoffer := range models.IterationStorage {
		var rows *sqlx.Rows
		rows, _ = DB.Queryx("SELECT offer_id, rate, name FROM rates WHERE offer_id = ? ", singleoffer.ID)
		for rows.Next() {

			err := rows.StructScan(&models.SingleRate)
			singleoffer.Rates = append(singleoffer.Rates, models.SingleRate)

			if err != nil {
				log.Fatal(err)
			}
		}
		models.GetOffersFromDB = append(models.GetOffersFromDB, singleoffer)

	}
	// fmt.Println(models.GetOffersFromDB)
	return models.GetOffersFromDB
}

func GetOfferFromDB(req int) error {

	rows, _ := DB.Queryx("SELECT offer_id, rate, name FROM rates WHERE offer_id = ? ", req)
	for rows.Next() {
		err := rows.StructScan(&models.RateFromDB)
		models.OfferFromDB.Rates = append(models.OfferFromDB.Rates, models.RateFromDB)
		if err != nil {
			log.Fatalln(err)
		}
	}
	err := DB.Get(&models.OfferFromDB, "SELECT * FROM offers WHERE offer_id = ?", req)
	if err != nil {
		return err
	}
	return nil
}

func InsertOfferIntoDB(req models.Offer) error {
	_, insertErr := DB.NamedExec(`INSERT INTO offers (number, tv_channel, offer_id, title, has_statistics, time) VALUES (:number, :tv_channel, :offer_id, :title, :has_statistics, :time)`, req)
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	for _, singleRate := range req.Rates {
		singleRate.OfferID = req.ID
		_, insertErr = DB.NamedExec(`INSERT INTO rates (offer_id, name, rate) VALUES (:offer_id, :name, :rate)`, singleRate)
		if insertErr != nil {
			log.Fatal(insertErr)
		}
	}
	return nil
}