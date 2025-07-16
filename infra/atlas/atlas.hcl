variable "db_user" {
  type = string
  default = getenv("DB_USER")
}

variable "db_password" {
  type = string
  default = getenv("DB_PASSWORD")
}

variable "db_host" {
  type = string
  default = getenv("DB_HOST")
}

variable "db_port" {
  type = string
  default = getenv("DB_PORT")
}

variable "db_name" {
  type = string
  default = getenv("DB_NAME")
}

variable "db_sslmode" {
  type = string
  default = getenv("DB_SSLMODE")
}

env "local" {
  url = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=${var.db_sslmode}"

  dev = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=${var.db_sslmode}"

  migration {
    dir = "file://infra/atlas/migrations"
  }
}
