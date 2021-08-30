class S3File < ActiveRecord::Base
  has_one_attached :file
end
