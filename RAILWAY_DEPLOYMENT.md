# Railway Deployment Guide

## Variables de Entorno

### ✅ **Configuración Automática**

El archivo `railway.json` incluye todas las variables de entorno ya configuradas que sí pueden ser compartidas al repositorio remoto.

#### Configuración Básica (ya configuradas)
- `ENABLE_SQL_LOGS`: "false" ✅
- `ENABLE_SWAGGER`: "true" ✅
- `MAIN_PORT`: "8080" ✅

### 🔧 **Variables que DEBES configurar en Railway directamente:**

#### Base de Datos PostgreSQL (configurar según tu DB)
- `ASTRO_CAT_POSTGRES_HOST`: Host de tu base de datos
- `ASTRO_CAT_POSTGRES_PORT`: Puerto de PostgreSQL
- `ASTRO_CAT_POSTGRES_USER`: Usuario de la base de datos
- `ASTRO_CAT_POSTGRES_PASSWORD`: Contraseña de la base de datos
- `ASTRO_CAT_POSTGRES_NAME`: Nombre de la base de datos

#### JWT (OBLIGATORIO)
- `TOKEN_SIGNATURE_KEY`: Clave secreta para firmar tokens JWT

#### Email (configurar con tus credenciales)
- `EMAIL_HOST`: "smtp.gmail.com" ✅
- `EMAIL_PORT`: "587" ✅
- `EMAIL_USER`: Tu email de Gmail
- `EMAIL_PASSWORD`: Contraseña de aplicación de Gmail
- `EMAIL_FROM`: Tu email de Gmail

#### Twilio (configurar con tus credenciales)
- `TWILIO_ACCOUNT_SID`: Tu Account SID de Twilio
- `TWILIO_AUTH_TOKEN`: Tu Auth Token de Twilio
- `TWILIO_PHONE_NUMBER`: Tu número de teléfono de Twilio

#### AWS S3 (configurar con tus credenciales)
- `AWS_ACCESS_KEY_ID`: Tu Access Key ID de AWS
- `AWS_SECRET_ACCESS_KEY`: Tu Secret Access Key de AWS
- `AWS_SESSION_TOKEN`: Tu Session Token de AWS (si usas credenciales temporales)
- `AWS_REGION`: "us-east-1" ✅
- `S3_BUCKET_NAME`: Nombre de tu bucket S3

## Pasos para Deploy

1. **Edita `railway.json`**: Reemplaza los valores placeholder con tus credenciales reales
2. **Push a GitHub**: Sube los cambios a tu repositorio `main`
3. **Conecta a Railway**: Ve a Railway.app y conecta tu repositorio
4. **Deploy automático**: Railway usará la configuración de `railway.json`

## Health Check

La aplicación expone un endpoint de health check en `/health-check/` que Railway usará para verificar que la aplicación esté funcionando correctamente.

## Puertos

- Railway asignará automáticamente el puerto a través de la variable `PORT`
- La aplicación está configurada para usar el puerto 8080 por defecto si no se especifica otro

## Base de Datos

Para la base de datos PostgreSQL, puedes:
1. **Usar Railway PostgreSQL**: Las variables se configurarán automáticamente
2. **Usar base de datos externa**: Configura las variables en `railway.json`

## ⚠️ **Importante**

- **NO** subas credenciales reales a GitHub
- Usa variables de entorno en Railway Dashboard para credenciales sensibles
- El archivo `railway.json` debe estar en tu repositorio para que Railway lo use
