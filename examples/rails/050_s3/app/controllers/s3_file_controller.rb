class S3FileController < ActionController::API
  def download
    begin
      file = S3File.find(0)
    rescue ActiveRecord::RecordNotFound => e
      render plain: "You haven't uploaded anything yet.\n" and return
    end

    send_data file.file.download, filename: file.file.blob.filename.to_s, type: file.file.blob.content_type, disposition: "inline"
  end

  def upload
    unless params.dig(:file).present?
      render plain: "You didn't pass the file to upload :(\n" and return
    end

    begin
      file = S3File.find(0)
    rescue ActiveRecord::RecordNotFound => e
      file = S3File.new
    end
    file.update(id: 0)
    file.file.attach(params[:file])
    file.save()

    render plain: "File uploaded.\n"
  end
end
