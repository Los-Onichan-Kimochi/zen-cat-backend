package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/utils"
)

// Estructura para escanear resultados desde GORM
type Reminder struct {
	UserEmail   string
	UserName    string
	SessionName string
	SessionTime time.Time
	SessionLink string
	LocalName   string
	StreetName  string
	BuildingNum string
	District    string
}

// Esta función es llamada desde main.go como goroutine
func StartDailyReminderJob(env *schemas.EnvSettings) {
	c := cron.New()

	// Ejecutar todos los días a las 6:00 AM hora Lima
	_, err := c.AddFunc("0 7 * * *", func() {
		log.Println("📬 [CRON] Ejecutando recordatorios diarios de sesiones...")
		SendRemindersForToday(env)
	})
	if err != nil {
		log.Fatalf("Error al programar cronjob: %v", err)
	}

	c.Start()
	log.Println("[CRON] Recordatorio diario programado correctamente (6:30 AM America/Lima)")
}

// Consulta las reservas y envía emails
func SendRemindersForToday(env *schemas.EnvSettings) {
	today := time.Now().In(time.FixedZone("UTC-5", -5*3600)).Format("2006-01-02")

	var reminders []Reminder

	err := env.DB.Table("astro_cat_reservation r").
		Select(`
		u.email AS user_email,
		u.name AS user_name,
		r.name AS session_name,
		r.reservation_time AS session_time,
		s.session_link,
		l.local_name,
		l.street_name,
		l.building_number,
		l.district
	`).
		Joins("JOIN astro_cat_user u ON u.id = r.user_id").
		Joins("JOIN astro_cat_session s ON s.id = r.session_id").
		Joins("LEFT JOIN astro_cat_local l ON l.id = s.local_id").
		Where("DATE(r.reservation_time) = ? AND r.state = ?", today, "CONFIRMED").
		Scan(&reminders).Error
	if err != nil {
		log.Printf("Error obteniendo reservas: %v", err)
		return
	}

	if len(reminders) == 0 {
		log.Println("No hay reservas programadas para hoy.")
		return
	}

	for _, r := range reminders {
		var tipoInfo string
		var linkInfo string

		if r.SessionLink != "" {
			tipoInfo = "🌐 Tipo: Virtual"
			linkInfo = fmt.Sprintf("\n\n🔗 Enlace de acceso: %s", r.SessionLink)
		} else {
			tipoInfo = fmt.Sprintf(
				"📍 Lugar: %s, %s %s, %s",
				r.LocalName,
				r.StreetName,
				r.BuildingNum,
				r.District,
			)
		}

		body := fmt.Sprintf(`Hola %s,

Este es un recordatorio de tu sesión de hoy:

🧘 Sesión: %s
🕘 Hora: %s
%s%s

Gracias por ser parte de ZenCat 🌿`,
			r.UserName,
			r.SessionName,
			r.SessionTime.Format("15:04"),
			tipoInfo,
			linkInfo,
		)

		err := utils.SendEmail(env, r.UserEmail, "Recordatorio de tu sesión en ZenCat", body)
		if err != nil {
			log.Printf("❌ Error enviando correo a %s: %v", r.UserEmail, err)
		} else {
			log.Printf("✅ Correo enviado a %s", r.UserEmail)
		}
	}
}
