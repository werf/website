# frozen_string_literal: true

require 'rails_helper'

RSpec.describe Api::LabelsController, type: :routing do
  describe 'routing' do
    it 'routes to #index' do
      expect(get: '/api/labels').to route_to('api/labels#index', format: :json)
    end

    it 'routes to #show' do
      expect(get: '/api/labels/1').to route_to('api/labels#show', id: '1', format: :json)
    end

    it 'routes to #create' do
      expect(post: '/api/labels').to route_to('api/labels#create', format: :json)
    end

    it 'routes to #update via PUT' do
      expect(put: '/api/labels/1').to route_to('api/labels#update', id: '1', format: :json)
    end

    it 'routes to #update via PATCH' do
      expect(patch: '/api/labels/1').to route_to('api/labels#update', id: '1', format: :json)
    end

    it 'routes to #destroy' do
      expect(delete: '/api/labels/1').to route_to('api/labels#destroy', id: '1', format: :json)
    end
  end
end
