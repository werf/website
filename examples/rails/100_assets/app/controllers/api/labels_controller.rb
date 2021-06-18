# frozen_string_literal: true

module Api
  class LabelsController < ApiController
    before_action :find_label, only: %i[show update destroy]

    def index
      @labels = Label.all
    end

    def show; end

    def create
      @label = Label.new(label_params)

      if @label.save
        render :show, status: :created
      else
        render_error @label.errors, status: :unprocessable_entity
      end
    end

    def update
      if @label.update(label_params)
        render :show
      else
        render_error @label.errors, status: :unprocessable_entity
      end
    end

    def destroy
      if @label.destroy
        render json: { result: :ok }
      else
        render_error @label.errors, status: :unprocessable_entity
      end
    end

    private

    def find_label
      @label = Label.find(params[:id])
    end

    def label_params
      params.permit(
        :label
      )
    end
  end
end
