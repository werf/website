namespace :crons do
  task cleanup_labels: :environment do
    Label.all.each do |label|
      if Time.now.since(label.created_at) > 3.minute
        puts "Destroying label #{label.label.inspect} by id #{label.id}"
        label.destroy!
      end
    end
  end
end
