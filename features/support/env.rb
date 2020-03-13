require 'aruba/cucumber'
require 'docker'
require 'fileutils'
require 'forwardable'
require 'tmpdir'

bin_dir = File.expand_path('../fakebin', __FILE__)
aruba_dir = File.expand_path('../../..', __FILE__) + '/tmp/aruba'

Before do
  # increase process exit timeout from the default of 3 seconds
  @aruba_timeout_seconds = 10
  # don't be "helpful"
  @aruba_keep_ansi = true

  begin
    app = Docker::Container.get("dummyvault")
    app.delete(force: true)
  rescue Docker::Error::NotFoundError
  end
  FileUtils.rm_rf("#{aruba_dir}/vault-cert-info-int-test")
end

After do
  begin
    app = Docker::Container.get("dummyvault")
    app.delete(force: true)
  rescue Docker::Error::NotFoundError
  end
  FileUtils.rm_rf("#{aruba_dir}/vault-cert-info-int-test")
end

RSpec::Matchers.define :be_successful_command do
  match do |cmd|
    cmd.success?
  end

  failure_message do |cmd|
    %(command "#{cmd}" exited with status #{cmd.status}:) <<
      cmd.output.gsub(/^/, ' ' * 2)
  end
end

class SimpleCommand
  attr_reader :output
  extend Forwardable

  def_delegator :@status, :exitstatus, :status
  def_delegators :@status, :success?

  def initialize cmd
    @cmd = cmd
  end

  def to_s
    @cmd
  end

  def self.run cmd
    command = new(cmd)
    command.run
    command
  end

  def run
    @output = `#{@cmd} 2>&1`.chomp
    @status = $?
    $?.success?
  end
end

World Module.new {
  # If there are multiple inputs, e.g., type in username and then type in password etc.,
  # the Go program will freeze on the second input. Giving it a small time interval
  # temporarily solves the problem.
  # See https://github.com/cucumber/aruba/blob/7afbc5c0cbae9c9a946d70c4c2735ccb86e00f08/lib/aruba/api.rb#L379-L382
  def type(*args)
    super.tap { sleep 0.1 }
  end

  def run_silent cmd
    in_current_dir do
      command = SimpleCommand.run(cmd)
      expect(command).to be_successful_command
      command.output
    end
  end

  # Aruba unnecessarily creates new Announcer instance on each invocation
  def announcer
    @announcer ||= super
  end

  def shell_escape(message)
    message.to_s.gsub(/['"\\ $]/) { |m| "\\#{m}" }
  end

  %w[output_from stdout_from stderr_from all_stdout all_stderr].each do |m|
    define_method(m) do |*args|
      home = ENV['HOME'].to_s
      output = super(*args)
      if home.empty?
        output
      else
        output.gsub(home, '$HOME')
      end
    end
  end
}
