class TalkerController < ActionController::API
  def say
    begin
      talker = Talker.find(0)
      render plain: talker.answer + ", " + talker.name + "!\n"
    rescue ActiveRecord::RecordNotFound => e
      render plain: "I have nothing to say.\n"
    end
  end

  def remember
    unless params.has_key?(:answer)
      render plain: "You forgot the answer :(\n" and return
    end
    unless params.has_key?(:name)
      render plain: "You forgot the name :(\n" and return
    end

    begin
      talker = Talker.find(0)
    rescue ActiveRecord::RecordNotFound => e
      talker = Talker.new
    end
    talker.update(id: 0, answer: params[:answer], name: params[:name])

    render plain: "Got it.\n"
  end
end
