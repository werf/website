class CreateTalkers < ActiveRecord::Migration[6.1]
  def change
    create_table :talkers do |t|
      t.text :answer
    end
  end
end
