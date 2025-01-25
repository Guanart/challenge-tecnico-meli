# Challenge Técnico MELI - Gonzalo Benito
Repositorio para el challenge técnico de MELI: "API para Detección de Vulnerabilidades en Imágenes de Contenedores"

## Objetivo
El equipo de Seguridad de MercadoLibre está desarrollando un sistema para gestionar las vulnerabilidades de todos los contenedores corriendo en la infraestructura.

El objetivo será construir una API en Go (preferentemente) o Python que permita analizar e informar las vulnerabilidades encontradas para las imágenes de contenedores bajo el repositorio, para esto debe integrar algún scanner open source como Gripe o Trivy. La API deberá exponer endpoints que permitan:
1. Un endpoint POST que permita indicar la imagen a analizar.
2. Un endpoint GET que liste las imágenes indexadas.
3. Un endpoint GET que permita obtener el listado de vulnerabilidades para una determinada imagen.

## Solución propuesta
Para la solución, se utilizará el lenguaje Go con el framework Gin, junto a una base de datos sqlite para la persistencia. Para realizar el análisis de vulnerabilidades en las imágenes docker, se optó por el scanner Grype.

## Posibles mejoras
- Hacer un refactor en connection.go: hacerlo orientado a objetos y evitar hardcodear el path de la BD (cargar el valor de un archivo de configuración e inyectarlo a una clase Conexión)
- Desacoplar del archivo main.go las funciones asociadas a cada ruta. Crear un controller o service.