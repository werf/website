# frozen_string_literal: true

Rails.application.routes.draw do
  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html
  root to: 'labels#index'

  namespace :api, defaults: { format: :json } do
    resources :labels
  end
end
