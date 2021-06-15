# frozen_string_literal: true

json.array! @labels, partial: 'api/labels/label', as: :label
