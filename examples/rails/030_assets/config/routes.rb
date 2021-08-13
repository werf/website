Rails.application.routes.draw do
  get '/image', to: 'application#image'
  get '/ping', to: 'application#ping'
end
