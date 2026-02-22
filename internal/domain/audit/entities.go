package audit

import (
	"encoding/json"
	"time"

	userDomain "erp-api/internal/domain/user"
	"erp-api/internal/utils/dbtypes"
	"gorm.io/gorm"
)

type AuditModule string

type AuditAction string

const (
	AuditModuleAuth     AuditModule = "auth"
	AuditModuleTenant   AuditModule = "tenant"
	AuditModuleUser     AuditModule = "user"
	AuditModuleClient   AuditModule = "client"
	AuditModuleProduct  AuditModule = "product"
	AuditModuleQuote    AuditModule = "quote"
	AuditModuleSettings AuditModule = "settings"
	AuditModuleReports  AuditModule = "reports"
)

const (
	AuditActionCreate AuditAction = "create"
	AuditActionUpdate AuditAction = "update"
	AuditActionDelete AuditAction = "delete"
	AuditActionLogin  AuditAction = "login"
	AuditActionLogout AuditAction = "logout"
)

type Audit struct {
	ID dbtypes.UUID `json:"id" gorm:"primaryKey"`

	TenantID dbtypes.UUID `json:"tenant_id" gorm:"index;not null"`

	UserID *dbtypes.UUID    `json:"user_id,omitempty" gorm:"index"`
	User   *userDomain.User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`

	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty" gorm:"index"`
	UserRole  string `json:"user_role,omitempty" gorm:"index"`

	Module AuditModule `json:"module" gorm:"index;not null"`
	Action AuditAction `json:"action" gorm:"index;not null"`

	ObjectID   string `json:"object_id,omitempty" gorm:"index"`
	ObjectName string `json:"object_name,omitempty" gorm:"index"`

	Method     string `json:"method,omitempty" gorm:"index"`
	Path       string `json:"path,omitempty" gorm:"index"`
	IP         string `json:"ip,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
	StatusCode int    `json:"status_code,omitempty" gorm:"index"`

	Payload json.RawMessage `json:"payload,omitempty" gorm:"type:json"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (a *Audit) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = dbtypes.NewUUID()
	}
	return nil
}

func (Audit) TableName() string { return "audits" }
