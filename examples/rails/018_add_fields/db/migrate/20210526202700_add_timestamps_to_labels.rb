class AddTimestampsToLabels < ActiveRecord::Migration[6.0]
  def change
    add_column :labels, :created_at, :datetime
    add_column :labels, :updated_at, :datetime
  end
end
