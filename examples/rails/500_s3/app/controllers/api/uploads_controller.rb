module Api
  class UploadsController < ApplicationController    

    
    def upload_data
      puts "received upload_data request"
      name = params[:name]
      if name.nil?
        render json: { errors: ["name required!"] }, status: 422
        return
      end

      io = StringIO.new params[:value].to_s

      ActiveStorage::Blob.service.upload name, io
    end

  end
end
