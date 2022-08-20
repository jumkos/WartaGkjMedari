package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jumkos/WartaGkjMedari/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateRenungan(t *testing.T) {

	err := refreshUserAndRenunganTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"The title", "content": "the content", "author_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			title:        "The title",
			content:      "the content",
			author_id:    user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"title":"The title", "content": "the content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "title Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"title":"When no token is passed", "content": "the content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"title":"When incorrect token is passed", "content": "the content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"title": "", "content": "The content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required Title",
		},
		{
			inputJSON:    `{"title": "This is a title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required Content",
		},
		{
			inputJSON:    `{"title": "This is an awesome title", "content": "the content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"title": "This is an awesome title", "content": "the content", "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/renungan", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateRenungan)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.Bytes()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetRenungan(t *testing.T) {

	err := refreshUserAndRenunganTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndRenungan()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/renungan", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetAllRenungan)
	handler.ServeHTTP(rr, req)
	
	var renungan []models.Renungan
	err = json.Unmarshal([]byte(rr.Body.Bytes()),&renungan)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(renungan), 2)
}
func TestGetRenunganByID(t *testing.T) {

	err := refreshUserAndRenunganTable()
	if err != nil {
		log.Fatal(err)
	}
	renungan, err := seedOneUserAndOneRenungan()
	if err != nil {
		log.Fatal(err)
	}
	renunganSample := []struct {
		id           string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(renungan.ID)),
			statusCode: 200,
			title:      renungan.Title,
			content:    renungan.Content,
			author_id:  renungan.AuthorID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range renunganSample {

		req, err := http.NewRequest("GET", "/renungan", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetRenungan)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.Bytes()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, renungan.Title, responseMap["title"])
			assert.Equal(t, renungan.Content, responseMap["content"])
			assert.Equal(t, float64(renungan.AuthorID), responseMap["author_id"]) //the response author id is float64
		}
	}
}

func TestUpdateRenungan(t *testing.T) {

	var RenunganUserEmail, RenunganUserPassword string
	var AuthRenunganAuthorID uint32
	var AuthRenunganID uint64

	err := refreshUserAndRenunganTable()
	if err != nil {
		log.Fatal(err)
	}
	users, renungan, err := seedUsersAndRenungan()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		RenunganUserEmail = user.Email
		RenunganUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(RenunganUserEmail, RenunganUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first post
	for _, post := range renungan {
		if post.ID == 2 {
			continue
		}
		AuthRenunganID = post.ID
		AuthRenunganAuthorID = post.AuthorID
	}
	// fmt.Printf("this is the auth post: %v\n", AuthRenunganID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthRenunganID)),
			updateJSON:   `{"title":"The updated post", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   200,
			title:        "The updated post",
			content:      "This is the updated content",
			author_id:    AuthRenunganAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthRenunganID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthRenunganID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Title 2" belongs to post 2, and title must be unique
			id:           strconv.Itoa(int(AuthRenunganID)),
			updateJSON:   `{"title":"Title 2", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "title Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthRenunganID)),
			updateJSON:   `{"title":"", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required Title",
		},
		{
			id:           strconv.Itoa(int(AuthRenunganID)),
			updateJSON:   `{"title":"Awesome title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required Content",
		},
		{
			id:           strconv.Itoa(int(AuthRenunganID)),
			updateJSON:   `{"title":"This is another title", "content": "This is the updated content"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthRenunganID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/renungan", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateRenungan)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.Bytes()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteRenungan(t *testing.T) {

	var RenunganUserEmail, RenunganUserPassword string
	var RenunganUserID uint32
	var AuthRenunganID uint64

	err := refreshUserAndRenunganTable()
	if err != nil {
		log.Fatal(err)
	}
	users, renunganlist, err := seedUsersAndRenungan()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		RenunganUserEmail = user.Email
		RenunganUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(RenunganUserEmail, RenunganUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second post
	for _, renungan := range renunganlist {
		if renungan.ID == 1 {
			continue
		}
		AuthRenunganID = renungan.ID
		RenunganUserID = renungan.AuthorID
	}
	postSample := []struct {
		id           string
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthRenunganID)),
			author_id:    RenunganUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthRenunganID)),
			author_id:    RenunganUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthRenunganID)),
			author_id:    RenunganUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			author_id:    1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range postSample {

		req, _ := http.NewRequest("GET", "/renungan", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteRenungan)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.Bytes()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
