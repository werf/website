require "liquid"

module Jekyll
  module SnippetCut
    class SnippetCutTag < Liquid::Block
      def initialize(tag_name, params_as_string, tokens)
        super
        @params_as_string = params_as_string
      end

      def render(context)
        result = ""

        begin
          @unnamed_params, @named_params = Utils.parse_params(context, @params_as_string)
          Utils.validate_params(@unnamed_params, @named_params, {
            named: [
              { name: "name", required: true },
              { name: "url", required: true },
              { name: "limited", regex: /^(true)|(false)$/ },
            ]
          })

          unless @named_params.key?("limited")
            @named_params["limited"] = false
          end

          content = super
          site_config = context.registers[:site].config
          rendered_content = Jekyll::Converters::Markdown::KramdownParser.new(site_config).convert(content)

          result = %Q(
<div class="snippetcut#{@named_params["limited"] ? ' snippetcut_limited' : ''}" data-snippetcut>
<div class="snippetcut__title">#{if @named_params["url"] != '#' then "<a href=\""+@named_params["url"]+"\" target=\"_blank\" class=\"snippetcut__title-name\" data-snippetcut-name>" else "<span class=\"snippetcut__title-name-text\">" end}#{@named_params["name"]}#{if @named_params["url"] != '#' then "</a>" else "</span>" end}
<a href="javascript:void(0)" class="snippetcut__title-btn" data-snippetcut-btn-name>Copy filename</a>
<a href="javascript:void(0)" class="snippetcut__title-btn" data-snippetcut-btn-text>Copy text</a>
</div>
<div class="highlight">
#{rendered_content}
</div>
<div class="snippetcut__raw" data-snippetcut-text>#{CGI::escapeHTML(remove_excessive_newlines(content.strip.gsub(%r!^```[a-zA-Z0-9]*!,'')))}</div>
</div>
)
        rescue => e
          Jekyll.logger.abort_with("[snippetcut] FATAL:", e.message)
        end

        result
      end

      private

      def remove_excessive_newlines(text)
        text.sub(/^(\s*\R)*/, "").rstrip
      end
    end
  end
end

Liquid::Template.register_tag('snippetcut', Jekyll::SnippetCut::SnippetCutTag)
