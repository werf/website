class NotificationsMailer < ApplicationMailer
  default from: "notifications@example.com"

  def labels_count_report_email
    @labels_count = params[:labels_count]
    mail(to: "admin@example.com", subject: "Labels count notification #{Time.now.utc.to_s}")
  end    
end
