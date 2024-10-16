package connector

import (
	"errors"
	"fmt"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

const (
	DummyTest                  = "dummy-test"
	DummyTestUser              = "dummy-user"
	DummyTestUserErrorOnCreate = "dummy-user-create-error"
	DummyTestUserErrorOnGrant  = "dummy-user-grant-error"
	DummyTestUserErrorOnRemove = "dummy-user-remove-error"
	InstanceDummyTestError     = "instance-dummy-test-error"
)

var (
	ErrCreateUser      = errors.New("error creating user")
	ErrGrantConnect    = errors.New("error granting connect")
	ErrorRemoveUser    = errors.New("error revoking permissions and removing user")
	ErrorCreatingRoles = errors.New("error creating roles")
)

type DummyTestConnector struct {
	ConnectionData dto.ConnectionInputDTO
}

func newDummyTestConnector(connectionData dto.ConnectionInputDTO) *DummyTestConnector {
	return &DummyTestConnector{ConnectionData: connectionData}
}

func (d *DummyTestConnector) TestConnection() error {
	if d.ConnectionData.Instance == InstanceDummyTestError {
		return fmt.Errorf("error testing connection with %s", d.ConnectionData.Instance)
	}
	return nil
}

func (d *DummyTestConnector) ListDatabases() ([]*Database, error) {
	if d.ConnectionData.Instance == InstanceDummyTestError {
		return nil, fmt.Errorf("error listing databases from %s", d.ConnectionData.Instance)
	}
	databases := []*Database{
		{
			Name:        "dummy-db-1",
			CurrentSize: "1GB",
		},
		{
			Name:        "dummy-db-2",
			CurrentSize: "2GB",
		},
		{
			Name:        "dummy-db-5",
			CurrentSize: "5GB",
		},
	}
	return databases, nil
}

func (d *DummyTestConnector) Driver() string {
	return "dummy"
}

func (d *DummyTestConnector) URL() string {
	return fmt.Sprintf("dummy-test://%s:%s@%s:%s/%s?sslmode=disable",
		d.ConnectionData.User,
		d.ConnectionData.Password,
		d.ConnectionData.Host,
		d.ConnectionData.Port,
		d.Database(),
	)
}

func (d *DummyTestConnector) Database() string {
	databaseName := d.ConnectionData.Database
	if databaseName == "" {
		databaseName = d.DefaultDatabase()
	}
	return databaseName
}

func (d *DummyTestConnector) DefaultDatabase() string {
	return "dummy-test-db"
}

func (d *DummyTestConnector) CreateRoles(_ []*DatabaseRole) error {
	if d.ConnectionData.Instance == InstanceDummyTestError {
		return fmt.Errorf("%w: Instance(%s)", ErrorCreatingRoles, d.ConnectionData.Instance)
	}
	return nil
}

func (d *DummyTestConnector) SetupGrantsToRoles() error {
	if d.ConnectionData.Instance == InstanceDummyTestError {
		return fmt.Errorf("%w: Instance(%s)", ErrGrantConnect, d.ConnectionData.Instance)
	}
	return nil
}

func (d *DummyTestConnector) UserExists(username string) (bool, error) {
	if d.ConnectionData.Instance == InstanceDummyTestError {
		return false, fmt.Errorf("error checking user existence in %s", d.ConnectionData.Instance)
	}
	if username == DummyTestUser || username == DummyTestUserErrorOnCreate {
		return false, nil
	}
	return true, nil
}

func (d *DummyTestConnector) CreateUser(user *DatabaseUser) error {
	if d.ConnectionData.Instance == InstanceDummyTestError || user.Username == DummyTestUserErrorOnCreate {
		return fmt.Errorf("%w: Instance(%s) - User(%s)", ErrCreateUser, d.ConnectionData.Instance, user.Username)
	}
	return nil
}

func (d *DummyTestConnector) GrantConnect(username string) error {
	if d.ConnectionData.Instance == InstanceDummyTestError || username == DummyTestUserErrorOnGrant {
		return fmt.Errorf("%w: Instance(%s) - User(%s)", ErrGrantConnect, d.ConnectionData.Instance, username)
	}
	return nil
}

func (d *DummyTestConnector) RevokeUserPrivilegesAndRemove(username string) error {
	if d.ConnectionData.Instance == InstanceDummyTestError || username == DummyTestUserErrorOnRemove {
		return fmt.Errorf("%w: Instance(%s) - User(%s)", ErrorRemoveUser, d.ConnectionData.Instance, username)
	}
	return nil
}
