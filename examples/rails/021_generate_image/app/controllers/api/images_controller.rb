module Api
  class ImagesController < ApplicationController
    def generate_image
      logger.debug "received generate_image request"
      send_data generate_image_from_text params[:text].to_s
    end

    private

    def generate_image_from_text(text)
      logger.debug "start generating image for text #{text.inspect}"
      background = Magick::Image.new(3840, 2160) do
        self.background_color = "yellow"
        self.format = "PNG"
      end

      image = Magick::Image.read("caption:#{text}") do
        self.size = "3840x2160"
        self.background_color = "none"
        self.pointsize = (3840 / text.size).to_i / 0.75
        self.fill = "black"
        self.gravity = Magick::CenterGravity
      end.first

      background.composite!(image, Magick::NorthEastGravity, 0, 0, Magick::OverCompositeOp)

      result = background.to_blob

      logger.debug "finish generation image for text #{text.inspect}"

      result
    end
  end
end
