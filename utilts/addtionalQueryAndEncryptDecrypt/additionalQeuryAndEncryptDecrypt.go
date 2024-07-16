package addtionalQueryAndEncryptDecrypt

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"sendMail_git/utilts/decrypt"
	"sync"
)

type additionalTeamleadInfo struct {
	TeamleadUsername  string
	TeamleadFirstname string
	TeamleadSurname   string
}

func AdditionalQueryTeamlead(db *sql.DB, teamleadID string) (additionalTeamleadInfo, error) {
	var teamleadinfo additionalTeamleadInfo
	if teamleadID == "" {
		return teamleadinfo, errors.New("teammleadID must not be empty")
	}
	var teamleadIDforLog string
	err := db.QueryRow("SELECT orgmb_id FROM organize_member WHERE orgmb_id = $1", teamleadID).Scan(&teamleadIDforLog)
	if err != nil {
		log.Printf("orgmb_id does not match. Error: %v", err)
		return teamleadinfo, errors.New("teamleadID does not match")
	}
	// var teamleadEmail, teamleadFirstname, teamleadSurname string
	err = db.QueryRow("SELECT orgmb_email, orgmb_name, orgmb_surname FROM organize_member WHERE orgmb_id = $1", teamleadIDforLog).Scan(&teamleadinfo.TeamleadUsername, &teamleadinfo.TeamleadFirstname, &teamleadinfo.TeamleadSurname)
	if err != nil {
		log.Printf("Failed to query data By teamleadID(AdditionalQueryTeamlead). Error: %v", err)
		return teamleadinfo, errors.New("failed to query data By teamleadID(AdditionalQueryTeamlead)")
	}
	var (
		usernameDecrypt  []byte
		firstnameDecrypt []byte
		surnameDecrypt   []byte
		errChan          = make(chan error, 3)
		mutex            sync.Mutex
	)
	go func() {
		detokenizeUsername, err := decrypt.DetokenizationEmailForMasking(teamleadinfo.TeamleadUsername)
		if err != nil {
			errChan <- err
			return
		}
		usernameDecrypt, err = base64.StdEncoding.DecodeString(detokenizeUsername)
		errChan <- err // return nil when no error occurs
	}()
	go func() {
		detokenizeFirstname, err := decrypt.Detokenize(teamleadinfo.TeamleadFirstname)
		if err != nil {
			errChan <- err
			return
		}
		firstnameDecrypt, err = base64.StdEncoding.DecodeString(detokenizeFirstname)
		errChan <- err // return nil when no error occurs
	}()
	go func() {
		detokenizeSurname, err := decrypt.Detokenize(teamleadinfo.TeamleadSurname)
		if err != nil {
			errChan <- err
			return
		}
		surnameDecrypt, err = base64.StdEncoding.DecodeString(detokenizeSurname)
		errChan <- err // return nil when no error occurs
	}()
	for i := 0; i < 3; i++ {
		if err := <-errChan; err != nil {
			log.Println(err)
			return teamleadinfo, err
		}
	}
	mutex.Lock()
	// defer mutex.Unlock()
	teamleadinfo.TeamleadUsername = string(usernameDecrypt)
	teamleadinfo.TeamleadFirstname = string(firstnameDecrypt)
	teamleadinfo.TeamleadSurname = string(surnameDecrypt)
	defer mutex.Unlock()
	return teamleadinfo, nil
}

// func EmailForLineID(db *sql.DB, email_line string) (emailForLineIDResponse, error) { // move this logic into repositoriesvalidateOTP
// 	var emailLineIDInfo emailForLineIDResponse
// 	if email_line == "" {
// 		return emailLineIDInfo, errors.New("email must not be empty")
// 	}
// 	cipherUsername, err := encrypt.SendToFortanixSDKMSTokenizationEmailForMasking(email_line, keyUsername, keyPassword)
// 	if err != nil {
// 		log.Println(err)
// 		return emailLineIDInfo, err
// 	}
// 	err = db.QueryRow("SELECT orgmb_id, orgmb_line_id FROM organize_member WHERE orgmb_email = $1", cipherUsername).Scan(&emailLineIDInfo.ID, &emailLineIDInfo.LineID)
// 	if err != nil {
// 		log.Printf("ID or lineID does not match. Error: %v", err)
// 		return emailLineIDInfo, errors.New("id or lineID does not match")
// 	}
// 	return emailLineIDInfo, nil

// }

func CountTables(db *sql.DB) {
	var count int
	query := `SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public'`
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatalf("Failed to query table count: %s", err.Error())
	}
	fmt.Printf("There are %d tables in the database.\n", count)
}
