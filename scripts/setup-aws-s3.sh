#!/bin/bash

# Helper simple para actualizar las variables de entorno de AWS S3

ENV_FILE="src/server/.env"

echo "🔧 Configuración AWS S3 - zen-cat-backend"
echo "=========================================="

# Verificar que el archivo .env existe
if [ ! -f "$ENV_FILE" ]; then
    echo "❌ Error: No se encontró $ENV_FILE"
    echo "💡 Ejecuta 'make init-vscode' primero"
    exit 1
fi

echo "📝 Ingresa las credenciales de AWS S3:"
echo ""

# Leer credenciales
read -p "AWS Access Key ID: " aws_access_key
read -p "AWS Secret Access Key: " aws_secret_key
read -p "AWS Session Token: " aws_session_token
read -p "AWS Region (press Enter to default value: [us-east-1]): " aws_region
read -p "S3 Bucket Name (press Enter to default value: [astro-cat-ingesoft-20251]): " s3_bucket_name

# Valores por defecto
aws_region=${aws_region:-us-east-1}
s3_bucket_name=${s3_bucket_name:-astro-cat-ingesoft-20251}

echo ""
echo "🔄 Actualizando .env..."

# Función para actualizar una variable en el .env de forma segura
update_env_var() {
    local var_name="$1"
    local var_value="$2"

    # Crear archivo temporal
    local temp_file=$(mktemp)

    # Leer línea por línea y actualizar la variable correspondiente
    while IFS= read -r line; do
        if [[ $line =~ ^${var_name}[[:space:]]*= ]]; then
            echo "${var_name} = \"${var_value}\""
        else
            echo "$line"
        fi
    done < "$ENV_FILE" > "$temp_file"

    # Reemplazar el archivo original
    mv "$temp_file" "$ENV_FILE"
}

# Actualizar variables
update_env_var "AWS_ACCESS_KEY_ID" "$aws_access_key"
update_env_var "AWS_SECRET_ACCESS_KEY" "$aws_secret_key"
update_env_var "AWS_SESSION_TOKEN" "$aws_session_token"
update_env_var "AWS_REGION" "$aws_region"
update_env_var "S3_BUCKET_NAME" "$s3_bucket_name"

echo "✅ Variables de AWS S3 actualizadas!"
echo ""
echo "📋 Configuración:"
echo "   Access Key: ${aws_access_key:0:4}****"
echo "   Secret Key: ****"
echo "   Session Token: ${aws_session_token:+[CONFIGURADO]}"
echo "   Region: $aws_region"
echo "   Bucket: $s3_bucket_name"
