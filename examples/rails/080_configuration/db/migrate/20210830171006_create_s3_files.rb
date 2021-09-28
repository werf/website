class CreateS3Files < ActiveRecord::Migration[6.1]
  def change
    create_table :s3_files do |f|
    end
  end
end
