class AddNameToTalkers < ActiveRecord::Migration[6.1]
  def change
    add_column :talkers, :name, :string
  end
end
