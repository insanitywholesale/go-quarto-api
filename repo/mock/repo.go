package mock

import (
	"fmt"
	"gitlab.com/insanitywholesale/adise1941/models"
)

type MockDB struct {
	Users   []*models.User
	UserIds []*models.UserId
	Games   []*models.Game
}

var mymockdb *MockDB

func NewMockDB() (*MockDB, error) {
	mymockdb = &MockDB{
		Users:   []*models.User{},
		UserIds: []*models.UserId{},
		Games:   []*models.Game{},
	}
	return mymockdb, nil
}

func (m *MockDB) AddUser(u *models.User) error {
	m.Users = append(m.Users, u)
	return nil
}

func (m *MockDB) AddUserId(uid *models.UserId) error {
	m.UserIds = append(m.UserIds, uid)
	return nil
}

func (m *MockDB) GetUserIdFromUserId(userid string) (*models.UserId, error) {
	for _, u := range m.UserIds {
		if u.UserId == userid {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user with id", userid, "not found")
}

func (m *MockDB) GetUserIdFromUserName(username string) (*models.UserId, error) {
	for _, u := range m.UserIds {
		if u.UserName == username {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user with name", username, "not found")
}

func (m *MockDB) AddGame(g *models.Game) error {
	m.Games = append(m.Games, g)
	return nil
}

func (m *MockDB) GetGame(gameid string) (*models.Game, error) {
	for _, g := range m.Games {
		if g.GameId == gameid {
			return g, nil
		}
	}
	return nil, fmt.Errorf("game with id", gameid, "not found")
}

func (m *MockDB) ChangeGame(gameNew *models.Game, _ *models.GameMove) error {
	for _, gameOld := range m.Games {
		if gameOld.GameId == gameNew.GameId {
			gameOld = gameNew
			return nil
		}
	}
	return fmt.Errorf("game with id", gameNew.GameId, "not found")
}

func (m *MockDB) GetAllGames() ([]*models.Game, error) {
	return m.Games, nil
}

func (m *MockDB) InviteUser(userid string, gameid string) error {
	u, err := m.GetUserIdFromUserId(userid)
	if err != nil {
		return err
	}
	g, err := m.GetGame(gameid)
	if err != nil {
		return err
	}
	g.InvitedPlayers = append(g.InvitedPlayers, u)
	fmt.Println("inviting:", u.UserName, u.UserId)
	return nil
}

func (m *MockDB) JoinUser(userid string, gameid string) error {
	u, err := m.GetUserIdFromUserId(userid)
	if err != nil {
		return err
	}
	g, err := m.GetGame(gameid)
	if err != nil {
		return err
	}
	for _, ip := range g.InvitedPlayers {
		fmt.Println(ip.UserName, ip.UserId)
		if cap(g.ActivePlayers) == models.MaxPlayers {
			return fmt.Errorf("couldn't join because game is full")
		} else if cap(g.ActivePlayers) > models.MaxPlayers {
			return fmt.Errorf("I honestly don't know how this happened")
		} else if u.UserId == ip.UserId {
			g.ActivePlayers = append(g.ActivePlayers, u)
			//g.InvitedPlayers = g.InvitedPlayers[:len(g.InvitedPlayers)-1]
			return nil
		}
	}
	return fmt.Errorf("player with id", userid, "is not invited to game with id", gameid)
}
