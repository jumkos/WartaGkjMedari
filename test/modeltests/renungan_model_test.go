package modeltests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jumkos/WartaGkjMedari/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllRenungan(t *testing.T) {

	err := refreshUserAndRenunganTable()
	if err != nil {
		log.Fatalf("Error refreshing user and renungan table %v\n", err)
	}
	_, _, err = seedUsersAndRenungan()
	if err != nil {
		log.Fatalf("Error seeding user and renungan table %v\n", err)
	}
	renungan, err := renunganInstance.FindAllRenungan(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the renungan: %v\n", err)
		return
	}
	assert.Equal(t, len(*renungan), 2)
}

func TestSaveRenungan(t *testing.T) {

	err := refreshUserAndRenunganTable()
	if err != nil {
		log.Fatalf("Error user and renungan refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newRenungan := models.Renungan{
		ID:       1,
		Title:    "This is the title",
		Content:  "This is the content",
		AuthorID: user.ID,
	}
	savedRenungan, err := newRenungan.SaveRenungan(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the renungan: %v\n", err)
		return
	}
	assert.Equal(t, newRenungan.ID, savedRenungan.ID)
	assert.Equal(t, newRenungan.Title, savedRenungan.Title)
	assert.Equal(t, newRenungan.Content, savedRenungan.Content)
	assert.Equal(t, newRenungan.AuthorID, savedRenungan.AuthorID)

}

func TestGetRenunganByID(t *testing.T) {

	err := refreshUserAndRenunganTable()
	if err != nil {
		log.Fatalf("Error refreshing user and renungan table: %v\n", err)
	}
	post, err := seedOneUserAndOneRenungan()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundRenungan, err := renunganInstance.FindRenunganByID(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundRenungan.ID, post.ID)
	assert.Equal(t, foundRenungan.Title, post.Title)
	assert.Equal(t, foundRenungan.Content, post.Content)
}

func TestUpdateARenungan(t *testing.T) {

	err := refreshUserAndRenunganTable()
	if err != nil {
		log.Fatalf("Error refreshing user and renungan table: %v\n", err)
	}
	renungan, err := seedOneUserAndOneRenungan()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	postUpdate := models.Renungan{
		ID:       1,
		Title:    "modiUpdate",
		Content:  "modiupdate@gmail.com",
		AuthorID: renungan.AuthorID,
	}
	updatedRenungan, err := postUpdate.UpdateARenungan(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedRenungan.ID, postUpdate.ID)
	assert.Equal(t, updatedRenungan.Title, postUpdate.Title)
	assert.Equal(t, updatedRenungan.Content, postUpdate.Content)
	assert.Equal(t, updatedRenungan.AuthorID, postUpdate.AuthorID)
}

func TestDeleteARenungan(t *testing.T) {

	err := refreshUserAndRenunganTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	renungan, err := seedOneUserAndOneRenungan()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := renunganInstance.DeleteARenungan(server.DB, renungan.ID, renungan.AuthorID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
