package v1

import (
	"context"
	"strings"
	"time"

	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-util/xid"
	"github.com/go-pantheon/fabrica-util/xrand"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	"github.com/go-pantheon/lares/app/account/internal/pkg/security"
	v1 "github.com/go-pantheon/lares/gen/api/server/account/interface/account/v1"
	gatev1 "github.com/go-pantheon/lares/gen/api/server/gate/intra/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func (s *AccountInterface) genLoginInfo(_ context.Context, account *domain.Account, isRegister bool, state, color string) (info *v1.LoginInfo, err error) {
	info = &v1.LoginInfo{
		Register: isRegister,
		State:    state,
	}

	if info.AccountId, err = xid.EncodeID(account.Id); err != nil {
		return nil, err
	}

	if info.Session, info.SessionTimeout, err = genSession(account); err != nil {
		return nil, err
	}

	if info.Token, info.TokenTimeout, err = genToken(account, color); err != nil {
		return nil, err
	}

	return
}

func genToken(account *domain.Account, color string) (token string, expiredTime int64, err error) {
	rand, err := xrand.RandAlphaNumString(RandLength)
	if err != nil {
		return "", 0, errors.Wrapf(err, "rand string failed")
	}

	expiredTime = time.Now().Add(TokenExpiredDuration).Unix()

	auth := &gatev1.AuthToken{
		AccountId: account.Id,
		Timeout:   expiredTime,
		Rand:      rand,
		ServerId:  0,
		Location:  int32(profile.Zone()),
		Color:     getAccountColor(account, color),
		Status:    gatev1.OnlineStatus_ONLINE_STATUS_GATE,
	}

	bytes, err := proto.Marshal(auth)
	if err != nil {
		return "", 0, errors.Wrapf(err, "token encode failed")
	}

	token, err = security.EncryptToken(bytes)
	if err != nil {
		return "", 0, errors.Wrapf(err, "token encrypt failed")
	}

	return token, expiredTime, nil
}

func genSession(account *domain.Account) (session string, expiredTime int64, err error) {
	rand, err := xrand.RandAlphaNumString(RandLength)
	if err != nil {
		return "", 0, errors.Wrapf(err, "rand string failed")
	}

	expiredTime = time.Now().Add(SessionExpiredDuration).Unix()

	p := &v1.Session{
		AccountId: account.Id,
		Timeout:   expiredTime,
		Key:       rand,
	}

	bytes, err := proto.Marshal(p)
	if err != nil {
		return "", 0, errors.Wrapf(err, "session encode failed")
	}

	session, err = security.EncryptSession(bytes)
	if err != nil {
		return "", 0, errors.Wrapf(err, "session encrypt failed")
	}

	return session, expiredTime, nil
}

func unmarshalSession(secret string) (session *v1.Session, err error) {
	org, err := security.DecryptSession(secret)
	if err != nil {
		return nil, errors.Wrapf(err, "session decrypt failed")
	}

	session = &v1.Session{}
	if err = proto.Unmarshal(org, session); err != nil {
		return nil, errors.Wrapf(err, "session decode failed")
	}

	return session, nil
}

func getAccountColor(account *domain.Account, color string) string {
	color = strings.TrimSpace(color)
	if len(color) > 0 {
		return color
	}

	return account.DefaultColor
}
