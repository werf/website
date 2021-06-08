# frozen_string_literal: true

module Jekyll
  class FileExistsTag < Liquid::Tag
    def initialize(tag_name, path, tokens)
      super
      @path = path
    end

    def render(context)
      # Pipe parameter through Liquid to make additional replacements possible
      url = Liquid::Template.parse(@path).render context

      # Adds the site source, so that it also works with a custom one
      site_source = context.registers[:site].config['destination']
      file_path = site_source + '/' + url

      # Check if file exists (returns true or false)
      File.exist?(file_path.strip!).to_s
    end
  end
end

Liquid::Template.register_tag('file_exist', Jekyll::FileExistsTag)
