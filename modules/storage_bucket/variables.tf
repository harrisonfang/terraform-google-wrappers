variable "name" {
  description = "A name for the storage bucket."
  type        = string
}

variable "region" {
  description = "A region in which to create the storage bucket, e.g. `us-west1`."
  type        = string

  validation {
    condition     = strcontains(var.location, "-")
    error_message = "Please select a regional bucket location. Multiregion buckets incur higher networking costs."
  }
}

variable "object_ttl_days" {
  description = "The number of days to store objects. After which, objects will be automatically deleted."
  type        = number
  default     = null
}

variable "read_users" {
  description = "A list of users that should have read access to the storage bucket."
  type        = list(string)
  default     = []
}

variable "write_once_users" {
  description = <<-EOT
    A list of users that should only have write access to the storage bucket. It is called "write once" because they are
    not allowed to read or overwrite existing data.
    EOT
  type        = list(string)
  default     = []
}

variable "read_write_users" {
  description = <<-EOT
    A list of users that should have both read and write access to the storage bucket. They are allowed to overwrite
    existing data.
    EOT
  type        = list(string)
  default     = []
}
