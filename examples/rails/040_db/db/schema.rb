# This file is auto-generated from the current state of the database.

ActiveRecord::Schema.define(version: 2021_08_17_164348) do

  create_table "talkers", charset: "utf8", force: :cascade do |t|
    t.text "answer"
    t.string "name"
  end

end
