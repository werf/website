Rails.application.routes.draw do
  get '/remember', to: 'talker#remember'
  get '/say', to: 'talker#say'

  get '/image', to: 'application#image'
  get '/ping', to: 'application#ping'

  # [<snippet s3>]
  get '/download', to: 's3_file#download'
  post '/upload', to: 's3_file#upload'
  # [<endsnippet s3>]
end
