package business_unit_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nuts-foundation/go-did/did"
	"github.com/openebl/openebl/pkg/bu_server/business_unit"
	"github.com/openebl/openebl/pkg/bu_server/model"
	"github.com/openebl/openebl/pkg/bu_server/storage"
	mock_business_unit "github.com/openebl/openebl/test/mock/bu_server/business_unit"
	mock_storage "github.com/openebl/openebl/test/mock/bu_server/storage"
	"github.com/stretchr/testify/suite"
)

type BusinessUnitManagerTestSuite struct {
	suite.Suite
	ctx       context.Context
	ctrl      *gomock.Controller
	storage   *mock_business_unit.MockBusinessUnitStorage
	tx        *mock_storage.MockTx
	buManager business_unit.BusinessUnitManager
}

func TestBusinessUnitManager(t *testing.T) {
	suite.Run(t, new(BusinessUnitManagerTestSuite))
}

func (s *BusinessUnitManagerTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.storage = mock_business_unit.NewMockBusinessUnitStorage(s.ctrl)
	s.tx = mock_storage.NewMockTx(s.ctrl)
	s.buManager = business_unit.NewBusinessUnitManager(s.storage)
}

func (s *BusinessUnitManagerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *BusinessUnitManagerTestSuite) TestCreateBusinessUnit() {
	ts := time.Now().Unix()

	request := business_unit.CreateBusinessUnitRequest{
		Requester:     "requester",
		ApplicationID: "application-id",
		Name:          "name",
		Addresses:     []string{"address"},
		Emails:        []string{"email"},
		Status:        model.BusinessUnitStatusActive,
	}

	expectedBusinessUnit := model.BusinessUnit{
		Version:       1,
		ApplicationID: request.ApplicationID,
		Status:        request.Status,
		Name:          request.Name,
		Addresses:     request.Addresses,
		Emails:        request.Emails,
		CreatedAt:     ts,
		CreatedBy:     request.Requester,
		UpdatedAt:     ts,
		UpdatedBy:     request.Requester,
	}

	gomock.InOrder(
		s.storage.EXPECT().CreateTx(gomock.Any(), gomock.Len(2)).Return(s.tx, nil),
		s.storage.EXPECT().StoreBusinessUnit(gomock.Any(), s.tx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx storage.Tx, bu model.BusinessUnit) error {
				expectedBusinessUnit.ID = bu.ID
				s.Assert().Equal(expectedBusinessUnit, bu)
				return nil
			},
		),
		s.tx.EXPECT().Commit(gomock.Any()).Return(nil),
		s.tx.EXPECT().Rollback(gomock.Any()).Return(nil),
	)

	newBu, err := s.buManager.CreateBusinessUnit(s.ctx, ts, request)
	s.NoError(err)
	s.Assert().Equal(expectedBusinessUnit, newBu)
}

func (s *BusinessUnitManagerTestSuite) TestUpdateBusinessUnit() {
	ts := time.Now().Unix()

	request := business_unit.UpdateBusinessUnitRequest{
		Requester:     "requester",
		ApplicationID: "application-id",
		ID:            did.MustParseDID("did:openebl:u0e2345"),
		Name:          "name",
		Addresses:     []string{"address"},
		Emails:        []string{"email"},
	}

	oldBusinessUnit := model.BusinessUnit{
		ID:            request.ID,
		Version:       1,
		ApplicationID: request.ApplicationID,
		Status:        model.BusinessUnitStatusActive,
		Name:          "old-name",
		Addresses:     []string{"old-address"},
		Emails:        []string{"old-email"},
		CreatedAt:     ts - 100,
		CreatedBy:     "old-requester",
		UpdatedAt:     ts - 100,
		UpdatedBy:     "old-requester",
	}

	expectedBusinessUnit := model.BusinessUnit{
		ID:            request.ID,
		Version:       2,
		ApplicationID: request.ApplicationID,
		Status:        model.BusinessUnitStatusActive,
		Name:          request.Name,
		Addresses:     request.Addresses,
		Emails:        request.Emails,
		CreatedAt:     ts - 100,
		CreatedBy:     "old-requester",
		UpdatedAt:     ts,
		UpdatedBy:     request.Requester,
	}

	gomock.InOrder(
		s.storage.EXPECT().CreateTx(gomock.Any(), gomock.Len(2)).Return(s.tx, nil),
		s.storage.EXPECT().ListBusinessUnits(
			gomock.Any(),
			s.tx,
			business_unit.ListBusinessUnitsRequest{
				Limit:           1,
				ApplicationID:   request.ApplicationID,
				BusinessUnitIDs: []string{request.ID.String()},
			},
		).Return(business_unit.ListBusinessUnitsResult{
			Total: 1,
			Records: []business_unit.ListBusinessUnitsRecord{
				{
					BusinessUnit: oldBusinessUnit,
				},
			},
		}, nil),
		s.storage.EXPECT().StoreBusinessUnit(gomock.Any(), s.tx, expectedBusinessUnit).Return(nil),
		s.tx.EXPECT().Commit(gomock.Any()).Return(nil),
		s.tx.EXPECT().Rollback(gomock.Any()).Return(nil),
	)

	newBu, err := s.buManager.UpdateBusinessUnit(s.ctx, ts, request)
	s.NoError(err)
	s.Assert().Equal(expectedBusinessUnit, newBu)
}

func (s *BusinessUnitManagerTestSuite) TestListBusinessUnits() {
	request := business_unit.ListBusinessUnitsRequest{
		Offset:          1,
		Limit:           10,
		ApplicationID:   "application-id",
		BusinessUnitIDs: []string{"did:openebl:u0e2345"},
	}

	expectedBusinessUnit := model.BusinessUnit{
		ID:            did.MustParseDID("did:openebl:u0e2345"),
		Version:       1,
		ApplicationID: request.ApplicationID,
		Status:        model.BusinessUnitStatusActive,
		Name:          "name",
		Addresses:     []string{"address"},
		Emails:        []string{"email"},
		CreatedAt:     12345,
		CreatedBy:     "requester",
		UpdatedAt:     12345,
		UpdatedBy:     "requester",
	}

	listResult := business_unit.ListBusinessUnitsResult{
		Total: 1,
		Records: []business_unit.ListBusinessUnitsRecord{
			{
				BusinessUnit: expectedBusinessUnit,
			},
		},
	}

	gomock.InOrder(
		s.storage.EXPECT().CreateTx(gomock.Any(), gomock.Len(0)).Return(s.tx, nil),
		s.storage.EXPECT().ListBusinessUnits(
			gomock.Any(),
			s.tx,
			request,
		).Return(listResult, nil),
		s.tx.EXPECT().Rollback(gomock.Any()).Return(nil),
	)

	result, err := s.buManager.ListBusinessUnits(s.ctx, request)
	s.NoError(err)
	s.Assert().Equal(listResult, result)
}

func (s *BusinessUnitManagerTestSuite) TestSetBusinessUnitStatus() {
	ts := time.Now().Unix()

	request := business_unit.SetBusinessUnitStatusRequest{
		Requester:     "requester",
		ApplicationID: "application-id",
		ID:            did.MustParseDID("did:openebl:u0e2345"),
		Status:        model.BusinessUnitStatusInactive,
	}

	oldBusinessUnit := model.BusinessUnit{
		ID:            request.ID,
		Version:       1,
		ApplicationID: "application-id",
		Status:        model.BusinessUnitStatusActive,
		Name:          "name",
		Addresses:     []string{"address"},
		Emails:        []string{"email"},
		CreatedAt:     ts - 100,
		CreatedBy:     "old-requester",
		UpdatedAt:     ts - 100,
		UpdatedBy:     "old-requester",
	}

	expectedBusinessUnit := model.BusinessUnit{
		ID:            request.ID,
		Version:       2,
		ApplicationID: "application-id",
		Status:        model.BusinessUnitStatusInactive,
		Name:          "name",
		Addresses:     []string{"address"},
		Emails:        []string{"email"},
		CreatedAt:     ts - 100,
		CreatedBy:     "old-requester",
		UpdatedAt:     ts,
		UpdatedBy:     request.Requester,
	}

	gomock.InOrder(
		s.storage.EXPECT().CreateTx(gomock.Any(), gomock.Len(2)).Return(s.tx, nil),
		s.storage.EXPECT().ListBusinessUnits(
			gomock.Any(),
			s.tx,
			business_unit.ListBusinessUnitsRequest{
				Limit:           1,
				ApplicationID:   request.ApplicationID,
				BusinessUnitIDs: []string{request.ID.String()},
			},
		).Return(business_unit.ListBusinessUnitsResult{
			Total: 1,
			Records: []business_unit.ListBusinessUnitsRecord{
				{
					BusinessUnit: oldBusinessUnit,
				},
			},
		}, nil),
		s.storage.EXPECT().StoreBusinessUnit(gomock.Any(), s.tx, expectedBusinessUnit).Return(nil),
		s.tx.EXPECT().Commit(gomock.Any()).Return(nil),
		s.tx.EXPECT().Rollback(gomock.Any()).Return(nil),
	)

	newBu, err := s.buManager.SetStatus(s.ctx, ts, request)
	s.NoError(err)
	s.Assert().Equal(expectedBusinessUnit, newBu)
}

func (s *BusinessUnitManagerTestSuite) TestAddAuthentication() {
	ts := time.Now().Unix()

	request := business_unit.AddAuthenticationRequest{
		Requester:      "requester",
		ApplicationID:  "application-id",
		BusinessUnitID: did.MustParseDID("did:openebl:u0e2345"),
		PrivateKey:     "FAKE PEM PRIVATE KEY",
		Certificate:    "FAKE PEM CERT",
	}

	expectedAuthentication := model.BusinessUnitAuthentication{
		Version:      1,
		BusinessUnit: request.BusinessUnitID,
		Status:       model.BusinessUnitAuthenticationStatusActive,
		CreatedAt:    ts,
		CreatedBy:    request.Requester,
		PrivateKey:   request.PrivateKey,
		Certificate:  request.Certificate,
	}

	gomock.InOrder(
		s.storage.EXPECT().CreateTx(gomock.Any(), gomock.Len(2)).Return(s.tx, nil),
		s.storage.EXPECT().StoreAuthentication(gomock.Any(), s.tx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx storage.Tx, auth model.BusinessUnitAuthentication) error {
				expectedAuthentication.ID = auth.ID
				s.Assert().Equal(expectedAuthentication, auth)
				return nil
			},
		),
		s.tx.EXPECT().Commit(gomock.Any()).Return(nil),
		s.tx.EXPECT().Rollback(gomock.Any()).Return(nil),
	)

	newAuthentication, err := s.buManager.AddAuthentication(s.ctx, ts, request)
	s.NoError(err)
	s.Assert().Empty(newAuthentication.PrivateKey)
	newAuthentication.PrivateKey = expectedAuthentication.PrivateKey
	s.Assert().Equal(expectedAuthentication, newAuthentication)
}

func (s *BusinessUnitManagerTestSuite) TestRevokeAuthentication() {
	ts := time.Now().Unix()

	request := business_unit.RevokeAuthenticationRequest{
		Requester:        "requester",
		ApplicationID:    "application-id",
		BusinessUnitID:   did.MustParseDID("did:openebl:u0e2345"),
		AuthenticationID: "authentication-id",
	}

	oldAuthentication := model.BusinessUnitAuthentication{
		ID:           "authentication-id",
		Version:      1,
		BusinessUnit: request.BusinessUnitID,
		Status:       model.BusinessUnitAuthenticationStatusActive,
		CreatedAt:    ts - 100,
		CreatedBy:    "old-requester",
		PrivateKey:   "FAKE PEM PRIVATE KEY",
		Certificate:  "FAKE PEM CERT",
		RevokedAt:    0,
	}

	expectedAuthentication := model.BusinessUnitAuthentication{
		ID:           "authentication-id",
		Version:      2,
		BusinessUnit: request.BusinessUnitID,
		Status:       model.BusinessUnitAuthenticationStatusRevoked,
		CreatedAt:    ts - 100,
		CreatedBy:    "old-requester",
		PrivateKey:   "FAKE PEM PRIVATE KEY",
		Certificate:  "FAKE PEM CERT",
		RevokedAt:    ts,
		RevokedBy:    request.Requester,
	}

	gomock.InOrder(
		s.storage.EXPECT().CreateTx(gomock.Any(), gomock.Len(2)).Return(s.tx, nil),
		s.storage.EXPECT().ListAuthentication(
			gomock.Any(),
			s.tx,
			business_unit.ListAuthenticationRequest{
				Limit:             1,
				ApplicationID:     request.ApplicationID,
				BusinessUnitID:    request.BusinessUnitID.String(),
				AuthenticationIDs: []string{request.AuthenticationID},
			},
		).Return(business_unit.ListAuthenticationResult{
			Total:   1,
			Records: []model.BusinessUnitAuthentication{oldAuthentication},
		}, nil),
		s.storage.EXPECT().StoreAuthentication(gomock.Any(), s.tx, expectedAuthentication).Return(nil),
		s.tx.EXPECT().Commit(gomock.Any()).Return(nil),
		s.tx.EXPECT().Rollback(gomock.Any()).Return(nil),
	)

	newAuthentication, err := s.buManager.RevokeAuthentication(s.ctx, ts, request)
	s.NoError(err)
	s.Assert().Empty(newAuthentication.PrivateKey)
	newAuthentication.PrivateKey = expectedAuthentication.PrivateKey
	s.Assert().Equal(expectedAuthentication, newAuthentication)
}

func (s *BusinessUnitManagerTestSuite) TestListAuthentication() {
	request := business_unit.ListAuthenticationRequest{
		Offset:            1,
		Limit:             10,
		ApplicationID:     "application-id",
		BusinessUnitID:    "did:openebl:u0e2345",
		AuthenticationIDs: []string{"authentication-id"},
	}

	expectedAuthentication := model.BusinessUnitAuthentication{
		ID:           "authentication-id",
		Version:      1,
		BusinessUnit: did.MustParseDID(request.BusinessUnitID),
		Status:       model.BusinessUnitAuthenticationStatusActive,
		CreatedAt:    12345,
		CreatedBy:    "requester",
		PrivateKey:   "FAKE PEM PRIVATE KEY",
		Certificate:  "FAKE PEM CERT",
		RevokedAt:    0,
	}

	listResult := business_unit.ListAuthenticationResult{
		Total:   1,
		Records: []model.BusinessUnitAuthentication{expectedAuthentication},
	}

	gomock.InOrder(
		s.storage.EXPECT().CreateTx(gomock.Any(), gomock.Len(0)).Return(s.tx, nil),
		s.storage.EXPECT().ListAuthentication(
			gomock.Any(),
			s.tx,
			request,
		).Return(listResult, nil),
		s.tx.EXPECT().Rollback(gomock.Any()).Return(nil),
	)

	result, err := s.buManager.ListAuthentication(s.ctx, request)
	s.NoError(err)
	s.Require().NotEmpty(result.Records)
	s.Assert().Empty(result.Records[0].PrivateKey)
	s.Assert().Equal(listResult, result)
}
