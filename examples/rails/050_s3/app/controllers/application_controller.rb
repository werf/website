class ApplicationController < ActionController::Base
  def image
    render template: "layouts/image"
  end

  def ping
    render plain: "pong\n"
  end
end
