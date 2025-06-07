import uuid

class EquivalenceClasses:
    def __init__(self):
        pass


    def communityId_FETCH(self):
        valid_class = [
            "e804b95a-a388-4751-b246-96fe97232d35", #In BD
            "e804b95aa3884751b24696fe97232d35" #In BD
        ]

        invalid_class = [
            "e3ee1c1b-d6a8-4865-9b03-fd12e2edf644", #Not in BD
            "e3ee1c1bd6a848659b03fd12e2edf644", #Not in BD
            "e804b95a-a388-4751-b246-96fe97232d3",
            "e804b95a-a388-4751-b246-96fe97232d35X",
            "e8|4b@5a-#388-4751-b246-96fe97232d35"
        ]

        return [valid_class, invalid_class]

    def planId_FETCH(self):
        valid_class = [
            "d1694efe-9a13-42d7-a9e8-4d629f9f2f35", #In BD
            "d1694efe9a1342d7a9e84d629f9f2f35" #In BD
        ]

        invalid_class = [
            "5e0dc3ba-7e76-430b-9536-1d065b6c1cd0", #Not in BD
            "5e0dc3ba-7e76-430b-9536-1d065b6c1cd0", #Not in BD
            "e804b95a-a388-4751-b246-96fe97232d3",
            "e804b95a-a388-4751-b246-96fe97232d35X",
            "e8|4b@5a-#388-4751-b246-96fe97232d35"
        ]

        return [valid_class, invalid_class]


    def communityId_POST(self):
        valid_class = [
            "769350ca-bde9-4c77-ab95-69d3c9c83ab7" #In BD
        ]

        invalid_class = [
            "e3ee1c1b-d6a8-4865-9b03-fd12e2edf644", #Not in BD
            "e3ee1c1bd6a848659b03fd12e2edf644", #Not in BD
            "e804b95a-a388-4751-b246-96fe97232d3",
            "e804b95a-a388-4751-b246-96fe97232d35X",
            "e8|4b@5a-#388-4751-b246-96fe97232d35"
        ]

        return [valid_class, invalid_class]




    def planId_POST(self):
        valid_class = [
            "b038237d-16cd-411f-bfb7-f4077fbc92b5" #In BD
        ]

        invalid_class = [
            "5e0dc3ba-7e76-430b-9536-1d065b6c1cd0", #Not in BD
            "5e0dc3ba-7e76-430b-9536-1d065b6c1cd0", #Not in BD
            "e804b95a-a388-4751-b246-96fe97232d3",
            "e804b95a-a388-4751-b246-96fe97232d35X",
            "e8|4b@5a-#388-4751-b246-96fe97232d35"
        ]

        return [valid_class, invalid_class]

# Equivalencias para el campo Id (UUID v4) del modelo Plan

    def plan_id(self):
        valid_class = [
            # UUIDs válidos (formato estándar v4: 8-4-4-4-12 caracteres hex)
            "22222222-2222-4222-8222-222222222222",
            "33333333-3333-4333-8333-333333333333"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            "zzzzzzzz-zzzz-4zzz-8zzz-zzzzzzzzzzzz",  # Caracteres no hexadecimales
            "22222222-2222-3222-8222-222222222222",  # Tercer bloque no comienza con '4' (versión inválida)
            "22222222-2222-4222-1222-222222222222",  # Cuarto bloque no comienza con 8, 9, a o b (variant inválido)
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]




    # Equivalencias para el campo Fee (float64 > 0 y con máximo 2 decimales)

    def plan_fee(self):
        valid_class = [
            70.0,        # caso real usado en dummy_data.go
            1000.00,     # valor alto válido
            0.01,        # mínimo positivo permitido
            19.99,       # 2 decimales, válido
            150          # entero tratado como float
        ]

        invalid_class = [
            0.0,         # no se permite monto cero
            -10.0,       # monto negativo
            15.999,      # más de 2 decimales
            100.0001,    # más de 2 decimales aunque visualmente parezca correcto
            "100.0",     # string numérico: tipo incorrecto
            "veinte",    # texto no numérico
            None,        # nulo
            True,        # booleano (interpretable como 1.0, pero semánticamente inválido)
            [],          # lista vacía
            {},          # diccionario
            float("nan"), # NaN no representa valor monetario
            float("inf")  # infinito no aceptado como fee
        ]

        return [valid_class, invalid_class]

    # Equivalencias para el campo Type (enum: "MONTHLY", "ANUAL")

    def plan_type(self):
        valid_class = [
            "MONTHLY",   # plan mensual válido según enum
            "ANUAL"      # plan anual válido según enum
        ]
        invalid_class = [
            "monthly",       # minúscula, no coincide con el enum (Go es case-sensitive)
            "Mensual",       # en español, no pertenece al enum
            "ANNUAL",        # traducción común pero no definida en el backend
            "MENSUAL",       # mayúscula pero en español
            "",              # cadena vacía
            None,            # valor nulo
            123,             # tipo numérico inválido
            True,            # booleano
            [],              # lista vacía
            "MONTHLYY"       # error de tipeo: string similar pero incorrecto
        ]

        return [valid_class, invalid_class]


    # Equivalencias para ReservationLimit (*int en Go)
    # - Este campo puede ser un entero mayor o igual a 0 o None (puntero nulo)
    # - No se permiten decimales, negativos ni tipos incorrectos

    def plan_reservation_limit(self):
        valid_class = [
            0,      # mínimo válido, puede interpretarse como "sin reservas permitidas"
            1,      # mínimo operativo (una reserva)
            8,      # valor usado en dummy_data.go
            50,     # valor arbitrario válido
            999,    # valor alto permitido
            None    # puntero nulo en Go: significa "sin tope definido"
        ]
        invalid_class = [
            -1,               # valores negativos no son aceptados en límites
            -100,             # extremo negativo
            7.5,              # valor decimal: el tipo debe ser entero exacto
            3.1416,           # otro ejemplo de decimal inválido
            "10",             # string numérico: tipo incorrecto
            "unlimited",      # string no parseable a entero ni equivalente semántico
            "",               # string vacío
            "None",           # string que representa null pero no es el valor real None
            True,             # booleano (equivale a 1 en Python pero no es válido por tipo)
            False,            # igual que el caso anterior (representa 0 pero no debe aceptarse)
            [],               # lista vacía: tipo estructural no válido
            {},               # diccionario vacío
            object(self),         # instancia de objeto genérico
            float("nan"),     # not-a-number: matemáticamente inválido
            float("inf"),     # infinito positivo no representa límite finito
            float("-inf")     # infinito negativo no válido
        ]

        return [valid_class, invalid_class]

    #----------------------------------------------------------------
    #B. COMMUNITY -----------------------------------------------------------

    # Equivalencias para el campo Id (UUID v4) del modelo Community

    def community_id(self):
        valid_class = [
            # UUIDs válidos (v4, formato estándar)
            "11111111-1111-4111-8111-111111111111",
            "44444444-4444-4444-8444-444444444444"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo Name (string no vacío, semánticamente válido)

    def community_name(self):
        valid_class = [
            "Yoga Community",
            "Gym Group",
            "Red de Bienestar Mental"
        ]
        invalid_class = [
            "",           # vacío
            "  ",         # solo espacios
            None,
            "A",          # demasiado corto
            123           # tipo incorrecto
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo Purpose (string claro y descriptivo)

    def community_purpose(self):
        valid_class = [
            "Community for yoga enthusiasts",
            "Community for meditation practitioners",
            "Espacio para el bienestar emocional y físico"
        ]
        invalid_class = [
            "",            # vacío
            "###",         # sin sentido
            None,
            456            # tipo incorrecto
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo ImageUrl (string con ruta válida)

    def community_image_url(self):
        valid_class = [
            "test-image",
            "banner-zen-cat.png",
            "https://cdn.zen-cat.com/img/community.png"
        ]
        invalid_class = [
            "",         # vacío
            None,
            1234,       # no es string
            "////",     # sin semántica útil
            " "         # espacio solo
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo NumberSubscriptions (int ≥ 0)

    def community_number_subscriptions(self):
        valid_class = [
            0,
            5,
            100
        ]
        invalid_class = [
            -1,          # no se permiten negativos
            None,
            "ten",       # string
            5.5          # decimal
        ]
        return [valid_class, invalid_class]

     #----------------------------------------------------------------
    #C. USER -----------------------------------------------------------
    # Equivalencias para el campo Id (UUID v4) del modelo User

    def user_id(self):
        valid_class = [
            "55555555-5555-4555-8555-555555555555",
            "66666666-6666-4666-8666-666666666666"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo Name (string no vacío y semántico)

    def user_name(self):
        valid_class = [
            "Alice",
            "Bob",
            "Test-1"
        ]
        invalid_class = [
            "",       # vacío
            "  ",     # solo espacios
            None,
            123       # tipo incorrecto
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo FirstLastName (string obligatorio)

    def user_first_last_name(self):
        valid_class = [
            "Hurtado",
            "Doe",
            "User"
        ]
        invalid_class = [
            "",       # vacío
            None,
            456,
            " "       # espacios sin contenido
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo SecondLastName (opcional, puede ser null o string)

    def user_second_last_name(self):
        valid_class = [
            "Aparicio",
            None,
            ""    # en algunos casos se permite vacío
        ]
        invalid_class = [
            789,        # tipo incorrecto
            [],         # lista
            {},         # diccionario
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo Password (string no vacío)

    def user_password(self):
        valid_class = [
            "test123",
            "securePass!2024",
            "abcDEF123"
        ]
        invalid_class = [
            "",        # vacío
            "  ",
            None,
            123456     # tipo incorrecto
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo Email (formato correcto de correo)

    def user_email(self):
        valid_class = [
            "test-1@zen-cat.com",
            "user@example.com"
        ]
        invalid_class = [
            "user@",           # sin dominio
            "@zen-cat.com",    # sin nombre
            "correo",          # sin arroba
            "",                # vacío
            None
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo Rol (enum: "ADMINISTRATOR", "CLIENT")

    def user_rol(self):
        valid_class = [
            "ADMINISTRATOR",
            "CLIENT"
        ]
        invalid_class = [
            "admin",      # minúscula
            "cliente",    # fuera del enum
            "ADMIN",      # incorrecto
            "",           # vacío
            None
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo ImageUrl (string con ruta semántica)

    def user_image_url(self):
        valid_class = [
            "test-image",
            "user-profile.png",
            "https://cdn.zen-cat.com/img/u01.png"
        ]
        invalid_class = [
            "",         # vacío
            None,
            "   ",      # espacios
            123,        # no es string
            "///"       # sin semántica
        ]
        return [valid_class, invalid_class]
    #----------------------------------------------------------------
    #D. SERVICE -----------------------------------------------------------

    # Equivalencias para el campo Id (UUID v4) del modelo Service

    def service_id(self):
        valid_class = [
            "77777777-7777-4777-8777-777777777777",
            "88888888-8888-4888-8888-888888888888"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo Name (string no vacío y semántico)

    def service_name(self):
        valid_class = [
            "Yoga",
            "GYM",
            "Citas Médicas"
        ]
        invalid_class = [
            "",        # vacío
            None,
            123,       # tipo incorrecto
            "  ",      # espacios vacíos
            "A"        # demasiado corto
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo Description (string con sentido)

    def service_description(self):
        valid_class = [
            "Servicio de yoga",
            "Servicio de gimnasio",
            "Servicio online de citas médicas"
        ]
        invalid_class = [
            "",        # vacío
            None,
            "###",     # sin semántica
            456,       # tipo incorrecto
            "  "
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo ImageUrl (string con ruta válida o referencia semántica)

    def service_image_url(self):
        valid_class = [
            "test-image",
            "service-banner.jpg",
            "https://cdn.zen-cat.com/services/yoga.png"
        ]
        invalid_class = [
            "",         # vacío
            None,
            123,        # no es string
            "///",      # sin semántica
            "    "
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo IsVirtual (booleano explícito)

    def service_is_virtual(self):
        valid_class = [
            True,
            False
        ]
        invalid_class = [
            "true",     # string
            "false",
            None,
            1,
            0,
            "",         # vacío
            "sí"
        ]
        return [valid_class, invalid_class]

    #----------------------------------------------------------------
    #E. Professional -----------------------------------------------------------

    # Equivalencias para el campo Id (UUID v4) del modelo Professional

    def professional_id(self):
        valid_class = [
            "99999999-9999-4999-8999-999999999999",
            "aaaaaaa1-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para Name (string obligatorio)

    def professional_name(self):
        valid_class = [
            "John",
            "Jane",
            "María Fernanda"
        ]
        invalid_class = [
            "",      # vacío
            " ",     # espacios sin contenido
            None,
            555      # tipo incorrecto
        ]
        return [valid_class, invalid_class]

    # Equivalencias para FirstLastName (string obligatorio)

    def professional_first_last_name(self):
        valid_class = [
            "Doe",
            "Smith",
            "García"
        ]
        invalid_class = [
            "",      # vacío
            None,
            "  ",
            123
        ]
        return [valid_class, invalid_class]

    # Equivalencias para SecondLastName (opcional, puede ser null)

    def professional_second_last_name(self):
        valid_class = [
            "López",
            "",
            None
        ]
        invalid_class = [
            456,
            [],
            {}
        ]
        return [valid_class, invalid_class]

    # Equivalencias para Specialty (string libre, pero no vacío)

    def professional_specialty(self):
        valid_class = [
            "Yoga",
            "Cardiología",
            "Fisioterapia"
        ]
        invalid_class = [
            "",         # vacío
            " ",        # espacios
            None,
            789         # tipo incorrecto
        ]
        return [valid_class, invalid_class]

    # Equivalencias para Email (formato válido de correo electrónico)

    def professional_email(self):
        valid_class = [
            "john@gmail.com",
            "jane@zen-cat.com",
            "profesional@salud.pe"
        ]
        invalid_class = [
            "correo.com",        # sin arroba
            "@gmail.com",        # sin usuario
            "nombre@",           # sin dominio
            "",                  # vacío
            None
        ]
        return [valid_class, invalid_class]

    # Equivalencias para PhoneNumber (string numérica, no vacía)

    def professional_phone_number(self):
        valid_class = [
            "123456789",
            "987654321"
        ]
        invalid_class = [
            "",           # vacío
            "abc123",     # mezcla no numérica
            None,
            123456789     # tipo incorrecto (entero)
        ]
        return [valid_class, invalid_class]

    # Equivalencias para Type (enum: "MEDIC", "GYM_TRAINER", "YOGA_TRAINER")

    def professional_type(self):
        valid_class = [
            "MEDIC",
            "GYM_TRAINER",
            "YOGA_TRAINER"
        ]
        invalid_class = [
            "medico",        # español
            "trainer",       # fuera del enum
            "",              # vacío
            None,
            "YOGA"           # parcial
        ]
        return [valid_class, invalid_class]

    # Equivalencias para ImageUrl (string con valor semántico)

    def professional_image_url(self):
        valid_class = [
            "test-image",
            "profile.jpg",
            "https://cdn.zen-cat.com/img/p001.png"
        ]
        invalid_class = [
            "",          # vacío
            None,
            1000,        # tipo incorrecto
            "   ",       # espacios
            "///"        # sin sentido
        ]
        return [valid_class, invalid_class]

    #----------------------------------------------------------------
    #F. Local -----------------------------------------------------------

    # Equivalencias para el campo Id (UUID v4) del modelo Local

    def local_id(self):
        valid_class = [
            "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
            "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para LocalName (nombre del local)

    def local_name(self):
        valid_class = [
            "Local Gym",
            "Local Yoga",
            "Sala de Reuniones"
        ]
        invalid_class = [
            "",        # vacío
            "  ",      # espacios en blanco
            None,
            789        # tipo incorrecto
        ]
        return [valid_class, invalid_class]

    # Equivalencias para StreetName (nombre de la calle)

    def local_street_name(self):
        valid_class = [
            "Main St",
            "Downtown Ave",
            "Av. San Martín"
        ]
        invalid_class = [
            "",       # vacío
            None,
            456,      # tipo incorrecto
            "   "
        ]
        return [valid_class, invalid_class]

    # Equivalencias para BuildingNumber (número de edificio)

    def local_building_number(self):
        valid_class = [
            "123",
            "456",
            "B-21"
        ]
        invalid_class = [
            "",         # vacío
            None,
            [],         # tipo inválido
            " "         # solo espacios
        ]
        return [valid_class, invalid_class]

    # Equivalencias para District (nombre del distrito)

    def local_district(self):
        valid_class = [
            "Downtown",
            "Business",
            "Miraflores"
        ]
        invalid_class = [
            "",        # vacío
            None,
            123,       # tipo incorrecto
            "   "
        ]
        return [valid_class, invalid_class]

    # Equivalencias para Province

    def local_province(self):
        valid_class = [
            "Central",
            "Lima",
            "Cusco"
        ]
        invalid_class = [
            "",        # vacío
            None,
            888        # tipo incorrecto
        ]
        return [valid_class, invalid_class]

    # Equivalencias para Region

    def local_region(self):
        valid_class = [
            "Metropolitan",
            "Costa",
            "Sierra"
        ]
        invalid_class = [
            "",        # vacío
            None,
            "   ",
            555
        ]
        return [valid_class, invalid_class]

    # Equivalencias para Reference (referencia adicional)

    def local_reference(self):
        valid_class = [
            "Near Central Park",
            "Next to Admin Office",
            "Al lado de cafetería"
        ]
        invalid_class = [
            "",        # vacío
            None,
            {},        # tipo incorrecto
            "  "
        ]
        return [valid_class, invalid_class]

    # Equivalencias para Capacity (capacidad del local, entero ≥ 0)

    def local_capacity(self):
        valid_class = [
            0,
            10,
            50,
            100
        ]
        invalid_class = [
            -1,        # negativo
            7.5,       # decimal
            "20",      # string numérico
            "",        # vacío
            None
        ]
        return [valid_class, invalid_class]

    # Equivalencias para ImageUrl (string válido con ruta o nombre)

    def local_image_url(self):
        valid_class = [
            "test-image",
            "local01.png",
            "https://zen-cat.com/images/local-x.jpg"
        ]
        invalid_class = [
            "",        # vacío
            None,
            0,
            "///",     # sin semántica
            "   "
        ]
        return [valid_class, invalid_class]

    def membership_id(self):
        valid_class = [
            "cccccccc-cccc-4ccc-8ccc-cccccccccccc",
            "dddddddd-dddd-4ddd-8ddd-dddddddddddd"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para Description (string no vacío)

    def membership_description(self):
        valid_class = [
            "Monthly Yoga Membership",
            "Yearly Gym Membership",
            "Acceso completo a comunidad"
        ]
        invalid_class = [
            "",        # vacío
            None,
            123,       # tipo incorrecto
            "  "
        ]
        return [valid_class, invalid_class]

    # Equivalencias para StartDate (formato de fecha válido)

    def membership_start_date(self):
        valid_class = [
            "2024-06-01T00:00:00Z",
            "2025-01-15T12:00:00Z"
        ]
        invalid_class = [
            "2024-13-01",    # mes inválido
            "hoy",           # no es fecha
            "",              # vacío
            None
        ]
        return [valid_class, invalid_class]

    # Equivalencias para EndDate (fecha posterior a StartDate)

    def membership_end_date(self):
        valid_class = [
            "2024-07-01T00:00:00Z",
            "2026-01-15T12:00:00Z"
        ]
        invalid_class = [
            "2020-01-01T00:00:00Z",  # fecha en pasado (ejemplo inválido)
            "",                      # vacío
            "fecha",                 # string no válido
            None
        ]
        return [valid_class, invalid_class]

    # Equivalencias para Status (enum válido)

    def membership_status(self):
        valid_class = [
            "ACTIVE",
            "EXPIRED",
            "CANCELLED"
        ]
        invalid_class = [
            "active",     # minúscula
            "activa",     # español
            "",           # vacío
            None,
            "INACTIVO"
        ]
        return [valid_class, invalid_class]

    # Equivalencias para CommunityId (UUID v4)

    def membership_community_id(self):
        valid_class = [
            "11111111-1111-4111-8111-111111111111",
            "44444444-4444-4444-8444-444444444444"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para UserId (UUID v4)

    def membership_user_id(self):
        valid_class = [
            "55555555-5555-4555-8555-555555555555",
            "66666666-6666-4666-8666-666666666666"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para PlanId (UUID v4)

    def membership_plan_id(self):
        valid_class = [
            "22222222-2222-4222-8222-222222222222",
            "33333333-3333-4333-8333-333333333333"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]


    #----------------------------------------------------------------
    #H. CommunityService -----------------------------------------------------------

    # Equivalencias para el campo Id (UUID v4) del modelo CommunityService

    def community_service_id(self):
        valid_class = [
            "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
            "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo CommunityId (UUID v4)

    def community_service_community_id(self):
        valid_class = [
            "11111111-1111-4111-8111-111111111111",
            "44444444-4444-4444-8444-444444444444"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo ServiceId (UUID v4)

    def community_service_service_id(self):
        valid_class = [
            "77777777-7777-4777-8777-777777777777",
            "88888888-8888-4888-8888-888888888888"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    #----------------------------------------------------------------
    #I. CommunityPlan -----------------------------------------------------------

    # Equivalencias para el campo Id (UUID v4) del modelo CommunityPlan

    def community_plan_id(self):
        valid_class = [
            "eeeeeeee-eeee-4eee-8eee-eeeeeeeeeeee",
            "ffffffff-ffff-4fff-8fff-ffffffffffff"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo CommunityId (UUID v4)

    def community_plan_community_id(self):
        valid_class = [
            "11111111-1111-4111-8111-111111111111",
            "44444444-4444-4444-8444-444444444444"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Equivalencias para el campo PlanId (UUID v4)

    def community_plan_plan_id(self):
        valid_class = [
            "22222222-2222-4222-8222-222222222222",
            "33333333-3333-4333-8333-333333333333"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    #----------------------------------------------------------------
    #J. Session -----------------------------------------------------------
    # Id de la sesión: debe ser un UUID v4 válido

    def session_id(self):
        valid_class = [
            "aaaa1111-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
            "bbbb2222-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # Título de la sesión: texto obligatorio no vacío

    def session_title(self):
        valid_class = [
            "Morning Yoga",
            "Evening Gym",
            "Cardio Intensivo"
        ]
        invalid_class = [
            "",                   # vacío
            None,                 # nulo
            123,                  # tipo incorrecto
            "   "                 # solo espacios
        ]
        return [valid_class, invalid_class]

    # Fecha de la sesión: formato ISO 8601

    def session_date(self):
        valid_class = [
            "2025-06-01T09:00:00Z",
            "2025-12-24T18:00:00Z"
        ]
        invalid_class = [
            "fecha",              # valor no fecha
            "",                   # vacío
            None,                 # nulo
            "2025-13-01"          # mes inválido
        ]
        return [valid_class, invalid_class]

    # Hora de inicio: formato ISO 8601

    def session_start_time(self):
        valid_class = [
            "2025-06-01T08:00:00Z",
            "2025-06-02T10:30:00Z"
        ]
        invalid_class = [
            "",                   # vacío
            None,                 # nulo
            9,                    # número
            "8am"                 # formato no compatible
        ]
        return [valid_class, invalid_class]

    # Hora de fin: debe ser posterior a la de inicio

    def session_end_time(self):
        valid_class = [
            "2025-06-01T09:00:00Z",
            "2025-06-02T11:30:00Z"
        ]
        invalid_class = [
            "2020-01-01T00:00:00Z",  # antes del inicio
            "",                      # vacío
            None,                    # nulo
            "end"                    # cadena inválida
        ]
        return [valid_class, invalid_class]

    # Estado de la sesión: debe ser parte del enum definido

    def session_state(self):
        valid_class = [
            "SCHEDULED",
            "ONGOING",
            "COMPLETED",
            "CANCELLED",
            "RESCHEDULED"
        ]
        invalid_class = [
            "",                    # vacío
            "activo",              # traducción no válida
            "iniciada",            # español sin match
            "completed",           # incorrecto por minúscula
            None                   # nulo
        ]
        return [valid_class, invalid_class]

    # Cantidad de inscritos: debe ser un entero ≥ 0

    def session_registered_count(self):
        valid_class = [0, 5, 20]
        invalid_class = [
            -1,                   # negativo
            None,                 # nulo
            "10",                 # string numérico
            3.5                   # decimal
        ]
        return [valid_class, invalid_class]

    # Capacidad: debe ser entero estrictamente > 0

    def session_capacity(self):
        valid_class = [1, 15, 50]
        invalid_class = [
            0,                    # cero no permitido
            -5,                   # negativo
            "20",                 # string numérico
            None                  # nulo
        ]
        return [valid_class, invalid_class]

    # Enlace de sesión (opcional): puede ser string o null

    def session_link(self):
        valid_class = [
            None,
            "",
            "https://meet.zen-cat.com/session-01"
        ]
        invalid_class = [
            123,                  # número
            {},                   # objeto no válido
            True                  # booleano
        ]
        return [valid_class, invalid_class]

    # ID del profesional: UUID v4 obligatorio

    def session_professional_id(self):
        valid_class = [
            "44444444-4444-4444-8444-444444444444",
            "55555555-5555-4555-8555-555555555555"
        ]
        invalid_class = [
            "",  # Cadena vacía
            "123",  # Cadena con menos de 36 caracteres
            "12345678-1234-1234-1234-1234567890123",  # Cadena con más de 36 caracteres
            "12345678123412341234123456789012",  # Cadena sin guiones
            None,  # Valor nulo
            123456,  # Tipo incorrecto: número
            {},  # Tipo incorrecto: diccionario
            True  # Tipo incorrecto: booleano
        ]
        return [valid_class, invalid_class]

    # ID del local: puede ser UUID válido o null

    def session_local_id(self):
        valid_class = [
            "66666666-6666-4666-8666-666666666666",
            None
        ]
        invalid_class = [
            "",                   # vacío
            "12345678123412341234123456789012",       # mal formado
            {},                   # objeto
            0                     # número entero
        ]
        return [valid_class, invalid_class]

    #----------------------------------------------------------------
    #K. Reservation -----------------------------------------------------------

    # Id de la reserva: UUID v4 válido

    def reservation_id(self):
        valid_class = [
            "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
            "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
        ]
        invalid_class = [
            "",                  # cadena vacía
            "12345678123412341234123456789012",           # no cumple formato UUID
            None,                # valor nulo
            1234,                # tipo numérico
            {}                   # objeto no válido
        ]
        return [valid_class, invalid_class]

    # Nombre de la reserva: texto no vacío

    def reservation_name(self):
        valid_class = [
            "Yoga Class Reservation",
            "Gym Session Reservation"
        ]
        invalid_class = [
            "",                  # texto vacío
            None,                # nulo
            456,                 # número
            "   "                # solo espacios
        ]
        return [valid_class, invalid_class]

    # Fecha y hora de la reserva: formato ISO 8601

    def reservation_time(self):
        valid_class = [
            "2025-06-07T09:00:00Z",
            "2025-12-24T18:30:00Z"
        ]
        invalid_class = [
            "",                  # vacío
            None,                # nulo
            "ayer",              # texto informal
            777                  # valor numérico inválido
        ]
        return [valid_class, invalid_class]

    # Estado de la reserva: debe ser uno de los definidos en el enum

    def reservation_state(self):
        valid_class = [
            "DONE",
            "CONFIRMED",
            "CANCELLED",
            "ANULLED"
        ]
        invalid_class = [
            "",                  # vacío
            "activa",            # valor desconocido
            None,                # nulo
            "done"               # formato incorrecto (minúsculas)
        ]
        return [valid_class, invalid_class]

    # Fecha de última modificación: formato ISO 8601

    def reservation_last_modification(self):
        valid_class = [
            "2025-06-06T10:00:00Z",
            "2025-06-07T13:45:00Z"
        ]
        invalid_class = [
            "",                  # vacío
            "modificación",      # texto informal
            None,                # nulo
            0                    # tipo incorrecto
        ]
        return [valid_class, invalid_class]

    # ID del usuario asociado: UUID v4 requerido

    def reservation_user_id(self):
        valid_class = [
            "11111111-1111-4111-8111-111111111111",
            "22222222-2222-4222-8222-222222222222"
        ]
        invalid_class = [
            "",                  # vacío
            "12345678123412341234123456789012",           # no es UUID
            None,                # nulo
            9876                 # tipo numérico
        ]
        return [valid_class, invalid_class]

    # ID de la sesión reservada: UUID v4 requerido

    def reservation_session_id(self):
        valid_class = [
            "33333333-3333-4333-8333-333333333333",
            "44444444-4444-4444-8444-444444444444"
        ]
        invalid_class = [
            "",                  # vacío
            "12345678123412341234123456789012",          # mal formado
            None,                # nulo
            True                 # tipo booleano
        ]
        return [valid_class, invalid_class]

    #----------------------------------------------------------------
    #L. Template -----------------------------------------------------------

    # Id del template: debe ser UUID v4 válido

    def template_id(self):
        valid_class = [
            "aaaa1111-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
            "bbbb2222-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
        ]
        invalid_class = [
            "",                  # cadena vacía
            "12345678123412341234123456789012",        # formato no UUID
            None,                # valor nulo
            1234,                # tipo incorrecto
            {}                   # tipo objeto
        ]
        return [valid_class, invalid_class]

    # Link del template: debe ser una URL válida (requerido)

    def template_link(self):
        valid_class = [
            "https://example.com/template-001",
            "http://zen-cat.org/template/abc"
        ]
        invalid_class = [
            "",                  # vacío (no se permite)
            None,                # nulo
            "12345678123412341234123456789012",        # string sin formato URL
            123,                 # tipo numérico
            []                   # lista vacía
        ]
        return [valid_class, invalid_class]

    # Id del profesional asociado: debe ser UUID v4 válido

    def template_professional_id(self):
        valid_class = [
            "44444444-4444-4444-8444-444444444444",
            "55555555-5555-4555-8555-555555555555"
        ]
        invalid_class = [
            "",                  # cadena vacía
            "12345678123412341234123456789012",     # texto inválido
            None,                # nulo
            9999,                # entero
            True                 # booleano
        ]
        return [valid_class, invalid_class]

    #----------------------------------------------------------------
    #M. Onboarding  -----------------------------------------------------------


    # Id del onboarding: UUID v4 válido

    def onboarding_id(self):
        valid_class = [
            "33333333-1111-4111-8111-111111111111",
            "33333333-2222-4222-8222-222222222222"
        ]
        invalid_class = [
            "",                    # vacío
            "no-uuid",             # texto no UUID
            None,                  # nulo
            1234,                  # número
            {},                    # objeto
            "aaaaaaaa-aaaa-3aaa-8aaa-aaaaaaaaaaaa",  # UUID no v4
        ]
        return [valid_class, invalid_class]

    # Tipo de documento: enum válido

    def onboarding_document_type(self):
        valid_class = [
            "DNI",
            "FOREIGNER_CARD",
            "PASSPORT"
        ]
        invalid_class = [
            "",              # vacío
            None,            # nulo
            "dni",           # error por lowercase
            "CE",            # sigla mal representada
            123,             # tipo incorrecto
            [],              # lista
        ]
        return [valid_class, invalid_class]

    # Número de documento: texto no vacío

    def onboarding_document_number(self):
        valid_class = [
            "12345678",
            "87654321",
            "X1234567"
        ]
        invalid_class = [
            "",              # vacío
            None,            # nulo
            12345678,        # tipo numérico
            True,            # booleano
            [],              # lista
        ]
        return [valid_class, invalid_class]

    # Número de teléfono: string no vacío

    def onboarding_phone_number(self):
        valid_class = [
            "987654321",
            "912345678"
        ]
        invalid_class = [
            "",              # vacío
            None,
            987654321,       # número en vez de string
            [],              # estructura no válida
        ]
        return [valid_class, invalid_class]

    # Fecha de nacimiento: ISO string de fecha (opcional)

    def onboarding_birth_date(self):
        valid_class = [
            "2000-01-01",
            "1995-12-31",
            None             # puede omitirse
        ]
        invalid_class = [
            "31/12/1995",     # formato no ISO
            "",               # vacío
            None if False else "null",
            1995,
            {}
        ]
        return [valid_class, invalid_class]

    # Género: enum válido o null

    def onboarding_gender(self):
        valid_class = [
            "MALE",
            "FEMALE",
            "OTHER",
            None
        ]
        invalid_class = [
            "",              # vacío
            "masculino",     # fuera del enum
            "female",        # minúscula
            True,
            0,
        ]
        return [valid_class, invalid_class]

    # Ciudad: string obligatorio

    def onboarding_city(self):
        valid_class = [
            "Lima",
            "Arequipa"
        ]
        invalid_class = [
            "", None, 123, [], {}
        ]
        return [valid_class, invalid_class]

    # Código postal: string obligatorio

    def onboarding_postal_code(self):
        valid_class = [
            "15001",
            "04001"
        ]
        invalid_class = [
            "", None, 15001, True
        ]
        return [valid_class, invalid_class]

    # Distrito: string obligatorio

    def onboarding_district(self):
        valid_class = [
            "Downtown",
            "Business"
        ]
        invalid_class = [
            "", None, 999, {}, []
        ]
        return [valid_class, invalid_class]

    # Dirección completa: string obligatorio

    def onboarding_address(self):
        valid_class = [
            "Av. Principal 123",
            "Jr. Comercio 456"
        ]
        invalid_class = [
            "", None, 789, [], {}
        ]
        return [valid_class, invalid_class]

    # ID de usuario vinculado (único): UUID v4

    def onboarding_user_id(self):
        valid_class = [
            "99999999-9999-4999-8999-999999999999",
            "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
        ]
        invalid_class = [
            "", "uuid-mal", None, 5678
        ]
        return [valid_class, invalid_class]
