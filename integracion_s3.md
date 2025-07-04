# Integración con S3 para la Gestión de Imágenes

Esta guía describe los pasos para integrar la subida y descarga de imágenes desde y hacia un bucket de S3 para cualquier entidad dentro de este proyecto. La integración se basa en un servicio de S3 (`S3Service`) que abstrae las operaciones de subida (`UploadFile`) y descarga (`DownloadFile`), y sigue un patrón consistente para la modificación de los esquemas, la API y el enrutador.

A continuación, se detallan los pasos genéricos a seguir para una entidad `Entity`.

## 1. Actualización de Esquemas

El primer paso es asegurar que los esquemas de la entidad estén preparados para manejar los datos de las imágenes.

-   **Modificar los esquemas de creación y actualización**: En el archivo `src/server/schemas/entity.go`, añade un campo para los bytes de la imagen (`ImageBytes *[]byte`) a las estructuras de solicitud de creación y actualización (`CreateEntityRequest` y `UpdateEntityRequest`). Esto permitirá que la API reciba los datos de la imagen en formato de bytes.

-   **Crear un esquema para respuestas con imagen**: Define una nueva estructura, como `EntityWithImage`, que combine la información de la entidad con los bytes de su imagen. Este esquema se utilizará para las respuestas de los endpoints que devuelven tanto los datos de la entidad como la imagen asociada.

## 2. Modificación de la API

Con los esquemas actualizados, el siguiente paso es modificar la capa de la API para manejar la lógica de negocio relacionada con las imágenes.

-   **Crear un endpoint para obtener la entidad con su imagen**: En `src/server/api/entity.go`, implementa una nueva función, como `GetEntityWithImage`. Esta función debe:
    1.  Obtener los datos de la entidad desde la base de datos.
    2.  Utilizar el `S3Service` para descargar la imagen desde S3. Es recomendable que, si la descarga falla, la API no devuelva un error, sino que continúe la ejecución y devuelva la entidad sin los bytes de la imagen.
    3.  Devolver una respuesta utilizando el esquema `EntityWithImage`.

-   **Actualizar el endpoint de creación**: Modifica la función `CreateEntity` para que:
    1.  Genere una URL única para la imagen utilizando el `S3Service` antes de crear la entidad en la base de datos.
    2.  Después de crear la entidad, suba la imagen a S3 utilizando los `ImageBytes` recibidos en la solicitud. De manera similar a la descarga, un fallo en la subida no debería interrumpir el flujo; la entidad se crea de todos modos, y la imagen puede subirse más tarde.

-   **Actualizar el endpoint de actualización**: Modifica la función `UpdateEntity` para que:
    1.  Si se proporciona una nueva imagen, genere una nueva URL para ella.
    2.  Después de actualizar los datos de la entidad, suba la nueva imagen a S3. En este caso, un fallo en la subida sí debería devolver un error, ya que la actualización de la imagen es una parte explícita de la operación.

## 3. Registro de la Nueva Ruta

Finalmente, para que el nuevo endpoint sea accesible, debes registrarlo en el enrutador.

-   **Añadir la nueva ruta**: En `src/server/api/router.go`, dentro del grupo de rutas de la entidad, añade una nueva ruta `GET` para el endpoint `GetEntityWithImage`. Por ejemplo:

    ```go
    entity.GET("/:entityId/image/", a.GetEntityWithImage)
    ```

Siguiendo estos tres pasos, puedes integrar de manera consistente y robusta la gestión de imágenes con S3 para cualquier entidad del sistema, manteniendo el código limpio y modular.
