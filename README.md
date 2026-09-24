# terraform-provider-cloudlab

Terraform-провайдер для управления ресурсами [cloudlab](https://github.com/VNGritsai/cloudlab) — учебного REST API для виртуальных инстансов.

Это учебный проект: цель — разобраться, как устроены Terraform-провайдеры на Go, и собрать рабочий провайдер с нуля на Terraform Plugin Framework.

## Что умеет

Управляет ресурсом `cloudlab_instance`:

- **create** — создать инстанс
- **read** — читать состояние из API
- **update** — обновлять `name`, `cpu`, `memory`
- **delete** — удалять инстанс

## Требования

- Go 1.25+
- Terraform 1.0+
- Запущенный сервер `cloudlab` (см. [репозиторий сервера](https://github.com/VNGritsai/cloudlab))

## Установка (локально, для разработки)

```bash
git clone https://github.com/VNGritsai/terraform-provider-cloudlab.git
cd terraform-provider-cloudlab
go install .
```

Настройка `~/.terraformrc`:
```hcl
provider_installation {
  dev_overrides {
    "VNGritsai/cloudlab" = "/home/<user>/go/bin"
  }
  direct {}
}
```
## Использование

Пример `main.tf`:
```hcl
terraform {
  required_providers {
    cloudlab = {
      source = "VNGritsai/cloudlab"
    }
  }
}

provider "cloudlab" {
  endpoint = "http://localhost:8080"
}

resource "cloudlab_instance" "example" {
  name   = "my-vm"
  cpu    = 2
  memory = 1024
}
```
 
## Команды Terraform
```bash
terraform plan
terraform apply
terraform destroy
```
