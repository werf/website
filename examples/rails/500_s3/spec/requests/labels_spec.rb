# frozen_string_literal: true

require 'rails_helper'

RSpec.describe '/api/labels', type: :request do
  let(:valid_attributes) { { label: 'Label' } }
  let(:invalid_attributes) { {} }
  let!(:label) { create :label, valid_attributes }
  let(:last_label) { Label.order(:id).last }

  def json
    JSON.parse(response.body)
  end

  describe 'GET /index' do
    it 'should return a successful response' do
      get api_labels_url
      expect(response).to have_http_status(:ok)
    end

    it 'should return a proper json' do
      get api_labels_url
      expect(json.count).to eq(Label.count)
    end
  end

  describe 'GET /show' do
    it 'renders a successful response' do
      get api_label_url(label)
      expect(response).to have_http_status(:ok)
    end

    it 'should return a proper json' do
      get api_label_url(label)
      expect(json).to eq({ 'id' => label.id, 'label' => label.label })
    end

    context 'with invalid parameters' do
      context 'when entry not exist' do
        it 'renders a forbidden response' do
          get api_label_url(last_label.id + 1)
          expect(response).to have_http_status(:forbidden)
        end

        it 'should return a proper json' do
          get api_label_url(last_label.id + 1)
          expect(json).to eq({ 'result' => 'error', 'comment' => 'entry not found' })
        end
      end
    end
  end

  describe 'POST /create' do
    context 'with valid parameters' do
      it 'should return a successful response' do
        post api_labels_url, params: { label: valid_attributes }
        expect(response).to have_http_status(:created)
      end

      it 'creates a new Label' do
        expect { post api_labels_url, params: { label: valid_attributes } }.to change(Label, :count).by(1)
      end

      it 'should return a proper json' do
        post api_labels_url, params: { label: valid_attributes }
        expect(json).to eq({ 'id' => last_label.id, 'label' => last_label.label })
      end
    end
  end

  describe 'PATCH /update' do
    context 'with valid parameters' do
      let(:new_attributes) { { label: 'New label' } }

      it 'should return a successful response' do
        patch api_label_url(label), params: { label: new_attributes }
        expect(response).to have_http_status(:ok)
      end

      it 'updates the requested label' do
        patch api_label_url(label), params: { label: new_attributes }
        label.reload
        expect(json).to eq({ 'id' => label.id, 'label' => label.label })
      end
    end

    context 'with invalid parameters' do
      context 'when entry not exist' do
        it 'renders a forbidden response' do
          get api_label_url(last_label.id + 1)
          expect(response).to have_http_status(:forbidden)
        end

        it 'should return a proper json' do
          get api_label_url(last_label.id + 1)
          expect(json).to eq({ 'result' => 'error', 'comment' => 'entry not found' })
        end
      end
    end
  end

  describe 'DELETE /destroy' do
    it 'should return a successful response' do
      delete api_label_url(label)
      expect(response).to have_http_status(:ok)
    end

    it 'destroys the requested label' do
      expect { delete api_label_url(label) }.to change(Label, :count).by(-1)
    end

    it 'redirects to the labels list' do
      delete api_label_url(label)
      expect(json).to eq({ 'result' => 'ok' })
    end

    context 'with invalid parameters' do
      context 'when entry not exist' do
        it 'renders a forbidden response' do
          get api_label_url(last_label.id + 1)
          expect(response).to have_http_status(:forbidden)
        end

        it 'should return a proper json' do
          get api_label_url(last_label.id + 1)
          expect(json).to eq({ 'result' => 'error', 'comment' => 'entry not found' })
        end
      end
    end
  end
end
