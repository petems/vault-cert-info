require 'fileutils'

Given(/^I have "([^"]*)" command installed$/) do |command|
  is_present = system("which #{ command} > /dev/null 2>&1")
  raise "Command #{command} is not present in the system" if not is_present
end

Given("a build of vault-cert-info") do

end

Given(/nothings running on port "(\w+)"/) do |port|
  running_on_port = system("lsof -i TCP:#{port} > /dev/null 2>&1")
  if running_on_port
    system("lsof -i TCP:#{port}")
    raise "Something running on port #{port}"
  end
end

Given /^no old "(\w+)" containers exist$/ do |container_name|
  begin
    app = Docker::Container.get(container_name)
    app.delete(force: true)
  rescue Docker::Error::NotFoundError
  end
end

Given /^I have a dummy vault server running called "(\w+)" running on port "(\w+)" with root token "(\w+)"$/ do |container_name, port, token|
  steps %Q(
    Given no old "#{container_name}" containers exist
    When I successfully run `docker run --name='#{container_name}' -p #{port}:#{port} -d --cap-add=IPC_LOCK -e 'VAULT_DEV_ROOT_TOKEN_ID=#{token}' -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:#{port}' vault`
  )
  sleep 2
end

When("I run {string} with the first serial") do |string|
  require 'vault'
  Vault.address = "http://0.0.0.0:8008"
  Vault.token   = "ROOT"

  first_cert = Vault.logical.list("pki/certs").first

  Vault.shutdown()

  sleep 2
  steps %Q(
    When I successfully run `#{string} #{first_cert}`
  )
end

Given("I have the PKI backend enabled at {string} with a test cert") do |pki_path|
	require 'vault'
  Vault.address = "http://0.0.0.0:8008"
  Vault.token   = "ROOT"

  Vault.sys.mount(pki_path, "pki", "PKI mount")
  Vault.sys.mount_tune(pki_path, max_lease_ttl: '87600h')
  
  Vault.logical.write("#{pki_path}/root/generate/internal", "common_name": "arubatest.com", "ttl": "87600h")
  Vault.logical.write("#{pki_path}/config/urls", "issuing_certificates": "http://127.0.0.1:8200/v1/#{pki_path}/ca", "crl_distribution_points": "http://127.0.0.1:8200/v1/#{pki_path}/crl")

  Vault.logical.write("#{pki_path}/roles/arubatest-dot-com", "allowed_domains": "arubatest.com", "allow_subdomains": "true", "max_ttl": "87600h")

  Vault.shutdown()

  sleep 2
end

Given("I have a certificate that expires in {int} days") do |int_days|
  require 'vault'
  Vault.address = "http://0.0.0.0:8008"
  Vault.token   = "ROOT"

  string_days = int_days * 24

  Vault.logical.write("pki/issue/arubatest-dot-com", "common_name": "#{int_days}days.arubatest.com", "ttl": "#{string_days}h")

  Vault.shutdown()

  sleep 2
end