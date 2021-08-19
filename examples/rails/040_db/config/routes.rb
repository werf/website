Rails.application.routes.draw do
  get '/remember', to: 'talker#remember'
  get '/say', to: 'talker#say'

  get '/image', to: 'application#image'
  get '/ping', to: 'application#ping'
end
