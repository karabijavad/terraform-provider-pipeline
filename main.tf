provider banzaicloud {
    access_token = var.access_token
    api_url = var.api_url
    organization_id = var.organization_id
}

resource banzaicloud_aws_secret secret {
    name = "test-secret2"
    access_key_id = "abc"
    secret_access_key = "def"
    validate = false
}
