package jobs

import (
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

// MembershipExpirer is a background job that marks all ACTIVE memberships whose
// end_date ya pasó como EXPIRED. Se ejecuta todos los días a las 00:15.
type MembershipExpirer struct {
	cron   *cron.Cron
	logger logging.Logger
	db     *gorm.DB
}

// NewMembershipExpirer crea la instancia y registra el job en el scheduler, pero
// NO lo arranca; para eso hay que llamar Start().
func NewMembershipExpirer(logger logging.Logger, db *gorm.DB) *MembershipExpirer {
	c := cron.New()
	expirer := &MembershipExpirer{cron: c, logger: logger, db: db}

	// "15 0 * * *"  ->  At 00:15 (UTC/local timezone del contenedor) todos los días
	_, err := c.AddFunc("15 0 * * *", expirer.run)
	if err != nil {
		logger.Errorf("MembershipExpirer: error añadiendo cron job: %v", err)
	}

	return expirer
}

// Start inicia el scheduler.
func (m *MembershipExpirer) Start() {
	m.logger.Infoln("MembershipExpirer: cron iniciado (diario a las 00:15)")
	m.cron.Start()
}

// run contiene la lógica de negocio: actualizar estado de las membresías.
func (m *MembershipExpirer) run() {
	now := time.Now()

	res := m.db.Model(&model.Membership{}).
		Where("status = ? AND end_date < ?", model.MembershipStatusActive, now).
		Update("status", model.MembershipStatusExpired)

	if res.Error != nil {
		m.logger.Errorf("MembershipExpirer: fallo al actualizar membresías: %v", res.Error)
		return
	}

	if res.RowsAffected > 0 {
		m.logger.Infof("MembershipExpirer: %d membresías pasaron a EXPIRED", res.RowsAffected)
	}
}
