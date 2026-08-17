package service

import (
	"testing"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/util"
	"gorm.io/gorm"
)

func newAuthService(t *testing.T, db *gorm.DB) *AuthService {
	t.Helper()
	userRepo := repository.NewUserRepository(db)
	audit := NewAuditService(repository.NewAuditRepository(db), util.Log)
	return NewAuthService(userRepo, audit, "test-secret", 24, util.Log)
}

func TestRegisterAndLogin(t *testing.T) {
	db := newTestDB(t)
	svc := newAuthService(t, db)

	user, err := svc.Register(&dto.RegisterReq{Username: "dr_zhang", Password: "secret123", Name: "张医生"}, "127.0.0.1", "rid")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if user.Role != constants.UserRoleDoctor {
		t.Errorf("role = %s, want doctor", user.Role)
	}

	resp, err := svc.Login(&dto.LoginReq{Username: "dr_zhang", Password: "secret123"}, "127.0.0.1", "rid")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("token should not be empty")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	db := newTestDB(t)
	svc := newAuthService(t, db)
	if _, err := svc.Register(&dto.RegisterReq{Username: "dr_li", Password: "secret123", Name: "李医生"}, "127.0.0.1", "rid"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if _, err := svc.Login(&dto.LoginReq{Username: "dr_li", Password: "wrong-pass"}, "127.0.0.1", "rid"); err == nil {
		t.Fatal("expected wrong password error")
	}
}

func TestLoginDisabledUser(t *testing.T) {
	db := newTestDB(t)
	svc := newAuthService(t, db)
	user, err := svc.Register(&dto.RegisterReq{Username: "dr_wang", Password: "secret123", Name: "王医生"}, "127.0.0.1", "rid")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	user.Status = constants.UserStatusDisabled
	if err := repository.NewUserRepository(db).Update(user); err != nil {
		t.Fatalf("disable failed: %v", err)
	}
	if _, err := svc.Login(&dto.LoginReq{Username: "dr_wang", Password: "secret123"}, "127.0.0.1", "rid"); err == nil {
		t.Fatal("expected disabled user error")
	}
}

