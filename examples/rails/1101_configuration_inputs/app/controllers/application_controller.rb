# frozen_string_literal: true

class ApplicationController < ActionController::API
  rescue_from StandardError do |exception|
    render_error exception.message
  end

  rescue_from ActiveRecord::RecordNotFound do |_exception|
    render_error 'entry not found', status: :forbidden
  end

  # Renders error with passed data
  def render_error(comment = nil, result: :error, status: :ok)
    render json: { result: result, comment: comment }, status: status
  end
end
