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
}

// Esta función es llamada desde main.go como goroutine
func StartDailyReminderJob(env *schemas.EnvSettings) {
	loc, err := time.LoadLocation("America/Lima")
	if err != nil {
		log.Fatalf("Error cargando zona horaria: %v", err)
	}

	c := cron.New(cron.WithLocation(loc))

	// Ejecutar todos los días a las 6:30 AM hora Lima
	_, err = c.AddFunc("30 6 * * *", func() {
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

	err := env.DB.Table("reservations r").
		Select(`u.email as user_email, u.name as user_name, r.name as session_name, r.reservation_time as session_time`).
		Joins("JOIN users u ON u.id = r.user_id").
		Where("DATE(r.reservation_time) = ?", today).
		Where("r.state = ?", "ONGOING").
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
		body := fmt.Sprintf(`Hola %s,

Este es un recordatorio de tu sesión de hoy:

🧘 Sesión: %s
🕘 Hora: %s

Gracias por ser parte de ZenCat 🌿`,
			r.UserName,
			r.SessionName,
			r.SessionTime.Format("15:04"),
		)

		err := utils.SendEmail(env, r.UserEmail, "Recordatorio de tu sesión en ZenCat", body)
		if err != nil {
			log.Printf("Error enviando correo a %s: %v", r.UserEmail, err)
		} else {
			log.Printf("Correo enviado a %s", r.UserEmail)
		}
	}
}
