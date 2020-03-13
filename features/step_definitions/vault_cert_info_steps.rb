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
  sleep 3
end

When /^I run `bin\/vault-cert-info-int-test cert` with the first serial$/ do
  require 'vault'
  Vault.address = "http://0.0.0.0:8008"
  Vault.token   = "ROOT"

  first_cert = Vault.logical.list("pki/certs").first

  Vault.shutdown()
  steps %Q(
    When I successfully run `bin\/vault-cert-info-int-test cert #{first_cert}`
  )
end

Given("I have the PKI backend enabled with a test cert") do
  require 'vault'
  Vault.address = "http://0.0.0.0:8008"
  Vault.token   = "ROOT"

  Vault.sys.mount("pki", "pki", "PKI mount")
  Vault.logical.write("pki/root/generate/internal", "common_name": "arubatest.com", "ttl": "87600h")

  Vault.shutdown()
end