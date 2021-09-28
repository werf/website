require "active_support/core_ext/integer/time"
require "#{Rails.root}/lib/json_simple_formatter"

Rails.application.configure do
  # Settings specified here will take precedence over those in config/application.rb.

  # Code is not reloaded between requests.
  config.cache_classes = true

  # Eager load code on boot. This eager loads most of Rails and
  # your application in memory, allowing both threaded web servers
  # and those relying on copy on write to perform better.
  # Rake tasks automatically ignore this option for performance.
  config.eager_load = true

  # Full error reports are disabled and caching is turned on.
  config.consider_all_requests_local       = false
  config.action_controller.perform_caching = true

  # Disable serving static files from the `/public` folder by default since
  # Apache or NGINX already handles this.
  config.public_file_server.enabled = ENV['RAILS_SERVE_STATIC_FILES'].present?

  # Include generic and useful information about system operation, but avoid logging too much
  # information to avoid inadvertent exposure of personally identifiable information (PII).
  # [<snippet log_level>]
  config.log_level = ENV.fetch("RAILS_LOG_LEVEL", "info").downcase.strip.to_sym
  # [<endsnippet log_level>]

  # Prepend all log lines with the following tags.
  config.log_tags = [ :request_id ]


  # Enable locale fallbacks for I18n (makes lookups for any locale fall back to
  # the I18n.default_locale when a translation cannot be found).
  config.i18n.fallbacks = true

  # Send deprecation notices to registered listeners.
  config.active_support.deprecation = :notify

  # Log disallowed deprecations.
  config.active_support.disallowed_deprecation = :log

  # Tell Active Support which deprecation messages to disallow.
  config.active_support.disallowed_deprecation_warnings = []

  # [<snippet storage>]
  config.active_storage.service = :minio
  # [<endsnippet storage>]

  # [<snippet log_formatter>]
  config.log_formatter = JSONSimpleFormatter.new

  # [<snippet logger>]
  config.logger           = ActiveSupport::Logger.new(STDOUT)
  # [<endsnippet logger>]
  config.logger.formatter = config.log_formatter
  # [<endsnippet log_formatter>]

  # Do not dump schema after migrations.
  config.active_record.dump_schema_after_migration = false
end
