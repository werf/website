module Jekyll
  module FilesUsed
    class FilesUsedTag < Liquid::Block
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
              { name: "title" }
            ]
          })

          unless @named_params.key?("title")
            @named_params["title"] = "Файлы"
          end

          content = super
          site_config = context.registers[:site].config
          rendered_content = Jekyll::Converters::Markdown::KramdownParser.new(site_config).convert(content)

          result = %Q(
<div class="filesused">
<p><strong>#{@named_params["title"]}</strong></p>
#{rendered_content}
</div>
)
        rescue => e
          Jekyll.logger.abort_with("[filesused] FATAL:", e.message)
        end

        result
      end
    end
  end
end

Liquid::Template.register_tag('filesused', Jekyll::FilesUsed::FilesUsedTag)
