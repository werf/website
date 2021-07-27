class ApplicationController < ActionController::API
  def ping
    render plain: "pong\n"
  end
end
