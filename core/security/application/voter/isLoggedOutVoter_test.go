package voter

import (
	"context"
	"testing"

	"flamingo.me/flamingo/core/security/application/role/mocks"
	"flamingo.me/flamingo/core/security/domain"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/suite"
)

type (
	IsLoggedOutVoterTestSuite struct {
		suite.Suite

		voter       *IsLoggedOutVoter
		roleService *mocks.Service

		context context.Context
		session *sessions.Session
	}
)

func TestIsLoggedOutVoterTestSuite(t *testing.T) {
	suite.Run(t, &IsLoggedOutVoterTestSuite{})
}

func (t *IsLoggedOutVoterTestSuite) SetupSuite() {
	t.context = context.Background()
	t.session = sessions.NewSession(nil, "-")
}

func (t *IsLoggedOutVoterTestSuite) SetupTest() {
	t.roleService = &mocks.Service{}
	t.voter = &IsLoggedOutVoter{}
	t.voter.Inject(t.roleService)
}

func (t *IsLoggedOutVoterTestSuite) TearDownTest() {
	t.roleService.AssertExpectations(t.T())
	t.roleService = nil
	t.voter = nil
}

func (t *IsLoggedOutVoterTestSuite) TestVote_AccessAbstained() {
	t.Equal(AccessAbstained, t.voter.Vote(t.context, t.session, "SomePermission", nil))
}

func (t *IsLoggedOutVoterTestSuite) TestVote_AccessGranted() {
	t.roleService.On("All", t.context, t.session).Return([]domain.Role{
		domain.RoleAnonymous,
	}).Once()
	t.Equal(AccessGranted, t.voter.Vote(t.context, t.session, domain.RoleAnonymous.Permission(), nil))
}

func (t *IsLoggedOutVoterTestSuite) TestVote_AccessDenied() {
	t.roleService.On("All", t.context, t.session).Return([]domain.Role{
		domain.RoleUser,
	}).Once()
	t.Equal(AccessDenied, t.voter.Vote(t.context, t.session, domain.RoleAnonymous.Permission(), nil))
}
