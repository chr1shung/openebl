package manager_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/openebl/openebl/pkg/bu_server/auth"
	"github.com/openebl/openebl/pkg/bu_server/cert_authority"
	"github.com/openebl/openebl/pkg/bu_server/manager"
	"github.com/openebl/openebl/pkg/bu_server/model"
	"github.com/openebl/openebl/pkg/util"
	mock_auth "github.com/openebl/openebl/test/mock/bu_server/auth"
	mock_cert_authority "github.com/openebl/openebl/test/mock/bu_server/cert_authority"
	"github.com/stretchr/testify/suite"
)

func TestManagerAPIWithDB(t *testing.T) {
	t.Skip("Skipping test for now")
	dbConfig := util.PostgresDatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "xdlai",
		Database: "bu_server_test",
		SSLMode:  "disable",
		PoolSize: 5,
	}

	restConfig := manager.ManagerAPIConfig{
		Database:     dbConfig,
		LocalAddress: "localhost:9100",
	}

	rest, err := manager.NewManagerAPI(restConfig)
	if err != nil {
		t.Fatal(err)
	}
	rest.Run()
}

type ManagerAPITestSuite struct {
	suite.Suite
	ctx       context.Context
	ctrl      *gomock.Controller
	userMgr   *mock_auth.MockUserManager
	appMgr    *mock_auth.MockApplicationManager
	apiKeyMgr *mock_auth.MockAPIKeyAuthenticator
	ca        *mock_cert_authority.MockCertAuthority
}

func TestManagerAPI(t *testing.T) {
	suite.Run(t, new(ManagerAPITestSuite))
}

func (s *ManagerAPITestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.userMgr = mock_auth.NewMockUserManager(s.ctrl)
	s.appMgr = mock_auth.NewMockApplicationManager(s.ctrl)
	s.apiKeyMgr = mock_auth.NewMockAPIKeyAuthenticator(s.ctrl)
	s.ca = mock_cert_authority.NewMockCertAuthority(s.ctrl)
}

func (s *ManagerAPITestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ManagerAPITestSuite) TestLogin() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9201")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	expectedRequest := auth.AuthenticateUserRequest{
		UserID:   "username",
		Password: "password",
	}
	userToken := auth.UserToken{
		Token:  "toooooooooken",
		UserID: "username",
	}
	gomock.InOrder(
		s.userMgr.EXPECT().Authenticate(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(userToken, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodGet, "http://localhost:9201/login", nil)
	request.Header.Add("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=") // username:password
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(userToken.Token, string(body))
}

func (s *ManagerAPITestSuite) TestLoginWithInvalidCredentials() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9202")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	expectedRequest := auth.AuthenticateUserRequest{
		UserID:   "username",
		Password: "password",
	}
	gomock.InOrder(
		s.userMgr.EXPECT().Authenticate(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(auth.UserToken{}, model.ErrUserAuthenticationFail),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodGet, "http://localhost:9202/login", nil)
	request.Header.Add("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=") // username:password
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusUnauthorized, response.StatusCode)
}

func (s *ManagerAPITestSuite) TestCreateUser() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9203")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	expectedRequest := auth.CreateUserRequest{
		RequestUser: "request_user",
		UserID:      "user_id",
		Password:    "password",
		Name:        "name",
		Emails:      []string{"email"},
		Note:        "note",
	}
	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "request_user",
	}
	newUser := auth.User{
		ID: "new_user_id",
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.userMgr.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(newUser, nil),
	)

	createUserRequest := auth.CreateUserRequest{
		UserID:   "user_id",
		Password: "password",
		Name:     "name",
		Emails:   []string{"email"},
		Note:     "note",
	}
	request, _ := http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9203/user", util.StructToJSONReader(createUserRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(newUser), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestListUser() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9204")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "user_id",
	}

	expectedRequest := auth.ListUserRequest{
		Offset: 1,
		Limit:  2,
	}
	listResult := auth.ListUserResult{
		Total: 999,
		Users: []auth.User{
			{
				ID: "user_id_1",
			},
			{
				ID: "user_id_2",
			},
		},
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.userMgr.EXPECT().ListUsers(gomock.Any(), gomock.Eq(expectedRequest)).Return(listResult, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodGet, "http://localhost:9204/user?offset=1&limit=2", nil)
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(listResult), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestGetUser() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9205")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "user_id",
	}

	expectedRequest := auth.ListUserRequest{
		Limit: 1,
		IDs:   []string{"user_id"},
	}
	listResult := auth.ListUserResult{
		Total: 1,
		Users: []auth.User{
			{
				ID: "user_id",
			},
		},
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.userMgr.EXPECT().ListUsers(gomock.Any(), gomock.Eq(expectedRequest)).Return(listResult, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodGet, "http://localhost:9205/user/user_id", nil)
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(listResult.Users[0]), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestUpdateUser() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9206")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "request_user",
	}

	expectedRequest := auth.UpdateUserRequest{
		RequestUser: "request_user",
		UserID:      "user_id",
		Name:        "name",
		Emails:      []string{"email"},
		Note:        "note",
	}
	updatedUser := auth.User{
		ID: "user_id",
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.userMgr.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(updatedUser, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9206/user/user_id", util.StructToJSONReader(expectedRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(updatedUser), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestUpdateUserStatus() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9207")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "request_user",
	}

	expectedRequest := auth.ActivateUserRequest{
		RequestUser: "request_user",
		UserID:      "user_id",
	}
	updatedUser := auth.User{
		ID: "user_id",
	}

	// Test Activate User
	restRequest := map[string]string{
		"status": "active",
	}
	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.userMgr.EXPECT().ActivateUser(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(updatedUser, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9207/user/user_id/status", util.StructToJSONReader(restRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(updatedUser), strings.TrimSpace(string(body)))
	// End of Test Activate User

	// Test Deactivate User
	restRequest = map[string]string{
		"status": "inactive",
	}
	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.userMgr.EXPECT().DeactivateUser(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(updatedUser, nil),
	)

	request, _ = http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9207/user/user_id/status", util.StructToJSONReader(restRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err = http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ = io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(updatedUser), strings.TrimSpace(string(body)))
	// End of Test Deactivate User
}

func (s *ManagerAPITestSuite) TestChangeUserPassword() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9208")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}
	expectedRequest := auth.ChangePasswordRequest{
		UserID:      "user_id",
		OldPassword: "old password",
		Password:    "new password",
	}
	newUser := auth.User{
		ID: "user_id",
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.userMgr.EXPECT().ChangePassword(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(newUser, nil),
	)

	restRequest := map[string]string{
		"old_password": "old password",
		"password":     "new password",
	}
	request, _ := http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9208/user/user_id/change_password", util.StructToJSONReader(restRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
}

func (s *ManagerAPITestSuite) TestResetUserPassword() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9209")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}
	expectedRequest := auth.ResetPasswordRequest{
		RequestUser: "requester_id",
		UserID:      "user_id",
		Password:    "new password",
	}
	newUser := auth.User{
		ID: "user_id",
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.userMgr.EXPECT().ResetPassword(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(newUser, nil),
	)

	restRequest := map[string]string{
		"password": "new password",
	}
	request, _ := http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9209/user/user_id/reset_password", util.StructToJSONReader(restRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
}

func (s *ManagerAPITestSuite) TestCreateApplication() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9210")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}

	expectedRequest := auth.CreateApplicationRequest{
		RequestUser: auth.RequestUser{
			User: "requester_id",
		},
		Name:         "name",
		CompanyName:  "company name",
		Addresses:    []string{"address"},
		Emails:       []string{"email"},
		PhoneNumbers: []string{"phone number"},
	}

	app := auth.Application{
		ID:      "app_id",
		Version: 1,
		Status:  auth.ApplicationStatusActive,
		Name:    "name",
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.appMgr.EXPECT().CreateApplication(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(app, nil),
	)

	restRequest := map[string]any{
		"name":          "name",
		"company_name":  "company name",
		"addresses":     []string{"address"},
		"emails":        []string{"email"},
		"phone_numbers": []string{"phone number"},
	}

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9210/application", util.StructToJSONReader(restRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(app), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestListApplication() {
	restServer, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9211")
	s.Require().NoError(err)
	go func() { restServer.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer restServer.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}

	expectedRequest := auth.ListApplicationRequest{
		Offset: 1,
		Limit:  2,
	}
	listResult := auth.ListApplicationResult{
		Total: 999,
		Applications: []auth.Application{
			{
				ID: "app_id_1",
			},
		},
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.appMgr.EXPECT().ListApplications(gomock.Any(), gomock.Eq(expectedRequest)).Return(listResult, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodGet, "http://localhost:9211/application?offset=1&limit=2", nil)
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(listResult), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestGetApplication() {
	rest, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9212")
	s.Require().NoError(err)
	go func() { rest.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer rest.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}

	expectedRequest := auth.ListApplicationRequest{
		Limit: 1,
		IDs:   []string{"app_id"},
	}

	listResult := auth.ListApplicationResult{
		Total: 1,
		Applications: []auth.Application{
			{
				ID: "app_id",
			},
		},
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.appMgr.EXPECT().ListApplications(gomock.Any(), gomock.Eq(expectedRequest)).Return(listResult, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodGet, "http://localhost:9212/application/app_id", nil)
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(listResult.Applications[0]), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestUpdateApplication() {
	rest, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9213")
	s.Require().NoError(err)
	go func() { rest.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer rest.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}

	expectedRequest := auth.UpdateApplicationRequest{
		CreateApplicationRequest: auth.CreateApplicationRequest{
			RequestUser: auth.RequestUser{
				User: "requester_id",
			},
			Name:         "name",
			CompanyName:  "company name",
			Addresses:    []string{"address"},
			Emails:       []string{"email"},
			PhoneNumbers: []string{"phone number"},
		},
		ID: "app_id",
	}

	app := auth.Application{
		ID: "app_id",
	}

	restRequest := map[string]any{
		"name":          "name",
		"company_name":  "company name",
		"addresses":     []string{"address"},
		"emails":        []string{"email"},
		"phone_numbers": []string{"phone number"},
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.appMgr.EXPECT().UpdateApplication(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(app, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9213/application/app_id", util.StructToJSONReader(restRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()
	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(app), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestUpdateApplicationStatus() {
	rest, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9214")
	s.Require().NoError(err)
	go func() { rest.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer rest.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}

	expectedRequest := auth.ActivateApplicationRequest{
		RequestUser: auth.RequestUser{
			User: "requester_id",
		},
		ApplicationID: "app_id",
	}

	app := auth.Application{
		ID: "app_id",
	}

	// Test Activate Application
	restRequest := map[string]any{
		"status": auth.ApplicationStatusActive,
	}
	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.appMgr.EXPECT().ActivateApplication(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(app, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9214/application/app_id/status", util.StructToJSONReader(restRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()
	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(app), strings.TrimSpace(string(body)))
	// End of Test Activate Application

	// Test Inactivate Application
	restRequest = map[string]any{
		"status": auth.ApplicationStatusInactive,
	}
	app.Status = auth.ApplicationStatusInactive
	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.appMgr.EXPECT().DeactivateApplication(gomock.Any(), gomock.Any(), gomock.Eq(auth.DeactivateApplicationRequest(expectedRequest))).Return(app, nil),
	)

	request, _ = http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9214/application/app_id/status", util.StructToJSONReader(restRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err = http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()
	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ = io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(app), strings.TrimSpace(string(body)))
	// End of Test Inactivate Application
}

func (s *ManagerAPITestSuite) TestCreateAPIKey() {
	rest, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9215")
	s.Require().NoError(err)
	go func() { rest.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer rest.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}

	expectedRequest := auth.CreateAPIKeyRequest{
		RequestUser: auth.RequestUser{
			User: "requester_id",
		},
		ApplicationID: "app_id",
		Scopes:        []auth.APIKeyScope{auth.APIKeyScopeAll},
	}

	apiKey := auth.APIKey{
		ID: "api_key_id",
	}
	apiKeyString := auth.APIKeyString("api_key_id:secret")

	restRequest := map[string]any{
		"scopes": []string{"all"},
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.apiKeyMgr.EXPECT().CreateAPIKey(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(apiKey, apiKeyString, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9215/application/app_id/api_key", util.StructToJSONReader(restRequest))
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()
	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(string(apiKeyString), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestListAPIKey() {
	rest, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9216")
	s.Require().NoError(err)
	go func() { rest.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer rest.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}

	expectedRequest := auth.ListAPIKeysRequest{
		Offset:         1,
		Limit:          2,
		ApplicationIDs: []string{"app_id"},
	}

	listResult := auth.ListAPIKeysResult{
		Total: 999,
		Keys: []auth.ListAPIKeyRecord{
			{
				APIKey: auth.APIKey{
					ID:            "api_key_id_1",
					ApplicationID: "app_id_1",
				},
				Application: auth.Application{
					ID: "app_id_1",
				},
			},
		},
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.apiKeyMgr.EXPECT().ListAPIKeys(gomock.Any(), gomock.Eq(expectedRequest)).Return(listResult, nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodGet, "http://localhost:9216/application/app_id/api_key?offset=1&limit=2", nil)
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()
	s.Assert().Equal(http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(listResult), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestRevokeAPIKey() {
	rest, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9217")
	s.Require().NoError(err)
	go func() { rest.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer rest.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}

	expectedRequest := auth.RevokeAPIKeyRequest{
		RequestUser: auth.RequestUser{
			User: "requester_id",
		},
		ApplicationID: "app_id",
		ID:            "api_key_id",
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.apiKeyMgr.EXPECT().RevokeAPIKey(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(nil),
	)

	request, _ := http.NewRequestWithContext(s.ctx, http.MethodDelete, "http://localhost:9217/application/app_id/api_key/api_key_id", nil)
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)
	defer response.Body.Close()
	s.Assert().Equal(http.StatusOK, response.StatusCode)
}

func (s *ManagerAPITestSuite) TestAddCACertificate() {
	rest, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9218")
	s.Require().NoError(err)
	go func() { rest.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer rest.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}

	request := cert_authority.AddCertificateRequest{
		Cert:       "cert",
		PrivateKey: "private key",
	}

	expectedRequest := request
	expectedRequest.Requester = userToken.UserID

	cert := model.Cert{
		ID:              "cert_id",
		Version:         1,
		Status:          model.CertStatusActive,
		Type:            model.CACert,
		PrivateKey:      "private key",
		Certificate:     "cert",
		CertFingerPrint: "fingerprint",
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.ca.EXPECT().AddCertificate(gomock.Any(), gomock.Any(), expectedRequest).Return(cert, nil),
	)

	httpRequest, err := http.NewRequestWithContext(s.ctx, http.MethodPost, "http://localhost:9218/ca/certificate", util.StructToJSONReader(request))
	s.Require().NoError(err)

	httpRequest.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(httpRequest)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, response.StatusCode)
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON(cert), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestListCACertificates() {
	rest, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9219")
	s.Require().NoError(err)
	go func() { rest.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer rest.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}
	expectedRequest := cert_authority.ListCertificatesRequest{
		Offset: 1,
		Limit:  2,
	}
	cert := model.Cert{
		ID: "cert_id",
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.ca.EXPECT().ListCertificates(gomock.Any(), gomock.Eq(expectedRequest)).Return([]model.Cert{cert}, nil),
	)

	httpRequest, err := http.NewRequestWithContext(s.ctx, http.MethodGet, "http://localhost:9219/ca/certificate?offset=1&limit=2", nil)
	s.Require().NoError(err)

	httpRequest.Header.Add("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(httpRequest)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, response.StatusCode)
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	s.Assert().Equal(util.StructToJSON([]model.Cert{cert}), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestGetCACertificate() {
	rest, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9220")
	s.Require().NoError(err)
	go func() { rest.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer rest.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}

	expectedRequest := cert_authority.ListCertificatesRequest{
		Limit: 1,
		IDs:   []string{"cert_id"},
	}
	cert := model.Cert{
		ID: "cert_id",
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.ca.EXPECT().ListCertificates(gomock.Any(), gomock.Eq(expectedRequest)).Return([]model.Cert{cert}, nil),
	)

	httpRequest, err := http.NewRequestWithContext(s.ctx, http.MethodGet, "http://localhost:9220/ca/certificate/cert_id", nil)
	s.Require().NoError(err)

	httpRequest.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(httpRequest)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	s.Assert().Equal(util.StructToJSON(cert), strings.TrimSpace(string(body)))
}

func (s *ManagerAPITestSuite) TestRevokeCACertificate() {
	rest, err := manager.NewManagerAPIWithControllers(s.userMgr, s.appMgr, s.apiKeyMgr, s.ca, "localhost:9221")
	s.Require().NoError(err)
	go func() { rest.Run() }()
	time.Sleep(200 * time.Millisecond)
	defer rest.Close()

	token := "user token"
	userToken := auth.UserToken{
		Token:  token,
		UserID: "requester_id",
	}
	expectedRequest := cert_authority.RevokeCertificateRequest{
		Requester: "requester_id",
		CertID:    "cert_id",
	}
	cert := model.Cert{
		ID: "cert_id",
	}

	gomock.InOrder(
		s.userMgr.EXPECT().TokenAuthorization(gomock.Any(), gomock.Any(), gomock.Eq(token)).Return(userToken, nil),
		s.ca.EXPECT().RevokeCertificate(gomock.Any(), gomock.Any(), gomock.Eq(expectedRequest)).Return(cert, nil),
	)

	httpRequest, err := http.NewRequestWithContext(s.ctx, http.MethodDelete, "http://localhost:9221/ca/certificate/cert_id", nil)
	s.Require().NoError(err)

	httpRequest.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(httpRequest)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	s.Assert().Equal(util.StructToJSON(cert), strings.TrimSpace(string(body)))
}
